package rpc

import (
	"github.com/go-redis/redis/v8"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/service"
)

type RpcxServer struct {
	servicesCfg   []*net.ServerConfig // 当前rpcx服务配置
	etcdEndpoints []string            // etcd服务IP端口配置
	lobby         client.XClient      // 大厅客户端
	center        client.XClient      // 账号服客户端
	redis         *redis.Client       // redis client
}

var temp *RpcxServer

func Default() *RpcxServer {
	return temp
}

func Init(servicesCfg []*net.ServerConfig, etcdEndpoints []string, redis *redis.Client) {
	temp = &RpcxServer{
		servicesCfg:   servicesCfg,
		etcdEndpoints: etcdEndpoints,
		redis:         redis,
	}
}

func (rs *RpcxServer) Start() {
	for _, cfg := range rs.servicesCfg {
		servicePath := service.ParseServicePath(cfg.Type, cfg.Id)
		d, err := etcdClient.NewEtcdV3Discovery(cfg.Path, servicePath, rs.etcdEndpoints, true, nil)
		if err != nil {
			log.Sugar().Errorf("NewEtcdV3Discovery err:%v", err)
			return
		}

		if cfg.Type == service.CenterService {
			rs.center = client.NewXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		} else if cfg.Type == service.LobbyService {
			rs.lobby = client.NewXClient(servicePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		}
	}
}

func (rs *RpcxServer) LobbyClient() client.XClient {
	return rs.lobby
}

func (rs *RpcxServer) CenterClient() client.XClient {
	return rs.center
}
