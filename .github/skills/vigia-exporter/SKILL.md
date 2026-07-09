---
name: vigia-exporter
description: |
  Use esta skill para implementar ou estender o vigia-exporter (CLI Go público).
  Acione para: novos providers SGBD, formato snapshot.json, codegen de scripts,
  flags não-interativas, cross-compile e migração do legado DbaMonitor.
metadata:
  author: Vigia
  version: "0.1.0"
  repository: vigia-exporter
  visibility: public
---

# Vigia Exporter

## Visão geral

O **vigia-exporter** é o componente de coleta pública do Vigia (licença source-available — não MIT). Binário CLI read-only que produz `snapshot.json` para a plataforma Vigia.

Repositório público: [`github.com/vigia-io/vigia-exporter`](https://github.com/vigia-io/vigia-exporter). Licença: **Vigia Source Available v1.0** (ADR-0008). Uso em produção requer conta Vigia.

## Capacidades

- **Providers:** SQL Server, Azure SQL, MySQL/MariaDB (fase 1); PostgreSQL, Oracle (fase 2)
- **Modos:** interativo (stdin) e não-interativo (flags + env)
- **Saída:** `snapshot.json` schema v1.0.0
- **Build:** cross-compile linux/amd64, windows/amd64, darwin/arm64

## Pré-requisitos

- Go 1.22+
- Acesso ao repo `vigia-scripts` para codegen (build time)
- Schema em `vigia-spec/schemas/snapshot-v1.json`

## Quick Start

| Cenário | Referência |
|---------|------------|
| Adicionar novo provider | `references/add-database-provider.md` |
| Migrar output.json → snapshot.json | `references/migrate-snapshot-format.md` |
| Implementar CLI não-interativo | `references/non-interactive-cli.md` |
| Pipeline codegen scripts.go | `references/scripts-codegen-pipeline.md` |
| Licença e uso permitido | `references/licensing-and-usage.md` |

## Workflow para agentes

**Triggers:** "exporter", "vigia-export", "snapshot.json", "provider", "go-mssqldb", "coleta offline"

1. Identifique se a tarefa é novo provider, métrica ou formato de saída.
2. Consulte ADR-0003 (offline-first), ADR-0006 (snapshot schema) e ADR-0008 (licença).
3. Use **apenas** padrões das referências — não invente campos no snapshot.
4. Scripts SQL vêm de `vigia-scripts` — exporter só embute, não edita SQL inline.
5. Mantenha módulo Go como `github.com/vigia-io/vigia-exporter`.

## Regras

1. Nunca persistir credenciais em disco além do processo
2. Somente queries SELECT / DMV read-only
3. `go.mod` obrigatório com versões fixadas
4. Testes de integração com Testcontainers quando possível
5. Mensagens CLI em PT-BR; logs estruturados em inglês
6. **Nunca** documentar ou implementar fluxo de produção sem ingest Vigia — licença source-available

## Referências

| Arquivo | Tarefa |
|---------|--------|
| `add-database-provider.md` | Novo engine no exporter |
| `migrate-snapshot-format.md` | Compatibilidade DbaMonitor → Vigia |
| `non-interactive-cli.md` | Flags, env vars, CI |
| `scripts-codegen-pipeline.md` | Gerar scripts.go do catálogo |
| `licensing-and-usage.md` | Licença source-available e fluxo Vigia |
