package data

// ContinuationRaw represents the continuation statistical data for course (or subject)
type ContinuationRaw struct {
	AggregationLevel             int      `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents             int      `bson:"number_of_students,omitempty"`
	ContinuingWithProvider       int      `bson:"proportion_of_students_continuing_with_provider_after_first_year_on_course,omitempty"`
	Dormant                      int      `bson:"proportion_of_students_dormant_after_first_year_on_course,omitempty"`
	GainingIntendedAwardOrHigher int      `bson:"proportion_of_students_gaining_intended_award_or_higher,omitempty"`
	GainedLowerAward             int      `bson:"proportion_of_students_gained_lower_award,omitempty"`
	LeavingCourse                int      `bson:"proportion_of_students_leaving_course,omitempty"`
	Subject                      *Subject `bson:"subject,omitempty"`
	Unavailable                  string   `bson:"unavailable,omitempty"`
}

// EmploymentRaw represents the employment statistical data for course (or subject)
type EmploymentRaw struct {
	AggregationLevel           int      `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents           int      `bson:"number_of_students,omitempty"`
	AssumedToBeUnemployed      int      `bson:"proportion_of_students_assumed_to_be_unemployed,omitempty"`
	InStudy                    int      `bson:"proportion_of_students_in_study,omitempty"`
	InWork                     int      `bson:"proportion_of_students_in_work,omitempty"`
	InWorkAndStudy             int      `bson:"proportion_of_students_in_work_and_study,omitempty"`
	InWorkOrStudy              int      `bson:"proportion_of_students_in_work_or_study,omitempty"`
	NotAvailableForWorkOrStudy int      `bson:"proportion_of_students_who_are_not_available_for_work_or_study,omitempty"`
	ResponseRate               int      `bson:"response_rate,omitempty"`
	Subject                    *Subject `bson:"subject,omitempty"`
	Unavailable                string   `bson:"unavailable,omitempty"`
}

// JobTypeRaw represents the job type statistical data for course (or subject)
type JobTypeRaw struct {
	AggregationLevel                int      `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents                int      `bson:"number_of_students,omitempty"`
	ProfessionalOrManagerialJobs    int      `bson:"proportion_of_students_in_professional_or_managerial_jobs,omitempty"`
	NonProfessionalOrManagerialJobs int      `bson:"proportion_of_students_in_non_professional_or_managerial_jobs,omitempty"`
	UnknownProfessions              int      `bson:"proportion_of_students_in_unknown_professions,omitempty"`
	ResponseRate                    int      `bson:"response_rate,omitempty"`
	Subject                         *Subject `bson:"subject,omitempty"`
	Unavailable                     string   `bson:"unavailable,omitempty"`
}

// LEORaw represents the LEO statistical data for course (or subject)
type LEORaw struct {
	AggregationLevel    int      `bson:"aggregation_level,omitempty"` // enum
	HigherQuartileRange int      `bson:"higher_quartile_range,omitempty"`
	LowerQuartileRange  int      `bson:"lower_quartile_range,omitempty"`
	Median              int      `bson:"median,omitempty"`
	NumberOfGraduates   int      `bson:"number_of_graduates,omitempty"`
	Subject             *Subject `bson:"subject,omitempty"`
	Unavailable         string   `bson:"unavailable,omitempty"`
}

// SalaryRaw represents the salary statistical data for course (or subject) stored in its raw state
type SalaryRaw struct {
	AggregationLevel                                int      `bson:"aggregation_level,omitempty"`                                     // SALAGG
	InstitutionCourseSalarySixMonthsAfterGraduation *Stats   `bson:"institution_course_salary_six_months_after_graduation,omitempty"` // INST
	KISMode                                         string   `bson:"kis_mode"`
	KISCourseID                                     string   `bson:"kis_course_id"`
	NumberOfStudents                                int      `bson:"number_of_students,omitempty"` // SALPOP
	PublicUKPRN                                     string   `bson:"public_ukprn"`
	ResponseRate                                    int      `bson:"response_rate"`                                       // SALRESP_RATE
	Subject                                         *Subject `bson:"subject,omitempty"`                                   // SALSBJ
	SubjectSalaryFortyMonthsAfterGraduation         *Stats   `bson:"subject_salary_40_months_after_graduation,omitempty"` // LD
	SubjectSalarySixMonthsAfterGraduation           *Stats   `bson:"subject_salary_six_months_after_graduation,omitempty"`
	UKPRN                                           string   `bson:"ukprn"`
	Unavailable                                     string   `bson:"unavailable,omitempty"`
}
