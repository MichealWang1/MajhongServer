package game

import (
	"errors"
	"fmt"
	"kxmj.common/codes"
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/internal/ai"
	"kxmj.game.mjxlch/pb"
	"math/rand"
	"time"
)

// 处理玩家出牌
func (d *Desk) handleOutCardReq(req *pb.GamePlayerOutCardRequest) {
	player := d.getPlayerBySeat(req.SeatId)
	outCard := lib.Card(req.GetCard())
	log.Sugar().Infof("handleOutCardReq req:%v,outCard:%v", req, outCard)
	// 校验位子是否可以出牌
	if err, code := d.checkOutCardReq(player.SeatId, outCard); err != nil {
		log.Sugar().Error("checkOutCardReq error:", err)
		d.sendErrorResponse(player.UserId, code)
		return
	}
	// 如果是托管
	if req.GetIsAuto() != player.getHostState() {
		//d.changePlayerHost(player.SeatId, req.GetIsAuto())
	}

	// 如果玩家有动作
	if player.getOperationalActions().Len() > 0 {
		err := d.onOperatePass(player.SeatId)
		if err != nil {
			log.Sugar().Errorf("handleOutCardReq error: %v", err)
		}
	}
	// 记录玩家出牌
	player.addDiscards(outCard)
	player.addOutCount()
	// 从手牌里删除这张牌
	err := player.deleteHandCard(outCard)
	if err != nil {
		log.Sugar().Errorf("handleOutCardReq error: %v", err)
	}
	// 重置玩家摸牌
	player.resetCatchCard()
	// 广播玩家出牌
	d.broadcastPlayerOutCardResponse(player.SeatId, outCard, req.GetIndex())

	// 提示玩家听牌提示
	d.calculateTingCards3N1(player.SeatId)

	handCards := player.getHandCards().Copy()
	handCards.Sort()
	log.Sugar().Infof("seatId:%v,handCards:%v,discards:%v", player.SeatId, handCards, player.getDiscards())
	// 设置当前玩家操作完成
	d.setCurrentStateCanOver(true)

	// 检测其他人的动作
	if actionTable := d.detectionActionOnOther(player.SeatId, outCard); len(actionTable) > 0 {
		for seat, actions := range actionTable {
			//d.sendPlayerOperateNotify(seat, actions)
			p := d.getPlayerBySeat(seat)
			p.setOperationalActions(actions)
		}
		return
	}
	// 设置是否检测胡杠
	d.setCurrentStateCheck(CheckGangHu)
	d.resetPrevActionGangCards()
	// 获取下一个可操作玩家
	nextPlayer := d.getNextOperatePlayer(player.SeatId)
	// 切换下一个玩家出牌
	d.setOperateSeatId(nextPlayer.SeatId)
	// 下一个玩家摸牌
	err = d.catchCard(nextPlayer.SeatId)
	if err != nil {
		log.Sugar().Infof("%v", err)
		return
	}
}

// 校验出牌位子
func (d *Desk) checkOutCardReq(seatId uint32, outCard lib.Card) (error, int) {
	if d.status != Playing {
		return errors.New(fmt.Sprintf("curstatus is %v not Playing", d.status)), codes.NowStateNotOperate
	}
	if !d.checkOperateSeatId(seatId) {
		return errors.New(fmt.Sprintf("operateSeatId is %v not %v", d.runtimeData.operateSeatId, seatId)), codes.PlayerNotHaveOutCardAuth
	}
	player := d.getPlayerBySeat(seatId)
	// 检测时间
	if player.getOpTime()+DurationOut > time.Now().UnixMilli() {
		return errors.New(fmt.Sprintf("op time fast !!!")), codes.PlayerNotHaveOutCardAuth
	}

	if !outCard.IsValid() {
		return errors.New(fmt.Sprintf("outCard:%v is not valid", outCard)), codes.OperateReqInfoError
	}

	if !player.getHandCards().In(outCard) {
		return errors.New(fmt.Sprintf("handCards:%v is not have card:%v", player.getHandCards(), outCard)), codes.OperateReqInfoError
	}
	// 手上有缺牌不能打其他牌
	//chooseMissCards := player.getHandCards().GetAllColorCards(ChooseMissTypeToLib(player.getChooseMissType()))
	//if chooseMissCards.Len() > 0 && !chooseMissCards.In(outCard) {
	//	return errors.New(fmt.Sprintf("You should first pick the missing card from your hand and play it！")), codes.PlayerNotHaveOutCardAuth
	//}
	return nil, 0
}

