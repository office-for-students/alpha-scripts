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

// CreateDegreeClass ...
func (c *Common) CreateDegreeClass(database, collection, fileName string, counter chan int) error {
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

		degreeClass := &data.DegreeClass{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			SubjectCode: line[7],
			UKPRN:       line[1],
		}

		if line[6] != "" {
			degreeClass.AggregationLevel, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			degreeClass.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[13] != "" {
			degreeClass.ProportionOfStudentsGainDistinction, err = strconv.Atoi(line[13])
			if err != nil {
				return err
			}
		}

		if line[8] != "" {
			degreeClass.ProportionOfStudentsGainFirstClass, err = strconv.Atoi(line[8])
			if err != nil {
				return err
			}
		}

		if line[10] != "" {
			degreeClass.ProportionOfStudentsGainLowerSecondClass, err = strconv.Atoi(line[10])
			if err != nil {
				return err
			}
		}

		if line[14] != "" {
			degreeClass.ProportionOfStudentsGainMerit, err = strconv.Atoi(line[14])
			if err != nil {
				return err
			}
		}

		if line[15] != "" {
			degreeClass.ProportionOfStudentsGainPass, err = strconv.Atoi(line[15])
			if err != nil {
				return err
			}
		}

		if line[12] != "" {
			degreeClass.ProportionOfStudentsGainOrdinaryDegree, err = strconv.Atoi(line[12])
			if err != nil {
				return err
			}
		}

		if line[11] != "" {
			degreeClass.ProportionOfStudentsGainOtherHonoursDegree, err = strconv.Atoi(line[11])
			if err != nil {
				return err
			}
		}

		if line[16] != "" {
			degreeClass.ProportionOfStudentsGainUnclassifiedDegree, err = strconv.Atoi(line[16])
			if err != nil {
				return err
			}
		}

		if line[9] != "" {
			degreeClass.ProportionOfStudentsGainUpperSecondClass, err = strconv.Atoi(line[9])
			if err != nil {
				return err
			}
		}

		if line[4] == "1" {
			degreeClass.Unavailable = true
		}

		if err := m.AddDegreeClass(database, collection, degreeClass); err != nil {
			log.ErrorC("failed to add degree class resource", err, log.Data{"line_count": count, "degree_class_resource": degreeClass})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created degree class resources", nil)

	return nil
}
