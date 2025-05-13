package kxmj_core

import "encoding/json"

type Device struct {
	Id           uint32 `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"`
	DeviceId     string `json:"device_id" redis:"device_id" gorm:"column:device_id"`          // 设备号
	Os           uint8  `json:"os" redis:"os" gorm:"column:os"`                               // 系统: 0 未知 1 安卓 2 IOS 3 其它
	Brand        string `json:"brand" redis:"brand" gorm:"column:brand"`                      // 品牌商
	Version      string `json:"version" redis:"version" gorm:"column:version"`                // 系统版本
	Model        string `json:"model" redis:"model" gorm:"column:model"`                      // 手机型号
	Width        uint32 `json:"width" redis:"width" gorm:"column:width"`                      // 宽度
	Height       uint32 `json:"height" redis:"height" gorm:"column:height"`                   // 分辨率
	Manufacturer string `json:"manufacturer" redis:"manufacturer" gorm:"column:manufacturer"` // 设备制造商
	AndroidSdk   uint8  `json:"android_sdk" redis:"android_sdk" gorm:"column:android_sdk"`    // android sdk 版本
	AndroidId    string `json:"android_id" redis:"android_id" gorm:"column:android_id"`       // android id
	AndroidImei  string `json:"android_imei" redis:"android_imei" gorm:"column:android_imei"` // 自定义IMEI码
	IosUuid      string `json:"ios_uuid" redis:"ios_uuid" gorm:"column:ios_uuid"`             // IOS设备ID
	Organic      uint8  `json:"organic" redis:"organic" gorm:"column:organic"`                // 自然量 1是，2非
	BundleId     string `json:"bundle_id" redis:"bundle_id" gorm:"column:bundle_id"`          // 分包ID
	CreatedAt    uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`       // 创建时间
	UpdatedAt    uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`       // 更新时间
}

func (d *Device) TableName() string {
	return "device"
}

func (d *Device) Schema() string {
	return "kxmj_core"
}

func (d *Device) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Device) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
