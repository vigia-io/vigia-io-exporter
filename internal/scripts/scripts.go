// Code synced from vigia-io-scripts (output/scripts.generated.go).
// Source: PR #2 — includes *.storage.table_row_count (M2-01).

package scripts

import "strings"

// Scripts maps metric id to read-only SQL query text.
var Scripts = map[string]string{
	"mysql.config.buffer_pool":            "SELECT CAST(VARIABLE_VALUE AS UNSIGNED) AS innodb_buffer_pool_size\nFROM performance_schema.global_variables\nWHERE VARIABLE_NAME = 'innodb_buffer_pool_size';",
	"mysql.health.uptime":                 "SELECT CAST(VARIABLE_VALUE AS UNSIGNED) AS uptime_seconds\nFROM performance_schema.global_status\nWHERE VARIABLE_NAME = 'Uptime';",
	"mysql.health.version":                "SELECT VERSION() AS version_string;",
	"mysql.maintenance.replication":       "SELECT\n    CHANNEL_NAME,\n    SERVICE_STATE,\n    LAST_ERROR_NUMBER,\n    LAST_ERROR_MESSAGE\nFROM performance_schema.replication_connection_status;",
	"mysql.performance.slow_indicator":    "SELECT CAST(VARIABLE_VALUE AS UNSIGNED) AS slow_queries\nFROM performance_schema.global_status\nWHERE VARIABLE_NAME = 'Slow_queries';",
	"mysql.sessions.active":               "SELECT CAST(VARIABLE_VALUE AS UNSIGNED) AS threads_connected\nFROM performance_schema.global_status\nWHERE VARIABLE_NAME = 'Threads_connected';",
	"mysql.sessions.locks":                "SELECT\n    ENGINE_LOCK_ID,\n    ENGINE_TRANSACTION_ID,\n    OBJECT_SCHEMA,\n    OBJECT_NAME,\n    LOCK_TYPE,\n    LOCK_MODE,\n    LOCK_STATUS\nFROM performance_schema.data_locks;",
	"mysql.storage.database_size":         "SELECT\n    table_schema AS schema_name,\n    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS size_mb\nFROM information_schema.tables\nWHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')\nGROUP BY table_schema\nORDER BY size_mb DESC;",
	"mysql.storage.table_row_count":       "SELECT\n    TABLE_SCHEMA AS `database`,\n    TABLE_SCHEMA AS `schema`,\n    TABLE_NAME AS `table`,\n    CAST(TABLE_ROWS AS SIGNED) AS row_count\nFROM information_schema.TABLES\nWHERE TABLE_TYPE = 'BASE TABLE'\n  AND TABLE_SCHEMA NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')\nORDER BY row_count DESC;",
	"sqlserver.config.max_memory":         "SELECT\n    CAST(value_in_use AS BIGINT) AS max_server_memory_mb,\n    CAST(value AS BIGINT) AS configured_max_server_memory_mb\nFROM sys.configurations\nWHERE name = 'max server memory (MB)';",
	"sqlserver.health.uptime":             "SELECT\n    DATEDIFF(SECOND, sqlserver_start_time, GETUTCDATE()) AS uptime_seconds,\n    sqlserver_start_time AS start_time_utc\nFROM sys.dm_os_sys_info;",
	"sqlserver.health.version":            "SELECT\n    CAST(SERVERPROPERTY('ProductVersion') AS NVARCHAR(128)) AS product_version,\n    CAST(SERVERPROPERTY('ProductLevel') AS NVARCHAR(128)) AS product_level,\n    CAST(SERVERPROPERTY('Edition') AS NVARCHAR(128)) AS edition,\n    CAST(@@VERSION AS NVARCHAR(4000)) AS version_string;",
	"sqlserver.maintenance.agent_jobs":    "SELECT\n    j.name AS job_name,\n    h.run_date,\n    h.run_time,\n    h.run_status,\n    CAST(h.message AS NVARCHAR(4000)) AS message\nFROM msdb.dbo.sysjobs j\nINNER JOIN msdb.dbo.sysjobhistory h ON j.job_id = h.job_id\nWHERE h.run_status = 0\n  AND h.step_id = 0\n  AND h.run_date >= CONVERT(int, CONVERT(varchar(8), DATEADD(day, -1, GETDATE()), 112))\nORDER BY h.run_date DESC, h.run_time DESC;",
	"sqlserver.maintenance.backup_status": "SELECT\n    d.name AS database_name,\n    MAX(CASE WHEN b.type = 'D' THEN b.backup_finish_date END) AS last_full_backup,\n    MAX(CASE WHEN b.type = 'I' THEN b.backup_finish_date END) AS last_diff_backup,\n    MAX(CASE WHEN b.type = 'L' THEN b.backup_finish_date END) AS last_log_backup\nFROM sys.databases d\nLEFT JOIN msdb.dbo.backupset b ON d.name = b.database_name\nWHERE d.database_id > 4\nGROUP BY d.name;",
	"sqlserver.performance.wait_stats":    "SELECT TOP 10\n    wait_type,\n    wait_time_ms,\n    waiting_tasks_count\nFROM sys.dm_os_wait_stats\nWHERE wait_type NOT IN (\n    'CLR_SEMAPHORE', 'LAZYWRITER_SLEEP', 'RESOURCE_QUEUE', 'SLEEP_TASK',\n    'SLEEP_SYSTEMTASK', 'SQLTRACE_BUFFER_FLUSH', 'WAITFOR', 'LOGMGR_QUEUE',\n    'CHECKPOINT_QUEUE', 'REQUEST_FOR_DEADLOCK_SEARCH', 'XE_TIMER_EVENT',\n    'BROKER_TO_FLUSH', 'BROKER_TASK_STOP', 'CLR_MANUAL_EVENT', 'CLR_AUTO_EVENT',\n    'DISPATCHER_QUEUE_SEMAPHORE', 'FT_IFTS_SCHEDULER_IDLE_WAIT',\n    'XE_DISPATCHER_WAIT', 'XE_DISPATCHER_JOIN', 'BROKER_EVENTHANDLER',\n    'TRACEWRITE', 'FT_IFTSHC_MUTEX', 'SQLTRACE_INCREMENTAL_FLUSH_SLEEP',\n    'DIRTY_PAGE_POLL', 'HADR_FILESTREAM_IOMGR_IOCOMPLETION', 'SP_SERVER_DIAGNOSTICS_SLEEP',\n    'QDS_PERSIST_TASK_MAIN_LOOP_SLEEP', 'QDS_ASYNC_QUEUE'\n)\nORDER BY wait_time_ms DESC;",
	"sqlserver.sessions.active":           "SELECT\n    status,\n    COUNT(*) AS session_count\nFROM sys.dm_exec_sessions\nWHERE is_user_process = 1\nGROUP BY status;",
	"sqlserver.sessions.blocking":         "SELECT\n    blocking.session_id AS blocker_spid,\n    blocked.session_id AS blocked_spid,\n    blocked.wait_type,\n    blocked.wait_time AS wait_time_ms\nFROM sys.dm_exec_requests blocked\nINNER JOIN sys.dm_exec_sessions blocking\n    ON blocked.blocking_session_id = blocking.session_id\nWHERE blocked.blocking_session_id <> 0;",
	"sqlserver.storage.database_size":     "SELECT\n    DB_NAME(database_id) AS database_name,\n    CAST(SUM(CASE WHEN type = 0 THEN size ELSE 0 END) * 8.0 / 1024 AS DECIMAL(18, 2)) AS data_mb,\n    CAST(SUM(CASE WHEN type = 1 THEN size ELSE 0 END) * 8.0 / 1024 AS DECIMAL(18, 2)) AS log_mb,\n    CAST(SUM(size) * 8.0 / 1024 AS DECIMAL(18, 2)) AS total_mb\nFROM sys.master_files\nWHERE database_id > 4\nGROUP BY database_id\nORDER BY total_mb DESC;",
	"sqlserver.storage.table_row_count":   "SELECT\n    DB_NAME() AS [database],\n    s.name AS [schema],\n    t.name AS [table],\n    SUM(p.row_count) AS row_count\nFROM sys.dm_db_partition_stats AS p\nINNER JOIN sys.tables AS t\n    ON p.object_id = t.object_id\nINNER JOIN sys.schemas AS s\n    ON t.schema_id = s.schema_id\nWHERE p.index_id IN (0, 1)\n  AND t.is_ms_shipped = 0\nGROUP BY s.name, t.name\nORDER BY row_count DESC;",
	"sqlserver.storage.volume_free":       "SELECT DISTINCT\n    vs.volume_mount_point,\n    vs.logical_volume_name,\n    CAST(vs.total_bytes / 1048576.0 AS DECIMAL(18, 2)) AS total_mb,\n    CAST(vs.available_bytes / 1048576.0 AS DECIMAL(18, 2)) AS available_mb,\n    CAST(100.0 * vs.available_bytes / NULLIF(vs.total_bytes, 0) AS DECIMAL(5, 2)) AS pct_free\nFROM sys.master_files mf\nCROSS APPLY sys.dm_os_volume_stats(mf.database_id, mf.file_id) vs;",
}

// Init is kept for call sites; Scripts is already populated at package init.
func Init() {}

// scriptPrefix returns the catalog prefix for an engine.
// azuresql reuses sqlserver scripts (same T-SQL catalog).
func scriptPrefix(engine string) string {
	switch strings.ToLower(engine) {
	case "azuresql":
		return "sqlserver"
	case "mariadb":
		return "mysql"
	default:
		return strings.ToLower(engine)
	}
}

// ForEngine returns scripts whose id starts with the engine prefix.
func ForEngine(engine string) map[string]string {
	prefix := scriptPrefix(engine) + "."
	out := make(map[string]string)
	for id, sql := range Scripts {
		if strings.HasPrefix(id, prefix) {
			out[id] = sql
		}
	}
	return out
}
