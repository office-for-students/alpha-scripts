package data

// Leo represents a course longitudinal education ourcomes
type Leo struct {
	AggregationLevel    int    `bson:"aggregation_level,omitempty"`     // LEOAGG
	HigherQuartileRange int    `bson:"higher_quartile_range,omitempty"` // LEOHQ
	KISMode             string `bson:"kis_mode"`
	KISCourseID         string `bson:"kis_course_id"`
	LowerQuartileRange  int    `bson:"lower_quartile_range,omitempty"` // LEOLQ
	Median              int    `bson:"median,omitempty"`               // LEOMED
	NumberOfGraduates   int    `bson:"number_of_graduates,omitempty"`  // LEOPOP
	PublicUKPRN         string `bson:"public_ukprn"`
	SubjectCode         string `bson:"subject_code,omitempty"` // LEOSBJ
	UKPRN               string `bson:"ukprn"`
	Unavailable         string `bson:"unavailable,omitempty"` // LEOUNAVAILREASON
}
