package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	migration *migrate.Migrate
}

func NewMigration(sourceDir, url string) (*Migration, error) {
	// sourceDir e.g. "file://etc/migrations"
	m, err := migrate.New(sourceDir, url)
	if err != nil {
		return nil, fmt.Errorf("migrate.New: %+v", err)
	}
	log.Println("Migration instance created")
	return &Migration{migration: m}, nil
}

func (m *Migration) MigrateOneStepUp() error {
	err := m.migration.Steps(1) // up 1
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no change")
			return nil
		}
		return err
	}
	return nil
}

func (m *Migration) MigrateOneStepDown() error {
	err := m.migration.Steps(-1) // down 1
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no change")
			return nil
		}
		return err
	}
	return nil
}

func (m *Migration) Up() error {
    err := m.migration.Up()
    if err != nil && err != migrate.ErrNoChange {
        return err
    }
    if err == migrate.ErrNoChange {
        log.Println("No migration changes needed.")
    } else {
        log.Println("Migration Up successful.")
    }
    return nil
}


func (m *Migration) ForceMigrate(version int) error {
	err := m.migration.Force(version)
	if err != nil {
		log.Println("force migration error:", err)
		return err
	}
	return nil
}

// ConstructMigrationUrl: postgres://user:pass@host:port/dbname?options
func ConstructMigrationUrl(user, password, host, port, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
}
