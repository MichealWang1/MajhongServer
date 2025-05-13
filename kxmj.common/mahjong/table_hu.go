package lib

import (
	"fmt"
	"sync"
)

// HuTable 各座位胡牌结果记录
type HuTable struct {
	sync.RWMutex
	table map[int32]*HuResult
}

// NewHuTable 初始化胡牌记录表
func NewHuTable() *HuTable {
	t := &HuTable{table: make(map[int32]*HuResult)}
	return t
}

// Add 添加座位胡牌结果记录
func (this *HuTable) Add(seatId int32, result *HuResult) error {
	this.Lock()
	defer this.Unlock()

	if result == nil {
		return fmt.Errorf("HuTable: Add: seatId:%v result is nil", seatId)
	}

	if _, ok := this.table[seatId]; ok {
		return fmt.Errorf("HuTable: Add: seatId:%v already set", seatId)
	}

	this.table[seatId] = result
	return nil
}

// Reset 重置
func (this *HuTable) Reset() {
	this.Lock()
	defer this.Unlock()

	this.table = make(map[int32]*HuResult)
}

// RemoveExclude 删除
func (this *HuTable) RemoveExclude(excludeSeatId int32) {
	this.Lock()
	defer this.Unlock()

	removeSeatIds := make([]int32, 0)
	for tblSeatId := range this.table {
		if tblSeatId != excludeSeatId {
			removeSeatIds = append(removeSeatIds, tblSeatId)
		}
	}

	for _, tmpSeatId := range removeSeatIds {
		delete(this.table, tmpSeatId)
	}
}

// Get 获取座位胡牌记录
func (this *HuTable) Get(seatId int32) (*HuResult, error) {
	this.RLock()
	defer this.RUnlock()

	result, ok := this.table[seatId]
	if !ok {
		return nil, fmt.Errorf("HuTable:Get: seatId:%v no huResult", seatId)
	}
	return result, nil
}

// 有多少个座位可以胡？
func (this *HuTable) CalculateHuSeatNums() int {
	this.RLock()
	defer this.RUnlock()

	n := 0
	for _, result := range this.table {
		if result != nil {
			n += 1
		}
	}
	return n
}

// ToMap 转换为Map表
func (this *HuTable) ToMap() map[int32]*HuResult {
	this.RLock()
	defer this.RUnlock()

	return this.table
}

// ToKindSlice 转换为胡牌牌型列表
func (this *HuTable) ToKindSlice(seatNums int32) []int32 {
	this.RLock()
	defer this.RUnlock()

	iSlice := make([]int32, 0)
	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		tmpKind := int32(0)
		if huResult, err := this.Get(tmpSeatId); err == nil {
			tmpKind = int32(huResult.GetKind())
		}
		iSlice = append(iSlice, tmpKind)
	}
	return iSlice
}

// ToRightSlice 转换为胡牌牌型列表
func (this *HuTable) ToRightSlice(seatNums int32) []int32 {
	this.RLock()
	defer this.RUnlock()

	iSlice := make([]int32, 0)
	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		tmpRight := int32(0)
		if huResult, err := this.Get(tmpSeatId); err == nil {
			tmpRight = int32(huResult.GetRight())
		}
		iSlice = append(iSlice, tmpRight)
	}
	return iSlice
}

// HasSeatHu 座位是否胡牌
func (this *HuTable) HasSeatHu(seatId int32) bool {
	this.RLock()
	defer this.RUnlock()

	for tmpSeatId, huResult := range this.table {
		if tmpSeatId == seatId && huResult != nil {
			return true
		}
	}
	return false
}
