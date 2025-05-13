package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
	"math/rand"
	"sort"
	"time"
)

// 随机换牌类型
func (d *Desk) randomSwapType() {
	rand.Seed(time.Now().UnixMilli())
	d.runtimeData.swapType = pb.SwapType(rand.Intn(3) + 1)
}

// 是否所有人都换完牌了
func (d *Desk) swapStateCanOver() bool {
	res := true
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		if p.getSwapCards().Len() == 0 {
			res = false
			break
		}
	}
	return res
}

// 进行换牌
func (d *Desk) swapStart() {
	for _, p := range d.players {
		handCards := p.getHandCards()
		swapCards := make(lib.Cards, 0, 3)
		switch d.runtimeData.swapType {
		case pb.SwapType_SWAP_TYPE_NEXT: // 换给下一个玩家
			swapCards = d.getPrevSeatPlayer(p.SeatId).getSwapCards()
		case pb.SwapType_SWAP_TYPE_PREV: // 换给上一个玩家
			swapCards = d.getNextSeatPlayer(p.SeatId).getSwapCards()
		case pb.SwapType_SWAP_TYPE_OPPOSITE: // 换给对家
			swapCards = d.getOppositeSeatPlayer(p.SeatId).getSwapCards()
		}
		handCards = handCards.AddCards(swapCards)
		p.setHandCards(handCards)
		// 通知换牌类型
		d.sendPlayerSwapTypeNotify(p.SeatId, swapCards)
		log.Sugar().Infof("seatId:%d,swapCards:%v,handCards:%v", p.SeatId, swapCards, p.getHandCards())
	}
	// 如果把摸得牌换了处理
	//opPlayer := d.getPlayerBySeat(d.getOperateSeatId())
	//if !opPlayer.getHandCards().In(opPlayer.getCatchCard()) {
	//	opPlayer.setCatchCard(opPlayer.getHandCards()[opPlayer.getHandCards().Len()-1])
	//}

	//d.broadcastAllHandCardsResponse(pb.UpdateMahjongType_UPDATE_SWAP)
}

// 处理时间到了玩家还没进行换牌策略
func (d *Desk) onSwapTimeOver() {
	for _, p := range d.players {
		//if !p.IsRobot {
		//	continue
		//}
		if p.getSwapCards().Len() != 0 {
			continue
		}
		swapCards := GetSmallColorCards(p.getHandCards(), 3)
		err := p.deleteHandCards(swapCards)
		if err != nil {
			log.Sugar().Errorf("onSwapTimeOver error: %v", err)
		}
		p.setSwapCards(swapCards)
		// 广播玩家换牌响应
		d.broadcastPlayerSwapResponse(p.SeatId)
		// 更新玩家手牌
		//d.broadcastUserCardsResponse(p.SeatId, pb.UpdateMahjongType_UPDATE_SWAP)
		time.Sleep(1 * time.Second)
	}
}

// 机器人换牌策略
func (d *Desk) onRobotSwap() {

}

func GetSmallColorCards(cards lib.Cards, n int) lib.Cards {
	swapCards := make(lib.Cards, 0, n)
	colorCards := make([]lib.Cards, 3)
	colorCards[lib.Card_Color_Wan] = cards.GetAllColorCards(lib.Card_Color_Wan)
	colorCards[lib.Card_Color_Tiao] = cards.GetAllColorCards(lib.Card_Color_Tiao)
	colorCards[lib.Card_Color_Tong] = cards.GetAllColorCards(lib.Card_Color_Tong)
	sort.Slice(colorCards, func(a, b int) bool {
		return colorCards[a].Len() < colorCards[b].Len()
	})
	for _, c := range colorCards {
		if c.Len() == 0 {
			continue
		}
		for i := 0; swapCards.Len() < n && i < c.Len(); i++ {
			swapCards = append(swapCards, c[i])
		}
	}
	return swapCards
}
