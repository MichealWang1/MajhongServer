package game

import (
	"kxmj.common/codes"
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/utils"
	"kxmj.game.mjxlch/pb"
	"time"
)

// 广播错误信息
func (d *Desk) sendErrorResponse(userId uint32, code int) {
	if d.getPlayer(userId).IsRobot {
		return
	}
	pbResponse := &pb.GameErrorResponse{
		Code: uint32(code),
		Msg:  codes.GetMessage(code),
	}
	log.Sugar().Errorf("GameErrorResponse = %v", pbResponse)
	d.game.SendMessage(userId, uint16(pb.MID_ERROR_MESSAGE), pbResponse)
}

// BroadcastUserInfo 广播玩家信息
func (d *Desk) broadcastEnterInfoNotify(userId uint32) {
	player := d.getPlayer(userId)
	pbPlayer := &pb.PlayerInfo{
		SeatId:     player.SeatId,
		UserId:     player.UserId,
		Gold:       player.Gold,
		Nickname:   player.Nickname,
		AvatarAddr: player.AvatarAddr,
		IsOnline:   player.IsOnline,
	}

	pbResponse := &pb.GamePlayerEnterInfoNotify{
		Player: pbPlayer,
	}
	log.Sugar().Infof("GamePlayerEnterInfoNotify:", d.roundId, "  pbResponse:", pbResponse)
	d.sendOtherPlayerMessage(player.UserId, pb.MID_PLAYER_ENTER_INFO_NOTIFY, pbResponse)
}

func (d *Desk) broadcastDeskInfo(userId uint32) {
	playerInfos := make([]*pb.PlayerInfo, 0)
	// 获取其他位子信息
	for _, p := range d.players {
		if p.UserId <= 0 {
			continue
		}
		cards := &pb.CardsInfo{
			HandCards: lib.Cards{}.ToUint32(),
			HandCount: uint32(p.getHandCards().Len()),
			OpCard:    lib.INVALID_CARD.ToUint32(),
			Actions:   d.toMahjongAction(p.SeatId),
			HuResult:  d.toHuResult(p.SeatId),
			Discards:  p.getDiscards().ToUint32(),
		}
		swapDefaultCards := lib.Cards{}.ToUint32()
		if p.UserId == userId {
			cards.HandCards = p.getHandCards().ToUint32()
			cards.OpCard = p.getCatchCard().ToUint32()
			if d.status == Swap && p.getSwapCards().Len() == 0 {
				swapDefaultCards = GetSmallColorCards(p.getHandCards(), 3).ToUint32()
			}
		}

		otherPlayer := &pb.PlayerInfo{
			SeatId:     p.SeatId,
			UserId:     p.UserId,
			Gold:       p.Gold,
			Nickname:   p.Nickname,
			AvatarAddr: p.AvatarAddr,
			IsOnline:   p.IsOnline,

			Cards:              cards,
			SwapStatus:         p.getSwapCards().Len() > 0,
			SwapDefaultCards:   swapDefaultCards,
			ChooseMissStatus:   p.getChooseMissType(),
			OperationalActions: UACardsToPtl(p.getOperationalActions()),
			HostStatus:         p.getHostState(),
		}
		playerInfos = append(playerInfos, otherPlayer)
	}

	pbResponse := &pb.GameDeskInfoResponse{
		Player:       playerInfos,
		State:        pb.GameState(d.status),
		OperateSeat:  d.getOperateSeatId(),
		Duration:     uint32((d.nextTime - time.Now().UnixMilli()) / 1000),
		StackCount:   uint32(d.runtimeData.cardStack.GetResidueCardsNum()),
		BankerSeatId: d.runtimeData.bankerSeatId,
	}
	log.Sugar().Infof("GameDeskInfoResponse pbResponse: %v", pbResponse)
	d.game.SendMessage(userId, uint16(pb.MID_DESK_INFO), pbResponse)
}

// 广播玩家准备
func (d *Desk) broadcastPlayerReadyResponse() {
	state := make([]bool, PLAYER_COUNT)
	for _, p := range d.players {
		if p.getPlayerReady() {
			state[p.SeatId] = true
		}
	}
	pbResponse := &pb.GamePlayerReadyResponse{
		State: state,
	}
	log.Sugar().Infof("GamePlayerReadyResponse pbResponse: ", pbResponse)
	d.sendDeskMessage(pb.MID_PLAYER_READY, pbResponse)
}

