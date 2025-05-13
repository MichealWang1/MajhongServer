package lib

import (
	"fmt"
)

// Scores 玩家分数
type Scores struct {
	table map[int32]int32
}

// NewScores 初始化积分
func NewScores() *Scores {
	t := make(map[int32]int32)
	return &Scores{table: t}
}

// NewScoresByMap 从map生成积分对象
func NewScoresByMap(m map[int32]int32, seatNums int32) (*Scores, error) {
	// 判断长度
	if len(m) != int(seatNums) {
		return nil, fmt.Errorf("newScoresByMap: map length is not equal to seatNums")
	}

	// 检查索引是否正确
	for n := int32(0); n < seatNums; n++ {
		if _, ok := m[n]; !ok {
			return nil, fmt.Errorf("newScoresByMap: index:%v not in map", n)
		}
	}

	obj := &Scores{table: m}
	return obj, nil
}

// Get 获取座位积分
func (this *Scores) Get(seatId int32) int32 {
	if n, ok := this.table[seatId]; ok {
		return n
	}
	return 0
}

// WinAll 赢所有人
func (this *Scores) WinAll(seatNums int32, seatId int32, score int32) {
	winScore := score * (seatNums - 1)
	this.table[seatId] += winScore

	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		if _, ok := this.table[tmpSeatId]; !ok {
			this.table[tmpSeatId] = 0
		}

		if tmpSeatId != seatId {
			this.table[tmpSeatId] -= score
		}
	}
}

// WinOne 赢一个人
func (this *Scores) WinOne(seatNums int32, winSeatId int32, loseSeatId int32, score int32) {
	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		if _, ok := this.table[tmpSeatId]; !ok {
			this.table[tmpSeatId] = 0
		}
	}

	this.table[winSeatId] += score
	this.table[loseSeatId] -= score
}

// Multiple 积分翻倍 - n为翻倍倍数，n应该大于1以上。(1被为自身，结果将直接在原有对象中保存)
func (this *Scores) Multiple(seatNums int32, n int32) {
	// 如果倍数小于等于1，直接return
	if n <= 1 {
		return
	}

	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		if _, ok := this.table[tmpSeatId]; ok {
			oldScore := this.table[tmpSeatId]
			this.table[tmpSeatId] = oldScore * n
			// this.table[tmpSeatId] += (this.table[tmpSeatId] * n)
		}
	}
}

// ToInt32Slice 转换为int32列表
func (this *Scores) ToInt32Slice(seatNums int32) []int32 {
	scores := make([]int32, 0)
	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		tmpScore, _ := this.table[tmpSeatId]
		scores = append(scores, tmpScore)
	}

	return scores
}

// ToMap 返回积分Map
func (this *Scores) ToMap() map[int32]int32 {
	return this.table
}

// String 实现Stringer接口
func (this *Scores) String() string {
	return fmt.Sprintf("Scores{table:%v}", this.table)
}

// InplaceAdd scores结构体相加，直接在现有的结构体上增加
func (this *Scores) InplaceAdd(seatNums int32, obj *Scores) {
	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		this.table[tmpSeatId] += obj.Get(tmpSeatId)
	}
}

// Add scores结构体相加, 返回一个新的Scores
func (this *Scores) Add(seatNums int32, obj *Scores) *Scores {
	retScores := NewScores()
	retScores.Reset(seatNums)

	otherScoresTable := obj.ToMap()

	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		score, _ := this.table[tmpSeatId]
		addScore, _ := otherScoresTable[tmpSeatId]
		retScores.table[tmpSeatId] += (score + addScore)
	}
	return retScores
}

// Adds 组合多个积分
func (this *Scores) Adds(seatNums int32, scores ...*Scores) *Scores {
	retScores := NewScores()
	retScores.Reset(seatNums)

	for tmpSeatId := int32(0); tmpSeatId < seatNums; tmpSeatId++ {
		score, _ := this.table[tmpSeatId]
		retScores.table[tmpSeatId] += score

		for _, tmpObj := range scores {
			retScores.table[tmpSeatId] += tmpObj.Get(tmpSeatId)
		}
	}
	return retScores
}

// Reset 重置
func (this *Scores) Reset(seatNums int32) {
	t := make(map[int32]int32)
	for n := int32(0); n < seatNums; n++ {
		t[n] = 0
	}
	this.table = t
}

// Transfer 转移分数
func (this *Scores) Transfer(srcSeatId int32, dstSeatId int32) {
	srcScore, _ := this.table[srcSeatId]
	this.table[dstSeatId] += srcScore
	this.table[srcSeatId] = 0
}
