package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
	"time"
)

// 判断是否选缺结束
func (d *Desk) chooseMissStateCanOver() bool {
	res := true
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		if p.getChooseMissType() == pb.MissType_MISS_NULL {
			res = false
			break
		}
	}
	return res
}

// 时间到了处理没完成选缺策略
func (d *Desk) onChooseMissTimeOver() {
	for _, p := range d.players {
		if p.getChooseMissType() != pb.MissType_MISS_NULL {
			continue
		}
		// 设置玩家选缺
		p.setChooseMissType(getHandSmallColor(p.getHandCards()))
		// 广播玩家选缺响应
		d.broadcastPlayerChooseMissResponse()
		time.Sleep(1 * time.Second)
	}
	// 设置庄家摸的牌为最右边牌
	bankerPlayer := d.getPlayerBySeat(d.runtimeData.bankerSeatId)
	handCards := bankerPlayer.getHandCards().Copy()
	handCards.Sort()
	if cards := handCards.GetAllColorCards(ChooseMissTypeToLib(bankerPlayer.getChooseMissType())); cards.Len() > 0 {
		cards.Sort()
		bankerPlayer.setCatchCard(cards[cards.Len()-1])
	} else {
		bankerPlayer.setCatchCard(handCards[handCards.Len()-1])
	}
	log.Sugar().Infof("banker handCards:%v,catch card :%v", bankerPlayer.getHandCards(), bankerPlayer.getCatchCard())
}

// 机器人选缺
func (d *Desk) onRobotChooseMiss() {

}

// 选缺选花色牌最少的
func getHandSmallColor(cards lib.Cards) pb.MissType {
	wanCards := cards.GetAllColorCards(lib.Card_Color_Wan)
	tiaoCards := cards.GetAllColorCards(lib.Card_Color_Tiao)
	tongCards := cards.GetAllColorCards(lib.Card_Color_Tong)
	res := pb.MissType_MISS_WAN

	if tiaoCards.Len() < wanCards.Len() {
		res = pb.MissType_MISS_TIAO
		if tongCards.Len() < tiaoCards.Len() {
			res = pb.MissType_MISS_TONG
		}
	} else if tongCards.Len() < wanCards.Len() {
		res = pb.MissType_MISS_TONG
	}
	return res
}
