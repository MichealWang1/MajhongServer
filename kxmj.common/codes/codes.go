package codes

const (
	UnKnowError              = 0  // 未知错误
	Success                  = 1  // 成功
	UserNameNotNull          = 2  // 用户名不能为空
	PasswordNotNull          = 3  // 密码不能为空
	UserNotExist             = 4  // 用户不存在
	PasswordError            = 5  // 密码错误
	DbError                  = 6  // 数据库错误
	AuthorizeFailed          = 7  // 授权失败
	InsufficientGold         = 8  // 金币不足
	FailedOperate            = 9  // 操作失败
	ServerNetErr             = 10 // 网络错误
	ServerNotFound           = 11 // 服务未找到
	ParamError               = 12 // 参数错误
	BundleIdNotNull          = 13 // 包名不能为空
	MarshalJsonErr           = 14 // Json序列化错误
	UnMarshalJsonErr         = 15 // Json反序列化错误
	MarshalPbErr             = 16 // Pb序列化错误
	UnMarshalPbErr           = 17 // Pb反序列化错误
	RuntimeErr               = 18 // Pb反序列化错误
	BundleIdErr              = 19 // 分包ID错误
	UserIsFrozen             = 20 // 用户被冻结
	UserExisted              = 21 // 用户已存在
	GameDeskNotExist         = 22 // 游戏桌子不存在
	GameServerClose          = 23 // 游戏服务已关闭
	GetUserInfoFail          = 24 // 获取用户信息失败
	SvrMatchFail             = 25 // 服务匹配失败
	RoomMatchFail            = 26 // 房间匹配失败
	GetUserGoldFail          = 27 // 获取用户金豆失败
	DeviceNotFound           = 28 // 设备信息未找到
	CreateUserIdFailed       = 29 // 生成用户ID失败
	BundleNotNull            = 30 // 分包信息未找到
	CreateUserFailed         = 31 // 创建用户失败
	SendSmsFailed            = 32 // 发送短信验证码失败
	InvalidTelNumber         = 33 // 非法手机号
	CheckSmsCodeFailed       = 34 // 短信验证码校验失败
	InvalidSmsCode           = 35 // 短信验证码不正确
	TelNumberExisted         = 36 // 手机号码已存在
	CanNotRepeatSendSms      = 37 // 重复发送
	GetGoodsConfigFailed     = 38 // 获取商品配置失败
	GetItemConfigFailed      = 39 // 获取物品配置失败
	GetUserInfoFailed        = 40 // 获取用户信息失败
	GetRoomConfigFailed      = 41 // 获取房间配置失败
	GetWalletInfoFailed      = 42 // 获取钱包信息失败
	UpdateWalletInfoFailed   = 43 // 更新钱包信息失败
	CheckWalletInfoFailed    = 44 // 检查钱包信息失败
	GoodsIsDeleteFailed      = 45 // 商品已经删除
	GoodsIsOffShelf          = 46 // 商品已下架
	PayProviderNotExist      = 47 // 支付渠道不存在
	ThirdCallFailed          = 48 // 调用第三方支付接口失败
	CreateGoodsOrderFailed   = 49 // 创建商品订单失败
	ParseBaseValueItemFailed = 50 // 解析物品数据失败
	GoodsShopTypeError       = 51 // 商品销售类型错误
	AddUserItemFailed        = 52 // 增加用户背包物品失败
	InsufficientDiamond      = 53 // 没有足够的钻石
	InsufficientGoldBean     = 54 // 没有足够的金豆
	AddUserBpFailed          = 55 // 增加用户经验值失败
	InvalidItem              = 56 // 无效物品
	AddUserWalletFailed      = 57 // 增加用户钱包失败
	AddUserRechargeFailed    = 58 // 增加用户累计充值失败
	GoodsOnlyBuyOnce         = 59 // 该商品只能购买一次
	NotCanTakeItems          = 60 // 没有可领取物品
	TheItemAlreadyTake       = 61 // 该物品已经被领取
)

