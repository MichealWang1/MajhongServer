package db

import (
	"context"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
)

func AddDevice(ctx context.Context, device *kxmj_core.Device) error {
	return mysql.CoreMaster().WithContext(ctx).Create(device).Error
}
