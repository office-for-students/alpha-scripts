package handlers

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// CreateUCASCourseID ...
func (c *Common) CreateUCASCourseID(database, collection, fileName string, counter chan int) error {
	m := c.Mongo

	// Remove data from collection
	if err := m.DropCollection(database, collection); err != nil {
		os.Exit(1)
	}

	csvFile, err := os.Open(c.RelativeFileLocation + fileName + fileExtension)
	if err != nil {
		log.ErrorC("encountered error immediately when attempting to open file", err, log.Data{"file name": fileName})
		return err
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	// Scan header row (not needed)
	_, err = csvReader.Read()
	if err != nil {
		log.ErrorC("encountered error immediately when processing header row", err, nil)
		return err
	}

	count := 0
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.ErrorC("encountered error reading csv", err, log.Data{"line_count": count, "csv_line": line})
			return err
		}

		ucasCourseID := &data.UCASCourseID{
			KISCourseID:  line[2],
			KISMode:      line[3],
			LocationID:   line[4],
			PublicUKPRN:  line[0],
			UKPRN:        line[1],
			UCASCourseID: line[5],
		}

		if err := m.AddUCASCourseID(database, collection, ucasCourseID); err != nil {
			log.ErrorC("failed to add ucas course id resource", err, log.Data{"line_count": count, "ucas_course_id_resource": ucasCourseID})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created ucas course id resources", nil)

	return nil
}
