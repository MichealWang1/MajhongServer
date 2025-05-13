package kxmj_game

import "encoding/json"

type UserGameRecord struct {
	Id        int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`             // 日志ID
	RoundId   string `json:"round_id" redis:"round_id" gorm:"column:round_id"`       // 局号
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`          // 用户ID
	GameId    uint16 `json:"game_id" redis:"game_id" gorm:"column:game_id"`          // 游戏Id
	RoomId    uint8  `json:"room_id" redis:"room_id" gorm:"column:room_id"`          // 房间ID
	RoomLevel uint8  `json:"room_level" redis:"room_level" gorm:"column:room_level"` // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	TableId   uint32 `json:"table_id" redis:"table_id" gorm:"column:table_id"`       // 桌子ID
	WinLose   int64  `json:"win_lose" redis:"win_lose" gorm:"column:win_lose"`       // 玩家输赢
	Duration  uint32 `json:"duration" redis:"duration" gorm:"column:duration"`       // 游戏时长
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"` // 创建时间
	SettledAt uint32 `json:"settled_at" redis:"settled_at" gorm:"column:settled_at"` // 结算时间
}

func (u *UserGameRecord) TableName() string {
	return "user_game_record"
}

func (u *UserGameRecord) Schema() string {
	return "kxmj_game"
}

func (u *UserGameRecord) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserGameRecord) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
