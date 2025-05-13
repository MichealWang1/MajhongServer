package game_core

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/smallnest/rpcx/client"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/game_core/server"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/mq"
	"kxmj.common/net"
	"time"
)

type Server struct {
	svrType      uint16             // 服务类型
	svrId        uint16             // 服务ID
	server       *server.RpcxServer // RPCX服务
	game         *Game              // 游戏管理类实例
	redis        *redis.Client      // redis客户端实例
	deskTemplate IDesk              // 桌子接口模板
}

func NewServer(self *net.ServerConfig, etcdEndpoints []string, lobbyConfig *server.RpcxServerConfig, centerConfig *server.RpcxServerConfig, redis *redis.Client, deskTemplate IDesk) *Server {
	return &Server{
		server:       server.NewRpcxServer(self, etcdEndpoints, lobbyConfig, centerConfig, redis),
		svrType:      self.Type,
		svrId:        self.Id,
		redis:        redis,
		deskTemplate: deskTemplate,
	}
}

func (s *Server) Start() {
	game := NewGame(s)
	s.game = game
	s.server.Start(game.Context)
	game.Start()
}

func (s *Server) Close() {
	s.game.Close()

	// 安全关闭
	s.game.WaitClose()

	// 关闭rpcx服务
	s.server.Close()
}

func (s *Server) Template() IDesk {
	return s.deskTemplate
}

func (s *Server) SvrType() uint16 {
	return s.svrType
}

func (s *Server) SvrId() uint16 {
	return s.svrId
}

func (s *Server) GetLobby() client.XClient {
	return s.server.Lobby()
}

func (s *Server) GetCenter() client.XClient {
	return s.server.Center()
}

func (s *Server) GetUserInfo(userId uint32) (*center.GetUserInfoResp, error) {
	reply := &center.GetUserInfoResp{}
	err := s.server.Center().Call(context.Background(), "GetUserInfo", &center.GetUserInfoReq{UserId: userId}, reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Server) CheckUserGold(userId uint32) (*center.CheckUserGoldResp, error) {
	reply := &center.CheckUserGoldResp{}
	err := s.server.Center().Call(context.Background(), "CheckUserGold", &center.CheckUserGoldReq{UserId: userId}, reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Server) GetUserGold(userId uint32, roomId uint32, roomLevel uint8) (*center.GetUserGoldResp, error) {
	reply := &center.GetUserGoldResp{}
	err := s.server.Center().Call(context.Background(), "GetUserGold", &center.GetUserGoldReq{
		UserId:    userId,
		RoomId:    roomId,
		GameId:    s.svrType,
		GameType:  1,
		RoomLevel: roomLevel,
	}, reply)

	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Server) SetUserGold(userId uint32, gold string, roomId uint32, roomLevel uint8) error {
	reply := &center.SetUserGoldResp{}
	err := s.server.Center().Call(context.Background(), "SetUserGold", &center.SetUserGoldReq{
		UserId:    userId,
		Gold:      gold,
		RoomId:    roomId,
		GameId:    s.svrType,
		GameType:  1,
		RoomLevel: roomLevel,
	}, reply)

	if err != nil {
		return err
	}

	if reply.Code != codes.Success {
		return codes.New(reply.Code, reply.Msg)
	}

	return nil
}

func (s *Server) CheckUserDiamond(userId uint32) (string, error) {
	reply := &center.CheckUserDiamondResp{}
	err := s.server.Center().Call(context.Background(), "CheckUserDiamond", &center.CheckUserDiamondReq{
		UserId: userId,
	}, reply)

	if err != nil {
		return "0", err
	}

	if reply.Code != codes.Success {
		return "0", codes.New(reply.Code, reply.Msg)
	}

	return reply.Data.Diamond, nil
}

func (s *Server) AddRecord(record interface{}) {
	err := mq.AddGameLogger(record)
	if err != nil {
		log.Sugar().Errorf("AddGameLogger:%v err:%v", record, err)
	}
}

func (s *Server) GetRoomConfigList(gameId uint16) (*center.RoomConfigListResp, error) {
	reply := &center.RoomConfigListResp{}
	err := s.server.Center().Call(context.Background(), "GetRoomConfigList", &center.RoomConfigListReq{GameId: gameId}, reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Server) GetRoomConfig(gameId uint16, roomId uint32) (*center.RoomConfigResp, error) {
	reply := &center.RoomConfigResp{}
	err := s.server.Center().Call(context.Background(), "GetRoomConfig", &center.RoomConfigReq{GameId: gameId, RoomId: roomId}, reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Server) UpdateStatistics(statistics []*UserStatistics) {
	if len(statistics) <= 0 {
		return
	}

	for _, d := range statistics {
		if d.TotalTimes <= 0 {
			continue
		}

		err := mq.UpdateStatistics(&kxmj_core.UserGameStatistics{
			Id:            0,
			UserId:        d.UserId,
			GameId:        d.GameId,
			GameType:      d.GameType,
			RoomLevel:     d.RoomLevel,
			PlayType:      d.PlayType,
			TotalTimes:    d.TotalTimes,
			TotalWinLoss:  d.TotalWinLoss,
			TotalDuration: d.TotalDuration,
			CreatedAt:     uint32(time.Now().Unix()),
			UpdatedAt:     uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("UpdateStatistics err:%v", err)
		}
	}
}
