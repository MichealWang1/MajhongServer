package goods

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/entities/kxmj_report"
	"kxmj.common/item"
	"kxmj.common/log"
	"kxmj.common/model/center"
	"kxmj.common/model/lobby"
	"kxmj.common/model/payment"
	"kxmj.common/redis_cache"
	"kxmj.common/utils"
	"kxmj.common/web"
	"kxmj.shop/internal/db"
	"kxmj.shop/internal/dto"
	"kxmj.shop/internal/model"
	"kxmj.shop/internal/server/business"
	"kxmj.shop/internal/server/rpc"
	"sort"
	"time"
)

// GetGoodsList 获取商品列表
// @Description GOODS
// @Tags GOODS
// @Summary 获取商品列表
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GoodsListResp} "请求成功"
// @Router	/goods/goods-list [GET]
func (s *Service) GetGoodsList(ctx *gin.Context) {
	//userId := web.GetUserId(ctx)

	resp := &model.GoodsListResp{}
	goodsMap, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	itemMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	var list []*kxmj_core.Goods
	for _, v := range goodsMap {
		if v.Status == 2 {
			continue
		}

		if v.IsDeleted == 1 {
			continue
		}

		itemInfo, has := itemMap[v.ItemId]
		if has == false {
			continue
		}

		// 复活卡类商品前端不显示
		if item.IsRise(item.Type(itemInfo.ItemType)) {
			continue
		}

		list = append(list, v)
	}

	// 排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})

	// 商品分类名
	categoryNames := map[uint8]string{
		1: "钻石",
		2: "金币",
		3: "装扮",
	}

	for _, v := range list {
		rewardAdded, _ := utils.Sub(v.RealCount, v.OriginalCount)
		mul, _ := utils.Mul(v.OriginalCount, "100")
		temp, _ := utils.Quo(mul.String(), v.RealCount)
		resp.List = append(resp.List, &model.GoodsData{
			GoodsId:        v.GoodsId,
			Name:           v.Name,
			ItemId:         v.ItemId,
			ShopType:       v.ShopType,
			Price:          v.Price,
			OriginalPrice:  v.OriginalPrice,
			RealCount:      v.RealCount,
			OriginalCount:  v.OriginalCount,
			RewardAdded:    rewardAdded.String(),
			IncomeTimes:    fmt.Sprintf("%f", float64(temp.Int64())/100),
			Recommend:      v.Recommend,
			FirstBuyDouble: v.FirstBuyDouble,
			ExpireTime:     v.ExpireTime,
			Category:       v.Category,
			CategoryName:   categoryNames[v.Category],
		})
	}
	web.RespSuccess(ctx, resp)
}

