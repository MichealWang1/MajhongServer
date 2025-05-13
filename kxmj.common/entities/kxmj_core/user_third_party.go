package kxmj_core

import "encoding/json"

type UserThirdParty struct {
	Id           int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                         // 主键ID
	UserId       uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                      // 用户ID
	DeviceId     string `json:"device_id" redis:"device_id" gorm:"column:device_id"`                // 设备ID
	WechatOpenId string `json:"wechat_open_id" redis:"wechat_open_id" gorm:"column:wechat_open_id"` // 微信openID
	TiktokId     string `json:"tiktok_id" redis:"tiktok_id" gorm:"column:tiktok_id"`                // 抖音ID
	HuaweiId     string `json:"huawei_id" redis:"huawei_id" gorm:"column:huawei_id"`                // 华为ID
	UpdatedAt    uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
}

func (u *UserThirdParty) TableName() string {
	return "user_third_party"
}

func (u *UserThirdParty) Schema() string {
	return "kxmj_core"
}

func (u *UserThirdParty) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserThirdParty) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
