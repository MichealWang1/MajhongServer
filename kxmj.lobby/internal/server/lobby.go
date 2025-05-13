package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kxmj.common/codes"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/lobby"
	"kxmj.common/net"
	"kxmj.common/proto/gateway_pb"
	"kxmj.common/proto/lobby_pb"
	"kxmj.common/redis_cache"
	"kxmj.common/service"
	"math/rand"
	"sort"
	"sync"
)

type LobbyServer struct {
	redis      *redis.Client                // redis服务
	server     *RpcxServer                  // rpcx服务
	users      map[uint32]*User             // 用户位置信息
	games      map[string]*Endpoint         // 游戏服连接session
	gateways   map[string]*Endpoint         // 网关连接session
	mutex      sync.Mutex                   // 同步锁
	onlineMaps map[uint16]map[uint32]uint32 // 在线人数
}

func NewLobbyServer(self *net.ServerConfig, etcdEndpoints []string, redis *redis.Client) *LobbyServer {
	l := &LobbyServer{
		redis:      redis,
		users:      make(map[uint32]*User, 0),
		games:      make(map[string]*Endpoint, 0),
		gateways:   make(map[string]*Endpoint, 0),
		mutex:      sync.Mutex{},
		onlineMaps: make(map[uint16]map[uint32]uint32, 0),
	}
	l.server = NewRpcxServer(self, etcdEndpoints, redis, l)
	return l
}

func (ls *LobbyServer) Start() {
	ls.server.RegisterRedis()
	ls.server.StartRpcxServer(ls.OnInput)
}

func (ls *LobbyServer) Close() {
	ls.server.Close()
}

func (ls *LobbyServer) OnInput(ctx net.MsgContext) {
	log.Sugar().Infof("[M:%dT:%dID:%dU:%d] <--", ctx.Request().MsgId, ctx.Request().SvrType, ctx.Request().SvrId, ctx.Request().UserId)
	switch ctx.Request().MsgId {
	case uint16(lobby_pb.MID_REGISTER):
		ls.register(ctx)
	case uint16(lobby_pb.MID_UNREGISTER):
		ls.unRegister(ctx)
	case uint16(lobby_pb.MID_ON_LINE):
		ls.online(ctx)
	case uint16(lobby_pb.MID_OFF_LINE):
		ls.offline(ctx)
	case uint16(gateway_pb.MID_MATCH):
		ls.match(ctx)
	case uint16(lobby_pb.MID_ENTER_DESK):
		ls.enter(ctx)
	case uint16(lobby_pb.MID_LEAVE_DESK):
		ls.leave(ctx)
	case uint16(gateway_pb.MID_LOCATION):
		ls.checkLocation(ctx)
	}
}

