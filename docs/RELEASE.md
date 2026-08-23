# Release — vigia-export

## Goreleaser

Pré-requisito: [Goreleaser](https://goreleaser.com/install/) instalado.

```bash
# Snapshot local (sem publicar)
goreleaser release --snapshot --clean

# Release oficial (tag v* no GitHub)
git tag v0.1.0
goreleaser release --clean
```

### Arquiteturas (M1-01)

| OS | Arch |
|----|------|
| linux | amd64 |
| windows | amd64 |
| darwin | arm64 |

Binário: `vigia-export`

## Homologação manual

Com banco de homologação disponível:

```bash
go build -o vigia-export ./cmd/vigia-export

VIGIA_DB_PASSWORD='***' ./vigia-export sqlserver \
  --host sql.homolog.internal --port 1433 \
  --user vigia_reader --password-env VIGIA_DB_PASSWORD \
  --alias homolog -o /tmp/snapshot.json -q

go test ./internal/snapshot/...   # valida contra schema em testdata/
```

Não commitar `snapshot.json` com dados reais de cliente.
