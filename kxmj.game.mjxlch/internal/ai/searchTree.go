package ai

// 搜索树

// 3N+2节点
type SearchNode3N2 struct {
	shanten  int32                   // 向听数
	children map[int32]SearchNode3N1 // 衍生的3N+1的节点（向听数不变的）
}

func (s *SearchNode3N2) GetChildren() map[int32]SearchNode3N1 {
	return s.children
}

func (s *SearchNode3N2) GetShanten() int32 {
	return s.shanten
}

// 3N+1节点
type SearchNode3N1 struct {
	handTiles    []int32                 // 手牌
	residueTiles []int32                 // 剩余牌
	shanten      int32                   // 向听数
	waits        map[int32]int32         // 能向前牌的剩余数
	children     map[int32]SearchNode3N2 // 能向前牌对应的3N+2节点
}

func (s *SearchNode3N1) GetWaits() map[int32]int32 {
	return s.waits
}

func (s *SearchNode3N1) GetChildren() map[int32]SearchNode3N2 {
	return s.children
}

func (s *SearchNode3N1) GetShenten() int32 {
	return s.shanten
}

func (s *SearchNode3N1) GetHandTiles() []int32 {
	return s.handTiles
}

// BuildSearchNode3N2 构建3N+2节点(shanten:当前向听数, toShanten:目标向听数, handTiles:手牌, residueTiles:剩余牌,tm:向听数的缓存)
func BuildSearchNode3N2(shanten, toShanten int32, handTiles, residueTiles []int32, ruleConfig *RuleConfig, tm map[string]CacheSC) SearchNode3N2 {
	children3N1 := make(map[int32]SearchNode3N1)
	// 如果向听数=-1已经胡牌则直接返回
	if shanten != -1 {
		// 删牌后向听数不变则删除这张牌
		for i, v := range handTiles {
			if v == 0 {
				continue
			}
			handTiles[i]--
			minShanten, _ := CalculateMinShanten(handTiles, ruleConfig, tm)
			if minShanten == shanten {
				children3N1[int32(i)] = BuildSearchNode3N1(shanten, toShanten, handTiles, residueTiles, ruleConfig, tm)
			}
			handTiles[i]++
		}
	}
	return SearchNode3N2{
		shanten:  shanten,
		children: children3N1,
	}
}

// BuildSearchNode3N1 构建3N+1节点(shanten:当前向听数, toShanten:目标向听数, handTiles:手牌, residueTiles:剩余牌,tm:向听数的缓存)
func BuildSearchNode3N1(shanten, toShanten int32, handTiles, residueTiles []int32, ruleConfig *RuleConfig, tm map[string]CacheSC) SearchNode3N1 {
	handTilesCopy := make([]int32, len(handTiles))
	copy(handTilesCopy, handTiles)
	res := SearchNode3N1{
		handTiles:    handTilesCopy,
		residueTiles: residueTiles,
		waits:        make(map[int32]int32),
		children:     make(map[int32]SearchNode3N2),
	}
	// 3N+1得从剩余牌里摸牌达到3N+2
	for i, v := range residueTiles {
		if v == 0 {
			continue
		}
		handTiles[i]++
		if minShanten, _ := CalculateMinShanten(handTiles, ruleConfig, tm); minShanten < shanten {
			// 向听前进了
			res.shanten = minShanten
			res.waits[int32(i)] = v
			if minShanten > toShanten {
				search3N2 := BuildSearchNode3N2(minShanten, toShanten, handTiles, residueTiles, ruleConfig, tm)
				res.children[int32(i)] = search3N2
			}
		}
		handTiles[i]--
	}
	return res
}

// AnalysisNode 分析节点
func (tree SearchNode3N2) AnalysisNode() Hand3N2AnalysisResultList {
	resList := make(Hand3N2AnalysisResultList, 0, len(tree.children))
	for card, t3n1 := range tree.children {
		node := t3n1.AnalysisNode()

		res := Hand3N2AnalysisResult{
			DiscardTile:      card,
			DiscardTileValue: float64(node.WaitsCount),
			Result13:         node,
		}

		resList = append(resList, res)
	}
	return resList
}

func (tree SearchNode3N1) AnalysisNode() Hand3N1AnalysisResult {
	handTiles := make([]int32, len(tree.handTiles))
	residueTiles := make([]int32, len(tree.residueTiles))
	copy(handTiles, tree.handTiles)
	copy(residueTiles, tree.residueTiles)
	waits := tree.waits
	waitCount := CountInt32Map(waits)

	res := Hand3N1AnalysisResult{
		HandTiles:    handTiles,
		ResidueTiles: residueTiles,
		Shanten:      tree.shanten,
		Waits:        waits,
		WaitsCount:   waitCount,
	}

	// 分析后续向前听

	return res
}
