package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"net/http"
)

func Limit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.RequestURI
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		realUrl, _ := ctx.Get("RealURI")
		raw := ctx.Request.URL.RawQuery
		userAgent := ctx.Request.UserAgent()
		ip := ctx.ClientIP()
		header, _ := json.Marshal(ctx.Request.Header)

		var payload []byte
		if ctx.Request.Body != nil {
			var err error
			payload, err = io.ReadAll(ctx.Request.Body)
			if err == nil {
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(payload))
			}
		}

		key := utils.Md5(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s", url, realUrl, path, method, raw, userAgent, ip, string(header), string(payload)))

		// 相同业务在同一个生命周期以内只能请求一次
		if redis_cache.GetCache().GetLimitCache().ExistRequestKey(ctx, key) {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// 设置当前接口请求唯一key
		err := redis_cache.GetCache().GetLimitCache().SetHttpRequestKey(ctx, key)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("SetHttpRequestKey err:%v", err))
		}

		// 捕获全局异常，当发生了错误，立即移除掉分布式缓存限制
		defer func() {
			err := recover()
			if err != nil {
				redis_cache.GetCache().GetLimitCache().DelHttpRequestKey(ctx, key)
				panic(err)
			}
		}()

		ctx.Next()

		// 请求成功移除接口限制Key
		redis_cache.GetCache().GetLimitCache().DelHttpRequestKey(ctx, key)
	}
}
