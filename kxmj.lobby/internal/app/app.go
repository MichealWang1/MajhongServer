package app

import (
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.lobby/config"
	"kxmj.lobby/internal/server"
)

type App struct {
	server *server.LobbyServer
}

func NewApp() *App {
	return &App{
		server: server.NewLobbyServer(config.Default.Self, config.Default.EtcdEndpoints, redis_core.Default()),
	}
}

func (a *App) Start(redis *redis.Client) {
	// 实例化redis缓存服务
	redis_cache.InitReadCache(redis)

	a.server.Start()
}

func (a *App) Close() {
	a.server.Close()
}