// 获取下一个可操作的位子
func (d *Desk) getNextOperatePlayer(seatId uint32) *Player {
	p := d.getNextSeatPlayer(seatId)
	if p.getHandCards().Len() == 0 {
		return d.getNextOperatePlayer(p.SeatId)
	}
	return p
}

// 检测自己的动作(杠、胡)
func (d *Desk) detectionActionOnSelf(seatId uint32, checkType CheckPlaying) lib.UserActions {
	player := d.getPlayerBySeat(seatId)
	handCards := player.getHandCards().Copy()
	opCard := player.getCatchCard()
	actionCards := player.getActionCardsTable()
	res := make(lib.UserActions, 0)
	// 校验手牌是否是3n+2
	if handCards.Len()%3 != 2 || checkType == NotCheck {
		return res
	}
	newCards, _ := handCards.DeleteCard(opCard)
	checkAction := lib.CheckAction{
		SeatId:      int32(player.SeatId),
		OutSeatId:   int32(player.SeatId),
		HandCards:   newCards,
		OutCard:     opCard,
		UserActions: actionCards,
		ZiStraight:  false,
	}
	// 检测杠
	if !player.isHu() {
		if actions := checkAction.CheckActionGang(); actions.Len() > 0 {
			log.Sugar().Infof("actions:%v", actions)
			for _, action := range actions {
				if action.OutCard.GetColor() != ChooseMissTypeToLib(player.getChooseMissType()) {
					res = append(res, action)
				}
			}
		}
	}

	// 检测胡
	if checkType == CheckGangHu {
		analysisHu := &AnalysisHu{
			SeatId:           seatId,
			OutSeatId:        seatId,
			BankerSeatId:     d.runtimeData.bankerSeatId,
			OpCard:           opCard,
			HandCards:        newCards,
			ChooseMissColor:  ChooseMissTypeToLib(player.getChooseMissType()),
			ActionCards:      player.getActionCardsTable(),
			OutCount:         player.getOutCount(),
			IsQiangGang:      false,
			PrevActionIsGang: d.getPrevActionIsGang() > 0,
			IsHaiDi:          d.runtimeData.cardStack.GetResidueCardsNum() == 0,
		}

		canHu, result := analysisHu.analysisHuResult()
		if canHu {
			huData := &HuData{
				seatId:      analysisHu.SeatId,
				outSeatId:   analysisHu.OutSeatId,
				opCard:      analysisHu.OpCard,
				handCards:   analysisHu.HandCards.Copy(),
				actionCards: analysisHu.ActionCards.Copy(),
				result:      result,
			}
			// 记录玩家能胡的数据
			player.setCurrentHuResult(huData)
			huAction := &lib.UserAction{
				SeatId:     int32(player.SeatId),
				OutSeatId:  int32(player.SeatId),
				ActionType: lib.ActionType_Hu,
				OutCard:    opCard,
			}
			res = append(res, huAction)
		}

	}
	log.Sugar().Infof(fmt.Sprintf("detectionActionOnSelf actions: %#v", res))
	return res
}

