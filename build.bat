set GOOS=windows
set GOARCH=amd64
go build -o build/mydaemon-v1.0-windows-amd64.exe

set GOOS=linux
set GOARCH=amd64
go build -o build/mydaemon-v1.0-linux-amd64

set GOOS=windows
set GOARCH=arm64
go build -o build/mydaemon-v1.0-windows-arm64.exe

set GOOS=linux
set GOARCH=arm64
go build -o build/mydaemon-v1.0-linux-arm64
