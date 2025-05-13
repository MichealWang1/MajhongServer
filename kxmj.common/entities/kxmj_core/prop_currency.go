package kxmj_core

import "encoding/json"

type PropCurrency struct {
	PropId           uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id;primary_key"`                      // 道具ID
	CurrencyPropType uint8  `json:"currency_prop_type" redis:"currency_prop_type" gorm:"column:currency_prop_type"` // 货币道具类型：1 钻石；2 金豆；
	Rate             string `json:"rate" redis:"rate" gorm:"column:rate"`                                           // 汇率：（1RMB兑换数量）
	IsDeleted        uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                         // 是否删除： 1 是；2 否
	CreatedAt        uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                         // 创建时间
	UpdatedAt        uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                         // 更新时间
}

func (p *PropCurrency) TableName() string {
	return "prop_currency"
}

func (p *PropCurrency) Schema() string {
	return "kxmj_core"
}

func (p *PropCurrency) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PropCurrency) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
