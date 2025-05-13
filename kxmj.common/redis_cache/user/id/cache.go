package id

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/redis_cache/redis_core"
)

type Cache struct {
	client *redis.Client
}

func (c *Cache) EventHandler(ctx context.Context, e *redis_core.EventParams) {
	if e == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}

	if e.Data == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}

	params := e.Data.(*kxmj_core.User)
	if params.UserId <= 0 {
		return
	}

	if e.Action == redis_core.InsertAction {
		err := c.Set(ctx, params.UserId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Set  params:%v err:%v", e, err))
		}
	} else if e.Action == redis_core.DeleteAction {
		err := c.Del(ctx, params.UserId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
		}
	}
}

func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.User{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Exists(ctx context.Context, userId uint32) bool {
	key := c.generateKey()
	field := fmt.Sprintf("%d", userId)
	return c.client.HExists(ctx, key, field).Val()
}

func (c *Cache) Set(ctx context.Context, userId uint32) error {
	key := c.generateKey()
	field := fmt.Sprintf("%d", userId)
	return c.client.HSet(ctx, key, map[string]string{field: fmt.Sprintf("%d", 1)}).Err()
}

func (c *Cache) BulkSet(ctx context.Context, users []*kxmj_core.User) error {
	key := c.generateKey()
	maps := make(map[string]string, 0)
	for _, u := range users {
		field := fmt.Sprintf("%d", u.UserId)
		maps[field] = fmt.Sprintf("%d", 1)
	}

	if len(maps) <= 0 {
		return nil
	}

	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) Del(ctx context.Context, userId uint32) error {
	key := c.generateKey()
	c.client.HDel(ctx, key, fmt.Sprintf("%d", userId))
	return nil
}

func (c *Cache) DelAll(ctx context.Context) error {
	key := c.generateKey()
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) generateKey() string {
	return keys.UserIdFormatKey
}
