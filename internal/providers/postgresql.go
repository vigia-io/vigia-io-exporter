package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgreSQL conecta em PostgreSQL (role read-only recomendada).
type PostgreSQL struct{}

func (PostgreSQL) Engine() string { return "postgresql" }

func (PostgreSQL) Open(cfg Config) (*sql.DB, error) {
	port := cfg.Port
	if port == 0 {
		port = 5432
	}
	dbName := cfg.Database
	if dbName == "" {
		dbName = "postgres"
	}

	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, port, cfg.User, cfg.Password, dbName,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (PostgreSQL) DetectVersion(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	return version, err
}
