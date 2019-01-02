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

// CreateJobList ...
func (c *Common) CreateJobList(database, collection, fileName string, counter chan int) error {
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

		jobList := &data.JobList{
			Job:         line[5],
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			SubjectCode: line[4],
			UKPRN:       line[1],
		}

		if line[6] != "" {
			jobList.PercentageOfStudents, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[7] != "" {
			jobList.Order, err = strconv.Atoi(line[7])
			if err != nil {
				return err
			}
		}

		if err := m.AddJobList(database, collection, jobList); err != nil {
			log.ErrorC("failed to add job list resource", err, log.Data{"line_count": count, "job_list_resource": jobList})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created job list resources", nil)

	return nil
}
