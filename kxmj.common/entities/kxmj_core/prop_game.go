package kxmj_core

import "encoding/json"

type PropGame struct {
	PropId       uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`                      // 道具ID
	GamePropType uint8  `json:"game_prop_type" redis:"game_prop_type" gorm:"column:game_prop_type"` // 游戏道具类型：
	IsDeleted    uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`             // 是否删除： 1 是；2 否
	CreatedAt    uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`             // 创建时间
	UpdatedAt    uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

func (p *PropGame) TableName() string {
	return "prop_game"
}

func (p *PropGame) Schema() string {
	return "kxmj_core"
}

func (p *PropGame) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PropGame) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
