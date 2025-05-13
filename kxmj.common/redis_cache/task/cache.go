package task

import (
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/task/login"
)

type Cache struct {
	loginCache *login.Cache // 登陆奖品缓存
}

var temp *Cache

func NewCache(client *redis.Client) *Cache {
	temp = &Cache{}

	// 缓存同步不走canal
	temp.loginCache = login.NewCache(client)
	return temp
}

func (c *Cache) LoginCache() *login.Cache {
	return c.loginCache
}