// BroadcastGameStart 广播游戏开始
func (d *Desk) broadcastGameStart() {
	log.Sugar().Infof("GameStartNotify game start")
	d.sendDeskMessage(pb.MID_GAME_START, &pb.GameStartNotify{BankerSeatId: d.runtimeData.bankerSeatId})
}

// 广播玩家金币数量
func (d *Desk) broadcastPlayerGoldNumber() {
	golds := make([]string, PLAYER_COUNT)
	for _, p := range d.players {
		if p.IsRobot {
			golds[p.SeatId] = p.Gold
			continue
		}
		// 获取玩家金币
		gold, _ := d.game.Server().GetUserGold(p.UserId, d.room.ID(), d.room.RoomLevel())
		if !p.GoldStatus { // 第一次设置金额
			p.Gold = gold.Data.Gold
			p.GoldStatus = true
		} else {
			p.Gold, _ = utils.AddToString(p.Gold, gold.Data.Gold)
		}
		golds[p.SeatId] = p.Gold
	}
	pbResponse := &pb.GameUpdateGoldNumberNotify{
		Gold: golds,
	}
	log.Sugar().Infof("GameUpdateGoldNumberNotify pbResponse: ", pbResponse)
	d.sendDeskMessage(pb.MID_UPDATE_GOLD_NUMBER, pbResponse)
}

// 广播游戏状态
//func (d *Desk) broadcastGameStatusNotify() {
//	pbResponse := &pb.GameStateNotify{
//		State: pb.GameState(d.status),
//	}
//	log.Sugar().Infof("GameStatusNotify pbResponse:", pbResponse)
//	d.sendDeskMessage(pb.MID_GAME_STATE, pbResponse)
//}

func (d *Desk) broadcastDealHandCards() {
	for _, p1 := range d.players {
		if p1.IsRobot {
			continue
		}
		allHandCardsInfo := make([]*pb.HandCardsInfo, 0, PLAYER_COUNT)
		for _, p2 := range d.players {
			handCardsInfo := &pb.HandCardsInfo{
				HandCards: lib.Cards{}.ToUint32(),
				CatchCard: lib.INVALID_CARD.ToUint32(),
				Count:     uint32(p2.getHandCards().Len()),
			}
			allHandCardsInfo = append(allHandCardsInfo, handCardsInfo)
		}
		allHandCardsInfo[p1.SeatId].HandCards = p1.getHandCards().ToUint32()
		allHandCardsInfo[p1.SeatId].CatchCard = p1.getCatchCard().ToUint32()
		pbResponse := &pb.GameDealHandCardsNotify{
			SeatId:        p1.SeatId,
			HandCardsInfo: allHandCardsInfo,
			StackCount:    uint32(d.runtimeData.cardStack.GetResidueCardsNum()),
		}
		log.Sugar().Infof("GameDealHandCardsNotify pbResponse:", pbResponse)
		d.game.SendMessage(p1.UserId, uint16(pb.MID_DEAL_HAND_CARDS_NOTIFY), pbResponse)
	}
}

func (d *Desk) sendPlayerCardsInfo(seatId uint32, updateType pb.UpdateMahjongType) {

}

// BroadcastUserCardsResponse 广播玩家手牌信息给所有人(对自己要知道所有信息，其他人则不需要知道手牌)
func (d *Desk) broadcastUserCardsResponse(seatId uint32, updateType pb.UpdateMahjongType) {
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		ptl := d.toUpdateMahjongPtl(seatId, p.SeatId == seatId, updateType)
		log.Sugar().Infof("GameUpdateMahjongResponse userId:%v,pbResponse:%v", p.UserId, ptl)
		d.game.SendMessage(p.UserId, uint16(pb.MID_UPDATE_PLAYER_CARDS_DETAIL_NOTIFY), ptl)
	}
}

