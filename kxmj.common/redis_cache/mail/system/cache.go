package system

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/log"
	"kxmj.common/redis_cache/keys"
	"strconv"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// generateSystemMailKey 获取玩家系统邮件前缀
func generateSystemMailKey(userId uint32) string {
	return fmt.Sprintf(keys.UserSystemMailFormatKey, userId)
}

// GetAll 获取玩家系统邮件
func (c *Cache) GetAll(ctx context.Context, userId uint32) (map[uint32]uint8, error) {
	key := generateSystemMailKey(userId)
	valueMap, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	mailList := make(map[uint32]uint8, 0)
	for k, v := range valueMap {
		mailId, _ := strconv.ParseUint(k, 10, 32)
		state, _ := strconv.ParseUint(v, 10, 8)
		if mailId > 0 && state > 0 {
			mailList[uint32(mailId)] = uint8(state)
		}
	}
	if len(mailList) <= 0 {
		return nil, errors.New(fmt.Sprintf("GetSystemAllMail len(mailList)==0 userId:%d ", userId))
	}
	return mailList, nil
}

// Set 设置玩家单个邮件的状态 HSet
func (c *Cache) Set(ctx context.Context, userId uint32, mailId uint32, state uint8) error {
	key := generateSystemMailKey(userId)
	field := fmt.Sprintf("%d", mailId)
	maps := make(map[string]string, 0)
	maps[field] = fmt.Sprintf("%d", state)
	return c.client.HMSet(ctx, key, maps).Err()
}

// Get 获取玩家指定 系统邮件的状态 HGet
func (c *Cache) Get(ctx context.Context, userId uint32, mailId uint32) (uint8, error) {
	key := generateSystemMailKey(userId)
	field := fmt.Sprintf("%d", mailId)
	value, err := c.client.HGet(ctx, key, field).Int()
	if err != nil {
		return 0, err
	}
	return uint8(value), nil
}

// DelAll 删除玩家所有的 福利邮件
func (c *Cache) DelAll(ctx context.Context, userId uint32) error {
	key := generateSystemMailKey(userId)
	mailNum, err := c.client.Del(ctx, key).Result()
	if err != nil {
		log.Sugar().Errorf("DelAllUserSystemMail userId:%d fali:%v", userId, err)
		return err
	}
	log.Sugar().Infof("DelAllUserSystemMail userId:%d suc mailNum:%d", userId, mailNum)
	return nil
}

// Del 删除玩家单封福利邮件
func (c *Cache) Del(ctx context.Context, userId uint32, mailId uint32) error {
	key := generateSystemMailKey(userId)
	field := fmt.Sprintf("%d", mailId)
	_, err := c.client.HDel(ctx, key, field).Result()
	if err != nil {
		log.Sugar().Errorf("DelUserSystemMailByMailId userId:%d  mailId%d fail 失败:%v", userId, mailId, err)
		return err
	}
	log.Sugar().Infof("DelUserSystemMailByMailId userId:%d suc mailId:%d", userId, mailId)
	return nil
}