// 检测其他人的吃碰杠胡
func (d *Desk) detectionActionOnOther(outSeatId uint32, outCard lib.Card) map[uint32]lib.UserActions {
	res := make(map[uint32]lib.UserActions, 0)
	for _, p := range d.players {
		// 不检测自己
		if p.SeatId == outSeatId {
			continue
		}
		userActions := make(lib.UserActions, 0)
		// 构造吃碰杠结构体
		handCards := p.getHandCards().Copy()
		actionCards := p.getActionCardsTable()
		checkAction := lib.CheckAction{
			SeatId:      int32(p.SeatId),
			OutSeatId:   int32(outSeatId),
			HandCards:   handCards,
			OutCard:     outCard,
			UserActions: actionCards,
			ZiStraight:  false,
		}
		// 检测吃(吃只检测下家)
		//if d.getNextSeatPlayer(outSeatId).SeatId != p.SeatId {
		//	if actions := checkAction.CheckActionChi(); actions.Len() > 0 {
		//		userActions = append(userActions, actions...)
		//	}
		//}
		// 检测碰
		if actions := checkAction.CheckActionPeng(); !p.isHu() && actions.Len() > 0 {
			for _, action := range actions {
				log.Sugar().Infof("actions:%v", actions)
				if action.OutCard.GetColor() != ChooseMissTypeToLib(p.getChooseMissType()) {
					userActions = append(userActions, action)
				}
			}
		}
		// 检测杠
		if actions := checkAction.CheckActionGang(); actions.Len() > 0 {
			for _, action := range actions {
				log.Sugar().Infof("actions:%v", actions)
				if action.OutCard.GetColor() != ChooseMissTypeToLib(p.getChooseMissType()) {
					userActions = append(userActions, action)
				}
			}
		}
		// todo:检测胡
		analysisHu := &AnalysisHu{
			SeatId:           p.SeatId,
			OutSeatId:        outSeatId,
			BankerSeatId:     d.runtimeData.bankerSeatId,
			OpCard:           outCard,
			HandCards:        handCards,
			ChooseMissColor:  ChooseMissTypeToLib(p.getChooseMissType()),
			ActionCards:      p.getActionCardsTable(),
			OutCount:         p.getOutCount(),
			IsQiangGang:      false,
			PrevActionIsGang: d.getPrevActionIsGang() > 0,
			IsHaiDi:          d.runtimeData.cardStack.GetResidueCardsNum() == 0,
		}

		canHu, result := analysisHu.analysisHuResult()
		if canHu {
			huData := &HuData{
				seatId:      analysisHu.SeatId,
				outSeatId:   analysisHu.OutSeatId,
				opCard:      analysisHu.OpCard,
				handCards:   analysisHu.HandCards.Copy(),
				actionCards: analysisHu.ActionCards.Copy(),
				result:      result,
			}
			// 记录玩家能胡的数据
			p.setCurrentHuResult(huData)
			huAction := &lib.UserAction{
				SeatId:     int32(p.SeatId),
				OutSeatId:  int32(outSeatId),
				ActionType: lib.ActionType_Hu,
				OutCard:    outCard,
			}
			userActions = append(userActions, huAction)
		}
		if userActions.Len() > 0 {
			log.Sugar().Infof(fmt.Sprintf("detectionActionOnOther seatId:%v, actions: %#v", p.SeatId, userActions))
			res[p.SeatId] = userActions
		}
	}
	log.Sugar().Infof("detectionActionOnOther actions:%#v", res)
	return res
}

// 只检测胡（补杠时）
func (d *Desk) detectHuOnOther(outSeatId uint32, outCard lib.Card) map[uint32]lib.UserActions {
	res := make(map[uint32]lib.UserActions, 0)
	for _, p := range d.players {
		if p.SeatId == outSeatId {
			continue
		}
		userActions := make(lib.UserActions, 0)
		// 检测胡
		analysisHu := &AnalysisHu{
			SeatId:           p.SeatId,
			OutSeatId:        outSeatId,
			BankerSeatId:     d.runtimeData.bankerSeatId,
			OpCard:           outCard,
			HandCards:        p.getHandCards().Copy(),
			ChooseMissColor:  ChooseMissTypeToLib(p.getChooseMissType()),
			ActionCards:      p.getActionCardsTable(),
			OutCount:         p.getOutCount(),
			IsQiangGang:      true,
			PrevActionIsGang: d.getPrevActionIsGang() > 0,
			IsHaiDi:          d.runtimeData.cardStack.GetResidueCardsNum() == 0,
		}

		canHu, result := analysisHu.analysisHuResult()
		if canHu {
			huData := &HuData{
				seatId:      analysisHu.SeatId,
				outSeatId:   analysisHu.OutSeatId,
				opCard:      analysisHu.OpCard,
				handCards:   analysisHu.HandCards.Copy(),
				actionCards: analysisHu.ActionCards.Copy(),
				result:      result,
			}
			// 记录玩家能胡的数据
			p.setCurrentHuResult(huData)
			huAction := &lib.UserAction{
				SeatId:     int32(p.SeatId),
				OutSeatId:  int32(outSeatId),
				ActionType: lib.ActionType_Hu,
				OutCard:    outCard,
			}
			userActions = append(userActions, huAction)
		}
		if userActions.Len() > 0 {
			res[p.SeatId] = userActions
		}
	}
	log.Sugar().Infof("detectHuOnOther actions:%v", res)
	return res
}

