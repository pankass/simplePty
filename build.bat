@echo off
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -v -o client.exe ./client/

set GOOS=linux
go build -v -o client ./client/
go build -v -o server ./server/
