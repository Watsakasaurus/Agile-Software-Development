package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"medicare-api/db"
	"medicare-api/server"
	"medicare-api/utils"
)

// Config represents server configuration
type Config struct {
	Postgres db.Config
	Env      string `envconfig:"ENV" default:"dev"`
}

func main() {
	var conf Config
	var migrate = flag.Bool("migrate", false, "do db migration")

	flag.Parse()

	log.SetLevel(log.InfoLevel)

	setupConf(&conf)

	if *migrate {
		migrateDB(conf.Postgres)
		os.Exit(0)
	}

	db, err := db.Connect(conf.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to Medicare DB")
	}

	server.Run(server.ContextParams{
		DB: db,
	})

	log.Exit(0)
}

func migrateDB(dbConf db.Config) {
	log.Info("Migrating Database Schema")
	dbc, err := db.Connect(dbConf)
	if err != nil {
		log.Fatalf("failed to connect to talaria db: %s", err)
	}

	err = dbc.Migrate("migrations", db.MigrationTargetLatest)
	if err != nil {
		log.Fatalf("failed to create schema: %s", err)
	}

	log.Info("Completed")
}

func setupConf(conf *Config) {
	// Load from env first as fall-back
	if err := envconfig.Process("medicare", conf); err != nil {
		log.Fatalf("Failed to load server config from env: %s", err)
	}

	if conf.Env == "prod" {
		// Get secret from AWS
		secret := utils.GetSecret()
		// If secret was successfully populated, unmarshal it into config
		if len(secret) > 0 {
			err := json.Unmarshal(secret, &conf.Postgres)
			if err != nil {
				log.Fatalf("Failed to unmarshal secret into conf: %s", err)
			}
		}
	}
}
