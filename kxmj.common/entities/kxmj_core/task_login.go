package kxmj_core

import "encoding/json"

type TaskLogin struct {
	TaskId    uint32 `json:"task_id" redis:"task_id" gorm:"column:task_id;primary_key;auto_increment"` // 任务ID
	TaskType  uint8  `json:"task_type" redis:"task_type" gorm:"column:task_type"`                      // 任务类型：1 登录天数；2 对局数；3 赢金币；4 赢倍数；5 充值数；
	Condition uint32 `json:"condition" redis:"condition" gorm:"column:condition"`                      // 连续登陆天数
	ItemId    uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id"`                            // 赠送物品ID
	ItemCount uint32 `json:"item_count" redis:"item_count" gorm:"column:item_count"`                   // 赠送物品数量
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                   // 创建时间
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                   // 更新时间
}

func (t *TaskLogin) TableName() string {
	return "task_login"
}

func (t *TaskLogin) Schema() string {
	return "kxmj_core"
}

func (t *TaskLogin) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TaskLogin) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
