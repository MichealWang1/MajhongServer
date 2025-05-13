package app

import (
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.common/web"
	"kxmj.email/config"
	"kxmj.email/internal/server/http"
	"kxmj.email/internal/server/rpc"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(redis *redis.Client) {
	// 实例化redis缓存服务
	redis_cache.InitReadCache(redis)

	// 初始化rpcx服务
	rpc.Init(config.Default.Services, config.Default.EtcdEndpoints, redis_core.Default())

	// 启动rpcx服务
	rpc.Default().Start()

	// 创建API路由
	engine := http.CreateRouter()

	// 启动Http服务
	go web.StartHttpServer(engine, config.Default.HttpPort)
}

func (a *App) Close() {

}
