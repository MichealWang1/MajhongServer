package game

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"kxmj.common/game_core"
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/net"
	"kxmj.game.mjxlch/pb"
	"math/rand"
	"time"
)

type Desk struct {
	id        uint32              // 桌子ID
	roundId   string              // 局号
	game      game_core.IGame     // 游戏管理类实例
	room      game_core.IRoom     // 房间
	players   []*Player           // 玩家列表
	msgChan   chan net.MsgContext // 消息管道
	runChan   chan struct{}       // 状态机管道
	markTime  int64               // 状态机执行时间
	nextTime  int64               // 下一个游戏状态开始时间
	status    Status              // 游戏状态
	closeChan chan struct{}       // 关闭管道
	close     bool                // 关闭标志

	waiting *WaitingEvent // 等待切换下一个状态

	runtimeData *DeskRunTimeData // 桌子游戏数据
}

func NewTemplate() *Desk {
	return &Desk{}
}

func (d *Desk) New(game game_core.IGame, room game_core.IRoom) game_core.IDesk {
	desk := &Desk{
		id:        room.GenerateDeskId(),
		game:      game,
		room:      room,
		players:   make([]*Player, PLAYER_COUNT),
		msgChan:   make(chan net.MsgContext, 100),
		runChan:   make(chan struct{}, 100),
		markTime:  time.Now().UnixMilli(),
		nextTime:  time.Now().UnixMilli(),
		status:    Match,
		closeChan: make(chan struct{}, 0),
		close:     false,
		waiting:   NewWaitingEvent(Match, MatchWaitDuration, true),
	}
	desk.runtimeData = NewDeskRunTimeData()

	// 初始化座位
	for i := uint32(0); i < PLAYER_COUNT; i++ {
		desk.players[i] = &Player{
			SeatId:      i,
			RunTimeData: NewPlayerRunTimeData(),
		}
	}
	// 初始化庄家和操作位子
	desk.randomBankerSeatId()
	desk.setOperateSeatId(desk.runtimeData.bankerSeatId)

	desk.roundId = desk.newRoundId()
	return desk
}

// 玩家离开游戏
func (d *Desk) PlayerOnLeave(seatId uint32) {
	player := d.getPlayerBySeat(seatId)
	gold := "0"
	if player.GoldStatus {
		gold = player.Gold
	}
	log.Sugar().Infof("Player:%v OnLeave gold:%v", player, gold)
	// 通知玩家离开
	d.game.OnLeave(&game_core.LeaveDesk{
		UserId:  player.UserId,
		Gold:    gold,
		SeatId:  player.SeatId,
		IsRobot: player.IsRobot,
		Desk:    d,
	})
}

func (d *Desk) addRobot() {
	robot := d.game.GetRobot(d.room.BaseScore())
	player := d.getNullPlayerBySeat()
	player.UserId = robot.UserId
	player.Gold = robot.Gold
	player.Nickname = robot.Nickname
	player.AvatarAddr = robot.AvatarAddr
	player.IsRobot = true
	player.IsOnline = true
	player.RunTimeData = NewPlayerRunTimeData()
	player.setPlayerReady()
	// 通知大厅用户进入游戏
	d.game.OnEnter(&game_core.EnterDesk{
		Desk: d,
		Player: &game_core.PlayerInfo{
			UserId:     player.UserId,
			IsRobot:    player.IsRobot,
			SeatId:     player.SeatId,
			Gold:       player.Gold,
			Nickname:   player.Nickname,
			AvatarAddr: player.AvatarAddr,
			IconStyle:  0,
		},
	})
	d.broadcastEnterInfoNotify(player.UserId)
}

func (d *Desk) ID() uint32 {
	return d.id
}

func (d *Desk) RoundId() string {
	return d.roundId
}

func (d *Desk) Room() game_core.IRoom {
	return d.room
}

func (d *Desk) Start() {
	go func() {
		for {
			select {
			case <-d.closeChan:
				d.game.OnDeskClose(d)
				return
			case ctx := <-d.msgChan:
				d.handler(ctx)
			case <-d.runChan:
				d.run()
			}
		}
	}()
}