// Buy 商品购买
// @Description GOODS
// @Tags GOODS
// @Summary 商品购买
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.BuyReq true "JSON"
// @Success 200 {object} web.Response{data=model.BuyResp} "请求成功"
// @Router	/goods/buy [POST]
func (s *Service) Buy(ctx *gin.Context) {
	payload := &model.BuyReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.UnMarshalJsonErr)
		return
	}

	userId := uint32(web.GetUserId(ctx))
	wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.GetWalletInfoFailed)
		return
	}

	var isBuy bool
	if len(wallet.OnlyOneGoods) > 0 {
		var buyGoods []string
		_ = json.Unmarshal([]byte(wallet.OnlyOneGoods), &buyGoods)
		for _, id := range buyGoods {
			if payload.GoodsId == id {
				isBuy = true
				break
			}
		}
	}
	if isBuy {
		web.RespFailed(ctx, codes.GoodsOnlyBuyOnce)
		return
	}

	goodsMap, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	goods, has := goodsMap[payload.GoodsId]
	if has == false {
		web.RespFailed(ctx, codes.GoodsIsDeleteFailed)
		return
	}

	if goods.ShopType < 1 || goods.ShopType > 4 {
		web.RespFailed(ctx, codes.GoodsShopTypeError)
		return
	}

	var resp *model.BuyResp
	if goods.ShopType == 1 {
		reply := &payment.CreateOrderResp{}
		err = rpc.Default().PaymentClient().Call(ctx, "CreateOrder", &payment.CreateOrderReq{
			UserId:  userId,
			GoodsId: payload.GoodsId,
			Type:    payload.Type,
		}, reply)

		if err != nil {
			log.Sugar().Errorf("CreateOrder user:%d goods:%s err:%v", userId, payload.GoodsId, err)
			web.RespFailed(ctx, codes.ServerNetErr)
			return
		}

		if reply.Code != codes.Success {
			log.Sugar().Errorf("CreateOrder user:%d goods:%s err:%v", userId, payload.GoodsId, codes.New(reply.Code, reply.Msg))
			web.RespFailed(ctx, reply.Code, reply.Msg)
			return
		}

		resp = &model.BuyResp{
			OrderId: fmt.Sprintf("%d", reply.Data.OrderId),
			PayUrl:  reply.Data.PayUrl,
		}
	} else {
		wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
		if goods.ShopType == 2 {
			if utils.Cmp(wallet.Diamond, goods.Price) < 0 {
				web.RespFailed(ctx, codes.InsufficientDiamond)
				return
			}
		} else if goods.ShopType == 3 {
			if utils.Cmp(wallet.Gold, goods.Price) < 0 {
				web.RespFailed(ctx, codes.InsufficientGold)
				return
			}
		} else if goods.ShopType == 4 {
			if utils.Cmp(wallet.GoldBean, goods.Price) < 0 {
				web.RespFailed(ctx, codes.InsufficientGoldBean)
				return
			}
		}

		itemMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
		if err != nil {
			web.RespFailed(ctx, codes.GetGoodsConfigFailed)
			return
		}

		dataItem, has := itemMap[goods.ItemId]
		if has == false {
			web.RespFailed(ctx, codes.GoodsIsDeleteFailed)
			return
		}

		orderId := utils.Snowflake.Generate().Int64()
		args := &center.SubUserWalletReq{
			UserId:       userId,
			OrderId:      orderId,
			Diamond:      "0",
			Gold:         "0",
			GoldBean:     "0",
			BusinessType: 1,
		}

		if goods.ShopType == 2 {
			args.Diamond = goods.Price
		} else if goods.ShopType == 3 {
			args.Gold = goods.Price
		} else if goods.ShopType == 4 {
			args.GoldBean = goods.Price
		}

		reply := &center.SubUserWalletResp{}
		err = rpc.Default().CenterClient().Call(ctx, "SubUserWallet", args, reply)
		if err != nil {
			log.Sugar().Errorf("SubUserWallet user:%d goods:%s err:%v", userId, payload.GoodsId, err)
			web.RespFailed(ctx, codes.SubUserGoldFailed)
			return
		}

		if reply.Code != codes.Success {
			web.RespFailed(ctx, reply.Code, reply.Msg)
			return
		}

		err = db.CreateOrder(ctx, &kxmj_report.OrderGoods{
			OrderId:       orderId,
			TradeId:       "",
			GoodsId:       goods.GoodsId,
			UserId:        userId,
			ShopType:      goods.ShopType,
			ItemId:        goods.ItemId,
			ItemType:      dataItem.ItemType,
			ItemCount:     goods.RealCount,
			Price:         goods.Price,
			OriginalPrice: goods.OriginalPrice,
			PaymentType:   payload.Type,
			OrderStatus:   3,
			CompletedAt:   uint32(time.Now().Unix()),
			CreatedAt:     uint32(time.Now().Unix()),
			UpdatedAt:     uint32(time.Now().Unix()),
		})

		if err != nil {
			log.Sugar().Errorf("CreateOrder user:%d goods:%s err:%v", userId, payload.GoodsId, err)
			web.RespFailed(ctx, codes.CreateGoodsOrderFailed)
			return
		}

		// 把物品数据库定义类型转成业务值类型
		valueItem := item.GetValueItem(dataItem)
		// 解析物品基础类型
		values, err := valueItem.ParseBaseValueItems(goods.RealCount)
		if err != nil {
			log.Sugar().Errorf("ParseBaseValueItems user:%d goods:%s order:%d err:%v", userId, orderId, orderId, err)
			web.RespFailed(ctx, codes.ParseBaseValueItemFailed)
			return
		}

		err = business.AddUserGoldItems(ctx, &dto.AddParameter{
			UserId:       userId,
			OrderId:      orderId,
			BusinessType: 1,
			CenterClient: rpc.Default().CenterClient(),
			Items:        values,
		})
		if err != nil {
			log.Sugar().Errorf("AddUserGoldItems err:%v", err)
		}

		err = business.AddUserPropItems(ctx, &dto.AddParameter{
			UserId:       userId,
			OrderId:      orderId,
			BusinessType: 1,
			CenterClient: rpc.Default().CenterClient(),
			Items:        values,
		})
		if err != nil {
			log.Sugar().Errorf("AddUserPropItems err:%v", err)
		}

		_, err = business.AddUserBpItems(ctx, &dto.AddBpParameter{
			UserId:       userId,
			OrderId:      orderId,
			CenterClient: rpc.Default().CenterClient(),
			Items:        values,
		})
		if err != nil {
			log.Sugar().Errorf("AddUserBpItems err:%v", err)
		}

		shoppingReply := &lobby.ShoppingSuccessResp{}
		// 通知用户购买成功
		err = rpc.Default().LobbyClient().Call(ctx, "ShoppingSuccess", &lobby.ShoppingSuccessReq{
			UserId:         userId,
			GoodsId:        payload.GoodsId,
			Name:           goods.Name,
			ShopType:       goods.ShopType,
			Price:          goods.Price,
			RealCount:      goods.RealCount,
			FirstBuyDouble: goods.FirstBuyDouble,
			ItemId:         goods.ItemId,
		}, shoppingReply)

		if err != nil {
			log.Sugar().Errorf("ShoppingSuccess user:%d goods:%s order:%d err:%v", userId, goods.GoodsId, orderId, err)
		}
	}

	web.RespSuccess(ctx, resp)
}
