package game

import (
	lib "kxmj.common/mahjong"
	"kxmj.common/utils"
	"kxmj.game.mjxlch/pb"
	"time"
)

type Player struct {
	SeatId     uint32 // 座位号
	UserId     uint32 // 用户ID
	Gold       string // 金豆
	Nickname   string // 昵称
	AvatarAddr string // 头像
	IsOnline   bool   // 是否在线
	IsRobot    bool   // 是否机器人
	GoldStatus bool   // 金币是否转入游戏

	RunTimeData *PlayerRunTimeData // 玩家游戏中的数据
}

func (p *Player) Reset() {
	p.UserId = 0
	p.Gold = ""
	p.Nickname = ""
	p.AvatarAddr = ""
	p.IsOnline = false
	p.IsRobot = false
	p.GoldStatus = false
	p.RunTimeData = NewPlayerRunTimeData()
}

type PlayerRunTimeData struct {
	isReady bool // 是否准备

	handCards        lib.Cards       // 手牌
	discards         lib.Cards       // 弃牌
	actionCardsTable lib.UserActions // 动作表
	catchCard        lib.Card        // 摸的牌(不是自己摸牌则为0xff)
	outCount         int32           // 出牌次数
	huData           []*HuData       // 胡牌数据
	totalScore       string          // 玩家总得分

	operationalActions lib.UserActions // 可操作动作（进行吃碰杠胡选择的动作）
	respOperateAction  *lib.UserAction // 玩家响应的动作
	currentHuResult    *HuData         // 玩家可胡牌的result(选择胡牌后暂时存储)

	swapCards      lib.Cards   // 换三张牌
	chooseMissType pb.MissType // 选缺类型

	hostState bool  // 托管状态
	opTime    int64 // 当前操作时间
}

func NewPlayerRunTimeData() *PlayerRunTimeData {
	return &PlayerRunTimeData{
		isReady:            false,
		handCards:          make(lib.Cards, 0, 14),
		discards:           make(lib.Cards, 0),
		actionCardsTable:   make(lib.UserActions, 0),
		catchCard:          lib.INVALID_CARD,
		outCount:           0,
		huData:             make([]*HuData, 0),
		totalScore:         "0",
		operationalActions: make(lib.UserActions, 0),
		respOperateAction:  nil,
		currentHuResult:    nil,
		swapCards:          make(lib.Cards, 0, 3),
		chooseMissType:     pb.MissType_MISS_NULL,
	}
}

// --------------------准备-------------------------
func (p *Player) setPlayerReady() {
	p.RunTimeData.isReady = true
}

func (p *Player) getPlayerReady() bool {
	return p.RunTimeData.isReady
}

// ----------------------手牌------------------------

// SetHandCards 设置玩家手牌
func (p *Player) setHandCards(cards lib.Cards) {
	p.RunTimeData.handCards = cards
}

// AddHandCards 添加玩家手牌
func (p *Player) addHandCards(card lib.Card) {
	newCards := p.RunTimeData.handCards.AddCard(card)
	p.RunTimeData.handCards = newCards
}

// DeleteHandCards 删除玩家手牌
func (p *Player) deleteHandCard(card lib.Card) error {
	newCards, err := p.RunTimeData.handCards.DeleteCard(card)
	if err != nil {
		return err
	}
	p.RunTimeData.handCards = newCards
	return nil
}

// DeleteSomeHandCards 删除一些牌
func (p *Player) deleteHandCards(cards lib.Cards) error {
	newCards, err := p.RunTimeData.handCards.DeleteCards(cards)
	if err != nil {
		return err
	}
	p.RunTimeData.handCards = newCards
	return nil
}

// GetHandCards 获取玩家手牌
func (p *Player) getHandCards() lib.Cards {
	return p.RunTimeData.handCards
}

// ----------------------------出牌次数-------------------------

// 玩家出牌次数自增
func (p *Player) addOutCount() {
	p.RunTimeData.outCount++
}

// 获取玩家出牌次数
func (p *Player) getOutCount() int32 {
	return p.RunTimeData.outCount
}

// -------------------------弃牌------------------------------

// AddDiscards 增加弃牌
func (p *Player) addDiscards(card lib.Card) {
	newCards := p.RunTimeData.discards.AddCard(card)
	p.RunTimeData.discards = newCards
}

// DeleteDiscards 删除弃牌
func (p *Player) deleteDiscards(card lib.Card) error {
	newCards, err := p.RunTimeData.discards.DeleteCard(card)
	if err != nil {
		return err
	}
	p.RunTimeData.discards = newCards
	return nil
}

// GetDiscards 获取玩家弃牌
func (p *Player) getDiscards() lib.Cards {
	return p.RunTimeData.discards
}

// --------------------------动作牌------------------------

