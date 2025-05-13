package lib

// IsPingHu 检测平胡（3n+2）
func IsPingHu(handCards Cards, outCard Card, ziSequence bool) bool {
	cards := handCards.AddCard(outCard)
	// 提取将牌
	cardMap := cards.ToMap()
	combinationCards := make([]Cards, 0)
	for card, count := range cardMap {
		if count >= 2 {
			combinationCards = append(combinationCards, card.Repeat(2))
		}
	}
	// 循环删除组合牌
	for i := 0; i < len(combinationCards); i++ {
		// 删除一组牌
		newCards, _ := cards.DeleteCards(combinationCards[i])
		if Is3n(newCards, ziSequence) {
			return true
		}
	}
	return false
}

// Is3n 检测3n(ziSequence:字是否成顺)
func Is3n(cards Cards, ziSequence bool) bool {
	if cards.Len() == 0 {
		return true
	}
	cardMap := cards.ToMap()
	combinationCards := make([]Cards, 0)
	// 取刻子
	for card, count := range cardMap {
		if count >= 3 {
			combinationCards = append(combinationCards, card.Repeat(3))
		}
	}
	// 去重
	uniqueCards := cards.ToUnique()
	// 取顺子 (取当前牌加后两张 x--)
	for i := 0; i < uniqueCards.Len(); i++ {
		// 万条筒的顺子
		if uniqueCards[i].GetColor() < Card_Color_Zi {
			if uniqueCards[i].GetValue() > byte(7) {
				continue
			}
			nextCard := uniqueCards[i].NextCard()
			nextNextCard := nextCard.NextCard()
			if uniqueCards.IsContain(Cards{nextCard, nextNextCard}) {
				combinationCards = append(combinationCards, Cards{uniqueCards[i], nextCard, nextNextCard})
			}
		} else if ziSequence && uniqueCards[i].GetColor().IsZi() { // 字成顺取顺子
			if uniqueCards[i].GetValue() == byte(1) { // 东（东南西、东南北、东西北）
				fengCards := Cards{0x32, 0x33, 0x34}
				for i := 0; i < fengCards.Len(); i++ {
					for j := i + 1; j < fengCards.Len(); j++ {
						if uniqueCards.IsContain(Cards{fengCards[i], fengCards[j]}) {
							combinationCards = append(combinationCards, Cards{uniqueCards[i], fengCards[i], fengCards[j]})
						}
					}
				}
			} else if uniqueCards[i].GetValue() == byte(2) { // 西（西南北）
				if uniqueCards.IsContain(Cards{0x33, 0x34}) {
					combinationCards = append(combinationCards, Cards{uniqueCards[i], 0x33, 0x34})
				}
			} else if uniqueCards[i].GetValue() == byte(5) { // 中（中发白）
				if uniqueCards.IsContain(Cards{0x36, 0x37}) {
					combinationCards = append(combinationCards, Cards{uniqueCards[i], 0x36, 0x37})
				}
			}
		}
	}
	// 循环删除组合牌
	for i := 0; i < len(combinationCards); i++ {
		// 删除一组牌
		newCards, _ := cards.DeleteCards(combinationCards[i])
		if Is3n(newCards, ziSequence) {
			return true
		}
	}
	return false
}

// 组合
type combinations struct {
	Cards     Cards
	kingCount int // 需要精牌数量
}

// IsKingHu 带精平胡
func IsKingHu(handCards, kingCards Cards, huCard Card, isZiMo, ziSequence bool) bool {
	// 获取手牌中精牌数量
	kingCount := 0
	for _, card := range kingCards {
		kingCount += handCards.GetCount(card)
	}
	// 把赖子牌移除
	removeKingCards := handCards.RemoveCards(kingCards)
	// 如果huCard是赖子牌，得判断是否是摸到的还是别人打的（自己摸得才能当赖子牌，别人打的只能当打出的牌本身）
	if kingCards.In(huCard) && isZiMo {
		kingCount++
	} else {
		removeKingCards = removeKingCards.AddCard(huCard)
	}
	// 如果手上全是精牌可以胡
	if removeKingCards.Len() == 0 {
		return true
	}

	cardMap := removeKingCards.ToMap()
	// 组合将牌
	combinationCards := make([]combinations, 0)
	for card, count := range cardMap {
		if count >= 2 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(2)})
		} else if int(count)+kingCount >= 2 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(int(count)), kingCount: 2 - int(count)})
		}
	}

	// 根据组合进行删牌
	for _, combination := range combinationCards {
		newCards := removeKingCards.RemoveCards(combination.Cards)
		if IsKing3n(newCards, kingCount-combination.kingCount, ziSequence) {
			return true
		}
	}
	return false
}

