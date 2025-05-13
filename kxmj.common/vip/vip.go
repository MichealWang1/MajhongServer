package vip

type Config struct {
	Level          uint32 // VIP等级
	RequiredBP     uint32 // 需要经验
	LoginRewardMul uint32 // 登陆奖励倍数
}

var configs = map[uint32]*Config{
	0: {Level: 0, RequiredBP: 0, LoginRewardMul: 1},
	1: {Level: 1, RequiredBP: 100, LoginRewardMul: 2},
	2: {Level: 2, RequiredBP: 200, LoginRewardMul: 3},
	3: {Level: 3, RequiredBP: 300, LoginRewardMul: 4},
	4: {Level: 4, RequiredBP: 400, LoginRewardMul: 5},
	5: {Level: 5, RequiredBP: 500, LoginRewardMul: 6},
	6: {Level: 6, RequiredBP: 600, LoginRewardMul: 7},
	7: {Level: 7, RequiredBP: 700, LoginRewardMul: 8},
	8: {Level: 8, RequiredBP: 800, LoginRewardMul: 9},
	9: {Level: 9, RequiredBP: 900, LoginRewardMul: 10},
}

func GetConfig(level uint32) *Config {
	config, has := configs[level]
	if has == false {
		return nil
	}
	return config
}
