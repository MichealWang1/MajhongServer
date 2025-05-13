package payment

type Type uint8

const (
	WechatType Type = 1 // 微信支付
	AlipayType Type = 2 // 支付宝支付
)

// CreateOrder

type CreateOrderReq struct {
	UserId  uint32 // 用户ID
	GoodsId string // 商品ID
	Type    uint8  // 支付类型：1 微信支付；2 支付宝支付；
}

type OrderInfo struct {
	OrderId int64  // 订单号
	PayUrl  string // 支付地址
}

type CreateOrderResp struct {
	Code int
	Msg  string
	Data *OrderInfo
}
