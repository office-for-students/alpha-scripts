package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
)

var (
	mongoURI string

	database   = "courses"
	collection = "courses"
	filename   = "institutions.json"
	mongoSize  = 500
)

func main() {
	flag.StringVar(&mongoURI, "mongo-uri", mongoURI, "mongoDB URI")
	flag.StringVar(&filename, "filename", filename, "filename")
	flag.IntVar(&mongoSize, "size", mongoSize, "size")
	flag.Parse()

	if mongoURI == "" {
		log.Error(errors.New("missing mongo-url flag"), nil)
		os.Exit(1)
	}

	listOfInstitutionNames, err := getInstitutionNames(mongoSize)
	if err != nil {
		log.ErrorC("error creating list of institution names", err, nil)
		os.Exit(1)
	}

	institutionObjects := []InstitutionNameObject{}
	for _, institutionObject := range listOfInstitutionNames {
		institutionObjects = append(institutionObjects, institutionObject)
	}

	// log.Debug("institutions", log.Data{"institutions": institutionObjects})
	// log.Debug("what are the total number of unique institutions?", log.Data{"number_of_institutions": len(institutionObjects)})

	file, err := json.MarshalIndent(institutionObjects, "", "  ")
	if err != nil {
		log.ErrorC("error writing list of institution names to json file", err, log.Data{"filename": filename})
		os.Exit(1)
	}

	if err = ioutil.WriteFile(filename, file, 0644); err != nil {
		log.ErrorC("error writing list of institution names to json file", err, log.Data{"filename": filename})
		os.Exit(1)
	}

	log.Info("completed creation of list of random courses", nil)

}

// InstitutionNameObject represents a document containing information on an institution
// to correctly sort a list of institutions
type InstitutionNameObject struct {
	Alphabet    string `json:"alphabet"`
	Name        string `json:"name"`
	OrderByName string `json:"order_by_name"`
}

func getInstitutionNames(size int) (map[string]InstitutionNameObject, error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return nil, err
	}
	defer session.Close()

	listOfInstitutionNames := make(map[string]InstitutionNameObject)

	it := session.DB(database).C(collection).Find(bson.M{}).Batch(size).Iter()
	count := 0
	for {
		itx := 0
		for ; itx < size; itx++ {
			count++
			course := data.Course{}

			if !it.Next(&course) {
				break
			}

			if course.Institution.UKPRNName == "" {
				continue
			} else if _, ok := listOfInstitutionNames[course.Institution.UKPRNName]; ok {
				continue
			} else {
				listOfInstitutionNames[course.Institution.UKPRNName] = createInstitutionNameObject(course)
			}
		}
		if itx == 0 { // No results read from iterator. Nothing more to do.
			time.Sleep(time.Second * 5)
			break
		}
	}

	// log.Debug("what is the number of courses iterated over?", log.Data{"count": count})

	return listOfInstitutionNames, nil
}

func createInstitutionNameObject(course data.Course) (ino InstitutionNameObject) {
	ino.Name = course.Institution.UKPRNName

	// lowercase institution name
	institutionName := strings.ToLower(ino.Name)

	// Remove unwanted prefixes
	institutionName = strings.TrimPrefix(institutionName, "university of ")
	institutionName = strings.TrimPrefix(institutionName, "the university of ")

	ino.OrderByName = institutionName
	ino.Alphabet = institutionName[:1]

	return
}
