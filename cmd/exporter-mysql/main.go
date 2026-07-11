package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/vigia-io/vigia-io-exporter/cmd/exporter-mysql/scripts"
	"github.com/vigia-io/vigia-io-exporter/internal/data"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	scripts.InitScripts()
}

func main() {
	var host, user, password, database string = "", "", "", ""

	var port int

	println("Bem-vindo ao Vigia Exporter para MySQL")

	println("")

	tryGetDataWithDefault("Informe o servidor do banco (localhost):", "localhost", &host)

	tryGetDataWithDefault("Informe o servidor do banco (3306):", 3306, &port)

	user = tryGetData("Informe o usuário para conectar:")

	password = tryGetData("Informe a senha para conectar:")

	tryGetDataWithDefault("Informe a base padrão para conexão (sys):", "sys", &database)

	db, err := createProvider(host, user, password, port, database)

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

func createProvider(host string, user string, password string, port int, database string) (*sql.DB, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)

	ret, err := sql.Open("mysql", connection)

	if err != nil {
		return nil, err
	}

	err = ret.Ping()

	if err != nil {
		return nil, err
	}

	return ret, err
}
