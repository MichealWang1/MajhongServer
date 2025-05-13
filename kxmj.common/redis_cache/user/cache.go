package user

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/redis_core"
	"kxmj.common/redis_cache/user/detail"
	"kxmj.common/redis_cache/user/equip"
	"kxmj.common/redis_cache/user/id"
	"kxmj.common/redis_cache/user/item"
	"kxmj.common/redis_cache/user/tel"
	"kxmj.common/redis_cache/user/vip"
	"kxmj.common/redis_cache/user/wallet"
	"kxmj.common/redis_cache/user/wechat_id"
	"reflect"
	"time"
)

type Cache struct {
	caches        []redis_core.ICache  // 缓存列表
	handler       *redis_core.Handler  // 事件处理器
	register      *redis_core.Register // 事件订阅注册器
	idCache       *id.Cache            // 用户id缓存
	telCache      *tel.Cache           // 手机号缓存
	detailCache   *detail.Cache        // 用户信息缓存
	walletCache   *wallet.Cache        // 用户钱包缓存
	wechatIdCache *wechat_id.Cache     // 微信ID缓存
	equipCache    *equip.Cache         // 用户装备缓存
	itemCache     *item.Cache          // 用户背包物品缓存
	vipCache      *vip.Cache           // 用户VIP缓存
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

	temp.telCache = tel.NewCache(client)
	temp.caches = append(temp.caches, temp.telCache)

	temp.detailCache = detail.NewCache(client)
	temp.caches = append(temp.caches, temp.detailCache)

	temp.walletCache = wallet.NewCache(client)
	temp.caches = append(temp.caches, temp.walletCache)

	temp.wechatIdCache = wechat_id.NewCache(client)
	temp.caches = append(temp.caches, temp.wechatIdCache)

	temp.vipCache = vip.NewCache(client)
	temp.caches = append(temp.caches, temp.vipCache)

	// 缓存同步不走canal
	temp.equipCache = equip.NewCache(client)
	temp.itemCache = item.NewCache(client)

	return temp
}

func (c *Cache) LoadCache() {
	// 加载用户缓存
	err := c.loadUserCache()
	if err != nil {
		panic(err)
	}

	// 加载用户第三方数据缓存
	err = c.loadUserThirdPartyCache()
	if err != nil {
		panic(err)
	}

	// 加载用户装备缓存
	err = c.loadEquipCache()
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

func (c *Cache) IdCache() *id.Cache {
	return c.idCache
}

func (c *Cache) TelCache() *tel.Cache {
	return c.telCache
}

func (c *Cache) DetailCache() *detail.Cache {
	return c.detailCache
}

func (c *Cache) WalletCache() *wallet.Cache {
	return c.walletCache
}

func (c *Cache) WeChatCache() *wechat_id.Cache {
	return c.wechatIdCache
}

func (c *Cache) EquipCache() *equip.Cache {
	return c.equipCache
}

func (c *Cache) ItemCache() *item.Cache {
	return c.itemCache
}

func (c *Cache) VIPCache() *vip.Cache {
	return c.vipCache
}

func (c *Cache) loadUserCache() error {
	ctx := context.Background()

	startId := int64(0)
	count := 5000
	startTime := time.Now().UnixMilli()
	// 添加用户手机号缓存
	for {
		queries := make([]*kxmj_core.User, 0)
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

		err = c.telCache.BulkSet(ctx, queries)
		if err != nil {
			return err
		}

		log.Sugar().Info(fmt.Sprintf("loadUserCache %d success", len(queries)))
		if len(queries) < count {
			break
		}
		startId = queries[len(queries)-1].Id
	}

	endTime := time.Now().UnixMilli()
	log.Sugar().Info(fmt.Sprintf("loadUserCache elapsed: %dms  ", endTime-startTime))
	return nil
}

func (c *Cache) loadUserThirdPartyCache() error {
	ctx := context.Background()

	startId := int64(0)
	count := 5000
	startTime := time.Now().UnixMilli()
	// 添加用户手机号缓存
	for {
		queries := make([]*kxmj_core.UserThirdParty, 0)
		err := mysql.CoreMaster().WithContext(ctx).
			Where("id > ?", startId).
			Limit(count).
			Find(&queries).Error

		if err != nil {
			return err
		}

		err = c.wechatIdCache.BulkSet(ctx, queries)
		if err != nil {
			return err
		}

		log.Sugar().Info(fmt.Sprintf("loadUserThirdPartyCache %d success", len(queries)))
		if len(queries) < count {
			break
		}
		startId = queries[len(queries)-1].Id
	}

	endTime := time.Now().UnixMilli()
	log.Sugar().Info(fmt.Sprintf("loadUserThirdPartyCache elapsed: %dms  ", endTime-startTime))
	return nil
}

func (c *Cache) loadEquipCache() error {
	ctx := context.Background()

	startId := int64(0)
	count := 5000
	startTime := time.Now().UnixMilli()
	// 添加用户手机号缓存
	for {
		queries := make([]*kxmj_core.UserEquip, 0)
		err := mysql.CoreMaster().WithContext(ctx).
			Where("id > ?", startId).
			Limit(count).
			Find(&queries).Error

		if err != nil {
			return err
		}

		err = c.equipCache.BulkSet(ctx, queries)
		if err != nil {
			return err
		}

		log.Sugar().Info(fmt.Sprintf("loadEquipCache %d success", len(queries)))
		if len(queries) < count {
			break
		}
		startId = queries[len(queries)-1].Id
	}

	endTime := time.Now().UnixMilli()
	log.Sugar().Info(fmt.Sprintf("loadEquipCache elapsed: %dms  ", endTime-startTime))
	return nil
}
