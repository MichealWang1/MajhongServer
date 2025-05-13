package server

import (
	"context"
	"kxmj.center/internal/db"
	"kxmj.center/internal/dto"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_logger"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/mq"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"time"
)

// GetUserInfo 获取用户信息
func (rs *RpcxServer) GetUserInfo(ctx context.Context, args *center.GetUserInfoReq, reply *center.GetUserInfoResp) error {
	info, err := redis_cache.GetCache().GetUserCache().DetailCache().Get(ctx, args.UserId)
	if err != nil {
		reply.Code = codes.GetUserInfoFailed
		reply.Msg = codes.GetMessage(codes.GetUserInfoFailed)
		log.Sugar().Errorf("GetUserInfo userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.GetUserInfoData{
		UserId:      info.UserId,
		Nickname:    info.Nickname,
		Gender:      info.Gender,
		AvatarAddr:  info.AvatarAddr,
		AvatarFrame: info.AvatarFrame,
		RealName:    info.RealName,
		UserMod:     info.UserMod,
		Vip:         info.Vip,
		TelNumber:   info.TelNumber,
		Status:      info.Status,
	}

	return nil
}

// CheckUserGold 检查用户金币数
func (rs *RpcxServer) CheckUserGold(ctx context.Context, args *center.CheckUserGoldReq, reply *center.CheckUserGoldResp) error {
	info, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, args.UserId)
	if err != nil {
		reply.Code = codes.CheckWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.CheckWalletInfoFailed)
		log.Sugar().Errorf("CheckUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.CheckUserGoldData{
		UserId: info.UserId,
		Gold:   info.Gold,
	}

	return nil
}

