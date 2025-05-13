package lib

// UserAction 玩家动作 如：玩家1打了张3万，玩家2碰了
// UserAction{SeatId:2,OutSeatId:1,ActionType:ActionType_Peng,OutCard:0x03,DeleteCards:Cards{0x03,0x03},CombineCards:Cards{0x03,0x03,0x03},IsDrop:false}
type UserAction struct {
	SeatId          int32           // 动作对象
	OutSeatId       int32           // 出牌对象
	ActionType      ActionType      // 动作类型
	ExtraActionType ExtraActionType // 其他动作(明杠、补杠、暗杠)
	OutCard         Card            // 动作牌
	DeleteCards     Cards           // 删除牌
	CombineCards    Cards           // 组合牌
	IsDrop          bool            // 是否被抢杠
}

func (u *UserAction) Copy() *UserAction {
	return &UserAction{
		SeatId:          u.SeatId,
		OutSeatId:       u.OutSeatId,
		ActionType:      u.ActionType,
		ExtraActionType: u.ExtraActionType,
		OutCard:         u.OutCard,
		DeleteCards:     u.DeleteCards,
		CombineCards:    u.CombineCards,
		IsDrop:          u.IsDrop,
	}
}

// SetUserActionRob 设置这个动作被抢杠
func (a *UserAction) SetUserActionRob() {
	a.IsDrop = true
}
