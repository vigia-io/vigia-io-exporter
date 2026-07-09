package scripts

//Scripts Contém os scripts a serem rodados
var Scripts map[string]string

//InitScripts inicializa os scripts
func InitScripts() {
	Scripts = make(map[string]string)
	Scripts["Teste1"] = `Select 1 AS Campo1, 'Sql' AS Campo2, null as Campo3`
	Scripts["Teste2"] = `Select 2 AS Campo1, 'Sql' AS Campo2, null as Campo3`
}
