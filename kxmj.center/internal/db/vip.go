package db

import (
	"context"
	"gorm.io/gorm"
	"kxmj.center/internal/dto"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache"
	"kxmj.common/vip"
	"time"
)

// AddUserBp 增加用户VIP经验值
func AddUserBp(ctx context.Context, parameter *dto.AddVipBpParameter) (*dto.AddVipBpResult, error) {
	userVip := &kxmj_core.UserVip{}
	err := mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", parameter.UserId).First(userVip).Error
		if err != nil {
			return err
		}

		nextConfig := vip.GetConfig(userVip.Level + 1)
		var isUpgrade bool
		if nextConfig != nil {
			// 如果可以升级
			if userVip.CurBp+parameter.BP >= nextConfig.RequiredBP {
				userVip.Level++
				isUpgrade = true
			}
		}
		userVip.CurBp += parameter.BP
		userVip.UpdatedAt = uint32(time.Now().Unix())

		err = tx.Where("user_id = ?", parameter.UserId).
			Select("level", "cur_bp", "updated_at").
			Updates(userVip).Error

		if err != nil {
			return err
		}

		// 数据一致性
		err = redis_cache.GetCache().GetUserCache().VIPCache().Del(ctx, parameter.UserId)
		if err != nil {
			return err
		}

		// 如果VIP发生升级，同步用户表
		if isUpgrade {
			err = tx.Where("user_id = ?", parameter.UserId).
				Select("vip", "updated_at").
				Updates(&kxmj_core.User{
					Vip:       uint8(userVip.Level),
					UpdatedAt: uint32(time.Now().Unix()),
				}).Error

			if err != nil {
				return err
			}

			// 数据一致性
			err = redis_cache.GetCache().GetUserCache().DetailCache().Del(ctx, parameter.UserId)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &dto.AddVipBpResult{Level: userVip.Level}, err
}
