package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	handlers "github.com/ofs/alpha-scripts/mongo/get-random-courses/handlers"
	"github.com/ofs/alpha-scripts/mongo/load-data/institution-builder/data"
)

var (
	authToken    string
	authPassword = "anything"

	mongoURI string

	database             = "institutions"
	collection           = "institutions"
	relativeFileLocation = "../files/"
	ukprnLookupFileName  = "UNISTATS_UKPRN_lookup_20160901"
	institutionFileName  = "INSTITUTION"
	locationFileName     = "LOCATION"
	fileExtension        = ".csv"

	institutionURL = "https://data.unistats.ac.uk/api/v4/KIS/Institution/"
)

func main() {
	flag.StringVar(&authToken, "auth-token", authToken, "authentication token or username")
	flag.StringVar(&authPassword, "auth-password", authPassword, "authentication password")
	flag.StringVar(&mongoURI, "mongo-uri", mongoURI, "mongoDB URI")
	flag.StringVar(&relativeFileLocation, "relative-file-location", relativeFileLocation, "relative location of files")
	flag.Parse()

	if mongoURI == "" {
		log.Error(errors.New("missing mongo-uri flag"), nil)
		os.Exit(1)
	}

	if authToken == "" {
		log.Error(errors.New("missing auth-header flag"), nil)
		os.Exit(1)
	}

	// Remove data from institution collection
	if err := dropCollection(); err != nil {
		os.Exit(1)
	}

	if err := createInstitutions(authToken, authPassword, ukprnLookupFileName); err != nil {
		os.Exit(1)
	}

	if err := updateInstitutions(institutionFileName); err != nil {
		os.Exit(1)
	}

	if err := updateLocations(locationFileName); err != nil {
		os.Exit(1)
	}

	log.Info("Successfully loaded institution data", nil)
}

func createInstitutions(authToken, authPassword, fileName string) error {
	csvFile, err := os.Open(relativeFileLocation + fileName + fileExtension)
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

	var client *http.Client

	client = http.DefaultClient

	request := handlers.Request{
		Authorization: &handlers.Authorization{
			Username: authToken,
			Password: authPassword,
		},
		Client: client,
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

		// Should get the institution names from ukrlp API and not unistats API; will need to apply for a free account?
		institutionName, _ := request.GetInstitutionName(institutionURL + line[0])

		if institutionName == "" {
			institutionName = strings.Replace(line[1], "/", ",", -1)
		}

		institution := &data.Institution{
			Country:     &data.Country{},
			Links:       &data.LinkList{},
			Name:        institutionName,
			PublicUKPRN: line[0],
		}

		if err := addResource(institution); err != nil {
			log.ErrorC("failed to add institution resource", err, log.Data{"line_count": count, "institution_resource": institution})
			return err
		}

		count++
	}

	log.Info("Created institution resources", log.Data{"count": count})

	return nil
}

func updateInstitutions(institutionFile string) error {
	csvFile, err := os.Open(relativeFileLocation + institutionFile + fileExtension)
	if err != nil {
		log.ErrorC("encountered error immediately when attempting to open file", err, log.Data{"file name": institutionFile})
		return err
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	// Scan header row (not needed)
	_, err = csvReader.Read()
	if err != nil {
		log.ErrorC("encountered error immediately when processing header row", err, nil)
		return err
	}

	// Possibly validate headers but for now we know the structure of file so continue

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

		country, err := countryCodeToName(line[2])
		if err != nil {
			log.Error(err, log.Data{"line_count": count, "csv_line": line})
			return err
		}

		publicUKPRN := line[0]

		institution := &data.Institution{
			APROutcome: line[5],
			Country: &data.Country{
				Code: line[2],
				Name: country,
			},
			Links: &data.LinkList{
				Courses: "https://localhost:10000/institutions/" + line[1] + "/courses",
				InstitutionStudentUnion: &data.Language{
					English: line[6],
					Welsh:   line[7],
				},
				Self: "https://localhost:10000/institutions/" + line[1],
			},
			TEFOutcome: line[4],
			UKPRN:      line[1],
		}

		// Manually add institution name for ukprn 10008173
		if publicUKPRN == "10008173" {
			var locations []*data.Location
			location := &data.Location{
				Latitude:  "51.453256", // latitude and longitude taken from google
				Longitude: "-0.963443",
				Name: &data.Language{
					English: institution.Name,
				},
			}
			locations = append(locations, location)

			institution.Locations = locations
			institution.Name = "University College of Estate Management"
			institution.PublicUKPRN = publicUKPRN
		}

		if err := upsertResource(publicUKPRN, institution); err != nil {
			log.ErrorC("failed to update institution resource", err, log.Data{"line_count": count, "institution_resource": institution, "line": line})
			return err
		}

		count++
	}

	log.Info("Updated many institution resources", log.Data{"count": count})

	return nil
}

