package server

import "kxmj.common/net"

type User struct {
	UserId  uint32      // 用户ID
	SvrType uint16      // 服务类型
	SvrId   uint16      // 服务分布式ID(由服务端返回)
	RoomId  uint32      // 房间ID
	DeskId  uint32      // 桌子ID
	Gateway net.Session // 用户当前连接网关
	Game    net.Session // 用户当前所在游戏
}

type Endpoint struct {
	SvrType uint16
	SvrId   uint16
	Addr    string
	Port    uint32
	Users   uint32
	Session net.Session
}
