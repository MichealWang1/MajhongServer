package business

import (
	"context"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/redis_cache"
	"kxmj.core.api/internal/db"
)

func ExistDevice(ctx context.Context, deviceId string) bool {
	return redis_cache.GetCache().GetDeviceCache().GetIdCache().Exists(ctx, deviceId)
}

func AddDevice(ctx context.Context, device *kxmj_core.Device) error {
	return db.AddDevice(ctx, device)
}

func GetDevice(ctx context.Context, deviceId string) (*kxmj_core.Device, error) {
	return redis_cache.GetCache().GetDeviceCache().GetDetailCache().Get(ctx, deviceId)
}
