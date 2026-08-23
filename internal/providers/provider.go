package providers

import "database/sql"

// Config comum de conexão.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Provider abre conexão read-only e detecta versão do engine.
type Provider interface {
	Engine() string
	Open(cfg Config) (*sql.DB, error)
	DetectVersion(db *sql.DB) (string, error)
}
