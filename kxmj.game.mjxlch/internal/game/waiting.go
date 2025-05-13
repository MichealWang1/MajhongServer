package game

import "time"

type WaitingEvent struct {
	nextStatus Status // 下一个状态
	nextTime   int64  // 等待时间
	has        bool   // 是否有事件
}

func NewWaitingEvent(nextStatus Status, nextTime int64, has bool) *WaitingEvent {
	return &WaitingEvent{nextStatus: nextStatus, nextTime: nextTime, has: has}
}

func (e *WaitingEvent) set(nextStatus Status, duration int64, has bool) {
	e.nextStatus = nextStatus
	e.nextTime = time.Now().UnixMilli() + duration
	e.has = true
}

func (e *WaitingEvent) wait() bool {
	if !e.has {
		return false
	}
	return time.Now().UnixMilli() < e.nextTime
}

func (e *WaitingEvent) next() Status {
	return e.nextStatus
}

func (e *WaitingEvent) hasEvent() bool {
	return e.has
}
