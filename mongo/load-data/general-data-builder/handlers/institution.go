package handlers

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// CreateInstitution ...
func (c *Common) CreateInstitution(database, collection, fileName string, counter chan int) error {
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

		institution := &data.Institution{
			APROutcome:             line[5],
			CountryCode:            line[2],
			PublicUKPRNCountryCode: line[3],
			PublicUKPRN:            line[0],
			StudentUnionURL:        line[6],
			StudentUnionURLWelsh:   line[7],
			TEFOutcome:             line[4],
			UKPRN:                  line[1],
		}

		if err := m.AddInstitution(database, collection, institution); err != nil {
			log.ErrorC("failed to add raw institution resource", err, log.Data{"line_count": count, "institution_resource": institution})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created raw institution resources", nil)

	return nil
}
