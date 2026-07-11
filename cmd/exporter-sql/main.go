package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/vigia-io/vigia-io-exporter/cmd/exporter-sql/scripts"
	"github.com/vigia-io/vigia-io-exporter/internal/data"

	_ "github.com/microsoft/go-mssqldb"
)

func init() {
	scripts.InitScripts()
}

func main() {
	var host, user, password string

	var port int

	println("Bem-vindo ao Vigia Exporter para SQL Server")

	println("")

	tryGetDataWithDefault("Informe o servidor do banco (localhost):", "localhost", &host)

	tryGetDataWithDefault("Informe o servidor do banco (1433):", 1433, &port)

	tryGetDataWithDefault("Informe o usuário para conectar (Windows Authentication):", "", &user)

	if len(user) > 0 {
		password = tryGetData("Informe a senha para conectar:")
	} else {
		password = ""
	}

	db, err := createProvider(host, user, password, port)

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

func tryGetDataWithDefault(text string, defaultValue interface{}, ret interface{}) {
	println(text)

	n, _ := fmt.Scanln(ret)

	if n == 0 {
		ret = defaultValue
	}

	return
}

func createProvider(host string, user string, password string, port int) (*sql.DB, error) {

	var connection string

	if len(user) > 0 {
		connection = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;", host, user, password, port)
	} else {
		connection = fmt.Sprintf("server=%s;trusted_connection=yes;port=%d;", host, port)
	}

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
