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

// CreateNHSNSS ...
func (c *Common) CreateNHSNSS(database, collection, fileName string, counter chan int) error {
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

		nhsNSS := &data.NHSNSS{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			UKPRN:       line[1],
			Unavailable: line[4],
		}

		if line[7] != "" {
			nhsNSS.AggregationLevel, err = strconv.Atoi(line[7])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			nhsNSS.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[6] != "" {
			nhsNSS.ResponseRate, err = strconv.Atoi(line[6])
		}

		if line[8] != "" {
			subjectObject, err := m.GetCAHCode("courses", "cah-codes", line[8])
			if err != nil {
				log.ErrorC("failed to find cah code resource", err, log.Data{"line_count": count, "nhs_nss_resource": nhsNSS})
			}

			if subjectObject != nil {
				nhsNSS.SubjectObject = subjectObject
			}
		}

		var surveys = []*data.Survey{}
		number := 9

		for k, v := range data.NHSNSSQuestions {
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

		nhsNSS.Surveys = surveys

		if err := m.AddNHSNSS(database, collection, nhsNSS); err != nil {
			log.ErrorC("failed to add nhs nss resource", err, log.Data{"line_count": count, "nhs_nss_resource": nhsNSS})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created nhs nss resources", nil)

	return nil
}
