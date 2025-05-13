package recharge

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"kxmj.common/codes"
	"kxmj.common/item"
	"kxmj.common/recharge"
	"kxmj.common/redis_cache"
	"kxmj.common/redis_cache/goods/continue_goods"
	"kxmj.common/utils"
	"kxmj.common/web"
	"kxmj.core.api/internal/dto"
	"kxmj.core.api/internal/model"
	"kxmj.core.api/internal/server/business"
	"kxmj.core.api/internal/server/rpc"
)

// GetFirstGiftPack 获取首充礼包信息
// @Description RECHARGE
// @Tags RECHARGE
// @Summary 获取首充礼包信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetFistRechargeGiftPackResp} "请求成功"
// @Router	/recharge/get-first-pack [GET]
func (s *Service) GetFirstGiftPack(ctx *gin.Context) {
	userId := uint32(web.GetUserId(ctx))

	itemMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	var valItem *item.ValueItem
	for _, i := range itemMap {
		if i.ItemType != uint16(item.GiftPack) {
			continue
		}

		value := item.GetValueItem(i)
		if len(value.Extra) <= 0 {
			continue
		}

		gType, has := value.Extra[1]
		if has == false {
			continue
		}

		if gType == uint32(recharge.FirstRecharge) {
			valItem = value
			break
		}
	}

	if valItem == nil {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	goodsMap, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	var goodsId string
	var goodsCount string
	for _, goods := range goodsMap {
		if goods.ItemId == valItem.ItemId {
			goodsId = goods.GoodsId
			goodsCount = goods.RealCount
			break
		}
	}

	if len(goodsId) <= 0 {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	valItems, err := valItem.ParseBaseValueItems(goodsCount)
	if err != nil {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}
	var items []*recharge.Item
	for _, i := range valItems {
		items = append(items, &recharge.Item{
			Id:    i.ItemId,
			Count: i.Count,
		})
	}

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
			if goodsId == id {
				isBuy = true
				break
			}
		}
	}

	resp := &model.GetFistRechargeGiftPackResp{
		IsBuy:   isBuy,
		GoodsId: goodsId,
		Items:   items,
	}
	web.RespSuccess(ctx, resp)
}

// GetContinueGiftPack 首充连续领取礼包信息
// @Description RECHARGE
// @Tags RECHARGE
// @Summary 首充连续领取礼包信息
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Success 200 {object} web.Response{data=model.GetContinueGiftPackResp} "请求成功"
// @Router	/recharge/get-continue-pack [GET]
func (s *Service) GetContinueGiftPack(ctx *gin.Context) {
	userId := uint32(web.GetUserId(ctx))

	itemMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	var valItems []*item.ValueItem
	for _, i := range itemMap {
		if i.ItemType != uint16(item.GiftPack) {
			continue
		}

		value := item.GetValueItem(i)
		if len(value.Extra) <= 0 {
			continue
		}

		gType, has := value.Extra[1]
		if has == false {
			continue
		}

		if gType == uint32(recharge.ContinueRecharge) {
			_, has = value.Extra[2]
			// 如果商品扩展属性配置不正确，返回给前端
			if has == false {
				web.RespFailed(ctx, codes.GetItemConfigFailed)
				return
			}

			valItems = append(valItems, value)
		}
	}

	if len(valItems) <= 0 {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	goodsMap, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	// 获取已经购买到的商品状态
	continueGoodsStatus, _ := redis_cache.GetCache().GetGoodsCache().GetContinueCache().Get(ctx, userId)

	var goodsList []*model.ContinueGiftPack
	for _, goods := range goodsMap {
		for _, i := range valItems {
			if goods.ItemId == i.ItemId {
				witchPack := i.Extra[2]
				packItems := recharge.ContinueGiftPack[witchPack]
				for witchDay, pItems := range packItems {
					// 计算物品数量
					var tempItems []*recharge.Item
					if utils.Cmp(goods.RealCount, utils.Zero().String()) > 1 {
						for _, tempItem := range pItems {
							count, _ := utils.Mul(goods.RealCount, tempItem.Count)
							tempItems = append(tempItems, &recharge.Item{
								Id:    tempItem.Id,
								Count: count.String(),
							})
						}
					} else {
						tempItems = pItems
					}

					// 解析出基础物品
					var parseItems []*recharge.Item
					for _, tempItem := range tempItems {
						data := itemMap[tempItem.Id]
						values, err := item.GetValueItem(data).ParseBaseValueItems(tempItem.Count)
						if err != nil {
							web.RespFailed(ctx, codes.GetItemConfigFailed)
							return
						}
						for _, v := range values {
							parseItems = append(parseItems, &recharge.Item{
								Id:    v.ItemId,
								Count: v.Count,
							})
						}
					}

					var status uint32
					if len(continueGoodsStatus) > 0 {
						cachePack, has := continueGoodsStatus[goods.GoodsId]
						if has {
							for _, cp := range cachePack {
								if witchDay == cp.WitchDay {
									status = cp.Status
									break
								}
							}
						}
					}

					goodsList = append(goodsList, &model.ContinueGiftPack{
						WitchDay: witchDay,
						IsBuy:    false,
						GoodsId:  goods.GoodsId,
						Status:   status,
						Items:    parseItems,
					})
				}
				break
			}
		}
	}

	if len(goodsList) <= 0 {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	wallet, err := redis_cache.GetCache().GetUserCache().WalletCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.GetWalletInfoFailed)
		return
	}

	if len(wallet.OnlyOneGoods) > 0 {
		var buyGoods []string
		_ = json.Unmarshal([]byte(wallet.OnlyOneGoods), &buyGoods)
		for _, id := range buyGoods {
			for _, goods := range goodsList {
				if id == goods.GoodsId {
					goods.IsBuy = true
					break
				}
			}
		}
	}

	resp := &model.GetContinueGiftPackResp{
		List: goodsList,
	}
	web.RespSuccess(ctx, resp)
}

