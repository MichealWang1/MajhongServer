package limit

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
		expire: 8,
	}
}

func (ls *Cache) SetHttpRequestKey(ctx context.Context, httpReqKey string) error {
	key := ls.generateRequestKey(httpReqKey)
	err := ls.client.Set(ctx, key, 1, time.Duration(ls.expire)*time.Second).Err()
	return err
}

func (ls *Cache) ExistRequestKey(ctx context.Context, httpReqKey string) bool {
	key := ls.generateRequestKey(httpReqKey)
	value, err := ls.client.Get(ctx, key).Int()
	if err != nil {
		return false
	}

	return value == 1
}

func (ls *Cache) DelHttpRequestKey(ctx context.Context, httpReqKey string) {
	key := ls.generateRequestKey(httpReqKey)
	ls.client.Del(ctx, key)
}

func (ls *Cache) generateRequestKey(httpReqKey string) string {
	return fmt.Sprintf(keys.HttpRequestLimitKey, httpReqKey)
}
