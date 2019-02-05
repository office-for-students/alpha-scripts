package data

// JobType contains information relating to the types of profession entered by students
type JobType struct {
	AggregationLevel                                  int            `bson:"aggregation_level,omitempty"` // JOBAGG
	KISMode                                           string         `bson:"kis_mode"`
	KISCourseID                                       string         `bson:"kis_course_id"`
	NumberOfStudents                                  int            `bson:"number_of_students,omitempty"`                                            // JOBPOP
	ProportionOfStudentsInProfessionalOrManagerial    int            `bson:"proportion_of_students_in_professional_or_managerial_jobs,omitempty"`     // PROFMAN
	ProportionOfStudentsInNonProfessionalOrManagerial int            `bson:"proportion_of_students_in_non_professional_or_managerial_jobs,omitempty"` // OTHERJOB
	ProportionOfStudentsInUnknownProfessions          int            `bson:"proportion_of_students_in_unknown_professions,omitempty"`                 // UNKWN
	PublicUKPRN                                       string         `bson:"public_ukprn"`
	ResponseRate                                      int            `bson:"response_rate,omitempty"` // JONRESP_RATE
	SubjectObject                                     *SubjectObject `bson:"subject,omitempty"`       // JOBSBJ
	UKPRN                                             string         `bson:"ukprn"`
	Unavailable                                       string         `bson:"unavailable,omitempty"`
}
