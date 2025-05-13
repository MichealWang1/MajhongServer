package middleware

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"net/http"
)

type UnauthorizedContent struct {
	Status uint32 `json:"status"` // 状态 1=token not null; 2=token timeout; 3=invalid token;
	Msg    string `json:"msg"`    // 消息描述
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(web.Authorization)
		if len(token) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &UnauthorizedContent{
				Status: 1,
				Msg:    "token not null",
			})
			return
		}

		userId, err := redis_cache.GetCache().GetTokenCache().GetUserId(ctx, token)
		if err != nil {
			log.Sugar().Errorf("Get token:%s userId err:%v", token, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &UnauthorizedContent{
				Status: 2,
				Msg:    "token timeout",
			})
			return
		}

		if userId == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &UnauthorizedContent{
				Status: 2,
				Msg:    "invalid token",
			})
			return
		}

		// 用户登陆奖励统一在API处理，解决用户在线时拿不到登陆奖励的问题
		err = redis_cache.GetCache().GetTaskCache().LoginCache().Set(ctx, uint32(userId))
		if err != nil {
			log.Sugar().Errorf("Set user:%d login cache err:%v", userId, err)
		}

		ctx.Set(web.UserIdKey, userId)
		ctx.Set(web.Authorization, token)
		ctx.Next()
	}
}

func GetToken(ctx *gin.Context) string {
	return ctx.GetString(web.Authorization)
}

func GetUserId(ctx *gin.Context) uint32 {
	return uint32(ctx.GetInt(web.UserIdKey))
}
