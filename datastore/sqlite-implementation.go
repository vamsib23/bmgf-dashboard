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

func (s *SQLiteStore) InsertSample(rec datatypes.SampleRecord) error {
	_, err := s.db.Exec(`
	INSERT INTO samples 
	(sample_id, sample_type, category, sampling_site, milk_union, district, collection_date, rtpcr)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.SampleID,
		rec.SampleType,
		rec.Category,
		rec.SamplingSite,
		rec.MilkUnion,
		rec.District,
		rec.CollectionDate.Format("2006-01-02"),
		rec.RTPCR,
	)
	return err
}

func (s *SQLiteStore) BulkInsert(samples []datatypes.SampleRecord) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	INSERT INTO samples 
	(sample_id, sample_type, category, sampling_site, milk_union, district, collection_date, rtpcr)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	for _, rec := range samples {
		_, err := stmt.Exec(
			rec.SampleID,
			rec.SampleType,
			rec.Category,
			rec.SamplingSite,
			rec.MilkUnion,
			rec.District,
			rec.CollectionDate.Format("2006-01-02"),
			rec.RTPCR,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
