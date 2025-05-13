package kxmj_core

import "encoding/json"

type Item struct {
	ItemId        uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id;primary_key"`          // 物品ID
	Name          string `json:"name" redis:"name" gorm:"column:name"`                               // 物品名称
	ItemType      uint16 `json:"item_type" redis:"item_type" gorm:"column:item_type"`                // 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；303 复活卡；304 王者归来；
	ServiceLife   uint32 `json:"service_life" redis:"service_life" gorm:"column:service_life"`       // 使用寿命（秒为单位）
	Content       string `json:"content" redis:"content" gorm:"column:content"`                      // 礼包、特权卡类道具内容：json格式：[{"id":101001,"type":101,"count":"66", "min":"0", "max":"0"},{"id":102002,"type":102,"count":"6000000", "min":"0", "max":"0"}]
	Extra         string `json:"extra" redis:"extra" gorm:"column:extra"`                            // 扩展属性(攻击等属性)：json格式：{"1":1,"2":1}
	GiftType      uint8  `json:"gift_type" redis:"gift_type" gorm:"column:gift_type"`                // 礼包类型 1 充值礼包；2 钻石礼包；3 抽奖礼包；4 推荐礼包；
	AdornmentType uint8  `json:"adornment_type" redis:"adornment_type" gorm:"column:adornment_type"` // 装扮物品类型：1 头部；2 衣服；
	CreatedAt     uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`             // 创建时间
	UpdatedAt     uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

func (i *Item) TableName() string {
	return "item"
}

func (i *Item) Schema() string {
	return "kxmj_core"
}

func (i *Item) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Item) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, i)
}
