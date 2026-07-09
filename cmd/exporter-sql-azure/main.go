package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/vigia-io/vigia-exporter/cmd/exporter-sql-azure/scripts"
	"github.com/vigia-io/vigia-exporter/internal/data"

	_ "github.com/microsoft/go-mssqldb"
)

func init() {
	scripts.InitScripts()
}

func main() {
	var host, user, password string

	println("Bem-vindo ao Vigia Exporter para Azure SQL")

	println("")

	println("Informe o servidor do banco (localhost):")

	n, err := fmt.Scanln(&host)

	if n == 0 {
		host = "localhost"
	}

	user = tryGetData("Informe o usuário para conectar:")

	password = tryGetData("Informe a senha para conectar:")

	db, err := createProvider(host, user, password)

	if err != nil {

		log.Fatalln("Falha ao conectar no banco. Erro: ", err.Error())
	}

	defer db.Close()

	result, err := data.GetData(db, scripts.Scripts)

	if err != nil {
		log.Fatalln("Falha ao buscar os dados. Erro: ", err.Error())
	}

	text, _ := json.Marshal(result)

	err = ioutil.WriteFile("output.json", text, 0644)

	if err != nil {
		log.Fatalln("Falha ao gerar o arquivo de saida. Erro: ", err.Error())
	}

	log.Println("Arquivo output.json gerado com sucesso.")
}

func tryGetData(text string) (ret string) {
	n := 0

	for n == 0 {
		println(text)

		n, _ = fmt.Scanln(&ret)
	}

	return
}

func createProvider(host string, user string, password string) (*sql.DB, error) {
	connection := fmt.Sprintf("server=%s;user id=%s;password=%s;", host, user, password)

	ret, err := sql.Open("sqlserver", connection)

	if err != nil {
		return nil, err
	}

	err = ret.Ping()

	if err != nil {
		return nil, err
	}

	return ret, err
}
