package logic

import (
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
	"math/rand"
	"time"
)

type SwapStateData struct {
	SwapCards map[int32]lib.Cards // 玩家选择换三张的牌
	SwapType  pb.SwapType         // 换牌类型
}

func NewSwapStateData() *SwapStateData {
	return &SwapStateData{SwapCards: make(map[int32]lib.Cards), SwapType: pb.SwapType_SWAP_TYPE_NEXT}
}

func (s *SwapStateData) Reset() {
	s.SwapCards = make(map[int32]lib.Cards)
}

// 设置玩家换的牌
func (s *SwapStateData) SetUserSwapCards(seatId int32, cards lib.Cards) {
	s.SwapCards[seatId] = cards
}

// 获取玩家换的牌
func (s *SwapStateData) GetUserSwapCards(seatId int32) lib.Cards {
	if _, ok := s.SwapCards[seatId]; !ok {
		return lib.Cards{}
	}
	return s.SwapCards[seatId]
}

// 获取玩家换牌状态
func (s *SwapStateData) GetUserSwapState(seatId int32) bool {
	if _, ok := s.SwapCards[seatId]; !ok {
		return true
	}
	return false
}

// 获取各个玩家换牌状态
func (s *SwapStateData) GetSwapStateToSlice() []bool {
	res := make([]bool, 4)
	for i, _ := range s.SwapCards {
		res[i] = true
	}
	return res
}

// 获取换牌类型
func (s *SwapStateData) GetSwapType() pb.SwapType {
	return s.SwapType
}

// 随机换牌类型
func (s *SwapStateData) RandomSwapType() {
	rand.Seed(time.Now().UnixMilli())
	r := rand.Int() % 3
	s.SwapType = pb.SwapType(r)
}
