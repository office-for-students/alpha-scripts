package data

// CommonData contains information relating to common job types obtained by students taking a course
type CommonData struct {
	AggregationLevel int    `bson:"aggregation_level,omitempty"` // COMAGG
	KISMode          string `bson:"kis_mode"`
	KISCourseID      string `bson:"kis_course_id"`
	NumberOfStudents int    `bson:"number_of_students,omitempty"` // COMPOP
	PublicUKPRN      string `bson:"public_ukprn"`
	ResponseRate     int    `bson:"response_rate,omitempty"` // COMRESP_RATE
	SubjectCode      string `bson:"subject_code,omitempty"`  // COMSBJ
	UKPRN            string `bson:"ukprn"`
	Unavailable      bool   `bson:"unavailable"` // COMUNAVAILREASON
}
