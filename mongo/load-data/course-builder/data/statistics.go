package data

// Statistics represents an object containing a list of statistical data for course (or subject)
type Statistics struct {
	Continuation []*Continuation `bson:"continuation,omitempty"`
	Employment   []*Employment   `bson:"employment,omitempty"`
	JobList      *JobList        `bson:"job_list,omitempty"`
	JobType      []*JobType      `bson:"job_type,omitempty"`
	LEO          []*LEO          `bson:"leo,omitempty"`
	Salary       []*Salary       `bson:"salary,omitempty"`
}

// Common represents the metadata relative to the job list statistical data for course (or subject)
type Common struct {
	AggregationLevel int    `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents int    `bson:"number_of_students,omitempty"`
	ResponseRate     int    `bson:"response_rate,omitempty"`
	SubjectCode      string `bson:"subject_code,omitempty"`
	Unavailable      string `bson:"unavailable,omitempty"`
}

// Continuation represents the continuation statistical data for course (or subject)
type Continuation struct {
	AggregationLevel             int    `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents             int    `bson:"number_of_students,omitempty"`
	ContinuingWithProvider       int    `bson:"proportion_of_students_continuing_with_provider_after_first_year_on_course,omitempty"`
	Dormant                      int    `bson:"proportion_of_students_dormant_after_first_year_on_course,omitempty"`
	GainingIntendedAwardOrHigher int    `bson:"proportion_of_students_gaining_intended_award_or_higher,omitempty"`
	GainedLowerAward             int    `bson:"proportion_of_students_gained_lower_award,omitempty"`
	LeavingCourse                int    `bson:"proportion_of_students_leaving_course,omitempty"`
	SubjectCode                  string `bson:"subject_code,omitempty"`
	Unavailable                  string `bson:"unavailable,omitempty"`
}

// Employment represents the employment statistical data for course (or subject)
type Employment struct {
	AggregationLevel           int    `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents           int    `bson:"number_of_students,omitempty"`
	AssumedToBeUnemployed      int    `bson:"proportion_of_students_assumed_to_be_unemployed,omitempty"`
	InStudy                    int    `bson:"proportion_of_students_in_study,omitempty"`
	InWork                     int    `bson:"proportion_of_students_in_work,omitempty"`
	InWorkAndStudy             int    `bson:"proportion_of_students_in_work_and_study,omitempty"`
	InWorkOrStudy              int    `bson:"proportion_of_students_in_work_or_study,omitempty"`
	NotAvailableForWorkOrStudy int    `bson:"proportion_of_students_who_are_not_available_for_work_or_study,omitempty"`
	ResponseRate               int    `bson:"response_rate,omitempty"`
	SubjectCode                string `bson:"subject_code,omitempty"`
	Unavailable                string `bson:"unavailable,omitempty"`
}

// JobList represents the job list statistical data for course
type JobList struct {
	AggregationLevel int       `bson:"aggregation_level,omitempty"` // enum
	Items            []JobItem `bson:"items,omitempty"`
	NumberOfStudents int       `bson:"number_of_students,omitempty"`
	ResponseRate     int       `bson:"response_rate,omitempty"`
	SubjectCode      string    `bson:"subject_code,omitempty"`
	Unavailable      string    `bson:"unavailable,omitempty"`
}

// JobItem represents a single item within a job list
type JobItem struct {
	List  []Job `bson:"list,omitempty"`
	Order int   `bson:"order,omitempty"`
}

// Job represents statistical data of the number of students in a job after taking course (or subject)
type Job struct {
	Job                  string `bson:"job"`
	PercentageOfStudents int    `bson:"percentage_of_students"`
}

// JobType represents the job type statistical data for course (or subject)
type JobType struct {
	AggregationLevel                int    `bson:"aggregation_level,omitempty"` // enum
	NumberOfStudents                int    `bson:"number_of_students,omitempty"`
	ProfessionalOrManagerialJobs    int    `bson:"proportion_of_students_in_professional_or_managerial_jobs,omitempty"`
	NonProfessionalOrManagerialJobs int    `bson:"proportion_of_students_in_non_professional_or_managerial_jobs,omitempty"`
	UnknownProfessions              int    `bson:"proportion_of_students_in_unknown_professions,omitempty"`
	ResponseRate                    int    `bson:"response_rate,omitempty"`
	SubjectCode                     string `bson:"subject_code,omitempty"`
	Unavailable                     string `bson:"unavailable,omitempty"`
}

// JobOrder represents statistical data of the number of students in a job after taking course (or subject)
type JobOrder struct {
	Order                int    `bson:"order"`
	Job                  string `bson:"job"`
	PercentageOfStudents int    `bson:"percentage_of_students"`
}

// LEO represents the LEO statistical data for course (or subject)
type LEO struct {
	AggregationLevel    int    `bson:"aggregation_level,omitempty"` // enum
	HigherQuartileRange int    `bson:"higher_quartile_range,omitempty"`
	LowerQuartileRange  int    `bson:"lower_quartile_range,omitempty"`
	Median              int    `bson:"median,omitempty"`
	NumberOfGraduates   int    `bson:"number_of_graduates,omitempty"`
	SubjectCode         string `bson:"subject_code,omitempty"`
	Unavailable         string `bson:"unavailable,omitempty"`
}

// Salary represents the salary statistical data for course (or subject)
type Salary struct {
	AggregationLevel    int    `bson:"aggregation_level,omitempty"` // enum
	HigherQuartileRange int    `bson:"higher_quartile_range,omitempty"`
	LowerQuartileRange  int    `bson:"lower_quartile_range,omitempty"`
	Median              int    `bson:"median,omitempty"`
	NumberOfGraduates   int    `bson:"number_of_graduates,omitempty"`
	ResponseRate        int    `bson:"response_rate,omitempty"`
	SubjectCode         string `bson:"subject_code,omitempty"`
	Unavailable         string `bson:"unavailable,omitempty"`
}

// RawSalary represents the salary statistical data for course (or subject) stored in its raw state
type RawSalary struct {
	AggregationLevel                                int    `bson:"aggregation_level,omitempty"`                                     // SALAGG
	InstitutionCourseSalarySixMonthsAfterGraduation *Stats `bson:"institution_course_salary_six_months_after_graduation,omitempty"` // INST
	KISMode                                         string `bson:"kis_mode"`
	KISCourseID                                     string `bson:"kis_course_id"`
	NumberOfStudents                                int    `bson:"number_of_students,omitempty"` // SALPOP
	PublicUKPRN                                     string `bson:"public_ukprn"`
	ResponseRate                                    int    `bson:"response_rate"`                                       // SALRESP_RATE
	SubjectCode                                     string `bson:"subject_code,omitempty"`                              // SALSBJ
	SubjectSalaryFortyMonthsAfterGraduation         *Stats `bson:"subject_salary_40_months_after_graduation,omitempty"` // LD
	SubjectSalarySixMonthsAfterGraduation           *Stats `bson:"subject_salary_six_months_after_graduation,omitempty"`
	UKPRN                                           string `bson:"ukprn"`
	Unavailable                                     string `bson:"unavailable,omitempty"`
}

// Stats contains a set of values for different statistical measurements of a dataset
type Stats struct {
	LowerQuartile int `bson:"lower_quartile_salary,omitempty"`             // LQ
	Median        int `bson:"median,omitempty"`                            // MED
	UpperQuartile int `bson:"upper_quartile_salary_for_subject,omitempty"` // UQ
}
