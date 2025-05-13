package libhash3n

import lib "kxmj.common/mahjong"

// 3n+2中的头
type Tou struct {
	card1 lib.Card
	card2 lib.Card
}

func NewTou(card1 lib.Card, card2 lib.Card) *Tou {
	tou := &Tou{card1: card1, card2: card2}
	return tou
}

func (tou *Tou) GetKingCount() int {
	count := 0
	if tou.card1 == lib.INVALID_CARD {
		count++
	}
	if tou.card2 == lib.INVALID_CARD {
		count++
	}
	return count
}
