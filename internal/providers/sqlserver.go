package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// SQLServer conecta em instâncias SQL Server on-prem.
type SQLServer struct{}

func (SQLServer) Engine() string { return "sqlserver" }

func (SQLServer) Open(cfg Config) (*sql.DB, error) {
	port := cfg.Port
	if port == 0 {
		port = 1433
	}

	var conn string
	if cfg.User != "" {
		conn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;",
			cfg.Host, cfg.User, cfg.Password, port)
	} else {
		conn = fmt.Sprintf("server=%s;trusted_connection=yes;port=%d;", cfg.Host, port)
	}

	db, err := sql.Open("sqlserver", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (SQLServer) DetectVersion(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT @@VERSION").Scan(&version)
	return version, err
}
