package game_manager

import (
	"kxmj-gm/internal/server/rpc"
)

type Service struct {
}

var (
	Controller = &Service{}
)

func (s *Service) XServer() *rpc.RpcxServer {
	return rpc.Default()
}
