package data

// Course represents a course resource
type Course struct {
	ApplicationProvider string             `bson:"application_provider,omitempty"`
	Country             *Country           `bson:"country"`
	DistanceLearning    *DistanceLearning  `bson:"distance_learning"`
	Foundation          string             `bson:"foundation_year_availability"` // enum
	Honours             bool               `bson:"honours_award_provision"`
	ID                  string             `bson:"_id"`
	Institution         *InstitutionObject `bson:"institution"`
	KISCourseID         string             `bson:"kis_course_id"`
	Length              *LengthObject      `bson:"length_of_course"`
	Links               *LinkList          `bson:"links"`
	Location            *Location          `bson:"location"`
	Mode                *Mode              `bson:"mode"` // enum - part time, full time, both
	NHSFunded           *NHSFunded         `bson:"nhs_funded,omitempty"`
	Qualification       *Qualification     `bson:"qualification"`
	SandwichYear        *Availability      `bson:"sandwich_year"`
	Statistics          *Statistics        `bson:"statistics,omitempty"`
	Subject             *Subject           `bson:"subject"`
	Title               *Language          `bson:"title"`
	UCASCode            string             `bson:"ucas_code_id,omitempty"`
	YearAbroad          *Availability      `bson:"year_abroad"`
}

// Availability represents an object referring to the availability
type Availability struct {
	Code  string `bson:"code"`
	Label string `bson:"label"` // enum , 0-2
}

// Country represents a country object
type Country struct {
	Code string `bson:"code"`
	Name string `bson:"name"`
}

// DistanceLearning represents an object referring
// to the course available through distance learning
type DistanceLearning struct {
	Code  string `bson:"code"`
	Label string `bson:"label"`
}

// InstitutionObject represents institution data related to course
type InstitutionObject struct {
	PublicUKPRNName string `bson:"public_ukprn_name"`
	PublicUKPRN     string `bson:"public_ukprn"`
	UKPRNName       string `bson:"ukprn_name"`
	UKPRN           string `bson:"ukprn"`
}

// Language represents an object containing english or welsh strings
type Language struct {
	English string `bson:"english,omitempty"`
	Welsh   string `bson:"welsh,omitempty"`
}

// LengthObject represents an object referring to the course length
type LengthObject struct {
	Code  string `bson:"code"`
	Label string `bson:"label"`
}

// LinkList represents a list of links related to resource
type LinkList struct {
	Accommodation       *Language `bson:"accommodation,omitempty"`
	AssessmentMethod    *Language `bson:"assessment_method,omitempty"`         // ASSURL
	CoursePage          *Language `bson:"course_page,omitempty"`               // CRSEURL
	EmploymentDetails   *Language `bson:"employment_details,omitempty"`        // EMPLOYURL
	FinancialSupport    *Language `bson:"financial_support_details,omitempty"` // SUPPORTURL
	Institution         string    `bson:"institution"`
	LearningAndTeaching *Language `bson:"learning_and_teaching_methods,omitempty"` // LTURL
	Self                string    `bson:"self"`
	StudentUnion        *Language `bson:"student_union,omitempty"`
}

// Location represents an object containing fields to enable one to locate institution
type Location struct {
	Changes   bool      `bson:"changes"`
	Latitude  string    `bson:"latitude"`
	Longitude string    `bson:"longitude"`
	Name      *Language `bson:"name"`
}

// Mode represents an object referring to the type of course
type Mode struct {
	Code  string `bson:"code"`
	Label string `bson:"label"`
}

// NHSFunded represents an object referring to the course having any NHS funded students
type NHSFunded struct {
	Code  string `bson:"code,omitempty"`
	Label string `bson:"label,omitempty"`
}

// Qualification represents an object referring to the qualification received from course
type Qualification struct {
	Code  string `bson:"code"`
	Label string `bson:"label"`
	Level string `bson:"level"`
	Name  string `bson:"name"`
}

// LocationIDObject represents a course location object
type LocationIDObject struct {
	ID string `bson:"id"`
}
