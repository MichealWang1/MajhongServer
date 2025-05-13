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
	"time"
)

type Cache struct {
	client   *redis.Client
	local    map[string]*kxmj_core.Goods // 5s本地缓存
	lastTime int64
	duration int64
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
	return &kxmj_core.Goods{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client:   client,
		local:    make(map[string]*kxmj_core.Goods, 0),
		lastTime: 0,
		duration: 10 * 1000,
	}
}

func (c *Cache) GetAll(ctx context.Context) (map[string]*kxmj_core.Goods, error) {
	if time.Now().UnixMilli()-c.lastTime < c.duration {
		return c.local, nil
	}

	key := generateKey()
	value := make(map[string]*kxmj_core.Goods, 0)
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) <= 0 {
		err = c.SetAll(ctx)
		if err != nil {
			return nil, err
		}
		return c.GetAll(ctx)
	} else {
		for k, v := range data {
			val := &kxmj_core.Goods{}
			err = json.Unmarshal([]byte(v), val)
			if err != nil {
				return nil, err
			}
			value[k] = val
		}
	}
	return value, nil
}

func (c *Cache) SetAll(ctx context.Context) error {
	var queries []*kxmj_core.Goods
	err := mysql.CoreMaster().WithContext(ctx).Where("1 = 1").Find(&queries).Error
	if err != nil {
		return err
	}

	if len(queries) <= 0 {
		return errors.New("goods is null")
	}

	maps := make(map[string]string, 0)
	for _, goods := range queries {
		data, err := json.Marshal(goods)
		if err != nil {
			return err
		}

		field := goods.GoodsId
		maps[field] = string(data)
	}

	key := generateKey()
	err = c.client.HMSet(ctx, key, maps).Err()
	if err != nil {
		return err
	}

	c.lastTime = time.Now().UnixMilli()
	cache := make(map[string]*kxmj_core.Goods, 0)
	for _, item := range queries {
		cache[item.GoodsId] = item
	}
	c.local = cache
	return nil
}

func (c *Cache) DelAll(ctx context.Context) error {
	key := generateKey()
	return c.client.Del(ctx, key).Err()
}

func generateKey() string {
	return fmt.Sprintf(keys.GoodsFormatKey)
}
