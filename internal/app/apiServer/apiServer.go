package apiServer

import (
	"Chat/internal/app/store/sqlStore"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

const migrationsPath = "file://migrations"

// Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseType, config.DatabaseURL, config.DatabaseAutoMigration)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlStore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionsKey))

	srv := newServer(store, sessionStore)

	srv.logger.Infof("Server started on port %v", config.Port)

	return http.ListenAndServe(config.Port, srv)
}

// newDB ...
func newDB(dbType, dbBaseURL string, dbAutoMigrations bool) (*sql.DB, error) {
	db, err := sql.Open(dbType, dbBaseURL)

	if err != nil {
		return nil, err
	}

	if dbAutoMigrations {
		if err := migrateDB(db, dbType); err != nil {
			return nil, err
		}
	}

	return db, db.Ping()
}

// migrateDB ...
func migrateDB(db *sql.DB, databaseType string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		databaseType,
		driver,
	)

	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
