package kxmj_report

import "encoding/json"

type ReportUserGameDaily struct {
	Id            uint32 `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"`          // ID
	UserId        uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                      // 用户ID
	GameId        uint16 `json:"game_id" redis:"game_id" gorm:"column:game_id"`                      // 游戏ID
	GameType      uint8  `json:"game_type" redis:"game_type" gorm:"column:game_type"`                // 游戏类型：1 麻将；2 斗地主
	RoomLevel     uint8  `json:"room_level" redis:"room_level" gorm:"column:room_level"`             // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	PlayType      uint8  `json:"play_type" redis:"play_type" gorm:"column:play_type"`                // 对局类型：1 PVE；2 PVP；3 混合；
	TotalTimes    uint32 `json:"total_times" redis:"total_times" gorm:"column:total_times"`          // 总局数
	TotalWinLoss  string `json:"total_win_loss" redis:"total_win_loss" gorm:"column:total_win_loss"` // 总输赢
	TotalDuration uint32 `json:"total_duration" redis:"total_duration" gorm:"column:total_duration"` // 总时长
	CreatedAt     uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`             // 创建时间
	UpdatedAt     uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`             // 更新时间
	Date          uint32 `json:"date" redis:"date" gorm:"column:date"`                               // 日期
}

func (r *ReportUserGameDaily) TableName() string {
	return "report_user_game_daily"
}

func (r *ReportUserGameDaily) Schema() string {
	return "kxmj_report"
}

func (r *ReportUserGameDaily) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReportUserGameDaily) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}
