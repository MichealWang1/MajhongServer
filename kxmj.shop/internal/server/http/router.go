package http

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kxmj.common/log"
	"kxmj.common/web/middleware"
	"kxmj.shop/config"
	_ "kxmj.shop/docs"
	"kxmj.shop/internal/server/http/goods"
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

	apiV1 := router.Group("/shop/v1")
	{
		if config.Default.UseSwagger == "yes" {
			apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}

		// 暂时测试先注释掉，正式开发时需要去掉注释，商城所有接口必须授权
		apiV1.Use(middleware.Auth())
		goodsApi := apiV1.Group("goods")
		{
			goodsApi.GET("goods-list", goods.Controller.GetGoodsList)
			goodsApi.POST("buy", goods.Controller.Buy)
		}
	}
	return router
}
