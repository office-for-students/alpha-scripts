package data

// Institution represents an institution resource
type Institution struct {
	APROutcome  string      `bson:"apr_outcome"`
	Country     *Country    `bson:"country"`
	Links       *LinkList   `bson:"links"`
	Locations   []*Location `bson:"locations"`
	Name        string      `bson:"name"`
	TEFOutcome  string      `bson:"tef_outcome"`
	PublicUKPRN string      `bson:"public_ukprn"`
	UKPRN       string      `bson:"ukprn"`
}

// Country represents a country object
type Country struct {
	Code string `bson:"code"`
	Name string `bson:"name"`
}

// LinkList represents a list of links related to resource
type LinkList struct {
	Courses                 string    `bson:"courses"`
	InstitutionStudentUnion *Language `bson:"institution_student_union,omitempty"`
	Self                    string    `bson:"self"`
}

// Location represents an object containing fields to enable one to locate institution
type Location struct {
	ID        string         `bson:"id"`
	Latitude  string         `bson:"latitude"`
	Links     *LocationLinks `bson:"links"`
	Longitude string         `bson:"longitude"`
	Name      *Language      `bson:"name,omitempty"`
}

// LocationLinks represents a list of links related to location
type LocationLinks struct {
	Accommodation *Language `bson:"accommodation,omitempty"`
	StudentUnion  *Language `bson:"student_union,omitempty"`
}

// Language represents an object containing english or welsh strings
type Language struct {
	English string `bson:"english,omitempty"`
	Welsh   string `bson:"welsh,omitempty"`
}
