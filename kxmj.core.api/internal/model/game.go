package model

// RoomListReq 获取房间列表
type RoomListReq struct {
	GameId uint16 `json:"gameId" binding:"required"` // 游戏ID (游戏类型)
}

// RoomInfo 房间信息
type RoomInfo struct {
	RoomId      uint32 `json:"roomId"`      // 房间ID
	RoomType    uint8  `json:"roomType"`    // 房间类型：1 巅峰赛；2 教技场；3 比赛场；
	GameId      uint16 `json:"gameId"`      // 游戏ID
	GameType    uint8  `json:"gameType"`    // 游戏类型：1 麻将；2 斗地主
	RoomLevel   uint8  `json:"roomLevel"`   // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	Tags        string `json:"tags"`        // 标签：1 最热；2 推荐；
	Extra       string `json:"extra"`       // 扩展玩法
	MinLimit    string `json:"minLimit"`    // 最小进场限制：0 代表不限制
	MaxLimit    string `json:"maxLimit"`    // 最大进场限制：0 代表不限制
	BaseScore   string `json:"baseScore"`   // 底分
	MaxMultiple uint32 `json:"maxMultiple"` // 最大倍数
	Ticket      string `json:"ticket"`      // 门票
	MatchTime   uint32 `json:"matchTime"`   // 匹配最大时长(秒)
	CurPlayers  uint32 `json:"curPlayers"`  // 当前房间玩家数量
}

// RoomListResp 获取房间列表
type RoomListResp struct {
	List []*RoomInfo `json:"list"` // 房间信息列表
}
