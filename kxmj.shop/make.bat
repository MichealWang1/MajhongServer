echo off

echo "start build"

swag init --parseDependency --parseInternal

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\shop\kxmj-shop  main.go

scp ..\build\shop\kxmj-shop centos@192.168.0.64:/home/kxmj/server/shop
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/shop/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/shop
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/shop

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"