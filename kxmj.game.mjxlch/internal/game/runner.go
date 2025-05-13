package game

import (
	"kxmj.common/log"
	"math/rand"
	"time"
)

// swapStatus 状态切换
func (d *Desk) swapStatus(status Status) {
	d.status = status
	d.markTime = time.Now().UnixMilli()
	switch status {
	case Match:
		d.nextTime = d.markTime + DurationMatch + rand.Int63n(1000)
		d.toMatch()
	case Ready:
		d.nextTime = d.markTime + DurationReady
		d.toReady()
	case Dice:
		d.nextTime = d.markTime + DurationDice
		d.toDice()
	case Deal:
		d.nextTime = d.markTime + DurationDealCard
		d.toDeal()
	case Swap:
		d.nextTime = d.markTime + DurationSwap
		d.toSwap()
	case ChooseMiss:
		d.nextTime = d.markTime + DurationChooseMiss
		d.toChooseMiss()
	case Playing:
		d.nextTime = d.markTime + DurationPlaying
		d.toPlaying()
	case Operate:
		d.nextTime = d.markTime + DurationOperate
		d.toOperate()
	case Settle:
		d.nextTime = d.markTime + DurationSettle
		d.toSettle()
	case End:
		d.nextTime = d.markTime + DurationEnd
		d.toEnd()
	case Pause:
		d.nextTime = d.markTime + DurationPause
		d.toPause()
	}
	// 广播游戏状态
	//d.broadcastGameStatusNotify()
}

func (d *Desk) toMatch() {
}

func (d *Desk) toReady() {
}

func (d *Desk) toDice() {
	// 定庄
	//d.runTimeData.DiceData.DetermineBankerSeatId(len(d.players))
}

func (d *Desk) toDeal() {
	d.dealCards()
}

func (d *Desk) toSwap() {
	// 提示玩家换三张
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		d.sendPlayerSwapCardNotify(p.SeatId)
	}
	// 随机换牌
	d.randomSwapType()
}

func (d *Desk) toChooseMiss() {
	// 提示玩家选缺
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		d.sendPlayerChooseMissNotify(p.SeatId)
	}
}

func (d *Desk) toPlaying() {
	if d.canOver() {
		d.waiting.set(End, EndWaitDuration, true)
		return
	}
	seatId := d.getOperateSeatId()

	// 清除玩家当前胡牌数据
	d.resetHuResult()

	//d.getPlayerBySeat(seatId).setOpTime(time.Now().UnixMilli())
	log.Sugar().Infof("operate seat:%d,handCards:%v", seatId, d.getPlayerBySeat(seatId).getHandCards())
	// 提示玩家出牌
	d.sendPlayerOutCardNotify(seatId)
	// 提示玩家听牌信息
	d.calculateTingCards3N2(seatId)
	// 检测玩家动作
	if actions := d.detectionActionOnSelf(seatId, d.getCurrentStateCheck()); actions.Len() > 0 {
		//d.sendPlayerOperateNotify(uint32(seatId), actions)
		p := d.getPlayerBySeat(seatId)
		p.setOperationalActions(actions)
		// 设置当前状态
		d.setCurrentStateCanOver(true)
		return
	}
}

func (d *Desk) toOperate() {
	d.setCurrentStateCanOver(false)
	// 通知有操作玩家进行操作
	for _, p := range d.players {
		if p.IsRobot {
			continue
		}
		if p.getOperationalActions().Len() > 0 {
			d.sendPlayerOperateNotify(p.SeatId, p.getOperationalActions())
		}
	}
}

func (d *Desk) toSettle() {
	// 进行结算
	d.calculateGold(d.getWaitSettlement())

	// 广播玩家算分
	d.broadcastBureauSettlementNotify()
}

func (d *Desk) toEnd() {
	d.broadcastEndSettlementNotify()
}

func (d *Desk) toPause() {
}

func (d *Desk) isWait() bool {
	nowTime := time.Now().UnixMilli()
	if nowTime < d.nextTime {
		return true
	}
	return false
}

func (d *Desk) run() {

	if d.waiting.hasEvent() {
		if !d.waiting.wait() {
			d.waiting.has = false
			d.swapStatus(d.waiting.next())
		}
		return
	}

	if d.canEnd() {
		log.Sugar().Infof("player offline end")
		d.onEnd()
		return
	}
	switch d.status {
	case Match:
		d.onMatch()
	case Ready:
		d.onReady()
	case Dice:
		d.onDice()
	case Deal:
		d.onDeal()
	case Swap:
		d.onSwap()
	case ChooseMiss:
		d.onChooseMiss()
	case Playing:
		d.onPlaying()
	case Operate:
		d.onOperate()
	case Settle:
		d.onSettle()
	case End:
		d.onEnd()
	case Pause:
		d.onPause()
	}
}

