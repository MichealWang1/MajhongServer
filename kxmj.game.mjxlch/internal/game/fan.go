package game

// 2023.7.28更新
// 单番型
type Fan struct {
	Id        uint32                 `json:"id"`       // 番型id
	Name      string                 `json:"name"`     // 番型名
	Multiple  string                 `json:"multiple"` // 番型倍数
	Mutex     []uint32               `json:"mutex"`    // 互斥番型
	Group     []uint32               `json:"group"`    // 组合id
	checkFunc func(*AnalysisHu) bool // 番型查找方法
}

var FanMap = map[uint32]*Fan{
	89:  {Id: 89, Name: "平胡", Multiple: "1", Mutex: []uint32{}, checkFunc: IsPingHu},
	47:  {Id: 47, Name: "碰碰和", Multiple: "2", Mutex: []uint32{89, 18, 119, 4}, checkFunc: IsPengPengHu},
	21:  {Id: 21, Name: "清一色", Multiple: "4", Mutex: []uint32{89}, checkFunc: IsQingYiSe},
	18:  {Id: 18, Name: "七对", Multiple: "4", Mutex: []uint32{89, 61, 55, 47, 54}, checkFunc: Is7Pair},
	119: {Id: 119, Name: "金钩钩", Multiple: "4", Mutex: []uint32{89, 47}, checkFunc: IsJinGouGou},
	120: {Id: 120, Name: "龙七对", Multiple: "8", Mutex: []uint32{89, 61, 55, 47, 18, 93, 54}, checkFunc: IsLongQiDui},
	121: {Id: 121, Name: "双龙七对", Multiple: "16", Mutex: []uint32{89, 61, 55, 47, 18, 93, 120, 54}, checkFunc: IsShuangLongQiDui},
	124: {Id: 124, Name: "三龙七对", Multiple: "32", Mutex: []uint32{89, 61, 55, 47, 18, 93, 120, 121}, checkFunc: IsSanLongQiDui},
	4:   {Id: 4, Name: "十八罗汉", Multiple: "64", Mutex: []uint32{89, 18, 93, 47, 61, 55, 119, 120, 121, 123, 125, 124}, checkFunc: IsShiBa},
	97:  {Id: 97, Name: "杠上开花", Multiple: "2", Mutex: []uint32{}, checkFunc: IsGangKai},
	98:  {Id: 98, Name: "杠上炮", Multiple: "2", Mutex: []uint32{}, checkFunc: IsGangPao},
	96:  {Id: 96, Name: "抢杠胡", Multiple: "2", Mutex: []uint32{}, checkFunc: IsQiangGang},
	79:  {Id: 79, Name: "自摸", Multiple: "2", Mutex: []uint32{}, checkFunc: IsZiMo},
	93:  {Id: 93, Name: "根", Multiple: "2", Mutex: []uint32{}, checkFunc: IsGen},
	43:  {Id: 43, Name: "海底捞月", Multiple: "2", Mutex: []uint32{}, checkFunc: IsHaiDi},
	103: {Id: 103, Name: "地胡", Multiple: "32", Mutex: []uint32{89, 79}, checkFunc: IsDiHu},
	102: {Id: 102, Name: "天胡", Multiple: "32", Mutex: []uint32{89, 79}, checkFunc: IsTianHu},
	126: {Id: 126, Name: "将对", Multiple: "8", Mutex: []uint32{89, 47, 67, 18}, checkFunc: IsJiangDui},
	123: {Id: 123, Name: "将七对", Multiple: "16", Mutex: []uint32{89, 18, 120, 121, 124, 61, 55, 67}, checkFunc: IsJiangLong},
	174: {Id: 174, Name: "将双七对", Multiple: "64", Mutex: []uint32{89, 93, 18, 120, 121, 123, 124, 61, 55, 67}, checkFunc: IsJiangShuangLong},
	175: {Id: 175, Name: "将三龙七对", Multiple: "126", Mutex: []uint32{89, 93, 174, 18, 120, 121, 123, 124, 61, 55, 67}, checkFunc: IsJiangSanLong},
	72:  {Id: 72, Name: "幺九", Multiple: "4", Mutex: []uint32{89}, checkFunc: IsYaoJiu},
	67:  {Id: 67, Name: "断幺九", Multiple: "2", Mutex: []uint32{89}, checkFunc: IsDuanYaoJiu},
}

