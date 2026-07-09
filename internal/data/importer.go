package data

import (
	"database/sql"
)

//GetData busca as informações no banco em questão
func GetData(db *sql.DB, scripts map[string]string) (ret map[string][]map[string]interface{}, err error) {

	ret = make(map[string][]map[string]interface{})

	for key, value := range scripts {
		result := make([]map[string]interface{}, 0)

		rows, err := db.Query(value)

		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()

		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			vals := make([]interface{}, len(columns))
			row := make(map[string]interface{})

			for i, value := range columns {
				vals[i] = new(sql.RawBytes)
				row[value] = nil
			}

			err = rows.Scan(vals...)

			for i, value := range columns {
				point := vals[i].(*sql.RawBytes)

				row[value] = string(*point)

				*point = nil
			}

			if err != nil {
				return nil, err
			}

			result = append(result, row)
		}

		ret[key] = result
	}

	return
}
