package game_core

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/net"
	"kxmj.common/proto/gateway_pb"
	"kxmj.common/proto/lobby_pb"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/gm"
	"kxmj.common/service"
	"kxmj.common/utils"
	"math/big"
	"sort"
	"time"
)

type Game struct {
	server        IServer                                 // 服务接口
	msgChan       chan net.MsgContext                     // 接收消息管道
	sendChan      chan *net.Message                       // 发送消息管道
	enterChan     chan *EnterDesk                         // 进入游戏回调管道
	leaveChan     chan *LeaveDesk                         // 离开游戏回调管道
	deskCloseChan chan IDesk                              // 桌子关闭事件管道
	desks         map[uint32]map[uint32]IDesk             // 房间桌子信息
	userDesks     map[uint32]IDesk                        // 用户桌子信息
	deskUsers     map[uint32]map[uint32]map[uint32]uint32 // 房间桌子用户信息
	rooms         map[uint32]IRoom                        // 房间列表
	userGateways  map[uint32]net.Session                  // 网关通讯Session
	closeChan     chan struct{}                           // 游戏服务关闭事件管道
	runChan       chan struct{}                           // 游戏运行时消息管道
	closed        bool                                    // 游戏服务关闭标志
	updateChan    chan []*center.RoomConfig               // 更新房间信息管道
}

func NewGame(server IServer) *Game {
	return &Game{
		server:        server,
		msgChan:       make(chan net.MsgContext, 4096),
		sendChan:      make(chan *net.Message, 4096),
		enterChan:     make(chan *EnterDesk, 4096),
		leaveChan:     make(chan *LeaveDesk, 4096),
		deskCloseChan: make(chan IDesk, 1024),
		desks:         make(map[uint32]map[uint32]IDesk, 0),
		userDesks:     make(map[uint32]IDesk, 0),
		deskUsers:     make(map[uint32]map[uint32]map[uint32]uint32, 0),
		rooms:         make(map[uint32]IRoom, 0),
		userGateways:  make(map[uint32]net.Session, 0),
		closeChan:     make(chan struct{}, 0),
		runChan:       make(chan struct{}, 100),
		closed:        false,
		updateChan:    make(chan []*center.RoomConfig, 10),
	}
}

func (g *Game) Start() {
	resp, err := g.Server().GetRoomConfigList(g.Server().SvrType())
	if err != nil {
		panic(err)
	}

	for _, config := range resp.Data {
		g.rooms[config.RoomId] = NewRoom(config)
	}

	// 5秒同步一次房间配置
	go func() {
		for {
			if g.closed == true {
				break
			}

			time.Sleep(time.Second * 5)
			resp, err = g.Server().GetRoomConfigList(g.Server().SvrType())
			if err != nil {
				continue
			}
			g.updateChan <- resp.Data
		}
	}()

	go func() {
		for {
			select {
			case <-g.closeChan:
				return
			case <-g.runChan:
				g.run()
			case ctx := <-g.msgChan:
				g.handler(ctx)
			case msg := <-g.sendChan:
				g.sendMessage(msg)
			case enterInfo := <-g.enterChan:
				g.enter(enterInfo)
			case leaveInfo := <-g.leaveChan:
				g.leave(leaveInfo)
			case desk := <-g.deskCloseChan:
				g.deskClose(desk)
			case configs := <-g.updateChan:
				g.updateRooms(configs)
			}
		}
	}()

	go func() {
		for {
			g.runChan <- struct{}{}
			time.Sleep(time.Millisecond * 200)
		}
	}()
}

func (g *Game) Close() {
	if g.closed {
		return
	}

	g.closed = true
}

func (g *Game) WaitClose() {
	<-g.closeChan
}

func (g *Game) Server() IServer {
	return g.server
}

func (g *Game) Context(ctx net.MsgContext) {
	g.msgChan <- ctx
}

