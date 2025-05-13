package model

// TelRegisterReq 手机注册
type TelRegisterReq struct {
	TelNumber string `json:"telNumber" binding:"required"` // 手机号码
	SMSCode   string `json:"SMSCode" binding:"required"`   // 短信验证码
	Password  string `json:"password" binding:"required"`  // 密码
}

// TelRegisterResp 手机注册
type TelRegisterResp struct {
	UserId uint32 `json:"userId"` // user id
}

// LoginReq 登录请求
type LoginReq struct {
	TelNumber string `json:"telNumber" binding:"required"` // 手机号码
	Password  string `json:"password" binding:"required"`  // 密码
}

// LoginResp 登录回复
type LoginResp struct {
	UserId uint32 `json:"userId"` // user id
	Token  string `json:"token"`  // token
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	TelNumber   string `json:"telNumber" binding:"required"`   // 手机号码
	SMSCode     string `json:"SMSCode" binding:"required"`     // 短信验证码
	NewPassword string `json:"NewPassword" binding:"required"` //新的密码
}

// ChangePasswordResp 修改密码回复
type ChangePasswordResp struct {
	Message string `json:"Message"` //修改成功通知
}

// BindPhoneNumReq 绑定请求
type BindPhoneNumReq struct {
	TelNumber string `json:"telNumber" binding:"required"` // 手机号码
	SMSCode   string `json:"SMSCode" binding:"required"`   //短信验证码
	UserId    uint32 `json:"userId" binding:"required"`    //user id
}

// BindPhoneNumResp 绑定回复
type BindPhoneNumResp struct {
	Message string `json:"Message"` //绑定成功通知
}

// GetInfoResp 获得用户信息回复
type GetInfoResp struct {
	UserId      uint32 `json:"userId"`      //UID
	Nickname    string `json:"nickname"`    //昵称
	AvatarAddr  string `json:"avatarAddr"`  //头像地址
	Vip         uint8  `json:"vip"`         //VIP等级
	AvatarFrame uint8  `json:"avatarFrame"` //头像框
	Diamond     string `json:"diamond"`     //钻石数
	Gold        string `json:"gold"`        //金币数
	GoldBean    string `json:"goldBean"`    //金豆数
	Gender      uint8  `json:"gender"`      //性别：1 男；2 女
}
