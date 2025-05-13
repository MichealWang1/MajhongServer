package game

import (
	"kxmj.common/log"
	lib "kxmj.common/mahjong"
	"kxmj.common/utils"
	"sort"
	"strconv"
)

type HuData struct {
	seatId      uint32          // 检测玩家
	outSeatId   uint32          // 出牌玩家
	opCard      lib.Card        // 操作牌
	handCards   lib.Cards       // 玩家手牌
	actionCards lib.UserActions // 玩家动作牌
	result      *HuResult       // 胡牌番型倍数
}

func (h *HuData) copy() *HuData {
	return &HuData{
		seatId:      h.seatId,
		outSeatId:   h.outSeatId,
		opCard:      h.opCard,
		handCards:   h.handCards.Copy(),
		actionCards: h.actionCards.Copy(),
		result:      h.result.Copy(),
	}
}

// AnalysisHu 定义胡结构体
type AnalysisHu struct {
	SeatId           uint32          // 检测玩家
	OutSeatId        uint32          // 出牌玩家
	BankerSeatId     uint32          // 庄家位子
	OpCard           lib.Card        // 操作牌
	ChooseMissColor  lib.CardColor   // 牌花色
	HandCards        lib.Cards       // 玩家手牌(3n+1)
	ActionCards      lib.UserActions // 玩家动作牌
	OutCount         int32           // 玩家出牌次数
	IsQiangGang      bool            // 是否是抢杠
	PrevActionIsGang bool            // 上一个动作是不是杠(判断杠上开花,杠上炮)
	IsHaiDi          bool            // 是否是海底
}

// 检测胡和胡牌权位
func (c *AnalysisHu) analysisHuResult() (bool, *HuResult) {
	log.Sugar().Debugf("analysisHuResult AnalysisHu:%#v", c)
	// 有缺牌不能胡
	if c.HandCards.GetAllColorCards(c.ChooseMissColor).Len() != 0 {
		return false, nil
	}
	// 检测胡
	hu, huResult := analysisHuKinds(c)
	// 检测其他权位(胡牌才检测)
	if hu {
		// 分析其他番型
		analysisHuRights(c, huResult)
	}
	log.Sugar().Infof("hu:%v,huResult:%v", hu, huResult)
	return hu, huResult
}

// 检测胡
func analysisHuKinds(c *AnalysisHu) (bool, *HuResult) {
	res := NewHuResult()
	// 检测平胡
	if FanMap[Ping_Hu].checkFunc(c) {
		log.Sugar().Infof("huResult set Kind_Ping_Hu")
		res.setFan(Ping_Hu)
	}
	// 检测七对
	if FanMap[Qi_Dui].checkFunc(c) {
		log.Sugar().Infof("huResult set Kind_Qi_Dui")
		res.setFan(Qi_Dui)
	}

	return len(res.haveFan) != 0, res
}

// 计算倍数
func CalculateHuResult(c *AnalysisHu, result *HuResult) {
	multiple := "1"
	for _, fan := range FanMap {
		if result.hasFan(fan.Id) {
			multiple, _ = utils.MulToString(multiple, fan.Multiple)
			log.Sugar().Infof("result has fan: %#v multiple*%d = %d", fan, fan.Multiple, multiple)
			// 根额外计算
			if fan.Id == Gen {
				m1 := strconv.Itoa(1 << calculateGenNums(c.HandCards, c.OpCard, c.ActionCards))
				multiple, _ = utils.MulToString(multiple, m1)
			}
		}
	}
	for _, fan := range GroupFanMap {
		if result.hasFan(fan.Id) {
			multiple, _ = utils.MulToString(multiple, fan.Multiple)
			log.Sugar().Infof("result has fan: %#v multiple*%d = %d", fan, fan.Multiple, multiple)
		}
	}
	result.multiple = multiple
}

