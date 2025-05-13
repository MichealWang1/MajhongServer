package game

import (
	"kxmj.common/codes"
	"kxmj.common/game_core"
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/net"
	"kxmj.common/proto/gateway_pb"
	"kxmj.common/proto/lobby_pb"
	"kxmj.game.mjxlch/pb"
)

func (d *Desk) handler(ctx net.MsgContext) {
	switch ctx.Request().MsgId {
	case uint16(lobby_pb.MID_ON_LINE):
		d.online(ctx)
	case uint16(lobby_pb.MID_OFF_LINE):
		d.offline(ctx)
	case uint16(gateway_pb.MID_MATCH):
		d.enter(ctx)
	case uint16(gateway_pb.MID_LEAVE_DESK):
		d.leave(ctx)
	case uint16(lobby_pb.MID_RISE_BUY_SUCCESS):
		d.riseSuccess(ctx)
	case uint16(gateway_pb.MID_GIVE_UP_RISE):
		d.giveUpRise(ctx)
	case uint16(pb.MID_DESK_INFO):
		d.onDeskInfoRequest(ctx)
	case uint16(pb.MID_PLAYER_READY):
		d.onPlayerReadyRequest(ctx)
	case uint16(pb.MID_SWAP_INFO):
		d.onSwapRequest(ctx)
	case uint16(pb.MID_CHOOSE_MISS_INFO):
		d.onChooseMissRequest(ctx)
	case uint16(pb.MID_OUT_CARD_INFO):
		d.onOutCardRequest(ctx)
	case uint16(pb.MID_ACTIONS_INFO):
		d.onOperateRequest(ctx)
	case uint16(pb.MID_FORCE_END):
		d.onForceEndRequest(ctx)
	}
}

