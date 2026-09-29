package scripts_test

import (
	"testing"

	"github.com/vigia-io/vigia-io-exporter/internal/scripts"
)

var postgresqlCoreIDs = []string{
	"postgresql.health.version",
	"postgresql.storage.database_size",
	"postgresql.sessions.active",
	"postgresql.sessions.blocking",
	"postgresql.maintenance.replication",
	"postgresql.storage.table_row_count",
	"postgresql.performance.bloat",
}

func TestForEngineFiltersByPrefix(t *testing.T) {
	scripts.Init()

	ss := scripts.ForEngine("sqlserver")
	if len(ss) == 0 {
		t.Fatal("expected sqlserver scripts")
	}
	for id := range ss {
		if len(id) < 10 || id[:10] != "sqlserver." {
			t.Fatalf("unexpected id for sqlserver filter: %q", id)
		}
	}
	if _, ok := ss["sqlserver.storage.table_row_count"]; !ok {
		t.Fatal("missing sqlserver.storage.table_row_count")
	}
	if _, ok := ss["mysql.storage.table_row_count"]; ok {
		t.Fatal("mysql script leaked into sqlserver filter")
	}

	my := scripts.ForEngine("mysql")
	if _, ok := my["mysql.storage.table_row_count"]; !ok {
		t.Fatal("missing mysql.storage.table_row_count")
	}
	if _, ok := my["sqlserver.health.version"]; ok {
		t.Fatal("sqlserver script leaked into mysql filter")
	}

	az := scripts.ForEngine("azuresql")
	if _, ok := az["sqlserver.storage.table_row_count"]; !ok {
		t.Fatal("azuresql should use sqlserver scripts")
	}
	if _, ok := az["mysql.health.version"]; ok {
		t.Fatal("mysql script leaked into azuresql filter")
	}

	pg := scripts.ForEngine("postgresql")
	if len(pg) < 5 {
		t.Fatalf("postgresql ForEngine: got %d scripts, want ≥5", len(pg))
	}
	for _, id := range postgresqlCoreIDs {
		if _, ok := pg[id]; !ok {
			t.Fatalf("missing %s", id)
		}
	}
	if _, ok := pg["mysql.health.version"]; ok {
		t.Fatal("mysql script leaked into postgresql filter")
	}
}

func TestCatalogIncludesTableRowCount(t *testing.T) {
	scripts.Init()
	for _, id := range []string{
		"sqlserver.storage.table_row_count",
		"mysql.storage.table_row_count",
		"postgresql.storage.table_row_count",
	} {
		if _, ok := scripts.Scripts[id]; !ok {
			t.Fatalf("missing embedded script %q", id)
		}
	}
}

func TestForEnginePostgreSQLHasExpectedIDs(t *testing.T) {
	scripts.Init()
	pg := scripts.ForEngine("postgresql")
	if len(pg) < 5 {
		t.Fatalf("got %d postgresql metrics, want ≥5", len(pg))
	}
	for _, id := range postgresqlCoreIDs {
		sql, ok := pg[id]
		if !ok {
			t.Fatalf("missing expected id %q", id)
		}
		if sql == "" {
			t.Fatalf("empty SQL for %q", id)
		}
	}
}
