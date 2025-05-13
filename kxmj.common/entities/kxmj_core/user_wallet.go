package kxmj_core

import "encoding/json"

type UserWallet struct {
	Id            int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                         // 主键ID
	UserId        uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                      // 用户ID
	Diamond       string `json:"diamond" redis:"diamond" gorm:"column:diamond"`                      // 钻石数
	Gold          string `json:"gold" redis:"gold" gorm:"column:gold"`                               // 金币数
	GoldBean      string `json:"gold_bean" redis:"gold_bean" gorm:"column:gold_bean"`                // 金豆数
	TotalRecharge string `json:"total_recharge" redis:"total_recharge" gorm:"column:total_recharge"` // 累计充值
	RechargeTimes uint32 `json:"recharge_times" redis:"recharge_times" gorm:"column:recharge_times"` // 累计充值笔数
	OnlyOneGoods  string `json:"only_one_goods" redis:"only_one_goods" gorm:"column:only_one_goods"` // 只允许购买一次商品列表：json:["1010011006", "1010011002"]
	UpdatedAt     uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

func (u *UserWallet) TableName() string {
	return "user_wallet"
}

func (u *UserWallet) Schema() string {
	return "kxmj_core"
}

func (u *UserWallet) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserWallet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