func (d *Desk) online(ctx net.MsgContext) {
	payload := &lobby_pb.Online{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	p := d.getPlayer(payload.UserId)
	log.Sugar().Infof("player:", p, "online!!!")
	if p != nil {
		p.IsOnline = true
	}
}

func (d *Desk) offline(ctx net.MsgContext) {
	payload := &lobby_pb.Offline{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	p := d.getPlayer(payload.UserId)
	log.Sugar().Infof("player:", p, "offline!!!")
	if p != nil {
		p.IsOnline = false
	}
}

func (d *Desk) enter(ctx net.MsgContext) {
	payload := &gateway_pb.MatchReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		d.game.OnLeave(&game_core.LeaveDesk{
			UserId:  payload.UserId,
			Gold:    "0",
			SeatId:  0,
			IsRobot: false,
			Desk:    d,
		})
		d.game.SendErrMessage(ctx, codes.UnMarshalPbErr)
		return
	}

	user, err := d.game.Server().GetUserInfo(payload.UserId)
	if err != nil {
		log.Sugar().Errorf("GetUserInfo user:%d err:%v", payload.UserId, err)
		d.game.OnLeave(&game_core.LeaveDesk{
			UserId:  payload.UserId,
			Gold:    "0",
			SeatId:  0,
			IsRobot: false,
			Desk:    d,
		})
		d.game.SendErrMessage(ctx, codes.GetUserInfoFail)
		return
	}

	if user.Code != codes.Success {
		log.Sugar().Errorf("GetUserInfo user:%d err:%v", payload.UserId, err)
		d.game.OnLeave(&game_core.LeaveDesk{
			UserId:  payload.UserId,
			Gold:    "0",
			SeatId:  0,
			IsRobot: false,
			Desk:    d,
		})
		d.game.SendErrMessage(ctx, user.Code)
		return
	}

	gold, err := d.game.Server().CheckUserGold(payload.UserId)
	if err != nil {
		log.Sugar().Errorf("GetUserGold user:%d err:%v", payload.UserId, err)
		d.game.SendErrMessage(ctx, codes.GetUserGoldFail)
		return
	}

	if gold.Code != codes.Success {
		log.Sugar().Errorf("GetUserGold user:%d err:%v", payload.UserId, err)
		d.game.OnLeave(&game_core.LeaveDesk{
			UserId:  payload.UserId,
			Gold:    "0",
			SeatId:  0,
			IsRobot: false,
			Desk:    d,
		})
		d.game.SendErrMessage(ctx, gold.Code)

		return
	}

	player := d.getNullPlayerBySeat()
	player.UserId = user.Data.UserId
	player.Gold = gold.Data.Gold
	player.Nickname = user.Data.Nickname
	player.AvatarAddr = user.Data.AvatarAddr
	player.IsOnline = true
	player.GoldStatus = false
	player.RunTimeData = NewPlayerRunTimeData()

	// 通知大厅用户进入游戏
	d.game.OnEnter(&game_core.EnterDesk{
		Desk: d,
		Player: &game_core.PlayerInfo{
			UserId:     player.UserId,
			IsRobot:    player.IsRobot,
			SeatId:     player.SeatId,
			Gold:       player.Gold,
			Nickname:   player.Nickname,
			AvatarAddr: player.AvatarAddr,
			IconStyle:  0,
		},
	})
	d.broadcastEnterInfoNotify(player.UserId)
}

func (d *Desk) leave(ctx net.MsgContext) {
	// todo 离开游戏，业务实现
	player := d.getPlayer(ctx.Request().UserId)
	d.PlayerOnLeave(player.SeatId)
}

func (d *Desk) riseSuccess(ctx net.MsgContext) {
	// todo 复活卡购买成功，业务实现
}

func (d *Desk) giveUpRise(ctx net.MsgContext) {
	// todo 放弃复活，业务实现
}

// 处理玩家信息请求
func (d *Desk) onDeskInfoRequest(ctx net.MsgContext) {
	req := &pb.GameDeskInfoRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onDeskInfoRequest, deskId:%v, req:%v", d.roundId, req)
	d.broadcastDeskInfo(ctx.Request().UserId)
}

func (d *Desk) onPlayerReadyRequest(ctx net.MsgContext) {
	req := &pb.GamePlayerReadyRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onPlayerReadyRequest req:%v", req)

	player := d.getPlayer(ctx.Request().UserId)
	if player.getPlayerReady() {
		d.sendErrorResponse(player.UserId, codes.PlayerIsReady)
		return
	}
	// 记录玩家准备
	player.setPlayerReady()
	// 广播玩家准备
	d.broadcastPlayerReadyResponse()
}

// 玩家换牌请求
func (d *Desk) onSwapRequest(ctx net.MsgContext) {
	req := &pb.GamePlayerSwapRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onSwapRequest req:%v", req)
	player := d.getPlayer(ctx.Request().UserId)
	// 检测玩家是否已经选择
	if player.getSwapCards().Len() != 0 {
		d.sendErrorResponse(player.UserId, codes.PlayerIsSwap)
		return
	}

	// 设置玩家换牌
	cards := lib.Uint32ToCards(req.GetCards())
	err = player.deleteHandCards(cards)
	if err != nil {
		d.sendErrorResponse(player.UserId, codes.SwapReqInfoError)
		//d.restoreDeskInfo(player.SeatId)
		return
	}
	player.setSwapCards(cards)

	// 广播玩家换牌响应
	d.broadcastPlayerSwapResponse(player.SeatId)

	// 更新玩家手牌
	//d.broadcastUserCardsResponse(player.SeatId, pb.UpdateMahjongType_UPDATE_SWAP)
}

func (d *Desk) onChooseMissRequest(ctx net.MsgContext) {
	req := &pb.GamePlayerChooseMissRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onSwapRequest req:%v", req)
	player := d.getPlayer(ctx.Request().UserId)
	// 检测玩家是否已操作
	if player.getChooseMissType() != pb.MissType_MISS_NULL {
		d.sendErrorResponse(player.UserId, codes.PlayerIsChooseMiss)
		return
	}

	// 设置玩家选缺
	player.setChooseMissType(req.MissType)

	// 广播玩家选缺响应
	d.broadcastPlayerChooseMissResponse()
}

// 玩家出牌请求
func (d *Desk) onOutCardRequest(ctx net.MsgContext) {
	req := &pb.GamePlayerOutCardRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onOutCardRequest req:%v", req)
	player := d.getPlayer(ctx.Request().UserId)
	if player.isHu() {
		log.Sugar().Error("Player:", player, " isHu no out card")
		d.sendErrorResponse(player.UserId, codes.PlayerNotHaveOutCardAuth)
	}
	d.handleOutCardReq(req)
}

func (d *Desk) onOperateRequest(ctx net.MsgContext) {
	req := &pb.GamePlayerActionRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onOperateRequest req:%v", req)
	d.handleOperateReq(req)
}

// 强制游戏结束
func (d *Desk) onForceEndRequest(ctx net.MsgContext) {
	req := &pb.GameForceEndRequest{}
	err := ctx.ShouldBind(req)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}
	log.Sugar().Infof("onForceEndRequest req:%v", req)

	pbResponse := &pb.GameForceEndResponse{}
	d.sendDeskMessage(pb.MID_FORCE_END, pbResponse)

	d.toEnd()
}
