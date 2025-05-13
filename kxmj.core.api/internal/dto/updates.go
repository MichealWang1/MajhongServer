package dto

import (
	"github.com/smallnest/rpcx/client"
	"kxmj.common/item"
)

type AddParameter struct {
	UserId       uint32            // 用户Id
	OrderId      int64             // 订单号
	BusinessType uint8             // 业务类型：1 商城；2 任务；3 邮件；
	CenterClient client.XClient    // 中心服务客户端
	Items        []*item.ValueItem // 物品数组
}

type AddBpParameter struct {
	UserId       uint32            // 用户Id
	OrderId      int64             // 订单号
	CenterClient client.XClient    // 中心服务客户端
	Items        []*item.ValueItem // 物品数组
}

type AddBpResult struct {
	UpgradeLevel uint32 // 增加后VIP等级
}

type AddItem struct {
	Id    uint32 `json:"id"`    // 物品ID
	Count string `json:"count"` // 物品数量
}
