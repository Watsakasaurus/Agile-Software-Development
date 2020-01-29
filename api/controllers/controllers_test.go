package controllers_test

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"

	"medicare-api/db"
	"medicare-api/server"
	"medicare-api/types"
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
	DROP TABLE IF EXISTS provider_procedures;
	DROP TABLE IF EXISTS providers;
	DROP TABLE IF EXISTS procedures;
	DROP TABLE IF EXISTS zip_code_lat_long;

	DROP TABLE IF EXISTS schema_migrations;
	`)

}

func cleanup() {
	dbx.MustExec(`
	DELETE FROM provider_procedures;
	DELETE FROM providers;
	DELETE FROM procedures;
	DELETE FROM zip_code_lat_long;
	`)

	// reset logging level, in case it was set in a sub-test
	log.SetLevel(log.FatalLevel)
}

func withCleanup(f func()) func() {
	return func() {
		Reset(cleanup)
		f()
	}
}

func insertProvider(payload types.Provider) *types.Provider {
	err := dbc.CreateProvider(payload)
	if err != nil {
		panic(err)
	}

	provider, err := dbc.GetProviderByID(payload.ID)
	if err != nil {
		panic(err)
	}

	return provider
}

func insertProcedure(payload types.Procedure) *types.Procedure {
	err := dbc.CreateProcedure(payload)
	if err != nil {
		panic(err)
	}

	procedure, err := dbc.GetProcedureByID(payload.ID)
	if err != nil {
		panic(err)
	}

	procedure.AverageTotalPayments = payload.AverageTotalPayments
	return procedure
}

func insertZipCode(payload types.ZipCodeLatLong) *types.ZipCodeLatLong {
	record, err := dbc.CreateZipCodeLatLong(payload)
	if err != nil {
		panic(err)
	}
	return record
}

func assignProcedureToProvider(payload types.ProvideRrocedure) {
	err := dbc.AssignProcedureToProvider(payload)
	if err != nil {
		panic(err)
	}
}
