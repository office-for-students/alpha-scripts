package data

// NHSNSSQuestions a map of questions in the national student survey (nss) on NHS funded courses
var NHSNSSQuestions = map[int]string{
	1: "I received sufficient preparatory information prior to my placement(s)",
	2: "I was allocated placement(s) suitable for my course",
	3: "I received appropriate supervision on placement(s)",
	4: "I was given opportunities to meet my required practice learning outcomes/competences",
	5: "My contribution during placement(s) as part of a clinical team was valued",
	6: "My practice supervisor(s) understood how my placement(s) related to the broader requirements of my course",
}

// NHSNSS contains the results for the questions on the NSS for students on NHS funded courses
type NHSNSS struct {
	AggregationLevel int       `bson:"aggregation_level,omitempty"` // NHSAGG
	KISMode          string    `bson:"kis_mode"`
	KISCourseID      string    `bson:"kis_course_id"`
	NumberOfStudents int       `bson:"number_of_students,omitempty"` // NHSPOP
	PublicUKPRN      string    `bson:"public_ukprn"`
	ResponseRate     int       `bson:"response_rate"` // NHSRESP_RATE
	Surveys          []*Survey `bson:"survey,omitempty"`
	SubjectCode      string    `bson:"subject_code,omitempty"` // NHSSBJ
	UKPRN            string    `bson:"ukprn"`
	Unavailable      string    `bson:"unavailable,omitempty"` // NHSUNAVAILREASON
}

// Survey contains a result for NSS question
type Survey struct {
	Number                    int    `bson:"question_number,omitempty"`
	ProportionOfStudentsAgree int    `bson:"proportion_of_students_agree_or_strongly_agree,omitempty"`
	Question                  string `bson:"question,omitempty"`
}
