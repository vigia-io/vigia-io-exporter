package data

import (
	"database/sql"
	"fmt"
)

// GetData executa scripts read-only e retorna linhas por script_id.
func GetData(db *sql.DB, scripts map[string]string) (map[string][]map[string]interface{}, error) {
	ret := make(map[string][]map[string]interface{})

	for key, query := range scripts {
		rows, err := db.Query(query)
		if err != nil {
			return nil, fmt.Errorf("script %s: %w", key, err)
		}

		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			return nil, err
		}

		result := make([]map[string]interface{}, 0)
		for rows.Next() {
			scanDest := make([]interface{}, len(columns))
			for i := range scanDest {
				scanDest[i] = new(sql.NullString)
			}

			if err := rows.Scan(scanDest...); err != nil {
				rows.Close()
				return nil, err
			}

			row := make(map[string]interface{}, len(columns))
			for i, col := range columns {
				ns := scanDest[i].(*sql.NullString)
				if ns.Valid {
					row[col] = ns.String
				} else {
					row[col] = nil
				}
			}
			result = append(result, row)
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()

		ret[key] = result
	}

	return ret, nil
}
