package data

// Location represents a course location resource
type Location struct {
	ID          string `bson:"id"`
	KISMode     string `bson:"kis_mode"`
	KISCourseID string `bson:"kis_course_id"`
	PublicUKPRN string `bson:"public_ukprn"`
	UKPRN       string `bson:"ukprn"`
}
