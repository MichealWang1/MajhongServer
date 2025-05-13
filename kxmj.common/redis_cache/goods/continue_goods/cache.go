package continue_goods

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/keys"
)

type Cache struct {
	client *redis.Client
}

type GoodsInfo struct {
	WitchDay uint32 `json:"witchDay" redis:"witchDay"` // 第几天领取物品
	Status   uint32 `json:"status" redis:"status"`     // 领取状态 0 未完成；1 已完成；2 已领取
	Date     uint32 `json:"date" redis:"date"`         // 可领取日期
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, userId uint32) (map[string][]*GoodsInfo, error) {
	key := generateKey(userId)
	maps, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*GoodsInfo, 0)
	for k, v := range maps {
		var ds []*GoodsInfo
		err := json.Unmarshal([]byte(v), &ds)
		if err != nil {
			return nil, err
		}
		result[k] = ds
	}
	return result, nil
}

func (c *Cache) Set(ctx context.Context, userId uint32, goodsList map[string][]*GoodsInfo) error {
	maps := make(map[string]string)
	for k, v := range goodsList {
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}
		maps[k] = string(d)
	}

	key := generateKey(userId)
	return c.client.HMSet(ctx, key, maps).Err()
}

func generateKey(userId uint32) string {
	return fmt.Sprintf(keys.GoodsContinueFormatKey, userId)
}
