package kxmj_core

import "encoding/json"

type ConfigApp struct {
	Id              uint32 `json:"id" redis:"id" gorm:"column:id;primary_key"`                                  // ID
	WechatSecretKey string `json:"wechat_secret_key" redis:"wechat_secret_key" gorm:"column:wechat_secret_key"` // 微信登陆应用密钥APPSecret
	WechatAppId     string `json:"wechat_app_id" redis:"wechat_app_id" gorm:"column:wechat_app_id"`             // 微信开放平台应用唯一标识
	HotRenewAddress string `json:"hot_renew_address" redis:"hot_renew_address" gorm:"column:hot_renew_address"` // 热更新包地址
	BrokenLine      string `json:"broken_line" redis:"broken_line" gorm:"column:broken_line"`                   // 破产线
	CreatedAt       uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                      // 创建时间
	UpdatedAt       uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                      // 更新时间
}

func (c *ConfigApp) TableName() string {
	return "config_app"
}

func (c *ConfigApp) Schema() string {
	return "kxmj_core"
}

func (c *ConfigApp) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigApp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
