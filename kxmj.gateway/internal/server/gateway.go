package server

import (
	"fmt"
	"kxmj.common/log"
	"kxmj.common/net"
)

type Gateway struct {
	Inner *Inner
	Outer *Outer
}

func (g *Gateway) Start() {
	g.Inner.Start()
	g.Outer.Start()
}

func (g *Gateway) Close() {
	g.Inner.Close()
	g.Outer.Close()
}

func (g *Gateway) ToInner(msg *net.Message) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%d:U:%d:L%d] ToInner <---", msg.MsgId, msg.SvrType, msg.SvrId, msg.UserId, msg.Length()))
	g.Inner.ToInner(msg)
}

func (g *Gateway) ToOuter(msg *net.Message) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%d:U:%d:L%d] ToOuter --->", msg.MsgId, msg.SvrType, msg.SvrId, msg.UserId, msg.Length()))
	g.Outer.ToOuter(msg)
}
