package keys

// token
const (
	TokenFormatKey = "account:token:%v"
	UserFormatKey  = "account:user:%v"
	RpcxFormatKey  = "service:rpcx:"
	SmsFormatKey   = "sms:%s:%d"
)

// HttpRequestLimitKey API接口防刷限制Key
const (
	HttpRequestLimitKey = "limit:%s" // API接口防刷限制Key
)

// user
const (
	UserIdFormatKey          = "user:id"
	UserTelFormatKey         = "user:tel"
	UserInfoFormatKey        = "user:info:%d"
	UserWalletFormatKey      = "user:wallet:%d"
	UserWeChatFormatKey      = "user:wechat"
	UserItemFormatKey        = "user:item:%d"
	UserItemLockerFormatKey  = "user:item:lock:%d"
	UserEquipFormatKey       = "user:equip"
	UserWelfareMailFormatKey = "user:welfare-mail:%d"
	UserSystemMailFormatKey  = "user:system-mail:%d"
	UserVIPFormatKey         = "user:vip:%d"
)

// device
const (
	DeviceIdFormatKey   = "device:id"
	DeviceInfoFormatKey = "device:info:%s"
)

// config
const (
	BundleFormatKey = "bundle:info:%s"
	RoomFormatKey   = "game:room"
	GoodsFormatKey  = "goods:info"
	ItemFormatKey   = "item:info"
	MailFormatKey   = "mail:info"
)

// GMConfigurationFormatKey gm key
const (
	GMConfigurationFormatKey = "gm:config:%d-%d-%d"
)

// task key
const (
	TaskDailyLoginFormatKey = "task:login:daily:%d"
	TaskTotalLoginFormatKey = "task:login:total:%d"
	TaskDateLoginFormatKey  = "task:login:date:%d"
)

// goods key
const (
	GoodsContinueFormatKey = "goods:continue:%d"
)
