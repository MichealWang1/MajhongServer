package server

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/proto/gateway_pb"
	"kxmj.common/proto/lobby_pb"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/service"
	"kxmj.gateway/config"
	"time"
)

type Inner struct {
	self           *config.InnerConfig               // 当前rpcx服务配置
	etcdEndpoints  []string                          // etcd服务IP端口配置
	redis          *redis.Client                     // redis client
	clients        map[string]client.XClient         // rpxc client
	clientsChan    map[string]chan *protocol.Message // 直连客户端接收消息管道
	lobbyClient    client.XClient                    // 大厅服务客户端
	inputChan      chan *net.Message                 // 内网接收消息管道
	outputChan     chan *net.Message                 // 内网发送消息管道
	addXClientChan chan *net.ServerConfig            // 新增客户端管道
	closeChan      chan struct{}                     // 关闭管道
	gateway        *Gateway                          // 网关管理类
}

func NewInner(self *config.InnerConfig, etcdEndpoints []string, redis *redis.Client, gateway *Gateway) *Inner {
	return &Inner{
		self:           self,
		etcdEndpoints:  etcdEndpoints,
		redis:          redis,
		clients:        make(map[string]client.XClient, 0),
		clientsChan:    make(map[string]chan *protocol.Message, 0),
		inputChan:      make(chan *net.Message, 4096),
		outputChan:     make(chan *net.Message, 4096),
		addXClientChan: make(chan *net.ServerConfig, 100),
		closeChan:      make(chan struct{}, 1),
		gateway:        gateway,
	}
}

