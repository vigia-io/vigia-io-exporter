package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// AzureSQL conecta em Azure SQL Database.
type AzureSQL struct{}

func (AzureSQL) Engine() string { return "azuresql" }

func (AzureSQL) Open(cfg Config) (*sql.DB, error) {
	conn := fmt.Sprintf("server=%s;user id=%s;password=%s;",
		cfg.Host, cfg.User, cfg.Password)

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

func (AzureSQL) DetectVersion(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT @@VERSION").Scan(&version)
	return version, err
}
