package handlers

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// CreateContinuation ...
func (c *Common) CreateContinuation(database, collection, fileName string, counter chan int) error {
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

		continuation := &data.Continuation{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			SubjectCode: line[7],
			UKPRN:       line[1],
		}

		if line[6] != "" {
			continuation.AggregationLevel, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			continuation.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[8] != "" {
			continuation.ProportionOfStudentsContinuing, err = strconv.Atoi(line[8])
			if err != nil {
				return err
			}
		}

		if line[9] != "" {
			continuation.ProportionOfStudentsDormant, err = strconv.Atoi(line[9])
			if err != nil {
				return err
			}
		}

		if line[10] != "" {
			continuation.ProportionOfStudentsGainExpectedOrHigherAward, err = strconv.Atoi(line[10])
			if err != nil {
				return err
			}
		}

		if line[12] != "" {
			continuation.ProportionOfStudentsGainLowerAward, err = strconv.Atoi(line[12])
			if err != nil {
				return err
			}
		}

		if line[11] != "" {
			continuation.ProportionOfStudentsLeft, err = strconv.Atoi(line[11])
			if err != nil {
				return err
			}
		}

		if line[4] == "1" {
			continuation.Unavailable = true
		}

		if err := m.AddContinuation(database, collection, continuation); err != nil {
			log.ErrorC("failed to add continuation resource", err, log.Data{"line_count": count, "continuation_resource": continuation})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created continuation resources", nil)

	return nil
}