// TakeContinueGiftPack 领取首充连续领取礼包
// @Description RECHARGE
// @Tags RECHARGE
// @Summary 领取首充连续领取礼包
// @Param BundleId	header string true "dev.kxmj.com"
// @Param DeviceId	header string true "0033-0000-9999-9999-9999-1111"
// @Param Authorization	header string true "9cb9aaf4e28d5094fbca383255ecbafbbabb95c40c758843054480391f448ee8"
// @Param request body model.TakeContinueGiftPackReq true "JSON"
// @Success 200 {object} web.Response{data=model.TakeContinueGiftPackResp} "请求成功"
// @Router	/recharge/take-continue-pack [GET]
func (s *Service) TakeContinueGiftPack(ctx *gin.Context) {
	payload := &model.TakeContinueGiftPackReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		web.RespFailed(ctx, codes.ParamError)
		return
	}

	userId := uint32(web.GetUserId(ctx))
	// 获取已经购买到的商品状态
	continueGoodsStatus, err := redis_cache.GetCache().GetGoodsCache().GetContinueCache().Get(ctx, userId)
	if err != nil {
		web.RespFailed(ctx, codes.NotCanTakeItems)
		return
	}

	statusList, has := continueGoodsStatus[payload.GoodsId]
	if has == false {
		web.RespFailed(ctx, codes.NotCanTakeItems)
		return
	}

	var status *continue_goods.GoodsInfo
	for _, d := range statusList {
		if d.WitchDay == payload.WitchDay {
			status = d
			break
		}
	}

	if status == nil {
		web.RespFailed(ctx, codes.NotCanTakeItems)
		return
	}

	if status.Status == 0 {
		web.RespFailed(ctx, codes.NotCanTakeItems)
		return
	}

	if status.Status == 2 {
		web.RespFailed(ctx, codes.TheItemAlreadyTake)
		return
	}

	goodsMap, err := redis_cache.GetCache().GetGoodsCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	goods, has := goodsMap[payload.GoodsId]
	if has == false {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	itemMap, err := redis_cache.GetCache().GetItemCache().GetDetailCache().GetAll(ctx)
	if err != nil {
		web.RespFailed(ctx, codes.GetGoodsConfigFailed)
		return
	}

	data, has := itemMap[goods.ItemId]
	if has == false {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	valueItem := item.GetValueItem(data)
	if len(valueItem.Extra) <= 0 {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	pType, has := valueItem.Extra[2]
	if has == false {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	cfg, has := recharge.ContinueGiftPack[pType]
	if has == false {
		web.RespFailed(ctx, codes.GetItemConfigFailed)
		return
	}

	var items []*item.ValueItem
	for witchDay, vs := range cfg {
		if witchDay != payload.WitchDay {
			continue
		}

		for _, v := range vs {
			d, has := itemMap[v.Id]
			if has == false {
				web.RespFailed(ctx, codes.GetItemConfigFailed)
				return
			}
			val := item.GetValueItem(d)

			var count string
			if utils.Cmp(goods.RealCount, utils.Zero().String()) > 0 {
				temp, _ := utils.Mul(v.Count, goods.RealCount)
				count = temp.String()
			} else {
				count = v.Count
			}

			valItems, err := val.ParseBaseValueItems(count)
			if err != nil {
				web.RespFailed(ctx, codes.GetItemConfigFailed)
				return
			}
			items = append(items, valItems...)
		}
	}

	orderId := utils.Snowflake.Generate().Int64()
	err = business.AddUserGoldItems(ctx, &dto.AddParameter{
		UserId:       userId,
		OrderId:      orderId,
		BusinessType: 2,
		CenterClient: rpc.Default().CenterClient(),
		Items:        items,
	})
	if err != nil {
		web.RespFailed(ctx, codes.AddUserWalletFailed)
		return
	}

	err = business.AddUserPropItems(ctx, &dto.AddParameter{
		UserId:       userId,
		OrderId:      orderId,
		BusinessType: 2,
		CenterClient: rpc.Default().CenterClient(),
		Items:        items,
	})
	if err != nil {
		web.RespFailed(ctx, codes.AddUserItemFailed)
		return
	}

	_, err = business.AddUserBpItems(ctx, &dto.AddBpParameter{
		UserId:       userId,
		OrderId:      orderId,
		CenterClient: rpc.Default().CenterClient(),
		Items:        items,
	})
	if err != nil {
		web.RespFailed(ctx, codes.AddUserBpFailed)
		return
	}

	resp := &model.TakeContinueGiftPackResp{}
	web.RespSuccess(ctx, resp)
}
