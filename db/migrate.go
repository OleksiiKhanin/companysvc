package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateSchema(storage *sql.DB, migrateFilesPath string) error {
	driver, err := postgres.WithInstance(storage, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("get driver for migration: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrateFilesPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
