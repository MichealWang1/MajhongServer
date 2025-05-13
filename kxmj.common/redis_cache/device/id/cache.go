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

	params := e.Data.(*kxmj_core.Device)
	if len(params.DeviceId) <= 0 {
		return
	}

	if e.Action == redis_core.InsertAction {
		err := c.Set(ctx, params.DeviceId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Set  params:%v err:%v", e, err))
		}
	} else if e.Action == redis_core.DeleteAction {
		err := c.Del(ctx, params.DeviceId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
		}
	}
}

func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.Device{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Exists(ctx context.Context, deviceId string) bool {
	key := c.generateKey()
	field := fmt.Sprintf("%s", deviceId)
	return c.client.HExists(ctx, key, field).Val()
}

func (c *Cache) Set(ctx context.Context, deviceId string) error {
	key := c.generateKey()
	field := fmt.Sprintf("%s", deviceId)
	return c.client.HSet(ctx, key, map[string]string{field: fmt.Sprintf("%d", 1)}).Err()
}

func (c *Cache) BulkSet(ctx context.Context, devices []*kxmj_core.Device) error {
	key := c.generateKey()
	maps := make(map[string]string, 0)
	for _, u := range devices {
		field := fmt.Sprintf("%s", u.DeviceId)
		maps[field] = fmt.Sprintf("%d", 1)
	}

	if len(maps) <= 0 {
		return nil
	}

	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) Del(ctx context.Context, deviceId string) error {
	key := c.generateKey()
	return c.client.HDel(ctx, key, fmt.Sprintf("%s", deviceId)).Err()
}

func (c *Cache) DelAll(ctx context.Context) error {
	key := c.generateKey()
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) generateKey() string {
	return keys.DeviceIdFormatKey
}
