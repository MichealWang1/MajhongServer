package usermail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/log"
	"kxmj.common/redis_cache/keys"
)

type Cache struct {
	client *redis.Client
}

type EmailUser struct {
	EmailId   uint32 `json:"emailId"`
	EmailType uint8  `json:"emailType"`
	Order     int64  `json:"order"`
	Status    uint8  `json:"status"`
	SendTime  uint32 `json:"sendTime"`
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// generateWelfareMailKey 获取玩家福利邮件前缀
func generateWelfareMailKey(userId uint32) string {
	return fmt.Sprintf(keys.UserWelfareMailFormatKey, userId)
}

// GetWelfareMail 获取玩家指定 福利邮件的状态 HGet
func (c *Cache) GetWelfareMail(ctx context.Context, userId uint32, mailId uint32) (*EmailUser, error) {
	key := generateWelfareMailKey(userId)
	field := fmt.Sprintf("%d", mailId)
	value, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}
	mailValue := &EmailUser{}
	err = json.Unmarshal([]byte(value), mailValue)
	if err != nil {
		log.Sugar().Errorf("Get User MailId %d Unmarshal  err:%v", mailId, err)
		return nil, err
	}
	return mailValue, nil
}

// GetAllWelfareMail 获取玩家所有的福利邮件
func (c *Cache) GetAllWelfareMail(ctx context.Context, userId uint32) (map[uint32]*EmailUser, error) {
	key := generateWelfareMailKey(userId)
	mailMap := make(map[uint32]*EmailUser, 0)
	maps, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range maps {
		value := &EmailUser{}
		err := json.Unmarshal([]byte(v), value)
		if err != nil {
			continue
		}
		mailMap[value.EmailId] = value
	}
	if len(mailMap) <= 0 {
		return nil, errors.New(fmt.Sprintf("GetAllWelfareMail len(mailList)==0 userId:%d ", userId))
	}
	return mailMap, nil
}

// SetWelfareMail 设置玩家单个福利邮件的状态 HMSet
func (c *Cache) SetWelfareMail(ctx context.Context, userId uint32, mailId uint32, data *EmailUser) error {
	key := generateWelfareMailKey(userId)
	field := fmt.Sprintf("%d", mailId)

	val, _ := json.Marshal(data)
	maps := make(map[string]string, 0)
	maps[field] = string(val)
	return c.client.HMSet(ctx, key, maps).Err()
}

// DelUserWelfareMail 删除玩家所有的 福利邮件
func (c *Cache) DelUserWelfareMail(ctx context.Context, userId uint32) error {
	key := generateWelfareMailKey(userId)
	mailNum, err := c.client.Del(ctx, key).Result()
	if err != nil {
		log.Sugar().Errorf("DelUserAllMail userId:%d fali:%v", userId, err)
		return err
	}
	log.Sugar().Infof("DelUserAllMail userId:%d suc mailNum:%d", userId, mailNum)
	return nil
}

// DelUserWelfareMailByMailId 删除玩家单封福利邮件
func (c *Cache) DelUserWelfareMailByMailId(ctx context.Context, userId uint32, mailId uint32) error {
	key := generateWelfareMailKey(userId)
	field := fmt.Sprintf("%d", mailId)
	_, err := c.client.HDel(ctx, key, field).Result()
	if err != nil {
		log.Sugar().Errorf("DelUserMailByMailId userId:%d  mailId%d fail 失败:%v", userId, mailId, err)
		return err
	}
	log.Sugar().Infof("DelUserMailByMailId userId:%d suc mailId:%d", userId, mailId)
	return nil
}
