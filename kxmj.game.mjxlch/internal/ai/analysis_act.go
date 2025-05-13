package ai

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
)

// 行为分析

// AnalysisPlayerOutCard 出牌行为分析
func AnalysisPlayerOutCard(robotArgs *EstimateArgs) lib.Card {
	// 构造分析参数
	args := NewRobotEstimateHuArgs()
	args.Build(robotArgs)
	// 手牌是3n+2才能出牌
	if args.HandCardsLen%3 != 2 {
		return lib.INVALID_CARD
	}
	outCard := lib.INVALID_CARD
	if args.MustCanOutCards.Len() > 0 {
		return args.MustCanOutCards[0]
	}

	// 根据机器人的等级
	switch robotArgs.Level {
	case RobotLevel1:
		// 获取最佳牌型相关信息
		n2, _ := CalculateOpCombOf3N2(args, robotArgs.RuleConfig, robotArgs.ScoreConfig)
		//if shanten == -1 {
		//	return lib.INVALID_CARD
		//}
		for _, comb := range n2 {
			cards := Int32ToCards(comb.Combination)
			if cards.Intersection(robotArgs.NoOutCards).Len() == 0 {
				return cards[0]
			}
		}
		outCard = Int32ToCards(n2[0].Combination)[0]
	case RobotLevel2:
		cards := CalculateResultOf3N2(args, robotArgs.RuleConfig)
		log.Sugar().Infof("=================== handCards:%v,recommend cards:%v ========================", robotArgs.HandCards, cards)
		for _, card := range cards {
			if !robotArgs.NoOutCards.In(card) {
				return card
			}
		}
	}
	// 分析向听数和最优组合
	//n2 := (args.HandTiles, args.ResidueTiles, make(map[string]CacheSC))
	return outCard
}
