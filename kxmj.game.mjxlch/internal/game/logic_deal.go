package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"math/rand"
	"time"
)

// 游戏开始发牌
func (d *Desk) dealCards() {
	for _, p := range d.players {
		n := 13
		if p.SeatId == d.runtimeData.bankerSeatId {
			n = 14
		}
		cards, err := d.runtimeData.cardStack.PickUpCards(n) // 每人拿13张,庄家拿14张
		if err != nil {
			log.Sugar().Errorf("deal hand cards error %v", err)
			return
		}
		p.setHandCards(cards)
		log.Sugar().Infof("player:%v,handCards:%v", p, p.getHandCards())
	}
	// 设置庄家起手摸牌
	bankerPlayer := d.getPlayerBySeat(d.runtimeData.bankerSeatId)
	bankerPlayer.setCatchCard(bankerPlayer.getHandCards()[bankerPlayer.getHandCards().Len()-1])
	// 设置操作位子
	d.setOperateSeatId(bankerPlayer.SeatId)

	d.broadcastDealHandCards()
}

// 玩家摸牌
func (d *Desk) catchCard(seatId uint32) error {
	player := d.getPlayerBySeat(seatId)
	catchCard, err := d.runtimeData.cardStack.Put()
	if err != nil {
		// 没有牌了进入结算
		d.setCanOver(true)
		return err
	}
	player.setCatchCard(catchCard)
	player.addHandCards(catchCard)
	d.broadcastPlayerCatchCard(seatId, catchCard)
	return nil
}

// 做牌
func (d *Desk) setCardStack() {
	if TEST_DECK == false {
		return
	}

	bankerSeatId := uint32(0)
	cards := lib.Cards{
		0x01, 0x02, 0x03, 0x01, 0x02, 0x03, 0x01, 0x02, 0x03, 0x01, 0x02, 0x07, 0x07, 0x08, // 0号位手牌
		0x11, 0x11, 0x11, 0x11, 0x12, 0x12, 0x12, 0x12, 0x13, 0x13, 0x03, 0x04, 0x04, // 1号位手牌
		0x21, 0x21, 0x21, 0x21, 0x22, 0x22, 0x22, 0x22, 0x23, 0x23, 0x13, 0x13, 0x14, // 2号位手牌
		0x05, 0x05, 0x05, 0x05, 0x06, 0x06, 0x06, 0x06, 0x07, 0x07, 0x23, 0x23, 0x24, // 3号位手牌
	}

	deleteCards, _ := AllCards.DeleteCards(cards)
	// 生成随机种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 洗牌
	for i := len(deleteCards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		deleteCards[i], deleteCards[j] = deleteCards[j], deleteCards[i]
	}
	newCards := cards.AddCards(deleteCards)
	log.Sugar().Infof("newCards:%#v", newCards)
	// 设置庄家位子
	d.runtimeData.bankerSeatId = bankerSeatId
	// 设置牌堆
	d.runtimeData.cardStack.Set(newCards)
}
