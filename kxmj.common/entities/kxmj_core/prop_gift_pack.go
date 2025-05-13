package kxmj_core

import "encoding/json"

type PropGiftPack struct {
	PropId           uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`                                     // 道具ID
	GiftPackPropType uint8  `json:"gift_pack_prop_type" redis:"gift_pack_prop_type" gorm:"column:gift_pack_prop_type"` // 礼包道具类型： 1 充值礼包；2 钻石礼包；3 抽奖礼包；
	Content          string `json:"content" redis:"content" gorm:"column:content"`                                     // 礼包内容：[{"prop_id":1,type":1,"count":100},{"prop_id":1,type":2,"count":100}]
	IsDeleted        uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                            // 是否删除： 1 是；2 否
	CreatedAt        uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                            // 创建时间
	UpdatedAt        uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                            // 更新时间
}

func (p *PropGiftPack) TableName() string {
	return "prop_gift_pack"
}

func (p *PropGiftPack) Schema() string {
	return "kxmj_core"
}

func (p *PropGiftPack) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PropGiftPack) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
