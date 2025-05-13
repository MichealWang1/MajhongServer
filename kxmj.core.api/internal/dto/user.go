package dto

type CreateUserParameter struct {
	Id                int64  // 主键ID
	UserId            uint32 // 用户ID
	Nickname          string // 昵称
	Gender            uint8  // 性别
	AvatarAddr        string // 头像地址
	AvatarFrame       uint8  // 头像框
	RealName          string // 实名
	IdCard            string // 身份证ID
	UserMod           uint8  // 人物样式
	AccountType       uint8  // 账号类型
	Vip               uint8  // VIP等级
	DeviceId          string // 注册的设备ID
	RegisterIp        string // 注册IP
	RegisterType      uint8  // 注册方式：1 人工创建；2 手机号；3 第三方登陆
	TelNumber         string // 手机号
	Status            uint8  // 状态。1 正常；2 冻结
	BindingAt         uint32 // 绑定手机时间
	LoginPassword     string // 登录密码
	LoginPasswordSalt string // 登录密码盐
	Remark            string // 备注
	BundleId          string // 分包ID
	BundleChannel     uint32 // 分包渠道：1 AppStore；2 华为；3 小米；4 OPPO；
	Organic           uint8  // 自然量 1是，2非
	WechatOpenId      string // 微信openID
	TiktokId          string // 抖音ID
	HuaweiId          string // 华为ID
	Diamond           string // 钻石数
	Gold              string // 金币数
	GoldBean          string // 金豆数
	TotalRecharge     string // 累计充值
	RechargeTimes     uint32 // 累计充值笔数
	Head              int64  // 头部
	Body              int64  // 身上
	Weapon            int64  // 武器
	CreatedAt         uint32 // 创建时间
	UpdatedAt         uint32 // 更新时间
}

type CreateUserResult struct {
	UserId uint32 // 用户ID
}
