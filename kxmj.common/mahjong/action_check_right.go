package lib

// IsPengPengHu 碰碰胡（手中牌和动作牌全是刻字、没有顺子）
func IsPengPengHu(handCards Cards, huCard Card, actions UserActions) bool {
	// 判断动作是否都是刻子（碰、杠）
	for _, action := range actions {
		if action.ActionType == ActionType_Hu {
			continue
		}
		if action.ActionType != ActionType_Peng && action.ActionType != ActionType_Gang {
			return false
		}
	}
	cards := handCards.AddCard(huCard)
	combinationCards := make([]Cards, 0, 4)
	cardMap := cards.ToMap()
	// 取将
	for card, count := range cardMap {
		if count >= 2 {
			combinationCards = append(combinationCards, card.Repeat(2))
		}
	}
	// 删除将牌后的牌都是刻子
	for _, combination := range combinationCards {
		newCards, _ := cards.DeleteCards(combination)
		newCardMap := newCards.ToMap()
		flag := true
		for _, count := range newCardMap {
			if count != 3 {
				flag = false
				break
			}
		}
		if flag {
			return true
		}
	}
	return false
}

// IsKingPengPengHu 带精碰碰胡
func IsKingPengPengHu(handCards, kingCards Cards, huCard Card, actions UserActions, isZiMo bool) bool {
	// 判断动作是否都是刻子（碰、杠）
	for _, action := range actions {
		if action.ActionType == ActionType_Hu {
			continue
		}
		if action.ActionType != ActionType_Peng && action.ActionType != ActionType_Gang {
			return false
		}
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
	combinationCards := make([]combinations, 0, 4)
	// 取将
	for card, count := range cardMap {
		if count >= 2 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(2)})
		} else if kingCount+int(count) >= 2 {
			combinationCards = append(combinationCards, combinations{Cards: card.Repeat(int(count)), kingCount: 2 - int(count)})
		}
	}

	// 删除将牌后的牌都是刻子
	for _, combination := range combinationCards {
		newCards, _ := removeKingCards.DeleteCards(combination.Cards)
		cpKingCount := kingCount - combination.kingCount
		flag := true
		for _, count := range newCards.ToMap() {
			if count != 3 {
				if int(count)%3+cpKingCount <= 3 {
					flag = false
					break
				}
				cpKingCount -= 3 - int(count)%3
			}
		}
		if flag {
			return true
		}
	}
	return false
}

// 一条龙（手中有1-9的同种所有牌）
func IsYiTiaoLong(handCards Cards, huCard Card) bool {
	cards := handCards.AddCard(huCard)
	// 去重后长度小于9
	uniqueCards := cards.ToUnique()
	if uniqueCards.Len() < 9 {
		return false
	}
	for _, color := range []CardColor{Card_Color_Wan, Card_Color_Tiao, Card_Color_Tong} {
		if uniqueCards.GetAllColorCards(color).Len() == 9 {
			return true
		}
	}
	return false
}

// 带精一条龙
func IsKingYiTiaoLong(handCards, kingCards Cards, huCard Card, isZiMo bool) bool {
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
	// 去重
	uniqueCards := removeKingCards.ToUnique()
	if uniqueCards.Len() < 9 {
		return false
	}
	for _, color := range []CardColor{Card_Color_Wan, Card_Color_Tiao, Card_Color_Tong} {
		if uniqueCards.GetAllColorCards(color).Len()+kingCount >= 9 {
			return true
		}
	}
	return false
}

// IsQingYiSe 清一色（全是一种花色）
func IsQingYiSe(handCards Cards, huCards Card, actions UserActions) bool {
	cards := handCards.AddCard(huCards)
	color := cards[0].GetColor()
	for _, action := range actions {
		if action.OutCard.GetColor() != color {
			return false
		}
	}
	return cards.IsSameColor()
}

// IsKingQingYiSe 带精清一色
func IsKingQingYiSe(handCards, kingCards Cards, huCard Card, isZiMo bool, actions UserActions) bool {
	removeKingCards := handCards.RemoveCards(kingCards)
	if !isZiMo || !kingCards.In(huCard) {
		removeKingCards = removeKingCards.AddCard(huCard)
	}
	if removeKingCards.Len() == 0 && actions.Len() > 0 {
		color := actions[0].OutCard.GetColor()
		for _, action := range actions {
			if action.OutCard.GetColor() != color {
				return false
			}
		}
	} else {
		color := removeKingCards[0].GetColor()
		for _, action := range actions {
			if action.OutCard.GetColor() != color {
				return false
			}
		}
	}
	return removeKingCards.IsSameColor()
}

// 门清（没有动作牌）
func IsMenQing(handCards Cards, huCard Card) bool {
	cards := handCards.AddCard(huCard)
	return cards.Len() == 14
}

// 暗杠门清（除暗杠外没有其他动作牌）
func IsAnGangMenQing(handCards Cards, huCard Card, actions UserActions) bool {
	for _, action := range actions {
		if action.ActionType != ActionType_Gang {
			return false
		} else if action.ExtraActionType != ExtraActionType_An_Gang {
			return false
		}
	}
	return true
}

// IsHunYaoJiu 是否组合都带幺九(刻子、顺子、将)
func IsHunYaoJiu(handCards Cards, huCard Card, actions UserActions) bool {
	for _, action := range actions {
		haveYaoJiu := false
		for _, card := range action.CombineCards {
			if card.GetColor() != Card_Color_Zi && (card.GetValue() == 1 || card.GetValue() == 9) {
				haveYaoJiu = true
				break
			}
		}
		if !haveYaoJiu {
			return false
		}
	}
	cards := handCards.AddCard(huCard)
	combCards := make([]Cards, 0, 5)
	cardMap := cards.ToMap()
	// 取将
	for card, count := range cardMap {
		if count >= 2 && card.GetColor() != Card_Color_Zi && (card.GetValue() == 1 || card.GetValue() == 9) {
			combCards = append(combCards, card.Repeat(2))
		}
	}
	for _, comb := range combCards {
		copyCards, _ := cards.DeleteCards(comb)
		if getYaoJiu3n(copyCards) {
			return true
		}
	}
	return false
}

// 取幺九顺刻
func getYaoJiu3n(cards Cards) bool {
	if cards.Len() == 0 {
		return true
	}
	copyCards := cards.Copy()
	combCards := make([]Cards, 0, 10)
	for _, card := range cards.ToUnique() {
		if card.GetColor() != Card_Color_Zi {
			if card.GetValue() == 1 {
				combCards = append(combCards, card.Repeat(3))                                                        // 刻子
				combCards = append(combCards, Cards{card, NewCard(card.GetColor(), 2), NewCard(card.GetColor(), 3)}) // 顺子
			} else if card.GetValue() == 9 {
				combCards = append(combCards, card.Repeat(3))                                                        // 刻子
				combCards = append(combCards, Cards{NewCard(card.GetColor(), 7), NewCard(card.GetColor(), 8), card}) // 顺子
			}
		}
	}
	for _, comb := range combCards {
		newCards, err := copyCards.DeleteCards(comb)
		if err != nil {
			continue
		}
		if getYaoJiu3n(newCards) {
			return true
		}
	}
	return false
}
