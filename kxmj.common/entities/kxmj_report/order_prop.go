package kxmj_report

import "encoding/json"

type OrderProp struct {
	GoodsId   uint32 `json:"goods_id" redis:"goods_id" gorm:"column:goods_id"`       // 商品ID
	OrderId   int64  `json:"order_id" redis:"order_id" gorm:"column:order_id"`       // 订单号
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`          // 用户ID
	PropId    uint32 `json:"prop_id" redis:"prop_id" gorm:"column:prop_id"`          // 道具ID
	PropType  uint8  `json:"prop_type" redis:"prop_type" gorm:"column:prop_type"`    // 道具类型：1 币类；2 服装；3 礼包；4 游戏道具
	Price     string `json:"price" redis:"price" gorm:"column:price"`                // 道具价格(钻石数)
	PropCount uint32 `json:"prop_count" redis:"prop_count" gorm:"column:prop_count"` // 道具数量
	UseStatus uint8  `json:"use_status" redis:"use_status" gorm:"column:use_status"` // 使用状态：1 未使用；2 已使用
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (o *OrderProp) TableName() string {
	return "order_prop"
}

func (o *OrderProp) Schema() string {
	return "kxmj_report"
}

func (o *OrderProp) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OrderProp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
