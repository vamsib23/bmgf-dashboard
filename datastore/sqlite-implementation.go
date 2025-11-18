package datastore

import (
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
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS samples (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sample_id TEXT,
		sample_type TEXT,
		category TEXT,
		sampling_site TEXT,
		milk_union TEXT,
		district TEXT,
		collection_date TEXT,
		rtpcr enum('Positive', 'Negative', 'Untested', 'Suspected')	
	);`)
	return err
}
