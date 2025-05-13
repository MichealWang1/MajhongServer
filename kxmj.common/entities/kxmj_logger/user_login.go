package kxmj_logger

import "encoding/json"

type UserLogin struct {
	Id        int32  `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"` // 日志ID
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`             // 用户ID
	Ip        string `json:"ip" redis:"ip" gorm:"column:ip"`                            // 登录IP
	DeviceId  string `json:"device_id" redis:"device_id" gorm:"column:device_id"`       // 设备ID
	LoginType uint8  `json:"login_type" redis:"login_type" gorm:"column:login_type"`    // 登录方式；1 手机登录；2 微信登录；3 token 登录;
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`    // 创建时间
}

func (u *UserLogin) TableName() string {
	return "user_login"
}

func (u *UserLogin) Schema() string {
	return "kxmj_logger"
}

func (u *UserLogin) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserLogin) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
