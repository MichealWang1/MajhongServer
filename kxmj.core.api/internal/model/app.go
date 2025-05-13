package model

import "kxmj.common/item"

type GetGatewayResp struct {
	List []*GetGatewayInfo `json:"list"`
}

type GetGatewayInfo struct {
	SvrType uint16 `json:"svrType"`
	SvrId   uint16 `json:"svrId"`
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
}

type GetAppBaseInfoResp struct {
	WechatSecretKey string `json:"wechatSecretKey"` // 微信登陆应用密钥APPSecret
	WechatAppId     string `json:"wechatAppId"`     // 微信开放平台应用唯一标识
	HotRenewAddress string `json:"hotRenewAddress"` // 热更新包地址
}

type AndroidInfo struct {
	SDK  uint8  `json:"sdk"`  // android sdk 版本
	ID   string `json:"id"`   // android id
	IMEI string `json:"imei"` // 自定义imei
}

type SyncDeviceReq struct {
	DeviceId     string       `json:"deviceId" binding:"required"` // 设备ID
	OS           uint8        `json:"os"`                          // 系统: 0 未知 1 安卓 2 IOS 3 其它
	Brand        string       `json:"brand"`                       // 品牌
	Manufacturer string       `json:"manufacturer"`                // 制造商
	Version      string       `json:"version"`                     // 系统版本号
	Model        string       `json:"model"`                       // 型号
	Width        uint32       `json:"width"`                       // 宽度
	Height       uint32       `json:"height"`                      // 高度
	AndroidInfo  *AndroidInfo `json:"androidInfo"`                 // android 设备信息
	IosUUID      string       `json:"iosUUID"`                     // IOS设备ID
	Organic      uint8        `json:"organic"`                     // organic 1 自然、2 广告
}

type SyncDeviceResp struct {
	DeviceId string `json:"deviceId"` // 设备Id
}

type HomeUser struct {
	Nickname    string `json:"nickname"`    // 昵称
	Gender      uint8  `json:"gender"`      // 性别：0 女；1 男；
	AvatarAddr  string `json:"avatarAddr"`  // 头像地址
	AvatarFrame uint8  `json:"avatarFrame"` // 头像框
	Diamond     string `json:"diamond"`     // 钻石数
	Gold        string `json:"gold"`        // 金币数
	GoldBean    string `json:"goldBean"`    // 金豆数
}

type GetHomeResp struct {
	User   *HomeUser      `json:"user"`   // 用户首页信息
	Guides map[int]uint32 `json:"guides"` // 首页引导提示列表(key：类型：1 商城；2 背包；3 活动；4 福利；5 直播；6 签到任务；7 对局任务；8 赢金任务；9 充值任务； value：提示数量)
}

type ItemData struct {
	ItemId        uint32                  `json:"itemId"`        // 物品ID
	Name          string                  `json:"name"`          // 物品名称
	ItemType      uint16                  `json:"itemType"`      // 物品类型： 101 钻石；102 金币；103 金豆；201 装扮；202 武器 203 头像框；204 牌桌装饰品类；205 牌背装饰品类；206 入场动画类；207 交互道具类 301 特权包；302 礼包；
	ServiceLife   uint32                  `json:"serviceLife"`   // 使用寿命（秒为单位）
	Content       []*item.GiftPackContent `json:"content"`       // 礼包、特权卡类道具内容
	Extra         map[uint32]uint32       `json:"extra"`         // 扩展属性(攻击、复活等属性)
	GiftType      uint8                   `json:"giftType"`      // 礼包类型：0 未定义；1 充值礼包；2 钻石礼包；3 抽奖礼包；
	AdornmentType uint8                   `json:"adornmentType"` // 装扮物品类型：1 头部；2 衣服；
}

type GetItemResp struct {
	List []*ItemData
}