func updateLocations(locationFile string) error {
	csvFile, err := os.Open(relativeFileLocation + locationFile + fileExtension)
	if err != nil {
		log.ErrorC("encountered error immediately when attempting to open file", err, log.Data{"file name": locationFile})
		return err
	}
	csvReader := csv.NewReader(bufio.NewReader(csvFile))

	// Scan header row (not needed)
	_, err = csvReader.Read()
	if err != nil {
		log.ErrorC("encountered error immediately when processing header row", err, nil)
		return err
	}

	// Possibly validate headers but for now we know the structure of file so continue

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

		publicUKPRN := line[8]
		ukprn := line[0]

		location := &data.Location{
			ID: line[3],
			Links: &data.LocationLinks{
				Accommodation: &data.Language{
					English: line[1],
					Welsh:   line[2],
				},
				StudentUnion: &data.Language{
					English: line[10],
					Welsh:   line[11],
				},
			},
			Latitude:  line[6],
			Longitude: line[7],
			Name: &data.Language{
				English: line[4],
				Welsh:   line[5],
			},
		}

		if publicUKPRN == "" {
			publicUKPRN = ukprn
		}

		if err := insertLocation(publicUKPRN, location); err != nil {
			log.ErrorC("failed to update institution resource with location data", err, log.Data{"line_count": count, "location_resource": location, "line": line})
			return err
		}

		// Some documents are missing resource name, use the english or welsh name found in location file
		var name string
		if line[4] != "" {
			name = line[4]
		} else {
			name = line[5]
		}

		if name != "" {
			_ = addName(publicUKPRN, name)
		}

		count++
	}

	log.Info("Updated many institution resources with location data", log.Data{"count": count})

	return nil
}

func countryCodeToName(code string) (name string, err error) {
	switch code {
	case "XF":
		name = "England"
	case "XG":
		name = "Northern Ireland"
	case "XH":
		name = "Scotland"
	case "XI":
		name = "Wales"
	default:
		err = fmt.Errorf("Unknown code: [%s]", code)
	}

	return name, nil
}

func addResource(institution *data.Institution) (err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB(database).C(collection).Insert(institution); err != nil {
		log.ErrorC("failed to create institution resource", err, nil)
		return
	}

	return
}

func addName(publicUKPRN, name string) (err error) {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return
	}
	defer session.Close()

	if err = session.DB(database).C(collection).Update(bson.M{"public_ukprn": publicUKPRN, "name": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"name": name}}); err != nil {
		if err != mgo.ErrNotFound {
			log.ErrorC("failed to update institution resource with name", err, nil)
		}
	}

	return
}

func insertLocation(publicUKPRN string, location *data.Location) error {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return err
	}
	defer session.Close()

	selector := bson.M{"public_ukprn": publicUKPRN}

	query := createLocationUpdateQuery(location)
	log.Info("check my location data", log.Data{"location": location, "selector": selector, "query": query})
	if err = session.DB(database).C(collection).Update(selector, query); err != nil {
		log.ErrorC("failed to upsert institution resource", err, nil)
	}

	return nil
}

