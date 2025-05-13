package lobby

// ShoppingSuccess

type ShoppingSuccessReq struct {
	UserId         uint32 // 用户ID
	GoodsId        string // 商品ID
	Name           string // 商品名称
	ShopType       uint8  // 销售方式：1 RMB购买；2 钻石购买；3 金币购买
	Price          string // 价格
	RealCount      string // 实际获得数量
	FirstBuyDouble uint8  // 首购翻倍：1 是；2 否
	ItemId         uint32 // 物品ID
}

type ShoppingSuccessResp struct {
	Code int
	Msg  string
}
