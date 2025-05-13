package center

// GetUserInfo

type GetUserInfoReq struct {
	UserId uint32 // 用户ID
}

type GetUserInfoResp struct {
	Code int
	Msg  string
	Data *GetUserInfoData
}

type GetUserInfoData struct {
	UserId      uint32 // 用户ID
	Nickname    string // 昵称
	Gender      uint8  // 性别
	AvatarAddr  string // 头像地址
	AvatarFrame uint8  // 头像框
	RealName    string // 实名
	UserMod     uint8  // 人物样式
	Vip         uint8  // VIP等级
	TelNumber   string // 手机号
	Status      uint8  // 状态。1 正常；2 冻结
}

// CheckUserGold

type CheckUserGoldReq struct {
	UserId uint32 // 用户ID
}

type CheckUserGoldResp struct {
	Code int
	Msg  string
	Data *CheckUserGoldData
}

type CheckUserGoldData struct {
	UserId uint32 // 用户ID
	Gold   string // 金豆数
}

// GetUserGold

type GetUserGoldReq struct {
	UserId    uint32 // 用户ID
	RoomId    uint32 // 房间ID
	GameId    uint16 // 游戏ID
	GameType  uint8  // 游戏类型：1 麻将；2 斗地主
	RoomLevel uint8  // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
}

type GetUserGoldResp struct {
	Code int
	Msg  string
	Data *GetUserGoldData
}

type GetUserGoldData struct {
	UserId uint32 // 用户ID
	Gold   string // 金豆数
}

// SetUserGold

type SetUserGoldReq struct {
	UserId    uint32 // 用户ID
	Gold      string // 金豆数
	RoomId    uint32 // 房间ID
	GameId    uint16 // 游戏ID
	GameType  uint8  // 游戏类型：1 麻将；2 斗地主
	RoomLevel uint8  // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
}

type SetUserGoldResp struct {
	Code int
	Msg  string
}

// AddUserWallet

type AddUserWalletReq struct {
	UserId       uint32 // 用户ID
	OrderId      int64  // 订单号
	Diamond      string // 钻石数
	Gold         string // 金币数
	GoldBean     string // 金豆数
	BusinessType uint8  // 业务类型：1 商城；2 任务；3 邮件；
}

type AddUserWalletData struct {
	Diamond  string // 钻石数
	Gold     string // 金币数
	GoldBean string // 金豆数
}

type AddUserWalletResp struct {
	Code int
	Msg  string
	Data *AddUserWalletData
}

// SubUserWallet

type SubUserWalletReq struct {
	UserId       uint32 // 用户ID
	OrderId      int64  // 订单号
	Diamond      string // 钻石数
	Gold         string // 金币数
	GoldBean     string // 金豆数
	BusinessType uint8  // 业务类型：1 商城；2 任务；3 邮件；
}

type SubUserWalletData struct {
	Diamond  string // 钻石数
	Gold     string // 金币数
	GoldBean string // 金豆数
}

type SubUserWalletResp struct {
	Code int
	Msg  string
	Data *SubUserWalletData
}

// CheckUserDiamond

type CheckUserDiamondReq struct {
	UserId uint32 // 用户ID
}

type CheckUserDiamondResp struct {
	Code int
	Msg  string
	Data *CheckUserDiamondData
}

type CheckUserDiamondData struct {
	UserId  uint32 // 用户ID
	Diamond string // 钻石数
}

// AddUserBp

type AddUserBpReq struct {
	UserId uint32 // 用户ID
	BP     uint32
}

type AddUserBpData struct {
	UpgradeLevel uint32 // 增加后VIP等级
}

type AddUserBpResp struct {
	Code int
	Msg  string
	Data *AddUserBpData
}

// AddRecharge

type AddRechargeReq struct {
	UserId uint32 // 用户ID
	Amount string // 充值金额
}

type AddRechargeResp struct {
	Code int
	Msg  string
}

// AddOnlyOnceGoods

type AddOnlyOnceGoodsReq struct {
	UserId  uint32 // 用户ID
	GoodsId string // 商品ID
}

type AddOnlyOnceGoodsResp struct {
	Code int
	Msg  string
}
