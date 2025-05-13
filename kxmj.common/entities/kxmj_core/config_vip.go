package kxmj_core

import "encoding/json"

type ConfigVip struct {
	Id                 uint32 `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"`
	Level              uint8  `json:"level" redis:"level" gorm:"column:level"`                                           // VIP等级
	ExperienceRequired string `json:"experience_required" redis:"experience_required" gorm:"column:experience_required"` // 要求经验（ 每充值1元获得100点VIP经验）
	RewardContent      string `json:"reward_content" redis:"reward_content" gorm:"column:reward_content"`                // 奖励内容
	CreatedAt          uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                            // 创建时间
	UpdatedAt          uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                            // 更新时间
}

func (c *ConfigVip) TableName() string {
	return "config_vip"
}

func (c *ConfigVip) Schema() string {
	return "kxmj_core"
}

func (c *ConfigVip) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigVip) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