func analysisHuRights(c *AnalysisHu, huResult *HuResult) {
	for _, fan := range FanMap {
		if fan.checkFunc(c) {
			log.Sugar().Infof("huResult set huFan:%#v", fan)
			huResult.setFan(fan.Id)
		}
	}

	log.Sugar().Infof("huResult:%#v", huResult)
	// 组合番型
	CombinationFanGroup(huResult)
	log.Sugar().Infof("addCombination huResult:%#v", huResult)
	// 互斥番型
	MutexFan(huResult)
	log.Sugar().Infof("MutexFan huResult:%#v", huResult)
	// 计算倍数
	CalculateHuResult(c, huResult)
	log.Sugar().Infof("CalculateHuResult multiple:%v", huResult.multiple)
}

// 组合番型
func CombinationFanGroup(huResult *HuResult) {
	for _, fan := range GroupFanMap {
		if huResult.hasFans(fan.Group) {
			huResult.setFan(fan.Id)
		}
	}
}

// 互斥番型
func MutexFan(huResult *HuResult) {
	allFan := make([]*Fan, 0, len(FanMap)+len(GroupFanMap))
	for _, fan := range FanMap {
		allFan = append(allFan, fan)
	}
	for _, fan := range GroupFanMap {
		allFan = append(allFan, fan)
	}
	// 根据分数排序
	sort.Slice(allFan, func(i, j int) bool {
		return allFan[i].Multiple > allFan[j].Multiple
	})

	for _, f := range allFan {
		if huResult.hasFan(f.Id) {
			for _, fan := range f.Mutex {
				huResult.removeFan(fan)
			}
		}
	}
}

// 平胡
func IsPingHu(hu *AnalysisHu) bool {
	return lib.IsPingHu(hu.HandCards, hu.OpCard, false)
}

// 七对
func Is7Pair(hu *AnalysisHu) bool {
	return lib.Is7Pair(hu.HandCards, hu.OpCard)
}

// 自摸
func IsZiMo(hu *AnalysisHu) bool {
	return hu.SeatId == hu.OutSeatId
}

// 根
func IsGen(hu *AnalysisHu) bool {
	return calculateGenNums(hu.HandCards, hu.OpCard, hu.ActionCards) > 0
}

// 杠上开花
func IsGangKai(hu *AnalysisHu) bool {
	return hu.SeatId == hu.OutSeatId && hu.PrevActionIsGang
}

// 杠上炮
func IsGangPao(hu *AnalysisHu) bool {
	return hu.SeatId != hu.OutSeatId && hu.PrevActionIsGang
}

// 抢杠胡
func IsQiangGang(hu *AnalysisHu) bool {
	return hu.IsQiangGang
}

// 天胡
func IsTianHu(hu *AnalysisHu) bool {
	return hu.SeatId == hu.OutSeatId && hu.SeatId == hu.BankerSeatId && hu.OutCount == 0
}

// 地胡
func IsDiHu(hu *AnalysisHu) bool {
	return hu.SeatId == hu.OutSeatId && hu.SeatId != hu.BankerSeatId && hu.OutCount == 0
}

// 海底捞月
func IsHaiDi(hu *AnalysisHu) bool {
	return hu.IsHaiDi
}

// 碰碰胡
func IsPengPengHu(hu *AnalysisHu) bool {
	return lib.IsPengPengHu(hu.HandCards, hu.OpCard, hu.ActionCards)
}

// 断幺九
func IsDuanYaoJiu(hu *AnalysisHu) bool {
	for _, action := range hu.ActionCards {
		for _, card := range action.CombineCards {
			if card.GetColor() != lib.Card_Color_Zi && (card.GetValue() == 1 || card.GetValue() == 9) {
				return false
			}
		}
	}
	cards := hu.HandCards.AddCard(hu.OpCard)
	for _, card := range cards {
		if card.GetColor() != lib.Card_Color_Zi && (card.GetValue() == 1 || card.GetValue() == 9) {
			return false
		}
	}
	return true
}

