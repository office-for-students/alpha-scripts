package data

// Subject contains JACS level subject codes for each KISCourse
type Subject struct {
	KISMode     string `bson:"kis_mode"`
	KISCourseID string `bson:"kis_course_id"`
	PublicUKPRN string `bson:"public_ukprn"`
	SubjectCode string `bson:"subject_code,omitempty"` // SBJ
	UKPRN       string `bson:"ukprn"`
}
