package login

import (
	"kxmj.common/item"
	"math/rand"
)

type Prize struct {
	Weight uint32 // 权重
	Value  string // 奖励物品数
}

type PrizeType uint32

// 每天登陆奖励
var dailyPrize = map[uint32]map[item.Type]map[uint32]*Prize{
	1: { // 第一天
		item.Gold: { // 物品1(金币)
			1: {Weight: 100, Value: "100000000"},
		},
		item.Diamond: { // 物品2(钻石)
			1: {Weight: 0, Value: "0"},
		},
		item.BP: { // 物品3(BP经验)
			1: {Weight: 100, Value: "30"},
		},
	},
	2: { // 第二天
		item.Gold: { // 物品1(金币)
			1: {Weight: 100, Value: "300000000"},
		},
		item.Diamond: { // 物品2(钻石)
			1: {Weight: 0, Value: "0"},
		},
		item.BP: { // 物品3(BP经验)
			1: {Weight: 100, Value: "50"},
		},
	},
	3: { // 第三天
		item.Gold: { // 物品1(金币)
			1: {Weight: 30, Value: "300000000"}, // 权重1
			2: {Weight: 30, Value: "400000000"}, // 权重2
			3: {Weight: 30, Value: "500000000"}, // 权重3
		},
		item.Diamond: { // 物品2(钻石)
			1: {Weight: 70, Value: "5"},  // 权重1
			2: {Weight: 30, Value: "10"}, // 权重2
		},
		item.BP: { // 物品3(BP经验)
			1: {Weight: 100, Value: "100"},
		},
	},
}

// 累计登陆奖励
var grandTotalConfig = map[uint32]string{
	2: "300000000", // 累计2天 奖励金币
	4: "500000000", // 累计4天 奖励金币
	6: "800000000", // 累计6天 奖励金币
}

// PrizeItem 奖励物品
type PrizeItem struct {
	ItemId uint32 `json:"itemId"` // 物品ID
	Count  string `json:"count"`  // 物品数量
	IsRand bool   `json:"isRand"` // 数量是否随机
}

func calculateWeights(weights map[uint32]*Prize) *Prize {
	var list []uint32
	for k, v := range weights {
		for i := uint32(0); i < v.Weight; i++ {
			list = append(list, k)
		}
	}

	if len(list) <= 0 {
		for _, v := range weights {
			return v
		}
	}

	key := list[rand.Intn(len(list))]
	return weights[key]
}

// GetDailyPrizeConfig 获取每日登陆奖励配置
func GetDailyPrizeConfig() map[uint32][]*PrizeItem {
	result := make(map[uint32][]*PrizeItem, 0)
	for k, v := range dailyPrize {
		var list []*PrizeItem
		goldWeights := v[item.Gold]
		val := &PrizeItem{
			ItemId: 102002,
			Count:  calculateWeights(goldWeights).Value,
			IsRand: false,
		}

		if len(goldWeights) > 1 {
			val.IsRand = true
		}
		list = append(list, val)

		diamondWeights := v[item.Diamond]
		val = &PrizeItem{
			ItemId: 101001,
			Count:  calculateWeights(diamondWeights).Value,
			IsRand: false,
		}
		if len(diamondWeights) > 1 {
			val.IsRand = true
		}
		list = append(list, val)

		bpWeights := v[item.BP]
		val = &PrizeItem{
			ItemId: 401001,
			Count:  calculateWeights(bpWeights).Value,
			IsRand: false,
		}
		if len(bpWeights) > 1 {
			val.IsRand = true
		}
		list = append(list, val)
		result[k] = list
	}

	return result
}

// GetGrandTotalConfig 获取累计登陆奖励配置
func GetGrandTotalConfig() map[uint32]*PrizeItem {
	maps := make(map[uint32]*PrizeItem, 0)
	for k, v := range grandTotalConfig {
		maps[k] = &PrizeItem{
			ItemId: 102002,
			Count:  v,
			IsRand: false,
		}
	}
	return maps
}