// 清一色
func IsQingYiSe(hu *AnalysisHu) bool {
	return lib.IsQingYiSe(hu.HandCards, hu.OpCard, hu.ActionCards)
}

// 金钩钩
func IsJinGouGou(hu *AnalysisHu) bool {
	return hu.HandCards.Len() == 1
}

// 幺九
func IsYaoJiu(hu *AnalysisHu) bool {
	return lib.IsHunYaoJiu(hu.HandCards, hu.OpCard, hu.ActionCards)
}

// 龙七对
func IsLongQiDui(hu *AnalysisHu) bool {
	return Is7Pair(hu) && getHandFourCardCount(hu.HandCards, hu.OpCard) == 1
}

// 双龙七对
func IsShuangLongQiDui(hu *AnalysisHu) bool {
	return Is7Pair(hu) && getHandFourCardCount(hu.HandCards, hu.OpCard) == 2
}

// 三龙七对
func IsSanLongQiDui(hu *AnalysisHu) bool {
	return Is7Pair(hu) && getHandFourCardCount(hu.HandCards, hu.OpCard) == 3
}

// 将对
func IsJiangDui(hu *AnalysisHu) bool {
	return isJiangHu(hu.HandCards, hu.OpCard, hu.ActionCards) && lib.IsPengPengHu(hu.HandCards, hu.OpCard, hu.ActionCards)
}

// 将龙七对
func IsJiangLong(hu *AnalysisHu) bool {
	return isJiangHu(hu.HandCards, hu.OpCard, hu.ActionCards) && IsLongQiDui(hu)
}

// 将双七对
func IsJiangShuangLong(hu *AnalysisHu) bool {
	return isJiangHu(hu.HandCards, hu.OpCard, hu.ActionCards) && IsShuangLongQiDui(hu)
}

// 将三七对
func IsJiangSanLong(hu *AnalysisHu) bool {
	return isJiangHu(hu.HandCards, hu.OpCard, hu.ActionCards) && IsSanLongQiDui(hu)
}

// 十八罗汉
func IsShiBa(hu *AnalysisHu) bool {
	if hu.HandCards.Len() != 1 {
		return false
	}
	for _, action := range hu.ActionCards {
		if action.ActionType != lib.ActionType_Gang {
			return false
		}
	}
	return true
}

// 清碰
// 清七对
// 清龙七对
// 清双龙七对
// 清三龙七对
// 清金钩钩
// 清幺九
// 清十八罗汉

func isJiangHu(handCards lib.Cards, opCard lib.Card, actionCards lib.UserActions) bool {
	for _, action := range actionCards {
		for _, card := range action.CombineCards {
			if card.GetColor() == lib.Card_Color_Zi || (card.GetValue() != byte(2) && card.GetValue() != byte(5) && card.GetValue() != byte(8)) {
				return false
			}
		}
	}
	cards := handCards.AddCard(opCard)
	for _, card := range cards {
		if card.GetColor() == lib.Card_Color_Zi || (card.GetValue() != byte(2) && card.GetValue() != byte(5) && card.GetValue() != byte(8)) {
			return false
		}
	}
	return true
}

// 计算手中4张的数量(龙七对。。。)
func getHandFourCardCount(handCards lib.Cards, huCard lib.Card) int32 {
	cards := handCards.AddCard(huCard)
	res := int32(0)
	for _, count := range cards.ToMap() {
		if count == 4 {
			res++
		}
	}
	return res
}

// 计算根的数量
func calculateGenNums(handCards lib.Cards, opCard lib.Card, actionCards lib.UserActions) uint32 {
	cards := handCards.AddCard(opCard)
	for _, action := range actionCards {
		cards = cards.AddCards(action.CombineCards)
	}
	res := uint32(0)
	for _, c := range cards.ToMap() {
		if c >= 4 {
			res++
		}
	}
	return res
}
