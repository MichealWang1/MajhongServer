package ai

import (
	"sort"
)

// OpCombinations 最优组合
type OpCombinations []CombinationResult // 组合

type CombType int32

const (
	SingleCombs   CombType = iota + 1 // 单牌
	PairCombs                         // 对子
	StraightCombs                     // 顺子
	KeZiCombs                         // 刻子
)

type CombinationResult struct {
	Combination     []int32  // 组合牌
	CombinationType CombType // 组合牌类型
	IsConn          bool     // 是否相连
	Score           int32    // 组合分
}

func (op OpCombinations) SortToScore() {
	sort.Slice(op, func(i, j int) bool {
		return op[i].Score <= op[j].Score
	})
}

func AnalysisCombs(handTiles []int32, combs [][]int32, scores *ScoreConfig) OpCombinations {
	cpCombs := SortCombs(combs)
	res := make(OpCombinations, 0)
	// 拆分组合
	for _, comb := range cpCombs {
		res = append(res, AnalysisCombsType(comb)...)
	}
	// 判断是否相连
	for i := 0; i < len(res)-1; i++ {
		if res[i].CombinationType == res[i+1].CombinationType {
			if TileIsWTT(int(res[i+1].Combination[0])) {
				if res[i].Combination[0] == res[i+1].Combination[0] || res[i].Combination[0]+1 == res[i+1].Combination[0] {
					res[i].IsConn = true
					res[i+1].IsConn = true
				}
			}

		}
	}
	// 计算分数(牌型分+花色分)
	colorRank := TilesToColorSort(handTiles)
	for i := 0; i < len(res); i++ {
		score := int32(0)
		switch res[i].CombinationType { // 牌型
		case SingleCombs:
			score = scores.SingleScore
			if res[i].IsConn {
				score = scores.ConnSingleScore
			}
		case PairCombs:
			score = scores.PairScore
			if res[i].IsConn {
				score = scores.ConnPairScore
			}
		case StraightCombs:
			score = scores.StraightScore
		case KeZiCombs:
			score = scores.KeZiScore
			if res[i].IsConn {
				score = scores.ConnKeZiScore
			}
		}
		switch colorRank[res[i].Combination[0]/9] { // 花色
		case 0: // 主花
			score += scores.OneColorScore
		case 1: // 二花
			score += scores.TwoColorScore
		case 2: // 三花
			score += scores.ThreeColorScore
		case 3: // 风牌
			score += scores.FengColorScore
		}
		res[i].Score += score
	}
	res.SortToScore()
	return res
}

// AnalysisCombsType 分析组合类型
func AnalysisCombsType(combs []int32) OpCombinations {
	n := len(combs)
	res := make(OpCombinations, 0)
	switch n {
	case 1: // 单牌
		res = append(res, CombinationResult{Combination: []int32{combs[0]}, CombinationType: SingleCombs})
	case 2:
		if combs[0] == combs[1] { // 对子
			res = append(res, CombinationResult{Combination: []int32{combs[0], combs[1]}, CombinationType: PairCombs})
		} else {
			res = append(res, CombinationResult{Combination: []int32{combs[0]}, CombinationType: SingleCombs})
			res = append(res, CombinationResult{Combination: []int32{combs[1]}, CombinationType: SingleCombs})
		}
	case 3:
		if combs[0] == combs[1] { // 刻子
			res = append(res, CombinationResult{Combination: []int32{combs[0], combs[1], combs[2]}, CombinationType: KeZiCombs})
		} else { // 顺子
			res = append(res, CombinationResult{Combination: []int32{combs[0], combs[1], combs[2]}, CombinationType: StraightCombs})
		}
	}
	return res
}
