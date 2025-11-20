package datastore

import "bmgf-dashboard/datatypes"

// defines an interface for the datastore
type Datastore interface {
	InsertSample(record datatypes.SampleRecord) error
	BulkInsert(samples []datatypes.SampleRecord) error
	// This function must implemented to fetch data for public API
	// The following contact must be satisfied:
	/*

	 */
	GetInfoForPublicAPI() ([]datatypes.SampleRecord, error)
	GetActiveTheme() (string, error)
}
