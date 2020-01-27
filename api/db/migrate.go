package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// needed for golang-migrate to read from filesystem
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// this is needed to enable postgres database support
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// MigrationTargetLatest represents a migration target of the latest version
const MigrationTargetLatest = 0

// Migrate does db migration up to the latest version
func (c *Client) Migrate(path string, target uint) error {
	driver, err := postgres.WithInstance(c.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Failed to to get postgres driver: %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", driver)
	if err != nil {
		return fmt.Errorf("Failed to create migration client: %s", err)
	}

	if target == MigrationTargetLatest {
		err = m.Up()
	} else {
		err = m.Migrate(target)
	}

	if err == migrate.ErrNoChange {
		log.Infof("Migration found nothing to do")
		return nil
	}
	return err
}
