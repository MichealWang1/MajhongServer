package device

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/device/detail"
	"kxmj.common/redis_cache/device/id"
	"kxmj.common/redis_cache/redis_core"
	"reflect"
	"time"
)

type Cache struct {
	caches      []redis_core.ICache  // 缓存列表
	handler     *redis_core.Handler  // 事件处理器
	register    *redis_core.Register // 事件订阅注册器
	idCache     *id.Cache            // 用户id缓存
	detailCache *detail.Cache        // 手机号缓存
}

var temp *Cache

func NewCache(client *redis.Client, handler *redis_core.Handler, register *redis_core.Register) *Cache {
	temp = &Cache{
		caches:   make([]redis_core.ICache, 0),
		handler:  handler,
		register: register,
	}

	temp.idCache = id.NewCache(client)
	temp.caches = append(temp.caches, temp.idCache)

	temp.detailCache = detail.NewCache(client)
	temp.caches = append(temp.caches, temp.detailCache)

	return temp
}

func (c *Cache) LoadCache() {
	// 加载用户缓存
	err := c.loadDeviceCache()
	if err != nil {
		panic(err)
	}
}

func (c *Cache) Register() {
	for _, cache := range c.caches {
		entity := cache.GetTableTemplate()
		// 反射获取table name
		reflectValue := reflect.ValueOf(entity)
		method := reflectValue.MethodByName("TableName")
		values := method.Call(nil)
		tableName := values[0].String()

		method = reflectValue.MethodByName("Schema")
		values = method.Call(nil)
		schema := values[0].String()

		// 注册到事件处理器
		c.handler.Register(schema, tableName, cache)

		ctx := context.Background()

		// 注册到订阅器
		err := c.register.Set(ctx, schema, tableName)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Cache) GetIdCache() *id.Cache {
	return c.idCache
}

func (c *Cache) GetDetailCache() *detail.Cache {
	return c.detailCache
}

func (c *Cache) loadDeviceCache() error {
	ctx := context.Background()

	startId := uint32(0)
	count := 5000
	startTime := time.Now().UnixMilli()
	// 添加用户手机号缓存
	for {
		queries := make([]*kxmj_core.Device, 0)
		err := mysql.CoreMaster().WithContext(ctx).
			Where("id > ?", startId).
			Limit(count).
			Find(&queries).Error

		if err != nil {
			return err
		}

		err = c.idCache.BulkSet(ctx, queries)
		if err != nil {
			return err
		}

		log.Sugar().Info(fmt.Sprintf("load device %d success", len(queries)))
		if len(queries) < count {
			break
		}
		startId = queries[len(queries)-1].Id
	}

	endTime := time.Now().UnixMilli()
	log.Sugar().Info(fmt.Sprintf("loadDeviceCache elapsed: %dms  ", endTime-startTime))
	return nil
}
