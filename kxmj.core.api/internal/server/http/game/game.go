package game

import (
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/model/lobby"
	"kxmj.common/web"
	"kxmj.core.api/internal/model"
)

// GetRoomList 获取房间列表信息
// @Description GAME
// @Tags GAME
// @Summary 获取房间列表信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.RoomListReq true "JSON"
// @Success 200 {object} web.Response{data=model.RoomListResp} "请求成功"
// @Router	/game/room-list [POST]
func (s *Service) GetRoomList(ctx *gin.Context) {
	payload := &model.RoomListReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	reply := &center.RoomConfigListResp{}
	err = s.XServer().CenterClient().Call(ctx, "GetRoomConfigList", &center.RoomConfigListReq{GameId: payload.GameId}, reply)
	if err != nil {
		log.Sugar().Errorf("GetRoomConfigList err:%v", err)
		web.RespFailed(ctx, codes.GetRoomConfigListFailed)
		return
	}

	if reply.Code != codes.Success {
		web.RespFailed(ctx, reply.Code, reply.Msg)
		return
	}

	onlineArgs := &lobby.GetRoomsOnlineReq{GameId: payload.GameId}
	onlineReply := &lobby.GetRoomsOnlineResp{}
	for _, d := range reply.Data {
		onlineArgs.RoomIds = append(onlineArgs.RoomIds, d.RoomId)
	}

	err = s.XServer().LobbyClient().Call(ctx, "GetRoomsOnline", onlineArgs, onlineReply)
	if err != nil {
		log.Sugar().Errorf("GetRoomsOnline err:%v", err)
	}

	resp := &model.RoomListResp{}
	for _, d := range reply.Data {
		info := &model.RoomInfo{
			RoomId:      d.RoomId,
			RoomType:    d.RoomType,
			GameId:      d.GameId,
			GameType:    d.GameType,
			RoomLevel:   d.RoomLevel,
			Tags:        d.Tags,
			Extra:       d.Extra,
			MinLimit:    d.MinLimit,
			MaxLimit:    d.MaxLimit,
			BaseScore:   d.BaseScore,
			MaxMultiple: d.MaxMultiple,
			Ticket:      d.Ticket,
			MatchTime:   d.MatchTime,
		}

		for _, o := range onlineReply.OnlineUsers {
			if d.RoomId == o.RoomId {
				info.CurPlayers = o.Users
				break
			}
		}

		resp.List = append(resp.List, info)
	}
	web.RespSuccess(ctx, resp)
}
