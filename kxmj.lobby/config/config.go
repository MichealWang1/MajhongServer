package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/redis_cache/redis_core"
)

var Default *Config

type Config struct {
	Self          *net.ServerConfig       `yaml:"self"`                               // 本地服务配置
	EtcdEndpoints []string                `json:"etcdEndpoints" yaml:"etcdEndpoints"` // Etcd 服务配置
	Redis         *redis_core.RedisConfig `yaml:"redis"`                              // redis配置配置
	Logger        *log.Config             `yaml:"logger"`                             // 日志配置
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
