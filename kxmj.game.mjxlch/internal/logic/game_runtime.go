package logic

type GameRunTimeData struct {
	RuleConfig     *GameRuleData        // 游戏规则
	DiceData       *DiceStateData       // 掷骰子数据(庄家、精牌...)
	DealData       *DealStateData       // 发牌数据
	SwapData       *SwapStateData       // 换牌数据
	ChooseMissData *ChooseMissStateData // 选缺数据
	PlayData       *PlayStateData       // 出牌数据
}

func NewDeskRunTimeData() *GameRunTimeData {
	g := &GameRunTimeData{
		RuleConfig:     NewGameRuleData(),
		DiceData:       NewDiceStateData(),
		DealData:       NewDealStateData(),
		SwapData:       NewSwapStateData(),
		ChooseMissData: NewChooseMissStateData(),
		PlayData:       NewPlayStateData(),
	}
	return g
}

func (g *GameRunTimeData) Reset() {
	g.RuleConfig.Reset()
	g.DiceData.Reset()
	g.DealData.Reset()
	g.SwapData.Reset()
	g.ChooseMissData.Reset()
	g.PlayData.Reset()
}
