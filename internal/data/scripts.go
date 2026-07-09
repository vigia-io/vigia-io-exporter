package data

import "database/sql"

//GetScripts retorna os scripts salvos no banco
func GetScripts(db *sql.DB) (ret []Script, err error) {
	ret = make([]Script, 0)

	query := "select ScriptProvider, ScriptName, Script from ScriptOffline"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		item := Script{}

		err = rows.Scan(&item.Provider, &item.Name, &item.Script)

		if err != nil {
			return nil, err
		}

		ret = append(ret, item)
	}

	return
}

//Script é a estrutura de um script
type Script struct {
	Provider string
	Name     string
	Script   string
}
