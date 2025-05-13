package mail

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/mail/config"
	"kxmj.common/redis_cache/mail/system"
	"kxmj.common/redis_cache/mail/welfare"
	"kxmj.common/redis_cache/redis_core"
	"reflect"
)

type Cache struct {
	caches          []redis_core.ICache  // 缓存列表
	handler         *redis_core.Handler  // 事件处理器
	register        *redis_core.Register // 事件订阅注册器
	mailConfigCache *config.Cache        // 邮箱配置
	welfareCache    *welfare.Cache       // 玩家福利邮件
	systemCache     *system.Cache        // 玩家系统邮件
}

// 当前邮箱中全局的redis Cache
var redisCache *Cache

// 创建 redis Cache
func NewCache(client *redis.Client, handler *redis_core.Handler, register *redis_core.Register) *Cache {
	redisCache = &Cache{
		caches:   make([]redis_core.ICache, 0),
		handler:  handler,
		register: register,
	}
	redisCache.mailConfigCache = config.NewCache(client)
	redisCache.welfareCache = welfare.NewCache(client)
	redisCache.systemCache = system.NewCache(client)
	redisCache.caches = append(redisCache.caches, redisCache.mailConfigCache)
	return redisCache
}

func (c *Cache) GetDetailCache() *config.Cache {
	return c.mailConfigCache
}

func (c *Cache) GetWelfareMailCache() *welfare.Cache {
	return c.welfareCache
}

func (c *Cache) GetSystemMailCache() *system.Cache {
	return c.systemCache
}

func (c *Cache) Register() {
	for _, cache := range c.caches {
		// 数据库表模板实例
		entity := cache.GetTableTemplate()
		// 反射获取 结构体 name
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
