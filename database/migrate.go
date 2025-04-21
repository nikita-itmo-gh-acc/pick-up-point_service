package database

import (
	"pvz_service/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrationTool 	*migrate.Migrate
}

func (m *Migrator) Init(path string, url string) error {
	newMigrationTool, err := migrate.New("file://" + path, url)
	if err != nil {
		logger.Err.Println("Error occured trying to init migration tool - ", err)
		return err
	}

	m.migrationTool = newMigrationTool
	return nil
}

func (m *Migrator) Apply() error {
	err := m.migrationTool.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Err.Println("Can't apply migrations - ", err)
		return err
	}
	
	if err == migrate.ErrNoChange {
		logger.Debug.Println("Nothing to apply")
	}

	logger.Debug.Println("Migration(s) applied successfully")
	return nil
}

func (m *Migrator) RollBack(steps int) error {
	err := m.migrationTool.Steps(-steps)
	if err != nil && err != migrate.ErrNoChange {
		logger.Err.Println("Rollback failed - ", err)
		return err
	}
	
	if err == migrate.ErrNoChange {
		logger.Debug.Println("Nothing to rollback")
	}

	logger.Debug.Println("Rollback complete")
	return nil
}
