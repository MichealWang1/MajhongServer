package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/utils"
	"kxmj.game.mjxlch/pb"
	"strconv"
)

type InningScores struct {
	handCards      lib.Cards         // 玩家手牌
	userActions    lib.UserActions   //玩家动作牌
	settlementType pb.SettlementType // 结算类型
	opCard         lib.Card          // 操作牌
	winSeatId      uint32            // 赢分玩家
	scores         []string          // 各位子得分
	winLoseGolds   []string          // 真实输赢金币
	huResult       *HuResult         // 胡牌权位
	multiple       string            // 倍数
	isCeiling      bool              // 是否封顶
	isBankruptcy   []bool            // 破产
}

func NewInningScores(handCards lib.Cards, actions lib.UserActions, settlementType pb.SettlementType, opCard lib.Card, winSeatId uint32, huResult *HuResult, multiple string, isCeiling bool) *InningScores {
	res := &InningScores{
		handCards:      handCards.Copy(),
		userActions:    actions.Copy(),
		settlementType: settlementType,
		opCard:         opCard,
		winSeatId:      winSeatId,
		scores:         make([]string, PLAYER_COUNT),
		winLoseGolds:   make([]string, PLAYER_COUNT),
		huResult:       huResult,
		multiple:       multiple,
		isCeiling:      isCeiling,
		isBankruptcy:   make([]bool, PLAYER_COUNT),
	}
	for i := uint32(0); i < PLAYER_COUNT; i++ {
		res.scores[i] = "0"
		res.winLoseGolds[i] = "0"
	}
	return res
}

func (i *InningScores) Copy() *InningScores {
	res := &InningScores{
		handCards:      i.handCards.Copy(),
		userActions:    i.userActions.Copy(),
		settlementType: i.settlementType,
		opCard:         i.opCard,
		winSeatId:      i.winSeatId,
		scores:         make([]string, PLAYER_COUNT),
		winLoseGolds:   make([]string, PLAYER_COUNT),
		multiple:       i.multiple,
		isCeiling:      i.isCeiling,
		isBankruptcy:   make([]bool, PLAYER_COUNT),
	}
	if i.huResult != nil {
		res.huResult = i.huResult.Copy()
	}
	copy(res.scores, i.scores)
	copy(res.winLoseGolds, i.winLoseGolds)
	copy(res.isBankruptcy, i.isBankruptcy)
	return res
}

// 计算*底分
func (s *InningScores) MulBaseScores(baseScore string) {
	for i, v := range s.scores {
		res, _ := utils.MulToString(v, baseScore)
		s.scores[i] = res
	}
}

// 计算杠分
func (d *Desk) calculateGangScore(action *lib.UserAction) *InningScores {
	player := d.getPlayerBySeat(uint32(action.SeatId))
	if action.ActionType != lib.ActionType_Gang {
		return nil
	}
	inningScores := NewInningScores(player.getHandCards(), player.getActionCardsTable(), pb.SettlementType_SETTLEMENT_INVALID, action.OutCard, uint32(action.SeatId), nil, "1", false)

	winSeatId := action.SeatId
	loseSeatId := action.OutSeatId
	mult := "1"

	switch action.ExtraActionType {
	case lib.ExtraActionType_Ming_Gang:
		// 计算倍数
		mult = "2"
		inningScores.settlementType = pb.SettlementType_MING_GANG
		// 计算每个玩家出多少
		inningScores.scores[int(winSeatId)] = mult
		inningScores.scores[int(loseSeatId)], _ = utils.MulToString(mult, "-1")
	case lib.ExtraActionType_Bu_Gang:
		// 计算倍数
		mult = "1"
		inningScores.settlementType = pb.SettlementType_BU_GANG
		// 计算每个玩家出多少
		inningScores.scores[int(winSeatId)], _ = utils.MulToString(mult, "3")
		for _, p := range d.players {
			if p.SeatId == uint32(winSeatId) {
				continue
			}
			inningScores.scores[int(p.SeatId)], _ = utils.MulToString(mult, "-1")
		}
	case lib.ExtraActionType_An_Gang:
		// 计算倍数
		mult = "2"
		inningScores.settlementType = pb.SettlementType_AN_GANG
		// 计算每个玩家出多少
		inningScores.scores[int(winSeatId)], _ = utils.MulToString(mult, "3")
		for _, p := range d.players {
			if p.SeatId == uint32(winSeatId) {
				continue
			}
			inningScores.scores[int(p.SeatId)], _ = utils.MulToString(mult, "-1")
		}
	}
	inningScores.multiple = mult
	inningScores.MulBaseScores(d.room.BaseScore())
	// 以小博大
	d.calculateYiXiaoBoDa(inningScores)

	// 计算金币
	d.calculateWinLoseGold(inningScores)
	log.Sugar().Infof("score: %v,winloseGold:%v", inningScores.scores, inningScores.winLoseGolds)
	return inningScores
}

