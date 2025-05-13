package user

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"kxmj.common/web"
	"kxmj.common/web/middleware"
	"kxmj.core.api/internal/dto"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
	"time"
)

// TelRegister 手机注册
// @Description USER
// @Tags USER
// @Summary 手机注册
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param request body model.TelRegisterReq true "JSON"
// @Success 200 {object} web.Response{data=model.TelRegisterResp} "请求成功"
// @Router	/user/tel-register [POST]
func (s *Service) TelRegister(ctx *gin.Context) {
	payload := &model.TelRegisterReq{}
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

	ok, err := business.CheckSms(ctx, &dto.CheckSmsParameter{
		TelNumber: payload.TelNumber,
		Type:      1,
		Code:      payload.SMSCode,
	})

	if err != nil {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.CheckSmsCodeFailed)
		return
	}

	if !ok {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.InvalidSmsCode)
		return
	}

	exist := redis_cache.GetCache().GetUserCache().TelCache().Exists(ctx, payload.TelNumber)
	if exist {
		web.RespFailed(ctx, codes.TelNumberExisted)
		return
	}

	device, err := business.GetDevice(ctx, middleware.GetDeviceId(ctx))
	if err != nil {
		log.Sugar().Errorf("GetDevice err:%v", err)
		web.RespFailed(ctx, codes.DeviceNotFound)
		return
	}

	userId, err := business.CreateUserId(ctx)
	if err != nil {
		log.Sugar().Errorf("CreateUserId err:%v", err)
		web.RespFailed(ctx, codes.CreateUserIdFailed)
		return
	}

	bundle, err := business.GetBundle(ctx, middleware.GetBundleId(ctx))
	if err != nil {
		log.Sugar().Errorf("GetBundle err:%v", err)
		web.RespFailed(ctx, codes.BundleNotNull)
		return
	}

	salt := business.RandString(10)
	now := uint32(time.Now().Unix())
	err = business.CreateUser(ctx, &dto.CreateUserParameter{
		Id:                utils.Snowflake.Generate().Int64(),
		UserId:            userId,
		Nickname:          business.GetNickname(userId),
		Gender:            0,
		AvatarAddr:        "",
		AvatarFrame:       1,
		RealName:          "",
		IdCard:            "",
		UserMod:           0,
		AccountType:       0,
		Vip:               0,
		DeviceId:          middleware.GetDeviceId(ctx),
		RegisterIp:        ctx.ClientIP(),
		RegisterType:      2,
		TelNumber:         payload.TelNumber,
		Status:            1,
		BindingAt:         now,
		LoginPassword:     business.GetPassword(payload.Password, salt),
		LoginPasswordSalt: salt,
		Remark:            "",
		BundleId:          middleware.GetBundleId(ctx),
		BundleChannel:     bundle.BundleChannel,
		Organic:           device.Organic,
		WechatOpenId:      "",
		TiktokId:          "",
		HuaweiId:          "",
		Diamond:           "100000000",
		Gold:              "0",
		GoldBean:          "0",
		TotalRecharge:     "0",
		RechargeTimes:     0,
		Head:              0,
		Body:              0,
		Weapon:            0,
		CreatedAt:         now,
		UpdatedAt:         now,
	})

	if err != nil {
		log.Sugar().Errorf("CreateUser err:%v", err)
		web.RespFailed(ctx, codes.CreateUserFailed)
		return
	}

	web.RespSuccess(ctx, model.TelRegisterResp{
		UserId: userId,
	})
}
