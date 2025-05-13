package libhash3n

var PAI_GROUP_MAP_INT_COUNT_BITS []int = []int{
	0x0,
	0x1,
	0x3,
	0x7,
	0xF,
	0x1F,
	0x3F,
	0x7F,
}

func CardCounterToInt(color CardType, cardCounter CardCounter) int {
	cardCounterLen := len(cardCounter)
	if GetCardLen(color) != cardCounterLen {
		panic("cardCounter len err")
	}
	var returnInt int = 0
	for i := cardCounterLen - 1; i >= 0; i-- {
		leftShiftCount := cardCounter[i]
		if leftShiftCount != 0 {
			returnInt = returnInt << leftShiftCount
			returnInt = returnInt | PAI_GROUP_MAP_INT_COUNT_BITS[leftShiftCount]
		}
		if i != 0 {
			returnInt <<= 1
		}
	}
	return returnInt
}

// cardCounter 掩码
const Card_Counter_Mask int = 0xFFFFF

// threen数 to counter数
func ThreeNIntToCounterInt(threeNInt int) int {
	return threeNInt & Card_Counter_Mask
}

// 通过threen数获得是否是3n结构
func Is3NByInt(cardType CardType, threeNInt int) bool {
	return GetMin3nNeedKingCount(cardType, threeNInt) == 0
}

// 万条筒最小精个数掩码及右移数
const wnt_Min_3n_king_count_mask = 0x78000000
const wnt_Min_3n_king_count_Shift_count uint = 27

// 风最小精个数掩码及右移数
const feng_Min_3n_king_count_mask = 0x7800000
const feng_Min_3n_king_count_Shift_count uint = 23

// 箭最小精个数掩码及右移数
const arrow_Min_3n_king_count_mask = 0x700000
const arrow_Min_3n_king_count_Shift_count uint = 20

// 获取threen数达到3n结构需要的精个数
func GetMin3nNeedKingCount(cardType CardType, threeNInt int) int {
	switch cardType {
	case Card_Wan, Card_Tiao, Card_Tong:
		return (threeNInt & wnt_Min_3n_king_count_mask) >> wnt_Min_3n_king_count_Shift_count
	case Card_Feng:
		return (threeNInt & feng_Min_3n_king_count_mask) >> feng_Min_3n_king_count_Shift_count
	case Card_Arrow:
		return (threeNInt & arrow_Min_3n_king_count_mask) >> arrow_Min_3n_king_count_Shift_count
	}
	return 15
}
