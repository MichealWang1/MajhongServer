package libhash3n

import lib "kxmj.common/mahjong"

var mahjongsMap map[lib.Card]int
var (
	mahjongs = lib.Cards{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
		0x31, 0x32, 0x33, 0x34,
		0x35, 0x36, 0x37,
	} // 所有万的定义
	wanCards   = lib.Cards{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09} // 所有万的定义
	tiaoCards  = lib.Cards{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19} // 所有条的定义
	tongCards  = lib.Cards{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29} // 所有筒的定义
	fengCards  = lib.Cards{0x31, 0x32, 0x33, 0x34}                               // 所有字牌定义
	arrowCards = lib.Cards{0x35, 0x36, 0x37}
)

func init() {
	mahjongsMap = make(map[lib.Card]int)
	for _, card := range mahjongs {
		mahjongsMap[card] = 1
	}
}

// card是否是库接口认可的精牌
func isValidKingCard(card lib.Card) bool {
	if card == lib.INVALID_CARD {
		return true
	}
	if _, exist := mahjongsMap[card]; exist {
		return true
	}
	return false
}

// card是否是库接口认可的card
func isValidHandCard(card lib.Card) bool {
	if _, exist := mahjongsMap[card]; exist {
		return true
	}
	return false
}
