package data

// InstitutionLocation contains details for each teaching location
type InstitutionLocation struct {
	AccommodationURL      string `bson:"accommodation_url, omitempty"`       // ACCOMURL
	AccommodationURLWelsh string `bson:"accommodation_url_welsh, omitempty"` // ACCOMURLW
	CountryCode           string `bson:"country_code"`                       // LOCCOUNTRY
	Latitude              string `bson:"latitude"`                           // LATITUDE
	LocationID            string `bson:"location_id"`                        // LONGITUDE
	LocationName          string `bson:"location_name, omitempty"`           // LOCNAME
	LocationNameWelsh     string `bson:"location_name_welsh, omitempty"`     // LOCNAMEW
	LocationUKPRN         string `bson:"location_ukprn"`                     // LOCUKPRN
	Longitude             string `bson:"longitude"`                          // LONGITUDE
	StudentUnionURL       string `bson:"student_union_url, omitempty"`       // SUURL
	StudentUnionURLWelsh  string `bson:"student_union_url_welsh, omitempty"` // SUURLW
	UKPRN                 string `bson:"ukprn"`
}
