setlocal

set GOARCH=amd64

set GOOS=linux
call go build

set GOOS=windows
call go build