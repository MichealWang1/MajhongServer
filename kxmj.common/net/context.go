package net

import (
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"google.golang.org/protobuf/proto"
	"net"
)

type MsgContext interface {
	Session() Session
	Request() *Message
	Send(msg *Message) error
	Response(body proto.Message) error
	ShouldBind(data proto.Message) error
}

type InnerServerContext struct {
	session *InnerServerSession
	request *Message
}

func NewInnerServerContext(server *server.Server, conn net.Conn, msg *Message) *InnerServerContext {
	return &InnerServerContext{
		session: NewInnerServerSession(conn.RemoteAddr().String(), server, conn),
		request: msg,
	}
}

func (ctx *InnerServerContext) Session() Session {
	return ctx.session
}

func (ctx *InnerServerContext) Request() *Message {
	return ctx.request
}

func (ctx *InnerServerContext) Send(msg *Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	return ctx.session.Send(data)
}

func (ctx *InnerServerContext) Response(body proto.Message) error {
	msg := &Message{
		MsgId:   ctx.request.MsgId,
		SvrType: ctx.request.SvrType,
		SvrId:   ctx.request.SvrId,
		UserId:  ctx.request.UserId,
		Data:    Marshal(body),
	}
	return ctx.Send(msg)
}

func (ctx *InnerServerContext) ShouldBind(data proto.Message) error {
	return ctx.request.Decode(data)
}

type InnerClientContext struct {
	session *InnerClientSession
	request *Message
}

func NewInnerClientContext(servicePath string, client client.XClient, msg *Message) *InnerClientContext {
	return &InnerClientContext{
		session: NewInnerClientSession(servicePath, client),
		request: msg,
	}
}

func (ctx *InnerClientContext) Session() Session {
	return ctx.session
}

func (ctx *InnerClientContext) Request() *Message {
	return ctx.request
}

func (ctx *InnerClientContext) Send(msg *Message) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	return ctx.session.Send(data)
}

func (ctx *InnerClientContext) Response(body proto.Message) error {
	msg := &Message{
		MsgId:   ctx.request.MsgId,
		SvrType: ctx.request.SvrType,
		SvrId:   ctx.request.SvrId,
		UserId:  ctx.request.UserId,
		Data:    Marshal(body),
	}
	return ctx.Send(msg)
}

func (ctx *InnerClientContext) ShouldBind(data proto.Message) error {
	return ctx.request.Decode(data)
}

type OuterContext struct {
	session *OuterSession
	request *Message
	codec   *OuterCodec
}

func NewOuterContext(session *OuterSession, request *Message, codec *OuterCodec) *OuterContext {
	return &OuterContext{
		session: session,
		request: request,
		codec:   codec,
	}
}

func (ctx *OuterContext) Session() Session {
	return ctx.session
}

func (ctx *OuterContext) Request() *Message {
	return ctx.request
}

func (ctx *OuterContext) Send(msg *Message) error {
	data, err := ctx.codec.Encode(msg)
	if err != nil {
		return err
	}
	return ctx.session.Send(data)
}

func (ctx *OuterContext) Response(body proto.Message) error {
	msg := &Message{
		MsgId:   ctx.request.MsgId,
		SvrType: ctx.request.SvrType,
		SvrId:   ctx.request.SvrId,
		UserId:  ctx.request.UserId,
		Data:    Marshal(body),
	}
	return ctx.Send(msg)
}

func (ctx *OuterContext) ShouldBind(data proto.Message) error {
	return ctx.request.Decode(data)
}
