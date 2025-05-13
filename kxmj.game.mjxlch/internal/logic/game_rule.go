package logic

type GameRuleData struct {
	PlayerNums int32 // 游戏人数
}

func NewGameRuleData() *GameRuleData {
	r := &GameRuleData{}
	r.Reset()
	return r
}

func (r *GameRuleData) Reset() {
	r.PlayerNums = 4
}

func (r *GameRuleData) GetPlayerNums() int32 {
	return r.PlayerNums
}
