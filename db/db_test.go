package db

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

type testDB struct {
	*DB
	container testcontainers.Container
}

func (tdb *testDB) TearDown() error {
	if err := tdb.Close(); err != nil {
		return err
	}
	// remove test container
	if err := tdb.container.Terminate(context.Background()); err != nil {
		return err
	}
	return nil
}

// Ensure the test database can open & close.
func TestDB(t *testing.T) {
	t.Parallel()

	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

// MustOpenDB returns a new, open DB. Fatal on error.
func MustOpenDB(tb testing.TB) *testDB {
	tb.Helper()

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, err := createMysqlContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	p, err := container.MappedPort(ctx, "3306")
	if err != nil {
		log.Fatalf("failed to get container external port: %v", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container external port: %v", err)
	}

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", DbUser, DbPass, hostIP, p.Port(), DbName)
	time.Sleep(time.Second)

	db := NewDB(DSN)
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}

	cancel()

	return &testDB{DB: db, container: container}
}

func createMysqlContainer(ctx context.Context) (testcontainers.Container, error) {
	mysqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8"),
		mysql.WithDatabase(DbName),
		mysql.WithUsername(DbUser),
		mysql.WithPassword(DbPass),
	)
	if err != nil {
		panic(err)
	}
	return mysqlContainer, nil
}

// MustCloseDB closes the DB. Fatal on error.
func MustCloseDB(tb testing.TB, db *testDB) {
	tb.Helper()
	if err := db.TearDown(); err != nil {
		tb.Fatal(err)
	}
}

// unit test NullTime methodspackage main
func TestNullTime_Scan(t *testing.T) {
	t.Run("scan nil", func(t *testing.T) {
		var nt NullTime
		err := nt.Scan(nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !time.Time(nt).IsZero() {
			t.Errorf("Expected zero time, got %v", nt)
		}
	})

	t.Run("scan byte slice", func(t *testing.T) {
		var nt NullTime
		err := nt.Scan([]byte("2022-01-01 00:00:00"))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
		if !time.Time(nt).Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, nt)
		}
	})

	t.Run("scan unsupported type", func(t *testing.T) {
		var nt NullTime
		err := nt.Scan(123)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestNullTime_Value(t *testing.T) {
	t.Run("value nil", func(t *testing.T) {
		var nt *NullTime
		v, err := nt.Value()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != nil {
			t.Errorf("Expected nil, got %v", v)
		}
	})

	t.Run("value zero time", func(t *testing.T) {
		nt := NullTime(time.Time{})
		v, err := nt.Value()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != nil {
			t.Errorf("Expected nil, got %v", v)
		}
	})

	t.Run("value non-zero time", func(t *testing.T) {
		nt := NullTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))
		v, err := nt.Value()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "2022-01-01 00:00:00"
		if v != expected {
			t.Errorf("Expected %v, got %v", expected, v)
		}
	})

	t.Run("PST time", func(t *testing.T) {
		nt := NullTime(time.Date(2022, 1, 1, 0, 0, 0, 0, time.FixedZone("PST", -8*60*60)))
		v, err := nt.Value()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expected := "2022-01-01 08:00:00"
		if v != expected {
			t.Errorf("Expected %v, got %v", expected, v)
		}
	})
}
