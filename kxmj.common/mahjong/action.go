package lib

type ActionType int32

const (
	ActionType_Pass    ActionType = 0   // 过
	ActionType_Chi     ActionType = 1   // 吃
	ActionType_Peng    ActionType = 2   // 碰
	ActionType_Gang    ActionType = 3   // 杠
	ActionType_Hu      ActionType = 4   // 胡
	ActionType_Ting    ActionType = 5   // 听
	ActionType_Unknown ActionType = 255 // 未知动作
)

type ExtraActionType int32

const (
	ExtraActionType_Null      ExtraActionType = 0 // 无
	ExtraActionType_Ming_Gang ExtraActionType = 1 // 明杠
	ExtraActionType_Bu_Gang   ExtraActionType = 2 // 补杠
	ExtraActionType_An_Gang   ExtraActionType = 3 // 暗杠
)
