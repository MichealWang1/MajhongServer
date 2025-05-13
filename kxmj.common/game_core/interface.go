package game_core

import (
	"github.com/smallnest/rpcx/client"
	"google.golang.org/protobuf/proto"
	"kxmj.common/model/center"
	"kxmj.common/net"
	"kxmj.common/redis_cache/gm"
)

type IDesk interface {
	ID() uint32                       // 桌子ID
	RoundId() string                  // 当前局号
	Room() IRoom                      // 房间实例
	New(game IGame, room IRoom) IDesk // 创建桌子
	Start()                           // 开始
	Close()                           // 关闭
	Run()                             // 状态机事件
	OnMessage(ctx net.MsgContext)     // 接收消息
}

type IRoom interface {
	Update(config *center.RoomConfig) // 更新配置
	Close()                           // 关闭房间
	ID() uint32                       // 房间ID
	GameId() uint16                   // 游戏ID
	GameType() uint8                  // 游戏类型：1 麻将；2 斗地主
	RoomLevel() uint8                 // 房间级别：1 初级场；2 中级场；3 高级场；4 大师场；5 圣雀场；
	MinLimit() string                 // 最小进场限制：0 代表不限制
	MaxLimit() string                 // 最大进场限制：0 代表不限制
	BaseScore() string                // 底分
	MaxMultiple() uint32              // 最大倍数
	Ticket() string                   // 门票
	MatchRobot() bool                 // 是否匹配机器人
	GenerateDeskId() uint32           // 生成桌子ID
	RemoveDeskId(deskId uint32)       // 移除桌子ID
}

type IGame interface {
	Server() IServer                                                      // 服务对象
	OnEnter(enterInfo *EnterDesk)                                         // 进入游戏事件
	OnLeave(leaveInfo *LeaveDesk)                                         // 离开游戏事件
	OnDeskClose(desk IDesk)                                               // 桌子关闭事件
	SendMessage(userId uint32, msgId uint16, data proto.Message)          // 发送消息
	SendErrMessage(ctx net.MsgContext, code int)                          // 发送错误消息
	GetRobot(baseScore string) *RobotInfo                                 // 获取机器人信息
	NotifyRise(userId uint32, desk IDesk)                                 // 通知用户复活
	GetManualConfig(userId uint32, desk IDesk) (*gm.CardStackData, error) // 手动配置
}

type IServer interface {
	Start()                                                                                     // 启动服务
	Close()                                                                                     // 关闭服务
	Template() IDesk                                                                            // 获取桌子模板
	SvrType() uint16                                                                            // 服务类型
	SvrId() uint16                                                                              // 服务ID
	GetLobby() client.XClient                                                                   // 获取大厅客户端
	GetCenter() client.XClient                                                                  // 获取中心客户端
	GetUserInfo(userId uint32) (*center.GetUserInfoResp, error)                                 // 获取用户信息
	CheckUserGold(userId uint32) (*center.CheckUserGoldResp, error)                             // 检查用户金币
	GetUserGold(userId uint32, roomId uint32, roomLevel uint8) (*center.GetUserGoldResp, error) // 获取用户金豆(钱包转到游戏)
	SetUserGold(userId uint32, gold string, roomId uint32, roomLevel uint8) error               // 带出用户金豆(游戏转到钱包)
	CheckUserDiamond(userId uint32) (string, error)                                             // 检查用户钻石数
	AddRecord(record interface{})                                                               // 写游戏日志
	GetRoomConfigList(gameId uint16) (*center.RoomConfigListResp, error)                        // 获取房间配置列表
	GetRoomConfig(gameId uint16, roomId uint32) (*center.RoomConfigResp, error)                 // 获取房间配置
	UpdateStatistics(statistics []*UserStatistics)                                              // 更新用户游戏统计
}
