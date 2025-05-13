package logic

import (
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
)

type PlayStateData struct {
	UserState        map[int32]pb.PlayerState // 玩家状态
	HandCardsTable   *lib.CardsTable          // 玩家手牌表
	DiscardTable     *lib.CardsTable          // 玩家弃牌表
	ActionCardsTable *lib.ActionCardsTable    // 玩家动作表
	CatchCardTable   *lib.CatchCardTable      // 玩家抓牌表

	// 玩家行为判断
	actionInspectTable *lib.ActionInspect // 动作行为检测
	huTable            *lib.HuTable       // 胡牌行为检测

}

func NewPlayStateData() *PlayStateData {
	p := &PlayStateData{}
	p.UserState = make(map[int32]pb.PlayerState, 4)
	p.HandCardsTable = lib.NewCardsTable()
	p.DiscardTable = lib.NewCardsTable()
	p.ActionCardsTable = lib.NewActionCardsTable()
	p.CatchCardTable = lib.NewCatchCardTable()

	p.actionInspectTable = lib.NewActionInspect()
	p.huTable = lib.NewHuTable()

	return p
}

func (p *PlayStateData) Reset() {
	p.UserState = make(map[int32]pb.PlayerState, 4)
	p.HandCardsTable.Reset()
	p.DiscardTable.Reset()
	p.ActionCardsTable.Reset()
	p.CatchCardTable.Reset()

	p.actionInspectTable.Reset()
	p.huTable.Reset()
}

// SetUserState 设置玩家状态
func (p *PlayStateData) SetUserState(seatId int32, state pb.PlayerState) {
	p.UserState[seatId] = state
}

// GetUserState 获取玩家状态
func (p *PlayStateData) GetUserState(seatId int32) pb.PlayerState {
	return p.UserState[seatId]
}

// GetAllUserStateToSlice 获取所有玩家状态到slice里
func (p *PlayStateData) GetAllUserStateToSlice() []pb.PlayerState {
	res := make([]pb.PlayerState, len(p.UserState))
	for i := 0; i < len(res); i++ {
		res[i] = p.UserState[int32(i)]
	}
	return res
}

// 获取玩家手牌表
func (p *PlayStateData) GetUserHandCardsTable() *lib.CardsTable {
	return p.HandCardsTable
}

// 添加玩家手牌
func (p *PlayStateData) AddUserHandCards(seatId int32, card lib.Card) {
	p.HandCardsTable.Add(seatId, card)
}

// 设置玩家手牌信息
func (p *PlayStateData) SetUserHandCards(seatId int32, cards lib.Cards) {
	p.HandCardsTable.Set(seatId, cards)
}

// 获取玩家手牌
func (p *PlayStateData) GetUserHandCards(seatId int32) lib.Cards {
	return p.HandCardsTable.Get(seatId)
}

// 增加玩家弃牌
func (p *PlayStateData) AddUserDiscard(seatId int32, card lib.Card) {
	p.DiscardTable.Add(seatId, card)
}

// 获取玩家的弃牌
func (p *PlayStateData) GetUserDiscards(seatId int32) lib.Cards {
	return p.DiscardTable.Get(seatId)
}

// 获取所有玩家弃牌表
func (p *PlayStateData) GetDiscardTable() *lib.CardsTable {
	return p.DiscardTable
}

// 添加玩家动作
func (p *PlayStateData) AddUserActionCards(seatId int32, action *lib.UserAction) {
	p.ActionCardsTable.Add(seatId, action)
}

// 获取玩家动作
func (p *PlayStateData) GetUserActionCards(seatId int32) lib.UserActions {
	return p.ActionCardsTable.Get(seatId)
}

// 获取玩家动作表
func (p *PlayStateData) GetActionCardsTable() *lib.ActionCardsTable {
	return p.ActionCardsTable
}

// 设置玩家抓牌
func (p *PlayStateData) SetUserCatchCard(seatId int32, card lib.Card) {
	p.CatchCardTable.Set(seatId, card)
}

// 获取玩家抓牌表
func (p *PlayStateData) GetCatchCardTable() *lib.CatchCardTable {
	return p.CatchCardTable
}
