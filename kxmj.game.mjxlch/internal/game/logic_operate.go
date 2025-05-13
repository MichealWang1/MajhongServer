package game

import (
	"errors"
	"kxmj.common/codes"
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
)

// 处理玩家操作
func (d *Desk) handleOperateReq(req *pb.GamePlayerActionRequest) {
	player := d.getPlayerBySeat(req.GetSeatId())
	action := &lib.UserAction{
		SeatId:          int32(req.GetSeatId()),
		OutSeatId:       int32(d.getOperateSeatId()),
		ActionType:      ActionTypePbToLib(req.GetAction().GetActionType()),
		ExtraActionType: ExtActionTypePbToLib(req.GetAction().GetExtensionType()),
		OutCard:         lib.Card(req.GetAction().OpCard),
		DeleteCards:     lib.Uint32ToCards(req.GetAction().GetDeleteCards()),
		CombineCards:    lib.Uint32ToCards(req.GetAction().GetCombineCards()),
	}
	// 判断是否有操作动作
	if err := d.checkOperateReq(player.SeatId, action); err != nil {
		log.Sugar().Errorf("onOperateAction error: %v", err)
		d.sendErrorResponse(player.UserId, codes.PlayerNotHaveOperateAuth)
		return
	}

	// 处理动作
	if err := d.onOperateAction(player.SeatId, action); err != nil {
		log.Sugar().Errorf("onOperateAction error: %v", err)
		d.sendErrorResponse(player.UserId, codes.OperateReqInfoError)
		return
	}

	// 判断是否可以结束，可以则返回最大操作
	if !d.operateActionCanOver() {
		return
	}
	// 获取最大的操作
	maxResponse := d.getMaxRespOperateAction()
	// 清除所有玩家的操作
	d.resetAllOperateActions()
	d.setCurrentStateCanOver(true)

	if maxResponse.SeatId == maxResponse.OutSeatId {
		d.onOperateSelf(maxResponse)
	} else {
		d.onOperateOther(maxResponse)
	}
}

// 校验出牌位子
func (d *Desk) checkOperateReq(seatId uint32, action *lib.UserAction) error {
	// 判断这个动作是否合理
	if action.ActionType == lib.ActionType_Unknown {
		return errors.New("this action type is not in the list of actions")
	}
	return nil
}

// 是否有可操作的动作
func (d *Desk) hasOperateActions() bool {
	for _, p := range d.players {
		if p.getOperationalActions().Len() > 0 {
			return true
		}
	}
	return false
}

// 玩家过这个动作
func (d *Desk) onOperatePass(seatId uint32) error {
	passAction := &lib.UserAction{
		SeatId:     int32(seatId),
		ActionType: lib.ActionType_Pass,
	}
	return d.onOperateAction(seatId, passAction)
}

// 玩家操作动作
func (d *Desk) onOperateAction(seatId uint32, action *lib.UserAction) error {
	p := d.getPlayerBySeat(seatId)
	// 判断这个玩家是否有动作
	if p.getOperationalActions().Len() == 0 {
		return errors.New("this player not have operational actions")
	}
	// 如果玩家是过
	if action.ActionType == lib.ActionType_Pass {
		// 记录弃胡、弃碰
		p.resetCurrentHuResult()
	} else {
		// 判断是否存在这个动作
		if !p.getOperationalActions().HasAction(action) {
			return errors.New("this action not in operational actions")
		}
	}

	// 记录玩家操作
	p.setRespOperateAction(action)
	// 重置
	p.resetOperationalActions()
	return nil
}

// 清除所有玩家操作
func (d *Desk) resetAllOperateActions() {
	for _, p := range d.players {
		p.resetOperationalActions()
		p.resetRespOperateAction()
	}
}

// 获取最大响应权限
func (d *Desk) getMaxRespOperateAction() *lib.UserAction {
	var max *lib.UserAction
	for _, p := range d.players {
		if p.getRespOperateAction() != nil {
			if max == nil {
				max = p.getRespOperateAction()
			} else if max.ActionType > p.getRespOperateAction().ActionType {
				max = p.getRespOperateAction()
			}
		}
	}
	return max
}

// 这个操作是否可以结束
func (d *Desk) operateActionCanOver() bool {
	for _, p := range d.players {
		if p.getOperationalActions().Len() > 0 {
			return false
		}
	}
	return true
}

