package datatypes

import "time"

// for reading and writing to the samples table in the database
type SampleRecord struct {
	ID             int
	SampleUniqueID string
	SpecimenType   string
	SampleCategory string
	SamplingSite   string
	District       string
	MilkUnion      string
	CollectionDate time.Time
	RTPCRResult    string
}
