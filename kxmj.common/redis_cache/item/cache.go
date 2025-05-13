package item

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/item/detail"
	"kxmj.common/redis_cache/redis_core"
	"reflect"
)

type Cache struct {
	caches      []redis_core.ICache  // 缓存列表
	handler     *redis_core.Handler  // 事件处理器
	register    *redis_core.Register // 事件订阅注册器
	detailCache *detail.Cache        // 用户信息缓存
}

var temp *Cache

func NewCache(client *redis.Client, handler *redis_core.Handler, register *redis_core.Register) *Cache {
	temp = &Cache{
		caches:   make([]redis_core.ICache, 0),
		handler:  handler,
		register: register,
	}

	temp.detailCache = detail.NewCache(client)
	temp.caches = append(temp.caches, temp.detailCache)

	return temp
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

func (c *Cache) GetDetailCache() *detail.Cache {
	return c.detailCache
}
