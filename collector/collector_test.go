package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProviderForEngine(t *testing.T) {
	provider, err := ProviderForEngine("sqlserver")
	if err != nil {
		t.Fatal(err)
	}
	if provider.Engine() != "sqlserver" {
		t.Fatalf("engine = %q", provider.Engine())
	}
}

func TestRunRequiresConnection(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "snapshot.json")
	code, err := Run(Options{
		Engine:    "sqlserver",
		Host:      "127.0.0.1",
		Port:      1,
		User:      "sa",
		Password:  "invalid",
		Output:    output,
		CollectorName: "vigia-agent",
		Version:   "test",
		Quiet:     true,
	})
	if err == nil {
		t.Fatal("expected connection error")
	}
	if code != ExitConnection {
		t.Fatalf("code = %d", code)
	}
	if _, statErr := os.Stat(output); statErr == nil {
		t.Fatal("output file should not exist")
	}
}
