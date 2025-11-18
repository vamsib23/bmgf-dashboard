package datastore

import "bmgf-dashboard/datatypes"

// defines an interface for the datastore
type Datastore interface {
	InsertSample(record datatypes.SampleRecord) error
	BulkInsert(samples []datatypes.SampleRecord) error
	GetAllSamples() ([]datatypes.SampleRecord, error)
}
