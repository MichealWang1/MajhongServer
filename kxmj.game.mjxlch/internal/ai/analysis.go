package ai

import (
	"fmt"
	lib "kxmj.common/mahjong"
)

// 获取下一个向听数(快速听牌)
func GetNextShanten(s int32) int32 {
	if s < 3 {
		return s - 2
	}
	return s - 1
}

// 3k+2 张牌，计算向听数、进张、改良、向听倒退等
func CalculateResultOf3N2(playerInfo *RobotEstimateHuArgs, ruleConfig *RuleConfig) lib.Cards {
	// 创建缓存向听map
	tm := make(map[string]CacheSC)
	// 获取当前手牌的向听数
	shanten, _ := CalculateMinShanten(playerInfo.HandTiles, ruleConfig, tm)
	fmt.Println("向听数：", shanten)
	// 计算目标向听数
	toShanten := GetNextShanten(shanten)
	// 构造搜索树
	n2 := BuildSearchNode3N2(shanten, toShanten, playerInfo.HandTiles, playerInfo.ResidueTiles, ruleConfig, tm)

	// 分析搜索树
	res := n2.AnalysisNode()
	return res.AnalysisOutCard()
	//return n2
}

// CalculateOpCombOf3N2 获取手牌最佳组合
func CalculateOpCombOf3N2(playerInfo *RobotEstimateHuArgs, ruleConfig *RuleConfig, scores *ScoreConfig) (OpCombinations, int32) {
	// 创建缓存向听map
	tm := make(map[string]CacheSC)
	// 获取当前手牌最佳组合
	shanten, combs := CalculateMinShanten(playerInfo.HandTiles, ruleConfig, tm)
	//fmt.Println("组合：", combs, "向听数:", shanten)
	fmt.Printf("向听数：%v\n", shanten)
	// 分析组合
	opCombs := AnalysisCombs(playerInfo.HandTiles, combs, scores)
	return opCombs, shanten
}
