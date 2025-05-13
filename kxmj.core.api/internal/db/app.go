package db

import (
	"context"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
)

func GetAppConfig(ctx context.Context) (*kxmj_core.ConfigApp, error) {
	value := &kxmj_core.ConfigApp{}
	err := mysql.CoreSlave().WithContext(ctx).Where("1 = 1").First(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
