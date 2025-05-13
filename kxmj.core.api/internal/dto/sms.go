package dto

type SendSmsParameter struct {
	TelNumber string `json:"telNumber" binding:"required"` // 手机号
	Type      uint8  `json:"type"`                         // 短信类型：1 注册；2 绑定手机号；3 修改密码
}

type SendSmsResult struct {
	TTL uint32 `json:"ttl"` // 过期时间(秒)
}

type CheckSmsParameter struct {
	TelNumber string `json:"telNumber" binding:"required"` // 手机号
	Type      uint8  `json:"type"`                         // 短信类型：1 注册；2 绑定手机号；3 修改密码
	Code      string `json:"code"`                         // 检查结果
}
