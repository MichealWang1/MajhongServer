echo off

echo "start build"

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build -o ..\build\game-mjxlch\kxmj-game-mjxlch  main.go

scp ..\build\game-mjxlch\kxmj-game-mjxlch centos@192.168.0.64:/home/kxmj/server/game-mjxlch
scp .\config\config.yaml centos@192.168.0.64:/home/kxmj/server/game-mjxlch/config.yaml
scp .\kill.sh centos@192.168.0.64:/home/kxmj/server/game-mjxlch
scp .\run.sh centos@192.168.0.64:/home/kxmj/server/game-mjxlch

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo "build success"