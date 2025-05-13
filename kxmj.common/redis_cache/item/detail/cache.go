package detail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/redis_cache/redis_core"
	"strconv"
	"time"
)

type Cache struct {
	client   *redis.Client
	local    map[uint32]*kxmj_core.Item // 5s本地缓存
	lastTime int64                      // 最后一次同步时间
	duration int64                      // 本地缓存存活时间
}

func (c *Cache) EventHandler(ctx context.Context, e *redis_core.EventParams) {
	if e == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}

	err := c.DelAll(ctx)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
	}
}

func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.Item{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client:   client,
		local:    make(map[uint32]*kxmj_core.Item, 0),
		lastTime: 0,
	}
}

func (c *Cache) GetAll(ctx context.Context) (map[uint32]*kxmj_core.Item, error) {
	if time.Now().UnixMilli()-c.lastTime < 5000 {
		return c.local, nil
	}

	key := generateKey()
	value := make(map[uint32]*kxmj_core.Item, 0)
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) <= 0 {
		err = c.SetAll(ctx)
		if err != nil {
			return nil, err
		}
		return c.GetAll(ctx)
	} else {
		for k, v := range data {
			item := &kxmj_core.Item{}
			err = json.Unmarshal([]byte(v), item)
			if err != nil {
				return nil, err
			}

			itemId, err := strconv.Atoi(k)
			if err != nil {
				log.Sugar().Errorf("Atoi err:%v", err)
			}
			value[uint32(itemId)] = item
		}
	}
	return value, nil
}

func (c *Cache) SetAll(ctx context.Context) error {
	var queries []*kxmj_core.Item
	err := mysql.CoreMaster().WithContext(ctx).Where("1 = 1").Find(&queries).Error
	if err != nil {
		return err
	}

	if len(queries) <= 0 {
		return errors.New("item is null")
	}

	maps := make(map[string]string, 0)
	for _, item := range queries {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		field := fmt.Sprintf("%d", item.ItemId)
		maps[field] = string(data)
	}

	key := generateKey()
	err = c.client.HMSet(ctx, key, maps).Err()
	if err != nil {
		return err
	}

	c.lastTime = time.Now().UnixMilli()
	c.local = make(map[uint32]*kxmj_core.Item, 0)
	for _, item := range queries {
		c.local[item.ItemId] = item
	}
	return nil
}

func (c *Cache) DelAll(ctx context.Context) error {
	key := generateKey()
	return c.client.Del(ctx, key).Err()
}

func generateKey() string {
	return fmt.Sprintf(keys.ItemFormatKey)
}
