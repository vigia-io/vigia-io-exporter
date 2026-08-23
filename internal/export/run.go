package export

import (
	"errors"
	"fmt"
	"os"

	"github.com/vigia-io/vigia-io-exporter/internal/data"
	"github.com/vigia-io/vigia-io-exporter/internal/providers"
	"github.com/vigia-io/vigia-io-exporter/internal/scripts"
	"github.com/vigia-io/vigia-io-exporter/internal/snapshot"
)

// Exit codes (non-interactive-cli.md).
const (
	ExitOK = iota
	ExitConnection
	ExitScript
	ExitWrite
	ExitValidation
)

// Options de coleta.
type Options struct {
	Provider       providers.Provider
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	HostAlias      string
	Output         string
	Format         string
	CollectorName  string
	Version        string
	Quiet          bool
	SkipVersionSQL bool
}

// Run executa a coleta e grava o arquivo de saída.
func Run(opts Options) (int, error) {
	if opts.CollectorName == "" {
		opts.CollectorName = "vigia-exporter"
	}
	if opts.Version == "" {
		opts.Version = "0.1.0"
	}
	if opts.Output == "" {
		opts.Output = "snapshot.json"
	}
	if opts.Format == "" {
		opts.Format = "v1"
	}
	if opts.HostAlias == "" {
		opts.HostAlias = "default"
	}

	scripts.Init()

	cfg := providers.Config{
		Host:     opts.Host,
		Port:     opts.Port,
		User:     opts.User,
		Password: opts.Password,
		Database: opts.Database,
	}

	db, err := opts.Provider.Open(cfg)
	if err != nil {
		return ExitConnection, fmt.Errorf("falha ao conectar no banco: %w", err)
	}
	defer db.Close()

	rows, err := data.GetData(db, scripts.Scripts)
	if err != nil {
		return ExitScript, fmt.Errorf("falha ao executar scripts: %w", err)
	}

	if opts.Format == "legacy" {
		if err := snapshot.WriteLegacy(opts.Output, rows); err != nil {
			return ExitWrite, fmt.Errorf("falha ao gravar saída: %w", err)
		}
		if !opts.Quiet {
			fmt.Fprintf(os.Stderr, "Arquivo %s gerado (formato legado).\n", opts.Output)
		}
		return ExitOK, nil
	}

	engineVersion := ""
	if !opts.SkipVersionSQL {
		if v, err := opts.Provider.DetectVersion(db); err == nil {
			engineVersion = v
		}
	}

	doc := snapshot.NewBuilder(opts.CollectorName, opts.Version, opts.Provider.Engine(), opts.HostAlias).
		Scripts(rows).
		Metrics([]snapshot.MetricPoint{}).
		EngineVersion(engineVersion).
		Build()

	if err := snapshot.Write(opts.Output, doc); err != nil {
		return ExitWrite, fmt.Errorf("falha ao gravar snapshot: %w", err)
	}

	if !opts.Quiet {
		fmt.Fprintf(os.Stderr, "Arquivo %s gerado com sucesso.\n", opts.Output)
	}
	return ExitOK, nil
}

// ErrHelp indica uso incorreto da CLI.
var ErrHelp = errors.New("help")