// 处理玩家最终操作

// 处理自己
func (d *Desk) onOperateSelf(maxAction *lib.UserAction) {
	player := d.getPlayerBySeat(uint32(maxAction.SeatId))
	// 广播玩家操作
	d.broadcastPlayerOperateResponse(player.SeatId, maxAction)

	actionType := maxAction.ActionType
	switch actionType {
	case lib.ActionType_Pass:
		// 设置操作玩家
		d.setOperateSeatId(player.SeatId)
		// 设置是否检测胡杠
		d.setCurrentStateCheck(NotCheck)
		// 清除玩家可操作胡
		player.resetCurrentHuResult()

	case lib.ActionType_Gang:
		// 补杠
		if maxAction.ExtraActionType == lib.ExtraActionType_Bu_Gang {
			// 判断其他人有没有胡(抢杠胡)
			if actionTable := d.detectHuOnOther(player.SeatId, maxAction.OutCard); len(actionTable) > 0 {
				// 有人能胡
				for seat, actions := range actionTable {
					p := d.getPlayerBySeat(seat)
					p.setOperationalActions(actions)
				}
				// 记录补杠
				d.setCurrentBuGangAction(maxAction)
				return
			}
			// 记录玩家补杠
			player.actionBuGang(maxAction.OutCard)
		} else { // 暗杠
			// 记录玩家动作
			player.addActionCards(maxAction)
		}
		// 删除玩家手牌
		player.deleteHandCards(maxAction.DeleteCards)
		// 记录上一个动作是杠
		d.setPrevActionGangCard(maxAction.OutCard)
		// 更新玩家手牌信息
		d.broadcastUserCardsResponse(player.SeatId, pb.UpdateMahjongType_UPDATE_ACTION)
		// 设置操作玩家
		d.setOperateSeatId(player.SeatId)
		// 算分
		score := d.calculateGangScore(maxAction)
		// 记录结算分
		d.addWaitSettlement(score)
		// 记录玩家分
		d.addTotalScore(score)
		// 设置操作玩家
		d.setOperateSeatId(player.SeatId)
	case lib.ActionType_Hu:
		// 删除手上这张胡的牌
		player.deleteHandCard(maxAction.OutCard)
		// 计算胡牌倍数
		mult := player.getCurrentHuResult().result.multiple
		// 记录玩家胡牌倍数
		player.addHuData(player.getCurrentHuResult().copy())
		// 计算玩家输赢分
		score := d.calculateHuScores(player.SeatId, player.SeatId, mult, maxAction.OutCard, player.getCurrentHuResult().result)
		// 记录结算分
		d.addWaitSettlement(score)
		// 记录玩家分
		d.addTotalScore(score)
		// 设置下家出牌
		d.setOperateSeatId(d.getNextOperatePlayer(player.SeatId).SeatId)
	}
}