// ToUpdateMahjongPtl 将手牌转换成pb数据
func (d *Desk) toUpdateMahjongPtl(seatId uint32, isSendSelf bool, updateType pb.UpdateMahjongType) *pb.GameUpdateMahjongResponse {
	p := d.getPlayerBySeat(seatId)

	handCards := p.getHandCards() // 手牌
	catchCard := p.getCatchCard() // 抓的牌
	discards := p.getDiscards()   // 弃牌堆

	actions := d.toMahjongAction(seatId)
	pbResponse := &pb.GameUpdateMahjongResponse{
		SeatId:     seatId,
		OpCard:     lib.INVALID_CARD.ToUint32(),
		UpdateType: updateType,
		Actions:    actions,
		Count:      uint32(handCards.Len()),
		Discard:    discards.ToUint32(),
		HuResult:   d.toHuResult(seatId),
	}
	// 如果是自己则带上手牌信息和操作牌信息
	if isSendSelf {
		pbResponse.HandCard = handCards.ToUint32()
		pbResponse.OpCard = catchCard.ToUint32()
	}
	return pbResponse
}

// 将胡牌数据转换成pb
func (d *Desk) toHuResult(seatId uint32) []*pb.HuResult {
	p := d.getPlayerBySeat(seatId)
	huData := p.getHuData()
	res := make([]*pb.HuResult, 0, len(huData))
	for _, data := range huData {
		tmp := HuDataToPtl(data)
		res = append(res, tmp)
	}
	return res
}

func HuDataToPtl(data *HuData) *pb.HuResult {
	return &pb.HuResult{
		OutSeatId:  data.outSeatId,
		HuCard:     data.opCard.ToUint32(),
		HuPosition: data.result.getFan(),
		Multiple:   data.result.multiple,
	}
}

// ToMahjongAction 将玩家所有动作牌转换成pb数据
func (d *Desk) toMahjongAction(seatId uint32) []*pb.MahjongAction {
	p := d.getPlayerBySeat(seatId)
	items := p.getActionCardsTable() // 玩家动作
	return UACardsToPtl(items)
}

// UACardsToPtl 将所有动作牌转换成pb数据
func UACardsToPtl(items lib.UserActions) []*pb.MahjongAction {
	ptl := make([]*pb.MahjongAction, 0, len(items))
	for _, item := range items {
		obj := UAToPtl(item)
		ptl = append(ptl, obj)
	}
	return ptl
}

// UAToPtl 将动作牌转换成pb数据
func UAToPtl(item *lib.UserAction) *pb.MahjongAction {
	return &pb.MahjongAction{
		OutSeatId:     uint32(item.OutSeatId),
		ActionType:    ActionTypeLibToPb(item.ActionType),
		ExtensionType: ExtActionTypeLibToPb(item.ExtraActionType),
		OpCard:        item.OutCard.ToUint32(),
		DeleteCards:   item.DeleteCards.ToUint32(),
		CombineCards:  item.CombineCards.ToUint32(),
	}
}

// BroadcastUserCatchCard 广播玩家抓牌
func (d *Desk) broadcastPlayerCatchCard(seatId uint32, catchCard lib.Card) {
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		pbResponse := &pb.GamePlayerCatchCard{
			SeatId:    seatId,
			Card:      lib.INVALID_CARD.ToUint32(),
			LeftCount: uint32(d.runtimeData.cardStack.GetResidueCardsNum()),
		}
		if seatId == p.SeatId {
			pbResponse.Card = catchCard.ToUint32()
		}
		d.game.SendMessage(p.UserId, uint16(pb.MID_CATCH_CARD_NOTIFY), pbResponse)
	}
	log.Sugar().Infof("GamePlayerCatchCard seatId:%v,catchCard:%v\n", seatId, catchCard)
}

// SendUserSwapCardNotify 通知玩家换牌
func (d *Desk) sendPlayerSwapCardNotify(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	pbResponse := &pb.GamePlayerSwapNotify{
		SeatId:       seatId,
		Duration:     DurationSwap / 1000,
		DefaultCards: GetSmallColorCards(player.getHandCards(), 3).ToUint32(),
	}
	log.Sugar().Infof("GamePlayerSwapNotify pbResponse:%v", pbResponse)
	d.game.SendMessage(player.UserId, uint16(pb.MID_SWAP_NOTIFY), pbResponse)
}

