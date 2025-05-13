package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoMethod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"msg": ctx.Request.RequestURI + " no method"})
	}
}

func NoRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": ctx.Request.RequestURI + " not found"})
	}
}
