package dto

import (
	"kxmj.shop/internal/model"
	"sync"
)

// ------------------------------------测试代码 start -----------------------------------
// 保存发给前端数据 因当前还没有把数据放入Redis中 因此 需要用到锁
// 后面直接去掉 把数据存入Redis
type SendGoodsData struct {
	goodsList []model.GoodsData
	listLock  sync.Mutex
}

// 创建一个保存 发送给前段的所有商品数据
func NewSendGoodsData() *SendGoodsData {
	return &SendGoodsData{
		goodsList: make([]model.GoodsData, 0),
		listLock:  sync.Mutex{},
	}
}

// 更新所有商品数据
func (g *SendGoodsData) UpdataGoods(list []model.GoodsData) {
	if len(list) <= 0 {
		return
	}
	g.listLock.Lock()
	defer g.listLock.Unlock()

	g.goodsList = list
}

//------------------------------------测试代码 end -----------------------------------
