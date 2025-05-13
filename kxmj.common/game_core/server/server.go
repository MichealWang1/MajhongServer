package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/proto/lobby_pb"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/service"
	tcpNet "net"
	"time"
)

type RpcxServerConfig struct {
	Path string `json:"path" yaml:"path"` // rpcx注册服务路径
	Type uint16 `json:"type" yaml:"type"` // 服务器类型
	Id   uint16 `json:"id" yaml:"id"`     // 服务Id
}

type RpcxServer struct {
	self           *net.ServerConfig                 // 当前rpcx服务配置
	etcdEndpoints  []string                          // etcd服务IP端口配置
	lobbyConfig    *RpcxServerConfig                 // 大厅rpcx服务配置
	centerConfig   *RpcxServerConfig                 // 账号rpcx服务配置
	redis          *redis.Client                     // redis client
	lobby          client.XClient                    // rpxc client
	center         client.XClient                    // rpxc client
	clientsChan    map[string]chan *protocol.Message // 直连客户端接收消息管道
	server         *server.Server                    // RPCX服务
	clientReceiver func(msg net.MsgContext)          // RPCX客户端消息接收器
	serverReceiver func(ctx net.MsgContext)          // RPCX服务端消息接收器
}

func NewRpcxServer(self *net.ServerConfig, etcdEndpoints []string, lobbyConfig *RpcxServerConfig, centerConfig *RpcxServerConfig, redis *redis.Client) *RpcxServer {
	return &RpcxServer{
		self:          self,
		etcdEndpoints: etcdEndpoints,
		lobbyConfig:   lobbyConfig,
		centerConfig:  centerConfig,
		redis:         redis,
		clientsChan:   make(map[string]chan *protocol.Message, 0),
	}
}

func (rs *RpcxServer) Start(receiver func(ctx net.MsgContext)) {
	rs.serverReceiver = receiver

	// 注册服务
	err := rs.register()
	if err != nil {
		log.Sugar().Errorf("register service err:%v", err)
		panic(err)
	}

	// 注册服务到etcd
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + fmt.Sprintf("%s:%d", rs.self.Ip, rs.self.Port),
		EtcdServers:    rs.etcdEndpoints,
		BasePath:       rs.self.Path,
		UpdateInterval: time.Minute,
	}

	rs.server = server.NewServer()
	err = r.Start()
	if err != nil {
		panic(err)
	}
	rs.server.Plugins.Add(r)

	servicePath := service.ParseServicePath(rs.self.Type, rs.self.Id)
	err = rs.server.RegisterName(servicePath, rs, "")
	if err != nil {
		panic(err)
	}

	// 启动rpcx服务监听
	go func() {
		err = rs.server.Serve("tcp", fmt.Sprintf(":%d", rs.self.Port))
		if err != nil {
			panic(err)
		}
	}()

	// 创建大厅服务连接(长连接)
	rs.connectLobby(receiver)

	// 创建账号服务连接(短连接)
	rs.connectAccount()
}

func (rs *RpcxServer) Close() {
	err := rs.unregister()
	if err != nil {
		log.Sugar().Errorf("unregister err:%v", err)
	}

	err = rs.lobby.Call(context.Background(), "OnGateway", &net.Message{
		MsgId:   uint16(lobby_pb.MID_UNREGISTER),
		SvrType: rs.self.Type,
		SvrId:   rs.self.Id,
		UserId:  0,
		Data: net.Marshal(&lobby_pb.UnRegister{
			SvrType: uint32(rs.self.Type),
			SvrId:   uint32(rs.self.Id),
		}),
	}, nil)

	if err != nil {
		log.Sugar().Errorf("Call err:%v", err)
	}
}

func (rs *RpcxServer) OnGateway(ctx context.Context, args *net.Message, reply *net.RpcxReply) error {
	if rs.serverReceiver != nil {
		conn := ctx.Value(server.RemoteConnContextKey)
		mCtx := net.NewInnerServerContext(rs.server, conn.(tcpNet.Conn), args)
		mCtx.Session().SetSvrType(service.GatewayService)
		rs.serverReceiver(mCtx)
	}
	return nil
}

func (rs *RpcxServer) connectLobby(callBack func(msg net.MsgContext)) {
	rs.clientReceiver = callBack
	servicePath := service.ParseServicePath(rs.lobbyConfig.Type, rs.lobbyConfig.Id)
	d, err := etcdClient.NewEtcdV3Discovery(rs.lobbyConfig.Path, servicePath, rs.etcdEndpoints, true, nil)
	if err != nil {
		log.Sugar().Errorf("NewEtcdV3Discovery err:%v", err)
		return
	}

	ch := make(chan *protocol.Message, 4096)
	rs.lobby = client.NewBidirectionalXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	rs.clientsChan[servicePath] = ch

	// 注册服务
	go func() {
		for {
			err = rs.lobby.Call(context.Background(), "OnGateway", &net.Message{
				MsgId:   uint16(lobby_pb.MID_REGISTER),
				SvrType: rs.self.Type,
				SvrId:   rs.self.Id,
				UserId:  0,
				Data: net.Marshal(&lobby_pb.Register{
					SvrType: uint32(rs.self.Type),
					SvrId:   uint32(rs.self.Id),
					Addr:    rs.self.Ip,
					Port:    uint32(rs.self.Port),
				}),
			}, nil)

			if err != nil {
				log.Sugar().Errorf("Call err:%v", err)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for msg := range ch {
			data, err := net.Unpack(msg.Payload)
			if err != nil {
				log.Sugar().Errorf("Unpack err:%v", err)
				continue
			}

			if rs.clientReceiver != nil {
				ctx := net.NewInnerClientContext(servicePath, rs.lobby, data)
				ctx.Session().SetSvrType(rs.lobbyConfig.Type)
				rs.clientReceiver(ctx)
			}
		}
	}()
}

func (rs *RpcxServer) connectAccount() {
	servicePath := service.ParseServicePath(rs.centerConfig.Type, rs.centerConfig.Id)
	d, err := etcdClient.NewEtcdV3Discovery(rs.centerConfig.Path, servicePath, rs.etcdEndpoints, true, nil)
	if err != nil {
		log.Sugar().Errorf("NewEtcdV3Discovery err:%v", err)
		return
	}

	rs.center = client.NewXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
}

func (rs *RpcxServer) Lobby() client.XClient {
	return rs.lobby
}

func (rs *RpcxServer) Center() client.XClient {
	return rs.center
}

func (rs *RpcxServer) register() error {
	data, err := json.Marshal(rs.self)
	if err != nil {
		return err
	}

	key := keys.RpcxFormatKey
	field := service.ParseServicePathKey(rs.self.Path, rs.self.Type, rs.self.Id)
	ctx := context.Background()
	return rs.redis.HSet(ctx, key, field, string(data)).Err()
}

func (rs *RpcxServer) unregister() error {
	key := keys.RpcxFormatKey
	field := service.ParseServicePathKey(rs.self.Path, rs.self.Type, rs.self.Id)
	ctx := context.Background()
	return rs.redis.HDel(ctx, key, field).Err()
}