// AddActionCards 增加玩家动作
func (p *Player) addActionCards(action *lib.UserAction) {
	p.RunTimeData.actionCardsTable.Add(action.Copy())
}

// ActionBuGang 玩家进行补杠操作
func (p *Player) actionBuGang(opCard lib.Card) {
	for _, action := range p.RunTimeData.actionCardsTable {
		if action.ActionType == lib.ActionType_Peng && action.OutCard == opCard {
			action.ActionType = lib.ActionType_Gang
			action.ExtraActionType = lib.ExtraActionType_Bu_Gang
			action.CombineCards = opCard.Repeat(4)
		}
	}
}

// GetActionCardsTable 获取玩家动作表
func (p *Player) getActionCardsTable() lib.UserActions {
	return p.RunTimeData.actionCardsTable
}

// ---------------------------摸牌----------------------------

// SetCatchCard 设置玩家摸牌
func (p *Player) setCatchCard(card lib.Card) {
	p.RunTimeData.catchCard = card
}

// ResetCatchCard 重置玩家摸牌
func (p *Player) resetCatchCard() {
	p.RunTimeData.catchCard = lib.INVALID_CARD
}

// GetCatchCard 获取玩家的摸牌
func (p *Player) getCatchCard() lib.Card {
	return p.RunTimeData.catchCard
}

// --------------------------换三张------------------------------

// SetSwapCards 设置玩家换三张牌
func (p *Player) setSwapCards(cards lib.Cards) {
	p.RunTimeData.swapCards = cards
}

// GetSwapCards 获取玩家换三张牌
func (p *Player) getSwapCards() lib.Cards {
	return p.RunTimeData.swapCards
}

// --------------------------选缺--------------------------------

// SetChooseMissType 设置玩家选缺类型
func (p *Player) setChooseMissType(t pb.MissType) {
	p.RunTimeData.chooseMissType = t
}

// GetChooseMissType 获取玩家选缺类型
func (p *Player) getChooseMissType() pb.MissType {
	return p.RunTimeData.chooseMissType
}

// -------------------------托管----------------------------------

// 设置玩家托管状态
func (p *Player) setHostState(state bool) {
	p.RunTimeData.hostState = state
}

// 获取玩家托管状态
func (p *Player) getHostState() bool {
	return p.RunTimeData.hostState
}

// ---------------------------操作动作------------------------

// 设置玩家可操作的动作
func (p *Player) setOperationalActions(actions lib.UserActions) {
	p.RunTimeData.opTime = time.Now().UnixMilli()
	p.RunTimeData.operationalActions = actions
}

// 获取玩家可操作动作
func (p *Player) getOperationalActions() lib.UserActions {
	return p.RunTimeData.operationalActions
}

func (p *Player) resetOperationalActions() {
	p.RunTimeData.operationalActions = make(lib.UserActions, 0)
}

// 设置玩家响应动作
func (p *Player) setRespOperateAction(action *lib.UserAction) {
	p.RunTimeData.respOperateAction = action
}

// 获取玩家响应动作
func (p *Player) getRespOperateAction() *lib.UserAction {
	return p.RunTimeData.respOperateAction
}

// 重置响应
func (p *Player) resetRespOperateAction() {
	p.RunTimeData.respOperateAction = nil
}

// ---------------------------胡牌数据--------------------------------

// 设置当前可操作的胡result
func (p *Player) setCurrentHuResult(d *HuData) {
	p.RunTimeData.currentHuResult = d
}

// 获取当前可胡result
func (p *Player) getCurrentHuResult() *HuData {
	return p.RunTimeData.currentHuResult
}

// 重置当前可胡result
func (p *Player) resetCurrentHuResult() {
	p.RunTimeData.currentHuResult = nil
}

// 增加玩家胡牌数据
func (p *Player) addHuData(data *HuData) {
	p.RunTimeData.huData = append(p.RunTimeData.huData, data)
}

// 获取玩家胡牌数据
func (p *Player) getHuData() []*HuData {
	return p.RunTimeData.huData
}

// 玩家是否胡牌
func (p *Player) isHu() bool {
	return len(p.RunTimeData.huData) > 0
}

// -----------------------操作时间---------------------------
func (p *Player) setOpTime(time int64) {
	p.RunTimeData.opTime = time
}

func (p *Player) getOpTime() int64 {
	return p.RunTimeData.opTime
}

// 托管玩家自动操作
func (p *Player) canAutoOperation() bool {
	nowTime := time.Now().UnixMilli()
	if (p.getHostState() || p.IsRobot) && nowTime > p.getOpTime()+DurationAuto {
		return true
	}
	return false
}

// --------------------结算------------------------
func (p *Player) addTotalScore(score string) {
	p.RunTimeData.totalScore, _ = utils.AddToString(p.RunTimeData.totalScore, score)
}

func (p *Player) getTotalScore() string {
	return p.RunTimeData.totalScore
}