func (d *Desk) Close() {
	d.close = true
}

func (d *Desk) Run() {
	d.runChan <- struct{}{}
}

func (d *Desk) OnMessage(ctx net.MsgContext) {
	d.msgChan <- ctx
}

func (d *Desk) newRoundId() string {
	return fmt.Sprintf("%d%d%d%d%d", time.Now().Unix(), d.game.Server().SvrType(), d.game.Server().SvrId(), d.room.ID(), d.id)
}

func (d *Desk) getPlayer(userId uint32) *Player {
	var p *Player
	for _, v := range d.players {
		if v.UserId == userId {
			p = v
			break
		}
	}
	return p
}

// 随机获取一个空player
func (d *Desk) getNullPlayerBySeat() *Player {
	seatId := rand.Uint32() % PLAYER_COUNT
	for i := uint32(0); i < PLAYER_COUNT; i++ {
		seatId = (seatId + i) % PLAYER_COUNT
		p := d.getPlayerBySeat(seatId)
		if p.UserId == 0 {
			break
		}
	}
	log.Sugar().Infof("null seat:%d ", seatId)
	return d.getPlayerBySeat(seatId)
}

func (d *Desk) getPlayerBySeat(seatId uint32) *Player {
	var p *Player
	for _, v := range d.players {
		if v.SeatId == seatId {
			p = v
			break
		}
	}
	return p
}

// 获取当前位子的上一个位子
func (d *Desk) getPrevSeatPlayer(seatId uint32) *Player {
	prev := (seatId + PLAYER_COUNT - 1) % PLAYER_COUNT
	return d.getPlayerBySeat(prev)
}

// 获取当前位子的下一个位子
func (d *Desk) getNextSeatPlayer(seatId uint32) *Player {
	next := (seatId + 1) % PLAYER_COUNT
	return d.getPlayerBySeat(next)
}

// 获取当前位子的对家
func (d *Desk) getOppositeSeatPlayer(seatId uint32) *Player {
	opposite := (seatId + 2) % PLAYER_COUNT
	return d.getPlayerBySeat(opposite)
}

// SendDeskMessage 广播给桌子上的所有玩家
func (d *Desk) sendDeskMessage(id pb.MID, message proto.Message) {
	for _, p := range d.players {
		if p.IsOnline == false {
			continue
		}
		if p.IsRobot {
			continue
		}
		if p.UserId <= 0 {
			continue
		}
		d.game.SendMessage(p.UserId, uint16(id), message)
	}
}

// SendDeskMessage 广播给桌子上的其他玩家
func (d *Desk) sendOtherPlayerMessage(userId uint32, id pb.MID, message proto.Message) {
	for _, p := range d.players {
		if p.IsOnline == false {
			continue
		}
		if p.IsRobot {
			continue
		}
		if p.UserId <= 0 || p.UserId == userId {
			continue
		}
		d.game.SendMessage(p.UserId, uint16(id), message)
	}
}

func (d *Desk) getGamePlayerCount() int {
	count := 0
	for _, p := range d.players {
		if p.UserId != 0 {
			count++
		}
	}
	return count
}

func (d *Desk) canStart() bool {
	for _, p := range d.players {
		if !p.getPlayerReady() {
			return false
		}
	}
	return true
}

// todo:后续去除
func (d *Desk) canEnd() bool {
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		if p.IsOnline {
			return false
		}
	}
	return true
}

type DeskRunTimeData struct {
	bankerSeatId uint32         // 庄家位子
	cardStack    *lib.CardStack // 牌堆

	operateSeatId       uint32          // 当前可操作位子
	currentStateCanOver bool            // 当前状态是否可以结束
	currentStateCheck   CheckPlaying    // 当前状态是否检测杠胡
	currentBuGangAction *lib.UserAction // 记录补杠操作（用于补杠被抢玩家点过恢复）

	prevActionGangCards lib.Cards // 上一个动作杠的牌(>0是杠，数值是杠的牌)

	swapType pb.SwapType // 换三张类型

	settlementData []*InningScores // 玩家结算细分
	waitSettlement []*InningScores // 待结算分

	canOver bool // 能否结束游戏
}