func (g *Game) OnEnter(enterInfo *EnterDesk) {
	g.enterChan <- enterInfo

	if enterInfo.Player.IsRobot {
		return
	}

	err := g.Server().GetLobby().Call(context.Background(), "OnGateway", &net.Message{
		MsgId:   uint16(lobby_pb.MID_ENTER_DESK),
		SvrType: g.Server().SvrType(),
		SvrId:   g.Server().SvrId(),
		UserId:  enterInfo.Player.UserId,
		Data: net.Marshal(&lobby_pb.EnterDesk{
			UserId:  enterInfo.Player.UserId,
			SvrType: uint32(g.Server().SvrType()),
			SvrId:   uint32(g.Server().SvrId()),
			RoomId:  enterInfo.Desk.Room().ID(),
			DeskId:  enterInfo.Desk.ID(),
		}),
	}, nil)

	if err != nil {
		log.Sugar().Errorf("Calll err:%v", err)
	}
}

func (g *Game) OnLeave(leaveInfo *LeaveDesk) {
	goldInt, ok := new(big.Int).SetString(leaveInfo.Gold, 10)
	zeroInt := new(big.Int).SetInt64(0)

	if ok && goldInt.Cmp(zeroInt) > 0 && leaveInfo.IsRobot == false {
		err := g.Server().SetUserGold(leaveInfo.UserId, leaveInfo.Gold, leaveInfo.Desk.Room().ID(), leaveInfo.Desk.Room().RoomLevel())
		if err != nil {
			log.Sugar().Errorf("SetUserGold data:%v err:%v", leaveInfo, err)
			return
		}
	}

	g.leaveChan <- leaveInfo

	if leaveInfo.IsRobot {
		return
	}

	if len(leaveInfo.Statistics) <= 0 {
		return
	}

	var userStatistics []*UserStatistics
	for _, d := range leaveInfo.Statistics {
		userStatistics = append(userStatistics, &UserStatistics{
			UserId:        leaveInfo.UserId,
			RoomId:        leaveInfo.Desk.Room().ID(),
			GameId:        leaveInfo.Desk.Room().GameId(),
			GameType:      leaveInfo.Desk.Room().GameType(),
			RoomLevel:     leaveInfo.Desk.Room().RoomLevel(),
			PlayType:      d.PlayType,
			TotalTimes:    d.TotalTimes,
			TotalWinLoss:  d.TotalWinLoss,
			TotalDuration: d.TotalDuration,
		})
	}

	g.Server().UpdateStatistics(userStatistics)
}

func (g *Game) OnDeskClose(desk IDesk) {
	g.deskCloseChan <- desk
}

func (g *Game) SendMessage(userId uint32, msgId uint16, data proto.Message) {
	if userId < uint32(MaxRobotId) {
		return
	}

	g.sendChan <- g.parseMessage(userId, msgId, data)
}

func (g *Game) SendErrMessage(ctx net.MsgContext, code int) {
	if ctx.Request().UserId < uint32(MaxRobotId) {
		return
	}

	msg := g.parseMessage(ctx.Request().UserId, uint16(gateway_pb.MID_ERR), &gateway_pb.Err{
		Code:        uint32(code),
		Msg:         codes.GetMessage(code),
		OriginMsgId: uint32(ctx.Request().MsgId),
	})

	g.send(ctx.Session(), msg)
}

func (g *Game) GetRobot(baseScore string) *RobotInfo {
	return CreateRobot(baseScore)
}

