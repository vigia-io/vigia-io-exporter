# Vigia Exporter

[![License: Vigia Source Available](https://img.shields.io/badge/License-Vigia%20Source-teal.svg)](LICENSE)

**Coleta offline do [Vigia](https://github.com/vigia-io/vigia-io-platform)** — observabilidade para bancos relacionais.

CLI em Go que conecta em SQL Server, MySQL, Azure SQL (e em breve PostgreSQL e Oracle), executa scripts read-only e gera `snapshot.json` para ingest na **plataforma Vigia**.

> Repositório **público** para auditoria do código — **não** é software livre (MIT). Uso em produção requer [conta Vigia](https://getvigia.com) ou licença comercial. Ver [LICENSE](./LICENSE) e [resumo em PT](./docs/legal/exporter-license.md).

---

## Ecossistema

| Repositório | Visibilidade | Função |
|-------------|--------------|--------|
| **vigia-io-exporter** (este) | Público (source-available) | Coleta offline |
| [vigia-io-platform](https://github.com/vigia-io/vigia-io-platform) | Privado | API SaaS, alertas, billing |
| vigia-io-console | Privado (futuro) | Dashboard |
| vigia-io-scripts | Privado (futuro) | Catálogo SQL |

---

## Quick start

```bash
go build -o vigia-export ./cmd/exporter-sql
./vigia-export
```

Configure sua API key Vigia para upload:

```bash
export VIGIA_API_KEY=vigia_sk_live_...
# ver documentação de ingest em getvigia.com/docs
```

### Providers

| Binário | Engine |
|---------|--------|
| `cmd/exporter-sql` | SQL Server |
| `cmd/exporter-sql-azure` | Azure SQL |
| `cmd/exporter-mysql` | MySQL / MariaDB |

---

## Saída

`snapshot.json` conforme schema v1.0.0 — contratos de ingest em [`vigia-io-platform/docs/schemas/`](https://github.com/vigia-io/vigia-io-platform/tree/main/docs/schemas).

---

## Licença

**Vigia Source Available License v1.0** — ver [LICENSE](./LICENSE).

Uso em produção = ingest autorizada na plataforma Vigia. Avaliação limitada em homologação conforme LICENSE.

**Vigia** — *Enxergue seus bancos antes que o incidente chegue.*
