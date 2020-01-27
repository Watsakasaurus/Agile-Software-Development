package controllers_test

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"medicare-api/db"
	"medicare-api/server"
)

var (
	dbc *db.Client
	dbx *sqlx.DB

	r *echo.Echo

	params *server.ContextParams
)

func setup() {
	var conf db.Config
	// Load from env first then overwrite
	if err := envconfig.Process("medicare", &conf); err != nil {
		log.Fatalf("Failed to load server config from env: %s", err)
	}

	var err error
	dbc, err = db.Connect(db.Config{
		User:     "test",
		DBName:   "test_cms",
		Password: "test",
		Host:     conf.Host,
		Port:     conf.Port,
	})
	if err != nil {
		panic(err)
	}
	dbx = dbc.DB()

	wipeDB()
	err = dbc.Migrate("../migrations", db.MigrationTargetLatest)
	if err != nil {
		panic(err)
	}

	params = &server.ContextParams{
		DB: dbc,
	}

	r = server.CreateRouter(*params)
}

func TestMain(m *testing.M) {
	setup()
	cleanup()
	os.Exit(m.Run())
}

func wipeDB() {
	dbx.MustExec(`
	DROP TABLE IF EXISTS inpatient_charge_data;

	DROP TABLE IF EXISTS schema_migrations;
	`)

}

func cleanup() {
	dbx.MustExec(`
	DELETE FROM inpatient_charge_data;
	`)

	// reset logging level, in case it was set in a sub-test
	log.SetLevel(log.FatalLevel)
}