// 出牌错误数据恢复
func (d *Desk) restoreDeskInfo(seatId uint32) {
	d.broadcastUserCardsResponse(seatId, pb.UpdateMahjongType_UPDATE_ERROR)
}

// 改变玩家托管状态
func (d *Desk) changePlayerHost(seatId uint32, status bool) {
	player := d.getPlayerBySeat(seatId)
	if player.getHostState() == status { // 状态相同不处理
		return
	}
	player.setHostState(status)
	d.broadcastPlayerHostNotify()
}

// 出牌超时处理
func (d *Desk) onOutCardTimeOver() {
	log.Sugar().Infof("onOutCardTimeOver")
	player := d.getPlayerBySeat(d.getOperateSeatId())
	outCard := player.getCatchCard()
	swapColor := ChooseMissTypeToLib(player.getChooseMissType())
	//for _, c := range player.getHandCards() {
	//	if c.GetColor() == swapColor {
	//		outCard = c
	//	}
	//}

	if outCard == lib.INVALID_CARD || !player.getHandCards().In(outCard) {
		outCard = player.getHandCards()[0]
	}

	if player.IsRobot && !player.isHu() {
		robotArgs := ai.NewEstimateArgs()
		robotArgs.SeatId = int32(player.SeatId)
		robotArgs.OutSeatId = int32(player.SeatId)
		robotArgs.UserCount = player.getOutCount()
		robotArgs.HandCards = player.getHandCards()
		robotArgs.DisCards = player.getDiscards()
		robotArgs.MustCanOutCards = player.getHandCards().GetAllColorCards(swapColor)

		if card := ai.AnalysisPlayerOutCard(robotArgs); card != lib.INVALID_CARD {
			outCard = card
		}
	}

	req := &pb.GamePlayerOutCardRequest{
		SeatId: player.SeatId,
		Card:   outCard.ToUint32(),
		IsAuto: true,
		Index:  uint32(player.getHandCards().Len()),
	}
	d.handleOutCardReq(req)
}

// 机器人出牌策略
func (d *Desk) onRobotOutCard() {
	player := d.getPlayerBySeat(d.getOperateSeatId())
	if !player.IsRobot {
		return
	}
	outCard := player.getCatchCard()
	swapColor := ChooseMissTypeToLib(player.getChooseMissType())
	for _, c := range player.getHandCards() {
		if c.GetColor() == swapColor {
			outCard = c
		}
	}
	if outCard == lib.INVALID_CARD {
		outCard = player.getHandCards()[0]
	}

	req := &pb.GamePlayerOutCardRequest{
		SeatId: player.SeatId,
		Card:   outCard.ToUint32(),
		Index:  uint32(rand.Int31n(int32(player.getHandCards().Len()))),
	}
	d.handleOutCardReq(req)
}

// pb选缺类型转换成lib
func ChooseMissTypeToLib(missType pb.MissType) lib.CardColor {
	switch missType {
	case pb.MissType_MISS_WAN:
		return lib.Card_Color_Wan
	case pb.MissType_MISS_TIAO:
		return lib.Card_Color_Tiao
	case pb.MissType_MISS_TONG:
		return lib.Card_Color_Tong
	default:
		return lib.Card_Color_Unknown
	}
}