func upsertResource(publicUKPRN string, institution *data.Institution) error {
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
		return err
	}
	defer session.Close()

	selector := bson.M{"public_ukprn": publicUKPRN}

	update := createInstitutionUpdateQuery(institution)
	if _, err = session.DB(database).C(collection).Upsert(selector, bson.M{"$set": update}); err != nil {
		log.ErrorC("failed to upsert institution resource", err, nil)
	}

	return nil
}

func createInstitutionUpdateQuery(institution *data.Institution) bson.M {
	setUpdates := make(bson.M)

	if institution.APROutcome != "" {
		setUpdates["apr_outcome"] = institution.APROutcome
	}

	if institution.Country != nil {
		if institution.Country.Code != "" {
			setUpdates["country.code"] = institution.Country.Code
		}

		if institution.Country.Name != "" {
			setUpdates["country.name"] = institution.Country.Name
		}
	}

	if institution.Links != nil {
		if institution.Links.Courses != "" {
			setUpdates["links.courses"] = institution.Links.Courses
		}

		if institution.Links.InstitutionStudentUnion != nil {
			if institution.Links.InstitutionStudentUnion.English != "" {
				setUpdates["links.institution_student_union.english"] = institution.Links.InstitutionStudentUnion.English
			}

			if institution.Links.InstitutionStudentUnion.Welsh != "" {
				setUpdates["links.institution_student_union.welsh"] = institution.Links.InstitutionStudentUnion.Welsh
			}
		}

		if institution.Links.Self != "" {
			setUpdates["links.self"] = institution.Links.Self
		}
	}

	if institution.Locations != nil && len(institution.Locations) > 0 {
		setUpdates["locations"] = institution.Locations
	}

	if institution.Name != "" {
		setUpdates["name"] = institution.Name
	}

	if institution.PublicUKPRN != "" {
		setUpdates["public_ukprn"] = institution.PublicUKPRN
	}

	if institution.UKPRN != "" {
		setUpdates["ukprn"] = institution.UKPRN
	}

	if institution.TEFOutcome != "" {
		setUpdates["tef_outcome"] = institution.TEFOutcome
	}

	return setUpdates
}

func createLocationUpdateQuery(location *data.Location) (query bson.M) {
	setUpdates := make(bson.M)

	if location.ID != "" {
		setUpdates["id"] = location.ID
	}

	if location.Latitude != "" {
		setUpdates["latitude"] = location.Latitude
	}

	if location.Longitude != "" {
		setUpdates["longitude"] = location.Longitude
	}

	if location.Name != nil {
		setName := make(bson.M)
		if location.Name.English != "" {
			setName["english"] = location.Name.English
		}

		if location.Name.Welsh != "" {
			setName["welsh"] = location.Name.Welsh
		}

		setUpdates["name"] = setName
	}

	if location.Links != nil {
		setLinks := make(bson.M)
		if location.Links.Accommodation != nil {
			setAccommodation := make(bson.M)
			if location.Links.Accommodation.English != "" {
				setAccommodation["english"] = location.Links.Accommodation.English
			}

			if location.Links.Accommodation.Welsh != "" {
				setAccommodation["welsh"] = location.Links.Accommodation.Welsh
			}

			setLinks["accommodation"] = setAccommodation
		}

		if location.Links.StudentUnion != nil {
			setStudentUnion := make(bson.M)
			if location.Links.StudentUnion.English != "" {
				setStudentUnion["english"] = location.Links.StudentUnion.English
			}

			if location.Links.StudentUnion.Welsh != "" {
				setStudentUnion["welsh"] = location.Links.StudentUnion.Welsh
			}

			setLinks["student_union"] = setStudentUnion
		}

		setUpdates["links"] = setLinks
	}

	query = bson.M{"$addToSet": bson.M{"locations": setUpdates}}

	return
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
