package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/mahjong/libhash3n"
	"kxmj.game.mjxlch/pb"
)

// 3n+2听牌
func (d *Desk) calculateTingCards3N2(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	handCards := player.getHandCards().Copy()
	// 不是3n+2牌型不检测
	if handCards.Len()%3 != 2 || player.isHu() {
		return
	}
	tingCardsInfos := make([]*pb.TingCardInfo, 0)

	// 去重
	uniqueCards := handCards.ToUnique()
	for _, card := range uniqueCards {
		newCards, _ := handCards.DeleteCard(card)
		// 听牌判断
		tingCards := d.getTingCardsInfo(player.SeatId, newCards)
		if len(tingCards) != 0 {
			tingInfo := &pb.TingCardInfo{
				OutCard:  card.ToUint32(),
				TingCard: tingCards,
			}
			tingCardsInfos = append(tingCardsInfos, tingInfo)
		}
	}
	if len(tingCardsInfos) > 0 {
		pbResponse := &pb.GameTing3N2InfoNotify{
			TingCards: tingCardsInfos,
		}
		log.Sugar().Infof("ProGameTing3N2InfoNotify pbResponse:", pbResponse)
		d.game.SendMessage(player.UserId, uint16(pb.MID_TING_3N2_NOTIFY), pbResponse)
	}
}

// 3n+1听牌
func (d *Desk) calculateTingCards3N1(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	handCards := player.getHandCards().Copy()
	if handCards.Len()%3 != 1 {
		return
	}
	// 听牌判断
	tingCards := d.getTingCardsInfo(player.SeatId, handCards)
	if len(tingCards) != 0 {
		// 获取听的牌的数量
		pbResponse := &pb.GameTing3N1InfoNotify{
			TingCard: tingCards,
		}
		log.Sugar().Infof("GameTing3N1InfoNotify pbResponse:", pbResponse)
		d.game.SendMessage(player.UserId, uint16(pb.MID_TING_3N1_NOTIFY), pbResponse)
	}
}

// 玩家的听牌信息
func (d *Desk) getTingCardsInfo(seatId uint32, handCards lib.Cards) []*pb.TingCard {
	res := make([]*pb.TingCard, 0)
	// 手牌有缺牌不听牌
	if handCards.GetAllColorCards(ChooseMissTypeToLib(d.getPlayerBySeat(seatId).getChooseMissType())).Len() != 0 {
		return res
	}
	// 判断哪些牌能听
	tingCards := d.analysisTingCards(handCards)
	for _, tCard := range tingCards {
		cardCount := &pb.TingCard{
			Card:     tCard.ToUint32(),
			Count:    uint32(d.getCardNowResidueCount(seatId, tCard)),
			Multiple: d.analysisTingMultiple(seatId, handCards, tCard),
		}
		res = append(res, cardCount)
	}
	return res
}

// 分析玩家听的牌
func (d *Desk) analysisTingCards(handCards lib.Cards) lib.Cards {
	tingCards := make(lib.Cards, 0)
	analyzer := libhash3n.NewTingHuAnalyzer(handCards, lib.Cards{}, false, false)
	cards := analyzer.GetAllTingCards()

	// 过滤牌
	for _, card := range cards {
		tingCards = tingCards.AddCard(card)
	}
	return tingCards
}

// 分析玩家听牌倍数
func (d *Desk) analysisTingMultiple(seatId uint32, handCards lib.Cards, huCard lib.Card) string {
	player := d.getPlayerBySeat(seatId)
	analysisHu := &AnalysisHu{
		SeatId:           seatId,
		OutSeatId:        SEAT_UNKNOWN,
		BankerSeatId:     d.runtimeData.bankerSeatId,
		OpCard:           huCard,
		HandCards:        handCards,
		ChooseMissColor:  ChooseMissTypeToLib(player.getChooseMissType()),
		ActionCards:      player.getActionCardsTable(),
		OutCount:         player.getOutCount(),
		IsQiangGang:      false,
		PrevActionIsGang: d.getPrevActionIsGang() > 0,
		IsHaiDi:          false,
	}
	huResult := NewHuResult()
	analysisHuRights(analysisHu, huResult)
	// 计算胡牌倍数
	return huResult.multiple
}

// 获取牌在当前玩家视角下还剩多少张
func (d *Desk) getCardNowResidueCount(seatId uint32, card lib.Card) int {
	// 所有玩家弃牌、动作牌、胡牌数据和自己手牌里里这张牌数量
	n := 0
	for _, p := range d.players {
		if p.SeatId == seatId { // 自己手牌
			n += p.getHandCards().GetCount(card)
		}
		// 弃牌
		n += p.getDiscards().GetCount(card)
		// 动作牌
		for _, action := range p.getActionCardsTable() {
			n += action.CombineCards.GetCount(card)
		}
		// 胡牌数据
		for _, data := range p.getHuData() {
			if data.opCard == card {
				n++
			}
		}
	}
	return 4 - n
}
