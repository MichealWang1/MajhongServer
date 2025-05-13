package kxmj_core

import "encoding/json"

type UserProp struct {
	Id         int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                // 主键ID
	UserId     uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`             // 用户ID
	PropId     uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`             // 道具ID
	PropType   uint8  `json:"prop_type" redis:"prop_type" gorm:"column:prop_type"`       // 道具类型
	PropCount  uint32 `json:"prop_count" redis:"prop_count" gorm:"column:prop_count"`    // 道具数量
	ExpireTime uint32 `json:"expire_time" redis:"expire_time" gorm:"column:expire_time"` // 过期时间：0 永不过期
	CreatedAt  uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`    // 创建时间
	UpdatedAt  uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`    // 更新时间
}

func (u *UserProp) TableName() string {
	return "user_prop"
}

func (u *UserProp) Schema() string {
	return "kxmj_core"
}

func (u *UserProp) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserProp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
