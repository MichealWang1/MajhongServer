package codes

const (
	GetUserVipInfoFailed          = 3001 // 获取用户VIP信息失败
	GetUserDailyLoginRewardFailed = 3002 // 获取用户每日登陆奖品信息失败
	GetUserTotalLoginRewardFailed = 3003 // 获取用户累计登陆奖品信息失败
	NotCanTakePrize               = 3004 // 没有可领取奖品
	UpdateDailyLoginCacheFailed   = 3005 // 更新每日登陆奖品缓存失败
	UpdateTotalLoginCacheFailed   = 3006 // 更新累计登陆奖品缓存失败
)

var taskMessage = map[int]string{
	GetUserVipInfoFailed:          "获取用户VIP信息失败",
	GetUserDailyLoginRewardFailed: "获取用户每日登陆奖品信息失败",
	GetUserTotalLoginRewardFailed: "获取用户累计登陆奖品信息失败",
	NotCanTakePrize:               "没有可领取奖品",
	UpdateDailyLoginCacheFailed:   "更新每日登陆奖品缓存失败",
	UpdateTotalLoginCacheFailed:   "更新累计登陆奖品缓存失败",
}
