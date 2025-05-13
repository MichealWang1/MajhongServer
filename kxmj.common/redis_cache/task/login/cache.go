package login

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/utils"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) GetDaily(ctx context.Context, userId uint32) (map[string]string, error) {
	dailyKey := generateDailyKey(userId)
	return c.client.HGetAll(ctx, dailyKey).Result()
}

func (c *Cache) SetDaily(ctx context.Context, userId uint32, data map[string]string) error {
	dailyKey := generateDailyKey(userId)
	return c.client.HMSet(ctx, dailyKey, data).Err()
}

func (c *Cache) GetTotal(ctx context.Context, userId uint32) (map[string]string, error) {
	totalKey := generateTotalKey(userId)
	return c.client.HGetAll(ctx, totalKey).Result()
}

func (c *Cache) SetTotal(ctx context.Context, userId uint32, data map[string]string) error {
	totalKey := generateTotalKey(userId)
	return c.client.HMSet(ctx, totalKey, data).Err()
}

func (c *Cache) Set(ctx context.Context, userId uint32) error {
	dateKey := generateDateKey(userId)
	today := utils.GetZeroUnix()
	date, _ := c.client.Get(ctx, dateKey).Int64()
	if date == today {
		return nil
	}

	err := c.client.Set(ctx, dateKey, today, time.Hour*24).Err()
	if err != nil {
		return err
	}

	// 每日登陆奖品逻辑
	dailyKey := generateDailyKey(userId)
	maps, _ := c.client.HGetAll(ctx, dailyKey).Result()
	// 如果是第一次登陆或者中间有间隔，那么用户只能领取第一天奖品
	if len(maps) <= 0 {
		// 每天登陆第一天奖品定义
		maps = map[string]string{
			"c": "0", // 连续登陆天数
			"1": "0", // 第一天 已完成(状态: 0 未完成；1 已完成；2 已领取)
			"2": "0", // 第二天 未完成(状态: 0 未完成；1 已完成；2 已领取)
			"3": "0", // 第三天 未完成(状态: 0 未完成；1 已完成；2 已领取)
		}
	} else {
		// 如果中间没有连续登陆，奖品清零
		if date == 0 {
			maps["c"] = "0"
			maps["1"] = "0"
			maps["2"] = "0"
			maps["3"] = "0"
		}
	}

	// 天数累计
	count := maps["c"]
	temp, _ := utils.Add(count, "1")
	count = temp.String()
	maps["c"] = count

	// 连续登陆3天循环
	if utils.Cmp(count, "4") >= 0 {
		maps["c"] = "1"
		maps["1"] = "0"
		maps["2"] = "0"
		maps["3"] = "0"
	}

	if maps["c"] == "1" {
		if maps["1"] != "1" {
			maps["1"] = "1"
		}
	} else if maps["c"] == "2" {
		if maps["2"] != "1" {
			maps["2"] = "1"
		}
	} else if maps["c"] == "3" {
		if maps["3"] != "1" {
			maps["3"] = "1"
		}
	}

	err = c.client.HMSet(ctx, dailyKey, maps).Err()
	if err != nil {
		return err
	}

	err = c.client.Expire(ctx, dailyKey, time.Hour*24).Err()
	if err != nil {
		return err
	}

	// 累计登陆奖品逻辑
	totalKey := generateTotalKey(userId)
	maps, _ = c.client.HGetAll(ctx, totalKey).Result()
	if len(maps) <= 0 {
		// 超过一个星期没登陆，奖品清零
		maps = map[string]string{
			"c": "0", // 次数
			"2": "0", // 累计2天 已完成(状态: 0 未完成；1 已完成；2 已领取)
			"4": "0", // 累计4天 未完成(状态: 0 未完成；1 已完成；2 已领取)
			"6": "0", // 累计6天 未完成(状态: 0 未完成；1 已完成；2 已领取)
		}
	}

	// 天数累计
	count = maps["c"]
	temp, _ = utils.Add(count, "1")
	count = temp.String()
	maps["c"] = count

	// 累计登陆超过6天，从第一天开始计算
	if utils.Cmp(count, "6") > 0 {
		// 每天登陆第一天奖品定义
		maps["c"] = "1"
	}

	if maps["c"] == "2" {
		if maps["2"] != "1" {
			maps["2"] = "1"
		}
	} else if maps["c"] == "4" {
		if maps["4"] != "1" {
			maps["4"] = "1"
		}
	} else if maps["c"] == "6" {
		if maps["6"] != "1" {
			maps["6"] = "1"
		}
	}

	err = c.client.HMSet(ctx, totalKey, maps).Err()
	if err != nil {
		return err
	}

	err = c.client.Expire(ctx, totalKey, time.Hour*24*7).Err()
	if err != nil {
		return err
	}

	return nil
}

func generateDateKey(userId uint32) string {
	return fmt.Sprintf(keys.TaskDateLoginFormatKey, userId)
}

func generateDailyKey(userId uint32) string {
	return fmt.Sprintf(keys.TaskDailyLoginFormatKey, userId)
}

func generateTotalKey(userId uint32) string {
	return fmt.Sprintf(keys.TaskTotalLoginFormatKey, userId)
}
