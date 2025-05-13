package kxmj_core

import "encoding/json"

type Goods struct {
	GoodsId        string `json:"goods_id" redis:"goods_id" gorm:"column:goods_id;primary_key"`             // 商品ID
	Name           string `json:"name" redis:"name" gorm:"column:name"`                                     // 商品名称
	Remark         string `json:"remark" redis:"remark" gorm:"column:remark"`                               // 商品描述
	Category       uint8  `json:"category" redis:"category" gorm:"column:category"`                         // 商品分类：0 不显示菜单；1 钻石；2 金币；3 装扮；
	ShopType       uint8  `json:"shop_type" redis:"shop_type" gorm:"column:shop_type"`                      // 销售类型：1 RMB购买；2 钻石购买；3 金币购买；4 金豆购买；
	ItemId         uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id"`                            // 物品ID
	Price          string `json:"price" redis:"price" gorm:"column:price"`                                  // 价格
	OriginalPrice  string `json:"original_price" redis:"original_price" gorm:"column:original_price"`       // 原价
	RealCount      string `json:"real_count" redis:"real_count" gorm:"column:real_count"`                   // 实际获得数量
	OriginalCount  string `json:"original_count" redis:"original_count" gorm:"column:original_count"`       // 原来获得数量
	Recommend      uint8  `json:"recommend" redis:"recommend" gorm:"column:recommend"`                      // 推荐商品：1 是；2 否；
	FirstBuyDouble uint8  `json:"first_buy_double" redis:"first_buy_double" gorm:"column:first_buy_double"` // 首购翻倍：1 是；2 否
	Status         uint8  `json:"status" redis:"status" gorm:"column:status"`                               // 商品状态：1 上架； 2 下架
	ExpireTime     uint32 `json:"expire_time" redis:"expire_time" gorm:"column:expire_time"`                // 过期时间：0 永不过期
	OnShelfTime    uint32 `json:"on_shelf_time" redis:"on_shelf_time" gorm:"column:on_shelf_time"`          // 上架时间
	OffShelfTime   uint32 `json:"off_shelf_time" redis:"off_shelf_time" gorm:"column:off_shelf_time"`       // 下架时间
	Sort           uint32 `json:"sort" redis:"sort" gorm:"column:sort"`                                     // 排序序号
	IsDeleted      uint8  `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                   // 是否删除： 1 是；2 否
	CreatedAt      uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                   // 创建时间
	UpdatedAt      uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                   // 更新时间
}

func (g *Goods) TableName() string {
	return "goods"
}

func (g *Goods) Schema() string {
	return "kxmj_core"
}

func (g *Goods) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Goods) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}
