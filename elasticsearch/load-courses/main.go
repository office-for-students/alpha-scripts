package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-dataset-api/models"
	"github.com/ofs/alpha-scripts/elasticsearch/load-courses/elasticsearch"
)

var (
	esDestURL        = "http://localhost:9200"
	esDestIndex      = "courses"
	esSignedRequests bool

	mongoURL        = "localhost:27017"
	mongoDatabase   = "courses"
	mongoCollection = "courses"
	mongoSize       = 500
)

var (
	wg         sync.WaitGroup
	sem        = make(chan int, 5)
	countCh    = make(chan int)
	insertedCh = make(chan int)
)

func main() {
	flag.StringVar(&mongoURL, "mongo-url", mongoURL, "mongoDB URL")
	flag.StringVar(&mongoDatabase, "mongo-database", mongoDatabase, "mongoDB database")
	flag.StringVar(&mongoCollection, "mongo-collection", mongoCollection, "mongoDB collection")
	flag.StringVar(&esDestURL, "es-dest-url", esDestURL, "elasticsearch destination URL")
	flag.StringVar(&esDestIndex, "es-dest-index", esDestIndex, "elasticsearch index")
	flag.BoolVar(&esSignedRequests, "es-signed-requests", esSignedRequests, "sign elasticsearch requests")
	flag.Parse()

	log.Namespace = "alpha-elasticsearch-loadinator"

	logData := log.Data{
		"mongo-url":          mongoURL,
		"mongo-database":     mongoDatabase,
		"mongo-collection":   mongoCollection,
		"es-dest-url":        esDestURL,
		"es-dest-index":      esDestIndex,
		"es-signed-requests": esSignedRequests,
	}

	ctx := context.Background()

	s, err := mgo.Dial(mongoURL)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "error creating mongoDB session"), logData)
	}

	client := http.DefaultClient
	elasticSearchAPI := elasticsearch.NewElasticSearchAPI(client, esDestURL, esSignedRequests)
	_, testStatus, err := elasticSearchAPI.CallElastic(context.Background(), esDestURL, "GET", nil)
	if err != nil {
		logData["http_status"] = testStatus
		log.ErrorCtx(ctx, errors.WithMessage(err, "unable to connect to elastic search instance"), logData)
		os.Exit(1)
	}

	// Delete index if it already exists
	apiStatus, err := elasticSearchAPI.DeleteSearchIndex(ctx, esDestIndex)
	if err != nil {
		if apiStatus != 404 {
			log.ErrorCtx(ctx, errors.WithMessage(err, "unable to remove index before creating new one"), log.Data{"status": apiStatus, "index": esDestIndex})
			os.Exit(1)
		}
	} else {
		logData["http_status"] = apiStatus
		log.InfoCtx(ctx, "index removed before creating new one", logData)
	}

	// Create new index
	apiStatus, err = elasticSearchAPI.CreateSearchIndex(ctx, esDestIndex)
	if err != nil {
		logData["http_status"] = apiStatus
		log.ErrorCtx(ctx, errors.WithMessage(err, "failure to create index"), logData)
		os.Exit(1)
	}

	// Allow for instance to be created (there can be a lag on the elasticsearch side
	// of things and hence could respond with 200 whilst not finishing the creation of the index)
	time.Sleep(1 * time.Second)

	go status(ctx)

	// Iterate mongo course data
	it := s.DB(mongoDatabase).C(mongoCollection).Find(bson.M{}).Batch(mongoSize).Iter()

	for {
		courses := make([]*models.Course, mongoSize)

		itx := 0
		for ; itx < len(courses); itx++ {
			result := models.Course{}

			if !it.Next(&result) {
				break
			}
			courses[itx] = &result
		}
		if itx == 0 { // No results read from iterator. Nothing more to do.
			time.Sleep(time.Second * 5)
			break
		}

		// This will block if we've reached our concurrecy limit (sem buffer size)
		sendToES(ctx, &courses, itx)
	}

}

