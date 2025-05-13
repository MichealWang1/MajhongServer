package kxmj_core

import "encoding/json"

type UserItem struct {
	Id         int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                // 主键ID
	UserId     uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`             // 用户ID
	ItemId     uint32 `json:"item_id" redis:"item_id" gorm:"column:item_id"`             // 物品ID
	ItemType   uint16 `json:"item_type" redis:"item_type" gorm:"column:item_type"`       // 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；
	ItemCount  string `json:"item_count" redis:"item_count" gorm:"column:item_count"`    // 物品数量
	ExpireTime uint32 `json:"expire_time" redis:"expire_time" gorm:"column:expire_time"` // 过期时间：0 永不过期
	CreatedAt  uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`    // 创建时间
	UpdatedAt  uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`    // 更新时间
}

func (u *UserItem) TableName() string {
	return "user_item"
}

func (u *UserItem) Schema() string {
	return "kxmj_core"
}

func (u *UserItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserItem) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
