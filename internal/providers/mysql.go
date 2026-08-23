package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL conecta em MySQL/MariaDB.
type MySQL struct{}

func (MySQL) Engine() string { return "mysql" }

func (MySQL) Open(cfg Config) (*sql.DB, error) {
	port := cfg.Port
	if port == 0 {
		port = 3306
	}
	dbName := cfg.Database
	if dbName == "" {
		dbName = "sys"
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.User, cfg.Password, cfg.Host, port, dbName)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func (MySQL) DetectVersion(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	return version, err
}
