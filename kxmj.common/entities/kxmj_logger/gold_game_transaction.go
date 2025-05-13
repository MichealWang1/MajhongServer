package kxmj_logger

import "encoding/json"

type GoldGameTransaction struct {
	Id        int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`             // 日志ID
	GameId    uint16 `json:"game_id" redis:"game_id" gorm:"column:game_id"`          // 游戏ID
	GameType  uint8  `json:"game_type" redis:"game_type" gorm:"column:game_type"`    // 游戏类型：1 麻将；2 斗地主
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`          // 用户ID
	RoomId    uint32 `json:"room_id" redis:"room_id" gorm:"column:room_id"`          // 房间ID
	RoomLevel uint8  `json:"room_level" redis:"room_level" gorm:"column:room_level"` // 房间等级
	Gold      string `json:"gold" redis:"gold" gorm:"column:gold"`                   // 金币
	Type      uint8  `json:"type" redis:"type" gorm:"column:type"`                   // 交易类型：1 带入游戏；2 游戏带出；
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"` // 创建时间
}

func (g *GoldGameTransaction) TableName() string {
	return "gold_game_transaction"
}

func (g *GoldGameTransaction) Schema() string {
	return "kxmj_logger"
}

func (g *GoldGameTransaction) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *GoldGameTransaction) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, g)
}
