package snapshot_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/vigia-io/vigia-io-exporter/internal/snapshot"
)

func TestGeneratedSnapshotMatchesSchema(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")
	schemaPath := filepath.Join(root, "testdata", "snapshot-v1.json")

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatal(err)
	}

	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaBytes))
	if err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("snapshot-v1.json", schema); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile("snapshot-v1.json")
	if err != nil {
		t.Fatal(err)
	}

	doc := snapshot.NewBuilder("vigia-exporter", "0.1.0", "sqlserver", "prod-erp").
		Scripts(map[string][]map[string]interface{}{
			"sqlserver.health.version": {{"Campo1": "1", "Campo2": "Sql", "Campo3": nil}},
		}).
		Metrics([]snapshot.MetricPoint{}).
		EngineVersion("Microsoft SQL Server 2019").
		Build()

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	var instance any
	if err := json.Unmarshal(raw, &instance); err != nil {
		t.Fatal(err)
	}

	if err := compiled.Validate(instance); err != nil {
		t.Fatalf("schema validation failed: %v", err)
	}
}

func TestExampleFixtureMatchesSchema(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")

	schemaBytes, err := os.ReadFile(filepath.Join(root, "testdata", "snapshot-v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaBytes))
	if err != nil {
		t.Fatal(err)
	}
	exampleBytes, err := os.ReadFile(filepath.Join(root, "testdata", "snapshot-v1.example.json"))
	if err != nil {
		t.Fatal(err)
	}

	var example any
	if err := json.Unmarshal(exampleBytes, &example); err != nil {
		t.Fatal(err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("snapshot-v1.json", schema); err != nil {
		t.Fatal(err)
	}
	compiled, err := compiler.Compile("snapshot-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := compiled.Validate(example); err != nil {
		t.Fatalf("example fixture invalid: %v", err)
	}
}
