package data

// UCASCourseID contains UCAS course identifiers for each COURSELOCATION
type UCASCourseID struct {
	KISMode      string `bson:"kis_mode"`
	KISCourseID  string `bson:"kis_course_id"`
	LocationID   string `bson:"location_id"`
	PublicUKPRN  string `bson:"public_ukprn"`
	UKPRN        string `bson:"ukprn"`
	UCASCourseID string `bson:"ucas_course_id"`
}
