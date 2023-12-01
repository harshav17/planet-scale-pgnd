package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

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

	// TODO split DATABASE_URL into vars for migration use cases with protocol: cockroachdb
	// if err := db.migrate(); err != nil {
	// 	return fmt.Errorf("migrate: %w", err)
	// }

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

// TODO add migration support
func (db *DB) migrate() error {
	m, err := migrate.New(
		"file://../../db/migrations",
		db.DSN)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
