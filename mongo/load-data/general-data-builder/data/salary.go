package data

// Salary contains salary information of students
type Salary struct {
	AggregationLevel                                int            `bson:"aggregation_level,omitempty"`                                     // SALAGG
	InstitutionCourseSalarySixMonthsAfterGraduation *Stats         `bson:"institution_course_salary_six_months_after_graduation,omitempty"` // INST
	KISMode                                         string         `bson:"kis_mode"`
	KISCourseID                                     string         `bson:"kis_course_id"`
	NumberOfStudents                                int            `bson:"number_of_students,omitempty"` // SALPOP
	PublicUKPRN                                     string         `bson:"public_ukprn"`
	ResponseRate                                    int            `bson:"response_rate"`                                       // SALRESP_RATE
	SubjectObject                                   *SubjectObject `bson:"subject,omitempty"`                                   // SALSBJ
	SubjectSalaryFortyMonthsAfterGraduation         *Stats         `bson:"subject_salary_40_months_after_graduation,omitempty"` // LD
	SubjectSalarySixMonthsAfterGraduation           *Stats         `bson:"subject_salary_six_months_after_graduation,omitempty"`
	UKPRN                                           string         `bson:"ukprn"`
	Unavailable                                     string         `bson:"unavailable,omitempty"` // SALUNAVAILREASON
}

// Stats contains a set of values for different statistical measurements of a dataset
type Stats struct {
	LowerQuartile int `bson:"lower_quartile_salary,omitempty"`             // LQ
	Median        int `bson:"median,omitempty"`                            // MED
	UpperQuartile int `bson:"upper_quartile_salary_for_subject,omitempty"` // UQ
}
