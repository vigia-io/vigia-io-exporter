package export_test

import (
	"testing"

	"github.com/vigia-io/vigia-io-exporter/internal/export"
)

func TestMaterializeTableRowCountsProducesNPoints(t *testing.T) {
	rows := map[string][]map[string]interface{}{
		"sqlserver.health.version": {
			{"product_version": "15.0"},
		},
		"sqlserver.storage.table_row_count": {
			{"database": "erp", "schema": "dbo", "table": "Orders", "row_count": "1500"},
			{"database": "erp", "schema": "dbo", "table": "OrderItems", "row_count": "4200"},
			{"database": "erp", "schema": "dbo", "table": "Customers", "row_count": "89"},
		},
	}

	points := export.MaterializeTableRowCounts(rows)
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3", len(points))
	}

	byTable := map[string]float64{}
	for _, p := range points {
		if p.ID != "sqlserver.storage.table_row_count" {
			t.Fatalf("unexpected id %q", p.ID)
		}
		if p.Unit != "count" {
			t.Fatalf("unit: got %q want count", p.Unit)
		}
		if p.Labels["database"] != "erp" || p.Labels["schema"] != "dbo" {
			t.Fatalf("labels: %+v", p.Labels)
		}
		tbl := p.Labels["table"]
		if tbl == "" {
			t.Fatal("missing table label")
		}
		byTable[tbl] = p.Value
	}

	want := map[string]float64{"Orders": 1500, "OrderItems": 4200, "Customers": 89}
	for tbl, val := range want {
		if byTable[tbl] != val {
			t.Fatalf("table %s: got %v want %v", tbl, byTable[tbl], val)
		}
	}
}

func TestMaterializeTableRowCountsMySQL(t *testing.T) {
	rows := map[string][]map[string]interface{}{
		"mysql.storage.table_row_count": {
			{"database": "app", "schema": "app", "table": "users", "row_count": "12"},
			{"database": "app", "schema": "app", "table": "sessions", "row_count": "34"},
		},
	}
	points := export.MaterializeTableRowCounts(rows)
	if len(points) != 2 {
		t.Fatalf("got %d points, want 2", len(points))
	}
	for _, p := range points {
		if p.ID != "mysql.storage.table_row_count" {
			t.Fatalf("id: %q", p.ID)
		}
	}
}

func TestMaterializeTableRowCountsSkipsMissingValue(t *testing.T) {
	rows := map[string][]map[string]interface{}{
		"sqlserver.storage.table_row_count": {
			{"database": "erp", "schema": "dbo", "table": "Bad", "row_count": nil},
			{"database": "erp", "schema": "dbo", "table": "Ok", "row_count": "1"},
		},
	}
	points := export.MaterializeTableRowCounts(rows)
	if len(points) != 1 {
		t.Fatalf("got %d points, want 1", len(points))
	}
	if points[0].Labels["table"] != "Ok" {
		t.Fatalf("unexpected point: %+v", points[0])
	}
}
