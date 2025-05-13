package user

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.common/web/middleware"
	"kxmj.core.api/internal/db"
	"kxmj.core.api/internal/dto"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
)

// BindPhoneNum 绑定手机号
// @Description USER
// @Tags USER
// @Summary 绑定手机号
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param request body model.BindPhoneNumReq true "JSON"
// @Success 200 {object} web.Response{data=model.BindPhoneNumResp} "请求成功"
// @Router	/user/bind-phone [POST]
func (s *Service) BindPhoneNum(ctx *gin.Context) {
	//解析请求参数
	payload := &model.BindPhoneNumReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	//检验手机号码的合法性
	valid, err := business.PhoneRegex("86", payload.TelNumber)
	if err != nil || !valid {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.InvalidTelNumber)
		return
	}

	//检查验证码的合法性
	ok, err := business.CheckSms(ctx, &dto.CheckSmsParameter{
		TelNumber: payload.TelNumber,
		Type:      2,
		Code:      payload.SMSCode,
	})

	if err != nil {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.InvalidSmsCode)
		return
	}

	if !ok {
		log.Sugar().Errorf("PhoneRegex err:%v", err)
		web.RespFailed(ctx, codes.CheckSmsCodeFailed)
		return
	}

	//根据用户ID来寻找到目标用户
	userId, err := redis_cache.GetCache().GetUserCache().TelCache().Get(ctx, payload.TelNumber)
	if err != nil {
		log.Sugar().Errorf("GetUserByUserID err:%v", err)
		web.RespFailed(ctx, codes.DbError)
		return
	}

	if userId <= 0 {
		web.RespFailed(ctx, codes.TelNumberExisted)
		return
	}

	//插入手机号码
	err = db.BindPhoneNum(ctx, payload.TelNumber, userId)
	if err != nil {
		web.RespFailed(ctx, codes.TelNumberExisted)
		return
	}

	//返回成功响应
	web.RespSuccess(ctx, model.BindPhoneNumResp{
		Message: "phone number has already bound",
	})
}

// GetUserInformation 获取用户信息
// @Description USER
// @Tags USER
// @Summary 获取用户信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetInfoResp} "请求成功"
// @Router	/user/get-info [GET]
func (s *Service) GetUserInformation(ctx *gin.Context) {
	//获取用户Id
	userId := middleware.GetUserId(ctx)

	//建立一个user实体接受从Redis中获取的信息
	type users struct {
		UserId      uint32 `json:"userId"`
		Nickname    string `json:"nickname"`
		AvatarAddr  string `json:"avatarAddr"`
		Vip         uint8  `json:"vip"`
		AvatarFrame uint8  `json:"avatarFrame"`
		Diamond     string `json:"diamond"`
		Gold        string `json:"gold"`
		GoldBean    string `json:"goldBean"`
		Gender      uint8  `json:"gender"`
	}

	//根据用户ID获取用户信息
	user, err := redis_cache.GetCache().GetUserCache().DetailCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.UserNotExist)
		return
	}

	//将获得的数据映射到users结构体上
	var userdata users
	userdata.UserId = user.UserId
	userdata.Nickname = user.Nickname
	userdata.AvatarAddr = user.AvatarAddr
	userdata.AvatarFrame = user.AvatarFrame
	userdata.Vip = user.Vip
	userdata.Gender = user.Gender

	//根据用户ID获得钱包信息
	wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.UserNotExist)
		return
	}

	//将获得的数据映射到users结构体上
	userdata.Diamond = wallet.Diamond
	userdata.Gold = wallet.Gold
	userdata.GoldBean = wallet.GoldBean

	//返回成功响应
	web.RespSuccess(ctx, model.GetInfoResp{
		UserId:      userdata.UserId,
		Nickname:    userdata.Nickname,
		AvatarAddr:  userdata.AvatarAddr,
		AvatarFrame: userdata.AvatarFrame,
		Diamond:     userdata.Diamond,
		Gold:        userdata.Gold,
		GoldBean:    userdata.GoldBean,
		Gender:      userdata.Gender,
		Vip:         userdata.Vip,
	})
}
