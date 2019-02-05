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

// CreateSalary ...
func (c *Common) CreateSalary(database, collection, fileName string, counter chan int) error {
	m := c.Mongo

	// Remove data from institution collection
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

		salary := &data.Salary{
			KISCourseID: line[2],
			KISMode:     line[3],
			PublicUKPRN: line[0],
			UKPRN:       line[1],
			Unavailable: line[4],
		}

		if line[7] != "" {
			salary.AggregationLevel, err = strconv.Atoi(line[7])
			if err != nil {
				return err
			}
		}

		if line[9] != "" {
			salary.SubjectSalaryFortyMonthsAfterGraduation = &data.Stats{}
			salary.SubjectSalaryFortyMonthsAfterGraduation.LowerQuartile, err = strconv.Atoi(line[9])
			if err != nil {
				return err
			}
		}

		if line[10] != "" {
			salary.SubjectSalaryFortyMonthsAfterGraduation.Median, err = strconv.Atoi(line[10])
			if err != nil {
				return err
			}
		}

		if line[11] != "" {
			salary.SubjectSalaryFortyMonthsAfterGraduation.UpperQuartile, err = strconv.Atoi(line[11])
			if err != nil {
				return err
			}
		}

		if line[12] != "" {
			salary.SubjectSalarySixMonthsAfterGraduation = &data.Stats{}
			salary.SubjectSalarySixMonthsAfterGraduation.LowerQuartile, err = strconv.Atoi(line[12])
			if err != nil {
				return err
			}
		}

		if line[13] != "" {
			salary.SubjectSalarySixMonthsAfterGraduation.Median, err = strconv.Atoi(line[13])
			if err != nil {
				return err
			}
		}

		if line[14] != "" {
			salary.SubjectSalarySixMonthsAfterGraduation.UpperQuartile, err = strconv.Atoi(line[14])
			if err != nil {
				return err
			}
		}

		if line[15] != "" {
			salary.InstitutionCourseSalarySixMonthsAfterGraduation = &data.Stats{}
			salary.InstitutionCourseSalarySixMonthsAfterGraduation.LowerQuartile, err = strconv.Atoi(line[15])
			if err != nil {
				return err
			}
		}

		if line[16] != "" {
			salary.InstitutionCourseSalarySixMonthsAfterGraduation.Median, err = strconv.Atoi(line[16])
			if err != nil {
				return err
			}
		}

		if line[17] != "" {
			salary.InstitutionCourseSalarySixMonthsAfterGraduation.UpperQuartile, err = strconv.Atoi(line[17])
			if err != nil {
				return err
			}
		}

		if line[5] != "" {
			salary.NumberOfStudents, err = strconv.Atoi(line[5])
			if err != nil {
				return err
			}
		}

		if line[6] != "" {
			salary.ResponseRate, err = strconv.Atoi(line[6])
			if err != nil {
				return err
			}
		}

		if line[8] != "" {
			subjectObject, err := m.GetCAHCode("courses", "cah-codes", line[8])
			if err != nil {
				log.ErrorC("failed to find cah code resource", err, log.Data{"line_count": count, "salary_resource": salary})
			}

			if subjectObject != nil {
				salary.SubjectObject = subjectObject
			}
		}

		if err := m.AddSalary(database, collection, salary); err != nil {
			log.ErrorC("failed to add salary resource", err, log.Data{"line_count": count, "salary_resource": salary})
			return err
		}

		count++
		if count%1000 == 0 {
			counter <- count
			count = 0
		}
	}

	counter <- count
	log.Info("created salary resources", nil)

	return nil
}
