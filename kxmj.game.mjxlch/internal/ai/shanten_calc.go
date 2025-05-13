package ai

import "fmt"

type CacheSC struct {
	shanten      uint8     // 最小向听数
	combinations [][]int32 // 最优组合牌
}

// CalculateMinShanten 计算最小向听数(handTiles 手牌信息,tm 缓存对应牌的向听数) / 最优组合
func CalculateMinShanten(handTiles []int32, ruleConfig *RuleConfig, tm map[string]CacheSC) (int32, [][]int32) {
	// 是否满足长度
	if CountOfTiles(handTiles) > 14 {
		fmt.Println("handLen > 14 is not allowed")
		return MaxShanten14, nil
	}
	// 构造手牌key值
	key := TilesToKey(handTiles)
	// 查询缓存里是否有
	if v, ok := tm[key]; ok {
		//fmt.Printf("key:%v,重复查询！！\n", key)
		return int32(v.shanten), v.combinations
	}
	// 构建向听结构体
	s := NewShanten(handTiles)
	// 计算七对向听
	if ruleConfig.CanHuQiDui {
		s.CalcQiDuiShanten()
	}
	// 分解手牌计算最小向听数
	s.DecomposeHandTile(0)
	// 拷贝最优组合
	combs := append(make([][]int32, 0, len(s.optimalCombination)), s.optimalCombination...)
	// 记录到缓存里
	tm[key] = CacheSC{
		shanten:      uint8(s.minShanten),
		combinations: combs,
	}
	return s.minShanten, combs
}
