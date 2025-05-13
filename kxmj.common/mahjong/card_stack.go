package lib

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// CardStack 牌堆
type CardStack struct {
	// 所有牌
	AllCards Cards
	// 洗完牌后的牌堆
	CurrentCards Cards
	// 现在拿到牌堆的第几张
	index int
	// 最后位置
	lastIndex int
}

// NewCardStack 创建牌堆
func NewCardStack(cards Cards) *CardStack {
	cardStack := &CardStack{
		AllCards:     cards.Copy(),
		CurrentCards: cards.Copy(),
		index:        0,
		lastIndex:    len(cards),
	}
	cardStack.Reset()
	return cardStack
}

// Reset 重置牌堆
func (s *CardStack) Reset() {
	s.shuffle()
	s.index = 0
	s.lastIndex = len(s.AllCards)
}

// shuffle 洗牌
func (s *CardStack) shuffle() {
	s.CurrentCards = s.AllCards.Copy()
	// 生成随机种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 洗牌
	for i := len(s.CurrentCards) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		s.CurrentCards[i], s.CurrentCards[j] = s.CurrentCards[j], s.CurrentCards[i]
	}
}

func (s *CardStack) Set(cards Cards) {
	s.CurrentCards = cards
}

// Put 从牌堆里取一张牌(按顺序)
func (s *CardStack) Put() (Card, error) {
	if s.index >= s.lastIndex {
		return INVALID_CARD, errors.New(fmt.Sprintf("CardStack is no Card"))
	}
	card := s.CurrentCards[s.index]
	s.index++
	return card, nil
}

// Get 从牌堆里取第n张牌(这里的n是指这个操作不会删牌)
func (s *CardStack) Get(n int) (Card, error) {
	if n >= s.lastIndex {
		return INVALID_CARD, errors.New(fmt.Sprintf("CardStack is no Card"))
	}
	return s.CurrentCards[n], nil
}

// PickUpCards 获取牌堆中n张牌
func (s *CardStack) PickUpCards(n int) (Cards, error) {
	if s.index+n > s.lastIndex {
		return nil, errors.New(fmt.Sprintf("CardStack is no Card"))
	}
	cards := s.CurrentCards[s.index : s.index+n]
	s.index += n
	return cards, nil
}

// GM工具 发自定义的手牌
func (s *CardStack) PickUpCardsByTable(cardTable Cards) []*Card {
	if len(cardTable) <= 0 {
		return nil
	}
	cards := make([]*Card, 0)
	for _, cardValue := range cardTable {
		card, err := s.GetFirstSameCard(cardValue)
		if err != nil {
		}
		cards = append(cards, &card)
	}
	return cards
}

// GetResidueCardsNum 获取剩余牌堆数量
func (s *CardStack) GetResidueCardsNum() int {
	return s.lastIndex - s.index
}

// GetResidueCards 获取剩余的牌堆
func (s *CardStack) GetResidueCards() Cards {
	return s.CurrentCards[s.index:s.lastIndex].Copy()
}

// GetFirstSameCard 获取顺序第一张的某张牌
func (s *CardStack) GetFirstSameCard(card Card) (Card, error) {
	for i, c := range s.CurrentCards {
		if c == card {
			s.CurrentCards[s.index], s.CurrentCards[i] = s.CurrentCards[i], s.CurrentCards[s.index]
			s.index++
			return c, nil
		}
	}
	return INVALID_CARD, errors.New("CardStack is no have this card")
}

// GetFinalSameCard 获取倒序第一张的某张牌
func (s *CardStack) GetFinalSameCard(card Card) (Card, error) {
	for i := s.GetResidueCardsNum() - 1; i >= 0; i-- {
		if s.CurrentCards[i] == card {
			s.CurrentCards[s.index], s.CurrentCards[i] = s.CurrentCards[i], s.CurrentCards[s.index]
			s.index++
			return card, nil
		}
	}
	return INVALID_CARD, errors.New("CardStack is no have this card")
}

// PutFinal 获取最后一张牌
func (s *CardStack) PutFinal() (Card, error) {
	if s.index >= s.lastIndex {
		return INVALID_CARD, errors.New(fmt.Sprintf("CardStack is no Card"))
	}
	card := s.CurrentCards[s.lastIndex-1]
	s.lastIndex--
	return card, nil
}
