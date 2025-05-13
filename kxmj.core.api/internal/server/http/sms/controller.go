package sms

import "kxmj.core.api/internal/server/rpc"

type Service struct {
}

var (
	Controller = &Service{}
)

func (s *Service) XServer() *rpc.RpcxServer {
	return rpc.Default()
}
