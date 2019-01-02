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

// CreateNSS ...
func (c *Common) CreateNSS(database, collection, fileName string, counter chan int) error {
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

		nss := &data.NSS{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			SubjectCode: line[8],
			UKPRN:       line[1],
		}

		if line[7] != "" {
			nss.AggregationLevel, err = strconv.Atoi(line[7])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			nss.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[6] != "" {
			nss.ResponseRate, err = strconv.Atoi(line[6])
		}

		var surveys = []*data.Survey{}
		number := 9

		for k, v := range data.NSSQuestions {
			if line[number] != "" {
				proportionOfStudentsAgree, err := strconv.Atoi(line[number])
				if err != nil {
					return err
				}

				survey := &data.Survey{
					Number:                    k,
					ProportionOfStudentsAgree: proportionOfStudentsAgree,
					Question:                  v,
				}

				surveys = append(surveys, survey)
			}
			number++
		}

		nss.Surveys = surveys

		if line[4] == "1" {
			nss.Unavailable = true
		}

		if err := m.AddNSS(database, collection, nss); err != nil {
			log.ErrorC("failed to add nss resource", err, log.Data{"line_count": count, "nss_resource": nss})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created nss resources", nil)

	return nil
}
