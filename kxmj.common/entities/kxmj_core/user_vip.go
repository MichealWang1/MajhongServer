package kxmj_core

import "encoding/json"

type UserVip struct {
	Id        int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`             // 主键ID
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`          // 用户ID
	Level     uint32 `json:"level" redis:"level" gorm:"column:level"`                // VIP等级
	CurBp     uint32 `json:"cur_bp" redis:"cur_bp" gorm:"column:cur_bp"`             // 当前经验
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (u *UserVip) TableName() string {
	return "user_vip"
}

func (u *UserVip) Schema() string {
	return "kxmj_core"
}

func (u *UserVip) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserVip) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
