package app

import (
	"github.com/go-redis/redis/v8"
	"kxmj.common/game_core"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.game.mjxlch/config"
	"kxmj.game.mjxlch/internal/game"
)

type App struct {
	server *game_core.Server
}

func NewApp() *App {
	return &App{
		server: game_core.NewServer(config.Default.Self, config.Default.EtcdEndpoints, config.Default.Lobby, config.Default.Center, redis_core.Default(), game.NewTemplate()),
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
