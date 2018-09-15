@ECHO OFF
setlocal

go get -u github.com/golang/dep/cmd/dep || goto :error
dep ensure -v || goto :error


exit

:error
exit /b %errorlevel%
