package item

import (
	"context"
	"encoding/json"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mq"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"time"
)

// Type 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；303 复活卡；304 王者归来；
type Type uint16

const (
	Diamond             Type = 101 // 钻石
	Gold                Type = 102 // 金币
	GoldBean            Type = 103 // 金豆
	Adornment           Type = 201 // 装扮
	Weapon              Type = 202 // 武器
	AvatarFrame         Type = 203 // 头像框
	CardTableDecoration Type = 204 // 牌桌装饰品类
	CardBackDecoration  Type = 205 // 牌背装饰品类
	EntranceAnimation   Type = 206 // 入场动画类
	InteractiveProps    Type = 207 // 交互道具类
	SpecialGiftPack     Type = 301 // 特权包
	GiftPack            Type = 302 // 礼包
	RisePack            Type = 303 // 复活卡
	KingBackPack        Type = 304 // 王者归来
	BP                  Type = 401 // BP经验
)

// 可叠加物品定义
var pileItems = []Type{
	Diamond,
	Gold,
	GoldBean,
	InteractiveProps,
	BP,
}

// 金币类物品定义
var goldItems = []Type{
	Diamond,
	Gold,
	GoldBean,
}

// 复活卡类物品定义
var riseItems = []Type{
	RisePack,
	KingBackPack,
}

// 包裹类物品定义
var packItems = []Type{
	SpecialGiftPack,
	GiftPack,
	RisePack,
	KingBackPack,
}

// BP经验类物品定义
var bpItems = []Type{
	BP,
}

// ValueItem 物品值类型对象
type ValueItem struct {
	ItemId        uint32             `json:"itemId"`        // 物品ID
	Name          string             `json:"name"`          // 物品名称
	ItemType      Type               `json:"itemType"`      // 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；303 复活卡；304 王者归来；
	ServiceLife   uint32             `json:"serviceLife"`   // 使用寿命（秒为单位）
	Content       []*GiftPackContent `json:"content"`       // 礼包、特权卡类道具内容
	Extra         map[uint32]uint32  `json:"extra"`         // 扩展属性(攻击等属性)：json格式：{"1":1, "2"：2}
	GiftType      GiftPackType       `json:"giftType"`      // 礼包类型：0 未定义；1 充值礼包；2 钻石礼包；3 抽奖礼包；
	AdornmentType AdornmentType      `json:"adornmentType"` // 装扮物品类型：1 头部；2 衣服；
	Count         string             `json:"count"`         // 物品数量
	CreatedAt     uint32             `json:"createdAt"`     // 创建时间
	UpdatedAt     uint32             `json:"updatedAt"`     // 更新时间
}

func ContentToJson(content []*GiftPackContent) string {
	if len(content) <= 0 {
		return ""
	}

	data, _ := json.Marshal(content)
	return string(data)
}

func JsonToContent(str string) []*GiftPackContent {
	if len(str) <= 0 {
		return nil
	}

	var content []*GiftPackContent
	err := json.Unmarshal([]byte(str), &content)
	if err != nil {
		return nil
	}
	return content
}

func ExtraToJson(extra map[uint32]uint32) string {
	if len(extra) <= 0 {
		return ""
	}

	data, _ := json.Marshal(extra)
	return string(data)
}

func JsonToExtra(str string) map[uint32]uint32 {
	if len(str) <= 0 {
		return nil
	}

	var content map[uint32]uint32
	err := json.Unmarshal([]byte(str), &content)
	if err != nil {
		return nil
	}
	return content
}

