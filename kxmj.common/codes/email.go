package codes

const (
	NoMail               = 10000 // 玩家没有邮件
	NotMailID            = 10001 // 邮件ID不存在
	EmailChangeStateFail = 10002 // 邮件改变状态失败
	NoTakeItemByMail     = 10003 // 没有邮件可以领取物品
	EmailAlreadyTake     = 10004 // 邮件已领取物品
	EmailAlreadyDelete   = 10005 // 邮件已删除
	EmailTakeGoldFail    = 10006 // 邮件领取金币失败
	EmailTakeItemFail    = 10007 // 邮件领取物品失败
	ReceiveMailStatus    = 10008 // 服务端收到设置邮件状态不正确
	EmailStatusError     = 10009 // 邮件状态错误
	NoChangeMailStatus   = 10010 // 邮件不需要改变状态
)

var emailMessage = map[int]string{
	NoMail:               "玩家没有邮件",
	NotMailID:            "邮件ID不存在",
	EmailChangeStateFail: "改变状态失败",
	NoTakeItemByMail:     "没有邮件可以领取物品",
	EmailAlreadyTake:     "邮件已领取物品",
	EmailAlreadyDelete:   "邮件已删除",
	EmailTakeGoldFail:    "邮件领取金币失败",
	EmailTakeItemFail:    "邮件领取物品失败",
	ReceiveMailStatus:    "收到设置邮件状态不正确",
	EmailStatusError:     "邮件状态错误",
	NoChangeMailStatus:   "邮件不需要改变状态",
}
