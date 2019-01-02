package data

// JobList contains information about common job types obtained by students
type JobList struct {
	Job                  string `bson:"job,omitempty"` // JOB
	KISMode              string `bson:"kis_mode"`
	KISCourseID          string `bson:"kis_course_id"`
	Order                int    `bson:"order,omitempty"`                  // ORDER
	PercentageOfStudents int    `bson:"percentage_of_students,omitempty"` // PERC
	PublicUKPRN          string `bson:"public_ukprn"`
	SubjectCode          string `bson:"subject_code,omitempty"` // COMSBJ
	UKPRN                string `bson:"ukprn"`
}
