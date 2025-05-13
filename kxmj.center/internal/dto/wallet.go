package dto

type AddUserWalletParameter struct {
	UserId   uint32 // 用户ID
	Diamond  string // 钻石数
	Gold     string // 金币数
	GoldBean string // 金豆数
}

type SubUserWalletParameter struct {
	UserId   uint32 // 用户ID
	Diamond  string // 钻石数
	Gold     string // 金币数
	GoldBean string // 金豆数
}

type AddRechargeParameter struct {
	UserId uint32 // 用户ID
	Amount string // 充值金额
}

type AddOnlyOnceGoodsParameter struct {
	UserId  uint32 // 用户ID
	GoodsId string // 商品ID
}
