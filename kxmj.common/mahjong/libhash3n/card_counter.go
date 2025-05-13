package libhash3n

import lib "kxmj.common/mahjong"

type CardCounter []uint

// card转成cardcounter
func ToCardCounter(cards lib.Cards) CardCounter {
	cardCounter := make([]uint, 34)
	for _, card := range cards {
		index := CardToIndex(card)
		cardCounter[index]++
	}
	return cardCounter
}
