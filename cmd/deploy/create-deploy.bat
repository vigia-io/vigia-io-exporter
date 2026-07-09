ECHO OFF

mkdir build

cd /d build

mkdir linux
mkdir windows

cd /d ..

ECHO Gerando os arquivos

cd /d create-script

go run ./main.go

cd /d ../../..

ECHO Verificando as bibliotecas

go get

cd /d cmd/exporter-mysql

ECHO Compilando o exporter-mysql

ECHO Compilando versao para linux
set GOOS=linux&&set GOARCH=amd64&&go build

ECHO Movendo a versao do linux
move /y ./exporter-mysql ../deploy/build/linux/

ECHO Compilando versao para windows
set GOOS=windows&&set GOARCH=386&& go build

ECHO Movendo a versao do windows
move /y ./exporter-mysql.exe ../deploy/build/windows/


cd /d ../exporter-sql

ECHO Compilando o exporter-sql

ECHO Compilando versao para linux
set GOOS=linux&&set GOARCH=amd64&&go build

ECHO Movendo a versao do linux
move /y ./exporter-sql ../deploy/build/linux/

ECHO Compilando versao para windows
set GOOS=windows&&set GOARCH=386&& go build

ECHO Movendo a versao do windows
move /y ./exporter-sql.exe ../deploy/build/windows/

cd /d ../exporter-sql-azure

ECHO Compilando o exporter-sql-azure

ECHO Compilando versao para linux
set GOOS=linux&&set GOARCH=amd64&&go build

ECHO Movendo a versao do linux
move /y ./exporter-sql-azure ../deploy/build/linux/

ECHO Compilando versao para windows
set GOOS=windows&&set GOARCH=386&& go build

ECHO Movendo a versao do windows
move /y ./exporter-sql-azure.exe ../deploy/build/windows/
