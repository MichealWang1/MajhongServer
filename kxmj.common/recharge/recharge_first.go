package recharge

// GiftPackType 充值礼包类型
type GiftPackType uint32

const (
	FirstRecharge    GiftPackType = 1 // 首充礼包
	ContinueRecharge GiftPackType = 2 // 首充连续领取礼包
)

// Item 物品定义
type Item struct {
	Id    uint32 `json:"id"`    // 物品ID
	Count string `json:"count"` // 物品数量
}

// 充值类礼包扩展字段说明:
// 1.首充礼包:{"1":1} key:1固定值,value:充值礼包类型
// 2.首充连续领取礼包:{"1":2,"2":1} key:1,value:充值礼包类型； key:2,value 礼包配表键值,如value=1对应首充连续领取礼包1

// ContinueGiftPack 首充连续领取礼包配表
var ContinueGiftPack = map[uint32]map[uint32][]*Item{
	1: { // 首充连续领取礼包1
		1: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第一天奖品
		2: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第二天奖品
		3: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第三天奖品
	},
	2: { // 首充连续领取礼包2
		1: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第一天奖品
		2: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第二天奖品
		3: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第三天奖品
	},
	3: { // 首充连续领取礼包3
		1: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第一天奖品
		2: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第二天奖品
		3: {{Id: 102002, Count: "10000"}, {Id: 102002, Count: "10000"}}, // 第三天奖品
	},
}
