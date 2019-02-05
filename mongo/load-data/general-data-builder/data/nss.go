package data

// NSSQuestions a map of questions in the national student survey (nss)
var NSSQuestions = map[int]string{
	1:  "Staff are good at explaining things",
	2:  "Staff have made the subject interesting",
	3:  "The course is intellectually stimulating",
	4:  "My course has challenged me to achieve my best work",
	5:  "My course has provided me with opportunities to explore ideas or concepts in depth",
	6:  "My course has provided me with opportunities to bring information and ideas together from different topics",
	7:  "My course has provided me with opportunities to apply what I have learnt",
	8:  "The criteria used in marking have been clear in advance",
	9:  "Marking and assessment has been fair",
	10: "Feedback on my work has been timely",
	11: "I have received helpful comments on my work",
	12: "I have been able to contact staff when I needed to",
	13: "I have received sufficient advice and guidance in relation to my course",
	14: "Good advice was available when I needed to make study choices on my course",
	15: "The course is well organised and running smoothly",
	16: "The timetable works efficiently for me",
	17: "Any changes in the course or teaching have been communicated effectively",
	18: "The IT resources and facilities provided have supported my learning well",
	19: "The library resources (e.g. books, online services and learning spaces) have supported my learning well",
	20: "I have been able to access course-specific resources (e.g. equipment, facilities, software, collections) when I needed to",
	21: "I feel part of a community of staff and students",
	22: "I have had the right opportunities to work with other students as part of my course",
	23: "I have had the right opportunities to provide feedback on my course",
	24: "Staff value students' views and opinions about the course",
	25: "It is clear how students' feedback on the course has been acted on",
	26: "The students' union (association or guild) effectively represents students' academic interests",
	27: "Overall, I am satisfied with the quality of the course",
}

// NSS contains the National Student Survey (NSS) results
type NSS struct {
	AggregationLevel int            `bson:"aggregation_level,omitempty"` // NSSAGG
	KISMode          string         `bson:"kis_mode"`
	KISCourseID      string         `bson:"kis_course_id"`
	NumberOfStudents int            `bson:"number_of_students,omitempty"` // NSSPOP
	PublicUKPRN      string         `bson:"public_ukprn"`
	ResponseRate     int            `bson:"response_rate"` // NSSRESP_RATE
	Surveys          []*Survey      `bson:"survey,omitempty"`
	SubjectObject    *SubjectObject `bson:"subject,omitempty"` // NSSSBJ
	UKPRN            string         `bson:"ukprn"`
	Unavailable      string         `bson:"unavailable,omitempty"` // NSSUNAVAILREASON
}