// 计算胡分
func (d *Desk) calculateHuScores(winSeatId, outSeatId uint32, score string, opCard lib.Card, huResult *HuResult) *InningScores {
	log.Sugar().Infof("winSeatId=%d outSeatId=%d score= %v opCard=%v huResult=%#v,maxMultiple=%v", winSeatId, outSeatId, score, opCard, huResult, d.room.MaxMultiple())
	player := d.getPlayerBySeat(winSeatId)
	isZiMo := winSeatId == outSeatId

	// 封顶处理
	isCeiling := false
	if d.room.MaxMultiple() != 0 {
		isCeiling = utils.Cmp(strconv.Itoa(int(d.room.MaxMultiple())), score) == -1
		if isCeiling {
			score = strconv.Itoa(int(d.room.MaxMultiple()))
		}
	}

	// 初始化分数
	inningScores := NewInningScores(player.getHandCards(), player.getActionCardsTable(), pb.SettlementType_HU, opCard, winSeatId, huResult.Copy(), score, isCeiling)
	// 自摸
	if isZiMo {
		inningScores.scores[int(winSeatId)], _ = utils.MulToString(score, "3")
		for _, p := range d.players {
			if p.SeatId == winSeatId {
				continue
			}
			inningScores.scores[int(p.SeatId)], _ = utils.MulToString(score, "-1")
		}
	} else {
		inningScores.scores[int(winSeatId)] = score
		inningScores.scores[int(outSeatId)], _ = utils.MulToString(score, "-1")
	}

	log.Sugar().Infof("inningScores:%v,baseScore: %v", inningScores.scores, d.room.BaseScore())
	inningScores.MulBaseScores(d.room.BaseScore())

	// 以小博大
	d.calculateYiXiaoBoDa(inningScores)
	// 计算金币
	d.calculateWinLoseGold(inningScores)
	log.Sugar().Infof("score: %v,winloseGold:%v", inningScores.scores, inningScores.winLoseGolds)

	return inningScores
}

// 以小博大
func (d *Desk) calculateYiXiaoBoDa(inningScores *InningScores) {
	winSeatId := inningScores.winSeatId
	winPlayer := d.getPlayerBySeat(winSeatId)
	winScore := "0"
	for i, score := range inningScores.scores {
		if uint32(i) == winSeatId {
			continue
		}
		loseScore, _ := utils.MulToString(score, "-1")
		if utils.Cmp(winPlayer.Gold, loseScore) == -1 {
			loseScore = winPlayer.Gold
		}
		winScore, _ = utils.AddToString(winScore, loseScore)
	}
	inningScores.scores[winSeatId] = winScore
}

// 计算需要结算的金币
func (d *Desk) calculateWinLoseGold(inningScores *InningScores) {
	winCount := 0 // 赢钱玩家有多少
	for _, score := range inningScores.scores {
		if utils.Cmp(score, "0") == 1 {
			winCount++
		} else {
			winCount--
		}
	}

	if winCount < 0 { // 赢家小于0 则为单个玩家赢其他人 [+,-,-,-] winCount = -2,[+,-,0,-] winCount = -2;[+,0,0,-] winCount = -2;
		d.calculateOneWinPlayer(inningScores)
	} else { // 赢家大于等于0 则是一家输多家 [-,+,+,0] winCount =0;[-,+,+,+] winCount =2;
		d.calculateOneLosePlayer(inningScores)
	}
}

// 一个玩家赢
func (d *Desk) calculateOneWinPlayer(inningScores *InningScores) {
	winSeatId := uint32(0)
	hasWin := false
	for i, score := range inningScores.scores {
		if utils.Cmp(score, "0") == 1 {
			winSeatId = uint32(i)
			hasWin = true
			break
		}
	}
	if !hasWin {
		return
	}
	winReduceGold := "0"
	for _, p := range d.players {
		if p.SeatId == winSeatId {
			continue
		}
		log.Sugar().Infof("Gold: %v", p.Gold)
		reduceGold := inningScores.scores[p.SeatId]
		if utils.Cmp("0", inningScores.scores[p.SeatId]) == 1 && utils.Cmp(p.Gold, inningScores.scores[p.SeatId][1:]) == -1 {
			reduceGold, _ = utils.MulToString(p.Gold, "-1")
			inningScores.isBankruptcy[p.SeatId] = true
		}
		winReduceGold, _ = utils.AddToString(winReduceGold, reduceGold[1:])
		inningScores.winLoseGolds[p.SeatId] = reduceGold
		log.Sugar().Infof("loseGold:%v", inningScores.winLoseGolds[p.SeatId])
	}
	inningScores.winLoseGolds[winSeatId] = winReduceGold
	log.Sugar().Infof("winGold:%v,winLoseGolds:%v", winReduceGold, inningScores.winLoseGolds)
}

