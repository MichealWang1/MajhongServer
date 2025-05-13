package db

import (
	"context"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
)

func GetBundle(ctx context.Context, bundleId string) (*kxmj_core.ConfigBundle, error) {
	query := &kxmj_core.ConfigBundle{}
	err := mysql.CoreSlave().WithContext(ctx).Where("bundle_id = ?", bundleId).First(query).Error
	return query, err
}
