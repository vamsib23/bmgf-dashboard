package main

import (
	"bmgf-dashboard/datatypes"
	"encoding/csv"
	"io"
	"log"
	"os"
)

//var db *sql.DB

func main() {
	// database will be injected into main, currently using sqlite, may be changed in future
	// _, err := datastore.NewSQLiteStore("bmgf-data.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	GetDataFromCSVFile("misclleaneous\\BMGF Dashboard LSDV 18-11-2025.csv")
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
	i := 0
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
		if i == 10 {
			break
		}
		i++

	}
	return samples
}
