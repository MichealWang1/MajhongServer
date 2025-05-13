package net

import (
	"context"
	"github.com/panjf2000/gnet/v2"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"net"
)

type Session interface {
	SessionId() string
	SvrType() uint16
	UserId() uint32
	SetSessionId(sessionId string)
	SetSvrType(svrType uint16)
	SetUserId(userId uint32)
	RemoteAddr() string
	Close() error
	Send(data []byte) error
}

type InnerServerSession struct {
	id         string
	server     *server.Server
	conn       net.Conn
	userId     uint32
	svrType    uint16
	remoteAddr string
}

func NewInnerServerSession(id string, server *server.Server, conn net.Conn) *InnerServerSession {
	return &InnerServerSession{
		id:         id,
		server:     server,
		conn:       conn,
		remoteAddr: conn.RemoteAddr().String(),
	}
}

func (s *InnerServerSession) SessionId() string {
	return s.id
}

func (s *InnerServerSession) SvrType() uint16 {
	return s.svrType
}

func (s *InnerServerSession) UserId() uint32 {
	return s.userId
}

func (s *InnerServerSession) SetSessionId(sessionId string) {
	s.id = sessionId
}

func (s *InnerServerSession) SetSvrType(svrType uint16) {
	s.svrType = svrType
}

func (s *InnerServerSession) SetUserId(userId uint32) {
	s.userId = userId
}

func (s *InnerServerSession) RemoteAddr() string {
	return s.remoteAddr
}

func (s *InnerServerSession) Close() error {
	return s.conn.Close()
}

func (s *InnerServerSession) Send(data []byte) error {
	err := s.server.SendMessage(s.conn.(net.Conn), "", "", nil, data)
	return err
}

type InnerClientSession struct {
	id         string
	client     client.XClient
	userId     uint32
	svrType    uint16
	remoteAddr string
}

func NewInnerClientSession(id string, client client.XClient) *InnerClientSession {
	conn, err := client.Stream(context.Background(), nil)
	var remoteAddr string
	if err == nil {
		remoteAddr = conn.RemoteAddr().String()
	}

	return &InnerClientSession{
		id:         id,
		client:     client,
		remoteAddr: remoteAddr,
	}
}

func (s *InnerClientSession) SessionId() string {
	return s.id
}

func (s *InnerClientSession) SvrType() uint16 {
	return s.svrType
}

func (s *InnerClientSession) UserId() uint32 {
	return s.userId
}

func (s *InnerClientSession) SetSessionId(sessionId string) {
	s.id = sessionId
}

func (s *InnerClientSession) SetSvrType(svrType uint16) {
	s.svrType = svrType
}

func (s *InnerClientSession) SetUserId(userId uint32) {
	s.userId = userId
}

func (s *InnerClientSession) RemoteAddr() string {
	return s.remoteAddr
}

func (s *InnerClientSession) Close() error {
	return s.client.Close()
}

func (s *InnerClientSession) Send(data []byte) error {
	err := s.client.Call(context.Background(), "", data, nil)
	return err
}

type OuterSession struct {
	id         string
	conn       gnet.Conn
	userId     uint32
	svrType    uint16
	remoteAddr string
}

func NewOuterSession(id string, conn gnet.Conn, userId uint32) *OuterSession {
	return &OuterSession{
		id:         id,
		conn:       conn,
		userId:     userId,
		remoteAddr: conn.RemoteAddr().String(),
	}
}

func (s *OuterSession) SessionId() string {
	return s.id
}

func (s *OuterSession) SvrType() uint16 {
	return s.svrType
}

func (s *OuterSession) UserId() uint32 {
	return s.userId
}

func (s *OuterSession) SetSessionId(sessionId string) {
	s.id = sessionId
}

func (s *OuterSession) SetSvrType(svrType uint16) {
	s.svrType = svrType
}

func (s *OuterSession) SetUserId(userId uint32) {
	s.userId = userId
}

func (s *OuterSession) RemoteAddr() string {
	return s.remoteAddr
}

func (s *OuterSession) Close() error {
	return s.conn.Close()
}

func (s *OuterSession) Send(data []byte) error {
	_, err := s.conn.Write(data)
	return err
}
