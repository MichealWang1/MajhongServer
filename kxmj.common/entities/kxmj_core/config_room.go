package kxmj_core

import "encoding/json"

type ConfigRoom struct {
	RoomId      uint32 `json:"room_id" redis:"room_id" gorm:"column:room_id;primary_key"`    // 房间ID
	RoomType    uint8  `json:"room_type" redis:"room_type" gorm:"column:room_type"`          // 房间类型：1:普通场； 2 巅峰赛；3 教技场；4 比赛场；
	RoomLevel   uint8  `json:"room_level" redis:"room_level" gorm:"column:room_level"`       // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	GameId      uint16 `json:"game_id" redis:"game_id" gorm:"column:game_id"`                // 游戏ID
	GameName    string `json:"game_name" redis:"game_name" gorm:"column:game_name"`          // 游戏名称
	GameType    uint8  `json:"game_type" redis:"game_type" gorm:"column:game_type"`          // 游戏类型：1 麻将；2 斗地主
	BaseScore   string `json:"base_score" redis:"base_score" gorm:"column:base_score"`       // 底分
	Ticket      string `json:"ticket" redis:"ticket" gorm:"column:ticket"`                   // 门票
	MinLimit    string `json:"min_limit" redis:"min_limit" gorm:"column:min_limit"`          // 最小进场限制：0 代表不限制
	MaxLimit    string `json:"max_limit" redis:"max_limit" gorm:"column:max_limit"`          // 最大进场限制：0 代表不限制
	MaxMultiple uint32 `json:"max_multiple" redis:"max_multiple" gorm:"column:max_multiple"` // 封顶倍数
	MatchTime   uint32 `json:"match_time" redis:"match_time" gorm:"column:match_time"`       // 匹配最大时长(秒)
	Tags        string `json:"tags" redis:"tags" gorm:"column:tags"`                         // 标签：1 最热；2 推荐；
	Extra       string `json:"extra" redis:"extra" gorm:"column:extra"`                      // 扩展玩法
	MatchRobot  uint8  `json:"match_robot" redis:"match_robot" gorm:"column:match_robot"`    // 是否匹配机器人：1 是；2否；
	CreatedAt   uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`       // 创建时间
	UpdatedAt   uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`       // 更新时间
}

func (c *ConfigRoom) TableName() string {
	return "config_room"
}

func (c *ConfigRoom) Schema() string {
	return "kxmj_core"
}

func (c *ConfigRoom) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ConfigRoom) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
