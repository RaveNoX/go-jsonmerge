@ECHO OFF
setlocal

set GOARCH=amd64

cd %~dp0
md artifacts

echo Linux
set GOOS=linux
call go build -o artifacts\jsonmerge .\cmd
if not %ERRORLEVEL% == 0 (exit %ERRORLEVEL%)

echo Windows
set GOOS=windows
call go build -o artifacts\jsonmerge.exe .\cmd
if not %ERRORLEVEL% == 0 (exit %ERRORLEVEL%)

echo Build done