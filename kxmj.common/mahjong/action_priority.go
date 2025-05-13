package lib

import (
	"errors"
	"fmt"
)

// ActionInspect 检查权限最大的动作
type ActionInspect struct {
	OutSeatId         int32        // 出牌玩家位置
	SeatOrder         []int32      // 位子顺序
	ActionOder        []ActionType // 动作大小
	UntreatedActions  UserActions  // 未处理的动作
	PossibleMaxAction *UserAction  // 可能的最大动作
	ResponseMaxAction *UserAction  // 响应的最大动作
	CanAllHu          bool         // 是否能胡多个(是否能一炮多响)
}

func (a *ActionInspect) SetSeatOrder(s []int32)       { a.SeatOrder = s }
func (a *ActionInspect) SetActionOder(t []ActionType) { a.ActionOder = t }
func (a *ActionInspect) SetCanAllHu(can bool)         { a.CanAllHu = can }

func NewActionInspect() *ActionInspect {
	return &ActionInspect{
		OutSeatId:         0,
		SeatOrder:         []int32{0, 1, 2, 3},
		ActionOder:        []ActionType{ActionType_Hu, ActionType_Gang, ActionType_Peng, ActionType_Chi, ActionType_Pass},
		UntreatedActions:  make(UserActions, 0, 4),
		PossibleMaxAction: nil,
		ResponseMaxAction: nil,
		CanAllHu:          false,
	}
}

func (a *ActionInspect) GetResponseMaxAction() *UserAction {
	return a.ResponseMaxAction
}

func (a *ActionInspect) GetPossibleMaxAction() *UserAction {
	return a.PossibleMaxAction
}

// HasActions 是否有动作
func (a *ActionInspect) HasActions() bool {
	return a.UntreatedActions.Len() > 0
}

// HasHu 是否有胡操作未操作
func (a *ActionInspect) HasHu() bool {
	for _, action := range a.UntreatedActions {
		if action.ActionType == ActionType_Hu {
			return true
		}
	}
	return false
}

// GetUserHasHu 玩家是否有胡操作未操作
func (a *ActionInspect) GetUserHasHu(seatId int32) (*UserAction, bool) {
	for _, action := range a.UntreatedActions {
		if action.SeatId == seatId && action.ActionType == ActionType_Hu {
			return action, true
		}
	}
	return nil, false
}

// AddActions 将要处理的动作加入到未处理动作
func (a *ActionInspect) AddActions(outSeatId int32, actions UserActions) {
	if a.UntreatedActions == nil {
		a.UntreatedActions = make(UserActions, 0, 4)
	}
	a.OutSeatId = outSeatId
	a.UntreatedActions = append(a.UntreatedActions, actions...)
}

// 玩家处理动作后将未处理的这个动作（包括这个玩家的动作去除）
// 玩家处理可能有：1、PASS；2、处理动作

func (a *ActionInspect) Pass(seatId int32) (*UserAction, bool) {
	action := &UserAction{
		SeatId:          seatId,
		OutSeatId:       a.OutSeatId,
		ActionType:      ActionType_Pass,
		ExtraActionType: ExtraActionType_Null,
		OutCard:         Card_Unknown,
		DeleteCards:     Cards{},
		CombineCards:    Cards{},
	}
	return a.Operate(action)
}

// Operate 处理动作(返回当前最大响应动作、是否可以结束当前操作)
func (a *ActionInspect) Operate(action *UserAction) (*UserAction, bool) {
	// 删除该玩家的未处理动作
	res := make(UserActions, 0)
	for _, untreatedAction := range a.UntreatedActions {
		if untreatedAction.SeatId != action.SeatId {
			res = append(res, untreatedAction)
		}
	}
	a.UntreatedActions = res
	fmt.Printf("UntreatedActions:%#v\n", a.UntreatedActions)
	// 更新最大响应动作
	a.IsMaxResponseAction(action)
	fmt.Printf("MaxResponseAction:%#v\n", a.ResponseMaxAction)
	// 更新最大可能动作
	a.UpdatePossibleMaxAction()
	fmt.Printf("PossibleMaxAction:%#v\n", a.PossibleMaxAction)

	return a.GetResponseMaxAction(), a.CanEnd()
}

// IsMaxResponseAction 判断是否是最大响应动作，是的话更新这个最大响应动作
func (a *ActionInspect) IsMaxResponseAction(action *UserAction) {
	if a.ResponseMaxAction == nil {
		a.ResponseMaxAction = action
	}
	maxAction := a.ComparePriority(a.ResponseMaxAction, action)
	a.ResponseMaxAction = maxAction
}

// ComparePriority 判断两个动作的优先级
func (a *ActionInspect) ComparePriority(action1, action2 *UserAction) *UserAction {
	if action1.ActionType == action2.ActionType {
		order := a.CompareSeatOrder(action1.SeatId, action2.SeatId)
		if order == action1.SeatId {
			return action1
		}
		return action2
	}

	for _, v := range a.ActionOder {
		if v == action1.ActionType {
			return action1
		}
		if v == action2.ActionType {
			return action2
		}
	}
	return action1
}

// CompareSeatOrder 判断两个位子优先级
func (a *ActionInspect) CompareSeatOrder(newSeatId, OldSeatId int32) int32 {
	tmp := a.OutSeatId
	userNums := len(a.SeatOrder)
	for i := 1; i < userNums; i++ {
		tmp = (tmp + int32(i)) % int32(userNums)
		if tmp == newSeatId {
			return newSeatId
		}
		if tmp == OldSeatId {
			return OldSeatId
		}
	}
	return newSeatId
}

// UpdatePossibleMaxAction 更新可能的最大动作(未处理的和响应的最大的动作)
func (a *ActionInspect) UpdatePossibleMaxAction() error {
	if a.UntreatedActions.Len() == 0 {
		return errors.New("no action!! Untreated actions is nil")
	}
	a.PossibleMaxAction = a.UntreatedActions[0]
	for _, action := range a.UntreatedActions {
		a.PossibleMaxAction = a.ComparePriority(a.PossibleMaxAction, action)
	}
	a.PossibleMaxAction = a.ComparePriority(a.PossibleMaxAction, a.ResponseMaxAction)
	return nil
}

// CanEnd 判断是否可以结束这个检测
func (a *ActionInspect) CanEnd() bool {
	// 未检测的动作没有了
	if !a.HasActions() {
		return true
	}
	// 有一炮多响时，未完成的操作里面还有胡操作
	if a.CanAllHu && a.HasHu() {
		return false
	}
	// 响应最大动作是可能的最大动作
	if a.ResponseMaxAction.SeatId == a.PossibleMaxAction.SeatId && a.ResponseMaxAction.ActionType == a.PossibleMaxAction.ActionType {
		return true
	}
	fmt.Printf("ResponseMaxAction:%#v,PossibleMaxAction:%#v\n", a.ResponseMaxAction, a.PossibleMaxAction)
	return false
}

// Reset 重置这个动作检测
func (a *ActionInspect) Reset() {
	a.UntreatedActions = make(UserActions, 0, 4)
	a.ResponseMaxAction = nil
	a.PossibleMaxAction = nil
}
