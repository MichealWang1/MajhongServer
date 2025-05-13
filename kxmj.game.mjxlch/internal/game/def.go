package game

import (
	lib "kxmj.common/mahjong"
	"kxmj.game.mjxlch/pb"
)

// 测试牌堆
const TEST_DECK = true

type Status uint32

const (
	DurationReady      = 2000   // 准备状态时长
	DurationMatch      = 1000   // 匹配机器人时长
	DurationDice       = 500    // 掷骰子状态时长
	DurationDealCard   = 2000   // 发牌状态时长
	DurationSwap       = 15000  // 换牌状态时长
	DurationChooseMiss = 15000  // 选缺状态时长
	DurationPlaying    = 15000  // 打牌状态
	DurationOperate    = 15000  // 动作操作状态时长
	DurationSettle     = 1000   // 结算状态
	DurationEnd        = 2000   // 结束状态时长
	DurationPause      = 100000 // 暂停状态时长

	DurationAuto = 2000 // 自动出牌操作时间
	DurationOut  = 500  // 出牌操作得等0.5s

	MatchWaitDuration      = 2000 // 匹配等待时间
	ReadyWaitDuration      = 0    // 准备等待时间
	DiceWaitDuration       = 0    // 掷骰子等待时间
	DealWaitDuration       = 1000 // 发牌等待时间
	SwapWaitDuration       = 1000 // 换牌等待时间
	ChooseMissWaitDuration = 1000 // 选缺等待时间
	PlayingWaitDuration    = 500  // 出牌等待时间
	OperateWaitDuration    = 500  // 操作等待时间
	SettleWaitDuration     = 500  // 结算等待时间
	EndWaitDuration        = 1000 // 结束等待时间
)

const (
	Match      Status = Status(pb.GameState_MATCH)           // 匹配状态
	Ready      Status = Status(pb.GameState_READY)           // 准备状态
	Dice       Status = Status(pb.GameState_DICE)            // 掷骰子状态
	Deal       Status = Status(pb.GameState_DEAL_HAND_CARDS) // 发牌状态
	Swap       Status = Status(pb.GameState_SWAP)            // 换牌状态（换三张）
	ChooseMiss Status = Status(pb.GameState_CHOOSE_MISS)     // 选缺状态
	Playing    Status = Status(pb.GameState_PLAY)            // 打牌状态
	Operate    Status = Status(pb.GameState_OPERATE)         // 动作操作状态
	Settle     Status = Status(pb.GameState_SETTLEMENT)      // 结算状态
	End        Status = Status(pb.GameState_END)             // 结束状态
	Pause      Status = Status(pb.GameState_PAUSE)           // 暂停状态
)

const PLAYER_COUNT uint32 = 4

// 未知座位号占位
const SEAT_UNKNOWN = 0xff

// 所有牌
var AllCards = lib.Cards{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // 1-9万

	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, // 1-9条

	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, // 1-9筒
	//
	//0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, // 东南西北中发白
	//0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, // 梅兰竹菊春夏秋冬
}

// toPlaying 检测类型
type CheckPlaying int32

const (
	CheckGangHu CheckPlaying = iota + 1 // 检测胡杠（摸牌后）
	CheckGang                           // 检测杠（吃碰后）
	NotCheck                            // 不检测胡杠(自身操作后)
)

// 胡牌权位
const (
	Ping_Hu = 89 // 平胡
	Qi_Dui  = 18 // 七对

	Gen      = 93 // 根
	Gang_Pao = 98 // 杠上炮
)
