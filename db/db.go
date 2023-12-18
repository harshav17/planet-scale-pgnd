package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DB struct {
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context

	// Datasource name.
	DSN string

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

func NewDB(DSN string) *DB {
	db := &DB{
		DSN: DSN,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() error {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	// Connect to the database.
	var err error
	if db.db, err = sql.Open("mysql", db.DSN); err != nil {
		return err
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) migrate() error {
	dbInstance, err := mysql.WithInstance(db.db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sourceInstance, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceInstance,
		"planetscale-pgnd", // TODO does this not matter at all?
		dbInstance,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
