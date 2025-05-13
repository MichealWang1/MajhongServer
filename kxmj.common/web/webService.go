package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kxmj.common/log"
	"net/http"
	"time"
)

type HttpServiceConfig struct {
	Port int `yaml:"port"`
}

func StartHttpServer(router *gin.Engine, port int) {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Sugar().Debug(fmt.Sprintf("http listen * %d", port))

	err := s.ListenAndServe()
	if err != nil {
		log.Sugar().Error("StartHttpServer", zap.Any("err", err))
		panic(err)
	}
}

func GetUserId(ctx *gin.Context) int {
	val := ctx.MustGet(UserIdKey)
	return val.(int)
}
