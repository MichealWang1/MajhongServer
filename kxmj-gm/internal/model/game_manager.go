package model

type CardStack struct {
	UserId    uint32   `json:"userId"`   // 玩家ID
	GameType  uint16   `json:"gameType"` // 游戏场次
	RoomLevel uint8    `json:"roomType"` // 房间类型
	Banker    uint8    `json:"banker"`   // 庄家 自己则是UserID对应的玩家 1自己 2下家 3对家 4上家
	Cards     []uint32 `json:"cards"`    // 配牌库 万(1-9) 条(17-25) 筒(33-41)
}

type CatchCard struct {
	UserId    uint32 `json:"userId"`   // 玩家ID
	GameType  uint16 `json:"gameType"` // 游戏场次
	RoomLevel uint8  `json:"roomType"` // 房间类型
	Card      uint32 `json:"card"`     // 牌 万(1-9) 条(17-25) 筒(33-41)
}

type DeleteCardStack struct {
	UserId   uint32 `json:"userId"`   // 玩家ID
	GameType uint16 `json:"gameType"` // 游戏场次
	RoomType uint8  `json:"roomType"` // 房间类型
}

type UserInfo struct {
	UserId uint32 `json:"userId"` // 玩家ID
}

type MatchPlayerType struct {
	UserId    uint32  `json:"userId"`    // 玩家ID
	GameType  uint16  `json:"gameType"`  // 游戏场次
	RoomLevel uint8   `json:"roomType"`  // 房间类型
	MatchType []uint8 `json:"matchType"` // 匹配类型 没选择的或者默认的填0即可
	MatchTime uint32  `json:"matchTime"` // 匹配时间
}