// 处理其他人
func (d *Desk) onOperateOther(maxAction *lib.UserAction) {
	player := d.getPlayerBySeat(uint32(maxAction.SeatId))
	// 广播玩家操作
	d.broadcastPlayerOperateResponse(player.SeatId, maxAction)
	switch maxAction.ActionType {
	case lib.ActionType_Pass:
		// 获取操作玩家
		opSeatId := d.getOperateSeatId()
		// 操作玩家下家摸牌
		nextPlayer := d.getNextOperatePlayer(uint32(opSeatId))
		// 设置操作玩家
		d.setOperateSeatId(nextPlayer.SeatId)
		// 清除玩家胡操作
		d.resetHuResult()
		// 可能是过补杠
		if d.getCurrentBuGangAction() != nil {
			action := d.getCurrentBuGangAction()
			p := d.getPlayerBySeat(uint32(action.SeatId))
			p.actionBuGang(action.OutCard)
			// 删除玩家手牌
			player.deleteHandCards(action.DeleteCards)
			// 设置操作玩家
			d.setOperateSeatId(p.SeatId)
			d.resetCurrentBuGangAction()
			// 记录上一个动作是杠
			d.setPrevActionGangCard(action.OutCard)
			// 更新玩家手牌信息
			d.broadcastUserCardsResponse(player.SeatId, pb.UpdateMahjongType_UPDATE_ACTION)
			// 算分
			score := d.calculateGangScore(action)
			// 记录结算分
			d.addWaitSettlement(score)
			// 记录玩家分
			d.addTotalScore(score)
			return
		}
		d.resetPrevActionGangCards()
		// 设置是否检测胡杠
		d.setCurrentStateCheck(CheckGangHu)
		err := d.catchCard(d.getPlayerBySeat(d.getOperateSeatId()).SeatId)
		if err != nil {
			log.Sugar().Infof("%v", err)
			return
		}
	case lib.ActionType_Peng:
		d.resetPrevActionGangCards()
		// 记录玩家动作
		player.addActionCards(maxAction)
		// 删除玩家手牌
		player.deleteHandCards(maxAction.DeleteCards)
		// 删除玩家弃牌
		d.getPlayerBySeat(uint32(maxAction.OutSeatId)).deleteDiscards(maxAction.OutCard)
		// 更新玩家手牌信息
		d.broadcastUserCardsResponse(player.SeatId, pb.UpdateMahjongType_UPDATE_ACTION)
		// 设置操作玩家
		d.setOperateSeatId(player.SeatId)
		// 设置是否检测胡杠
		d.setCurrentStateCheck(CheckGang)
	case lib.ActionType_Gang: // 明杠
		// 记录玩家动作
		player.addActionCards(maxAction)
		// 删除玩家手牌
		player.deleteHandCards(maxAction.DeleteCards)
		// 删除玩家弃牌
		d.getPlayerBySeat(uint32(maxAction.OutSeatId)).deleteDiscards(maxAction.OutCard)
		// 记录上一个动作是杠
		d.setPrevActionGangCard(maxAction.OutCard)
		// 更新玩家手牌信息
		d.broadcastUserCardsResponse(player.SeatId, pb.UpdateMahjongType_UPDATE_ACTION)
		// 设置操作玩家
		d.setOperateSeatId(player.SeatId)
		// 算分
		score := d.calculateGangScore(maxAction)
		// 记录结算分
		d.addWaitSettlement(score)
		// 记录玩家分
		d.addTotalScore(score)
	case lib.ActionType_Hu:
		loserPlayer := d.getPlayerBySeat(uint32(maxAction.OutSeatId))
		// 删除玩家弃牌
		loserPlayer.deleteDiscards(maxAction.OutCard)
		scores := make([]*InningScores, 0)
		// 从出牌玩家下家开始查胡
		for i := uint32(0); i < PLAYER_COUNT-1; i++ {
			p := d.getNextOperatePlayer(loserPlayer.SeatId + i)
			if p.getCurrentHuResult() == nil {
				continue
			}
			// 计算胡牌倍数
			mult := p.getCurrentHuResult().result.multiple
			// 记录玩家胡牌倍数
			p.addHuData(p.getCurrentHuResult().copy())
			// 计算玩家胡牌输赢分
			score := d.calculateHuScores(p.SeatId, uint32(maxAction.OutSeatId), mult, maxAction.OutCard, p.getCurrentHuResult().result)
			scores = append(scores, score)
			// 呼叫转移(不是一炮多响的杠上炮)
			if d.getCurrentHuCount() == 1 && p.getCurrentHuResult().result.hasFan(Gang_Pao) {
				transferScore := d.calculateTransfer(p.SeatId, loserPlayer.SeatId)
				scores = append(scores, transferScore)
			}
			// 设置下家出牌
			d.setOperateSeatId(d.getNextOperatePlayer(p.SeatId).SeatId)
		}
		// 分析一炮多响玩家破产问题
		if len(scores) > 1 {
			d.calculateLoseMostScores(scores)
		}
		for _, s := range scores {
			// 记录结算分
			d.addWaitSettlement(s)
			// 记录玩家分
			d.addTotalScore(s)
		}
		d.resetPrevActionGangCards()
	}
}

// 获取当前操作胡牌玩家人数
func (d *Desk) getCurrentHuCount() int {
	count := 0
	for _, p := range d.players {
		if p.getCurrentHuResult() != nil {
			count++
		}
	}
	return count
}

// 操作超时处理
func (d *Desk) onOperateTimeOver() {
	for _, p := range d.players {
		d.operate(p.SeatId)
	}
}

