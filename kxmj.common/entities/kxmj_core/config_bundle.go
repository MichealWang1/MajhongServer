package kxmj_core

import "encoding/json"

type ConfigBundle struct {
	BundleId        string `json:"bundle_id" redis:"bundle_id" gorm:"column:bundle_id;primary_key"`             // 分发包ID
	Version         string `json:"version" redis:"version" gorm:"column:version"`                               // 版本号
	AppName         string `json:"app_name" redis:"app_name" gorm:"column:app_name"`                            // 上架名称
	BundleChannel   uint32 `json:"bundle_channel" redis:"bundle_channel" gorm:"column:bundle_channel"`          // 分包渠道：1 AppStore；2 华为；3 小米；4 OPPO；
	HotRenewAddress string `json:"hot_renew_address" redis:"hot_renew_address" gorm:"column:hot_renew_address"` // 热更新包地址
	IsDeleted       uint32 `json:"is_deleted" redis:"is_deleted" gorm:"column:is_deleted"`                      // 是否删除： 1 是；2 否
	CreatedAt       uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                      // 创建时间
	UpdatedAt       uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                      // 更新时间
}

func (c *ConfigBundle) TableName() string {
	return "config_bundle"
}

func (c *ConfigBundle) Schema() string {
	return "kxmj_core"
}

func (c *ConfigBundle) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigBundle) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
