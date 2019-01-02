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

// CreateEntry ...
func (c *Common) CreateEntry(database, collection, fileName string, counter chan int) error {
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

		entry := &data.Entry{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			SubjectCode: line[7],
			UKPRN:       line[1],
		}

		if line[6] != "" {
			entry.AggregationLevel, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			entry.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[9] != "" {
			entry.ProportionOfStudentsWithALevel, err = strconv.Atoi(line[9])
			if err != nil {
				return err
			}
		}

		if line[8] != "" {
			entry.ProportionOfStudentsWithAccessCourse, err = strconv.Atoi(line[8])
			if err != nil {
				return err
			}
		}

		if line[10] != "" {
			entry.ProportionOfStudentsWithBaccalaureate, err = strconv.Atoi(line[10])
			if err != nil {
				return err
			}
		}

		if line[11] != "" {
			entry.ProportionOfStudentsWithDegree, err = strconv.Atoi(line[11])
			if err != nil {
				return err
			}
		}

		if line[12] != "" {
			entry.ProportionOfStudentsWithFoundation, err = strconv.Atoi(line[12])
			if err != nil {
				return err
			}
		}

		if line[13] != "" {
			entry.ProportionOfStudentsWithNoQuals, err = strconv.Atoi(line[13])
			if err != nil {
				return err
			}
		}

		if line[15] != "" {
			entry.ProportionOfStudentsWithOtherHEQuals, err = strconv.Atoi(line[15])
			if err != nil {
				return err
			}
		}

		if line[14] != "" {
			entry.ProportionOfStudentsWithOtherQuals, err = strconv.Atoi(line[14])
			if err != nil {
				return err
			}
		}

		if line[4] == "1" {
			entry.Unavailable = true
		}

		if err := m.AddEntry(database, collection, entry); err != nil {
			log.ErrorC("failed to add entry resource", err, log.Data{"line_count": count, "entry_resource": entry})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created entry resources", nil)

	return nil
}
