package kxmj_core

import "encoding/json"

type TaskRecharge struct {
	TaskId    uint32 `json:"task_id" redis:"task_id" gorm:"column:task_id;primary_key"` // 任务ID
	TaskType  uint8  `json:"task_type" redis:"task_type" gorm:"column:task_type"`       // 任务类型：1 登录天数；2 对局数；3 赢金币；4 赢倍数；5 充值数；
	Condition string `json:"condition" redis:"condition" gorm:"column:condition"`       // 累计充值消费的现金达到条件
	ItemId    uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id"`             // 赠送物品ID
	ItemCount uint32 `json:"item_count" redis:"item_count" gorm:"column:item_count"`    // 赠送物品数量
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`    // 创建时间
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`    // 更新时间
}

func (t *TaskRecharge) TableName() string {
	return "task_recharge"
}

func (t *TaskRecharge) Schema() string {
	return "kxmj_core"
}

func (t *TaskRecharge) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TaskRecharge) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
