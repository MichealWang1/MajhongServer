package app

import (
	"kxmj.center/config"
	"kxmj.center/internal/server"
	"kxmj.common/redis_cache/redis_core"
)

type App struct {
	server *server.RpcxServer
}

func NewApp() *App {
	return &App{
		server: server.NewRpcxServer(config.Default.Self, config.Default.EtcdEndpoints, redis_core.Default()),
	}
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Close() {
	a.server.Close()
}