// GetUserGold 获取用户金币
func (rs *RpcxServer) GetUserGold(ctx context.Context, args *center.GetUserGoldReq, reply *center.GetUserGoldResp) error {
	info, err := db.GetUserGold(ctx, args.UserId)
	if err != nil {
		reply.Code = codes.GetWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.GetWalletInfoFailed)
		log.Sugar().Errorf("GetUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.GetUserGoldData{
		UserId: info.UserId,
		Gold:   info.Gold,
	}

	if utils.Cmp(info.Gold, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.GoldGameTransaction{
			Id:        utils.Snowflake.Generate().Int64(),
			GameId:    args.GameId,
			GameType:  args.GameType,
			UserId:    args.UserId,
			RoomId:    args.RoomId,
			RoomLevel: args.RoomLevel,
			Gold:      info.Gold,
			Type:      1,
			CreatedAt: uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}
	return nil
}

// SetUserGold 设置用户金币
func (rs *RpcxServer) SetUserGold(ctx context.Context, args *center.SetUserGoldReq, reply *center.SetUserGoldResp) error {
	if utils.Cmp(args.Gold, utils.Zero().String()) == 0 {
		return nil
	}

	info, err := db.SetUserGold(ctx, args.UserId, args.Gold)
	if err != nil {
		reply.Code = codes.UpdateWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.UpdateWalletInfoFailed)
		log.Sugar().Errorf("GetUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)

	// 增加日志
	err = mq.AddLogger(&kxmj_logger.GoldGameTransaction{
		Id:        utils.Snowflake.Generate().Int64(),
		GameId:    args.GameId,
		GameType:  args.GameType,
		UserId:    args.UserId,
		RoomId:    args.RoomId,
		RoomLevel: args.RoomLevel,
		Gold:      info.Gold,
		Type:      2,
		CreatedAt: uint32(time.Now().Unix()),
	})

	if err != nil {
		log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
	}

	return nil
}

// AddUserWallet 增加用户钱包数据
func (rs *RpcxServer) AddUserWallet(ctx context.Context, args *center.AddUserWalletReq, reply *center.AddUserWalletResp) error {
	result, err := db.AddUserWallet(ctx, &dto.AddUserWalletParameter{
		UserId:   args.UserId,
		Diamond:  args.Diamond,
		Gold:     args.Gold,
		GoldBean: args.GoldBean,
	})

	if err != nil {
		reply.Code = codes.UpdateWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.UpdateWalletInfoFailed)
		log.Sugar().Errorf("GetUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.AddUserWalletData{
		Diamond:  result.Diamond,
		Gold:     result.Gold,
		GoldBean: result.GoldBean,
	}

	if utils.Cmp(args.Gold, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.GoldTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         1,
			BusinessType: args.BusinessType,
			Count:        args.Gold,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}

	if utils.Cmp(args.Diamond, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.DiamondTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         1,
			BusinessType: args.BusinessType,
			Count:        args.Diamond,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}

	if utils.Cmp(args.GoldBean, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.GoldBeanTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         1,
			BusinessType: args.BusinessType,
			Count:        args.GoldBean,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}
	return nil
}

// SubUserWallet 扣除用户钱包数据
func (rs *RpcxServer) SubUserWallet(ctx context.Context, args *center.SubUserWalletReq, reply *center.SubUserWalletResp) error {
	result, err := db.SubUserWallet(ctx, &dto.SubUserWalletParameter{
		UserId:   args.UserId,
		Diamond:  args.Diamond,
		Gold:     args.Gold,
		GoldBean: args.GoldBean,
	})

	if err != nil {
		reply.Code = codes.UpdateWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.UpdateWalletInfoFailed)
		log.Sugar().Errorf("GetUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.SubUserWalletData{
		Diamond:  result.Diamond,
		Gold:     result.Gold,
		GoldBean: result.GoldBean,
	}

	if utils.Cmp(args.Gold, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.GoldTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         2,
			BusinessType: args.BusinessType,
			Count:        args.Gold,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}

	if utils.Cmp(args.Diamond, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.DiamondTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         2,
			BusinessType: args.BusinessType,
			Count:        args.Diamond,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}

	if utils.Cmp(args.GoldBean, utils.Zero().String()) > 0 {
		// 增加日志
		err = mq.AddLogger(&kxmj_logger.GoldBeanTransaction{
			Id:           utils.Snowflake.Generate().Int64(),
			OrderId:      args.OrderId,
			UserId:       args.UserId,
			Type:         2,
			BusinessType: args.BusinessType,
			Count:        args.GoldBean,
			CreatedAt:    uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("AddLogger userId:%v err:%v", args.UserId, err)
		}
	}

	return nil
}

// CheckUserDiamond 检查用户钻石数
func (rs *RpcxServer) CheckUserDiamond(ctx context.Context, args *center.CheckUserDiamondReq, reply *center.CheckUserDiamondResp) error {
	info, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, args.UserId)
	if err != nil {
		reply.Code = codes.CheckWalletInfoFailed
		reply.Msg = codes.GetMessage(codes.CheckWalletInfoFailed)
		log.Sugar().Errorf("CheckUserGold userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.CheckUserDiamondData{
		UserId:  info.UserId,
		Diamond: info.Diamond,
	}

	return nil
}

// AddUserBp 增加用户VIP经验值
func (rs *RpcxServer) AddUserBp(ctx context.Context, args *center.AddUserBpReq, reply *center.AddUserBpResp) error {
	resul, err := db.AddUserBp(ctx, &dto.AddVipBpParameter{
		UserId: args.UserId,
		BP:     args.BP,
	})

	if err != nil {
		reply.Code = codes.AddUserBpFailed
		reply.Msg = codes.GetMessage(codes.AddUserBpFailed)
		log.Sugar().Errorf("AddUserBp userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)
	reply.Data = &center.AddUserBpData{
		UpgradeLevel: resul.Level,
	}

	return nil
}

// AddRecharge 增加用户累计充值
func (rs *RpcxServer) AddRecharge(ctx context.Context, args *center.AddRechargeReq, reply *center.AddRechargeResp) error {
	err := db.AddRecharge(ctx, &dto.AddRechargeParameter{
		UserId: args.UserId,
		Amount: args.Amount,
	})

	if err != nil {
		reply.Code = codes.AddUserRechargeFailed
		reply.Msg = codes.GetMessage(codes.AddUserRechargeFailed)
		log.Sugar().Errorf("AddRecharge userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)

	return nil
}

// AddOnlyOnceGoods 添加唯一购买商品
func (rs *RpcxServer) AddOnlyOnceGoods(ctx context.Context, args *center.AddOnlyOnceGoodsReq, reply *center.AddOnlyOnceGoodsResp) error {
	err := db.AddOnlyOnceGoods(ctx, &dto.AddOnlyOnceGoodsParameter{
		UserId:  args.UserId,
		GoodsId: args.GoodsId,
	})

	if err != nil {
		reply.Code = codes.AddUserRechargeFailed
		reply.Msg = codes.GetMessage(codes.AddUserRechargeFailed)
		log.Sugar().Errorf("AddRecharge userId:%v err:%v", args.UserId, err)
		return nil
	}

	reply.Code = codes.Success
	reply.Msg = codes.GetMessage(codes.Success)

	return nil
}
