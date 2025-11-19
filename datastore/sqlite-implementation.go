package datastore

import (
	"bmgf-dashboard/datatypes"
	"database/sql"

	_ "modernc.org/sqlite" // pure Go sqlite driver
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) init() error {
	err := s.createTables()
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) createTables() error {
	query := `CREATE TABLE IF NOT EXISTS samples (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sample_name TEXT UNIQUE,
    sample_type TEXT,
    sample_category TEXT,
    sampling_site TEXT,
    milk_union TEXT,
    district TEXT,
    collection_date TEXT,
    rtpcr TEXT CHECK (rtpcr IN ('Positive', 'Negative', 'Untested', 'Suspected')) DEFAULT 'Untested');
	
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT,
		role TEXT
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) BulkInsert(samples []datatypes.SampleRecord) error {
	query := `INSERT OR IGNORE INTO samples (sample_name, sample_type, sample_category, sampling_site, milk_union, district, collection_date, rtpcr) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, sample := range samples {
		_, err := stmt.Exec(sample.SampleID,
			sample.SpecimenType,
			sample.SampleCategory,
			sample.SamplingSite,
			sample.MilkUnion,
			sample.District,
			sample.CollectionDate,
			sample.RTPCRResult)
		if err != nil {
			return err
		}
	}
	return nil
}
