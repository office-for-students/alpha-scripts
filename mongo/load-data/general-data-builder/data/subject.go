package data

// Subject contains JACS level subject codes for each KISCourse
type Subject struct {
	KISMode       string         `bson:"kis_mode"`
	KISCourseID   string         `bson:"kis_course_id"`
	PublicUKPRN   string         `bson:"public_ukprn"`
	SubjectObject *SubjectObject `bson:"subject,omitempty"` // SBJ
	UKPRN         string         `bson:"ukprn"`
}

// SubjectObject contains relation between subject code and name for each KISCourse
type SubjectObject struct {
	SubjectCode string `bson:"code,omitempty"`
	SubjectName string `bson:"name,omitempty"`
}