func (g *Game) NotifyRise(userId uint32, desk IDesk) {
	// 获取用户钱包钻石数
	diamond, err := g.Server().CheckUserDiamond(userId)
	if err != nil {
		log.Sugar().Errorf("CheckUserDiamond user:%d err:%v", userId, err)
		return
	}

	// 房间级别
	roomLevel := desk.Room().RoomLevel()

	// 获取商品配置
	goodsConfig, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(context.Background())
	if err != nil {
		log.Sugar().Errorf("Get goods config err:%v", err)
		return
	}

	// 获取物品配置
	itemConfig, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(context.Background())

	// 复活卡购买逻辑
	// 1.先找到对应房间级别复活卡物品
	// 2.通过服务卡物品找到对应的复活卡商品
	// 3.把复活卡商品发送给前端
	// 4.前端显示复活卡购买界面，前端调用商城购买接口购买
	// 5.购买成功通知大厅服务转发，大厅服务判断购买的商品是否是复活卡，如果是转发给用户所在游戏桌子，并同时通知前端用户购买商品成功
	// 6.游戏收到用户购买复活卡成功消息，同步用户钱包金币到游戏并通知客户端刷新用户余额
	var goodsList []*RiseGoods
	for _, d := range itemConfig {
		if item.RisePack != item.Type(d.ItemType) {
			continue
		}

		v := item.GetValueItem(d)
		if len(v.Extra) <= 0 {
			continue
		}

		// 复活卡扩展字段说明{"1":2} key:房间级别, value:复活卡级别
		level, has := v.Extra[uint32(roomLevel)]
		if has {
			var goods *kxmj_core.Goods
			// 优先找到钻石购买商品
			for _, gVal := range goodsConfig {
				if d.ItemId == gVal.ItemId {
					if utils.Cmp(diamond, gVal.Price) > 0 && gVal.ShopType == 2 {
						goods = gVal
						break
					}
				}
			}

			// 如果没有找到钻石购买商品，那么找RMB购买商品
			if goods == nil {
				for _, gVal := range goodsConfig {
					if d.ItemId == gVal.ItemId {
						if gVal.ShopType == 1 {
							goods = gVal
							break
						}
					}
				}
			}

			goodsList = append(goodsList, &RiseGoods{
				GoodsId:       goods.GoodsId,
				Price:         goods.Price,
				RealCount:     goods.RealCount,
				OriginalCount: goods.OriginalCount,
				ShopType:      goods.ShopType,
				RiseLevel:     level,
			})
		}
	}

	sort.Slice(goodsList, func(i, j int) bool {
		return goodsList[i].RiseLevel < goodsList[j].RiseLevel
	})

	data := &gateway_pb.NotifyRise{}
	for _, goods := range goodsList {
		data.List = append(data.List, &gateway_pb.RiseGoods{
			GoodsId:       goods.GoodsId,
			Price:         goods.Price,
			RealCount:     goods.RealCount,
			OriginalCount: goods.OriginalCount,
			ShopType:      uint32(goods.ShopType),
			RiseLevel:     goods.RiseLevel,
		})
	}
	g.SendMessage(userId, uint16(gateway_pb.MID_NOTIFY_RISE), data)
}

func (g *Game) GetManualConfig(userId uint32, desk IDesk) (*gm.CardStackData, error) {
	return redis_cache.GetCache().GetGameManageCache().Get(context.Background(), g.Server().SvrType(), desk.Room().RoomLevel(), userId)
}

func (g *Game) parseMessage(userId uint32, msgId uint16, data proto.Message) *net.Message {
	return &net.Message{
		MsgId:   msgId,
		SvrType: g.Server().SvrType(),
		SvrId:   g.Server().SvrId(),
		UserId:  userId,
		Data:    net.Marshal(data),
	}
}

func (g *Game) sendMessage(msg *net.Message) {
	if msg.UserId < uint32(MaxRobotId) {
		return
	}

	session, has := g.userGateways[msg.UserId]
	if has == false {
		log.Sugar().Errorf("User:%d session not found", msg.UserId)
		return
	}
	g.send(session, msg)
}

func (g *Game) send(session net.Session, msg *net.Message) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%d:U:%d] --->", msg.MsgId, msg.SvrType, msg.SvrId, msg.UserId))

	d, err := msg.Encode()
	if err != nil {
		log.Sugar().Errorf("Encode msg:%v err:%v", msg, err)
	}

	err = session.Send(d)
	if err != nil {
		log.Sugar().Errorf("SendMessage msg:%v err:%v", msg, err)
	}
}

