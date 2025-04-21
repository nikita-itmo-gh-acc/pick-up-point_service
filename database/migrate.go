package database

import (
	"fmt"

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
		fmt.Println("Error occured trying to init migration tool - ", err)
		return err
	}

	m.migrationTool = newMigrationTool
	return nil
}

func (m *Migrator) Apply() error {
	err := m.migrationTool.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Can't apply migrations - ", err)
		return err
	}
	
	if err == migrate.ErrNoChange {
		fmt.Println("Nothing to apply")
	}

	fmt.Println("Migration(s) applied successfully")
	return nil
}

func (m *Migrator) RollBack(steps int) error {
	err := m.migrationTool.Steps(-steps)
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Rollback failed - ", err)
		return err
	}
	
	if err == migrate.ErrNoChange {
		fmt.Println("Nothing to rollback")
	}

	fmt.Println("Rollback complete")
	return nil
}
