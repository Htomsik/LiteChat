package Server

import (
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

const migrationsPath = "file://migrations"

// Start ...
func Start(config *Config) error {

	srv := newServer()

	srv.logger.Infof("Server started on port %v", config.Port)

	return http.ListenAndServe(config.Port, srv)
}
