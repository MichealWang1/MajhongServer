package ai

import (
	lib "kxmj.common/mahjong"
	"sort"
)

type Hand3N2AnalysisResultList []Hand3N2AnalysisResult

// 分析整个搜索数，得出最佳的路径
type Hand3N2AnalysisResult struct {
	DiscardTile      int32                 // 需要出的牌
	DiscardTileValue float64               // 牌的价值
	Result13         Hand3N1AnalysisResult // 出牌后的分析结果
}

// 3k+1 张手牌的分析结果
type Hand3N1AnalysisResult struct {
	HandTiles    []int32         // 原手牌
	ResidueTiles []int32         // 剩余牌
	Waits        map[int32]int32 // 向听前进牌--> 进张
	WaitsCount   int32           // 能进张的数量
	Shanten      int32           // 向听数

	//AvgImproveWaitsCount     float64         // 摸到非进张牌时的进张数的加权均值
	//AvgNextShantenWaitsCount float64         // 向听前进后的(最大)进张数的加权均值
	//MixedWaitsScore          float64         // 综合了进张与向听前进后进张的评分
	//AvgAgariRate             float64         // 听牌时的手牌和率
	//AvgHuPoint               float64         // 听牌时的平均胡牌分
	//HuTypes                  []string        // 胡牌类型
}

func (a Hand3N2AnalysisResultList) AnalysisOutCard() lib.Cards {
	res := make(lib.Cards, 0, len(a))
	sort.Slice(a, func(i, j int) bool {
		return a[i].DiscardTileValue > a[j].DiscardTileValue
	})
	for _, v := range a {
		res = append(res, TileToCardMap[v.DiscardTile])
	}
	return res
}
