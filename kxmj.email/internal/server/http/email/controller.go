package email

import "kxmj.email/internal/server/rpc"

type Service struct {
}

var (
	Controller = &Service{}
)

func (s *Service) XServer() *rpc.RpcxServer {
	return rpc.Default()
}
