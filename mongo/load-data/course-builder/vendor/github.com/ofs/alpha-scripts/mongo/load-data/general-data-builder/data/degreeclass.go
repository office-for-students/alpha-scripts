package data

// DegreeClass contains information relating to the degree classifications obtained by students
type DegreeClass struct {
	AggregationLevel                           int    `bson:"aggregation_level,omitempty"` // DEGAGG
	KISMode                                    string `bson:"kis_mode"`
	KISCourseID                                string `bson:"kis_course_id"`
	NumberOfStudents                           int    `bson:"number_of_students,omitempty"`                                  // DEGPOP
	ProportionOfStudentsGainDistinction        int    `bson:"proportion_of_students_gaining_distinction,omitempty"`          // UDISTINCTION
	ProportionOfStudentsGainFirstClass         int    `bson:"proportion_of_students_gaining_first_class,omitempty"`          // UFIRST
	ProportionOfStudentsGainLowerSecondClass   int    `bson:"proportion_of_students_gaining_lower_second_class,omitempty"`   // ULOWER
	ProportionOfStudentsGainMerit              int    `bson:"proportion_of_students_gaining_merit,omitempty"`                // UMERIT
	ProportionOfStudentsGainOrdinaryDegree     int    `bson:"proportion_of_students_gaining_ordinary_degree,omitempty"`      // UORDINARY
	ProportionOfStudentsGainOtherHonoursDegree int    `bson:"proportion_of_students_gaining_other_honours_degree,omitempty"` // UOTHER
	ProportionOfStudentsGainPass               int    `bson:"proportion_of_students_gaining_pass,omitempty"`                 // UPASS
	ProportionOfStudentsGainUnclassifiedDegree int    `bson:"proportion_of_students_gaining_unclassified_degree,omitempty"`  // UNA
	ProportionOfStudentsGainUpperSecondClass   int    `bson:"proportion_of_students_gaining_upper_second_class,omitempty"`   // UUPPER
	PublicUKPRN                                string `bson:"public_ukprn"`
	SubjectCode                                string `bson:"subject_code,omitempty"` // DEGSBJ
	UKPRN                                      string `bson:"ukprn"`
	Unavailable                                bool   `bson:"unavailable"` // DEGUNAVAILREASON
}
