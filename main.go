package main

import (
	"bmgf-dashboard/datastore"
	"bmgf-dashboard/datatypes"
	"encoding/csv"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"middlewares"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DB *datastore.SQLiteStore

const PORT = ":2025"

type ThemeConfig struct {
	ActiveTheme string `json:"active_theme"`
}

func LoadThemeConfig() (string, error) {
	// @TODO: replace reading from config file with reading from database
	data, err := os.ReadFile("themes/config.json")
	if err != nil {
		return "", err
	}

	var cfg ThemeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	return cfg.ActiveTheme, nil
}

func main() {
	// database will be injected into main, currently using sqlite, may be changed in future
	DB, err := datastore.NewSQLiteStore("bmgf-data.db")
	if err != nil {
		log.Fatal("Error getting the database:", err)
	}
	defer DB.Close()

	// data := GetDataFromCSVFile("miscellaneous\\BMGF Dashboard LSDV 18-11-2025.csv")
	// if err := DB.BulkInsert(data); err != nil {
	// 	log.Fatal("error inserting bulk data:", err)
	// }

	activeTheme, _ := LoadThemeConfig()

	themePath := filepath.Join("themes", activeTheme)

	// Serve static assets (CSS, JS, Images)
	assetsPath := filepath.Join(themePath, "assets")

	// Serve template files
	templates := template.Must(template.ParseGlob(filepath.Join(themePath, "*.html")))

	serv := http.NewServeMux()
	serv.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))
	serv.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index.html", nil)
	})
	serv.HandleFunc("/data/all", func(w http.ResponseWriter, r *http.Request) {
		samples, err := DB.GetInfoForPublicAPI()
		if err != nil {
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(samples)
	})
	servWithlogging := middlewares.Logging(serv)
	log.Println("Starting webserver at", PORT)
	log.Fatal(http.ListenAndServe(PORT, servWithlogging))
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
		SampleCollectionDt, err := time.Parse("02-01-2006 15:04", record[7])
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
			SampleUniqueID: record[1],
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
