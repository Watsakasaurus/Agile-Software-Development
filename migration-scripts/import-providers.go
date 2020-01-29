package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type config struct {
	User     string `envconfig:"postgres_user" default:"cms" json:"username"`
	Password string `envconfig:"postgres_password" default:"secret" json:"password"`
	DBName   string `default:"cms"`
	Host     string `envconfig:"postgres_host" default:"127.0.0.1" json:"host"`
	Port     int    `envconfig:"postgres_port" default:"5432" json:"port"`
}

func main() {
	const (
		code = 0
		lat  = 3
		long = 4
	)

	f, err := os.Open("providers.csv")
	checkErr(err)
	defer f.Close()

	r := csv.NewReader(f)

	var conf config
	if err := envconfig.Process("medicare", &conf); err != nil {
		log.Fatalf("Failed to load server config from env: %s", err)
	}

	cn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.User,
		conf.Password)

	db, err := sqlx.Open("postgres", cn)
	checkErr(err)
	defer db.Close()

	drgs := map[string]struct{}{}
	providers := map[string]struct{}{}

	// some zips are missing in zips csv
	missingZips := map[string]struct {
		Name   string
		Street string
		City   string
		State  string
	}{}

	sqlProcedures := "INSERT INTO procedures(total_discharges, average_total_payments, average_covered_charges, average_medicare_payments, drg_definition) VALUES($1, $2, $3, $4, $5)"
	sqlProviders := "INSERT INTO providers(id, name, street, city, state, zip_code, hrr_description) VALUES($1, $2, $3, $4, $5, $6, $7)"
	sqlProviderProcedure := "INSERT INTO provider_procedures(provider_id, procedure_id) VALUES($1, $2)"

	tx := db.MustBegin()

	record, err := r.Read() // skip title
	for {
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)

		data := strings.SplitN(record[0], " - ", 2)
		drgID := data[0]
		drgDefinition := record[0]
		id := record[1]
		name := record[2]
		street := record[3]
		city := record[4]
		state := record[5]
		zipCode := fmt.Sprintf("%05s", record[6])
		hrrDescription := record[7]
		totalDischarges := record[8]
		averageCoveredCharges := record[9]
		averageTotalPayments := record[10]
		averageMedicarePayments := record[11]

		if _, found := drgs[drgID]; !found {
			tx.MustExec(sqlProcedures, totalDischarges, averageTotalPayments, averageCoveredCharges, averageMedicarePayments, drgDefinition)
			drgs[drgID] = struct{}{}
		}

		if _, found := providers[id]; !found {
			zip := ""
			tx.Get(&zip, "SELECT zip_code FROM zip_code_lat_long WHERE zip_code=$1", zipCode)
			if zip == "" {
				missingZips[zipCode] = struct {
					Name   string
					Street string
					City   string
					State  string
				}{
					Name:   name,
					Street: street,
					City:   city,
					State:  state,
				}
				continue
			}
			tx.Exec(sqlProviders, id, name, street, city, state, zipCode, hrrDescription)
			providers[id] = struct{}{}
		}

		tx.MustExec(sqlProviderProcedure, id, drgID)
	}
	err = tx.Commit()
	checkErr(err)

	fmt.Println("These zip codes are misssing in zip_code_lat_long table:")
	for k, v := range missingZips {
		fmt.Printf("Zip: %s, State: %s, City: %s\n", k, v.State, v.City)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
