package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"

	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (blw bodyLogWriter) Write(b []byte) (int, error) {
	blw.body.Write(b)
	return blw.ResponseWriter.Write(b)
}

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.RequestURI
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		raw := ctx.Request.URL.RawQuery
		realUrl := ctx.GetString("RealURI")

		ip := ctx.ClientIP()
		headers, _ := json.Marshal(ctx.Request.Header)

		var payload []byte
		if ctx.Request.Body != nil {
			var err error
			payload, err = io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Error(
					url,
					zap.String("realUrl", realUrl),
					zap.String("path", path),
					zap.String("method", method),
					zap.Error(err),
				)
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(payload))
		}

		blw := &bodyLogWriter{
			ResponseWriter: ctx.Writer,
			body:           &bytes.Buffer{},
		}
		ctx.Writer = blw

		ctx.Next()

		end := time.Now()
		latencyTime := end.Sub(start)
		status := ctx.Writer.Status()
		userId := fmt.Sprintf("%d", GetUserId(ctx))

		if latencyTime > time.Millisecond*500 {
			log.Warn("api elapsed", zap.String("realUrl", realUrl), zap.Duration("elapsed", latencyTime))
		}

		log.Info(
			url,
			zap.Int("status", status),
			zap.Duration("latency", latencyTime),
			zap.String("ip", ip),
			zap.String("x-forwarded-for", ctx.Request.Header.Get("X-Forwarded-For")),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("realUrl", realUrl),
			zap.String("query", raw),
			zap.String("headers", string(headers)),
			zap.String("userId", userId),
			zap.ByteString("payload", payload),
			zap.ByteString("response", blw.body.Bytes()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("error", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
