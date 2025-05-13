echo off

echo "start build"

swag init --parseDependency --parseInternal

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\email\kxmj-email  main.go

scp ..\build\email\kxmj-email centos@192.168.0.64:/home/kxmj/server/email
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/email/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/email
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/email

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"