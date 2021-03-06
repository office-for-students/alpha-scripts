package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	generalData "github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/data"
	institutionData "github.com/ofs/alpha-scripts/mongo/load-data/institution-builder/data"
	uuid "github.com/satori/go.uuid"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/data"
	"github.com/ofs/alpha-scripts/mongo/load-data/course-builder/statistics"
)

var (
	mongoURI string

	database             = "courses"
	collection           = "courses"
	relativeFileLocation = "../files/"
	courseFileName       = "KISCOURSE"
	fileExtension        = ".csv"
)

func main() {
	flag.StringVar(&mongoURI, "mongo-uri", mongoURI, "mongoDB URI")
	flag.StringVar(&relativeFileLocation, "relative-file-location", relativeFileLocation, "relative location of files")
	flag.Parse()

	if mongoURI == "" {
		log.Error(errors.New("missing mongo-uri flag"), nil)
		os.Exit(1)
	}

	if err := dropCollection(); err != nil {
		os.Exit(1)
	}

	if err := createCourses(courseFileName); err != nil {
		os.Exit(1)
	}

	log.Info("Successfully loaded data", nil)
}

func createCourses(fileName string) error {
	csvFile, err := os.Open(relativeFileLocation + fileName + fileExtension)
	if err != nil {
		log.ErrorC("encountered error immediately when attempting to open file", err, log.Data{"file name": fileName})
		return err
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	// Scan header row (not used)
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

		distance, err := distanceLearningCodeToLabel(line[6])
		if err != nil {
			log.Error(err, log.Data{"func": "distanceLearningCodeToLabel", "line_count": count, "csv_line": line})
			return err
		}

		var honours bool
		if line[10] == "1" {
			honours = true
		}

		length, err := lengthCodeToLabel(line[25])
		if err != nil {
			log.Error(err, log.Data{"func": "lengthCodeToLabel", "line_count": count, "csv_line": line})
			return err
		}

		mode, err := modeCodeToLabel(line[17])
		if err != nil {
			log.Error(err, log.Data{"func": "modeCodeToLabel", "line_count": count, "csv_line": line})
			return err
		}

		sandwichYear, err := availabilityCodeToDescription(line[26])
		if err != nil {
			log.Error(err, log.Data{"func": "sandwich-availabilityCodeToDescription", "line_count": count, "csv_line": line})
			return err
		}

		yearAbroad, err := availabilityCodeToDescription(line[33])
		if err != nil {
			log.Error(err, log.Data{"func": "year-abroad-availabilityCodeToDescription", "line_count": count, "csv_line": line})
			return err
		}

		institution, err := getInstitution("ukprn", line[1])
		if err != nil {
			log.Error(err, log.Data{"func": "getInstitution", "line_count": count, "ukprn": line[1]})
			return err
		}

		publicInstitution, err := getInstitution("public_ukprn", line[0])
		if err != nil {
			log.Error(err, log.Data{"func": "getInstitution", "line_count": count, "public_ukprn": line[0]})
			return err
		}

		var missingLocationID bool
		courseLocationID := ""
		courseLocation, err := getCourseLocation(line[1], line[0], line[16], line[17])
		if err != nil {
			log.Error(err, log.Data{"func": "getLocationIDObject", "line_count": count, "public_ukprn": line[0], "course_id": line[16], "course_mode": line[17]})
			missingLocationID = true
		} else {
			courseLocationID = courseLocation.ID
		}

		// Find teaching location based on courseLocation.ID within locations array inside either publicInstitution or Institution resource
		teachingLocation, err := findTeachingLocation(missingLocationID, institution.Locations, publicInstitution.Locations, courseLocationID)
		if err != nil {
			log.Error(err, log.Data{"func": "findTeachingLocation", "line_count": count, "public_ukprn": line[0], "course_id": line[16], "course_mode": line[17]})
		}

		qualification, err := getQualification(line[34])
		if err != nil {
			log.Error(err, log.Data{"func": "getQualification", "line_count": count, "qualification_code": line[34]})
			return err
		}

		id, err := uuid.NewV4()
		if err != nil {
			log.Error(err, log.Data{"func": "uuid.NewV4", "line_count": count})
			return err
		}

		course := &data.Course{
			ApplicationProvider: line[32],
			Country: &data.Country{
				Code: institution.Country.Code,
				Name: institution.Country.Name,
			},
			DistanceLearning: &data.DistanceLearning{
				Code:  line[6],
				Label: distance,
			},
			Foundation: line[9],
			Honours:    honours,
			ID:         id.String(),
			Institution: &data.InstitutionObject{
				UKPRNName:       institution.Name,
				UKPRN:           line[1],
				PublicUKPRNName: publicInstitution.Name,
				PublicUKPRN:     line[0],
			},
			KISCourseID: line[16],
			Length: &data.LengthObject{
				Code:  line[25],
				Label: length,
			},
			Links: &data.LinkList{
				AssessmentMethod: &data.Language{
					English: line[2],
					Welsh:   line[3],
				},
				CoursePage: &data.Language{
					English: line[4],
					Welsh:   line[5],
				},
				EmploymentDetails: &data.Language{
					English: line[7],
					Welsh:   line[8],
				},
				FinancialSupport: &data.Language{
					English: line[27],
					Welsh:   line[28],
				},
				Institution: "https://localhost:10000/institutions/" + institution.PublicUKPRN,
				LearningAndTeaching: &data.Language{
					English: line[22],
					Welsh:   line[23],
				},
				Self: "https://localhost:10000/institutions/" + institution.PublicUKPRN + "/courses/" + line[16] + "/modes/" + line[17],
			},
			Location: &data.Location{},
			Mode: &data.Mode{
				Code:  line[17],
				Label: mode,
			},
			SandwichYear: &data.Availability{
				Code:  line[26],
				Label: sandwichYear,
			},
			Title: &data.Language{
				English: line[29],
				Welsh:   line[30],
			},
			UCASCode: line[31],
			YearAbroad: &data.Availability{
				Code:  line[33],
				Label: yearAbroad,
			},
		}

		if teachingLocation != nil && teachingLocation.Latitude != "" {
			course.Location.Latitude = teachingLocation.Latitude
			course.Location.Longitude = teachingLocation.Longitude
			course.Location.Name = &data.Language{
				English: teachingLocation.Name.English,
				Welsh:   teachingLocation.Name.Welsh,
			}

			if teachingLocation.Links != nil {
				if teachingLocation.Links.Accommodation != nil {
					course.Links.Accommodation = &data.Language{
						English: teachingLocation.Links.Accommodation.English,
						Welsh:   teachingLocation.Links.Accommodation.Welsh,
					}
				}

				if teachingLocation.Links.StudentUnion != nil {
					course.Links.StudentUnion = &data.Language{
						English: teachingLocation.Links.StudentUnion.English,
						Welsh:   teachingLocation.Links.StudentUnion.Welsh,
					}
				}
			}
		}

		if qualification != nil {
			course.Qualification = &data.Qualification{
				Code:  line[34],
				Label: qualification.Label,
				Level: qualification.Level,
				Name:  qualification.Name,
			}
		} else {
			course.Qualification = manualQualificationLookup(line[34])
		}

		if line[21] != "" {
			courseChange, err := courseChangeCodeToBool(line[21])
			if err != nil {
				log.Error(err, log.Data{"func": "courseChangeCodeToBool", "line_count": count, "csv_line": line})
				return err
			}

			course.Location.Changes = courseChange
		}

		if line[24] != "" {
			nhsFunded, err := nhsCodeToLabel(line[24])
			if err != nil {
				log.Error(err, log.Data{"func": "nhsCodeToLabel", "line_count": count, "csv_line": line})
				return err
			}

			course.NHSFunded = &data.NHSFunded{
				Code:  line[24],
				Label: nhsFunded,
			}
		}

		stats, subject, err := statistics.Get(mongoURI, line[0], line[16], line[17], institution.Country.Code)
		if err != nil {
			log.Error(err, log.Data{"func": "statistics.Get", "line_count": count, "csv_line": line})
			return err
		}

		course.Subject = &data.Subject{
			Code: subject.Code,
			Name: subject.Name,
		}
		course.Statistics = stats

		// Missing title for ucas code 'A16-H09'
		if line[31] == "A16-H09" {
			course.Title.English = "Law"
		}

		if err := addResource(course); err != nil {
			log.ErrorC("failed to add course resource", err, log.Data{"line_count": count, "course_resource": course})
			return err
		}

		count++
		if count%1000 == 0 {
			log.Info(fmt.Sprintf("Progress: %v", count), nil)
		}
	}

	log.Info("Created many course resources", log.Data{"count": count})

	return nil
}

