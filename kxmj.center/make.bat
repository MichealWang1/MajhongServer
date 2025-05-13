echo off

echo "start build"

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\center\kxmj-center  main.go

scp ..\build\center\kxmj-center centos@192.168.0.64:/home/kxmj/server/center
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/center/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/center
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/center

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"