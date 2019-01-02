package handlers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
)

// CreateQualifications ...
func (c *Common) CreateQualifications(database, collection, fileName string, counter chan int) error {
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

		code, err := getCode(line[0])
		if err != nil {
			log.Error(err, log.Data{"line_count": count, "csv_line": line})
			return err
		}

		qualification := &data.Qualification{
			Code:  code,
			Label: line[1],
			Level: line[2],
			Name:  line[3],
		}

		if err := m.AddQualification(database, collection, qualification); err != nil {
			log.ErrorC("failed to add qualification resource", err, log.Data{"line_count": count, "qualification_resource": qualification})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("Created qualification resources", nil)

	return nil
}

func getCode(code string) (newCode string, err error) {
	switch len(code) {
	case 1:
		newCode = "00" + code
	case 2:
		newCode = "0" + code
	case 3:
		newCode = code
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}
