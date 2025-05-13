package kxmj_core

import "encoding/json"

type ConfigNotice struct {
	Id        uint32 `json:"id" redis:"id" gorm:"column:id"`                         // ID
	Title     string `json:"title" redis:"title" gorm:"column:title"`                // 公告标题
	Content   string `json:"content" redis:"content" gorm:"column:content"`          // 公告内容
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (c *ConfigNotice) TableName() string {
	return "config_notice"
}

func (c *ConfigNotice) Schema() string {
	return "kxmj_core"
}

func (c *ConfigNotice) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigNotice) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
