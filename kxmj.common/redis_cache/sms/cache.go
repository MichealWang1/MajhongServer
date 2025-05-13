package sms

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/keys"
	"time"
)

type Cache struct {
	client *redis.Client
	expire int
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
		expire: 3600 * 48,
	}
}

func (ls *Cache) Exist(ctx context.Context, telNumber string, formatType uint8) bool {
	key := ls.generateKey(telNumber, formatType)
	return ls.client.Exists(ctx, key).Val() > 0
}

func (ls *Cache) Set(ctx context.Context, telNumber string, formatType uint8, code string) (uint32, error) {
	key := ls.generateKey(telNumber, formatType)
	err := ls.client.Set(ctx, key, code, time.Second*150).Err()
	if err != nil {
		return 0, err
	}
	return 150, nil
}

func (ls *Cache) Get(ctx context.Context, telNumber string, formatType uint8) (string, error) {
	key := ls.generateKey(telNumber, formatType)
	val, err := ls.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (ls *Cache) generateKey(telNumber string, formatType uint8) string {
	prefix := keys.SmsFormatKey
	return fmt.Sprintf(prefix, telNumber, formatType)
}
