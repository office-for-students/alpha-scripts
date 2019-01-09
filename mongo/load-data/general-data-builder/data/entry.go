package data

// Entry contains information relating to the entry qualifications of students
type Entry struct {
	AggregationLevel                      int    `bson:"aggregation_level,omitempty"` // ENTAGG
	KISMode                               string `bson:"kis_mode"`
	KISCourseID                           string `bson:"kis_course_id"`
	NumberOfStudents                      int    `bson:"number_of_students,omitempty"`                                                  // ENTPOP
	ProportionOfStudentsWithAccessCourse  int    `bson:"proportion_of_students_with_access_course,omitempty"`                           // ACCESS
	ProportionOfStudentsWithALevel        int    `bson:"proportion_of_students_with_a_level,omitempty"`                                 // ALEVEL
	ProportionOfStudentsWithBaccalaureate int    `bson:"proportion_of_students_with_baccalaureate,omitempty"`                           // BACC
	ProportionOfStudentsWithDegree        int    `bson:"proportion_of_students_with_degree,omitempty"`                                  // DEGREE
	ProportionOfStudentsWithFoundation    int    `bson:"proportion_of_students_with_foundation,omitempty"`                              // FOUNDTN
	ProportionOfStudentsWithNoQuals       int    `bson:"proportion_of_students_with_no_qualifications,omitempty"`                       // NOQUALS
	ProportionOfStudentsWithOtherQuals    int    `bson:"proportion_of_students_with_other_qualifications,omitempty"`                    // OTHER
	ProportionOfStudentsWithOtherHEQuals  int    `bson:"proportion_of_students_with_another_higher_education_qualifications,omitempty"` // OTHERHE
	PublicUKPRN                           string `bson:"public_ukprn"`
	ResponseRate                          int    `bson:"response_rate,omitempty"` // EMPRESP_RATE
	SubjectCode                           string `bson:"subject_code,omitempty"`  // ENTSBJ
	UKPRN                                 string `bson:"ukprn"`
	Unavailable                           string `bson:"unavailable,omitempty"` // ENTUNAVAILREASON
}
