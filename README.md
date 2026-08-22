# Vigia Exporter

[![License: Vigia Source Available](https://img.shields.io/badge/License-Vigia%20Source-teal.svg)](LICENSE)

**Coleta offline do [Vigia](https://github.com/vigia-io/vigia-io-spec)** — observabilidade para bancos relacionais.

CLI `vigia-export` em Go: conecta em SQL Server, MySQL/MariaDB e Azure SQL, executa scripts read-only e gera `snapshot.json` (schema v1) para ingest na plataforma Vigia.

> Repositório **público** para auditoria — **não** é MIT. Uso em produção requer [conta Vigia](https://getvigia.com). Ver [LICENSE](./LICENSE).

---

## Quick start

```bash
go build -o vigia-export ./cmd/vigia-export

# Não-interativo (cron / CI)
vigia-export sqlserver \
  --host localhost --port 1433 \
  --user vigia_reader --password-env VIGIA_DB_PASSWORD \
  --alias prod-erp -o snapshot.json -q

# Interativo (prompts em PT-BR)
vigia-export mysql
```

### Providers

| Comando | Engine |
|---------|--------|
| `vigia-export sqlserver` ou `sql` | SQL Server on-prem |
| `vigia-export mysql` | MySQL / MariaDB |
| `vigia-export azuresql` | Azure SQL |

### Flags e variáveis

| Flag | Env | Descrição |
|------|-----|-----------|
| `--host` | `VIGIA_DB_HOST` | Host do banco |
| `--port` | `VIGIA_DB_PORT` | Porta |
| `--user` | `VIGIA_DB_USER` | Usuário |
| `--password-env` | — | Nome da env com a senha (preferido) |
| `--alias` | `VIGIA_HOST_ALIAS` | Alias da instância |
| `-o` / `--output` | `VIGIA_OUTPUT` | Arquivo de saída (padrão: `snapshot.json`) |
| `--format` | — | `v1` (padrão) ou `legacy` (output.json deprecado) |
| `-q` / `--quiet` | — | Suprime mensagem de sucesso |

Precedência: flag → env → prompt interativo.

### Upload (cliente)

O exporter **não** faz upload HTTP. Exemplo com `curl`:

```bash
curl -sf -X POST https://api.getvigia.com/v1/snapshots \
  -H "Authorization: Bearer $VIGIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"instance_id\":\"$VIGIA_INSTANCE_ID\",\"payload\":$(cat snapshot.json)}"
```

---

## Saída

`snapshot.json` conforme [`vigia-io-spec/schemas/snapshot-v1.json`](https://github.com/vigia-io/vigia-io-spec/blob/main/schemas/snapshot-v1.json).

Validar localmente (com spec clonado):

```bash
go test ./...
```

---

## Release multi-arch

[Goreleaser](https://goreleaser.com/) — targets: `linux/amd64`, `windows/amd64`, `darwin/arm64`.

```bash
goreleaser release --snapshot --clean
```

Config: [`.goreleaser.yml`](./.goreleaser.yml).

---

## Legado

Binários antigos (`cmd/exporter-sql`, `exporter-mysql`, `exporter-sql-azure`) permanecem para compatibilidade DbaMonitor; use `vigia-export` para novos deploys.

---

## Licença

**Vigia Source Available License v1.0** — [LICENSE](./LICENSE).

**Vigia** — *Enxergue seus bancos antes que o incidente chegue.*
