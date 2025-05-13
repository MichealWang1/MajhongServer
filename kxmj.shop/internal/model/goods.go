package model

// 该部分代码只供参考，开发时删除掉

type GoodsData struct {
	GoodsId        string `json:"goodsId"`        // 商品ID
	Name           string `json:"name"`           // 商品名称
	ItemId         uint32 `json:"itemId"`         // 物品ID
	ShopType       uint8  `json:"shopType"`       // 销售方式：1 RMB购买；2 钻石购买；3 金币购买
	Price          string `json:"price"`          // 价格
	OriginalPrice  string `json:"originalPrice"`  // 原价
	RealCount      string `json:"realCount"`      // 实际获得数量
	OriginalCount  string `json:"originalCount"`  // 原来获得数量
	RewardAdded    string `json:"rewardAdded"`    // 加赠数量
	IncomeTimes    string `json:"incomeTimes"`    // 收益倍数
	Recommend      uint8  `json:"recommend"`      // 推荐商品：1 是；2 否；
	FirstBuyDouble uint8  `json:"firstBuyDouble"` // 首购翻倍：1 是；2 否
	ExpireTime     uint32 `json:"expireTime"`     // 过期时间：0 永不过期
	Category       uint8  `json:"category"`       // 商品分类：0 不显示菜单；1 钻石；2 金币；3 装扮；
	CategoryName   string `json:"categoryName"`   // 商品分类名称
}

type GoodsListResp struct {
	List []*GoodsData `json:"list"` // 商品列表
}

type BuyReq struct {
	GoodsId string `json:"goodsId"` // 商品ID
	Type    uint8  `json:"type"`    // 支付类型：1 微信支付；2 支付宝支付；
}

type BuyResp struct {
	OrderId string `json:"orderId"` // 订单号
	PayUrl  string `json:"payUrl"`  // 支付地址 (销售方式：1 RMB购买 有值)
}
