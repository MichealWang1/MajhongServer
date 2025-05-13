package sms

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.core.api/internal/dto"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
)

// Send 发送短信
// @Description SMS
// @Tags SMS
// @Summary 发送短信
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param request body model.SendSmsReq true "JSON"
// @Success 200 {object} web.Response{data=model.SendSmsResp} "请求成功"
// @Router	/sms/send [POST]
func (s *Service) Send(ctx *gin.Context) {
	payload := &model.SendSmsReq{}
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

	if redis_cache.GetCache().GetSmsCache().Exist(ctx, payload.TelNumber, payload.Type) {
		web.RespFailed(ctx, codes.CanNotRepeatSendSms)
		return
	}

	result, err := business.SendSms(ctx, &dto.SendSmsParameter{
		TelNumber: payload.TelNumber,
		Type:      payload.Type,
	})

	if err != nil {
		log.Sugar().Errorf("SendSms err:%v", err)
		web.RespFailed(ctx, codes.SendSmsFailed)
		return
	}

	web.RespSuccess(ctx, &model.SendSmsResp{
		TTL: result.TTL,
	})
}
