package data

// Employment contains information relating to student employment outcomes
type Employment struct {
	AggregationLevel                               int    `bson:"aggregation_level,omitempty"` // EMPAGG
	KISMode                                        string `bson:"kis_mode"`
	KISCourseID                                    string `bson:"kis_course_id"`
	NumberOfStudents                               int    `bson:"number_of_students,omitempty"`                                             // EMPPOP
	ProportionOfStudentsAssumedToBeUnemployed      int    `bson:"proportion_of_students_assumed_to_be_unemployed,omitempty"`                // ASSUNEMP
	ProportionOfStudentsInStudy                    int    `bson:"proportion_of_students_in_study,omitempty"`                                // STUDY
	ProportionOfStudentsInWork                     int    `bson:"proportion_of_students_in_work,omitempty"`                                 // WORK
	ProportionOfStudentsInWorkAndStudy             int    `bson:"proportion_of_students_in_work_and_study,omitempty"`                       // BOTH
	ProportionOfStudentsInWorkOrStudy              int    `bson:"proportion_of_students_in_work_or_study,omitempty"`                        // WORKSTUDY
	ProportionOfStudentsNotAvailableForWorkOrStudy int    `bson:"proportion_of_students_who_are_not_available_for_work_or_study,omitempty"` // NOAVAIL
	PublicUKPRN                                    string `bson:"public_ukprn"`
	ResponseRate                                   int    `bson:"response_rate,omitempty"` // EMPRESP_RATE
	SubjectCode                                    string `bson:"subject_code,omitempty"`  // EMPSBJ
	UKPRN                                          string `bson:"ukprn"`
	Unavailable                                    bool   `bson:"unavailable"` // EMPUNAVAILREASON
}
