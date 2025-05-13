package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/service"
	"time"
)

type RpcxServer struct {
	self          *net.ServerConfig // 当前rpcx服务配置
	etcdEndpoints []string          // etcd服务IP端口配置
	redis         *redis.Client     // redis client
	server        *server.Server    // RPCX服务
}

func NewRpcxServer(self *net.ServerConfig, etcdEndpoints []string, redis *redis.Client) *RpcxServer {
	return &RpcxServer{
		self:          self,
		etcdEndpoints: etcdEndpoints,
		redis:         redis,
	}
}

func (rs *RpcxServer) Start() {
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
}

func (rs *RpcxServer) Close() {
	err := rs.unregister()
	if err != nil {
		log.Sugar().Errorf("unregister err:%v", err)
	}
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