func sendToES(ctx context.Context, courses *[]*models.Course, length int) {
	// Wait on semaphore if we've reached our concurrency limit
	wg.Add(1)
	sem <- 1

	go func() {
		defer func() {
			<-sem
			wg.Done()
		}()
		countCh <- length

		var bulk []byte

		i := 0
		for i < length {
			course, courseID := mapResult(ctx, (*courses)[i])
			if course != nil {
				doc := &esDoc{
					Doc: course,
				}
				b, err := json.Marshal(doc)
				if err != nil {
					log.ErrorCtx(ctx, errors.WithMessage(err, "error marshal to json"), nil)
				}

				bulk = append(bulk, []byte("{ \"create\": { \"_index\" : \"courses\", \"_type\" : \"course\", \"_id\": \""+courseID+"\" } }\n")...)
				bulk = append(bulk, b...)
				bulk = append(bulk, []byte("\n")...)
			} else {
				log.ErrorCtx(ctx, errors.New("course empty"), nil)
			}

			i++
		}

		// Load course data into elasticsearch via bulk api
		r, err := http.Post(esDestURL+"/"+esDestIndex+"/_bulk", "application/json", bytes.NewReader(bulk))
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "error posting request"), log.Data{"bulk_json_body": string(bulk)})
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "error reading response body"), nil)
			return
		}

		log.InfoCtx(ctx, "reading response", nil)

		if r.StatusCode > 299 {
			log.ErrorCtx(ctx, errors.New("unexpected post response"), log.Data{"status": r.Status, "body": string(b), "bulk_json_body": string(bulk)})
			return
		}

		log.InfoCtx(ctx, "checked status", nil)

		var bulkRes esBulkResponse
		if err := json.Unmarshal(b, &bulkRes); err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "error unmarshaling json"), nil)
			return
		}

		if bulkRes.Errors {
			for _, r := range bulkRes.Items {
				if r["create"].Status != 201 {
					log.ErrorCtx(ctx, errors.New("error inserting doc"), log.Data{"error": r["create"].Result})
				}
			}
		}

		insertedCh <- length
	}()
}

type esBulkResponse struct {
	Took   int                  `json:"took"`
	Errors bool                 `json:"errors"`
	Items  []esBulkItemResponse `json:"items"`
}

type esBulkItemResponse map[string]esBulkItemResponseData

type esBulkItemResponseData struct {
	ID      string `json:"_id"`
	Index   string `json:"_index"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
	Type    string `json:"_type"`
	Version int    `json:"_version"`
}

type esDoc struct {
	Doc *esCourse `json:"doc"`
}

type esCourse struct {
	KISCourseID          string           `json:"kis_course_id"`
	EnglishTitle         string           `json:"english_title"`
	WelshTitle           string           `json:"welsh_title,omitempty"`
	Country              string           `json:"country"`
	CountryCode          string           `json:"country_code"`
	DistanceLearning     string           `json:"distance_learning"`
	DistanceLearningCode string           `json:"distance_learning_code"`
	FoundationYear       string           `json:"foundation_year"`
	HonoursAward         string           `json:"honours_award"`
	InstitutionName      string           `json:"institution_name"`
	Institution          *esInstitution   `json:"institution"`
	LengthOfCourse       string           `json:"length_of_course"`
	Link                 string           `json:"link"`
	Location             *esLocation      `json:"location"`
	Mode                 string           `json:"mode"`
	NHSFunded            string           `json:"nhs_funded,omitempty"`
	Qualification        *esQualification `json:"qualification"`
	SandwichYear         string           `json:"sandwich_year"`
	SubjectCode          string           `json:"subject_code"`
	SubjectName          string           `json:"subject_name"`
	YearAbroad           string           `json:"year_abroad"`
}

type esInstitution struct {
	PublicUKPRN     string `json:"public_ukprn"`
	PublicUKPRNName string `json:"public_ukprn_name"`
	UKPRN           string `json:"ukprn"`
	UKPRNName       string `json:"ukprn_name"`
	LCUKPRNName     string `json:"lc_ukprn_name"`
}

type esLocation struct {
	EnglishName string `json:"english_name"`
	WelshName   string `json:"welsh_name"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

type esQualification struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Level string `json:"level"`
	Name  string `json:"name"`
}

