package game_manager

import (
	"github.com/gin-gonic/gin"
	"kxmj-gm/internal/model"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/model/lobby"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/gm"
	"kxmj.common/web"
)

// checkPlayerPower 检测玩家是否有GM权限
func (s *Service) checkPlayerPower(ctx *gin.Context) bool {
	return true
	// 当前先屏蔽判断代码
	//userId := uint32(web.GetUserId(ctx))
	//reply := &center.GetUserInfoResp{}
	//err := s.XServer().CenterClient().Call(ctx, "GetUserInfo", &center.GetUserInfoReq{UserId: uint32(userId)}, reply)
	//if err != nil {
	//	return false
	//}
	//if reply.Data.Status < 3 {
	//	return false
	//}
	//return true
}

// SetCardStack 配置牌堆和庄家
// @Description GM
// @Tags GM
// @Summary 配置牌堆和庄家
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.CardStack true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-card-stack [POST]
func (s *Service) SetCardStack(ctx *gin.Context) {
	// 这里还要判断当前GM玩家操作权限判断 判断当前
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.CardStack{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("DoPairingCards ShouldBind PairingHandCards err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	if len(data.Cards) < 40 {
		log.Sugar().Infof("DoPairingCards len(data.Cards) < 40 ")
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	pairingData := &gm.CardStackData{
		Banker: data.Banker,
		Cards:  data.Cards,
	}
	err = redis_cache.GetCache().GetGameManageCache().Set(ctx, data.GameType, data.RoomLevel, data.UserId, pairingData)
	if err != nil {
		log.Sugar().Errorf("DoPairingCards GetGameManageCache().SetCardStack() err:%v", err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	web.RespSuccess(ctx, nil)
	return
}

// DelCardStack 删除配置牌堆
// @Description GM
// @Tags GM
// @Summary 删除配置牌堆
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.DeleteCardStack true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-del-card-stack [POST]
func (s *Service) DelCardStack(ctx *gin.Context) {
	// 这里还要判断当前GM玩家操作权限判断 判断当前
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.DeleteCardStack{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("DelCardStack ShouldBind PairingHandCards err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	err = redis_cache.GetCache().GetGameManageCache().Del(ctx, data.GameType, data.RoomType, data.UserId)
	if err != nil {
		log.Sugar().Infof(" DelCardStack data.userId:%d data.GameType:%d data.RoomType:%d err:%v ", data.UserId, data.GameType, data.RoomType, err)
	}
	return
}

// SetCatchCard 配下一张摸牌
// @Description GM
// @Tags GM
// @Summary 配下一张摸牌
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.CatchCard true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-catch-card [POST]
func (s *Service) SetCatchCard(ctx *gin.Context) {
	// 这里还要判断当前GM玩家操作权限判断 判断当前
	if s.checkPlayerPower(ctx) == false {
		return
	}

	data := &model.CatchCard{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Errorf("SetCatchCard ShouldBind CatchCard err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	if data.Card <= 0 {
		log.Sugar().Errorf("SetCatchCard data.Card:%d", data.Card)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	pairingData := &gm.CardStackData{
		CatchCard: data.Card,
	}
	err = redis_cache.GetCache().GetGameManageCache().Set(ctx, data.GameType, data.RoomLevel, data.UserId, pairingData)
	if err != nil {
		log.Sugar().Errorf("SetCatchCard GetGameManageCache().SetCatchCard() err:%v", err)
		web.RespFailed(ctx, codes.DbError)
		return
	}
	web.RespSuccess(ctx, nil)
	return
}

// SetPauseRoom 暂停房间
// @Description GM
// @Tags GM
// @Summary 暂停房间
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.UserInfo true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-pause-room [POST]
func (s *Service) SetPauseRoom(ctx *gin.Context) {
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.UserInfo{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("SetPauseRoom ShouldBind UserInfo err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	req := &lobby.PauseUserGameReq{
		UserId: data.UserId,
	}
	resp := &lobby.PauseUserGameResp{}
	s.XServer().LobbyClient().Call(ctx, "PauseUserGame", req, resp)
	web.RespSuccess(ctx, resp)
}

// SetResumeRoom 恢复房间(暂停恢复)
// @Description GM
// @Tags GM
// @Summary 恢复房间(暂停恢复)
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.UserInfo true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-resume-room [POST]
func (s *Service) SetResumeRoom(ctx *gin.Context) {
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.UserInfo{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("SetResumeRoom ShouldBind UserInfo err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	req := &lobby.ContinueGameReq{
		UserId: data.UserId,
	}
	resp := &lobby.ContinueGameResp{}
	s.XServer().LobbyClient().Call(ctx, "ContinueGame", req, resp)
	web.RespSuccess(ctx, resp)
}

// SetDismissRoom GM解散房间功能
// @Description GM
// @Tags GM
// @Summary GM解散房间功能
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.UserInfo true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-dismiss-room [POST]
func (s *Service) SetDismissRoom(ctx *gin.Context) {
	// 这里还要判断当前GM玩家操作权限判断 判断当前
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.UserInfo{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("SetResumeRoom ShouldBind UserInfo err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	// todo rpcx调用解散房间
	req := &lobby.PauseUserGameReq{
		UserId: data.UserId,
	}
	resp := &lobby.PauseUserGameResp{}
	s.XServer().LobbyClient().Call(ctx, "PauseUserGame", req, resp)
	web.RespSuccess(ctx, nil)
}

// SetMatchPlayer 调整匹配时 匹配的玩家
// @Description GM
// @Tags GM
// @Summary 调整匹配时 匹配的玩家
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.MatchPlayerType true "JSON"
// @Success 200 {object} web.Response{} "请求成功"
// @Router	/gm-set-match-player [POST]
func (s *Service) SetMatchPlayer(ctx *gin.Context) {
	if s.checkPlayerPower(ctx) == false {
		return
	}
	data := &model.MatchPlayerType{}
	err := ctx.ShouldBind(data)
	if err != nil {
		log.Sugar().Infof("SetResumeRoom ShouldBind UserInfo err:%v", err)
		web.RespFailed(ctx, codes.ParamError)
		return
	}
	s.XServer().LobbyClient()
	web.RespSuccess(ctx, nil)
}
