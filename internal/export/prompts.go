package export

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ResolveConnection preenche campos ausentes via env ou prompts interativos.
func ResolveConnection(opts *Options, requireUser bool) error {
	opts.Host = firstNonEmpty(opts.Host, os.Getenv("VIGIA_DB_HOST"))
	opts.User = firstNonEmpty(opts.User, os.Getenv("VIGIA_DB_USER"))
	opts.Password = firstNonEmpty(opts.Password, os.Getenv("VIGIA_DB_PASSWORD"))
	opts.HostAlias = firstNonEmpty(opts.HostAlias, os.Getenv("VIGIA_HOST_ALIAS"))
	opts.Output = firstNonEmpty(opts.Output, os.Getenv("VIGIA_OUTPUT"))

	if opts.Port == 0 {
		if p := os.Getenv("VIGIA_DB_PORT"); p != "" {
			if v, err := strconv.Atoi(p); err == nil {
				opts.Port = v
			}
		}
	}
	if opts.Database == "" {
		opts.Database = os.Getenv("VIGIA_DB_DATABASE")
	}

	if opts.Host != "" && (!requireUser || opts.User != "") && opts.Password != "" {
		return nil
	}
	if opts.Host != "" && !requireUser && opts.User == "" {
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	if opts.Host == "" {
		v, err := promptWithDefault(reader, "Informe o servidor do banco (localhost):", "localhost")
		if err != nil {
			return err
		}
		opts.Host = v
	}

	if opts.Port == 0 {
		defaultPort := "1433"
		if opts.Provider.Engine() == "mysql" {
			defaultPort = "3306"
		}
		v, err := promptWithDefault(reader, fmt.Sprintf("Informe a porta (%s):", defaultPort), defaultPort)
		if err != nil {
			return err
		}
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("porta inválida: %w", err)
		}
		opts.Port = p
	}

	if requireUser && opts.User == "" {
		v, err := promptRequired(reader, "Informe o usuário para conectar:")
		if err != nil {
			return err
		}
		opts.User = v
	} else if !requireUser && opts.User == "" {
		v, err := promptWithDefault(reader, "Informe o usuário (vazio = Windows Authentication):", "")
		if err != nil {
			return err
		}
		opts.User = v
	}

	if opts.User != "" && opts.Password == "" {
		v, err := promptRequired(reader, "Informe a senha para conectar:")
		if err != nil {
			return err
		}
		opts.Password = v
	}

	if opts.Provider.Engine() == "mysql" && opts.Database == "" {
		v, err := promptWithDefault(reader, "Informe a base padrão para conexão (sys):", "sys")
		if err != nil {
			return err
		}
		opts.Database = v
	}

	if opts.HostAlias == "" {
		v, err := promptWithDefault(reader, "Informe o alias da instância (default):", "default")
		if err != nil {
			return err
		}
		opts.HostAlias = v
	}

	if opts.Output == "" {
		opts.Output = "snapshot.json"
	}

	return nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func promptWithDefault(reader *bufio.Reader, label, def string) (string, error) {
	fmt.Fprintln(os.Stderr, label)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def, nil
	}
	return line, nil
}

func promptRequired(reader *bufio.Reader, label string) (string, error) {
	for {
		fmt.Fprintln(os.Stderr, label)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}
}

// PasswordFromEnv lê senha de variável nomeada por --password-env.
func PasswordFromEnv(envName string) (string, error) {
	if envName == "" {
		return "", nil
	}
	v := os.Getenv(envName)
	if v == "" {
		return "", fmt.Errorf("variável de ambiente %s não definida", envName)
	}
	return v, nil
}