func NewDeskRunTimeData() *DeskRunTimeData {
	r := &DeskRunTimeData{
		bankerSeatId:        SEAT_UNKNOWN,
		cardStack:           lib.NewCardStack(AllCards),
		operateSeatId:       SEAT_UNKNOWN,
		currentStateCanOver: false,
		currentStateCheck:   CheckGangHu,
		currentBuGangAction: nil,
		prevActionGangCards: make(lib.Cards, 0),
		swapType:            pb.SwapType_SWAP_TYPE_NEXT,
		settlementData:      make([]*InningScores, 0),
		waitSettlement:      make([]*InningScores, 0),
	}
	return r
}

func (d *Desk) randomBankerSeatId() {
	rand.Seed(time.Now().UnixMilli())
	d.runtimeData.bankerSeatId = uint32(rand.Intn(int(PLAYER_COUNT)))
}

func (d *Desk) setOperateSeatId(seatId uint32) {
	d.getPlayerBySeat(seatId).setOpTime(time.Now().UnixMilli())
	d.runtimeData.operateSeatId = seatId
}

func (d *Desk) getOperateSeatId() uint32 {
	return d.runtimeData.operateSeatId
}

func (d *Desk) checkOperateSeatId(seatId uint32) bool {
	return d.runtimeData.operateSeatId == seatId
}

func (d *Desk) setCurrentStateCanOver(canOver bool) {
	d.runtimeData.currentStateCanOver = canOver
}

func (d *Desk) getCurrentStateCanOver() bool {
	return d.runtimeData.currentStateCanOver
}

func (d *Desk) setCurrentStateCheck(c CheckPlaying) {
	d.runtimeData.currentStateCheck = c
}

func (d *Desk) getCurrentStateCheck() CheckPlaying {
	return d.runtimeData.currentStateCheck
}

func (d *Desk) setCurrentBuGangAction(action *lib.UserAction) {
	d.runtimeData.currentBuGangAction = action
}

func (d *Desk) getCurrentBuGangAction() *lib.UserAction {
	return d.runtimeData.currentBuGangAction
}

func (d *Desk) resetCurrentBuGangAction() {
	d.runtimeData.currentBuGangAction = nil
}

func (d *Desk) setPrevActionGangCard(card lib.Card) {
	d.runtimeData.prevActionGangCards.AddCard(card)
}

func (d *Desk) getPrevActionGangCards() lib.Cards {
	return d.runtimeData.prevActionGangCards
}

func (d *Desk) getPrevActionIsGang() int32 {
	return int32(d.runtimeData.prevActionGangCards.Len())
}

func (d *Desk) resetPrevActionGangCards() {
	d.runtimeData.prevActionGangCards = make(lib.Cards, 0)
}

// 记录胡牌分
func (d *Desk) addTotalScore(score *InningScores) {
	log.Sugar().Infof("adding score .....................")
	d.runtimeData.settlementData = append(d.runtimeData.settlementData, score.Copy())
	for i, v := range score.winLoseGolds {
		p := d.getPlayerBySeat(uint32(i))
		p.addTotalScore(v)
	}
}

func (d *Desk) getSettlementData() []*InningScores {
	return d.runtimeData.settlementData
}

// 添加胡牌位子
func (d *Desk) addWaitSettlement(score *InningScores) {
	d.runtimeData.waitSettlement = append(d.runtimeData.waitSettlement, score.Copy())
}

// 获取胡牌位子
func (d *Desk) getWaitSettlement() []*InningScores {
	return d.runtimeData.waitSettlement
}

// 重置胡牌位子
func (d *Desk) resetWaitSettlement() {
	d.runtimeData.waitSettlement = make([]*InningScores, 0)
}

// 游戏是否可以结束
func (d *Desk) setCanOver(f bool) {
	d.runtimeData.canOver = f
}

func (d *Desk) canOver() bool {
	return d.runtimeData.canOver
}

// 清除玩家所有huResult
func (d *Desk) resetHuResult() {
	// 清除玩家胡操作
	for _, p := range d.players {
		p.resetCurrentHuResult()
	}
}
