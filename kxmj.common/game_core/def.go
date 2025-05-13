package game_core

type LeaveDesk struct {
	UserId     uint32            // 用户ID
	Gold       string            // 金币
	SeatId     uint32            // 座位号
	IsRobot    bool              // 是否是机器人
	Desk       IDesk             // 桌子
	Statistics []*GameStatistics // 对局统计数据
}

type PlayerInfo struct {
	UserId     uint32 // 用户ID
	IsRobot    bool   // 是否是机器人
	SeatId     uint32 // 座位号
	Gold       string // 金豆
	Nickname   string // 昵称
	AvatarAddr string // 头像
	IconStyle  uint32 // 装饰
}

type EnterDesk struct {
	Desk   IDesk       // 桌子
	Player *PlayerInfo // 玩家信息
}

type GameStatistics struct {
	PlayType      uint8  // 对局类型：1 PVE；2 PVP；3 混合；
	TotalTimes    uint32 // 总局数
	TotalWinLoss  string // 总输赢
	TotalDuration uint32 // 总时长
}

type UserStatistics struct {
	UserId        uint32 // 用户ID
	RoomId        uint32 // 房间ID
	GameId        uint16 // 游戏ID
	GameType      uint8  // 游戏类型：1 麻将；2 斗地主
	RoomLevel     uint8  // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	PlayType      uint8  // 对局类型：1 PVE；2 PVP；3 混合；
	TotalTimes    uint32 // 总局数
	TotalWinLoss  string // 总输赢
	TotalDuration uint32 // 总时长
}

type RiseGoods struct {
	GoodsId       string // 商品Id
	Price         string // 商品价格
	RealCount     string // 实际获得数量
	OriginalCount string // 原来获得数量
	ShopType      uint8  // 销售类型：1 RMB购买；2 钻石购买；3 金币购买；4 金豆购买；
	RiseLevel     uint32 // 复活卡等级：1，2，3 级
}

type RiseGoodsInfo struct {
	List []*RiseGoods // 购买列表
}