// BroadcastUserSwapResponse 玩家换牌响应
func (d *Desk) broadcastPlayerSwapResponse(seatId uint32) {
	res := make([]bool, PLAYER_COUNT)
	for _, p := range d.players {
		if p.getSwapCards().Len() != 0 {
			res[int(p.SeatId)] = true
		}
	}
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		pbResponse := &pb.GamePlayerSwapResponse{
			SeatId: seatId,
		}
		if p.SeatId == seatId {
			pbResponse.SwapCards = p.getSwapCards().ToUint32()
		}
		log.Sugar().Infof("GamePlayerSwapResponse pbResponse:%v", pbResponse)
		d.game.SendMessage(p.UserId, uint16(pb.MID_SWAP_INFO), pbResponse)
	}
}

// sendPlayerSwapTypeNotify 通知玩家换牌类型
func (d *Desk) sendPlayerSwapTypeNotify(seatId uint32, cards lib.Cards) {
	p := d.getPlayerBySeat(seatId)
	if p.IsRobot {
		return
	}
	pbResponse := &pb.GameSwapResultNotify{
		SwapType:  d.runtimeData.swapType,
		SwapCards: cards.ToUint32(),
	}
	log.Sugar().Infof("GameSwapResultNotify pbResponse:%v", pbResponse)
	d.game.SendMessage(p.UserId, uint16(pb.MID_SWAP_RESULT_NOTIFY), pbResponse)
}

// 提示玩家选缺
func (d *Desk) sendPlayerChooseMissNotify(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	pbResponse := &pb.GamePlayerChooseMissNotify{
		SeatId:   seatId,
		Duration: DurationChooseMiss / 1000,
	}
	log.Sugar().Infof("GamePlayerChooseMissNotify pbResponse:%v", pbResponse)
	d.game.SendMessage(player.UserId, uint16(pb.MID_CHOOSE_MISS_NOTIFY), pbResponse)
}

// 玩家选缺响应
func (d *Desk) broadcastPlayerChooseMissResponse() {
	res := make([]bool, PLAYER_COUNT)
	for _, p := range d.players {
		if p.getChooseMissType() != pb.MissType_MISS_NULL {
			res[int(p.SeatId)] = true
		}
	}
	pbResponse := &pb.GamePlayerChooseMissResponse{
		Status: res,
	}
	log.Sugar().Infof("GamePlayerChooseMissResponse pbResponse:%v", pbResponse)
	d.sendDeskMessage(pb.MID_CHOOSE_MISS_INFO, pbResponse)
}

// 广播选缺结果
func (d *Desk) broadcastPlayerChooseMissResultNotify() {
	res := make([]pb.MissType, PLAYER_COUNT)
	for _, p := range d.players {
		res[int(p.SeatId)] = p.getChooseMissType()
	}
	pbResponse := &pb.GamePlayerChooseMissResultNotify{
		MissType: res,
	}
	log.Sugar().Infof("GamePlayerChooseMissResultNotify pbResponse:%v", pbResponse)
	d.sendDeskMessage(pb.MID_CHOOSE_MISS_RESULT_NOTIFY, pbResponse)
}

// 广播玩家托管
func (d *Desk) broadcastPlayerHostNotify() {
	hostState := make([]bool, PLAYER_COUNT)
	for _, p := range d.players {
		if p.getHostState() {
			hostState[int(p.SeatId)] = true
		}
	}
	pbResponse := &pb.GamePlayerAutoResponse{
		IsHost: hostState,
	}
	log.Sugar().Infof("GamePlayerAutoResponse pbResponse: ", pbResponse)
	d.sendDeskMessage(pb.MID_PLAYER_AUTO_INFO, pbResponse)
}

// 提示玩家出牌
func (d *Desk) sendPlayerOutCardNotify(seatId uint32) {
	//player := d.getPlayerBySeat(seatId)
	pbResponse := &pb.GamePlayerOutCardNotify{
		SeatId:   seatId,
		Duration: DurationPlaying / 1000,
	}
	log.Sugar().Infof("GamePlayerOutCardNotify pbResponse:%v", pbResponse)
	d.sendDeskMessage(pb.MID_OUT_CARD_NOTIFY, pbResponse)
}

// broadcastPlayerOutCardResponse 广播玩家出牌
func (d *Desk) broadcastPlayerOutCardResponse(seatId uint32, outCard lib.Card, index uint32) {
	pbResponse := &pb.GamePlayerOutCardResponse{
		SeatId: seatId,
		Card:   outCard.ToUint32(),
		Index:  index,
	}
	log.Sugar().Infof("GamePlayerOutCardResponse pbResponse: %v", pbResponse)
	d.sendDeskMessage(pb.MID_OUT_CARD_INFO, pbResponse)
}

