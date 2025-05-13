package db

import (
	"context"
	"gorm.io/gorm"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache"
	"kxmj.core.api/internal/dto"
	"time"
)

func ChangePassword(ctx context.Context, userId uint32, NewPassword string) error {
	return mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &kxmj_core.User{
			LoginPassword: NewPassword,
			UpdatedAt:     uint32(time.Now().Unix()),
		}

		err := tx.Where("user_id = ?", userId).Select("login_password", "updated_at").Updates(user).Error
		if err != nil {
			return err
		}

		// 保证数据一致性，先删redis
		return redis_cache.GetCache().GetUserCache().DetailCache().Del(ctx, userId)
	})
}

func BindPhoneNum(ctx context.Context, telNumber string, userID uint32) error {
	return mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &kxmj_core.User{
			UserId:    userID,
			TelNumber: telNumber,
			UpdatedAt: uint32(time.Now().Unix()),
		}

		err := tx.Where("user_id = ?", userID).
			Select("tel_number", "updated_at").
			Updates(user).Error

		if err != nil {
			return err
		}

		// 保证数据一致性，先删redis
		err = redis_cache.GetCache().GetUserCache().DetailCache().Del(ctx, userID)
		if err != nil {
			return err
		}

		// 保证数据一致性，先同步redis
		return redis_cache.GetCache().GetUserCache().TelCache().Set(ctx, telNumber, userID)
	})
}

func CreateUser(ctx context.Context, parameter *dto.CreateUserParameter) error {
	return mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &kxmj_core.User{
			Id:                parameter.Id,
			UserId:            parameter.UserId,
			Nickname:          parameter.Nickname,
			Gender:            parameter.Gender,
			AvatarAddr:        parameter.AvatarAddr,
			AvatarFrame:       parameter.AvatarFrame,
			RealName:          parameter.RealName,
			IdCard:            parameter.IdCard,
			UserMod:           parameter.UserMod,
			AccountType:       parameter.AccountType,
			Vip:               parameter.Vip,
			DeviceId:          parameter.DeviceId,
			RegisterIp:        parameter.RegisterIp,
			RegisterType:      parameter.RegisterType,
			TelNumber:         parameter.TelNumber,
			Status:            parameter.Status,
			BindingAt:         parameter.BindingAt,
			LoginPassword:     parameter.LoginPassword,
			LoginPasswordSalt: parameter.LoginPasswordSalt,
			Remark:            parameter.Remark,
			BundleId:          parameter.BundleId,
			BundleChannel:     parameter.BundleChannel,
			Organic:           parameter.Organic,
			CreatedAt:         parameter.CreatedAt,
			UpdatedAt:         parameter.UpdatedAt,
		}

		err := tx.Create(user).Error
		if err != nil {
			return err
		}

		userThirdParty := &kxmj_core.UserThirdParty{
			Id:           parameter.Id,
			UserId:       parameter.UserId,
			DeviceId:     parameter.DeviceId,
			WechatOpenId: parameter.WechatOpenId,
			TiktokId:     parameter.TiktokId,
			HuaweiId:     parameter.HuaweiId,
			UpdatedAt:    parameter.UpdatedAt,
		}
		err = tx.Create(userThirdParty).Error
		if err != nil {
			return err
		}

		userWallet := &kxmj_core.UserWallet{
			Id:            parameter.Id,
			UserId:        parameter.UserId,
			Diamond:       parameter.Diamond,
			Gold:          parameter.Gold,
			GoldBean:      parameter.GoldBean,
			TotalRecharge: parameter.TotalRecharge,
			RechargeTimes: parameter.RechargeTimes,
			UpdatedAt:     parameter.UpdatedAt,
		}
		err = tx.Create(userWallet).Error
		if err != nil {
			return err
		}

		equip := &kxmj_core.UserEquip{
			Id:        parameter.Id,
			UserId:    parameter.UserId,
			Head:      parameter.Head,
			Body:      parameter.Body,
			Weapon:    parameter.Weapon,
			UpdatedAt: parameter.UpdatedAt,
		}
		err = tx.Create(equip).Error
		if err != nil {
			return err
		}

		vip := &kxmj_core.UserVip{
			Id:        parameter.Id,
			UserId:    parameter.UserId,
			Level:     0,
			CurBp:     0,
			UpdatedAt: parameter.UpdatedAt,
		}
		err = tx.Create(vip).Error
		if err != nil {
			return err
		}

		// 保证数据一致性，先同步redis
		err = redis_cache.GetCache().GetUserCache().IdCache().Set(ctx, parameter.UserId)
		if err != nil {
			return err
		}

		if len(parameter.TelNumber) > 0 {
			// 保证数据一致性，先同步redis
			err = redis_cache.GetCache().GetUserCache().TelCache().Set(ctx, parameter.TelNumber, parameter.UserId)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
