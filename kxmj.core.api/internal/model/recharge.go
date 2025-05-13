package model

import "kxmj.common/recharge"

type GetFistRechargeGiftPackResp struct {
	IsBuy   bool             `json:"isBuy"`   // 是否已购买
	GoodsId string           `json:"goodsId"` // 礼包商品Id
	Items   []*recharge.Item `json:"items"`   // 礼包物品列表
}

type ContinueGiftPack struct {
	WitchDay uint32           `json:"witchDay"`  // 第几天领取
	IsBuy    bool             `json:"isBuy"`     // 是否已购买
	GoodsId  string           `json:"goodsId"`   // 礼包商品Id
	Status   uint32           `json:"status"`    // 领取状态 0 未完成；1 已完成；2 已领取
	Items    []*recharge.Item `json:"packItems"` // 礼包物品列表
}

type GetContinueGiftPackResp struct {
	List []*ContinueGiftPack `json:"list"` // 商品列表
}

type TakeContinueGiftPackReq struct {
	WitchDay uint32 `json:"witchDay"` // 第几天领取
	GoodsId  string `json:"goodsId"`  // 礼包商品Id
}

type TakeContinueGiftPackResp struct {
	Items []*recharge.Item `json:"packItems"` // 礼包物品列表
}
