package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"kxmj.common/log"
	"kxmj.common/mq"
	"kxmj.common/mysql"
	"kxmj.common/net"
	"kxmj.common/redis_cache/redis_core"
)

var Default *Config

type InnerConfig struct {
	Type uint16 `json:"type" yaml:"type"` // 服务器类型
	Id   uint16 `json:"id" yaml:"id"`     // 服务Id
}

type HttpServiceConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	HttpPort      int                     `yaml:"httpPort"`                           // web监听端口
	Services      []*net.ServerConfig     `yaml:"services"`                           // 内网分布式服务
	EtcdEndpoints []string                `json:"etcdEndpoints" yaml:"etcdEndpoints"` // Etcd 服务配置
	CoreMaster    *mysql.DbConfig         `yaml:"coreMaster"`                         // 核心业务主库
	CoreSlave     *mysql.DbConfig         `yaml:"coreSlave"`                          // 核心业务从库
	LoggerMaster  *mysql.DbConfig         `yaml:"loggerMaster"`                       // 日志主库
	LoggerSlave   *mysql.DbConfig         `yaml:"loggerSlave"`                        // 日志从库
	ReportMaster  *mysql.DbConfig         `yaml:"reportMaster"`                       // 业务报表主库(business owner)
	ReportSlave   *mysql.DbConfig         `yaml:"reportSlave"`                        // 业务报表从库(business owner)
	GameMaster    *mysql.DbConfig         `yaml:"gameMaster"`                         // 游戏记录主库
	GameSlave     *mysql.DbConfig         `yaml:"gameSlave"`                          // 游戏记录从库
	Redis         *redis_core.RedisConfig `yaml:"redis"`                              // redis配置配置
	MqConfig      *mq.RabbitmqConfig      `yaml:"mqConfig"`                           // MQ服务
	Logger        *log.Config             `yaml:"logger"`                             // 日志配置
	UseSwagger    string                  `yaml:"useSwagger"`                         // 是否启用swagger文档
}

func Create() *Config {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("./config")
	vp.AddConfigPath(".") // 添加搜索路径
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s", err))
	}

	config := &Config{}
	err = vp.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("配置文件解析失败: %s", err))
	}
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件变更")
		err := vp.Unmarshal(config)
		if err != nil {
			fmt.Println("配置文件更新失败")
		} else {
			fmt.Printf("配置文件更新成功: %+v \n", config)
		}
	})
	return config
}
