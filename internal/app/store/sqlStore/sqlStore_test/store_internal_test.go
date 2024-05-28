package sqlStore_test

import (
	"os"
	"testing"
)

var (
	databaseURL  string
	databaseType string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	databaseType = os.Getenv("DATABASE_TYPE")

	if databaseURL == "" {
		// TODO найти инфу как правильно указывать пути в go
		databaseURL = "../../../../../assets/db/test-store.db"
	}

	if databaseType == "" {
		databaseType = "sqlite3"
	}

	os.Exit(m.Run())
}
