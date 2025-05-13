package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/lobby"
	"kxmj.common/redis_cache"
	"kxmj.common/web"
	"kxmj.common/web/middleware"
	"kxmj.core.api/internal/db"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
	"net/http"
	"time"
)

// GetGateways 获取网关地址
// @Description APP
// @Tags APP
// @Summary 获取网关地址
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetGatewayResp} "请求成功"
// @Router	/app/get-gateways [GET]
func (s *Service) GetGateways(ctx *gin.Context) {
	reply := &lobby.GetGatewayResp{}
	err := s.XServer().LobbyClient().Call(ctx, "GetGateways", &lobby.GetGatewayReq{}, reply)
	if err != nil {
		log.Sugar().Errorf("GetGateways err:%v", err)
		web.RespFailed(ctx, codes.ServerNetErr)
		return
	}

	resp := &model.GetGatewayResp{}
	for _, d := range reply.List {
		resp.List = append(resp.List, &model.GetGatewayInfo{
			SvrType: d.SvrType,
			SvrId:   d.SvrId,
			Addr:    d.Addr,
			Port:    d.Port,
		})
	}

	web.RespSuccess(ctx, resp)
}

// GetAppBaseInfo 获取APP基础信息
// @Description APP
// @Tags APP
// @Summary 获取APP基础信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetAppBaseInfoResp} "请求成功"
// @Router	/app/get-base [GET]
func (s *Service) GetAppBaseInfo(ctx *gin.Context) {
	config, err := db.GetAppConfig(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.DbError)
		return
	}

	reply := &model.GetAppBaseInfoResp{
		WechatSecretKey: config.WechatSecretKey,
		WechatAppId:     config.WechatAppId,
		HotRenewAddress: config.HotRenewAddress,
	}

	web.RespSuccess(ctx, reply)
}

// SyncDevice 同步设备信息
// @Description APP
// @Tags APP
// @Summary 同步设备信息
// @Param BundleId	header string true "dev.kxmj.com"
// @param request body model.SyncDeviceReq true "JSON"
// @Success 200 {object} web.Response{data=model.SyncDeviceResp} "请求成功"
// @Router	/app/sync-device [POST]
func (s *Service) SyncDevice(ctx *gin.Context) {
	payload := &model.SyncDeviceReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	// 00000000-0000-0000-0000-000000000000
	if payload.DeviceId == "00000000-0000-0000-0000-000000000000" || len(payload.DeviceId) == 0 {
		payload.DeviceId = uuid.NewString()
	}

	exist := business.ExistDevice(ctx, payload.DeviceId)
	if exist {
		web.RespSuccess(ctx, &model.SyncDeviceResp{
			DeviceId: payload.DeviceId,
		})
		return
	}

	err = business.AddDevice(ctx, &kxmj_core.Device{
		DeviceId:     payload.DeviceId,
		Os:           payload.OS,
		Brand:        payload.Brand,
		Version:      payload.Version,
		Model:        payload.Model,
		Width:        payload.Width,
		Height:       payload.Height,
		Manufacturer: payload.Manufacturer,
		AndroidSdk:   payload.AndroidInfo.SDK,
		AndroidId:    payload.AndroidInfo.ID,
		AndroidImei:  payload.AndroidInfo.IMEI,
		IosUuid:      payload.IosUUID,
		Organic:      payload.Organic,
		BundleId:     middleware.GetBundleId(ctx),
		CreatedAt:    uint32(time.Now().Unix()),
		UpdatedAt:    uint32(time.Now().Unix()),
	})

	if err != nil {
		web.RespFailed(ctx, codes.DbError)
		return
	}

	web.RespSuccess(ctx, &model.SyncDeviceResp{
		DeviceId: payload.DeviceId,
	})
}

// GetHome 获取首页信息
// @Description APP
// @Tags APP
// @Summary 获取首页信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetHomeResp} "请求成功"
// @Router	/app/get-home [GET]
func (s *Service) GetHome(ctx *gin.Context) {
	userId := uint32(web.GetUserId(ctx))
	user, err := redis_cache.GetCache().GetUserCache().DetailCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.GetUserInfoFail)
		return
	}

	wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.GetWalletInfoFailed)
		return
	}

	resp := &model.GetHomeResp{
		User: &model.HomeUser{
			Nickname:    user.Nickname,
			Gender:      user.Gender,
			AvatarAddr:  user.AvatarAddr,
			AvatarFrame: user.AvatarFrame,
			Diamond:     wallet.Diamond,
			Gold:        wallet.Gold,
			GoldBean:    wallet.GoldBean,
		},
		Guides: make(map[int]uint32, 0),
	}

	web.RespSuccess(ctx, resp)
}

// GetItems 获取物品列表
// @Description APP
// @Tags APP
// @Summary 获取物品列表
// @Success 200 {object} []model.ItemData "请求成功"
// @Router	/app/get-items [GET]
func (s *Service) GetItems(ctx *gin.Context) {
	items, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	var list []*model.ItemData
	for _, i := range items {
		var content []*item.GiftPackContent
		if len(i.Content) > 0 {
			_ = json.Unmarshal([]byte(i.Content), &content)
		}

		var extra map[uint32]uint32
		if len(i.Extra) > 0 {
			_ = json.Unmarshal([]byte(i.Extra), &extra)
		}

		list = append(list, &model.ItemData{
			ItemId:        i.ItemId,
			Name:          i.Name,
			ItemType:      i.ItemType,
			ServiceLife:   i.ServiceLife,
			Content:       content,
			Extra:         extra,
			GiftType:      i.GiftType,
			AdornmentType: i.AdornmentType,
		})
	}

	ctx.JSON(http.StatusOK, list)
}