// 提示玩家处理动作
func (d *Desk) sendPlayerOperateNotify(seatId uint32, actions lib.UserActions) {
	player := d.getPlayerBySeat(seatId)
	pbAction := UACardsToPtl(actions)
	if player.getCurrentHuResult() != nil {
		for i, a := range pbAction {
			if a.ActionType == pb.ActionType_ACTION_HU {
				pbAction[i].HuMultiple = player.getCurrentHuResult().result.multiple
			}
		}
	}
	pbResponse := &pb.GamePlayerActionNotify{
		SeatId:   seatId,
		Actions:  pbAction,
		Duration: DurationPlaying / 1000,
	}
	log.Sugar().Infof("GamePlayerActionNotify pbResponse:%v", pbResponse)
	d.game.SendMessage(player.UserId, uint16(pb.MID_ACTIONS_NOTIFY), pbResponse)
}

// 玩家处理动作响应
func (d *Desk) broadcastPlayerOperateResponse(seatId uint32, actions *lib.UserAction) {
	pbResponse := &pb.GamePlayerActionResultResponse{
		SeatId: seatId,
		Action: UAToPtl(actions),
	}
	log.Sugar().Infof("GamePlayerActionResultResponse pbResponse: ", pbResponse)
	d.sendDeskMessage(pb.MID_ACTIONS_INFO, pbResponse)
}

// 结算数据转换成pb文件
func (d *Desk) settlementToPtl(settlements []*InningScores) []*pb.BureauSettlementInfo {
	data := make([]*pb.BureauSettlementInfo, 0)
	for _, settlement := range settlements {
		var huInfo *pb.HuInfo
		log.Sugar().Infof("settlement huResult:%#v", settlement.huResult)
		if settlement.settlementType == pb.SettlementType_HU {
			huInfo = &pb.HuInfo{
				HandCards:  settlement.handCards.ToUint32(),
				HuCard:     settlement.opCard.ToUint32(),
				Actions:    UACardsToPtl(settlement.userActions),
				HuPosition: settlement.huResult.getFan(),
			}
		}
		tmp := &pb.BureauSettlementInfo{
			WinSeatId:      settlement.winSeatId,
			SettlementType: settlement.settlementType,
			OpCard:         settlement.opCard.ToUint32(),
			HuCardsInfo:    huInfo,
			InningScores:   settlement.scores,
			RealityScores:  settlement.winLoseGolds,
			Multiple:       settlement.multiple,
			IsCeiling:      settlement.isCeiling,
			IsBankruptcy:   settlement.isBankruptcy,
		}
		data = append(data, tmp)
	}
	return data
}

// 广播玩家局内结算
func (d *Desk) broadcastBureauSettlementNotify() {
	golds := make([]string, PLAYER_COUNT)
	for _, p := range d.players {
		golds[p.SeatId] = p.Gold
	}

	pbResponse := &pb.GameBureauSettlementNotify{
		Data: d.settlementToPtl(d.getWaitSettlement()),
		Gold: golds,
	}
	log.Sugar().Infof("GameHuResultNotify pbResponse: %v", pbResponse)
	d.sendDeskMessage(pb.MID_SETTLEMENT_BUREAU, pbResponse)
}

// 广播最终结算
func (d *Desk) broadcastEndSettlementNotify() {
	data := make([]*pb.GameEndPersonalSettlementInfo, PLAYER_COUNT)
	for _, p := range d.players {
		endInfo := &pb.GameEndPersonalSettlementInfo{
			HandCards:   p.getHandCards().ToUint32(),
			Discards:    p.getDiscards().ToUint32(),
			Actions:     UACardsToPtl(p.getActionCardsTable()),
			HuResult:    d.toHuResult(p.SeatId),
			TotalScores: p.getTotalScore(),
			Gold:        p.Gold,
		}
		data[p.SeatId] = endInfo
	}
	pbResponse := &pb.GameEndSettlementNotify{
		EndType:           pb.EndType_NORMAL_END,
		Data:              data,
		BureauSettlements: d.settlementToPtl(d.getSettlementData()),
	}
	log.Sugar().Infof("GameEndSettlementNotify pbResponse: %v", pbResponse)
	d.sendDeskMessage(pb.MID_SETTLEMENT_END, pbResponse)
}
