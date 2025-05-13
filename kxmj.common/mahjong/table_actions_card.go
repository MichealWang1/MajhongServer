package lib

import (
	"fmt"
	"sync"
)

// 吃碰杠牌表数据管理

// ActionCardsTable 吃碰杠牌表数据管理
type ActionCardsTable struct {
	sync.RWMutex
	table map[int32]UserActions
}

// NewActionCardsTable 初始化ActionCardsTable
func NewActionCardsTable() *ActionCardsTable {
	return &ActionCardsTable{table: make(map[int32]UserActions)}
}

// Reset 重置
func (this *ActionCardsTable) Reset() {
	this.Lock()
	defer this.Unlock()

	this.table = make(map[int32]UserActions)
}

// Get 获取用户的吃碰杠操作
func (this *ActionCardsTable) Get(seatId int32) UserActions {
	this.RLock()
	defer this.RUnlock()

	actions, ok := this.table[seatId]
	if !ok {
		return UserActions{}
	}
	return actions
}

// Add 添加用户的吃碰杠操作
func (this *ActionCardsTable) Add(seatId int32, action *UserAction) (UserActions, error) {
	this.Lock()
	defer this.Unlock()

	if action == nil {
		return nil, fmt.Errorf("ActionCardsTable: Add: action is nil")
	}

	var newActions UserActions = make(UserActions, 0)
	if actions, ok := this.table[seatId]; ok {
		newActions = append(actions, action)
	} else {
		newActions = append(newActions, action)
	}

	this.table[seatId] = newActions
	return newActions, nil
}

// Set 设置用户吃碰杠数据
func (this *ActionCardsTable) Set(seatId int32, actions UserActions) {
	this.Lock()
	defer this.Unlock()

	this.table[seatId] = actions
}

// BuGangAction 补杠操作
func (this *ActionCardsTable) BuGangAction(seatId int32, opCard Card) error {
	this.Lock()
	defer this.Unlock()

	actions, ok := this.table[seatId]
	if !ok {
		return fmt.Errorf("ActionCardsTable: BuGangAction: seatId:%v opCard:%v not found", seatId, opCard)
	}

	for _, action := range actions {
		if action.ActionType == ActionType_Peng && action.OutCard == opCard {
			action.ActionType = ActionType_Gang
			action.ExtraActionType = ExtraActionType_Bu_Gang
			action.CombineCards = opCard.Repeat(4)
		}
	}
	return nil
}
