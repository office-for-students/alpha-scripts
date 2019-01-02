package handlers

import "github.com/ofs/alpha-scripts/mongo/load-data/general-data-builder/mongo"

var fileExtension = ".csv"

// Common ...
type Common struct {
	Mongo                *mongo.Mongo
	RelativeFileLocation string
}
