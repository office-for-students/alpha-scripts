package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
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
	mongoSize       = 2
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
			course := mapResult(ctx, (*courses)[i])
			if course != nil {
				doc := &esDoc{
					Doc: course,
				}
				b, err := json.Marshal(doc)
				if err != nil {
					log.ErrorCtx(ctx, errors.WithMessage(err, "error marshal to json"), nil)
				}

				bulk = append(bulk, []byte("{ \"create\": { \"_index\" : \"courses\", \"_type\" : \"course\", \"_id\": \""+(*courses)[i].ID+"\" } }\n")...)
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
		}

		if r.StatusCode > 299 {
			log.ErrorCtx(ctx, errors.New("unexpected post response"), log.Data{"status": r.Status, "r": r, "bulk_body": string(bulk)})
		}

		var bulkRes esBulkResponse
		if err := json.Unmarshal(b, &bulkRes); err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "error unmarshaling json"), nil)
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
	KISCourseID      string           `json:"kis_course_id"`
	EnglishTitle     string           `json:"english_title"`
	WelshTitle       string           `json:"welsh_title,omitempty"`
	DistanceLearning string           `json:"distance_learning"`
	FoundationYear   string           `json:"foundation_year"`
	Institution      *esInstitution   `json:"institution"`
	Link             string           `json:"link"`
	Location         *esLocation      `json:"location"`
	Mode             string           `json:"mode"`
	NHSFunded        string           `json:"nhs_funded,omitempty"`
	Qualification    *esQualification `json:"qualification"`
	SandwichYear     string           `json:"sandwich_year"`
	YearAbroad       string           `json:"year_abroad"`
}

type esInstitution struct {
	PublicUKPRN     string `json:"public_ukprn"`
	PublicUKPRNName string `json:"public_ukprn_name"`
	UKPRN           string `json:"ukprn"`
	UKPRNName       string `json:"ukprn_name"`
}

type esLocation struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type esQualification struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Level string `json:"level"`
	Name  string `json:"name"`
}

func mapResult(ctx context.Context, course *models.Course) *esCourse {

	esCourse := &esCourse{
		KISCourseID:      course.KISCourseID,
		EnglishTitle:     course.Title.English,
		WelshTitle:       course.Title.Welsh,
		DistanceLearning: course.DistanceLearning.Label,
		FoundationYear:   course.Foundation,
		Institution: &esInstitution{
			PublicUKPRN:     course.Institution.PublicUKPRN,
			PublicUKPRNName: course.Institution.PublicUKPRNName,
			UKPRN:           course.Institution.UKPRN,
			UKPRNName:       course.Institution.UKPRNName,
		},
		Link: course.Links.Self,
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
		YearAbroad:   course.YearAbroad.Label,
	}

	if course.NHSFunded != nil {
		esCourse.NHSFunded = course.NHSFunded.Label
	}

	return esCourse
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
