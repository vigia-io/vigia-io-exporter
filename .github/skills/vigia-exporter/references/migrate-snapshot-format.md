# Migrar formato output.json → snapshot.json

## Contexto

O DbaMonitor legado produzia:

```json
{
  "Teste1": [{ "Campo1": "1", "Campo2": "Sql", "Campo3": null }],
  "Teste2": [{ "Campo1": "2", "Campo2": "Sql", "Campo3": null }]
}
```

O Vigia encapsula isso em `scripts` e adiciona metadados.

## Formato Vigia v1.0.0

```json
{
  "schema_version": "1.0.0",
  "collector": {
    "name": "vigia-exporter",
    "version": "0.1.0",
    "hostname": "bastion-01"
  },
  "target": {
    "engine": "sqlserver",
    "host_alias": "prod-erp",
    "version": "Microsoft SQL Server 2019"
  },
  "collected_at": "2026-07-06T17:30:00Z",
  "duration_ms": 1240,
  "scripts": {
    "sqlserver.health.version": [
      { "edition": "Enterprise", "version": "15.0.2000" }
    ]
  },
  "metrics": [
    {
      "id": "sqlserver.sessions.blocking",
      "value": 2,
      "unit": "count",
      "severity": "critical"
    }
  ]
}
```

## Implementação Go

```go
type Snapshot struct {
    SchemaVersion string            `json:"schema_version"`
    Collector     CollectorMeta     `json:"collector"`
    Target        TargetMeta        `json:"target"`
    CollectedAt   time.Time         `json:"collected_at"`
    DurationMs    int64             `json:"duration_ms,omitempty"`
    Scripts       map[string][]map[string]interface{} `json:"scripts"`
    Metrics       []MetricPoint     `json:"metrics,omitempty"`
}
```

## Compatibilidade na ingest

A plataforma aceita upload legado por 6 meses via adapter:

```go
func AdaptLegacy(raw map[string][]map[string]interface{}) Snapshot {
    return Snapshot{
        SchemaVersion: "1.0.0",
        Scripts:       raw,
        // metadados mínimos preenchidos server-side
    }
}
```

## Flag de transição

```bash
vigia-export sqlserver --format=v1        # default
vigia-export sqlserver --format=legacy    # output.json puro (deprecated)
```

Remover `--format=legacy` na v1.0.0 do exporter.
