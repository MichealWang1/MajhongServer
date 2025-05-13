package user

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_logger"
	"kxmj.common/log"
	"kxmj.common/mq"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.common/web/middleware"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
	"time"
)

// Login 手机登录
// @Description USER
// @Tags USER
// @Summary 手机登录
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param request body model.LoginReq true "JSON"
// @Success 200 {object} web.Response{data=model.LoginResp} "请求成功"
// @Router	/user/login [POST]
func (s *Service) Login(ctx *gin.Context) {
	payload := &model.LoginReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	valid, err := business.PhoneRegex("86", payload.TelNumber)
	if err != nil || !valid {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.InvalidTelNumber)
		return
	}

	userId, err := redis_cache.GetCache().GetUserCache().TelCache().Get(ctx, payload.TelNumber)
	if err != nil || userId <= 0 {
		web.RespFailed(ctx, codes.UserNotExist)
		return
	}

	user, err := redis_cache.GetCache().GetUserCache().DetailCache().Get(ctx, userId)
	if err != nil {
		log.Sugar().Errorf("GetUserByTelNumber err:%v", err)
		web.RespFailed(ctx, codes.DbError)
		return
	}

	if business.GetPassword(payload.Password, user.LoginPasswordSalt) != user.LoginPassword {
		web.RespFailed(ctx, codes.PasswordError)
		return
	}

	token, err := redis_cache.GetCache().GetTokenCache().SetToken(ctx, int(userId))
	if err != nil {
		web.RespFailed(ctx, codes.DbError)
		return
	}

	// 写登陆日志
	err = mq.AddLogger(&kxmj_logger.UserLogin{
		UserId:    userId,
		Ip:        ctx.ClientIP(),
		DeviceId:  middleware.GetDeviceId(ctx),
		LoginType: 1,
		CreatedAt: uint32(time.Now().Unix()),
	})

	if err != nil {
		log.Sugar().Errorf("AddLogger err:%v", err)
	}

	web.RespSuccess(ctx, model.LoginResp{
		UserId: userId,
		Token:  token,
	})
}

// TokenLogin Token登录
// @Description USER
// @Tags USER
// @Summary Token登录
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.LoginResp} "请求成功"
// @Router	/user/token-login [GET]
func (s *Service) TokenLogin(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)
	token := middleware.GetToken(ctx)

	err := redis_cache.GetCache().GetTokenCache().Expired(ctx, token, int(userId))
	if err != nil {
		web.RespFailed(ctx, codes.DbError)
		log.Sugar().Errorf("Expired err:%v", err)
		return
	}

	web.RespSuccess(ctx, model.LoginResp{
		UserId: userId,
		Token:  token,
	})
}
