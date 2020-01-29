package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

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

	f, err := os.Open("zip-codes-new.csv")
	checkErr(err)
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

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

	sql := "INSERT INTO zip_code_lat_long(zip_code, latitude, longitude) VALUES($1, $2, $3)"

	tx := db.MustBegin()

	record, err := r.Read() // skip title
	for {
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)

		tx.MustExec(sql, record[code], record[lat], record[long])
	}
	err = tx.Commit()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
