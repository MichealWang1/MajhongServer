package game_core

import (
	"fmt"
	"kxmj.common/utils"
	"math/rand"
	"time"
)

type RobotInfo struct {
	UserId      uint32 // 用户ID
	Nickname    string // 昵称
	Gender      uint8  // 性别
	AvatarAddr  string // 头像地址
	AvatarFrame uint8  // 头像框
	UserMod     uint8  // 人物样式
	Vip         uint8  // VIP等级
	Gold        string // 金币
}

func GetNickname(userId uint32) string {
	userIdStr := fmt.Sprintf("%d", userId)
	return "Player" + userIdStr[len(userIdStr)-4:]
}

var MaxRobotId = 20000000

func createRobotId() uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint32(r.Intn(10000000) + 10000000)
}

func CreateRobot(baseScore string) *RobotInfo {
	multiple := rand.Intn(1000000) + 100
	gold, _ := utils.Mul(baseScore, fmt.Sprintf("%d", multiple))
	userId := createRobotId()
	return &RobotInfo{
		UserId:      userId,
		Nickname:    GetNickname(userId),
		Gender:      uint8(rand.Intn(2) + 1),
		AvatarAddr:  "",
		AvatarFrame: 0,
		UserMod:     0,
		Vip:         uint8(rand.Intn(6)),
		Gold:        gold.String(),
	}
}
