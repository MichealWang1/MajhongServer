package kxmj_core

import "encoding/json"

type PropClothing struct {
	PropId           uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`                                  // 道具ID
	ClothingPropType uint8  `json:"clothing_prop_type" redis:"clothing_prop_type" gorm:"column:clothing_prop_type"` // 服装道具类型：
	IsDeleted        uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                         // 是否删除： 1 是；2 否
	CreatedAt        uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                         // 创建时间
	UpdatedAt        uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                         // 更新时间
}

func (p *PropClothing) TableName() string {
	return "prop_clothing"
}

func (p *PropClothing) Schema() string {
	return "kxmj_core"
}

func (p *PropClothing) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PropClothing) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
