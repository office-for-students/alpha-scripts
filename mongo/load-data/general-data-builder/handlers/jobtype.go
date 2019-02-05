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

// CreateJobType ...
func (c *Common) CreateJobType(database, collection, fileName string, counter chan int) error {
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

		jobType := &data.JobType{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			UKPRN:       line[1],
			Unavailable: line[4],
		}

		if line[7] != "" {
			jobType.AggregationLevel, err = strconv.Atoi(line[7])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			jobType.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[10] != "" {
			jobType.ProportionOfStudentsInNonProfessionalOrManagerial, err = strconv.Atoi(line[10])
			if err != nil {
				return err
			}
		}

		if line[9] != "" {
			jobType.ProportionOfStudentsInProfessionalOrManagerial, err = strconv.Atoi(line[9])
			if err != nil {
				return err
			}
		}

		if line[11] != "" {
			jobType.ProportionOfStudentsInUnknownProfessions, err = strconv.Atoi(line[11])
			if err != nil {
				return err
			}
		}

		if line[6] != "" {
			jobType.ResponseRate, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[8] != "" {
			subjectObject, err := m.GetCAHCode("courses", "cah-codes", line[8])
			if err != nil {
				log.ErrorC("failed to find cah code resource", err, log.Data{"line_count": count, "job_type_resource": jobType})
			}

			if subjectObject != nil {
				jobType.SubjectObject = subjectObject
			}
		}

		if err := m.AddJobType(database, collection, jobType); err != nil {
			log.ErrorC("failed to add job type resource", err, log.Data{"line_count": count, "job_type_resource": jobType})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created job type resources", nil)

	return nil
}
