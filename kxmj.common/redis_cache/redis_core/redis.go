package redis_core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"kxmj.common/log"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`     //地址:端口
	Password string `yaml:"password"` //密码
}

var (
	client *redis.Client // db 0
)

func create(config *RedisConfig) *redis.Client {
	db := 0
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           db,
		PoolSize:     30,
		MinIdleConns: 30,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("初始化db:%d成功", db), zap.Any("db", db), zap.Any("err", err))
		return client
	}

	log.Sugar().Info(fmt.Sprintf("初始化db:%d成功", db), zap.Any("db", db), zap.Any("pong", pong))
	return client
}

func Init(config *RedisConfig) {
	client = create(config)
}

func Default() *redis.Client {
	return client
}
