package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/redis_cache/redis_core"
	"strconv"
)

type Cache struct {
	client *redis.Client
}

// NewCache 创建
func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func generateAgentKey() string {
	return fmt.Sprintf(keys.MailFormatKey)
}

// EventHandler 更新数据回调 当数据库里面数据更新时，会同步到redis中
func (c *Cache) EventHandler(ctx context.Context, e *redis_core.EventParams) {
	if e == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}
	if e.Data == nil {
		log.Sugar().Error(fmt.Sprintf("Invalid event params"))
		return
	}
	// 把数据转成结构体 判断 结构体中 OrderId 和 UserId 是否为空
	params := e.Data.(*kxmj_core.ConfigEmail)
	if params.EmailId <= 0 {
		return
	}
	// 判断 是否删除 还是 更新
	err := c.DelAll(ctx)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Del  params:%v err:%v", e, err))
	}
}

// GetTableTemplate 获取数据表类型
func (c *Cache) GetTableTemplate() interface{} {
	return &kxmj_core.ConfigEmail{}
}

// Del 删除一个 eMail 配置
func (c *Cache) Del(ctx context.Context, emailId uint32) error {
	// 获取 key
	key := generateAgentKey()
	// key  fields 是 emailId
	filed := fmt.Sprintf("%d", emailId)
	return c.client.HDel(ctx, key, filed).Err()
}

// DelAll 删除整个邮件配置
func (c *Cache) DelAll(ctx context.Context) error {
	key := generateAgentKey()
	return c.client.Del(ctx, key).Err()
}

// Set 设置单个邮件配置
func (c *Cache) Set(ctx context.Context, email *kxmj_core.ConfigEmail) error {
	data, err := json.Marshal(email)
	if err != nil {
		return err
	}
	key := generateAgentKey()
	filed := fmt.Sprintf("%d", email.EmailId)
	return c.client.HSet(ctx, key, map[string]string{filed: string(data)}).Err()
}

// SetAll 设置所有的邮件配置
func (c *Cache) SetAll(ctx context.Context) error {
	var mailConfig []*kxmj_core.ConfigEmail
	err := mysql.CoreMaster().WithContext(ctx).Where("email_id > 0").Find(&mailConfig).Error
	if err != nil {
		return err
	}
	if len(mailConfig) <= 0 {
		return errors.New("SetAll mailConfig is null")
	}
	maps := make(map[string]string, 0)
	for _, mailDate := range mailConfig {
		data, err := json.Marshal(mailDate)
		if err != nil {
			return err
		}
		field := fmt.Sprintf("%d", mailDate.EmailId)
		maps[field] = string(data)
	}
	key := generateAgentKey()
	return c.client.HMSet(ctx, key, maps).Err()
}

// Get 获取单个邮件 根据 邮件ID
func (c *Cache) Get(ctx context.Context, emildId uint32) (*kxmj_core.ConfigEmail, error) {
	// 获取 key
	key := generateAgentKey()
	value := &kxmj_core.ConfigEmail{}
	// 获取 field
	field := fmt.Sprintf("%d", emildId)
	// 在redis中获取 Email 配置
	data, err := c.client.HGet(ctx, key, field).Result()
	// redis中没有数据
	if err != nil || len(data) <= 0 {
		// 根据 email_id 从数据库 获取 信息
		err = mysql.CoreMaster().Where("email_id = ? ", emildId).First(value).Error
		if err != nil {
			return nil, err
		}
		// 从数据库获取到数据后 再存入Redis
		err = c.Set(ctx, value)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(data), value)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

// GetAll 获取所有邮件的配置
func (c *Cache) GetAll(ctx context.Context) (map[uint32]*kxmj_core.ConfigEmail, error) {
	key := generateAgentKey()
	value := make(map[uint32]*kxmj_core.ConfigEmail, 0)
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) <= 0 {
		err = c.SetAll(ctx)
		if err != nil {
			return nil, err
		}
		return c.GetAll(ctx)
	} else {
		for k, d := range data {
			mailConfig := &kxmj_core.ConfigEmail{}
			err = json.Unmarshal([]byte(d), mailConfig)
			if err != nil {
				continue
			}
			mailId, _ := strconv.ParseUint(k, 10, 32)
			// 判断邮件ID 是否 大于0
			if mailId > 0 {
				value[uint32(mailId)] = mailConfig
			}
		}
	}
	return value, nil
}
