package main

import (
	"kxmj-gm/config"
	"kxmj-gm/internal/app"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/mq"
	"kxmj.common/redis_cache/redis_core"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

// @title gm 服务接口
// @version 1.0
// @description gm 服务接口

// @BasePath /gm/v1

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

// @schemes http https

func main() {
	//初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 初始化配置文件
	config.Default = config.Create()

	// 初始化日志文件
	log.Init(config.Default.Logger)

	// 初始化redis服务
	redis_core.Init(config.Default.Redis)

	// 初始化mq服务
	mq.Init(config.Default.MqConfig)

	// 初始化错误代码
	codes.AddAllMessages()

	// 实例化app服务
	a := app.NewApp()

	// 启动服务
	log.Sugar().Info("gm 服务器启动...")
	a.Start(redis_core.Default())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务
	log.Sugar().Info("gm 服务器关闭...")
	a.Close()
}