// 一个玩家输
func (d *Desk) calculateOneLosePlayer(inningScores *InningScores) {
	loserSeatId := uint32(0)
	hasLoser := false
	for i, score := range inningScores.scores {
		if utils.Cmp("0", score) == 1 {
			loserSeatId = uint32(i)
			hasLoser = true
			break
		}
	}
	if !hasLoser {
		return
	}
	loserPlayer := d.getPlayerBySeat(loserSeatId)
	if utils.Cmp(loserPlayer.Gold, inningScores.scores[loserSeatId][1:]) != 1 {
		inningScores.isBankruptcy[loserSeatId] = true
		inningScores.winLoseGolds[loserSeatId], _ = utils.MulToString(loserPlayer.Gold, "-1")
		gold := loserPlayer.Gold
		maxScore := 0
		loseScore, _ := utils.MulToString(inningScores.scores[loserSeatId], "-1")
		for i, score := range inningScores.scores {
			if utils.Cmp("0", score) != -1 {
				continue
			}
			tmp, _ := utils.MulToString(score, loserPlayer.Gold)
			kou, _ := utils.QuoToString(tmp, loseScore)
			if utils.Cmp("0", kou) == 0 {
				kou = "1"
			}
			gold, _ = utils.SubToString(gold, kou)
			inningScores.winLoseGolds[i] = kou
			if utils.Cmp(inningScores.winLoseGolds[maxScore], kou) == -1 {
				maxScore = i
			}
		}
		inningScores.winLoseGolds[maxScore], _ = utils.AddToString(inningScores.winLoseGolds[maxScore], gold)
	} else {
		for i, score := range inningScores.scores {
			inningScores.winLoseGolds[i] = score
		}
	}
}

// 金币结算
func (d *Desk) calculateGold(scores []*InningScores) {
	for _, p := range d.players {
		for _, score := range scores {
			if score.winLoseGolds[p.SeatId] == "0" {
				continue
			}
			gold, _ := utils.AddToString(p.Gold, score.winLoseGolds[p.SeatId])
			p.Gold = gold
			log.Sugar().Infof("gold:%v", gold)
		}
	}
}

// 呼叫转移
func (d *Desk) calculateTransfer(winSeatId, loserSeatId uint32) *InningScores {
	log.Sugar().Infof("calculateTransfer winSeatId:%v,loserSeatId:%v", winSeatId, loserSeatId)
	var gangScores *InningScores
	// 上一次杠牌分
	for _, data := range d.getSettlementData() {
		if data.winSeatId == loserSeatId && data.opCard == d.getPrevActionGangCards()[d.getPrevActionGangCards().Len()-1] {
			gangScores = data.Copy()
			break
		}
	}
	log.Sugar().Infof("calculateTransfer gangScores:%v", gangScores)
	winPlayer := d.getPlayerBySeat(winSeatId)
	// 初始化分数
	inningScores := NewInningScores(winPlayer.getHandCards(), winPlayer.getActionCardsTable(), pb.SettlementType_ZHUANG_YI, gangScores.opCard, winSeatId, nil, gangScores.multiple, false)
	inningScores.scores[winSeatId] = gangScores.scores[loserSeatId]
	inningScores.scores[loserSeatId], _ = utils.MulToString("-", gangScores.scores[loserSeatId])
	// 计算金币
	d.calculateWinLoseGold(inningScores)
	log.Sugar().Infof("score: %v,winloseGold:%v", inningScores.scores, inningScores.winLoseGolds)

	return inningScores
}

// 查大叫、查花猪、退税
//func (d *Desk) calculateEndScores() []*InningScores {
//	// 判断哪些玩家没听牌、哪些玩家听牌
//	tingCards := make([]lib.Cards, PLAYER_COUNT)           // 玩家听的牌
//	haveChooseMissCardStatus := make([]bool, PLAYER_COUNT) // 是否是花猪玩家
//	hasStatus := false                                     // 是否要进行退税、查叫、查花猪
//	// 没听牌玩家退回全部杠牌
//	// 计算听牌玩家最大可能的倍数
//	// 没听牌玩家赔付听牌玩家最大倍数
//	for _, p := range d.players {
//		handCards := p.getHandCards().Copy()
//		if handCards.GetAllColorCards(ChooseMissTypeToLib(p.getChooseMissType())).Len() != 0 {
//			tingCards[p.SeatId] = lib.Cards{}
//			haveChooseMissCardStatus[p.SeatId] = true
//			hasStatus = true
//			continue
//		}
//		cards := d.analysisTingCards(handCards)
//		if cards.Len() == 0 {
//			hasStatus = true
//		}
//		tingCards[p.SeatId] = cards
//	}
//
//	if !hasStatus {
//		return nil
//	}
//
//	//
//	//for _, p := range d.players {
//	//}
//
//}

