package game_core

import "kxmj.common/model/center"

type Room struct {
	roomId      uint32            // 房间ID
	gameId      uint16            // 游戏ID
	gameType    uint8             // 游戏类型：1 麻将；2 斗地主
	roomLevel   uint8             // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	minLimit    string            // 最小进场限制：0 代表不限制
	maxLimit    string            // 最大进场限制：0 代表不限制
	baseScore   string            // 底分
	maxMultiple uint32            // 最大倍数
	ticket      string            // 门票
	matchTime   uint32            // 匹配最大时间 (秒)
	matchRobot  uint8             // 是否匹配机器人：1 是；2否；
	lastDeskId  uint32            // 最后一次分配桌子ID
	close       bool              // 关闭状态
	deskIds     map[uint32]uint32 // 当前已分配桌子ID
}

func NewRoom(config *center.RoomConfig) *Room {
	return &Room{
		roomId:      config.RoomId,
		gameId:      config.GameId,
		gameType:    config.GameType,
		roomLevel:   config.RoomLevel,
		minLimit:    config.MinLimit,
		maxLimit:    config.MaxLimit,
		baseScore:   config.BaseScore,
		maxMultiple: config.MaxMultiple,
		ticket:      config.Ticket,
		matchTime:   config.MatchTime,
		matchRobot:  config.MatchRobot,
		lastDeskId:  0,
		deskIds:     make(map[uint32]uint32, 0),
	}
}

func (r *Room) Update(config *center.RoomConfig) {
	r.gameId = config.GameId
	r.gameType = config.GameType
	r.roomLevel = config.RoomLevel
	r.minLimit = config.MinLimit
	r.maxLimit = config.MaxLimit
	r.baseScore = config.BaseScore
	r.maxMultiple = config.MaxMultiple
	r.ticket = config.Ticket
	r.matchTime = config.MatchTime
	r.matchRobot = config.MatchRobot
	r.lastDeskId = 0
}

func (r *Room) Close() {
	r.close = true
}

func (r *Room) ID() uint32 {
	return r.roomId
}

func (r *Room) GameId() uint16 {
	return r.gameId
}

func (r *Room) GameType() uint8 {
	return r.gameType
}

func (r *Room) RoomLevel() uint8 {
	return r.roomLevel
}

func (r *Room) MinLimit() string {
	return r.minLimit
}

func (r *Room) MaxLimit() string {
	return r.maxLimit
}

func (r *Room) BaseScore() string {
	return r.baseScore
}

func (r *Room) MaxMultiple() uint32 {
	return r.maxMultiple
}

func (r *Room) Ticket() string {
	return r.ticket
}

func (r *Room) MatchTime() uint32 {
	return r.matchTime
}

func (r *Room) MatchRobot() bool {
	return r.matchRobot == 1
}

func (r *Room) GenerateDeskId() uint32 {
	// 一个房间最多分配10万张桌子
	if r.lastDeskId >= 100000 {
		r.lastDeskId = 0
	}
	r.lastDeskId++

	for {
		_, has := r.deskIds[r.lastDeskId]
		if has {
			r.lastDeskId++
			continue
		}
		r.deskIds[r.lastDeskId] = r.lastDeskId
		break
	}
	return r.lastDeskId
}

func (r *Room) RemoveDeskId(deskId uint32) {
	delete(r.deskIds, deskId)
}
