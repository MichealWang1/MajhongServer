echo off

echo "start build"

swag init --parseDependency --parseInternal

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\gm\kxmj-gm  main.go

scp ..\build\gm\kxmj-gm centos@192.168.0.64:/home/kxmj/server/gm
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/gm/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/gm
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/gm

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"