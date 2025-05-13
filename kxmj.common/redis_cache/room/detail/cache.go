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
	"sort"
	"time"
)

type Cache struct {
	client   *redis.Client                    // redis client
	local    map[string]*kxmj_core.ConfigRoom // 本地缓存
	lastTime int64                            // 最后一次同步时间
	duration int64                            // 本地缓存存活时间
}

func (c *Cache) EventHandler(ctx context.Context, e *redis_core.EventParams) {
	if e == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}

	err := c.Del(ctx)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
	}
}

func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.ConfigRoom{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client:   client,
		local:    make(map[string]*kxmj_core.ConfigRoom, 0),
		lastTime: 0,
		duration: 10 * 1000,
	}
}

func (c *Cache) Get(ctx context.Context, gameId uint16, roomId uint32) (*kxmj_core.ConfigRoom, error) {
	if time.Now().UnixMilli()-c.lastTime < c.duration {
		for _, d := range c.local {
			if d.GameId == gameId && d.RoomId == roomId {
				return d, nil
			}
		}
	} else {
		temp, err := c.GetAll(ctx, gameId)
		if err != nil {
			return nil, err
		}
		for _, d := range temp {
			if d.GameId == gameId && d.RoomId == roomId {
				return d, nil
			}
		}
	}
	return nil, errors.New("data not exist")
}

func (c *Cache) GetAll(ctx context.Context, gameId uint16) ([]*kxmj_core.ConfigRoom, error) {
	var value []*kxmj_core.ConfigRoom
	var cache map[string]*kxmj_core.ConfigRoom
	if time.Now().UnixMilli()-c.lastTime < c.duration {
		cache = c.local
	} else {
		key := generateKey()
		data, err := c.client.HGetAll(ctx, key).Result()
		if err != nil || len(data) <= 0 {
			err = c.Set(ctx)
			if err != nil {
				return nil, err
			}
			return c.GetAll(ctx, gameId)
		}

		cache = make(map[string]*kxmj_core.ConfigRoom, 0)
		for _, d := range data {
			item := &kxmj_core.ConfigRoom{}
			err = json.Unmarshal([]byte(d), item)
			if err != nil {
				return nil, err
			}
			cache[generateField(item.GameId, item.RoomId)] = item
		}

		c.lastTime = time.Now().UnixMilli()
		c.local = cache
	}

	for _, d := range cache {
		if d.GameId == gameId {
			value = append(value, d)
		}
	}

	sort.Slice(value, func(i, j int) bool {
		return value[i].RoomId < value[j].RoomId
	})

	sort.Slice(value, func(i, j int) bool {
		return value[i].RoomLevel < value[j].RoomLevel
	})
	return value, nil
}

func (c *Cache) Set(ctx context.Context) error {
	var queries []*kxmj_core.ConfigRoom
	err := mysql.CoreMaster().WithContext(ctx).Where("1 = 1").Find(&queries).Error
	if err != nil {
		return err
	}

	cache := make(map[string]*kxmj_core.ConfigRoom)
	maps := make(map[string]string, 0)
	for _, room := range queries {
		data, err := json.Marshal(room)
		if err != nil {
			return err
		}

		field := generateField(room.GameId, room.RoomId)
		maps[field] = string(data)
		cache[field] = room
	}

	c.lastTime = time.Now().UnixMilli()
	c.local = cache
	key := generateKey()
	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) Del(ctx context.Context) error {
	key := generateKey()
	return c.client.Del(ctx, key).Err()
}

func generateField(gameId uint16, roomId uint32) string {
	return fmt.Sprintf("%d_%d", gameId, roomId)
}

func generateKey() string {
	return fmt.Sprintf(keys.RoomFormatKey)
}
