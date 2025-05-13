package logic

import (
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/internal/game"
	"math/rand"
	"time"
)

type DiceStateData struct {
	bankerSeatId uint32    // 庄家位子
	kingCards    lib.Cards // 精牌信息
}

func NewDiceStateData() *DiceStateData {
	d := &DiceStateData{}
	d.Reset()
	return d
}

// Reset 重置
func (d *DiceStateData) Reset() {
	d.bankerSeatId = game.SEAT_UNKNOWN
	d.kingCards = make(lib.Cards, 0)
}

// GetBankerSeatId 获取庄家位子
func (d *DiceStateData) GetBankerSeatId() uint32 {
	return d.bankerSeatId
}

// 获取精牌信息
func (d *DiceStateData) GetKingCards() lib.Cards {
	return d.kingCards
}

// DetermineBankerSeatId 定庄(随机)
func (d *DiceStateData) DetermineBankerSeatId(userNum int) {
	rand.Seed(time.Now().UnixMilli())
	d.bankerSeatId = rand.Uint32() % uint32(userNum)
}
