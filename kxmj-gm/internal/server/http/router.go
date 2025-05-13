package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kxmj-gm/config"
	_ "kxmj-gm/docs"
	"kxmj-gm/internal/server/http/game_manager"
	"kxmj.common/log"
	"kxmj.common/web/middleware"
)

func CreateRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.Logger(log.Default()))
	router.Use(middleware.Limit())
	router.Use(middleware.RealIP())
	router.Use(middleware.Header())
	router.Use(middleware.Cors())
	router.NoMethod(middleware.NoMethod())

	apiV1 := router.Group("/gm/v1")
	{
		if config.Default.UseSwagger == "yes" {
			apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
	}
	apiV1.Use(middleware.Auth())
	{
		apiV1.POST("gm-set-card-stack", game_manager.Controller.SetCardStack)
		apiV1.POST("gm-del-card-stack", game_manager.Controller.DelCardStack)
		apiV1.POST("gm-set-catch-card", game_manager.Controller.SetCatchCard)
		apiV1.POST("gm-set-pause-room", game_manager.Controller.SetPauseRoom)
		apiV1.POST("gm-set-resume-room", game_manager.Controller.SetResumeRoom)
		apiV1.POST("gm-set-dismiss-room", game_manager.Controller.SetDismissRoom)
		apiV1.POST("gm-set-match-player", game_manager.Controller.SetMatchPlayer)
	}
	return router
}
