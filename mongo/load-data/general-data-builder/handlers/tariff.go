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

// CreateTariff ...
func (c *Common) CreateTariff(database, collection, fileName string, counter chan int) error {
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

		tariff := &data.Tariff{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			UKPRN:       line[1],
			Unavailable: line[4],
		}

		if line[6] != "" {
			tariff.AggregationLevel, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			tariff.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[7] != "" {
			subjectObject, err := m.GetCAHCode("courses", "cah-codes", line[7])
			if err != nil {
				log.ErrorC("failed to find cah code resource", err, log.Data{"line_count": count, "tariff_resource": tariff})
			}

			if subjectObject != nil {
				tariff.SubjectObject = subjectObject
			}
		}

		var tariffs = []*data.TariffStats{}
		number := 8

		for code, description := range data.TariffDescriptions {
			if line[number] != "" {
				proportionOfEntrants, err := strconv.Atoi(line[number])
				if err != nil {
					return err
				}

				tariffStat := &data.TariffStats{
					Code:                 code,
					Description:          description,
					ProportionOfEntrants: proportionOfEntrants,
				}

				tariffs = append(tariffs, tariffStat)
			}
			number++
		}

		tariff.Tariffs = tariffs

		if err := m.AddTariff(database, collection, tariff); err != nil {
			log.ErrorC("failed to add tariff resource", err, log.Data{"line_count": count, "tariff_resource": tariff})
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