func (i *Inner) Start() {
	// 事件处理
	go func() {
		for {
			select {
			case <-i.closeChan:
				return
			case cfg := <-i.addXClientChan:
				i.addRpxcClient(cfg)
			case msg := <-i.inputChan:
				i.fromInner(msg)
			case msg := <-i.outputChan:
				i.toInner(msg)
			}
		}
	}()

	// 自动长连接
	go func() {
		for {
			registers, err := i.getRegisters()
			if err != nil {
				log.Sugar().Errorf("getRegisters err:%v", err)
			}

			for _, cfg := range registers {
				i.addXClientChan <- cfg
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func (i *Inner) Close() {
	close(i.closeChan)

	// 取消注册
	err := i.lobbyClient.Call(context.Background(), "OnGateway", &net.Message{
		MsgId:   uint16(lobby_pb.MID_UNREGISTER),
		SvrType: service.LobbyService,
		SvrId:   0,
		UserId:  0,
		Data: net.Marshal(&lobby_pb.UnRegister{
			SvrType: uint32(i.self.Type),
			SvrId:   uint32(i.self.Id),
		}),
	}, nil)

	if err != nil {
		log.Sugar().Errorf("Call err:%v", err)
	}
}

func (i *Inner) ToInner(msg *net.Message) {
	i.outputChan <- msg
}

func (i *Inner) getRegisters() ([]*net.ServerConfig, error) {
	key := keys.RpcxFormatKey
	ctx := context.Background()
	maps, err := i.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*net.ServerConfig, 0)
	for _, d := range maps {
		data := &net.ServerConfig{}
		err = json.Unmarshal([]byte(d), data)
		if err != nil {
			log.Sugar().Errorf("Decode err:%v", err)
			continue
		}
		result = append(result, data)
	}
	return result, nil
}

func (i *Inner) addRpxcClient(config *net.ServerConfig) {
	if i.self.Type == config.Type {
		return
	}

	servicePath := service.ParseServicePath(config.Type, config.Id)
	_, has := i.clients[servicePath]
	if has == true {
		return
	}

	d, err := etcdClient.NewEtcdV3Discovery(config.Path, servicePath, i.etcdEndpoints, true, nil)
	if err != nil {
		log.Sugar().Errorf("NewEtcdV3Discovery err:%v", err)
		return
	}

	ch := make(chan *protocol.Message, 4096)
	cli := client.NewBidirectionalXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	i.clients[servicePath] = cli
	i.clientsChan[servicePath] = ch

	// 大厅服务只有一个
	if config.Type == service.LobbyService {
		i.lobbyClient = cli

		// 注册网关
		go func() {
			for {
				err = i.lobbyClient.Call(context.Background(), "OnGateway", &net.Message{
					MsgId:   uint16(lobby_pb.MID_REGISTER),
					SvrType: service.LobbyService,
					SvrId:   0,
					UserId:  0,
					Data: net.Marshal(&lobby_pb.Register{
						SvrType: uint32(i.self.Type),
						SvrId:   uint32(i.self.Id),
						Addr:    i.gateway.Outer.config.Addr,
						Port:    uint32(i.gateway.Outer.config.Port),
					}),
				}, nil)

				if err != nil {
					log.Sugar().Errorf("Calll err:%v", err)
				}

				time.Sleep(time.Second * 2)
			}
		}()
	}

	go func() {
		for {
			select {
			case <-i.closeChan:
				return
			case msg := <-ch:
				data, err := net.Unpack(msg.Payload)
				if err != nil {
					log.Sugar().Errorf("Unpack err:%v", err)
					continue
				}
				i.inputChan <- data
			}
		}
	}()
}

func (i *Inner) fromInner(msg *net.Message) {
	if msg.MsgId == uint16(gateway_pb.MID_MATCH) {
		payload := &gateway_pb.MatchResp{}
		err := msg.Decode(payload)
		if err != nil {
			i.gateway.ToOuter(&net.Message{
				MsgId:   uint16(gateway_pb.MID_ERR),
				SvrType: msg.SvrType,
				SvrId:   msg.SvrId,
				UserId:  msg.UserId,
				Data: net.Marshal(&gateway_pb.Err{
					Code:        uint32(codes.UnMarshalPbErr),
					Msg:         codes.GetMessage(codes.UnMarshalPbErr),
					OriginMsgId: uint32(msg.MsgId),
				}),
			})
			return
		}

		// 如果没有匹配到可用服务，直接返回错误
		if payload.SvrId <= 0 {
			i.gateway.ToOuter(&net.Message{
				MsgId:   uint16(gateway_pb.MID_ERR),
				SvrType: msg.SvrType,
				SvrId:   msg.SvrId,
				UserId:  msg.UserId,
				Data: net.Marshal(&gateway_pb.Err{
					Code:        uint32(codes.SvrMatchFail),
					Msg:         codes.GetMessage(codes.SvrMatchFail),
					OriginMsgId: uint32(msg.MsgId),
				}),
			})
			return
		}

		// 服务已找到，开始匹配房间
		if payload.DeskId <= 0 {
			servicePath := service.ParseServicePath(msg.SvrType, msg.SvrId)
			cli, has := i.clients[servicePath]
			if has == false {
				// 如果服务未找到，给客户端返回错误消息
				i.gateway.ToOuter(&net.Message{
					MsgId:   uint16(gateway_pb.MID_ERR),
					SvrType: msg.SvrType,
					SvrId:   msg.SvrId,
					UserId:  msg.UserId,
					Data: net.Marshal(&gateway_pb.Err{
						Code:        uint32(codes.ServerNotFound),
						Msg:         codes.GetMessage(codes.ServerNotFound),
						OriginMsgId: uint32(msg.MsgId),
					}),
				})
				return
			}

			// 转发给内网游戏服务
			err = cli.Call(context.Background(), "OnGateway", &net.Message{
				MsgId:   uint16(gateway_pb.MID_MATCH),
				SvrType: uint16(payload.SvrType),
				SvrId:   uint16(payload.SvrId),
				UserId:  msg.UserId,
				Data: net.Marshal(&gateway_pb.MatchReq{
					UserId:  msg.UserId,
					SvrType: payload.SvrType,
					RoomId:  payload.RoomId,
				}),
			}, nil)

			if err != nil {
				i.gateway.ToOuter(&net.Message{
					MsgId:   uint16(gateway_pb.MID_ERR),
					SvrType: msg.SvrType,
					SvrId:   msg.SvrId,
					UserId:  msg.UserId,
					Data: net.Marshal(&gateway_pb.Err{
						Code:        uint32(codes.ServerNetErr),
						Msg:         codes.GetMessage(codes.ServerNetErr),
						OriginMsgId: uint32(msg.MsgId),
					}),
				})

				log.Sugar().Errorf("Call err:%v", err)
			}
			return
		}
	}

	i.gateway.ToOuter(msg)
}

func (i *Inner) toInner(msg *net.Message) {
	if msg.MsgId == uint16(gateway_pb.MID_MATCH) || msg.SvrType == service.LobbyService {
		err := i.lobbyClient.Call(context.Background(), "OnGateway", msg, nil)
		if err != nil {
			log.Sugar().Errorf("Call err:%v", err)
		}
		return
	}

	servicePath := service.ParseServicePath(msg.SvrType, msg.SvrId)
	cli, has := i.clients[servicePath]
	if has == false {
		// 如果服务未找到，给客户端返回错误消息
		i.gateway.ToOuter(&net.Message{
			MsgId:   uint16(gateway_pb.MID_ERR),
			SvrType: msg.SvrType,
			SvrId:   msg.SvrId,
			UserId:  msg.UserId,
			Data: net.Marshal(&gateway_pb.Err{
				Code:        uint32(codes.ServerNotFound),
				Msg:         codes.GetMessage(codes.ServerNotFound),
				OriginMsgId: uint32(msg.MsgId),
			}),
		})
		return
	}

	// 转发给内网游戏服务
	err := cli.Call(context.Background(), "OnGateway", msg, nil)
	if err != nil {
		i.gateway.ToOuter(&net.Message{
			MsgId:   uint16(gateway_pb.MID_ERR),
			SvrType: msg.SvrType,
			SvrId:   msg.SvrId,
			UserId:  msg.UserId,
			Data: net.Marshal(&gateway_pb.Err{
				Code:        uint32(codes.ServerNetErr),
				Msg:         codes.GetMessage(codes.ServerNetErr),
				OriginMsgId: uint32(msg.MsgId),
			}),
		})

		log.Sugar().Errorf("Call err:%v", err)
	}
}