var GroupFanMap = map[uint32]*Fan{
	132: {Id: 132, Name: "清龙七对", Multiple: "16", Mutex: []uint32{89, 61, 55, 47, 18, 93, 54, 21, 120, 133}, Group: []uint32{21, 120}},
	134: {Id: 134, Name: "清金钩钓", Multiple: "16", Mutex: []uint32{89, 47, 21, 119}, Group: []uint32{21, 119}},
	133: {Id: 133, Name: "清七对", Multiple: "16", Mutex: []uint32{89, 61, 55, 47, 54, 21, 18}, Group: []uint32{21, 18}},
	136: {Id: 136, Name: "清碰", Multiple: "8", Mutex: []uint32{89, 18, 119, 4, 21, 47}, Group: []uint32{21, 47}},
	128: {Id: 128, Name: "清十八罗汉", Multiple: "256", Mutex: []uint32{89, 18, 93, 47, 61, 55, 119, 120, 121, 123, 125, 124, 4, 21}, Group: []uint32{4, 21}},
	7:   {Id: 7, Name: "清幺九", Multiple: "16", Mutex: []uint32{89, 72, 21}, Group: []uint32{72, 21}},
	185: {Id: 185, Name: "清双龙七对", Multiple: "64", Mutex: []uint32{89, 61, 55, 47, 18, 93, 120, 54, 121, 21, 133}, Group: []uint32{121, 21}},
	186: {Id: 186, Name: "清三龙七对", Multiple: "128", Mutex: []uint32{89, 61, 55, 47, 18, 93, 120, 121, 124, 21, 133}, Group: []uint32{124, 21}},
}

// 胡牌番型封装
type HuResult struct {
	haveFan  map[uint32]struct{}
	fan      []uint32 // 番型
	multiple string   // 倍数
}

func NewHuResult() *HuResult {
	return &HuResult{
		haveFan:  make(map[uint32]struct{}, 0),
		fan:      make([]uint32, 0),
		multiple: "1",
	}
}

func (h *HuResult) Copy() *HuResult {
	haveFan := make(map[uint32]struct{}, 0)
	for k, _ := range h.haveFan {
		haveFan[k] = struct{}{}
	}
	fan := make([]uint32, len(h.fan))
	copy(fan, h.fan)
	return &HuResult{
		haveFan:  haveFan,
		fan:      fan,
		multiple: h.multiple,
	}
}

func (h *HuResult) setFan(x uint32) {
	if _, ok := h.haveFan[x]; !ok {
		h.haveFan[x] = struct{}{}
		h.fan = append(h.fan, x)
	}
}

func (h *HuResult) hasFan(x uint32) bool {
	if _, ok := h.haveFan[x]; !ok {
		return false
	}
	return true
}

func (h *HuResult) hasFans(fans []uint32) bool {
	for _, x := range fans {
		if _, ok := h.haveFan[x]; !ok {
			return false
		}
	}
	return true
}

func (h *HuResult) getFan() []uint32 {
	return h.fan
}

func (h *HuResult) removeFan(x uint32) {
	if _, ok := h.haveFan[x]; !ok {
		return
	}
	delete(h.haveFan, x)
	newFan := make([]uint32, len(h.fan)-1)
	for i, f := range h.fan {
		if f == x {
			newFan = append(h.fan[:i], h.fan[i+1:]...)
			break
		}
	}
	h.fan = newFan
}
