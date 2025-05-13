package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"kxmj.center/internal/dto"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/mysql"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"time"
)

// GetUserGold 获取用户金币
func GetUserGold(ctx context.Context, userId uint32) (*kxmj_core.UserWallet, error) {
	query := &kxmj_core.UserWallet{}
	err := mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", userId).First(query).Error
		if err != nil {
			return err
		}

		update := &kxmj_core.UserWallet{
			Gold:      "0",
			UpdatedAt: uint32(time.Now().Unix()),
		}

		err = tx.Where("user_id = ?", userId).Select("gold", "updated_at").Updates(update).Error
		if err != nil {
			return err
		}

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, userId)
	})

	return query, err
}

// SetUserGold 设置用户金币
func SetUserGold(ctx context.Context, userId uint32, gold string) (*kxmj_core.UserWallet, error) {
	query := &kxmj_core.UserWallet{}
	err := mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", userId).First(query).Error
		if err != nil {
			return err
		}

		val, ok := utils.AddToString(query.Gold, gold)
		if ok == false {
			return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.Gold, gold))
		}

		query.Gold = val
		query.UpdatedAt = uint32(time.Now().Unix())
		err = tx.Where("user_id = ?", userId).Select("gold", "updated_at").Updates(query).Error
		if err != nil {
			return err
		}

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, userId)
	})

	return query, err
}

// AddUserWallet 增加用户钱包数据
func AddUserWallet(ctx context.Context, parameter *dto.AddUserWalletParameter) (*kxmj_core.UserWallet, error) {
	query := &kxmj_core.UserWallet{}
	err := mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", parameter.UserId).First(query).Error
		if err != nil {
			return err
		}

		if utils.Cmp(parameter.Diamond, utils.Zero().String()) > 0 {
			val, ok := utils.AddToString(query.Diamond, parameter.Diamond)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.Diamond, parameter.Diamond))
			}
			query.Diamond = val
		}

		if utils.Cmp(parameter.Gold, utils.Zero().String()) > 0 {
			val, ok := utils.AddToString(query.Gold, parameter.Gold)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.Gold, parameter.Gold))
			}
			query.Gold = val
		}

		if utils.Cmp(parameter.GoldBean, utils.Zero().String()) > 0 {
			val, ok := utils.AddToString(query.GoldBean, parameter.GoldBean)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.GoldBean, parameter.GoldBean))
			}
			query.GoldBean = val
		}

		query.UpdatedAt = uint32(time.Now().Unix())
		err = tx.Where("user_id = ?", parameter.UserId).
			Select("diamond", "gold", "gold_bean", "updated_at").
			Updates(query).Error

		if err != nil {
			return err
		}

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, parameter.UserId)
	})

	return query, err
}

// SubUserWallet 扣除用户钱包数据
func SubUserWallet(ctx context.Context, parameter *dto.SubUserWalletParameter) (*kxmj_core.UserWallet, error) {
	query := &kxmj_core.UserWallet{}
	err := mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", parameter.UserId).First(query).Error
		if err != nil {
			return err
		}

		if utils.Cmp(parameter.Diamond, utils.Zero().String()) > 0 {
			val, ok := utils.SubToString(query.Diamond, parameter.Diamond)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.Diamond, parameter.Diamond))
			}
			query.Diamond = val
		}

		if utils.Cmp(parameter.Gold, utils.Zero().String()) > 0 {
			val, ok := utils.SubToString(query.Gold, parameter.Gold)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.Gold, parameter.Gold))
			}
			query.Gold = val
		}

		if utils.Cmp(parameter.GoldBean, utils.Zero().String()) > 0 {
			val, ok := utils.SubToString(query.GoldBean, parameter.GoldBean)
			if ok == false {
				return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.GoldBean, parameter.GoldBean))
			}
			query.GoldBean = val
		}

		query.UpdatedAt = uint32(time.Now().Unix())
		err = tx.Where("user_id = ?", parameter.UserId).
			Select("diamond", "gold", "gold_bean", "updated_at").
			Updates(query).Error

		if err != nil {
			return err
		}

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, parameter.UserId)
	})

	return query, err
}

// AddRecharge 扣除用户钱包数据
func AddRecharge(ctx context.Context, parameter *dto.AddRechargeParameter) error {
	return mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := &kxmj_core.UserWallet{}
		err := tx.Where("user_id = ?", parameter.UserId).First(query).Error
		if err != nil {
			return err
		}

		val, ok := utils.AddToString(query.TotalRecharge, parameter.Amount)
		if ok == false {
			return errors.New(fmt.Sprintf("AddToString s1:%s s2:%s err", query.TotalRecharge, parameter.Amount))
		}
		query.TotalRecharge = val
		query.RechargeTimes += 1

		query.UpdatedAt = uint32(time.Now().Unix())
		err = tx.Where("user_id = ?", parameter.UserId).
			Select("total_recharge", "recharge_times", "updated_at").
			Updates(query).Error

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, parameter.UserId)
	})
}

// AddOnlyOnceGoods 添加唯一购买商品
func AddOnlyOnceGoods(ctx context.Context, parameter *dto.AddOnlyOnceGoodsParameter) error {
	return mysql.CoreMaster().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := &kxmj_core.UserWallet{}
		err := tx.Where("user_id = ?", parameter.UserId).First(query).Error
		if err != nil {
			return err
		}

		var goodsList []string
		if len(query.OnlyOneGoods) > 0 {
			err = json.Unmarshal([]byte(query.OnlyOneGoods), &goodsList)
			if err != nil {
				return err
			}
		}

		var exist bool
		for _, id := range goodsList {
			if id == parameter.GoodsId {
				exist = true
				break
			}
		}

		if exist {
			return nil
		}

		goodsList = append(goodsList, parameter.GoodsId)
		data, err := json.Marshal(goodsList)
		if err != nil {
			return err
		}

		query.OnlyOneGoods = string(data)
		query.UpdatedAt = uint32(time.Now().Unix())
		err = tx.Where("user_id = ?", parameter.UserId).
			Select("only_one_goods", "updated_at").
			Updates(query).Error

		return redis_cache.GetCache().GetUserCache().WalletCache().Del(ctx, parameter.UserId)
	})
}