func (ls *LobbyServer) register(ctx net.MsgContext) {
	payload := &lobby_pb.Register{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	servicePath := service.ParseServicePath(uint16(payload.SvrType), uint16(payload.SvrId))
	if uint16(payload.SvrType) == service.GatewayService {
		endpoint, has := ls.gateways[servicePath]
		if has == false {
			endpoint = &Endpoint{
				SvrType: uint16(payload.SvrType),
				SvrId:   uint16(payload.SvrId),
				Addr:    payload.Addr,
				Port:    payload.Port,
				Users:   0,
				Session: ctx.Session(),
			}
		}
		ls.gateways[servicePath] = endpoint
	} else {
		endpoint, has := ls.games[servicePath]
		if has == false {
			endpoint = &Endpoint{
				SvrType: uint16(payload.SvrType),
				SvrId:   uint16(payload.SvrId),
				Addr:    payload.Addr,
				Port:    payload.Port,
				Users:   0,
				Session: ctx.Session(),
			}
		}
		ls.games[servicePath] = endpoint
	}
}

func (ls *LobbyServer) unRegister(ctx net.MsgContext) {
	payload := &lobby_pb.UnRegister{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	servicePath := service.ParseServicePath(uint16(payload.SvrType), uint16(payload.SvrId))
	if uint16(payload.SvrType) == service.GatewayService {
		delete(ls.gateways, servicePath)
	} else {
		delete(ls.games, servicePath)
	}
}

func (ls *LobbyServer) online(ctx net.MsgContext) {
	payload := &lobby_pb.Online{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	ls.mutex.Lock()
	user, has := ls.users[payload.UserId]
	if has == false {
		user = &User{
			UserId:  payload.UserId,
			SvrType: 0,
			SvrId:   0,
			RoomId:  0,
			DeskId:  0,
			Gateway: ctx.Session(),
			Game:    nil,
		}
	}
	user.Gateway = ctx.Session()
	ls.users[payload.UserId] = user

	ls.mutex.Unlock()

	// 用户上线通知游戏服务
	if user.Game != nil {
		data, err := ctx.Request().Encode()
		if err != nil {
			log.Sugar().Errorf("Encode err:%v", err)
			return
		}

		err = user.Game.Send(data)
		if err != nil {
			log.Sugar().Errorf("Send err:%v", err)
		}
	}
}

func (ls *LobbyServer) offline(ctx net.MsgContext) {
	payload := &lobby_pb.Offline{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	ls.mutex.Lock()

	user, has := ls.users[payload.UserId]
	if has == false {
		user = &User{
			UserId:  payload.UserId,
			SvrType: 0,
			SvrId:   0,
			RoomId:  0,
			DeskId:  0,
			Gateway: nil,
			Game:    nil,
		}
	}
	ls.users[payload.UserId] = user

	ls.mutex.Unlock()
	if user.Gateway != nil {
		if user.Gateway.SessionId() == ctx.Session().SessionId() {
			user.Gateway = nil
		}

		// 用户掉线通知游戏服务
		if user.Game != nil {
			data, err := ctx.Request().Encode()
			if err != nil {
				log.Sugar().Errorf("Encode err:%v", err)
				return
			}

			err = user.Game.Send(data)
			if err != nil {
				log.Sugar().Errorf("Send err:%v", err)
			}
		}
	}
}

func (ls *LobbyServer) match(ctx net.MsgContext) {
	payload := &gateway_pb.MatchReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	var endpoint *Endpoint
	for _, val := range ls.games {
		if val.SvrType == uint16(payload.SvrType) && val.Users <= 20000 {
			endpoint = val
			break
		}
	}

	if endpoint == nil {
		err = ctx.Send(&net.Message{
			MsgId:   uint16(gateway_pb.MID_ERR),
			SvrType: ctx.Request().SvrType,
			SvrId:   ctx.Request().SvrId,
			UserId:  ctx.Request().UserId,
			Data: net.Marshal(&gateway_pb.Err{
				Code:        uint32(codes.ServerNotFound),
				Msg:         codes.GetMessage(codes.ServerNotFound),
				OriginMsgId: uint32(ctx.Request().MsgId),
			}),
		})

		if err != nil {
			log.Sugar().Errorf("Send err:%v", err)
		}
		return
	}

	err = ctx.Send(&net.Message{
		MsgId:   uint16(gateway_pb.MID_MATCH),
		SvrType: ctx.Request().SvrType,
		SvrId:   endpoint.SvrId,
		UserId:  ctx.Request().UserId,
		Data: net.Marshal(&gateway_pb.MatchResp{
			UserId:  ctx.Request().UserId,
			SvrType: payload.SvrType,
			SvrId:   uint32(endpoint.SvrId),
			RoomId:  payload.RoomId,
			DeskId:  0,
		}),
	})

	if err != nil {
		log.Sugar().Errorf("Send err:%v", err)
	}
}

func (ls *LobbyServer) enter(ctx net.MsgContext) {
	payload := &lobby_pb.EnterDesk{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	// 房间随机在线人数
	gameOnline, has := ls.onlineMaps[uint16(payload.SvrType)]
	if has == false {
		gameOnline = make(map[uint32]uint32, 0)
	}
	ls.onlineMaps[uint16(payload.SvrType)] = gameOnline

	roomOnline := gameOnline[payload.RoomId]
	if roomOnline == 0 {
		roomOnline = uint32(rand.Intn(2000) + 100)
	}

	roomOnline++
	gameOnline[payload.RoomId] = roomOnline

	// 用户位置
	user, has := ls.users[payload.UserId]
	if has == false {
		user = &User{
			UserId:  0,
			SvrType: 0,
			SvrId:   0,
			RoomId:  0,
			DeskId:  0,
			Gateway: nil,
			Game:    nil,
		}
	}

	user.UserId = payload.UserId
	user.SvrType = uint16(payload.SvrType)
	user.SvrId = uint16(payload.SvrId)
	user.RoomId = payload.RoomId
	user.DeskId = payload.DeskId
	user.Game = ctx.Session()
	ls.users[payload.UserId] = user

	// 游戏服务实际在线人数
	servicePath := service.ParseServicePath(uint16(payload.SvrType), uint16(payload.SvrId))
	endpoint, has := ls.games[servicePath]
	if has {
		endpoint.Users++
	}
}

func (ls *LobbyServer) leave(ctx net.MsgContext) {
	payload := &lobby_pb.LeaveDesk{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Errorf("ShouldBind err:%v", err)
		return
	}

	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	// 房间随机在线人数
	gameOnline, has := ls.onlineMaps[uint16(payload.SvrType)]
	if has == false {
		gameOnline = make(map[uint32]uint32, 0)
	}
	ls.onlineMaps[uint16(payload.SvrType)] = gameOnline

	roomOnline := gameOnline[payload.RoomId]
	if roomOnline > 0 {
		roomOnline--
	}
	gameOnline[payload.RoomId] = roomOnline

	user, has := ls.users[payload.UserId]
	if has == false {
		user = &User{
			UserId:  0,
			SvrType: 0,
			SvrId:   0,
			RoomId:  0,
			DeskId:  0,
			Gateway: nil,
			Game:    nil,
		}
		return
	}

	if user.UserId == payload.UserId &&
		user.SvrType == uint16(payload.SvrType) &&
		user.SvrId == uint16(payload.SvrId) &&
		user.RoomId == payload.RoomId &&
		user.DeskId == payload.DeskId &&
		user.Game.SessionId() == ctx.Session().SessionId() {
		user.SvrType = 0
		user.SvrId = 0
		user.RoomId = 0
		user.DeskId = 0
		user.Game = nil
	}
	ls.users[payload.UserId] = user

	// 服务在线人数
	servicePath := service.ParseServicePath(uint16(payload.SvrType), uint16(payload.SvrId))
	endpoint, has := ls.games[servicePath]
	if has {
		endpoint.Users--
	}
}

func (ls *LobbyServer) checkLocation(ctx net.MsgContext) {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	userId := ctx.Request().UserId
	user, has := ls.users[userId]
	if has == false {
		err := ctx.Response(&gateway_pb.Location{})
		if err != nil {
			log.Sugar().Errorf("Response err:%v", err)
		}
		return
	}

	err := ctx.Response(&gateway_pb.Location{
		SvrType: uint32(user.SvrType),
		SvrId:   uint32(user.SvrId),
		RoomId:  user.RoomId,
		DeskId:  user.DeskId,
	})

	if err != nil {
		log.Sugar().Errorf("Response err:%v", err)
	}
}

func (ls *LobbyServer) getLocation(userId uint32) *lobby.LocationInfo {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	user, has := ls.users[userId]
	if has == false {
		return &lobby.LocationInfo{}
	}

	return &lobby.LocationInfo{
		UserId:  user.UserId,
		SvrType: user.SvrType,
		SvrId:   user.SvrId,
		RoomId:  user.RoomId,
		DeskId:  user.DeskId,
	}
}

func (ls *LobbyServer) getGateways() []*Endpoint {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	list := make([]*Endpoint, 0)
	for _, endpoint := range ls.gateways {
		list = append(list, endpoint)
	}

	if len(list) > 0 {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Users < list[j].Users
		})
	}

	return list
}

func (ls *LobbyServer) shoppingSuccess(parameter *lobby.ShoppingSuccessReq) {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	user, has := ls.users[parameter.UserId]
	if has == false {
		return
	}

	if user.Gateway == nil {
		return
	}

	msg := &net.Message{
		MsgId:   uint16(gateway_pb.MID_SHOPPING_SUCCESS),
		SvrType: ls.server.self.Type,
		SvrId:   0,
		UserId:  parameter.UserId,
		Data: net.Marshal(&gateway_pb.ShoppingSuccess{
			GoodsId:        parameter.GoodsId,
			Name:           parameter.Name,
			ShopType:       uint32(parameter.ShopType),
			Price:          parameter.Price,
			RealCount:      parameter.RealCount,
			FirstBuyDouble: uint32(parameter.FirstBuyDouble),
			ItemId:         parameter.ItemId,
		}),
	}

	data, err := msg.Encode()
	if err != nil {
		log.Sugar().Errorf("Encode msg:%d err:%v", msg.MsgId, err)
		return
	}

	err = user.Gateway.Send(data)
	if err != nil {
		log.Sugar().Errorf("Send msg:%d err:%v", msg.MsgId, err)
	}

	// 检查消息是否要转发到游戏
	if user.Game == nil {
		return
	}

	itemConfig, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(context.Background())
	if err != nil {
		log.Sugar().Errorf("Get item config msg:%d err:%v", msg.MsgId, err)
	}

	d, has := itemConfig[parameter.ItemId]
	if has == false {
		return
	}

	if item.Type(d.ItemType) == item.RisePack {
		riseMsg := &net.Message{
			MsgId:   uint16(lobby_pb.MID_RISE_BUY_SUCCESS),
			SvrType: ls.server.self.Type,
			SvrId:   0,
			UserId:  parameter.UserId,
			Data:    net.Marshal(&lobby_pb.RiseBuySuccess{}),
		}

		riseData, err := riseMsg.Encode()
		if err != nil {
			log.Sugar().Errorf("Encode msg:%d err:%v", riseMsg.MsgId, err)
			return
		}

		err = user.Game.Send(riseData)
		if err != nil {
			log.Sugar().Errorf("Send msg:%d err:%v", riseMsg.MsgId, err)
		}
	}
}

func (ls *LobbyServer) getRoomsOnline(gameId uint16, roomIds []uint32) map[uint32]uint32 {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	// 房间随机在线人数
	gameOnline, has := ls.onlineMaps[gameId]
	if has == false {
		gameOnline = make(map[uint32]uint32, 0)
	}
	ls.onlineMaps[gameId] = gameOnline

	for _, roomId := range roomIds {
		roomOnline := gameOnline[roomId]
		if roomOnline == 0 {
			roomOnline = uint32(rand.Intn(1000) + 100)
		}
		gameOnline[roomId] = roomOnline
	}
	return gameOnline
}

func (ls *LobbyServer) pauseUserGame(userId uint32) {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	user, has := ls.users[userId]
	if has == false {
		return
	}

	if user.Game == nil {
		return
	}

	msg := &net.Message{
		MsgId:   uint16(lobby_pb.MID_PAUSE_GAME),
		SvrType: ls.server.self.Type,
		SvrId:   0,
		UserId:  userId,
		Data: net.Marshal(&lobby_pb.PauseGame{
			UserId: userId,
		}),
	}

	data, err := msg.Encode()
	if err != nil {
		log.Sugar().Errorf("Encode msg:%d err:%v", msg.MsgId, err)
		return
	}

	err = user.Game.Send(data)
	if err != nil {
		log.Sugar().Errorf("Send msg:%d err:%v", msg.MsgId, err)
	}

}

func (ls *LobbyServer) continueUserGame(userId uint32) {
	defer ls.mutex.Unlock()
	ls.mutex.Lock()

	user, has := ls.users[userId]
	if has == false {
		return
	}

	if user.Game == nil {
		return
	}

	msg := &net.Message{
		MsgId:   uint16(lobby_pb.MID_CONTINUE_GAME),
		SvrType: ls.server.self.Type,
		SvrId:   0,
		UserId:  userId,
		Data: net.Marshal(&lobby_pb.ContinueGame{
			UserId: userId,
		}),
	}

	data, err := msg.Encode()
	if err != nil {
		log.Sugar().Errorf("Encode msg:%d err:%v", msg.MsgId, err)
		return
	}

	err = user.Game.Send(data)
	if err != nil {
		log.Sugar().Errorf("Send msg:%d err:%v", msg.MsgId, err)
	}

}
