package data

// Continuation contains continuation information for students on a course
type Continuation struct {
	AggregationLevel                              int            `bson:"aggregation_level,omitempty"` // CONTAGG
	KISMode                                       string         `bson:"kis_mode"`
	KISCourseID                                   string         `bson:"kis_course_id"`
	NumberOfStudents                              int            `bson:"number_of_students,omitempty"`                                                         // CONTPOP
	ProportionOfStudentsContinuing                int            `bson:"proportion_of_students_continuing_with_provider_after_first_year_on_course,omitempty"` // UCONT
	ProportionOfStudentsDormant                   int            `bson:"proportion_of_students_dormant_after_first_year_on_course,omitempty"`                  // UDORMANT
	ProportionOfStudentsGainExpectedOrHigherAward int            `bson:"proportion_of_students_gaining_intended_award_or_higher,omitempty"`                    // UGAINED
	ProportionOfStudentsGainLowerAward            int            `bson:"proportion_of_students_gained_lower_award,omitempty"`                                  // ULOWER
	ProportionOfStudentsLeft                      int            `bson:"proportion_of_students_leaving_course,omitempty"`                                      // ULEFT
	PublicUKPRN                                   string         `bson:"public_ukprn"`
	ResponseRate                                  int            `bson:"response_rate,omitempty"` // COMRESP_RATE
	SubjectObject                                 *SubjectObject `bson:"subject,omitempty"`       // CONTSBJ
	UKPRN                                         string         `bson:"ukprn"`
	Unavailable                                   string         `bson:"unavailable,omitempty"` // CONTUNAVAILREASON
}
