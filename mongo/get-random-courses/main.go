package main

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
)

var (
	mongoURI string

	database   = "courses"
	collection = "courses"
	filename   = "courses.csv"
	mongoSize  = 100
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

	connection, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.ErrorC("error opening file", err, log.Data{"filename": filename})
	}

	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	pipe := session.DB(database).C(collection).Pipe([]bson.M{bson.M{"$sample": bson.M{"size": mongoSize}}})
	it := pipe.Iter()

	writeToFile(connection, filename, "Course,Provider,Location,Length,Study Mode,NHS Funded,Distance Learning,Sandwich Year,Year Abroad,Qualification")

	courses := make([]*data.Course, mongoSize)

	for itx := 0; itx < len(courses); itx++ {
		result := data.Course{}

		if !it.Next(&result) {
			break
		}
		courses[itx] = &result

	}

	for _, course := range courses {
		length := getLength(course.Length.Code)

		nhsFunded := "n/a"
		if course.NHSFunded != nil {
			nhsFunded = course.NHSFunded.Label
		}

		line := course.Qualification.Label + " " + removeCommas(course.Title.English) + "," +
			removeCommas(course.Institution.PublicUKPRNName) + "," +
			removeCommas(course.Location.Name.English) + "," +
			length + "," +
			course.Mode.Label + "," +
			nhsFunded + "," +
			course.DistanceLearning.Label + "," +
			course.SandwichYear.Label + "," +
			course.YearAbroad.Label + "," +
			course.Qualification.Name

		writeToFile(connection, filename, line)
	}

	log.Info("completed creation of list of random courses", nil)

}

func getLength(code string) string {
	plural := " years"
	switch code {
	case "1":
		plural = " year"
	}

	return code + plural
}

func removeCommas(name string) string {
	newName := strings.Replace(name, ",", "", -1)
	return newName
}

func writeToFile(connection *os.File, filename string, line string) {
	_, err := connection.WriteString(line + "\n")
	if err != nil {
		log.ErrorC("error writing line to file", err, log.Data{"line": line, "filename": filename})
	}
}
