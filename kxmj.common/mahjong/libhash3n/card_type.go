package libhash3n

type CardType int

const (
	Card_Wan CardType = iota
	Card_Tiao
	Card_Tong
	Card_Feng
	Card_Arrow
)

func GetCardLen(color CardType) int {
	switch color {
	case Card_Wan, Card_Tiao, Card_Tong:
		return 9
	case Card_Feng:
		return 4
	case Card_Arrow:
		return 3
	}
	return 9
}
