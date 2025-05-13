package gm

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/utils"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// CardStackData struct
type CardStackData struct {
	Banker    uint8    `json:"banker" redis:"userId"`       // 庄家/地主
	CatchCard uint32   `json:"catchCard" redis:"catchCard"` // 摸的牌
	Cards     []uint32 `json:"cards" redis:"cards"`         // 牌堆
	HandCards []uint32 `json:"handCards" redis:"handCards"` // 手牌
}

func (c *Cache) Get(ctx context.Context, gameType uint16, roomLevel uint8, userId uint32) (*CardStackData, error) {
	key := generateKey(gameType, roomLevel, userId)
	var value *CardStackData
	err := c.client.HGetAll(ctx, key).Scan(value)
	return value, err
}

func (c *Cache) Set(ctx context.Context, gameType uint16, roomLevel uint8, userId uint32, config *CardStackData) error {
	maps, err := utils.StructToMap(config)
	if err != nil {
		return err
	}

	key := generateKey(gameType, roomLevel, userId)
	_, err = c.client.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range maps {
			err = pipeline.HSet(ctx, key, k, v).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (c *Cache) Del(ctx context.Context, gameType uint16, roomLevel uint8, userId uint32) error {
	key := generateKey(gameType, roomLevel, userId)
	return c.client.Del(ctx, key).Err()
}

func generateKey(gameType uint16, roomLevel uint8, userId uint32) string {
	return fmt.Sprintf(keys.GMConfigurationFormatKey, gameType, roomLevel, userId)
}
