package user

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.core.api/internal/db"
	"kxmj.core.api/internal/dto"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
)

// ChangePassword 修改密码
// @Description USER
// @Tags USER
// @Summary 修改密码
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param request body model.ChangePasswordReq true "JSON"
// @Success 200 {object} web.Response{data=model.ChangePasswordResp} "请求成功"
// @Router	/user/change-password [POST]
func (s *Service) ChangePassword(ctx *gin.Context) {
	//解析请求参数
	payload := &model.ChangePasswordReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	//验证电话号码的合法性
	valid, err := business.PhoneRegex("86", payload.TelNumber)
	if err != nil || !valid {
		log.Sugar().Errorf("PhoneRegexerr:%v", err)
		web.RespFailed(ctx, codes.InvalidTelNumber)
		return
	}

	//验证手机验证码是否正确
	ok, err := business.CheckSms(ctx, &dto.CheckSmsParameter{
		TelNumber: payload.TelNumber,
		Type:      3,
		Code:      payload.SMSCode})

	if err != nil {
		log.Sugar().Errorf("PhoneRegexerr:%v", err)
		web.RespFailed(ctx, codes.CheckSmsCodeFailed)
		return
	}

	if !ok {
		log.Sugar().Errorf("PhoneRegexerr:%v", err)
		web.RespFailed(ctx, codes.InvalidSmsCode)
		return
	}

	//检查手机号码是否已存在
	userId, err := redis_cache.GetCache().GetUserCache().TelCache().Get(ctx, payload.TelNumber)
	if err != nil {
		web.RespFailed(ctx, codes.UserNotExist)
		return
	}

	//根据手机号码获取用户信息
	user, err := redis_cache.GetCache().GetUserCache().DetailCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.UserNotExist)
		return
	}

	newPassword := business.GetPassword(payload.NewPassword, user.LoginPasswordSalt)
	//更新密码
	err = db.ChangePassword(ctx, userId, newPassword)
	if err != nil {
		web.RespFailed(ctx, codes.PasswordError)
		return
	}

	//返回成功响应
	web.RespSuccess(ctx, model.ChangePasswordResp{
		Message: "Password changed successfully",
	})

}
