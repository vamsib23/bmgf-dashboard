package datatypes

import "time"

// for reading and writing to the samples table in the database
type SampleRecord struct {
	ID             int       // ID representing the unique addressable id of the sample
	SampleID       string    // sample name provided by the project user
	SampleType     string    //
	Category       string    // Clinical / Environmental
	MilkUnion      string    // Name of milk union
	CollectionDate time.Time // only the date on which the sample was collected.
	RTPCR          string    // Positive / Negative/ untested/ suspected
}

// for reading and writing to the locationdetails database table.
type LocationDetails struct {
	ID           int
	SamplingSite string // Name of area in which sample was collected
	District     string
}

// for reading and writing to the milkunions database table.
type MilkUnion struct {
	ID        int
	MilkUnion string
}
