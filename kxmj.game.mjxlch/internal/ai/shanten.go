package ai

import "fmt"

const (
	MaxShanten14 = 8 //14张牌，最大向听
)

// Shanten 向听结构体
type Shanten struct {
	handLen    int32   // 手牌数量
	handCards  []int32 // 各牌数量(0-8:1-9万；9-17：1-9条；18-26：1-9筒；27-30：东南西北；31-33：中发白)
	shunKeNums int32   // 刻子和顺子数量(包括动作牌)
	daZiNums   int32   // 搭子数量（组成顺子或刻子差一张的组合）
	pairNums   int32   // 对子数量（将牌）

	combinations [][]int32 // 组合牌
	combScore    int32     // 组合分

	mustOutCardNums int32 // 必须打的牌数量
	signalCard      int64 // 单牌牌（34张牌用int64的34位）
	anGangCard      int64 // 暗杠牌（34张牌用int64的34位）

	optimalCombination [][]int32 // 最优组合牌
	OpCombScore        int32     // 最优组合分
	minShanten         int32     // 最小向听数
}

func (s *Shanten) String() string {
	return fmt.Sprintf("handLen=%d,handCards=%v,shunKeNums=%d,daZiNums=%d,pairNums=%d,combinations=%v,mustOutCardNums=%d,signalCard=%d,anGangCard=%d,optimalCombination=%v,minShanten=%d",
		s.handLen, s.handCards, s.shunKeNums, s.daZiNums, s.pairNums, s.combinations, s.mustOutCardNums, s.signalCard, s.anGangCard, s.optimalCombination, s.minShanten)
}

func NewShanten(handCards []int32) *Shanten {
	s := &Shanten{
		handCards:          handCards,
		combinations:       make([][]int32, 0, 5),
		minShanten:         MaxShanten14,
		optimalCombination: make([][]int32, 0, 5),
	}
	handLen := int32(0)
	for card, count := range handCards {
		handLen += count
		if count == 4 {
			s.anGangCard |= int64(1) << uint(card)
		}
	}
	s.handLen = handLen
	s.shunKeNums = (14 - handLen) / 3
	return s
}

func (s *Shanten) AddCombinations(cards []int32, score int32) {
	s.combinations = append(s.combinations, cards)
	s.combScore += score
}

func (s *Shanten) RemoveCombinations(score int32) {
	s.combinations = s.combinations[:len(s.combinations)-1]
	s.combScore -= score
}

// CalculateShanten 计算向听数（S = 最大向听数-2*顺刻数-搭子数-对子数 ）
func (s *Shanten) CalculateShanten() int32 {
	shanten := MaxShanten14 - 2*s.shunKeNums - s.daZiNums - s.pairNums
	// 面子数(除去一组对子以外的对子+搭子数+顺刻数)
	faceNums := s.daZiNums + s.shunKeNums + s.pairNums
	if s.pairNums > 0 {
		faceNums -= 1
	} else if s.signalCard|s.anGangCard == s.anGangCard { // 没有将牌且有单牌都是暗杠牌 如：5555 这个的向听数是1
		shanten++
	}
	// 面子数过多
	if faceNums >= 5 {
		shanten += faceNums - 4
	}
	mustOutCardNums := s.mustOutCardNums
	if shanten != -1 {
		if s.handLen%3 == 2 {
			mustOutCardNums--
		}
		if shanten < mustOutCardNums {
			return mustOutCardNums
		}
	}
	return shanten
}

// DecomposeHandTile 递归分解手牌(从0-33);depth 表示当前拆到第几张牌
func (s *Shanten) DecomposeHandTile(depth int) {
	// 当最小向听数 = -1 说明已经胡牌
	if s.minShanten == -1 {
		return
	}

	// 获取牌数量>0的下一张牌
	for ; depth < 34 && s.handCards[depth] == 0; depth++ {
	}

	// 如果已经拆完34张牌则返回计算向听数
	if depth >= 34 {
		// TODO:计算向听数
		shanten := s.CalculateShanten()
		if s.minShanten > shanten || (s.minShanten == shanten && s.combScore > s.OpCombScore) {
			s.minShanten = shanten
			combs := make([][]int32, 0, 5)
			combs = append(combs, s.combinations...)
			s.optimalCombination = combs
			s.OpCombScore = s.combScore
		}
		//fmt.Println(s.String(), shanten)
		return
	}

	if TileIsWTT(depth) { // 拆万筒条
		s.DecomposeWTT(depth)
	} else if TileIsDNXB(depth) { // 拆东南西北
		s.DecomposeDNXB(depth)
	} else { // 拆中发白
		s.DecomposeZFB(depth)
	}

}

