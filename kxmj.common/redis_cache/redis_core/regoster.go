package redis_core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Register struct {
	client     *redis.Client
	lastTime   int64
	localCache map[string]string
}

func NewRegister(client *redis.Client) *Register {
	return &Register{
		client:     client,
		lastTime:   time.Now().UnixMilli(),
		localCache: make(map[string]string, 0),
	}
}

func (c *Register) GetAll(ctx context.Context) (map[string]string, error) {
	if time.Now().UnixMilli()-c.lastTime < 1000 {
		return c.localCache, nil
	}

	c.lastTime = time.Now().UnixMilli()
	key := c.generateKey()
	result, err := c.client.HGetAll(ctx, key).Result()
	if err == nil {
		c.localCache = result
		return c.localCache, nil
	}
	return nil, err
}

func (c *Register) Set(ctx context.Context, schema string, table string) error {
	key := c.generateKey()
	field := fmt.Sprintf("%s:%s", schema, table)
	return c.client.HSet(ctx, key, map[string]string{field: table}).Err()
}

func (c *Register) generateKey() string {
	return fmt.Sprintf(RegisterEventFormatKey)
}
