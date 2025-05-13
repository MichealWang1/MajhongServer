package kxmj_report

import "encoding/json"

type OrderGoods struct {
	OrderId       int64  `json:"order_id" redis:"order_id" gorm:"column:order_id;primary_key"`       // 订单号
	TradeId       string `json:"trade_id" redis:"trade_id" gorm:"column:trade_id"`                   // 三方订单号
	GoodsId       string `json:"goods_id" redis:"goods_id" gorm:"column:goods_id"`                   // 商品ID
	UserId        uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                      // 用户ID
	ShopType      uint8  `json:"shop_type" redis:"shop_type" gorm:"column:shop_type"`                // 销售类型：1 RMB购买；2 钻石购买；3 金币购买；4 金豆购买；
	ItemId        uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id"`                      // 物品ID
	ItemType      uint16 `json:"item_type" redis:"item_type" gorm:"column:item_type"`                // 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；
	ItemCount     string `json:"item_count" redis:"item_count" gorm:"column:item_count"`             // 物品数量
	Price         string `json:"price" redis:"price" gorm:"column:price"`                            // 商品价格（RMB，钻石，金币）
	OriginalPrice string `json:"original_price" redis:"original_price" gorm:"column:original_price"` // 原价（RMB）
	PaymentType   uint8  `json:"payment_type" redis:"payment_type" gorm:"column:payment_type"`       // 类型：1 微信支付； 2 支付宝支付；
	OrderStatus   uint8  `json:"order_status" redis:"order_status" gorm:"column:order_status"`       // 订单状态：1 待付款； 2 失败； 3 成功；
	CompletedAt   uint32 `json:"completed_at" redis:"completed_at" gorm:"column:completed_at"`       // 完成时间
	CreatedAt     uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`             // 创建时间
	UpdatedAt     uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

func (o *OrderGoods) TableName() string {
	return "order_goods"
}

func (o *OrderGoods) Schema() string {
	return "kxmj_report"
}

func (o *OrderGoods) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OrderGoods) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