// 退税
func (d *Desk) calculateTuiShui(seatId uint32) *InningScores {
	player := d.getPlayerBySeat(seatId)
	inningScores := NewInningScores(player.getHandCards(), player.getActionCardsTable(), pb.SettlementType_TUI_SHUI, lib.INVALID_CARD, SEAT_UNKNOWN, nil, "1", false)
	// 查找该位子的杠牌
	for _, data := range d.getSettlementData() {
		if data.winSeatId == seatId && (data.settlementType == pb.SettlementType_MING_GANG || data.settlementType == pb.SettlementType_BU_GANG || data.settlementType == pb.SettlementType_AN_GANG) { // 没听牌肯定没胡牌,所以只要找到这个玩家有的结算分就是杠分了
			for _, p := range d.players {
				scores, _ := utils.MulToString(data.scores[p.SeatId], "-1")
				inningScores.scores[p.SeatId], _ = utils.AddToString(inningScores.scores[p.SeatId], scores)

				winGolds, _ := utils.MulToString(data.winLoseGolds[p.SeatId], "-1")
				inningScores.winLoseGolds[p.SeatId], _ = utils.AddToString(inningScores.winLoseGolds[p.SeatId], winGolds)
			}

		}
	}
	// 计算金币
	d.calculateWinLoseGold(inningScores)
	log.Sugar().Infof("score: %v,winloseGold:%v", inningScores.scores, inningScores.winLoseGolds)
	return inningScores
}

// 同一玩家输多个分数处理
func (d *Desk) calculateLoseMostScores(inningScores []*InningScores) {
	for _, p := range d.players {
		winLoseScore := "0"
		for _, data := range inningScores {
			winLoseScore, _ = utils.AddToString(winLoseScore, data.scores[p.SeatId])
		}
		if utils.Cmp("0", winLoseScore) == 1 {
			if utils.Cmp(winLoseScore[1:], p.Gold) == 1 {
				gold := p.Gold
				maxScore := 0
				for i, data := range inningScores {
					if utils.Cmp(data.scores[p.SeatId], "0") == 0 {
						inningScores[i].winLoseGolds[p.SeatId] = "0"
						continue
					}
					inningScores[i].isBankruptcy[p.SeatId] = true
					tmp, _ := utils.MulToString(data.scores[p.SeatId], p.Gold)
					kou, _ := utils.QuoToString(tmp, winLoseScore)
					kou, _ = utils.MulToString(kou, "-1")
					if utils.Cmp("0", kou) == 0 {
						kou = "-1"
					}
					gold, _ = utils.AddToString(gold, kou)
					inningScores[i].winLoseGolds[p.SeatId] = kou
					if utils.Cmp(kou, inningScores[maxScore].winLoseGolds[p.SeatId]) == -1 {
						maxScore = i
					}
				}
				inningScores[maxScore].winLoseGolds[p.SeatId], _ = utils.SubToString(inningScores[maxScore].winLoseGolds[p.SeatId], gold)
			} else {
				for i, data := range inningScores {
					inningScores[i].winLoseGolds[p.SeatId] = data.scores[p.SeatId]
				}
			}
		}
	}

	for i, data := range inningScores {
		loseGold := "0"
		winGold := "0"
		for j, score := range data.winLoseGolds {
			if utils.Cmp("0", score) == -1 {
				winGold, _ = utils.AddToString(winGold, inningScores[i].scores[j])
				continue
			}
			loseGold, _ = utils.AddToString(loseGold, score)
		}
		gold, _ := utils.MulToString(loseGold, "-1")
		maxScore := 0
		for j, score := range data.winLoseGolds {
			if utils.Cmp("0", score) == -1 {
				tmp, _ := utils.MulToString(score, loseGold)
				kou, _ := utils.QuoToString(tmp, winGold)
				kou, _ = utils.MulToString(kou, "-1")
				if utils.Cmp("0", kou) == 0 {
					kou = "1"
				}
				gold, _ = utils.SubToString(gold, kou)
				inningScores[i].winLoseGolds[j] = kou
				if utils.Cmp(inningScores[i].winLoseGolds[maxScore], kou) == -1 {
					maxScore = i
				}
			}
		}
		inningScores[i].winLoseGolds[maxScore], _ = utils.AddToString(inningScores[i].winLoseGolds[maxScore], gold)
		log.Sugar().Infof("calculateLoseMostScores end inningScores:%v", inningScores[i])
	}
}
