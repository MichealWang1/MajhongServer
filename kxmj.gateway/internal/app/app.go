package app

import (
	"github.com/go-redis/redis/v8"
	"kxmj.gateway/config"
	"kxmj.gateway/internal/server"
)

type App struct {
	gateway *server.Gateway
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(redis *redis.Client) {
	a.gateway = &server.Gateway{}
	a.gateway.Inner = server.NewInner(config.Default.Self, config.Default.EtcdEndpoints, redis, a.gateway)
	a.gateway.Outer = server.NewOuter(a.gateway, config.Default.Outer)
	a.gateway.Start()
}

func (a *App) Close() {
	a.gateway.Close()
}