func (g *Game) handler(ctx net.MsgContext) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%d:U:%d] <---", ctx.Request().MsgId, ctx.Request().SvrType, ctx.Request().SvrId, ctx.Request().UserId))

	// 绑定用户网关
	g.bindGateway(ctx)

	switch ctx.Request().MsgId {
	case uint16(gateway_pb.MID_MATCH):
		g.match(ctx)
	default:
		desk, has := g.userDesks[ctx.Request().UserId]
		if has == false {
			g.send(ctx.Session(), &net.Message{
				MsgId:   uint16(gateway_pb.MID_ERR),
				SvrType: ctx.Request().SvrType,
				SvrId:   ctx.Request().SvrId,
				UserId:  ctx.Request().UserId,
				Data: net.Marshal(&gateway_pb.Err{
					Code:        codes.GameDeskNotExist,
					Msg:         codes.GetMessage(codes.GameDeskNotExist),
					OriginMsgId: uint32(ctx.Request().MsgId),
				}),
			})
			return
		}
		desk.OnMessage(ctx)
	}
}

func (g *Game) match(ctx net.MsgContext) {
	if g.closed {
		g.SendErrMessage(ctx, codes.GameServerClose)
		return
	}

	payload := &gateway_pb.MatchReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		g.SendErrMessage(ctx, codes.MarshalPbErr)
		return
	}

	desk, has := g.userDesks[payload.UserId]
	if has {
		// 如果用户在游戏中，立即返回
		g.send(ctx.Session(), &net.Message{
			MsgId:   ctx.Request().MsgId,
			SvrType: ctx.Request().SvrType,
			SvrId:   ctx.Request().SvrId,
			UserId:  ctx.Request().UserId,
			Data: net.Marshal(&gateway_pb.MatchResp{
				UserId:  ctx.Request().UserId,
				SvrType: uint32(ctx.Request().SvrType),
				SvrId:   uint32(ctx.Request().SvrId),
				RoomId:  desk.Room().ID(),
				DeskId:  desk.ID(),
			}),
		})
		return
	}

	desks, has := g.desks[payload.RoomId]
	room := g.rooms[payload.RoomId]

	// 匹配逻辑
	userDesks, has := g.deskUsers[payload.RoomId]
	if has {
		for k, v := range userDesks {
			if len(v) < 4 {
				if has {
					d, has := desks[k]
					if has {
						desk = d
						break
					}
				}
			}
		}
	}

	// 如果没有匹配到新建一张桌子
	if desk == nil {
		desk = g.server.Template().New(g, room)
	}

	if has == false {
		desks = make(map[uint32]IDesk)
	}

	g.desks[room.ID()] = desks
	desks[desk.ID()] = desk
	g.userDesks[payload.UserId] = desk

	g.setDeskUser(room.ID(), desk.ID(), payload.UserId)

	desk.Start()
	desk.OnMessage(ctx)
}

func (g *Game) bindGateway(ctx net.MsgContext) {
	if ctx.Session().SvrType() == service.GatewayService && ctx.Request().UserId > 0 {
		g.userGateways[ctx.Request().UserId] = ctx.Session()
	}
}

func (g *Game) setDeskUser(roomId uint32, deskId uint32, userId uint32) {
	roomDeskUsers, has := g.deskUsers[roomId]
	if has == false {
		roomDeskUsers = make(map[uint32]map[uint32]uint32, 0)
	}
	g.deskUsers[roomId] = roomDeskUsers

	deskUsers, has := roomDeskUsers[deskId]
	if has == false {
		deskUsers = make(map[uint32]uint32, 0)
	}
	roomDeskUsers[deskId] = deskUsers
	deskUsers[userId] = userId
}

