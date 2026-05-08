@echo off
set Version=1.2

set GOOS=windows
set GOARCH=amd64
go build -o build/mydaemon-v%Version%-windows-amd64.exe

set GOOS=linux
set GOARCH=amd64
go build -o build/mydaemon-v%Version%-linux-amd64

set GOOS=windows
set GOARCH=arm64
go build -o build/mydaemon-v%Version%-windows-arm64.exe

set GOOS=linux
set GOARCH=arm64
go build -o build/mydaemon-v%Version%-linux-arm64
