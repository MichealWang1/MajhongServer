package business

import (
	"context"
	"fmt"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"kxmj.core.api/internal/db"
	"kxmj.core.api/internal/dto"
	"math/rand"
	"time"
)

func CreateUserId(ctx context.Context) (uint32, error) {
	var userId uint32
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		userId = uint32(r.Intn(80000000) + 20000000)
		exist := redis_cache.GetCache().GetUserCache().IdCache().Exists(ctx, userId)
		if exist {
			continue
		}

		break
	}
	return userId, nil
}

func GetPassword(password string, salt string) string {
	return utils.Md5(fmt.Sprintf("%s_%s", password, salt))
}

func GetNickname(userId uint32) string {
	userIdStr := fmt.Sprintf("%d", userId)
	return "Player" + userIdStr[len(userIdStr)-4:]
}

func CreateUser(ctx context.Context, parameter *dto.CreateUserParameter) error {
	return db.CreateUser(ctx, parameter)
}

// RandString 随机字符串
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
