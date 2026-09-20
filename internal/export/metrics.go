package export

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vigia-io/vigia-io-exporter/internal/snapshot"
)

const tableRowCountSuffix = ".storage.table_row_count"

// MaterializeTableRowCounts converts multi-row *.storage.table_row_count
// script results into one MetricPoint per table (labels: database/schema/table).
// Values come from catalog/DMV metadata only (P1 object names; no row content).
func MaterializeTableRowCounts(scriptRows map[string][]map[string]interface{}) []snapshot.MetricPoint {
	points := make([]snapshot.MetricPoint, 0)
	for id, rows := range scriptRows {
		if !strings.HasSuffix(id, tableRowCountSuffix) {
			continue
		}
		for _, row := range rows {
			if pt, ok := tableRowCountPoint(id, row); ok {
				points = append(points, pt)
			}
		}
	}
	return points
}

func tableRowCountPoint(id string, row map[string]interface{}) (snapshot.MetricPoint, bool) {
	value, ok := asFloat64(row["row_count"])
	if !ok {
		return snapshot.MetricPoint{}, false
	}

	labels := make(map[string]string, 3)
	for _, key := range []string{"database", "schema", "table"} {
		if s, ok := asString(row[key]); ok && s != "" {
			labels[key] = s
		}
	}

	return snapshot.MetricPoint{
		ID:     id,
		Value:  value,
		Unit:   "count",
		Labels: labels,
	}, true
}

func asString(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case fmt.Stringer:
		return t.String(), true
	default:
		return fmt.Sprint(t), true
	}
}

func asFloat64(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		s := fmt.Sprint(t)
		f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
}
