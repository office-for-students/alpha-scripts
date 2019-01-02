package data

// TariffDescriptions a list of tariff codes mapped to a description
var TariffDescriptions = map[string]string{
	"001": "less than 48 tariff points",
	"048": "between 48 and 63 tariff points",
	"064": "between 64 and 79 tariff points",
	"080": "between 80 and 95 tariff points",
	"096": "between 96 and 111 tariff points",
	"112": "between 112 and 127 tariff points",
	"128": "between 128 and 143 tariff points",
	"144": "between 144 and 159 tariff points",
	"160": "between 160 and 175 tariff points",
	"176": "between 176 and 191 tariff points",
	"192": "between 192 and 207 tariff points",
	"208": "between 208 and 223 tariff points",
	"224": "between 224 and 239 tariff points",
	"240": "more than 240 tariff points",
}

// Tariff contains information relating to the entry tariff points of students
type Tariff struct {
	AggregationLevel int            `bson:"aggregation_level,omitempty"` // TARAGG
	KISMode          string         `bson:"kis_mode"`
	KISCourseID      string         `bson:"kis_course_id"`
	NumberOfStudents int            `bson:"number_of_students,omitempty"` // TARPOP
	PublicUKPRN      string         `bson:"public_ukprn"`
	Tariffs          []*TariffStats `bson:"tariff,omitempty"`       // T**
	SubjectCode      string         `bson:"subject_code,omitempty"` // TARSBJ
	UKPRN            string         `bson:"ukprn"`
	Unavailable      bool           `bson:"unavailable"` // TARUNAVAILREASON
}

// TariffStats contains entry data for a particular tariff
type TariffStats struct {
	Code                 string `bson:"code"`
	Description          string `bson:"description"`
	ProportionOfEntrants int    `bson:"proportion_of_entrants"`
}
