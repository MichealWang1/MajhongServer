package detail

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.common/utils"
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

	params := e.Data.(*kxmj_core.ConfigBundle)
	if len(params.BundleId) <= 0 {
		return
	}

	if e.Action == redis_core.UpdateAction {
		err := c.Del(ctx, params.BundleId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
		}
	} else if e.Action == redis_core.DeleteAction {
		err := c.Del(ctx, params.BundleId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
		}
	}
}

func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.ConfigBundle{}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, bundleId string) (*kxmj_core.ConfigBundle, error) {
	key := generateKey(bundleId)
	value := &kxmj_core.ConfigBundle{}
	err := c.client.HGetAll(ctx, key).Scan(value)
	if err != nil || len(value.BundleId) <= 0 {
		err = mysql.CoreMaster().Where("bundle_id = ?", bundleId).First(value).Error
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

func (c *Cache) Set(ctx context.Context, bundle *kxmj_core.ConfigBundle) error {
	maps, err := utils.StructToMap(bundle)
	if err != nil {
		return err
	}

	key := generateKey(bundle.BundleId)
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

func (c *Cache) Del(ctx context.Context, bundleId string) error {
	key := generateKey(bundleId)
	c.client.Del(ctx, key)
	return nil
}

func generateKey(bundleId string) string {
	return fmt.Sprintf(keys.BundleFormatKey, bundleId)
}
