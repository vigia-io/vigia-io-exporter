package scripts

// MVP embutido — substituído por codegen de vigia-io-scripts (bloco M2-00).
var Scripts map[string]string

// Init carrega scripts read-only embutidos.
func Init() {
	Scripts = map[string]string{
		"sqlserver.health.version": `SELECT 1 AS Campo1, 'Sql' AS Campo2, NULL AS Campo3`,
		"sqlserver.sessions.blocking": `SELECT 2 AS Campo1, 'Sql' AS Campo2, NULL AS Campo3`,
	}
}
