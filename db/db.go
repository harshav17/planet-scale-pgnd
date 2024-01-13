package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"log"
	"strings"
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

// NullTime represents a helper wrapper for time.Time. It automatically converts
// time fields to/from RFC 3339 format. Also supports NULL for zero time.
type NullTime time.Time

// Scan reads a time value from the database.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if v, ok := value.([]byte); ok {
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}

		*(*time.Time)(n) = t
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

// Value formats a time value for the database.
func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format("2006-01-02 15:04:05"), nil
}

// TODO consider moving to a DB util class
type findWhereClause struct {
	columns []string
	values  []interface{}
}

func (w *findWhereClause) Add(column string, value interface{}) {
	w.columns = append(w.columns, column)
	w.values = append(w.values, value)
}

func (w *findWhereClause) ToClause() string {
	s := strings.Builder{}
	if len(w.columns) > 0 {
		s.WriteString("WHERE ")
		for i, column := range w.columns {
			if i > 0 {
				s.WriteString(" AND ")
			}
			s.WriteString(column + " = ?")
		}
	}
	return s.String()
}