func mapResult(ctx context.Context, course *models.Course) (*esCourse, string) {
	// Set honours variable
	honours := "Not available"
	if course.Honours {
		honours = "Available"
	}

	foundationYear := "Not available"
	if course.Foundation == "1" {
		foundationYear = "Optional"
	} else if foundationYear == "2" {
		foundationYear = "Compulsory"
	}

	institutionName := removeUniversityOf(course.Institution.UKPRNName)

	esCourse := &esCourse{
		KISCourseID:          course.KISCourseID,
		EnglishTitle:         course.Title.English,
		WelshTitle:           course.Title.Welsh,
		Country:              course.Country.Name,
		CountryCode:          course.Country.Code,
		DistanceLearning:     course.DistanceLearning.Label,
		DistanceLearningCode: course.DistanceLearning.Code,
		FoundationYear:       foundationYear,
		HonoursAward:         honours,
		InstitutionName:      institutionName,
		Institution: &esInstitution{
			PublicUKPRN:     course.Institution.PublicUKPRN,
			PublicUKPRNName: course.Institution.PublicUKPRNName,
			UKPRN:           course.Institution.UKPRN,
			UKPRNName:       strings.Replace(course.Institution.UKPRNName, ",", "", -1),
			LCUKPRNName:     strings.Replace(strings.ToLower(course.Institution.UKPRNName), ",", "", -1),
		},
		LengthOfCourse: course.Length.Code,
		Link:           course.Links.Self,
		Location: &esLocation{
			Latitude:  course.Location.Latitude,
			Longitude: course.Location.Longitude,
		},
		Mode: course.Mode.Label,
		Qualification: &esQualification{
			Code:  course.Qualification.Code,
			Label: course.Qualification.Label,
			Level: course.Qualification.Level,
			Name:  course.Qualification.Name,
		},
		SandwichYear: course.SandwichYear.Label,
		SubjectCode:  course.Subject.Code,
		SubjectName:  course.Subject.Name,
		YearAbroad:   course.YearAbroad.Label,
	}

	if course.Location.Name != nil {
		if course.Location.Name.English != "" {
			esCourse.Location.EnglishName = course.Location.Name.English
		}
		if course.Location.Name.Welsh != "" {
			esCourse.Location.WelshName = course.Location.Name.Welsh
		}
	}

	if course.NHSFunded != nil {
		esCourse.NHSFunded = course.NHSFunded.Label
	}

	courseID := course.Institution.PublicUKPRN + course.KISCourseID + course.Mode.Code

	return esCourse, courseID
}

func removeUniversityOf(institutionName string) string {
	// lowercase institution name
	institutionName = strings.ToLower(institutionName)

	// Remove unwanted prefixes
	institutionName = strings.TrimPrefix(institutionName, "university of ")
	institutionName = strings.TrimPrefix(institutionName, "the university of ")

	// remove unwanted commas
	institutionName = strings.Replace(institutionName, ",", "", -1)

	return institutionName
}

func status(ctx context.Context) {
	var (
		iteratedCounter = 0
		insertedCounter = 0
	)

	t := time.NewTicker(time.Second)

	for {
		select {
		case n := <-countCh:
			iteratedCounter += n
		case n := <-insertedCh:
			insertedCounter += n
		case <-t.C:
			log.InfoCtx(ctx, "Logged:", log.Data{"read": iteratedCounter, "written": insertedCounter})
		}
	}
}