func (d *Desk) onReady() {
	wait := d.isWait()
	// 玩家没有4人则添加机器人
	if !wait && d.getGamePlayerCount() < int(PLAYER_COUNT) {
		d.waiting.set(Match, MatchWaitDuration, true)
		return
	}
	if wait || !d.canStart() {
		return
	}
	// 获取玩家金币
	d.broadcastPlayerGoldNumber()
	// 设置牌堆
	d.setCardStack()
	// 通知所有玩家游戏开始
	d.broadcastGameStart()
	d.waiting.set(Dice, DiceWaitDuration, true)
}

func (d *Desk) onMatch() {
	if d.getGamePlayerCount() == int(PLAYER_COUNT) {
		d.waiting.set(Ready, ReadyWaitDuration, true)
		return
	}
	if d.isWait() {
		return
	}
	log.Sugar().Infof("match robot: %v", d.room.MatchRobot())
	if d.room.MatchRobot() {
		// 玩家没有4人则添加机器人
		d.addRobot()
		if d.getGamePlayerCount() != int(PLAYER_COUNT) {
			d.waiting.set(Match, MatchWaitDuration, true)
		}
	}
}

func (d *Desk) onDice() {
	if d.isWait() {
		return
	}

	d.waiting.set(Deal, DealWaitDuration, true)
}

func (d *Desk) onDeal() {
	if d.isWait() {
		return
	}

	d.waiting.set(Swap, SwapWaitDuration, true)
}

func (d *Desk) onSwap() {
	if d.isWait() && !d.swapStateCanOver() {
		return
	}
	// 机器人换牌
	//d.onRobotSwap()

	// 没有换牌玩家进行换牌
	d.onSwapTimeOver()

	// 进行换牌操作
	d.swapStart()

	d.waiting.set(ChooseMiss, ChooseMissWaitDuration, true)
	//d.toChooseMiss()
}

func (d *Desk) onChooseMiss() {
	if d.isWait() && !d.chooseMissStateCanOver() {
		return
	}
	// 机器人选缺
	//d.onRobotChooseMiss()
	// 处理超时玩家
	d.onChooseMissTimeOver()
	// 广播选缺结果
	d.broadcastPlayerChooseMissResultNotify()

	d.waiting.set(Playing, PlayingWaitDuration, true)
}

func (d *Desk) onPlaying() {
	wait := d.isWait()
	if !d.getCurrentStateCanOver() && (!wait || d.getPlayerBySeat(d.getOperateSeatId()).canAutoOperation() || d.getPlayerBySeat(d.getOperateSeatId()).isHu()) {
		log.Sugar().Infof("wait:", wait, ",time now:", time.Now().UnixMilli(), "plyer over time:", d.getPlayerBySeat(d.getOperateSeatId()).getOpTime()+DurationAuto)
		// 处理超时
		d.onOutCardTimeOver()
	}
	if wait && !d.getCurrentStateCanOver() {
		return
	}
	d.setCurrentStateCanOver(false)
	// 玩家有操作
	if d.hasOperateActions() {
		d.waiting.set(Operate, OperateWaitDuration, true)
		return
	}
	log.Sugar().Infof("nextSeatId:%v", d.getOperateSeatId())
	// 下一个人出牌
	d.waiting.set(Playing, PlayingWaitDuration, true)
}

func (d *Desk) onOperate() {
	wait := d.isWait()
	d.onOperateAuto(wait)
	if wait && !d.getCurrentStateCanOver() {
		return
	}
	d.setCurrentStateCanOver(false)

	// 玩家有操作
	if d.hasOperateActions() {
		d.waiting.set(Operate, OperateWaitDuration, true)
		return
	}
	// 有玩家点了胡、杠
	if len(d.getWaitSettlement()) > 0 {
		d.waiting.set(Settle, SettleWaitDuration, true)
		return
	}

	// 下一个人出牌
	d.waiting.set(Playing, PlayingWaitDuration, true)
}

func (d *Desk) onSettle() {
	if d.isWait() {
		return
	}
	// 清理结算数据
	d.resetWaitSettlement()

	if !d.canOver() {
		// 设置是否检测胡杠
		d.setCurrentStateCheck(CheckGangHu)
		// 玩家摸牌
		err := d.catchCard(d.getOperateSeatId())
		if err != nil {
			log.Sugar().Infof("%v", err)
			return
		}
		d.waiting.set(Playing, PlayingWaitDuration, true)
		return
	}

	d.roundId = d.newRoundId()
	d.waiting.set(End, EndWaitDuration, true)
	//d.toPause()
}

func (d *Desk) onEnd() {
	log.Sugar().Infof("onEnd ")
	for _, p := range d.players {
		// 通知玩家离开
		d.PlayerOnLeave(p.SeatId)
		// 清理桌子
		p.Reset()
	}
	log.Sugar().Infof("End....")
	d.game.OnDeskClose(d)
}

func (d *Desk) onPause() {

}
