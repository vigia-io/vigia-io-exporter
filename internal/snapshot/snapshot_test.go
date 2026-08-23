package snapshot_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vigia-io/vigia-io-exporter/internal/snapshot"
)

func TestBuilderProducesV1Document(t *testing.T) {
	scripts := map[string][]map[string]interface{}{
		"sqlserver.health.version": {
			{"edition": "Enterprise", "version": "15.0"},
		},
	}

	doc := snapshot.NewBuilder("vigia-exporter", "0.1.0", "sqlserver", "prod-erp").
		Scripts(scripts).
		Metrics([]snapshot.MetricPoint{}).
		EngineVersion("Microsoft SQL Server 2019").
		Build()

	if doc.SchemaVersion != snapshot.SchemaVersion {
		t.Fatalf("schema_version: got %q want %q", doc.SchemaVersion, snapshot.SchemaVersion)
	}
	if doc.Collector.Name != "vigia-exporter" {
		t.Fatalf("collector.name: %q", doc.Collector.Name)
	}
	if doc.Target.Engine != "sqlserver" || doc.Target.HostAlias != "prod-erp" {
		t.Fatalf("target: %+v", doc.Target)
	}
	if doc.Metrics == nil {
		t.Fatal("metrics must not be nil")
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"schema_version", "collector", "target", "collected_at", "scripts", "metrics"} {
		if _, ok := parsed[key]; !ok {
			t.Fatalf("missing key %q in JSON", key)
		}
	}

	if _, err := time.Parse(time.RFC3339, parsed["collected_at"].(string)); err != nil {
		t.Fatalf("collected_at RFC3339: %v", err)
	}
}

func TestWriteRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/snapshot.json"

	doc := snapshot.NewBuilder("vigia-exporter", "0.1.0", "mysql", "lab").
		Scripts(map[string][]map[string]interface{}{"mysql.health.version": {}}).
		Metrics([]snapshot.MetricPoint{}).
		Build()

	if err := snapshot.Write(path, doc); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"schema_version": "1.0.0"`) {
		t.Fatalf("unexpected content: %s", raw)
	}
}
