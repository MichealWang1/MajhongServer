package business

import (
	"context"
	"kxmj.common/log"
	"kxmj.shop/internal/db"
	"kxmj.shop/internal/server/dto"
)

var (
	sendCliendShopData *dto.SendGoodsData
)

func init() {
	sendCliendShopData = dto.NewSendGoodsData()
}

// 处理商城数据
func HandleShopData() {
	// 创建 context
	ctx := context.Background()
	// 获取商品列表
	goodlist, err := db.GetGoodsShowList(ctx)
	if err != nil {
		log.Sugar().Error(" 获取数据库商品列表失败")
		return
	}
	var sendList = []dto.GoodsData{}
	// 把数据整合发送给前段
	for _, goods := range goodlist {
		temp := dto.GoodsData{
			GoodsId:        goods.GoodsId,
			Name:           goods.Name,
			Remark:         goods.Remark,
			ShopType:       goods.ShopType,
			ItemId:         goods.ItemId,
			Price:          goods.Price,
			OriginalPrice:  goods.OriginalPrice,
			RealCount:      goods.RealCount,
			OriginalCount:  goods.OriginalCount,
			RewardAdded:    goods.RewardAdded,
			IncomeTimes:    goods.IncomeTimes,
			Recommend:      goods.Recommend,
			FirstBuyDouble: goods.FirstBuyDouble,
			Status:         goods.Status,
			ExpireTime:     goods.ExpireTime,
			OnShelfTime:    goods.OnShelfTime,
			OffShelfTime:   goods.OffShelfTime,
			Sort:           goods.Sort,
		}
		sendList = append(sendList, temp)
	}
	if len(sendList) > 0 {
		sendCliendShopData.UpdataGoods(sendList)
	}
}
