package lib

type UserActions []*UserAction

func (actions UserActions) Len() int { return len(actions) }

func (actions UserActions) Copy() UserActions {
	res := make(UserActions, len(actions))
	copy(res, actions)
	return res
}

// Add 添加
func (actions *UserActions) Add(action *UserAction) {
	*actions = append(*actions, action)
}

// Del 删除
func (actions *UserActions) Del(action *UserAction) {
	for i, v := range *actions {
		if v == action {
			*actions = append((*actions)[:i], (*actions)[i+1:]...)
			break
		}
	}
}

// Mod 修改
func (actions *UserActions) Mod(action *UserAction) {
	for i, v := range *actions {
		if v.ActionType == action.ActionType && v.OutCard == action.OutCard {
			*actions = append((*actions)[:i], append([]*UserAction{action}, (*actions)[i+1:]...)...)
			break
		}
	}
}

// HasAction 是否有这个动作
func (actions UserActions) HasAction(action *UserAction) bool {
	for _, untreatedAction := range actions {
		if untreatedAction.SeatId == action.SeatId && untreatedAction.CombineCards.IsContain(action.CombineCards) {
			return true
		}
	}
	return false
}

// 获取这个动作信息
//func (actions UserActions) GetAction(actionType ActionType,opCard Card,combinations Cards) UserAction {
//	for _, untreatedAction := range actions {
//		if untreatedAction.ActionType == actionType && untreatedAction.OutCard == opCard && untreatedAction.CombineCards.IsContain(combinations) {
//			return untreatedAction
//		}
//	}
//}

// HasHu 是否有胡
func (actions UserActions) HasHu() bool {
	for _, untreatedAction := range actions {
		if untreatedAction.ActionType == ActionType_Hu {
			return true
		}
	}
	return false
}

// GetHuAction 获取胡动作
func (actions UserActions) GetHuAction() *UserAction {
	for _, untreatedAction := range actions {
		if untreatedAction.ActionType == ActionType_Hu {
			return untreatedAction.Copy()
		}
	}
	return nil
}
