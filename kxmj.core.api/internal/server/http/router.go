package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kxmj.common/log"
	"kxmj.common/web/middleware"
	"kxmj.core.api/config"
	_ "kxmj.core.api/docs"
	"kxmj.core.api/internal/server/http/app"
	"kxmj.core.api/internal/server/http/game"
	"kxmj.core.api/internal/server/http/lobby"
	"kxmj.core.api/internal/server/http/recharge"
	"kxmj.core.api/internal/server/http/sms"
	"kxmj.core.api/internal/server/http/user"
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

	apiV1 := router.Group("/api/v1")
	{
		if config.Default.UseSwagger == "yes" {
			apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}

		userApi := apiV1.Group("user")
		{
			userApi.POST("tel-register", user.Controller.TelRegister)
			userApi.POST("login", user.Controller.Login)
			userApi.POST("bind-phone", user.Controller.BindPhoneNum)
			userApi.POST("change-password", user.Controller.ChangePassword)
			userApi.GET("get-info", middleware.Auth(), user.Controller.GetUserInformation)
			userApi.GET("token-login", middleware.Auth(), user.Controller.TokenLogin)
		}

		smsApi := apiV1.Group("sms")
		{
			smsApi.POST("send", sms.Controller.Send)
		}

		appApi := apiV1.Group("app")
		{
			appApi.POST("sync-device", app.Controller.SyncDevice)
			appApi.GET("get-gateways", middleware.Auth(), app.Controller.GetGateways)
			appApi.GET("get-base", middleware.Auth(), app.Controller.GetAppBaseInfo)
			appApi.GET("get-home", middleware.Auth(), app.Controller.GetHome)
			appApi.GET("get-items", app.Controller.GetItems)
		}

		apiV1.Use(middleware.Auth())
		gameApi := apiV1.Group("game")
		{
			gameApi.POST("room-list", game.Controller.GetRoomList)
		}

		lobbyApi := apiV1.Group("lobby")
		{
			lobbyApi.GET("get-wallet", lobby.Controller.GetWallet)
		}

		rechargeApi := apiV1.Group("recharge")
		{
			rechargeApi.GET("get-first-pack", recharge.Controller.GetFirstGiftPack)
			rechargeApi.GET("get-continue-pack", recharge.Controller.GetContinueGiftPack)
			rechargeApi.GET("take-continue-pack", recharge.Controller.TakeContinueGiftPack)
		}
	}

	return router
}
