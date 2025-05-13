package lib

import (
	"sync"
)

// CatchCardTable 抓牌管理表
type CatchCardTable struct {
	sync.RWMutex
	table map[int32]Card
}

// NewCatchCardTable 初始化抓牌表
func NewCatchCardTable() *CatchCardTable {
	return &CatchCardTable{table: make(map[int32]Card)}
}

// Set 设置座位抓牌
func (this *CatchCardTable) Set(seatId int32, card Card) {
	this.Lock()
	defer this.Unlock()

	this.table[seatId] = card
}

// Remove 移除座位抓牌
func (this *CatchCardTable) Remove(seatId int32) {
	this.Lock()
	defer this.Unlock()

	this.table[seatId] = Card_Unknown
}

// Get 获取座位抓牌数据
func (this *CatchCardTable) Get(seatId int32) Card {
	this.RLock()
	defer this.RUnlock()

	if card, ok := this.table[seatId]; ok {
		return card
	}
	return Card_Unknown
}

// Reset 重置表
func (this *CatchCardTable) Reset() {
	this.Lock()
	defer this.Unlock()

	this.table = make(map[int32]Card)
}
