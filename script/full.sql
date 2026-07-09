use master
go

create database DbaMonitor
go

use DbaMonitor
go

create table ScriptOffline
(
	ScriptOfflineId int identity(1,1) not null,
	ScriptProvider varchar(20) not null,
	ScriptName varchar(100) not null,
	Script varchar(max) not null,
	constraint PK_ScriptOfflineId primary key clustered (ScriptOfflineId)
)
go

insert into ScriptOffline (ScriptProvider, ScriptName, Script)
values 
('sql', 'Teste1', 'Select 1 AS Campo1, ''Sql'' AS Campo2, null as Campo3'),
('sql', 'Teste2', 'Select 2 AS Campo1, ''Sql'' AS Campo2, null as Campo3'),
('mysql', 'Teste1', 'Select 1 AS Campo1, ''MySql'' AS Campo2, null as Campo3'),
('mysql', 'Teste2', 'Select 2 AS Campo1, ''MySql'' AS Campo2, null as Campo3'),
('sql-azure', 'Teste1', 'Select 1 AS Campo1, ''Sql Azure'' AS Campo2, null as Campo3'),
('sql-azure', 'Teste2', 'Select 2 AS Campo1, ''Sql Azure'' AS Campo2, null as Campo3')