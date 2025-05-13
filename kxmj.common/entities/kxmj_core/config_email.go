package kxmj_core

import "encoding/json"

type ConfigEmail struct {
	EmailId      uint32 `json:"email_id" redis:"email_id" gorm:"column:email_id;primary_key;auto_increment"` // 邮件ID
	EmailType    uint8  `json:"email_type" redis:"email_type" gorm:"column:email_type"`                      // 邮件类型：1 福利发放；2 系统通知
	Title        string `json:"title" redis:"title" gorm:"column:title"`                                     // 邮件标题
	Remark       string `json:"remark" redis:"remark" gorm:"column:remark"`                                  // 描述
	IsReward     uint8  `json:"is_reward" redis:"is_reward" gorm:"column:is_reward"`                         // 是否奖励：1 是；2 否
	IsSingleSend int8   `json:"is_single_send" redis:"is_single_send" gorm:"column:is_single_send"`          // 是否单独发送 0群发 1单独发送
	ItemList     string `json:"item_list" redis:"item_list" gorm:"column:item_list"`                         // 奖励物品格式：id是物品ID，count是物品数量[{"id":1001,"count":30},{"id":1001,"count":30}]
	CreatedAt    uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                      // 创建时间
	ExpireAt     uint32 `json:"expire_at" redis:"expire_at" gorm:"column:expire_at"`                         // 过期时间
	UpdatedAt    uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                      // 更新时间
	SendAt       uint32 `json:"send_at" redis:"send_at" gorm:"column:send_at"`                               // 发送时间
}

func (c *ConfigEmail) TableName() string {
	return "config_email"
}

func (c *ConfigEmail) Schema() string {
	return "kxmj_core"
}

func (c *ConfigEmail) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigEmail) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
