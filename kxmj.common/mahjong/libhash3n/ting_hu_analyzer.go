package libhash3n

import (
	"fmt"
	lib "kxmj.common/mahjong"
)

var _ fmt.Formatter

type TingHuAnalyzer struct {
	fengSequence  bool      //风是否成顺
	arrowSequence bool      //前是否成顺
	argsIsValid   bool      //传入的参数是否正确
	cards         lib.Cards //手牌
	kingCards     lib.Cards //精牌
	kingCount     int       //精的数量

	notKingCardCounter CardCounter //去掉精后的counter

}

// cards:手牌，必须长度是3n+1，
// kingCards:精或宝,如果不支持宝的玩法，可传lib.INVALID_CARD
// fengSequence:风是否成顺
// arrowSequence:箭是否成顺
// 初始化的时候，将会对传入的参数做检测
func NewTingHuAnalyzer(cards lib.Cards, kingCards lib.Cards, fengSequence bool, arrowSequence bool) *TingHuAnalyzer {
	handCards := &TingHuAnalyzer{cards: cards, kingCards: kingCards, argsIsValid: false, fengSequence: fengSequence, arrowSequence: arrowSequence}
	if !handCards.checkCards() {
		return handCards
	}
	handCards.argsIsValid = true
	handCards.init()
	return handCards
}

// 对传入的card作检测
func (tingHuAnalyzer *TingHuAnalyzer) checkCards() bool {
	for _, card := range tingHuAnalyzer.kingCards {
		if !isValidKingCard(card) {
			return false
		}
	}
	for _, card := range tingHuAnalyzer.cards {
		if !isValidHandCard(card) {
			return false
		}
	}
	if len(tingHuAnalyzer.cards)%3 != 1 {
		return false
	}
	return true
}

func (tingHuAnalyzer *TingHuAnalyzer) init() {
	//计算出宝的数量
	for _, card := range tingHuAnalyzer.cards {
		for _, kingCard := range tingHuAnalyzer.kingCards {
			if kingCard == lib.INVALID_CARD {
				continue
			}
			if kingCard == card {
				tingHuAnalyzer.kingCount++
				break
			}
		}
	}

	//将宝位置上的数量抹0
	tingHuAnalyzer.notKingCardCounter = ToCardCounter(tingHuAnalyzer.cards)
	for _, kingCard := range tingHuAnalyzer.kingCards {
		if kingCard == lib.INVALID_CARD {
			continue
		}
		kingIndex := CardToIndex(kingCard)
		tingHuAnalyzer.notKingCardCounter[kingIndex] = 0
	}

}

func (tingHuAnalyzer *TingHuAnalyzer) getAllTou(opCard lib.Card) []*Tou {
	tous := make([]*Tou, 0)

	opCardIndex := CardToIndex(opCard)
	tingHuAnalyzer.notKingCardCounter[opCardIndex]++

	for index, count := range tingHuAnalyzer.notKingCardCounter {
		card := IndexToCard(index)
		if count >= 2 { //当前位置牌>=2 一对
			tous = append(tous, NewTou(card, card))
		}
		if tingHuAnalyzer.kingCount >= 1 && count >= 1 { //精和任一牌做头
			tous = append(tous, NewTou(card, lib.INVALID_CARD))
		}
	}

	if tingHuAnalyzer.kingCount >= 2 {
		tous = append(tous, NewTou(lib.INVALID_CARD, lib.INVALID_CARD))
	}
	tingHuAnalyzer.notKingCardCounter[opCardIndex]--
	return tous

}

func (tingHuAnalyzer *TingHuAnalyzer) isHu(opCard lib.Card) bool {
	tous := tingHuAnalyzer.getAllTou(opCard)
	if len(tous) <= 0 {
		return false
	}

	kingCount := tingHuAnalyzer.kingCount

	colors := []CardType{Card_Wan,
		Card_Tiao,
		Card_Tong,
		Card_Feng,
		Card_Arrow}

	opCardIndex := CardToIndex(opCard)
	for _, tou := range tous {
		tmpCardCounter := make(CardCounter, 34)
		copy(tmpCardCounter, tingHuAnalyzer.notKingCardCounter)

		//去掉头的牌
		tmpCardCounter[opCardIndex]++
		if tou.card1 != lib.INVALID_CARD {
			tou1Index := CardToIndex(tou.card1)
			tmpCardCounter[tou1Index]--
		}
		if tou.card2 != lib.INVALID_CARD {
			tou2Index := CardToIndex(tou.card2)
			tmpCardCounter[tou2Index]--
		}
		//cardCounter此时去掉了所有的精，以及去掉头了

		nowKingCount := kingCount - tou.GetKingCount() //此时精的数量
		is3n := true
		for _, color := range colors {
			needKingCount := 100
			switch color {
			case Card_Wan:
				someColorCardCounter := make(CardCounter, 9)
				copy(someColorCardCounter, tmpCardCounter[0:9])
				needKingCount = newThreeNGroup(color, someColorCardCounter, tingHuAnalyzer.fengSequence, tingHuAnalyzer.arrowSequence).getMinNeedKingCount()
			case Card_Tiao:
				someColorCardCounter := make(CardCounter, 9)
				copy(someColorCardCounter, tmpCardCounter[9:18])
				needKingCount = newThreeNGroup(color, someColorCardCounter, tingHuAnalyzer.fengSequence, tingHuAnalyzer.arrowSequence).getMinNeedKingCount()
			case Card_Tong:
				someColorCardCounter := make(CardCounter, 9)
				copy(someColorCardCounter, tmpCardCounter[18:27])
				needKingCount = newThreeNGroup(color, someColorCardCounter, tingHuAnalyzer.fengSequence, tingHuAnalyzer.arrowSequence).getMinNeedKingCount()
			case Card_Feng:
				someColorCardCounter := make(CardCounter, 4)
				copy(someColorCardCounter, tmpCardCounter[27:31])
				needKingCount = newThreeNGroup(color, someColorCardCounter, tingHuAnalyzer.fengSequence, tingHuAnalyzer.arrowSequence).getMinNeedKingCount()
			case Card_Arrow:
				someColorCardCounter := make(CardCounter, 3)
				copy(someColorCardCounter, tmpCardCounter[31:34])
				needKingCount = newThreeNGroup(color, someColorCardCounter, tingHuAnalyzer.fengSequence, tingHuAnalyzer.arrowSequence).getMinNeedKingCount()
			}

			nowKingCount -= needKingCount
			if nowKingCount < 0 { //精不够了，退出判断
				is3n = false
				break
			}
		}

		if is3n { //如果是3n结构，退出判断
			return true
		}
		//继续选取下一个头

	}
	return false
}

// 获取所有听的牌
func (tingHuAnalyzer *TingHuAnalyzer) GetAllTingCards() lib.Cards {
	tingCards := make(lib.Cards, 0)
	if !tingHuAnalyzer.argsIsValid {
		return tingCards
	}
	for _, card := range mahjongs {
		if tingHuAnalyzer.isHu(card) {
			tingCards = append(tingCards, card)
		}
	}
	return tingCards
}
