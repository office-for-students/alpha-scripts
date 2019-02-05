package data

// Qualification represents attributes associated with qualification
type Qualification struct {
	Code  string `bson:"code"`
	Label string `bson:"label"`
	Level string `bson:"level"`
	Name  string `bson:"name"`
}
