package business

import (
	"context"
	"kxmj.common/entities/kxmj_logger"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/mq"
	"kxmj.common/utils"
	"kxmj.shop/internal/dto"
	"strconv"
	"time"
)

// AddUserGoldItems 添加用户钱包
func AddUserGoldItems(ctx context.Context, parameter *dto.AddParameter) error {
	// 更新用户钱包数据
	goldItems := item.GetGoldItems(parameter.Items)
	for _, g := range goldItems {
		diamond := "0"
		gold := "0"
		goldBean := "0"

		if g.ItemType == item.Diamond {
			diamond = g.Count
		} else if g.ItemType == item.Gold {
			gold = g.Count
		} else if g.ItemType == item.GoldBean {
			goldBean = g.Count
		}

		args := &center.AddUserWalletReq{
			UserId:       parameter.UserId,
			OrderId:      parameter.OrderId,
			Diamond:      diamond,
			Gold:         gold,
			GoldBean:     goldBean,
			BusinessType: parameter.BusinessType,
		}

		reply := &center.AddUserWalletResp{}
		err := parameter.CenterClient.Call(ctx, "AddUserWallet", args, reply)
		if err != nil {
			log.Sugar().Errorf("AddUserWallet user:%d order:%d err:%v", parameter.UserId, parameter.OrderId, err)
			return err
		}
	}
	return nil
}

// AddUserPropItems 添加用户背包道具
func AddUserPropItems(ctx context.Context, parameter *dto.AddParameter) error {
	// 更新用户背包物品
	propItems := item.GetPropItems(parameter.Items)
	err := item.UpdateUserItems(ctx, propItems, parameter.UserId)
	if err != nil {
		log.Sugar().Errorf("UpdateUserItems user:%d order:%d err:%v", parameter.UserId, parameter.OrderId, err)
		return err
	}

	// 写物品购买日志
	for _, p := range propItems {
		err = mq.AddLogger(&kxmj_logger.ItemTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      parameter.OrderId,
			UserId:       parameter.UserId,
			ItemId:       p.ItemId,
			Count:        p.Count,
			Type:         1,
			BusinessType: parameter.BusinessType,
			CreatedAt:    uint32(time.Now().Unix()),
		})
	}

	return nil
}

func AddUserBpItems(ctx context.Context, parameter *dto.AddBpParameter) (*dto.AddBpResult, error) {
	// 更新用户钱包数据
	items := item.GetBPItems(parameter.Items)
	result := &dto.AddBpResult{}
	for _, g := range items {
		bp, _ := strconv.Atoi(g.Count)
		args := &center.AddUserBpReq{
			UserId: parameter.UserId,
			BP:     uint32(bp),
		}

		reply := &center.AddUserBpResp{}
		err := parameter.CenterClient.Call(ctx, "AddUserBp", args, reply)
		if err != nil {
			log.Sugar().Errorf("AddUserBp user:%d order:%d err:%v", parameter.UserId, parameter.OrderId, err)
			return nil, err
		}
		result.UpgradeLevel = reply.Data.UpgradeLevel
	}
	return result, nil
}
