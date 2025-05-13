package kxmj_core

import "encoding/json"

type UserEquip struct {
	Id        int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`             // 主键ID
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`          // 用户ID
	Head      int64  `json:"head" redis:"head" gorm:"column:head"`                   // 头部
	Body      int64  `json:"body" redis:"body" gorm:"column:body"`                   // 身上
	Weapon    int64  `json:"weapon" redis:"weapon" gorm:"column:weapon"`             // 武器
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (u *UserEquip) TableName() string {
	return "user_equip"
}

func (u *UserEquip) Schema() string {
	return "kxmj_core"
}

func (u *UserEquip) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserEquip) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