var message = map[int]string{
	UnKnowError:              "未知错误",
	Success:                  "成功",
	UserNameNotNull:          "用户名不能为空",
	PasswordNotNull:          "密码不能为空",
	UserNotExist:             "用户不存在",
	PasswordError:            "密码错误",
	DbError:                  "数据库错误",
	AuthorizeFailed:          "授权失败",
	InsufficientGold:         "金币不足",
	FailedOperate:            "操作失败",
	ServerNetErr:             "网络错误",
	ServerNotFound:           "服务未找到",
	ParamError:               "参数错误",
	BundleIdNotNull:          "包名不能为空",
	MarshalJsonErr:           "Json序列化错误",
	UnMarshalJsonErr:         "Json反序列化错误",
	MarshalPbErr:             "Pb序列化错误",
	UnMarshalPbErr:           "Pb反序列化错误",
	RuntimeErr:               "Pb反序列化错误",
	BundleIdErr:              "分包ID错误",
	UserIsFrozen:             "用户被冻结",
	UserExisted:              "用户已存在",
	GameDeskNotExist:         "游戏桌子不存在",
	GameServerClose:          "游戏服务已关闭",
	GetUserInfoFail:          "获取用户信息失败",
	SvrMatchFail:             "服务匹配失败",
	RoomMatchFail:            "房间匹配失败",
	GetUserGoldFail:          "获取用户金豆失败",
	DeviceNotFound:           "设备信息未找到",
	CreateUserIdFailed:       "生成用户ID失败",
	BundleNotNull:            "分包信息未找到",
	CreateUserFailed:         "创建用户失败",
	SendSmsFailed:            "发送短信验证码失败",
	InvalidTelNumber:         "非法手机号",
	CheckSmsCodeFailed:       "短信验证码校验失败",
	InvalidSmsCode:           "短信验证码不正确",
	TelNumberExisted:         "手机号码已存在",
	CanNotRepeatSendSms:      "重复发送",
	GetGoodsConfigFailed:     "获取商品配置失败",
	GetItemConfigFailed:      "获取物品配置失败",
	GetUserInfoFailed:        "获取用户信息失败",
	GetRoomConfigFailed:      "获取房间配置失败",
	GetWalletInfoFailed:      "获取钱包信息失败",
	UpdateWalletInfoFailed:   "更新钱包信息失败",
	CheckWalletInfoFailed:    "检查钱包信息失败",
	GoodsIsDeleteFailed:      "商品已经删除",
	GoodsIsOffShelf:          "商品已下架",
	PayProviderNotExist:      "支付渠道不存在",
	ThirdCallFailed:          "调用第三方支付接口失败",
	CreateGoodsOrderFailed:   "创建商品订单失败",
	ParseBaseValueItemFailed: "解析物品数据失败",
	GoodsShopTypeError:       "商品销售类型错误",
	AddUserItemFailed:        "增加用户背包物品失败",
	InsufficientDiamond:      "没有足够的钻石",
	InsufficientGoldBean:     "没有足够的金豆",
	AddUserBpFailed:          "增加用户经验值失败",
	InvalidItem:              "无效物品",
	AddUserWalletFailed:      "增加用户钱包失败",
	AddUserRechargeFailed:    "增加用户累计充值失败",
	GoodsOnlyBuyOnce:         "该商品只能购买一次",
	NotCanTakeItems:          "没有可领取物品",
	TheItemAlreadyTake:       "该物品已经被领取",
}

type IResultCode interface {
	InsertMessages(msg map[int]string)
	GetMessage(code int) string
}

func AddMessages(messages map[int]string) {
	for k, v := range messages {
		message[k] = v
	}
}

func AddAllMessages() {
	AddMessages(coreApiMessage)
	AddMessages(emailMessage)
	AddMessages(gameMessage)
	AddMessages(lobbyMessage)
	AddMessages(shopMessage)
	AddMessages(taskMessage)
}

func AllMessages() map[int]string {
	return message
}

func GetMessage(code int) string {
	msg, ok := message[code]
	if ok {
		return msg
	}
	return "unknown error"
}
