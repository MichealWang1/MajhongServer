package main

import (
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/mq"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.core.api/config"
	"kxmj.core.api/internal/app"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

// @title core-api 服务接口
// @version 1.0
// @description core-api 服务接口
// @BasePath /api/v1
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

	// 初始化核心业务主库
	mysql.InitCoreMaster(config.Default.CoreMaster)

	// 初始化核心业务从库
	mysql.InitCoreSlave(config.Default.CoreSlave)

	// 初始化日志主库
	mysql.InitLoggerMaster(config.Default.LoggerMaster)

	// 初始化日志从库
	mysql.InitLoggerSlave(config.Default.LoggerSlave)

	// 初始化业务报表主库(business owner)
	mysql.InitReportMaster(config.Default.ReportMaster)

	// 初始化业务报表从库(business owner)
	mysql.InitReportSlave(config.Default.ReportSlave)

	// 初始化游戏记录主库
	mysql.InitGameMaster(config.Default.GameMaster)

	// 初始化游戏记录从库
	mysql.InitGameSlave(config.Default.GameSlave)

	// 初始化redis服务
	redis_core.Init(config.Default.Redis)

	// 初始化mq服务
	mq.Init(config.Default.MqConfig)

	// 初始化错误代码
	codes.AddAllMessages()

	// 实例化app服务
	a := app.NewApp()

	// 启动服务
	log.Sugar().Info("core-api 服务器启动...")
	a.Start(redis_core.Default())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务
	log.Sugar().Info("core-api 服务器关闭...")
	a.Close()
}
