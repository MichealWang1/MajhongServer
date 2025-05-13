package net

// Message 服务消息包
type Message struct {
	MsgId   uint16 // 消息ID
	SvrType uint16 // 服务类型
	SvrId   uint16 // 服务分布式ID(由服务端返回)
	UserId  uint32 // 用户ID
	Data    []byte // 包体byte数组
}

// RpcxReply 占位符
type RpcxReply struct {
}

func (msg *Message) Length() int {
	return len(msg.Data)
}
