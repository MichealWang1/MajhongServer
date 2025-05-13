package ai

const (
	SignalScore   int32 = (iota + 1) * 5 // 单牌分
	DaZi13Score                          // 1-3搭子分
	DaZi12Score                          // 12-搭子分
	PairScore                            // 对子分
	StraightScore                        // 顺子分
	GroupScore                           // 刻子分

)

// 增加单牌
func (s *Shanten) increaseSignal(i int) {
	s.handCards[i]--
	s.signalCard |= int64(1) << uint(i)
	s.AddCombinations([]int32{int32(i)}, SignalScore)
}

// 删除单牌
func (s *Shanten) decreaseSignal(i int) {
	s.handCards[i]++
	s.signalCard &^= int64(1) << uint(i)
	s.RemoveCombinations(SignalScore)
}

// 增加对子
func (s *Shanten) increasePair(i int) {
	s.handCards[i] -= 2
	s.pairNums++
	s.AddCombinations([]int32{int32(i), int32(i)}, PairScore)
}

// 删除对子
func (s *Shanten) decreasePair(i int) {
	s.handCards[i] += 2
	s.pairNums--
	s.RemoveCombinations(PairScore)
}

// 增加一组刻子
func (s *Shanten) increaseGroup(i int) {
	s.handCards[i] -= 3
	s.shunKeNums++
	s.AddCombinations([]int32{int32(i), int32(i), int32(i)}, GroupScore)
}

// 删除一组刻子
func (s *Shanten) decreaseGroup(i int) {
	s.handCards[i] += 3
	s.shunKeNums--
	s.RemoveCombinations(GroupScore)
}

// 增加一组顺子
func (s *Shanten) increaseStraight(i int) {
	s.handCards[i]--
	s.handCards[i+1]--
	s.handCards[i+2]--
	s.shunKeNums++
	s.AddCombinations([]int32{int32(i), int32(i + 1), int32(i + 2)}, StraightScore)
}

// 删除一组顺子
func (s *Shanten) decreaseStraight(i int) {
	s.handCards[i]++
	s.handCards[i+1]++
	s.handCards[i+2]++
	s.shunKeNums--
	s.RemoveCombinations(StraightScore)
}

// 增加一组搭子12-
func (s *Shanten) increaseLUG12(i int) {
	s.handCards[i]--
	s.handCards[i+1]--
	s.daZiNums++
	s.AddCombinations([]int32{int32(i), int32(i + 1)}, DaZi12Score)
}

// 删除一组搭子12-
func (s *Shanten) decreaseLUG12(i int) {
	s.handCards[i]++
	s.handCards[i+1]++
	s.daZiNums--
	s.RemoveCombinations(DaZi12Score)
}

// 增加一组搭子1-3
func (s *Shanten) increaseLUG13(i int) {
	s.handCards[i]--
	s.handCards[i+2]--
	s.daZiNums++
	s.AddCombinations([]int32{int32(i), int32(i + 2)}, DaZi13Score)
}

// 删除一组搭子1-3
func (s *Shanten) decreaseLUG13(i int) {
	s.handCards[i]++
	s.handCards[i+2]++
	s.daZiNums--
	s.RemoveCombinations(DaZi13Score)
}