func (g *Game) delDeskUser(roomId uint32, deskId uint32, userId uint32) {
	roomDeskUsers, has := g.deskUsers[roomId]
	if has == false {
		return
	}

	deskUsers, has := roomDeskUsers[deskId]
	if has == false {
		return
	}

	delete(deskUsers, userId)
	if len(deskUsers) <= 0 {
		delete(roomDeskUsers, deskId)
	}
}

func (g *Game) getDeskUsers(roomId uint32, deskId uint32) []uint32 {
	roomDeskUsers, has := g.deskUsers[roomId]
	if has == false {
		return nil
	}

	deskUsers, has := roomDeskUsers[deskId]
	if has == false {
		return nil
	}

	var list []uint32
	for _, uId := range deskUsers {
		list = append(list, uId)
	}
	return list
}

func (g *Game) enter(info *EnterDesk) {
	g.userDesks[info.Player.UserId] = info.Desk

	desks, has := g.desks[info.Desk.Room().ID()]
	if has == false {
		desks = make(map[uint32]IDesk, 0)
	}
	g.desks[info.Desk.Room().ID()] = desks
	desks[info.Desk.ID()] = info.Desk

	g.setDeskUser(info.Desk.Room().ID(), info.Desk.ID(), info.Player.UserId)

	msg := g.parseMessage(info.Player.UserId, uint16(gateway_pb.MID_MATCH), &gateway_pb.MatchResp{
		UserId:  info.Player.UserId,
		SvrType: uint32(g.Server().SvrType()),
		SvrId:   uint32(g.Server().SvrId()),
		RoomId:  info.Desk.Room().ID(),
		DeskId:  info.Desk.ID(),
	})
	g.sendMessage(msg)
}

func (g *Game) leave(info *LeaveDesk) {
	delete(g.userDesks, info.UserId)
	desks, has := g.desks[info.Desk.Room().ID()]
	if has {
		delete(desks, info.Desk.ID())
	}

	users := g.getDeskUsers(info.Desk.Room().ID(), info.Desk.ID())
	for _, userId := range users {
		if userId < uint32(MaxRobotId) {
			continue
		}

		if userId != info.UserId && info.SeatId <= 0 {
			continue
		}

		msg := g.parseMessage(userId, uint16(gateway_pb.MID_LEAVE_DESK), &gateway_pb.LeaveDesk{
			SvrType: uint32(g.Server().SvrType()),
			SvrId:   uint32(g.Server().SvrId()),
			RoomId:  info.Desk.Room().ID(),
			DeskId:  info.Desk.ID(),
			UserId:  info.UserId,
		})
		g.sendMessage(msg)
	}

	g.delDeskUser(info.Desk.Room().ID(), info.Desk.ID(), info.UserId)
}

func (g *Game) deskClose(desk IDesk) {
	desk.Room().RemoveDeskId(desk.ID())

	desks, has := g.desks[desk.Room().ID()]
	if has {
		delete(desks, desk.ID())
	}

	userIds := g.getDeskUsers(desk.Room().ID(), desk.ID())
	for _, userId := range userIds {
		delete(g.userDesks, userId)
		g.delDeskUser(desk.Room().ID(), desk.ID(), userId)
	}
}

func (g *Game) updateRooms(configs []*center.RoomConfig) {
	for _, config := range configs {
		room, has := g.rooms[config.RoomId]
		if has == false {
			g.rooms[config.RoomId] = NewRoom(config)
			continue
		}
		room.Update(config)
	}
}

func (g *Game) run() {
	g.runClose()
	g.runDesks()
}

func (g *Game) runClose() {
	if g.closed {
		for _, v := range g.desks {
			if len(v) > 0 {
				return
			}

			for _, d := range v {
				d.Close()
			}
		}

		for _, v := range g.rooms {
			v.Close()
		}
		close(g.closeChan)
	}
}

func (g *Game) runDesks() {
	for _, desks := range g.desks {
		for _, desk := range desks {
			desk.Run()
		}
	}
}
