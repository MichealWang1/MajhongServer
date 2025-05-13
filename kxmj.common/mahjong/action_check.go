package lib

// CheckAction 检测玩家动作(吃、碰、杠)
type CheckAction struct {
	SeatId      int32       // 检测玩家Id
	OutSeatId   int32       // 打牌/摸牌玩家id
	HandCards   Cards       // 手牌
	OutCard     Card        // 出/摸的牌
	UserActions UserActions // 玩家的动作牌
	ZiStraight  bool        // 字是否成顺
}

// CheckActionChi 检测玩家吃
func (c *CheckAction) CheckActionChi() UserActions {
	res := make(UserActions, 0, 4)
	// 判断万条筒的吃
	if c.OutCard.GetColor() < Card_Color_Zi {
		// 左吃(左吃最多到7 x--)
		if c.OutCard.GetValue() <= byte(7) {
			nextCard := c.OutCard.NextCard()
			nextNextCard := nextCard.NextCard()
			if c.HandCards.IsContain(Cards{nextCard, nextNextCard}) {
				res = append(res, &UserAction{
					SeatId:       c.SeatId,
					OutSeatId:    c.OutSeatId,
					ActionType:   ActionType_Chi,
					OutCard:      c.OutCard,
					DeleteCards:  Cards{nextCard, nextNextCard},
					CombineCards: Cards{c.OutCard, nextCard, nextNextCard},
				})
			}
		}
		// 中吃(0,9不会参与中吃-x-)
		if c.OutCard.GetValue() != byte(0) && c.OutCard.GetValue() != byte(9) {
			prevCard := c.OutCard.PrevCard()
			nextCard := c.OutCard.NextCard()
			if c.HandCards.IsContain(Cards{prevCard, nextCard}) {
				res = append(res, &UserAction{
					SeatId:       c.SeatId,
					OutSeatId:    c.OutSeatId,
					ActionType:   ActionType_Chi,
					OutCard:      c.OutCard,
					DeleteCards:  Cards{prevCard, nextCard},
					CombineCards: Cards{prevCard, c.OutCard, nextCard},
				})
			}
		}
		// 右吃(右吃最多到3 --x)
		if c.OutCard.GetValue() >= byte(3) {
			prevCard := c.OutCard.PrevCard()
			prevPrevCard := prevCard.PrevCard()
			if c.HandCards.IsContain(Cards{prevPrevCard, prevCard}) {
				res = append(res, &UserAction{
					SeatId:       c.SeatId,
					OutSeatId:    c.OutSeatId,
					ActionType:   ActionType_Chi,
					OutCard:      c.OutCard,
					DeleteCards:  Cards{prevPrevCard, prevCard},
					CombineCards: Cards{prevPrevCard, prevCard, c.OutCard},
				})
			}
		}
	}
	// 字成顺
	if c.ZiStraight {
		if c.OutCard.GetColor().IsZi() {
			// 东西南北
			fengCards := Cards{0x31, 0x32, 0x33, 0x34}
			if fengCards.In(c.OutCard) {
				delCards, _ := fengCards.DeleteCard(c.OutCard) // 把打的那张牌删除
				inCards := c.HandCards.Intersection(delCards)  // 手牌里与除去打出那张牌的交集
				if inCards.Len() >= 2 {
					for i := 0; i < inCards.Len(); i++ {
						for j := i + 1; j < inCards.Len(); j++ {
							res = append(res, &UserAction{
								SeatId:       c.SeatId,
								OutSeatId:    c.OutSeatId,
								ActionType:   ActionType_Chi,
								OutCard:      c.OutCard,
								DeleteCards:  Cards{inCards[i], inCards[j]},
								CombineCards: Cards{inCards[i], inCards[j], c.OutCard},
							})
						}
					}
				}
			}
			// 中发白
			ziCards := Cards{0x35, 0x36, 0x37}
			if ziCards.In(c.OutCard) {
				delCards, _ := ziCards.DeleteCard(c.OutCard)  // 把打的那张牌删除
				inCards := c.HandCards.Intersection(delCards) // 手牌里与除去打出那张牌的交集
				if inCards.Len() == 2 {
					res = append(res, &UserAction{
						SeatId:       c.SeatId,
						OutSeatId:    c.OutSeatId,
						ActionType:   ActionType_Chi,
						OutCard:      c.OutCard,
						DeleteCards:  Cards{inCards[0], inCards[1]},
						CombineCards: Cards{inCards[0], inCards[1], c.OutCard},
					})
				}
			}
		}
	}
	return res
}

// CheckActionPeng 检测玩家碰
func (c *CheckAction) CheckActionPeng() UserActions {
	res := make(UserActions, 0, 4)
	if c.HandCards.GetCount(c.OutCard) >= 2 {
		res = append(res, &UserAction{
			SeatId:       c.SeatId,
			OutSeatId:    c.OutSeatId,
			ActionType:   ActionType_Peng,
			OutCard:      c.OutCard,
			DeleteCards:  c.OutCard.Repeat(2),
			CombineCards: c.OutCard.Repeat(3),
		})
	}
	return res
}

// CheckActionGang 检测玩家杠
func (c *CheckAction) CheckActionGang() UserActions {
	res := make(UserActions, 0, 4)
	if c.SeatId == c.OutSeatId {
		// 检测暗杠
		cardMap := c.HandCards.AddCard(c.OutCard).ToMap()
		for card, n := range cardMap {
			if n == 4 {
				res = append(res, &UserAction{
					SeatId:          c.SeatId,
					OutSeatId:       c.OutSeatId,
					ActionType:      ActionType_Gang,
					ExtraActionType: ExtraActionType_An_Gang,
					OutCard:         card,
					DeleteCards:     card.Repeat(4),
					CombineCards:    card.Repeat(4),
				})
			}
		}
		// 检测补杠
		for _, action := range c.UserActions {
			if action.ActionType == ActionType_Peng && cardMap[action.OutCard] > 0 {
				res = append(res, &UserAction{
					SeatId:          c.SeatId,
					OutSeatId:       c.OutSeatId,
					ActionType:      ActionType_Gang,
					ExtraActionType: ExtraActionType_Bu_Gang,
					OutCard:         action.OutCard,
					DeleteCards:     action.OutCard.Repeat(1),
					CombineCards:    action.OutCard.Repeat(4),
				})
			}
		}
	} else { // 明杠
		if c.HandCards.GetCount(c.OutCard) == 3 {
			res = append(res, &UserAction{
				SeatId:          c.SeatId,
				OutSeatId:       c.OutSeatId,
				ActionType:      ActionType_Gang,
				ExtraActionType: ExtraActionType_Ming_Gang,
				OutCard:         c.OutCard,
				DeleteCards:     c.OutCard.Repeat(3),
				CombineCards:    c.OutCard.Repeat(4),
			})
		}
	}
	return res
}
