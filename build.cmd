@echo off
setlocal

for /f %%i in ('git describe --tags --always --dirty') do set VER=%%i
for /f "tokens=2 delims= " %%i in ('date /t') do set DATE=%%i

echo %VER%
echo %DATE%

rem rsrc -manifest cryptopass.manifest -o rsrc.syso -arch="amd64"
go test -cover ./...
go build -ldflags "-X main.Version=%VER% -X main.Build=%DATE%"