// GetValueItem 获取物品值类型对象
func GetValueItem(item *kxmj_core.Item) *ValueItem {
	return &ValueItem{
		ItemId:        item.ItemId,
		Name:          item.Name,
		ItemType:      Type(item.ItemType),
		ServiceLife:   item.ServiceLife,
		Content:       JsonToContent(item.Content),
		Extra:         JsonToExtra(item.Extra),
		GiftType:      GiftPackType(item.GiftType),
		AdornmentType: AdornmentType(item.AdornmentType),
		Count:         "1",
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

// GetItem 获取所有基础物品，礼包、特权卡类拆解成所有基础道具
func (vt *ValueItem) GetItem() *kxmj_core.Item {
	return &kxmj_core.Item{
		ItemId:        vt.ItemId,
		Name:          vt.Name,
		ItemType:      uint16(vt.ItemType),
		ServiceLife:   vt.ServiceLife,
		Content:       ContentToJson(vt.Content),
		Extra:         ExtraToJson(vt.Extra),
		GiftType:      uint8(vt.GiftType),
		AdornmentType: uint8(vt.AdornmentType),
		CreatedAt:     vt.CreatedAt,
		UpdatedAt:     vt.UpdatedAt,
	}
}

// ParseBaseValueItems 解析所有基础物品，礼包、特权卡类拆解成所有基础道具
func (vt *ValueItem) ParseBaseValueItems(count string) ([]*ValueItem, error) {
	var values []*ValueItem
	if IsPack(vt.ItemType) == false {
		if utils.Cmp(count, utils.Zero().String()) > 0 {
			c, ok := utils.Mul(vt.Count, count)
			if ok {
				vt.Count = c.String()
			}
		}

		values = append(values, vt)
		return values, nil
	}

	itemsMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	maps := make(map[uint32]*ValueItem, 0)
	for _, content := range vt.Content {
		content.ParseBaseValueItems(itemsMap, maps, count)
	}

	for _, v := range maps {
		values = append(values, v)
	}

	return values, nil
}

// GetGoldItems 获取币类物品
func GetGoldItems(items []*ValueItem) []*ValueItem {
	var values []*ValueItem
	for _, item := range items {
		if IsGold(item.ItemType) {
			values = append(values, item)
		}
	}
	return values
}

// GetPropItems 获取道具类物品
func GetPropItems(items []*ValueItem) []*ValueItem {
	var values []*ValueItem
	for _, item := range items {
		if IsGold(item.ItemType) {
			continue
		}

		if IsBP(item.ItemType) {
			continue
		}

		values = append(values, item)
	}
	return values
}

// GetBPItems 获取BP经验类物品
func GetBPItems(items []*ValueItem) []*ValueItem {
	var values []*ValueItem
	for _, item := range items {
		if IsBP(item.ItemType) {
			values = append(values, item)
		}
	}
	return values
}

// UpdateUserItems 更新用户背包物品
func UpdateUserItems(ctx context.Context, values []*ValueItem, userId uint32) error {
	props := GetPropItems(values)
	if len(props) <= 0 {
		return nil
	}

	cache := redis_cache.GetCache().GetUserCache().ItemCache()
	cache.Lock(ctx, userId)

	items, err := cache.GetAll(ctx, userId)
	if err != nil {
		return err
	}

	now := uint32(time.Now().Unix())
	saves := make([]*kxmj_core.UserItem, 0)
	for _, p := range props {
		data, has := items[p.ItemId]
		if has {
			// 如果是可叠加商品，叠加数量
			if IsCanPile(Type(data.ItemType)) {
				count, ok := utils.Add(data.ItemCount, data.ItemCount)
				if ok == false {
					log.Sugar().Errorf("Add p:%v o:%v err", p, data)
				} else {
					data.ItemCount = count.String()
				}
			}
		} else {
			data = &kxmj_core.UserItem{
				Id:         utils.Snowflake.Generate().Int64(),
				UserId:     userId,
				ItemId:     p.ItemId,
				ItemType:   uint16(p.ItemType),
				ItemCount:  p.Count,
				ExpireTime: 0,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
		}

		// 物品过期时间计算
		var expireTime uint32
		if p.ServiceLife > 0 {
			// 时间叠加
			if data.ExpireTime > now {
				expireTime = data.ExpireTime + p.ServiceLife
			} else {
				expireTime = now + p.ServiceLife
			}
		} else {
			// 永久生效
			expireTime = 0
		}

		data.ExpireTime = expireTime
		saves = append(saves, data)
	}

	err = cache.BulkSet(ctx, saves, userId)
	if err != nil {
		return err
	}

	cache.Unlock(ctx, userId)

	for _, d := range saves {
		err = mq.SyncTable(d, mq.AddOrUpdate)
		if err != nil {
			log.Sugar().Errorf("SyncTable err:%v", err)
		}
	}

	return nil
}

// IsPack 是否是包裹类型物品
func IsPack(itemType Type) bool {
	for _, t := range packItems {
		if t == itemType {
			return true
		}
	}
	return false
}

// IsCanPile 是否是可叠加类型物品
func IsCanPile(itemType Type) bool {
	for _, t := range pileItems {
		if t == itemType {
			return true
		}
	}
	return false
}

// IsGold 是否是金币类物品
func IsGold(itemType Type) bool {
	for _, t := range goldItems {
		if t == itemType {
			return true
		}
	}
	return false
}

// IsRise 是否是复活卡类物品
func IsRise(itemType Type) bool {
	for _, t := range riseItems {
		if t == itemType {
			return true
		}
	}
	return false
}

// IsBP 是否是BP经验类物品
func IsBP(itemType Type) bool {
	for _, t := range bpItems {
		if t == itemType {
			return true
		}
	}
	return false
}
