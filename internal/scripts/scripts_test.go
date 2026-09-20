package scripts_test

import (
	"testing"

	"github.com/vigia-io/vigia-io-exporter/internal/scripts"
)

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
	if _, ok := pg["postgresql.health.version"]; !ok {
		t.Fatal("missing postgresql.health.version")
	}
	if _, ok := pg["postgresql.storage.database_size"]; !ok {
		t.Fatal("missing postgresql.storage.database_size")
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
		"postgresql.health.version",
		"postgresql.storage.database_size",
	} {
		if _, ok := scripts.Scripts[id]; !ok {
			t.Fatalf("missing embedded script %q", id)
		}
	}
}
