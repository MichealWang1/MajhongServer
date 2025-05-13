package lib

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

// Stater 状态机接口
type Stater interface {
	Enter() // 进入状态机
	Exit()  // 退出状态机
}

// State 状态
type State struct {
	State   string // 状态
	Handler Stater // 状态实现
}

// 状态列表
type States []State

// 事件回调
// TODO：完成这个func
type EventCallback func(proto.Message)

// 事件定义
type Event struct {
	Name     string        // 事件名
	State    string        // 来源状态
	Callback EventCallback // 事件回调方法
}

// 事件列表
type Events []Event

// 事件列表key
type eKey struct {
	Name  string // 事件名
	State string // 状态机内部名
}

// 状态机状态
type MachineStatusType int32

const (
	MachineRunning MachineStatusType = 0 // 状态机运行
	MachineSuspend MachineStatusType = 1 // 状态机挂起
)

// 状态机
type Machine struct {
	currentState string            // 当前状态
	defaultState string            // 默认状态
	states       map[string]Stater // 状态
	events       map[eKey]Event    // 事件
	status       MachineStatusType // 状态机状态(挂起、运行)
}

// NewMachine 创建状态机（state:初始状态；states:状态列表；events:事件列表）
func NewMachine(state string, states States, events Events) (*Machine, error) {
	stateTable := make(map[string]Stater)
	for _, s := range states {
		if _, ok := stateTable[s.State]; ok {
			return nil, fmt.Errorf("NewMachine: state:%v duplicate define", s.State)
		}
		stateTable[s.State] = s.Handler
	}

	eventTable := make(map[eKey]Event)
	for _, e := range events {
		key := eKey{Name: e.Name, State: e.State}
		if _, ok := eventTable[key]; ok {
			return nil, fmt.Errorf("NewMachine: event:%v duplicate define", key)
		}
		eventTable[key] = e
	}

	return &Machine{
		currentState: state,
		defaultState: state,
		states:       stateTable,
		events:       eventTable,
		status:       MachineRunning,
	}, nil
}

// 事件通知入口
func (m *Machine) Event(name string, message proto.Message) error {
	// 判断状态机状态是否挂起
	if m.status != MachineRunning {
		return fmt.Errorf("Machine status:%v is not Running\n", m.status)
	}
	// 获取对应事件的回调方法
	key := eKey{Name: name, State: m.currentState}
	event, ok := m.events[key]
	if !ok {
		return fmt.Errorf("Machine event:%v not define\n", key)
	}
	event.Callback(message)
	return nil
}

// 状态迁移
func (m *Machine) Transition(nextState string) error {
	// 判断要去的状态是否和现在状态一样
	if m.currentState == nextState {
		return nil
	}
	// 获取对应状态机的状态
	handle, ok := m.states[m.currentState]
	if !ok {
		return fmt.Errorf("Machine state:%v not define\n", m.currentState)
	}
	// 判断下一个状态是否存在
	nextHandle, ok := m.states[nextState]
	if !ok {
		return fmt.Errorf("Machine state:%v not define\n", nextState)
	}

	// 退出当前状态
	handle.Exit()
	// 状态机状态切换
	m.currentState = nextState
	// 进入下一个状态
	nextHandle.Enter()
	return nil
}

// 状态机挂起
func (m *Machine) Suspend() {
	if m.status == MachineSuspend {
		return
	}
	m.status = MachineSuspend
}

// 状态机恢复
func (m *Machine) Resume() {
	if m.status == MachineRunning {
		return
	}
	m.status = MachineRunning
}

// 状态机重置
func (m *Machine) Reset() error {
	// 获取当前状态的handle
	handle, ok := m.states[m.currentState]
	if !ok {
		return fmt.Errorf("Machine: Reset: state:%v not define\n", m.currentState)
	}
	// 结束当前状态
	handle.Exit()
	// 重置状态机
	m.currentState = m.defaultState
	m.status = MachineRunning
	return nil
}

// 获取当前状态
func (m *Machine) GetCurrentState() string {
	return m.currentState
}