// DecomposeWTT 拆万条筒
func (s *Shanten) DecomposeWTT(depth int) {
	// 判断是否是万条筒
	if !TileIsWTT(depth) {
		return
	}
	// 获取当前位置的牌值
	value := GetTilesValue(int32(depth))
	// 判断数量
	switch s.handCards[depth] {
	case 1: // 单张牌： 组成顺子它的值必须小于等于7
		// 与后面的两张牌组成顺子123
		// 与后面一张组成搭子12-
		// 与后面第二张组成搭子1-3
		// 单牌
		if value < 7 && s.handCards[depth+1] == 1 && s.handCards[depth+2] > 0 && s.handCards[depth+3] < 4 {
			// 组成顺子
			s.increaseStraight(depth)
			s.DecomposeHandTile(depth + 2)
			s.decreaseStraight(depth)

		} else {
			if value <= 7 && s.handCards[depth+2] > 0 {
				// 与后面的两张牌组成顺子123
				if s.handCards[depth+1] != 0 {
					s.increaseStraight(depth)
					s.DecomposeHandTile(depth + 1)
					s.decreaseStraight(depth)
				}
				// 与后面第二张组成搭子1-3
				s.increaseLUG13(depth)
				s.DecomposeHandTile(depth + 1)
				s.decreaseLUG13(depth)
			}
			// 与后面一张组成搭子12-
			if value <= 8 && s.handCards[depth+1] > 0 {
				s.increaseLUG12(depth)
				s.DecomposeHandTile(depth + 1)
				s.decreaseLUG12(depth)
			}
			// 单牌处理
			s.increaseSignal(depth)
			s.DecomposeHandTile(depth + 1)
			s.decreaseSignal(depth)
		}
	case 2: // 两张牌
		// 凑成将
		s.increasePair(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreasePair(depth)
		// 顺子
		if depth <= 7 && s.handCards[depth+1] > 0 && s.handCards[depth+2] > 0 {
			s.increaseStraight(depth)
			s.DecomposeHandTile(depth)
			s.decreaseStraight(depth)
		}

	case 3: // 三张牌
		// 三暗刻
		s.increaseGroup(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreaseGroup(depth)

		// 将+顺子、将+搭子
		s.increasePair(depth)
		if value <= 7 && s.handCards[depth+2] > 0 {
			// 与后面的两张牌组成顺子123
			if s.handCards[depth+1] != 0 {
				s.increaseStraight(depth)
				s.DecomposeHandTile(depth + 1)
				s.decreaseStraight(depth)
			}
			// 与后面第二张组成搭子1-3
			s.increaseLUG13(depth)
			s.DecomposeHandTile(depth + 1)
			s.decreaseLUG13(depth)
		}
		// 与后面一张组成搭子12-
		if value <= 8 && s.handCards[depth+1] > 0 {
			s.increaseLUG12(depth)
			s.DecomposeHandTile(depth + 1)
			s.decreaseLUG12(depth)
		}
		s.decreasePair(depth)

		// 两组顺子
		if value <= 7 && s.handCards[depth+1] >= 2 && s.handCards[depth+2] >= 2 {
			s.increaseStraight(depth)
			s.increaseStraight(depth)
			s.DecomposeHandTile(depth)
			s.decreaseStraight(depth)
			s.decreaseStraight(depth)
		}

	case 4:
		// 拆暗刻
		s.increaseGroup(depth)
		s.DecomposeHandTile(depth)
		s.decreaseGroup(depth)
		// 拆对子
		s.increasePair(depth)
		if value <= 7 && s.handCards[depth+2] > 0 {
			// 与后面的两张牌组成顺子123
			if s.handCards[depth+1] != 0 {
				s.increaseStraight(depth)
				s.DecomposeHandTile(depth)
				s.decreaseStraight(depth)
			}
			// 与后面第二张组成搭子1-3
			s.increaseLUG13(depth)
			s.DecomposeHandTile(depth)
			s.decreaseLUG13(depth)
		}
		// 与后面一张组成搭子12-
		if value <= 8 && s.handCards[depth+1] > 0 {
			s.increaseLUG12(depth)
			s.DecomposeHandTile(depth)
			s.decreaseLUG12(depth)
		}
		s.decreasePair(depth)
	}
}

// 拆东南西北
func (s *Shanten) DecomposeDNXB(depth int) {
	if !TileIsDNXB(depth) {
		return
	}
	// 判断数量
	switch s.handCards[depth] {
	case 1: // 单张
		s.increaseSignal(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreaseSignal(depth)
	case 2: // 对子
		s.increasePair(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreasePair(depth)
	case 3: //刻子
		s.increaseGroup(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreaseGroup(depth)
	case 4: //刻子+单张（这个单张为必须出的）
		s.increaseGroup(depth)
		s.increaseSignal(depth)
		s.mustOutCardNums++
		s.DecomposeHandTile(depth + 1)
		s.decreaseSignal(depth)
		s.mustOutCardNums--
		s.decreaseGroup(depth)
	}
}

// 拆中发白
func (s *Shanten) DecomposeZFB(depth int) {
	if !TileIsZFB(depth) {
		return
	}
	// 判断数量
	switch s.handCards[depth] {
	case 1: // 单张
		s.increaseSignal(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreaseSignal(depth)
	case 2: // 对子
		s.increasePair(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreasePair(depth)
	case 3: //刻子
		s.increaseGroup(depth)
		s.DecomposeHandTile(depth + 1)
		s.decreaseGroup(depth)
	case 4: //刻子+单张（这个单张为必须出的）
		s.increaseGroup(depth)
		s.increaseSignal(depth)
		s.mustOutCardNums++
		s.DecomposeHandTile(depth + 1)
		s.decreaseSignal(depth)
		s.mustOutCardNums--
		s.decreaseGroup(depth)
	}
}

// 计算七对
func (s *Shanten) CalcQiDuiShanten() {

	// 手牌小于13则返回
	if s.handLen < 13 {
		return
	}

}
