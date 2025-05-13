package token

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"kxmj.common/redis_cache/keys"
	"kxmj.common/utils"
	"time"
)

type Cache struct {
	client *redis.Client
	expire int
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
		expire: 3600 * 48,
	}
}

func (ls *Cache) GetUserId(ctx context.Context, token string) (int, error) {
	key := ls.generateTokenKey(token)
	value, err := ls.client.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	ls.client.Expire(ctx, key, time.Duration(ls.expire)*time.Second)
	return value, nil
}

func (ls *Cache) SetToken(ctx context.Context, userId int) (string, error) {
	userKey := ls.generateUserKey(userId)
	curToken, _ := ls.client.Get(ctx, userKey).Result()
	if len(curToken) > 0 {
		//ls.client.Expire(ctx, curToken, 120*time.Second)
		ls.client.Del(ctx, ls.generateTokenKey(curToken))
	}

	token := ls.generateToken(userId)
	tokenKey := ls.generateTokenKey(token)
	err := ls.client.Set(ctx, tokenKey, userId, time.Duration(ls.expire)*time.Second).Err()
	if err != nil {
		return "", err
	}

	err = ls.client.Set(ctx, userKey, token, time.Duration(ls.expire)*time.Second).Err()
	return token, err
}

func (ls *Cache) CheckToken(ctx context.Context, userId int, token string) (bool, error) {
	userKey := ls.generateUserKey(userId)
	curToken, err := ls.client.Get(ctx, userKey).Result()
	if err != nil {
		return false, err
	}
	return curToken == token, nil
}

func (ls *Cache) Expired(ctx context.Context, token string, userId int) error {
	if ls.expire == 0 {
		return nil
	}

	userKey := ls.generateUserKey(userId)
	err := ls.client.Expire(ctx, userKey, time.Duration(ls.expire)*time.Second).Err()
	if err != nil {
		return err
	}

	tokenKey := ls.generateTokenKey(token)
	err = ls.client.Expire(ctx, tokenKey, time.Duration(ls.expire)*time.Second).Err()
	return err
}

func (ls *Cache) Delete(ctx context.Context, token string) error {
	userId, err := ls.GetUserId(ctx, token)
	if err != nil {
		return err
	}

	userKey := ls.generateUserKey(userId)
	err = ls.client.Del(ctx, userKey).Err()
	if err != nil {
		return err
	}

	tokenKey := ls.generateTokenKey(token)
	err = ls.client.Del(ctx, tokenKey).Err()
	return err
}

func (ls *Cache) generateToken(userId int) string {
	value := fmt.Sprintf("%d_%s", userId, uuid.New().String())
	return utils.Sha256(value)
}

func (ls *Cache) generateTokenKey(token string) string {
	prefix := keys.TokenFormatKey
	return fmt.Sprintf(prefix, token)
}

func (ls *Cache) generateUserKey(userId int) string {
	prefix := keys.UserFormatKey
	return fmt.Sprintf(prefix, userId)
}
