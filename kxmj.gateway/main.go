package main

import (
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.gateway/config"
	"kxmj.gateway/internal/app"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	//初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 初始化配置文件
	config.Default = config.Create()

	// 初始化日志文件
	log.Init(config.Default.Logger)

	// 初始化redis服务
	redis_core.Init(config.Default.Redis)

	// 实例化redis缓存服务
	redis_cache.InitReadCache(redis_core.Default())

	// 实例化app服务
	a := app.NewApp()

	// 启动服务
	log.Sugar().Info("gateway服务器启动...")
	a.Start(redis_core.Default())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务
	log.Sugar().Info("gateway服务器关闭...")
	a.Close()
}
