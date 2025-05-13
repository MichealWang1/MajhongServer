package redis_cache

import (
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/bundle"
	"kxmj.common/redis_cache/device"
	"kxmj.common/redis_cache/gm"
	"kxmj.common/redis_cache/goods"
	"kxmj.common/redis_cache/item"
	"kxmj.common/redis_cache/limit"
	"kxmj.common/redis_cache/mail"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.common/redis_cache/room"
	"kxmj.common/redis_cache/sms"
	"kxmj.common/redis_cache/task"
	"kxmj.common/redis_cache/token"
	"kxmj.common/redis_cache/user"
)

type Cache struct {
	register    *redis_core.Register
	handler     *redis_core.Handler
	tokenCache  *token.Cache
	limitCache  *limit.Cache
	smsCache    *sms.Cache
	userCache   *user.Cache
	deviceCache *device.Cache
	bundleCache *bundle.Cache
	roomCache   *room.Cache
	itemCache   *item.Cache
	goodsCache  *goods.Cache
	mailCache   *mail.Cache
	gmCache     *gm.Cache
	taskCache   *task.Cache
}

var temp *Cache

// InitPublishCache 实例化发布者缓存
func InitPublishCache(client *redis.Client) {
	temp = &Cache{}
	temp.register = redis_core.NewRegister(client)
	temp.handler = redis_core.NewHandler(client)
}

// InitConsumerCache 实例化消费者缓存
func InitConsumerCache(client *redis.Client) {
	temp = &Cache{}
	temp.register = redis_core.NewRegister(client)
	temp.handler = redis_core.NewHandler(client)
	temp.initCache(client)

	// 加载全量永久缓存
	temp.userCache.LoadCache()
	temp.deviceCache.LoadCache()

	// 注册缓存更新订阅事件
	temp.userCache.Register()
	temp.deviceCache.Register()
	temp.bundleCache.Register()
	temp.roomCache.Register()
	temp.itemCache.Register()
	temp.goodsCache.Register()
	temp.mailCache.Register()
	// 启动消费者
	temp.handler.Consumer()
}

// InitReadCache 实例化只读缓存
func InitReadCache(client *redis.Client) {
	temp = &Cache{}
	temp.initCache(client)
}

func (c *Cache) initCache(client *redis.Client) {
	temp.tokenCache = token.NewCache(client)
	temp.limitCache = limit.NewCache(client)
	temp.smsCache = sms.NewCache(client)
	temp.userCache = user.NewCache(client, temp.handler, temp.register)
	temp.deviceCache = device.NewCache(client, temp.handler, temp.register)
	temp.bundleCache = bundle.NewCache(client, temp.handler, temp.register)
	temp.roomCache = room.NewCache(client, temp.handler, temp.register)
	temp.itemCache = item.NewCache(client, temp.handler, temp.register)
	temp.goodsCache = goods.NewCache(client, temp.handler, temp.register)
	temp.mailCache = mail.NewCache(client, temp.handler, temp.register)
	temp.gmCache = gm.NewCache(client)
	temp.taskCache = task.NewCache(client)
}

// GetCache 获取缓存实例
func GetCache() *Cache {
	return temp
}

// GetRegister 获取注册器实例
func (c *Cache) GetRegister() *redis_core.Register {
	return c.register
}

// GetHandler 获取事件处理器
func (c *Cache) GetHandler() *redis_core.Handler {
	return c.handler
}

// GetTokenCache 获取token缓存实例
func (c *Cache) GetTokenCache() *token.Cache {
	return c.tokenCache
}

// GetLimitCache 获取limit缓存实例
func (c *Cache) GetLimitCache() *limit.Cache {
	return c.limitCache
}

// GetSmsCache 获取sms缓存实例
func (c *Cache) GetSmsCache() *sms.Cache {
	return c.smsCache
}

// GetUserCache 获取用户缓存实例
func (c *Cache) GetUserCache() *user.Cache {
	return c.userCache
}

// GetDeviceCache 获取设备缓存实例
func (c *Cache) GetDeviceCache() *device.Cache {
	return c.deviceCache
}

// GetBundleCache 获取分包缓存实例
func (c *Cache) GetBundleCache() *bundle.Cache {
	return c.bundleCache
}

// GetRoomCache 获取游戏房间缓存实例
func (c *Cache) GetRoomCache() *room.Cache {
	return c.roomCache
}

// GetItemCache 获取物品缓存实例
func (c *Cache) GetItemCache() *item.Cache {
	return c.itemCache
}

// GetGoodsCache 获取商品缓存实例
func (c *Cache) GetGoodsCache() *goods.Cache {
	return c.goodsCache
}

// GetMailCache 获取邮箱实例
func (c *Cache) GetMailCache() *mail.Cache {
	return c.mailCache
}

// GetGameManageCache 获取游戏GM 实例
func (c *Cache) GetGameManageCache() *gm.Cache {
	return c.gmCache
}

// GetTaskCache 获取任务缓存实例
func (c *Cache) GetTaskCache() *task.Cache {
	return c.taskCache
}
