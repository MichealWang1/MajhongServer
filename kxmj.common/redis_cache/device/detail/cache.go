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
	"time"
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

	if e.Action == redis_core.UpdateAction {
		err := c.Del(ctx, params.DeviceId)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
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

func (c *Cache) Get(ctx context.Context, deviceId string) (*kxmj_core.Device, error) {
	key := generateKey(deviceId)
	value := &kxmj_core.Device{}
	err := c.client.HGetAll(ctx, key).Scan(value)
	if err != nil || value.Id <= 0 {
		err = mysql.CoreMaster().Where("device_id = ?", deviceId).First(value).Error
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

func (c *Cache) Set(ctx context.Context, device *kxmj_core.Device) error {
	maps, err := utils.StructToMap(device)
	if err != nil {
		return err
	}

	key := generateKey(device.DeviceId)
	_, err = c.client.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range maps {
			err = pipeline.HSet(ctx, key, k, v).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})

	c.client.Expire(ctx, key, time.Hour*8)
	return err
}

func (c *Cache) Del(ctx context.Context, deviceId string) error {
	key := generateKey(deviceId)
	c.client.Del(ctx, key)
	return nil
}

func generateKey(deviceId string) string {
	return fmt.Sprintf(keys.DeviceInfoFormatKey, deviceId)
}
