package db

import (
	"context"
	"kxmj.common/entities/kxmj_report"
	"kxmj.common/mysql"
)

func CreateOrder(ctx context.Context, order *kxmj_report.OrderGoods) error {
	return mysql.ReportMaster().WithContext(ctx).Create(order).Error
}
