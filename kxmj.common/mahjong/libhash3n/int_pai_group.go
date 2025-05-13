package libhash3n

type threeNGroup struct {
	fengSequence                 bool
	arrowSequence                bool
	color                        CardType
	cardCounter                  CardCounter
	reach_3n_min_need_king_count int //达到3n结构最小需要的精数量
	counterInt                   int
}

func newThreeNGroup(color CardType, cardCounter CardCounter, fengSequence bool, arrowSequence bool) *threeNGroup {
	threeNGroup := &threeNGroup{color: color, cardCounter: cardCounter, fengSequence: fengSequence, arrowSequence: arrowSequence}
	threeNGroup.init()
	return threeNGroup
}
func (threeNGroup *threeNGroup) init() {
	threeNGroup.counterInt = CardCounterToInt(threeNGroup.color, threeNGroup.cardCounter)
	threeNGroup.reach_3n_min_need_king_count = threeNGroup.get3nMinNeedKingCount()
}

func (intPaiGroup *threeNGroup) getMinNeedKingCount() int {
	return intPaiGroup.reach_3n_min_need_king_count
}

// 获取达到3n结构所需要的最少精的个数
func (intPaiGroup *threeNGroup) get3nMinNeedKingCount() int {
	loopCount := 0
	left, right := 0, 0
	if !intPaiGroup.fengSequence && !intPaiGroup.arrowSequence {
		right = 65393 - 1
	} else if !intPaiGroup.fengSequence && intPaiGroup.arrowSequence {
		right = 65393 - 1
	} else if intPaiGroup.fengSequence && !intPaiGroup.arrowSequence {
		right = 65413 - 1
	} else {
		right = 65413 - 1
	}

	for mid := (left + right) / 2; left <= right; mid = (left + right) / 2 {
		if loopCount > 1000 {
			return 8
		}
		ThreeInt := -1
		if !intPaiGroup.fengSequence && !intPaiGroup.arrowSequence {
			ThreeInt = Normal_ThreeNInts_Feng_false_Arrow_false[mid]
		} else if !intPaiGroup.fengSequence && intPaiGroup.arrowSequence {
			ThreeInt = Normal_ThreeNInts_Feng_false_Arrow_true[mid]
		} else if intPaiGroup.fengSequence && !intPaiGroup.arrowSequence {
			ThreeInt = Normal_ThreeNInts_Feng_true_Arrow_false[mid]
		} else {
			ThreeInt = Normal_ThreeNInts_Feng_true_Arrow_true[mid]
		}
		if ThreeInt < 0 {
			return 8
		}
		rawPatternInt := ThreeNIntToCounterInt(ThreeInt)
		if intPaiGroup.counterInt == rawPatternInt {
			return GetMin3nNeedKingCount(intPaiGroup.color, ThreeInt)
		}
		//因为数据都是降序的
		if intPaiGroup.counterInt < rawPatternInt {
			left = mid + 1
		} else {
			right = mid - 1
		}
		loopCount++
	}

	return 8
}