func availabilityCodeToDescription(code string) (description string, err error) {
	switch code {
	case "0":
		description = "Not available"
	case "1":
		description = "Optional"
	case "2":
		description = "Compulsory"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func distanceLearningCodeToLabel(code string) (label string, err error) {
	switch code {
	case "0":
		label = "Course is available other than by distance learning"
	case "1":
		label = "Course is only available through distance learning"
	case "2":
		label = "Course is optionally available through distance learning"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func courseChangeCodeToBool(code string) (changes bool, err error) {
	switch code {
	case "":
	case "0":
		changes = false
	case "1":
		changes = true
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func lengthCodeToLabel(code string) (label string, err error) {
	switch code {
	case "":
	case "1":
		label = "1 stage"
	case "2":
		label = "2 stages"
	case "3":
		label = "3 stages"
	case "4":
		label = "4 stages"
	case "5":
		label = "5 stages"
	case "6":
		label = "6 stages"
	case "7":
		label = "7 stages"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func manualQualificationLookup(code string) (qualification *data.Qualification) {
	switch code {
	case "189":
		qualification = &data.Qualification{
			Code:  code,
			Label: "MRad",
			Level: "U",
			Name:  "Diagnostic Radiography",
		}
	case "190":
		qualification = &data.Qualification{
			Code:  code,
			Label: "MDiet",
			Level: "U",
			Name:  "Dietetics",
		}
	case "191":
		qualification = &data.Qualification{
			Code:  code,
			Label: "MoOth",
			Level: "U",
			Name:  "Occupational Therapy",
		}
	}

	return
}

func modeCodeToLabel(code string) (label string, err error) {
	switch code {
	case "1":
		label = "Full-time"
	case "2":
		label = "Part-time"
	case "3":
		label = "Both"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func nhsCodeToLabel(code string) (label string, err error) {
	switch code {
	case "":
	case "0":
		label = "None"
	case "1":
		label = "Any"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return
}

func addResource(course *data.Course) (err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB(database).C(collection).Insert(course); err != nil {
		log.ErrorC("failed to create course resource", err, nil)
		return
	}

	return
}

func getInstitution(key, value string) (institution *institutionData.Institution, err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("institutions").C("institutions").Find(bson.M{key: value}).One(&institution); err != nil {
		log.ErrorC("failed to find institution resource", err, nil)
	}

	return
}

func getCourseLocation(ukprn, publicUKPRN, kisCourseID, kisMode string) (locationObject *generalData.Location, err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("courses").C("locations").Find(bson.M{"ukprn": ukprn, "public_ukprn": publicUKPRN, "kis_course_id": kisCourseID, "kis_mode": kisMode, "id": bson.M{"$ne": ""}}).One(&locationObject); err != nil {
		log.ErrorC("failed to find course location id resource", err, nil)
	}

	return
}

func getLocation(ukprn, locID string) (teachingLocation *generalData.InstitutionLocation, err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB("institutions").C("locations").Find(bson.M{"ukprn": ukprn, "location_id": locID}).One(&teachingLocation); err != nil {
		log.ErrorC("failed to find teaching location resource", err, nil)
	}

	return
}

func findTeachingLocation(missingLocationID bool, ukprnLocations, publicUKPRNLocations []*institutionData.Location, locationID string) (*institutionData.Location, error) {

	teachingLocation := &institutionData.Location{}
	for _, location := range ukprnLocations {
		if location.ID == locationID {
			return location, nil
		}

		teachingLocation = ukprnLocations[0]
	}

	for _, location := range publicUKPRNLocations {
		if location.ID == locationID {
			return location, nil
		}

		teachingLocation = publicUKPRNLocations[0]
	}

	// fallback on first result in array? (starting with public_ukprn institutions)
	if teachingLocation.ID != "" {
		return teachingLocation, nil
	}

	log.Debug("still no location", log.Data{"ukprn_locations": ukprnLocations, "public_ukprn_locations": publicUKPRNLocations})

	return nil, errors.New("teaching location not found in the possible locations associated with institution")
}

func getQualification(code string) (*generalData.Qualification, error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return nil, err
	}
	defer session.Close()

	var qualification *generalData.Qualification
	if err = session.DB(database).C("qualifications").Find(bson.M{"code": code}).One(&qualification); err != nil {
		log.ErrorC("failed to find qualification resource", err, nil)
	}

	return qualification, nil
}

func dropCollection() (err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if _, err = session.DB(database).C(collection).RemoveAll(nil); err != nil {
		return
	}

	return
}
