package libhash3n

import lib "kxmj.common/mahjong"

const Card_Type_Mask int = 0xF0
const Card_Value_Mask int = 0x0F

// card to 索引值
func CardToIndex(card lib.Card) int {
	cardInt := int(card)
	color := (cardInt & Card_Type_Mask) >> 4
	value := (cardInt & Card_Value_Mask) - 1
	return color*9 + value
}

// 索引值 to card
func IndexToCard(index int) lib.Card {
	colorIndex := index / 9
	valueIndex := index%9 + 1
	cardByte := (colorIndex << 4) + valueIndex
	return lib.Card(cardByte)
}
