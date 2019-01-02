package handlers

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// CreateInstitutionLocation ...
func (c *Common) CreateInstitutionLocation(database, collection, fileName string, counter chan int) error {
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

		institutionLocation := &data.InstitutionLocation{
			AccommodationURL:      line[1],
			AccommodationURLWelsh: line[2],
			CountryCode:           line[9],
			Latitude:              line[6],
			LocationID:            line[3],
			LocationName:          line[4],
			LocationNameWelsh:     line[5],
			LocationUKPRN:         line[8],
			Longitude:             line[7],
			StudentUnionURL:       line[10],
			StudentUnionURLWelsh:  line[11],
			UKPRN:                 line[0],
		}

		if err := m.AddInstitutionLocation(database, collection, institutionLocation); err != nil {
			log.ErrorC("failed to add institution location resource", err, log.Data{"line_count": count, "institution_location_resource": institutionLocation})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created institution location resources", nil)

	return nil
}