// 自动处理动作
func (d *Desk) onOperateAuto(wait bool) {
	operationActions := make(map[uint32]struct{}, 0)
	// 查询当前状态能操作的玩家进行操作
	for _, p := range d.players {
		if p.getOperationalActions().Len() > 0 {
			operationActions[p.SeatId] = struct{}{}
		}
	}
	for _, p := range d.players {
		if _, ok := operationActions[p.SeatId]; ok {
			if !wait || p.canAutoOperation() {
				d.operate(p.SeatId)
			}
		}
	}
}

func (d *Desk) operate(seatId uint32) {
	p := d.getPlayerBySeat(seatId)
	if p.getOperationalActions().Len() > 0 {
		actions := p.getOperationalActions()
		if actions.Len() == 0 {
			return
		}
		action := &pb.MahjongAction{
			OutSeatId:  uint32(actions[0].OutSeatId),
			ActionType: ActionTypeLibToPb(lib.ActionType_Pass),
			OpCard:     actions[0].OutCard.ToUint32(),
		}
		if actions.HasHu() {
			huAction := actions.GetHuAction()
			action = &pb.MahjongAction{
				OutSeatId:  uint32(huAction.OutSeatId),
				ActionType: ActionTypeLibToPb(lib.ActionType_Hu),
				OpCard:     huAction.OutCard.ToUint32(),
			}
		}

		req := &pb.GamePlayerActionRequest{
			SeatId: p.SeatId,
			Action: action,
		}
		d.handleOperateReq(req)
	}
}

// 机器人操作
func (d *Desk) onRobotOperate(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	if !player.IsRobot {
		return
	}
	actions := player.getOperationalActions()
	if actions.Len() == 0 {
		return
	}
	action := &pb.MahjongAction{
		OutSeatId:  uint32(actions[0].OutSeatId),
		ActionType: ActionTypeLibToPb(lib.ActionType_Pass),
		OpCard:     actions[0].OutCard.ToUint32(),
	}

	req := &pb.GamePlayerActionRequest{
		SeatId: seatId,
		Action: action,
	}
	d.handleOperateReq(req)
}

// ActionTypePbToLib 转换pb动作类型到lib
func ActionTypePbToLib(actionType pb.ActionType) lib.ActionType {
	switch actionType {
	case pb.ActionType_ACTION_PASS:
		return lib.ActionType_Pass
	case pb.ActionType_ACTION_CHI:
		return lib.ActionType_Chi
	case pb.ActionType_ACTION_PENG:
		return lib.ActionType_Peng
	case pb.ActionType_ACTION_GANG:
		return lib.ActionType_Gang
	case pb.ActionType_ACTION_HU:
		return lib.ActionType_Hu
	case pb.ActionType_ACTION_TING:
		return lib.ActionType_Ting
	default:
		return lib.ActionType_Unknown
	}
}

func ActionTypeLibToPb(actionType lib.ActionType) pb.ActionType {
	switch actionType {
	case lib.ActionType_Pass:
		return pb.ActionType_ACTION_PASS
	case lib.ActionType_Chi:
		return pb.ActionType_ACTION_CHI
	case lib.ActionType_Peng:
		return pb.ActionType_ACTION_PENG
	case lib.ActionType_Gang:
		return pb.ActionType_ACTION_GANG
	case lib.ActionType_Hu:
		return pb.ActionType_ACTION_HU
	case lib.ActionType_Ting:
		return pb.ActionType_ACTION_TING
	default:
		return pb.ActionType_ACTION_INVALID
	}
}

func ExtActionTypePbToLib(extActionType pb.ActionExtType) lib.ExtraActionType {
	switch extActionType {
	case pb.ActionExtType_ACTION_EXT_MING:
		return lib.ExtraActionType_Ming_Gang
	case pb.ActionExtType_ACTION_EXT_BU:
		return lib.ExtraActionType_Bu_Gang
	case pb.ActionExtType_ACTION_EXT_AN:
		return lib.ExtraActionType_An_Gang
	default:
		return lib.ExtraActionType_Null
	}
}

func ExtActionTypeLibToPb(extActionType lib.ExtraActionType) pb.ActionExtType {
	switch extActionType {
	case lib.ExtraActionType_Ming_Gang:
		return pb.ActionExtType_ACTION_EXT_MING
	case lib.ExtraActionType_Bu_Gang:
		return pb.ActionExtType_ACTION_EXT_BU
	case lib.ExtraActionType_An_Gang:
		return pb.ActionExtType_ACTION_EXT_AN
	default:
		return pb.ActionExtType_ACTION_EXT_NULL
	}
}
