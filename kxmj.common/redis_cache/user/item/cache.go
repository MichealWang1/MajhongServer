package item

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/redis_cache/keys"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, userId uint32, id int64) (*kxmj_core.UserItem, error) {
	key := generateKey(userId)
	value := &kxmj_core.UserItem{}
	field := fmt.Sprintf("%d", id)
	v, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}

	val := &kxmj_core.UserItem{}
	err = json.Unmarshal([]byte(v), val)
	if err != nil {
		log.Sugar().Errorf("Unmarshal err:%v", err)
		return nil, err
	}

	return value, nil
}

func (c *Cache) GetAll(ctx context.Context, userId uint32) (map[uint32]*kxmj_core.UserItem, error) {
	key := generateKey(userId)
	value := make(map[uint32]*kxmj_core.UserItem, 0)
	maps, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	for _, v := range maps {
		val := &kxmj_core.UserItem{}
		err = json.Unmarshal([]byte(v), val)
		if err != nil {
			log.Sugar().Errorf("Unmarshal err:%v", err)
			continue
		}
		value[val.ItemId] = val
	}
	return value, nil
}

func (c *Cache) Set(ctx context.Context, user *kxmj_core.UserItem) error {
	key := generateKey(user.UserId)
	maps := make(map[string]string, 0)
	field := fmt.Sprintf("%d", user.Id)
	val, _ := json.Marshal(user)
	maps[field] = string(val)
	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) BulkSet(ctx context.Context, users []*kxmj_core.UserItem, userId uint32) error {
	if len(users) <= 0 {
		return nil
	}

	key := generateKey(userId)
	maps := make(map[string]string, 0)
	for _, user := range users {
		field := fmt.Sprintf("%d", user.Id)
		val, _ := json.Marshal(user)
		maps[field] = string(val)
	}
	return c.client.HMSet(ctx, key, maps).Err()
}

func (c *Cache) Del(ctx context.Context, userId uint32) error {
	key := generateKey(userId)
	c.client.Del(ctx, key)
	return nil
}

func (c *Cache) Lock(ctx context.Context, userId uint32) {
	key := generateLockKey(userId)
	retryCount := 0
	for {
		if retryCount >= 100 {
			log.Sugar().Errorf("Lock user:%d failed...", userId)
			break
		}

		val, _ := c.client.Get(ctx, key).Int()
		if val > 0 {
			time.Sleep(time.Millisecond * 100)
			retryCount++
			continue
		}
		break
	}
	c.client.Set(ctx, key, 1, time.Second*10)
}

func (c *Cache) Unlock(ctx context.Context, userId uint32) {
	key := generateLockKey(userId)
	c.client.Del(ctx, key)
}

func generateKey(userId uint32) string {
	return fmt.Sprintf(keys.UserItemFormatKey, userId)
}

func generateLockKey(userId uint32) string {
	return fmt.Sprintf(keys.UserItemLockerFormatKey, userId)
}
