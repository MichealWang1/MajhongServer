package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		if len(origin) != 0 {
			context.Header("Access-Control-Allow-Origin", origin)
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			context.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, DeviceId, BundleId")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			context.Header("Access-Control-Allow-Credentials", "true")
			context.Header("Access-Control-Max-Age", "172800")
		}
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

const (
	DeviceIdKey = "DeviceId"
	BundleId    = "BundleId"
)

func Header() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(DeviceIdKey, ctx.Request.Header.Get(DeviceIdKey))
		ctx.Set(BundleId, ctx.Request.Header.Get(BundleId))
		ctx.Next()
	}
}

func GetDeviceId(ctx *gin.Context) string {
	return ctx.GetString(DeviceIdKey)
}

func GetBundleId(ctx *gin.Context) string {
	return ctx.GetString(BundleId)
}
