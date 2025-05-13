package service

const (
	GatewayService = uint16(1000)  // 网关服务
	CenterService  = uint16(2000)  // 中心服务
	TaskService    = uint16(3000)  // 任务服务
	ApiCoreService = uint16(4000)  // ApiCore服务
	LobbyService   = uint16(5000)  // 大厅服务
	PaymentService = uint16(6000)  // 支付服务
	ReportService  = uint16(7000)  // 业务报表服务
	SmsService     = uint16(8000)  // 短信服务
	ShopService    = uint16(9000)  // 商城服务
	EmailService   = uint16(10000) // 邮件服务
	GMService      = uint16(11000) // GM服务
)

const (
	GameMjXLCH = uint16(50001) // 麻将血流成河
	GameMjXZ   = uint16(50002) // 麻将血战玩法
	GameMjHZXL = uint16(50003) // 麻将红中血流
	GameMjDZXL = uint16(50004) // 麻将大众血流
)
