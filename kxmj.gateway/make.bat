echo off

echo "start build"

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\gateway\kxmj-gateway  main.go

scp ..\build\gateway\kxmj-gateway centos@192.168.0.64:/home/kxmj/server/gateway1
scp ..\build\gateway\kxmj-gateway centos@192.168.0.64:/home/kxmj/server/gateway2
scp ..\build\gateway\kxmj-gateway centos@192.168.0.64:/home/kxmj/server/gateway3

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"