package kxmj_core

import "encoding/json"

type Robot struct {
	Id          uint32 `json:"id" redis:"id" gorm:"column:id"`                               // ID
	UserId      uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                // 用户ID
	Nickname    string `json:"nickname" redis:"nickname" gorm:"column:nickname"`             // 昵称
	Gender      uint8  `json:"gender" redis:"gender" gorm:"column:gender"`                   // 性别
	AvatarAddr  string `json:"avatar_addr" redis:"avatar_addr" gorm:"column:avatar_addr"`    // 头像地址
	AvatarFrame uint8  `json:"avatar_frame" redis:"avatar_frame" gorm:"column:avatar_frame"` // 头像框
	CreatedAt   uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`       // 创建时间
	UpdatedAt   uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`       // 更新时间
}

func (r *Robot) TableName() string {
	return "robot"
}

func (r *Robot) Schema() string {
	return "kxmj_core"
}

func (r *Robot) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Robot) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
