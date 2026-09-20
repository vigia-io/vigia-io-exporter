// Package collector expõe a pipeline de coleta do exporter para consumo por vigia-agent.
package collector

import (
	"fmt"

	"github.com/vigia-io/vigia-io-exporter/internal/export"
	"github.com/vigia-io/vigia-io-exporter/internal/providers"
)

const (
	ExitOK         = export.ExitOK
	ExitConnection = export.ExitConnection
)

// Options de coleta compartilhadas entre CLI e agent.
type Options struct {
	Engine        string
	Host          string
	Port          int
	User          string
	Password      string
	Database      string
	HostAlias     string
	Output        string
	CollectorName string
	Version       string
	Quiet         bool
}

// Run executa a coleta e grava snapshot no caminho de saída.
func Run(opts Options) (int, error) {
	provider, err := ProviderForEngine(opts.Engine)
	if err != nil {
		return export.ExitConnection, err
	}
	if opts.CollectorName == "" {
		opts.CollectorName = "vigia-agent"
	}
	return export.Run(export.Options{
		Provider:      provider,
		Host:          opts.Host,
		Port:          opts.Port,
		User:          opts.User,
		Password:      opts.Password,
		Database:      opts.Database,
		HostAlias:     opts.HostAlias,
		Output:        opts.Output,
		Format:        "v1",
		CollectorName: opts.CollectorName,
		Version:       opts.Version,
		Quiet:         opts.Quiet,
	})
}

// ProviderForEngine resolve o provider pelo nome da engine.
func ProviderForEngine(engine string) (providers.Provider, error) {
	switch engine {
	case "sqlserver", "sql":
		return providers.SQLServer{}, nil
	case "mysql":
		return providers.MySQL{}, nil
	case "azuresql":
		return providers.AzureSQL{}, nil
	default:
		return nil, fmt.Errorf("engine não suportada: %s", engine)
	}
}
