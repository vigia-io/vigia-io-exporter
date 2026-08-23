package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

const SchemaVersion = "1.0.0"

// CollectorMeta identifica o binário que produziu o snapshot.
type CollectorMeta struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Hostname string `json:"hostname,omitempty"`
}

// TargetMeta descreve a instância alvo da coleta.
type TargetMeta struct {
	Engine    string `json:"engine"`
	HostAlias string `json:"host_alias"`
	Version   string `json:"version,omitempty"`
}

// MetricPoint é uma métrica normalizada para ingest em séries.
type MetricPoint struct {
	ID       string            `json:"id"`
	Value    float64           `json:"value"`
	Unit     string            `json:"unit,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
	Severity string            `json:"severity,omitempty"`
}

// Document é o payload snapshot.json v1.
type Document struct {
	SchemaVersion string                              `json:"schema_version"`
	Collector     CollectorMeta                       `json:"collector"`
	Target        TargetMeta                          `json:"target"`
	CollectedAt   time.Time                           `json:"collected_at"`
	DurationMs    int64                               `json:"duration_ms,omitempty"`
	Scripts       map[string][]map[string]interface{} `json:"scripts"`
	Metrics       []MetricPoint                       `json:"metrics"`
	Signature     string                              `json:"signature,omitempty"`
}

// Builder monta um Document v1.
type Builder struct {
	collectorName    string
	collectorVersion string
	engine           string
	hostAlias        string
	engineVersion    string
	scripts          map[string][]map[string]interface{}
	metrics          []MetricPoint
	startedAt        time.Time
}

// NewBuilder inicia um builder com metadados do collector.
func NewBuilder(collectorName, collectorVersion, engine, hostAlias string) *Builder {
	return &Builder{
		collectorName:    collectorName,
		collectorVersion: collectorVersion,
		engine:           engine,
		hostAlias:        hostAlias,
		scripts:          make(map[string][]map[string]interface{}),
		metrics:          []MetricPoint{},
		startedAt:        time.Now().UTC(),
	}
}

// Scripts define o mapa de resultados por script_id.
func (b *Builder) Scripts(scripts map[string][]map[string]interface{}) *Builder {
	b.scripts = scripts
	return b
}

// EngineVersion define a versão detectada do SGBD.
func (b *Builder) EngineVersion(version string) *Builder {
	b.engineVersion = version
	return b
}

// Metrics define métricas normalizadas (pode ser slice vazia).
func (b *Builder) Metrics(metrics []MetricPoint) *Builder {
	if metrics == nil {
		b.metrics = []MetricPoint{}
		return b
	}
	b.metrics = metrics
	return b
}

// Build finaliza o documento.
func (b *Builder) Build() Document {
	hostname, _ := os.Hostname()
	collectedAt := time.Now().UTC()
	duration := collectedAt.Sub(b.startedAt).Milliseconds()

	return Document{
		SchemaVersion: SchemaVersion,
		Collector: CollectorMeta{
			Name:     b.collectorName,
			Version:  b.collectorVersion,
			Hostname: hostname,
		},
		Target: TargetMeta{
			Engine:    b.engine,
			HostAlias: b.hostAlias,
			Version:   b.engineVersion,
		},
		CollectedAt: collectedAt,
		DurationMs:  duration,
		Scripts:     b.scripts,
		Metrics:     b.metrics,
	}
}

// MarshalJSON serializa com RFC3339 em collected_at.
func (d Document) MarshalJSON() ([]byte, error) {
	type alias Document
	return json.Marshal(&struct {
		CollectedAt string `json:"collected_at"`
		*alias
	}{
		CollectedAt: d.CollectedAt.UTC().Format(time.RFC3339),
		alias:       (*alias)(&d),
	})
}

// WriteLegacy grava apenas o mapa scripts (formato output.json deprecado).
func WriteLegacy(path string, scripts map[string][]map[string]interface{}) error {
	data, err := json.MarshalIndent(scripts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Write grava snapshot.json v1.
func Write(path string, doc Document) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
