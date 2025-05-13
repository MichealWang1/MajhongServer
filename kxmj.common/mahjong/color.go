package lib

type CardColor byte

const (
	Card_Color_Wan     CardColor = 0
	Card_Color_Tiao    CardColor = 1
	Card_Color_Tong    CardColor = 2
	Card_Color_Zi      CardColor = 3
	Card_Color_Hua     CardColor = 4
	Card_Color_Unknown CardColor = 5
)

func (c CardColor) IsWan() bool {
	return c == Card_Color_Wan
}

func (c CardColor) IsTiao() bool {
	return c == Card_Color_Tiao
}

func (c CardColor) IsTong() bool {
	return c == Card_Color_Tong
}

func (c CardColor) IsZi() bool {
	return c == Card_Color_Zi
}

func (c CardColor) IsHua() bool {
	return c == Card_Color_Hua
}
