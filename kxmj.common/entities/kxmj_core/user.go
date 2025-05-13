package kxmj_core

import "encoding/json"

type User struct {
	Id                int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                                        // 主键ID
	UserId            uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                                     // 用户ID
	Nickname          string `json:"nickname" redis:"nickname" gorm:"column:nickname"`                                  // 昵称
	Gender            uint8  `json:"gender" redis:"gender" gorm:"column:gender"`                                        // 性别：0 女；1 男
	AvatarAddr        string `json:"avatar_addr" redis:"avatar_addr" gorm:"column:avatar_addr"`                         // 头像地址
	AvatarFrame       uint8  `json:"avatar_frame" redis:"avatar_frame" gorm:"column:avatar_frame"`                      // 头像框
	RealName          string `json:"real_name" redis:"real_name" gorm:"column:real_name"`                               // 实名
	IdCard            string `json:"id_card" redis:"id_card" gorm:"column:id_card"`                                     // 身份证ID
	UserMod           uint8  `json:"user_mod" redis:"user_mod" gorm:"column:user_mod"`                                  // 人物样式
	AccountType       uint8  `json:"account_type" redis:"account_type" gorm:"column:account_type"`                      // 账号类型
	Vip               uint8  `json:"vip" redis:"vip" gorm:"column:vip"`                                                 // VIP等级
	DeviceId          string `json:"device_id" redis:"device_id" gorm:"column:device_id"`                               // 注册的设备ID
	RegisterIp        string `json:"register_ip" redis:"register_ip" gorm:"column:register_ip"`                         // 注册IP
	RegisterType      uint8  `json:"register_type" redis:"register_type" gorm:"column:register_type"`                   // 注册方式：1 人工创建；2 手机号；3 第三方登陆
	TelNumber         string `json:"tel_number" redis:"tel_number" gorm:"column:tel_number"`                            // 手机号
	Status            uint8  `json:"status" redis:"status" gorm:"column:status"`                                        // 状态。1 正常；2 冻结
	BindingAt         uint32 `json:"binding_at" redis:"binding_at" gorm:"column:binding_at"`                            // 绑定手机时间
	LoginPassword     string `json:"login_password" redis:"login_password" gorm:"column:login_password"`                // 登录密码
	LoginPasswordSalt string `json:"login_password_salt" redis:"login_password_salt" gorm:"column:login_password_salt"` // 登录密码盐
	Remark            string `json:"remark" redis:"remark" gorm:"column:remark"`                                        // 备注
	BundleId          string `json:"bundle_id" redis:"bundle_id" gorm:"column:bundle_id"`                               // 分包ID
	BundleChannel     uint32 `json:"bundle_channel" redis:"bundle_channel" gorm:"column:bundle_channel"`                // 分包渠道：1 AppStore；2 华为；3 小米；4 OPPO；
	Organic           uint8  `json:"organic" redis:"organic" gorm:"column:organic"`                                     // 自然量 1是，2非
	CreatedAt         uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`                            // 创建时间
	UpdatedAt         uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`                            // 更新时间
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Schema() string {
	return "kxmj_core"
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
