package lib

type Card byte

const (
	INVALID_CARD Card = 0xff
)

func (c Card) ToInt32() int32 {
	return int32(c)
}

func (c Card) ToUint32() uint32 {
	return uint32(c)
}

// GetColor 获取牌花色
func (c *Card) GetColor() CardColor {
	return CardColor(c.ToInt32() >> 4)
}

// GetValue 获取牌值
func (c *Card) GetValue() byte {
	return byte(c.ToInt32() & 0x0f)
}

// IsValid 牌是否有效
func (c *Card) IsValid() bool {
	if c.GetColor() > Card_Color_Hua || c.GetValue() > byte(9) {
		return false
	}
	return true
}

// PrevCard 上一张牌
func (c *Card) PrevCard() Card {
	if c.GetValue() == byte(1) {
		if c.GetColor().IsZi() {
			return NewCard(c.GetColor(), byte(7))
		}
		return NewCard(c.GetColor(), byte(9))
	}
	return NewCard(c.GetColor(), byte(c.GetValue()-1))
}

// NextCard 下一张牌
func (c *Card) NextCard() Card {
	if c.GetValue() == byte(9) {
		return NewCard(c.GetColor(), byte(1))
	}
	if c.GetColor().IsZi() && c.GetValue() == byte(7) {
		return NewCard(c.GetColor(), byte(1))
	}
	return NewCard(c.GetColor(), byte(c.GetValue()+1))
}

// NewCard 生成牌
func NewCard(color CardColor, value byte) Card {
	return Card(byte(color)<<4 | value)
}

// Repeat 生成重复的牌
func (c *Card) Repeat(n int) Cards {
	r := make(Cards, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, NewCard(c.GetColor(), c.GetValue()))
	}
	return r
}