// IsKing3n 带精3n
func IsKing3n(cards Cards, kingCount int, ziSequence bool) bool {
	if cards.Len() == 0 {
		return true
	}
	cardMap := cards.ToMap()
	combinationCards := make([]combinations, 0)
	// 取刻子
	for card, count := range cardMap {
		if count >= 3 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(3)})
		} else if int(count)+kingCount >= 3 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(int(count)), kingCount: 3 - int(count)})
		}
	}
	// 去重
	uniqueCards := cards.ToUnique()
	// 取顺子 (取当前牌加后两张 x--)
	for i := 0; i < uniqueCards.Len(); i++ {
		// 万条筒的顺子
		if uniqueCards[i].GetColor() < Card_Color_Zi {
			if uniqueCards[i].GetValue() > byte(8) {
				continue
			} else if uniqueCards[i].GetValue() == byte(8) {
				nextCard := uniqueCards[i].NextCard()
				if kingCount >= 1 && uniqueCards.IsContain(Cards{nextCard}) {
					combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], nextCard}, kingCount: 1})
				}
			}
			nextCard := uniqueCards[i].NextCard()
			nextNextCard := nextCard.NextCard()
			if uniqueCards.IsContain(Cards{nextCard, nextNextCard}) {
				combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], nextCard, nextNextCard}})
			} else if kingCount >= 1 {
				if uniqueCards.In(nextCard) {
					combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], nextCard}, kingCount: 1})
				}
				if uniqueCards.In(nextNextCard) {
					combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], nextNextCard}, kingCount: 1})
				}
			}

		} else if ziSequence && uniqueCards[i].GetColor().IsZi() { // 字成顺取顺子
			if uniqueCards[i].GetValue() <= byte(4) { // 东（东南西、东南北、东西北）
				fengCards := Cards{0x31, 0x32, 0x33, 0x34}
				fengCards, _ = fengCards.DeleteCard(uniqueCards[i])
				for k := 0; k < fengCards.Len(); k++ {
					for j := k + 1; j < fengCards.Len(); j++ {
						if uniqueCards.IsContain(Cards{fengCards[k], fengCards[j]}) {
							combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], fengCards[k], fengCards[j]}})
						}
					}
					if kingCount >= 1 {
						if uniqueCards.In(fengCards[k]) {
							combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], fengCards[k]}, kingCount: 1})
						}
					}
				}
			} else { // 中发白
				ziCards := Cards{0x35, 0x36, 0x37}
				ziCards, _ = ziCards.DeleteCard(uniqueCards[i])
				for k := 0; k < ziCards.Len(); k++ {
					for j := k + 1; j < ziCards.Len(); j++ {
						if uniqueCards.IsContain(Cards{ziCards[k], ziCards[j]}) {
							combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], ziCards[k], ziCards[j]}})
						}
					}
					if kingCount >= 1 {
						if uniqueCards.In(ziCards[k]) {
							combinationCards = append(combinationCards, combinations{Cards: Cards{uniqueCards[i], ziCards[k]}, kingCount: 1})
						}
					}
				}
			}
		}
	}

	// 循环删除组合牌
	for i := 0; i < len(combinationCards); i++ {
		// 删除一组牌
		newCards, _ := cards.DeleteCards(combinationCards[i].Cards)
		if IsKing3n(newCards, kingCount-combinationCards[i].kingCount, ziSequence) {
			return true
		}
	}
	return false
}

// Is7Pair 七对
func Is7Pair(handCards Cards, huCard Card) bool {
	// 七对不能有动作牌
	if handCards.Len() != 13 {
		return false
	}
	cards := handCards.AddCard(huCard)
	cardMap := cards.ToMap()
	for _, count := range cardMap {
		if count%2 != 0 {
			return false
		}
	}
	return true
}

// IsKing7Pair 带精七对
func IsKing7Pair(handCards, kingCards Cards, huCard Card, isZiMo bool) bool {
	// 七对不能有动作牌
	if handCards.Len() != 13 {
		return false
	}
	kingCount := 0
	for _, card := range kingCards {
		kingCount += handCards.GetCount(card)
	}
	removeKingCards := handCards.RemoveCards(kingCards)
	if isZiMo && kingCards.In(huCard) {
		kingCount++
	} else {
		removeKingCards = removeKingCards.AddCard(huCard)
	}
	cardMap := removeKingCards.ToMap()
	for _, count := range cardMap {
		if count%2 != 0 {
			if kingCount <= 0 {
				return false
			}
			kingCount--
		}
	}
	return true
}
