package lib

import (
	"fmt"
	"sync"
)

type CardsTable struct {
	sync.RWMutex
	table map[int32]Cards
}

// NewCardsTable 初始化麻将牌表
func NewCardsTable() *CardsTable {
	return &CardsTable{table: make(map[int32]Cards)}
}

// Reset 重置
func (this *CardsTable) Reset() {
	this.Lock()
	defer this.RUnlock()
	this.table = make(map[int32]Cards)
}

// Get 获取座位牌数据
func (this *CardsTable) Get(seatId int32) Cards {
	this.RLock()
	defer this.RUnlock()
	cards, ok := this.table[seatId]
	if !ok {
		return Cards{}
	}
	return cards
}

// Set 设置座位牌数据
func (this *CardsTable) Set(seatId int32, cards Cards) {
	this.Lock()
	defer this.Unlock()

	this.table[seatId] = cards
}

// Delete 删除牌 - 返回删除后剩余手牌
func (this *CardsTable) Delete(seatId int32, card Card) (Cards, error) {
	this.Lock()
	defer this.Unlock()

	cards, ok := this.table[seatId]
	if !ok {
		return Cards{}, fmt.Errorf("CardsTable: Delete: seatId:%v not found", seatId)
	}
	newCards, err := cards.DeleteCard(card)
	if err != nil {
		return Cards{}, err
	}
	this.table[seatId] = newCards
	return newCards, nil
}

// DeleteCards 删除多张牌 - 返回删除后剩余手牌
func (this *CardsTable) DeleteCards(seatId int32, cards Cards) (Cards, error) {
	this.Lock()
	defer this.Unlock()

	seatCards, ok := this.table[seatId]
	if !ok {
		return Cards{}, fmt.Errorf("CardsTable: Delete: seatId:%v not found", seatId)
	}
	newCards, err := seatCards.DeleteCards(cards)
	if err != nil {
		return Cards{}, err
	}
	this.table[seatId] = newCards
	return newCards, nil
}

// Add 抓牌 - 返回抓牌后手牌数据
func (this *CardsTable) Add(seatId int32, card Card) Cards {
	this.Lock()
	defer this.Unlock()

	cards, ok := this.table[seatId]
	if !ok {
		cards = Cards{}
	}
	newCards := cards.AddCard(card)
	this.table[seatId] = newCards
	return newCards
}

// ToSlice - 将表全部转换为列表
func (this *CardsTable) ToSlice(seatNums int32) [][]int32 {
	this.RLock()
	defer this.RUnlock()

	var tmpSeatId int32

	objs := make([][]int32, 0)
	for ; tmpSeatId < seatNums; tmpSeatId++ {
		if cards, ok := this.table[tmpSeatId]; ok {
			objs = append(objs, cards.ToInt32())
		} else {
			objs = append(objs, []int32{})
		}
	}

	return objs
}

// CountToSlice - 各座位牌数量
func (this *CardsTable) CountToSlice(seatNums int32) []int32 {
	this.RLock()
	defer this.RUnlock()

	var tmpSeatId int32

	objs := make([]int32, 0)
	for ; tmpSeatId < seatNums; tmpSeatId++ {
		if cards, ok := this.table[tmpSeatId]; ok {
			objs = append(objs, int32(len(cards)))
		} else {
			objs = append(objs, 0)
		}
	}

	return objs
}

// GetTable
func (this *CardsTable) GetTable() map[int32]Cards {
	return this.table
}
