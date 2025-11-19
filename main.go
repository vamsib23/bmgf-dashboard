package main

import (
	"bmgf-dashboard/datastore"
	"bmgf-dashboard/datatypes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var DB *datastore.SQLiteStore

func main() {
	// database will be injected into main, currently using sqlite, may be changed in future
	DB, err := datastore.NewSQLiteStore("bmgf-data.db")
	if err != nil {
		log.Fatal("Error getting the database:", err)
	}
	defer DB.Close()

	data := GetDataFromCSVFile("miscellaneous\\BMGF Dashboard LSDV 18-11-2025.csv")
	if err := DB.BulkInsert(data); err != nil {
		log.Fatal("error inserting bulk data:", err)
	}

}

func GetDataFromCSVFile(filename string) []datatypes.SampleRecord {
	// read and parse csv file into datatypes.SampleRecord
	var samples []datatypes.SampleRecord
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Could not read the file", err)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("error parsing the records", err)
			return nil
		}
		log.Println(record)
		SampleCollectionDt, err := time.Parse("02-01-2006", record[7])
		if err != nil {
			log.Println("error parsing the date", err)
			return nil
		}
		rtpcrresult := ""
		switch strings.ToLower(record[8]) {
		case "yes":
			rtpcrresult = "Positive"
		case "no":
			rtpcrresult = "Negative"
		case "":
			rtpcrresult = "Untested"
		default:
			log.Println("Unknown case", record[8])
			rtpcrresult = "Suspected"
		}
		samples = append(samples, datatypes.SampleRecord{
			SampleID:       record[1],
			SpecimenType:   record[2],
			SampleCategory: record[3],
			SamplingSite:   record[4],
			MilkUnion:      record[5],
			District:       record[6],
			CollectionDate: SampleCollectionDt,
			RTPCRResult:    rtpcrresult,
		})
	}
	return samples
}
