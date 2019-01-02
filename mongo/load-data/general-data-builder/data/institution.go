package data

// Institution contains information of the reporting institution
type Institution struct {
	APROutcome             string `bson:"apr_outcome"`  // APROutcome
	CountryCode            string `bson:"country_code"` // COUNTRY
	PublicUKPRN            string `bson:"public_ukprn"`
	PublicUKPRNCountryCode string `bson:"public_ukprn_country_code"` // PUBUKPRNCOUNTRY
	TEFOutcome             string `bson:"tef_outcome"`               // TEFOutcome
	StudentUnionURL        string `bson:"student_union_url"`         // SUURL
	StudentUnionURLWelsh   string `bson:"student_union_url_welsh"`   // SUURLW
	UKPRN                  string `bson:"ukprn"`
}
