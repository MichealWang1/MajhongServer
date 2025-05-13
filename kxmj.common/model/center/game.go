package center

// GetRoomConfig

type RoomConfigReq struct {
	GameId uint16 // 游戏ID
	RoomId uint32 // 房间ID
}

type RoomConfigResp struct {
	Code int
	Msg  string
	Data *RoomConfig
}

type RoomConfig struct {
	RoomId      uint32 // 房间ID
	RoomType    uint8  // 房间类型：1 巅峰赛；2 教技场；3 比赛场；
	GameId      uint16 // 游戏ID
	GameType    uint8  // 游戏类型：1 麻将；2 斗地主
	Tags        string // 标签：1 最热；2 推荐；
	Extra       string // 扩展玩法
	RoomLevel   uint8  // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	MinLimit    string // 最小进场限制：0 代表不限制
	MaxLimit    string // 最大进场限制：0 代表不限制
	BaseScore   string // 底分
	MaxMultiple uint32 // 最大倍数
	Ticket      string // 门票
	MatchTime   uint32 // 匹配最大时长
	MatchRobot  uint8  // 是否匹配机器人：1 是；2否；
}

// GetRoomConfigList

type RoomConfigListReq struct {
	GameId uint16 // 游戏ID
}

type RoomConfigListResp struct {
	Code int
	Msg  string
	Data []*RoomConfig
}
