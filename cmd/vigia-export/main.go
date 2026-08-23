package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vigia-io/vigia-io-exporter/internal/export"
	"github.com/vigia-io/vigia-io-exporter/internal/providers"
)

var (
	version = "0.1.0"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		if err == export.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "vigia-export",
		Short: "Coleta offline Vigia — gera snapshot.json",
		Long:  "CLI read-only do Vigia. Produz snapshot.json (schema v1) para ingest na plataforma.",
	}

	root.AddCommand(
		newProviderCmd("sqlserver", "SQL Server on-premises", providers.SQLServer{}, false),
		newProviderCmd("sql", "alias de sqlserver", providers.SQLServer{}, false),
		newProviderCmd("mysql", "MySQL / MariaDB", providers.MySQL{}, true),
		newProviderCmd("azuresql", "Azure SQL Database", providers.AzureSQL{}, true),
	)

	root.Version = version
	return root
}

type providerFlags struct {
	host         string
	port         int
	user         string
	password     string
	passwordEnv  string
	database     string
	alias        string
	output       string
	format       string
	quiet        bool
	skipVersion  bool
}

func newProviderCmd(use, short string, provider providers.Provider, requireUser bool) *cobra.Command {
	flags := &providerFlags{}

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := export.Options{
				Provider:       provider,
				Host:           flags.host,
				Port:           flags.port,
				User:           flags.user,
				Password:       flags.password,
				Database:       flags.database,
				HostAlias:      flags.alias,
				Output:         flags.output,
				Format:         flags.format,
				Version:        version,
				Quiet:          flags.quiet,
				SkipVersionSQL: flags.skipVersion,
			}

			if flags.passwordEnv != "" {
				pw, err := export.PasswordFromEnv(flags.passwordEnv)
				if err != nil {
					os.Exit(export.ExitConnection)
					return err
				}
				opts.Password = pw
			}

			if err := export.ResolveConnection(&opts, requireUser); err != nil {
				return err
			}

			code, err := export.Run(opts)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(code)
			return nil
		},
	}

	cmd.Flags().StringVar(&flags.host, "host", "", "Host do banco (ou VIGIA_DB_HOST)")
	cmd.Flags().IntVar(&flags.port, "port", 0, "Porta (ou VIGIA_DB_PORT)")
	cmd.Flags().StringVar(&flags.user, "user", "", "Usuário (ou VIGIA_DB_USER)")
	cmd.Flags().StringVar(&flags.password, "password", "", "Senha (preferir --password-env)")
	cmd.Flags().StringVar(&flags.passwordEnv, "password-env", "", "Nome da variável de ambiente com a senha")
	cmd.Flags().StringVar(&flags.database, "database", "", "Base MySQL (padrão: sys)")
	cmd.Flags().StringVar(&flags.alias, "alias", "", "Alias da instância (ou VIGIA_HOST_ALIAS)")
	cmd.Flags().StringVarP(&flags.output, "output", "o", "", "Arquivo de saída (padrão: snapshot.json)")
	cmd.Flags().StringVar(&flags.format, "format", "v1", "Formato: v1 ou legacy")
	cmd.Flags().BoolVarP(&flags.quiet, "quiet", "q", false, "Suprime mensagens de sucesso")
	cmd.Flags().BoolVar(&flags.skipVersion, "skip-version", false, "Não executar query de versão do engine")

	_ = cmd.Flags().MarkHidden("password")

	return cmd
}
