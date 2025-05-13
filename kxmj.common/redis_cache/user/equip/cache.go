package equip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/keys"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, userId uint32) (*kxmj_core.UserEquip, error) {
	key := generateKey()
	value := &kxmj_core.UserEquip{}
	field := fmt.Sprintf("%d", userId)
	err := c.client.HGet(ctx, key, field).Scan(value)
	if err != nil || value.Id <= 0 {
		err = mysql.CoreMaster().Where("user_id = ?", userId).First(value).Error
		if err != nil {
			return nil, err
		}

		err = c.Set(ctx, value)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func (c *Cache) Set(ctx context.Context, user *kxmj_core.UserEquip) error {
	key := generateKey()
	maps := make(map[string]string, 0)
	field := fmt.Sprintf("%d", user.UserId)
	val, _ := json.Marshal(user)
	maps[field] = string(val)
	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) BulkSet(ctx context.Context, users []*kxmj_core.UserEquip) error {
	key := generateKey()
	maps := make(map[string]string, 0)
	for _, u := range users {
		field := fmt.Sprintf("%d", u.UserId)
		val, _ := json.Marshal(u)
		maps[field] = string(val)
	}

	if len(maps) <= 0 {
		return nil
	}

	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) Del(ctx context.Context, userId uint32) error {
	key := generateKey()
	field := fmt.Sprintf("%d", userId)
	c.client.HDel(ctx, key, field)
	return nil
}

func generateKey() string {
	return fmt.Sprintf(keys.UserEquipFormatKey)
}
