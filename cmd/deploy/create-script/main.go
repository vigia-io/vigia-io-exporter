package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vigia-io/vigia-exporter/internal/data"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	db, err := createProvider("localhost", "usuarioteste", "1234", 1433)

	if err != nil {
		log.Fatalln("Falha ao conectar no banco. Erro: ", err.Error())
	}

	scripts, err := data.GetScripts(db)

	if err != nil {
		log.Fatalln("Falha ao buscar os dados. Erro: ", err.Error())
	}

	group := map[string][]data.Script{
		"sql":       make([]data.Script, 0),
		"mysql":     make([]data.Script, 0),
		"sql-azure": make([]data.Script, 0),
	}

	for _, item := range scripts {
		group[item.Provider] = append(group[item.Provider], item)
	}

	err = builMySQLScriptFile(group["mysql"])

	if err != nil {
		log.Fatalln("Falha ao montar o arquivo de script de Mysql. Erro: ", err.Error())
	}

	err = builSQLScriptFile(group["sql"])

	if err != nil {
		log.Fatalln("Falha ao montar o arquivo de script de Sql. Erro: ", err.Error())
	}

	err = builSQLAzureScriptFile(group["sql-azure"])

	if err != nil {
		log.Fatalln("Falha ao montar o arquivo de script de Sql Azure. Erro: ", err.Error())
	}

}

func createProvider(host string, user string, password string, port int) (*sql.DB, error) {

	var connection string

	if len(user) > 0 {
		connection = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", host, user, password, port, "DbaMonitor")
	} else {
		connection = fmt.Sprintf("server=%s;trusted_connection=yes;port=%d;database=%s;", host, port, "DbaMonitor")
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

func builSQLScriptFile(scripts []data.Script) (err error) {
	os.Remove("../../exporter-sql/scripts/scripts.go")
	file, err := os.Create("../../exporter-sql/scripts/scripts.go")

	if err != nil {
		return
	}

	defer file.Close()

	writeHeader(file)

	writScripts(file, scripts)

	writeFooter(file)

	return
}

func builMySQLScriptFile(scripts []data.Script) (err error) {
	os.Remove("../../exporter-mysql/scripts/scripts.go")
	file, err := os.Create("../../exporter-mysql/scripts/scripts.go")

	if err != nil {
		return
	}

	defer file.Close()

	writeHeader(file)

	writScripts(file, scripts)

	writeFooter(file)

	return
}

func builSQLAzureScriptFile(scripts []data.Script) (err error) {
	os.Remove("../../exporter-sql-azure/scripts/scripts.go")
	file, err := os.Create("../../exporter-sql-azure/scripts/scripts.go")

	if err != nil {
		return
	}

	defer file.Close()

	writeHeader(file)

	writScripts(file, scripts)

	writeFooter(file)

	return
}

func writeHeader(file *os.File) error {
	w := bufio.NewWriter(file)

	w.WriteString("package scripts\r\n\r\n")
	w.WriteString("//Scripts Contém os scripts a serem rodados\r\n")
	w.WriteString("var Scripts map[string]string\r\n\r\n")
	w.WriteString("//InitScripts inicializa os scripts\r\n")
	w.WriteString("func InitScripts() {\r\n")
	w.WriteString("\tScripts = make(map[string]string)\r\n")

	return w.Flush()
}

func writeFooter(file *os.File) error {
	w := bufio.NewWriter(file)

	w.WriteString("}\r\n")

	return w.Flush()
}

func writScripts(file *os.File, scripts []data.Script) error {
	w := bufio.NewWriter(file)

	for _, item := range scripts {
		w.WriteString("\tScripts[\"")
		w.WriteString(item.Name)
		w.WriteString("\"] = `")
		w.WriteString(replaceLiterals(item.Script))
		w.WriteString("`\r\n")
	}

	return w.Flush()
}

func replaceLiterals(value string) string {
	value = strings.Replace(value, "\"", "\\\"", -1)

	value = strings.Replace(value, "`", "", -1)

	return value
}
