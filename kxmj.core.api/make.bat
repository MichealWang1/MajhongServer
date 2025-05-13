echo off

echo "start build"

swag init --parseDependency --parseInternal

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\core-api\kxmj-core-api  main.go

scp ..\build\core-api\kxmj-core-api centos@192.168.0.64:/home/kxmj/server/core-api
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/core-api/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/core-api
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/core-api

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"