package item

import (
	"fmt"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/utils"
	"math/rand"
)

// GiftPackType 礼包类型 1 充值礼包；2 钻石礼包；3 抽奖礼包；4 推荐礼包；
type GiftPackType uint8

const (
	GiftPackRecharge   GiftPackType = 1 // 充值礼包
	GiftPackDiamond    GiftPackType = 2 // 钻石礼包
	GiftPackDrawReward GiftPackType = 3 // 抽奖礼包
	GiftPackRecommend  GiftPackType = 4 // 推荐礼包
)

type GiftPackContent struct {
	Id    uint32 `json:"id"`    // 物品ID
	Type  Type   `json:"type"`  // 物品类型
	Count string `json:"count"` // 物品数量 如果是随机值在配置时填0，计算时生成值
	Min   string `json:"min"`   // 随机数量最小值
	Max   string `json:"max"`   // 随机数量最大值
}

func (gp *GiftPackContent) IsRandom() bool {
	if utils.Cmp(gp.Max, utils.Zero().String()) > 0 {
		return true
	}

	return false
}

func (gp *GiftPackContent) ParseBaseValueItems(itemsMap map[uint32]*kxmj_core.Item, values map[uint32]*ValueItem, count string) {
	item, has := itemsMap[gp.Id]
	if has == false {
		return
	}

	// 如果是基础物品类型
	if IsPack(gp.Type) == false {
		// 设置物品随机数量
		gp.setRandomCount()

		value := GetValueItem(item)
		if utils.Cmp(count, utils.Zero().String()) > 0 {
			c, ok := utils.Mul(gp.Count, count)
			if ok {
				value.Count = c.String()
			}
		} else {
			value.Count = gp.Count
		}

		old, has := values[item.ItemId]
		if has {
			c, ok := utils.Add(old.Count, value.Count)
			if ok {
				old.Count = c.String()
			}
		} else {
			values[value.ItemId] = value
		}
		return
	}

	// 如果是礼包类型
	contents := JsonToContent(item.Content)
	for _, content := range contents {
		if gp.Id == content.Id {
			return
		}

		content.ParseBaseValueItems(itemsMap, values, count)
	}
}

func (gp *GiftPackContent) setRandomCount() {
	if gp.IsRandom() == false {
		return
	}

	if utils.Cmp(gp.Count, utils.Zero().String()) > 0 {
		return
	}

	val, ok := utils.Sub(gp.Max, gp.Min)
	if ok {
		randVal := rand.Int63n(val.Int64())
		result, ok := utils.Add(val.String(), fmt.Sprintf("%d", randVal))
		if ok {
			gp.Count = result.String()
		}
	}
}
