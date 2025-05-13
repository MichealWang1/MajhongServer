package model

// MailStatus 邮件状态
type MailStatus uint8

const (
	UnRead MailStatus = iota + 1
	Read
	Take
	Delete
)

// MailType 邮件类型 System系统邮件
type MailType uint8

const (
	Welfare MailType = iota + 1
	System
)

// 是否有奖励类型 HaveItem 有奖励
type IsRewardType uint8

const (
	HaveItem IsRewardType = iota + 1
)

// 是否单独发送 SingSend 单独发送
type IsSingleSendType int8

const (
	SingSend IsSingleSendType = iota + 1
)

type MailItem struct {
	Id    uint32 `json:"id"`    // 物品ID
	Count string `json:"count"` // 物品的数量
}

type MailData struct {
	EmailId   uint32      `json:"emailId"`   // 邮件ID
	EmailType uint8       `json:"emailType"` // 邮件类型：1 福利发放；2 系统通知
	Title     string      `json:"title"`     // 邮件标题
	Remark    string      `json:"remark"`    // 描述
	IsReward  uint8       `json:"isReward"`  // 是否有奖励：1 是；2 否
	ItemList  []*MailItem `json:"itemList"`  // 奖励物品 是 itemId
	Status    uint8       `json:"status"`    // 邮件状态 1.未读 2.已读 3.已领取
	CreatedAt uint32      `json:"createdAt"` // 邮件创建时间 到秒
}

type UserMailListResp struct {
	List []*MailData `json:"list"` // 邮件列表
}

type SetMailDataReq struct {
	EmailId uint32 `json:"emailId"` // 其中EmailId = 0则把所有邮件标记成已读
}

type TakeItemResp struct {
	ItemList []*MailItem `json:"itemList"` // 奖励物品 是 itemId
}

type Empty struct {
}
