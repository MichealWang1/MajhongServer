package kxmj_core

import "encoding/json"

type GoodsShow struct {
	GoodsId        uint32 `json:"goods_id" redis:"goods_id" gorm:"column:goods_id;primary_key;auto_increment"` // 商品ID
	PropId         uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`                               // 道具ID
	PropType       uint8  `json:"prop_type" redis:"prop_type" gorm:"column:prop_type"`                         // 道具类型：1 币类；2 服装；3 礼包；4 游戏道具
	Price          string `json:"price" redis:"price" gorm:"column:price"`                                     // 价格
	OriginalPrice  string `json:"original_price" redis:"original_price" gorm:"column:original_price"`          // 原价
	PropCount      uint32 `json:"prop_count" redis:"prop_count" gorm:"column:prop_count"`                      // 道具数量
	RewardAdded    string `json:"reward_added" redis:"reward_added" gorm:"column:reward_added"`                // 加赠数量
	FirstBuyDouble uint32 `json:"first_buy_double" redis:"first_buy_double" gorm:"column:first_buy_double"`    // 首购翻倍：1 是；2 否
	Remark         string `json:"remark" redis:"remark" gorm:"column:remark"`                                  // 商品描述
	OnShelfTime    uint32 `json:"on_shelf_time" redis:"on_shelf_time" gorm:"column:on_shelf_time"`             // 上架时间
	OffShelfTime   uint32 `json:"off_shelf_time" redis:"off_shelf_time" gorm:"column:off_shelf_time"`          // 下架时间
	Status         uint8  `json:"status" redis:"status" gorm:"column:status"`                                  // 商品状态：1 启用； 2 不启用
	ExpireTime     uint32 `json:"expire_time" redis:"expire_time" gorm:"column:expire_time"`                   // 过期时间：0 永不过期
	IsDeleted      uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                      // 是否删除： 1 是；2 否
	CreatedAt      uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                      // 创建时间
	UpdatedAt      uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                      // 更新时间
}

func (g *GoodsShow) TableName() string {
	return "goods_show"
}

func (g *GoodsShow) Schema() string {
	return "kxmj_core"
}

func (g *GoodsShow) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GoodsShow) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}
