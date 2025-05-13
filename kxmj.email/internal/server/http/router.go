package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kxmj.common/log"
	"kxmj.common/web/middleware"
	"kxmj.email/config"
	_ "kxmj.email/docs"
	"kxmj.email/internal/server/http/email"
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

	apiV1 := router.Group("/email/v1")
	{
		if config.Default.UseSwagger == "yes" {
			apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
		apiV1.Use(middleware.Auth())
		{
			apiV1.GET("email-list", email.Controller.GetUserAllMails)

			apiV1.POST("set-email-read", email.Controller.SetUserMailRead)
			apiV1.POST("take-email-item", email.Controller.TakeMailItem)
			apiV1.POST("del-email", email.Controller.DelUserMail)
		}
	}

	return router
}
