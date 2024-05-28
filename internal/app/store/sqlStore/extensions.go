package sqlStore

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

// TODO найти инфу как правильно указывать пути в go
const migrationsPath = "file://../../../../../migrations"

// TestDb ...
func TestDb(t *testing.T, databaseType, databaseURL string) (*sql.DB, func(string)) {
	t.Helper()

	db, err := sql.Open(databaseType, databaseURL)

	if err != nil {
		t.Fatal(err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		t.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		databaseType,
		driver,
	)

	if err != nil {
		t.Fatal(err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(table string) {
		_, _ = db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		_ = db.Close()
	}
}
