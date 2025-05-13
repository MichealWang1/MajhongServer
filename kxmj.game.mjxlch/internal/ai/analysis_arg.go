package ai

import lib "kxmj.common/mahjong"

// MustNeedArgs 机器人分析需要的参数(开始时的参数，不改变)
type MustNeedArgs struct {
	SeatId          int32                                          // 分析的seatId
	OutSeatId       int32                                          // huCard是哪家的
	UserCount       int32                                          // 玩家数量
	OutCardCount    int32                                          // 玩家出牌次数
	HandCards       lib.Cards                                      // SeatId用户的手牌
	DisCards        lib.Cards                                      // 所有用户的弃牌
	OpCard          lib.Card                                       // seatId的操作牌
	KingCards       lib.Cards                                      // 精牌，无精的传的是空的slice
	NoOutCards      lib.Cards                                      // 不能出的牌
	MustCanOutCards lib.Cards                                      // 必须出的牌
	CanHuFunc       func(args interface{}) ([]string, int32, bool) // 检测胡牌函数，返回胡牌得分
	//ActionTable  *lib.ActionTable                               // 当前operate状态
	//ActionItems  *lib.ActionCardsTable                          // 所有用户的actions保存
	//GoldCards    lib.Cards                                      // 金牌
}

// RuleConfig 游戏规则
type RuleConfig struct {
	// 是否可胡七对
	CanHuQiDui bool
}

type RobotLevel int32

const (
	RobotLevel1 RobotLevel = 1
	RobotLevel2 RobotLevel = 2
	RobotLevel3 RobotLevel = 3
	RobotLevel4 RobotLevel = 4
	RobotLevel5 RobotLevel = 5
	RobotLevel6 RobotLevel = 6
)

type RobotConfig struct {
	Level      RobotLevel // 机器人等级
	TargetType int32      // 目标牌型
	TargetFunc func(args interface{}) (lib.HuResult, lib.Cards, error)
	*ScoreConfig
}

// 分数配置
type ScoreConfig struct {
	TargetScore int32 // 目标牌型分
	// 花色
	OneColorScore   int32 // 主花分
	TwoColorScore   int32 // 二花分
	ThreeColorScore int32 // 三花分
	FengColorScore  int32 // 风牌分
	OneKingScore    int32 // 单张癞子分
	AllKingScore    int32 // 多癞子分
	// 牌型
	SingleScore     int32 // 不相连单牌分
	ConnSingleScore int32 // 相连单牌分
	PairScore       int32 // 不相连对子分
	ConnPairScore   int32 // 相连对子分
	StraightScore   int32 // 顺子分
	KeZiScore       int32 // 不相连刻子分
	ConnKeZiScore   int32 // 相连刻子分
}

type EstimateArgs struct {
	*MustNeedArgs // 继承
	*RuleConfig   // 游戏规则
	*RobotConfig  // 机器人规则
}

// RobotEstimateHuArgs 机器人分析参数(分析过程中的参数，会改变)
type RobotEstimateHuArgs struct {
	HandCardsLen      int32     // 手牌数量
	HandTiles         []int32   // 手牌，不含副露
	DiscardsTiles     []int32   // 所有用户的弃牌
	ResidueTiles      []int32   // 剩余牌
	KingIndex         []int32   // 精牌的位置数组，无精空
	TreeOpCard        lib.Card  // 搜索树胡牌时的操作牌
	TreeHandCards     lib.Cards // 搜索树胡牌时的手牌
	NotAllowedOutCard lib.Cards // 不允许出的牌
	MustCanOutCards   lib.Cards //必须出的牌
	//OutKingCount      int32     // 玩家飞宝数目
	//GoldIndex         []int32   // 金牌的位置数组，无空
}

// NewEstimateArgs 初始化需要的数据
func NewEstimateArgs() *EstimateArgs {
	mustNeedArgs := &MustNeedArgs{
		SeatId:       0,
		OutSeatId:    0,
		UserCount:    0,
		OutCardCount: 0,
		HandCards:    make(lib.Cards, 0, 14),
		DisCards:     make(lib.Cards, 0, 138),
		OpCard:       lib.INVALID_CARD,
		KingCards:    make(lib.Cards, 0, 4),
		NoOutCards:   make(lib.Cards, 0, 4),
		CanHuFunc:    nil,
	}

	ruleConfig := &RuleConfig{
		CanHuQiDui: false,
	}

	robotConfig := &RobotConfig{
		Level:       RobotLevel2,
		TargetType:  1,
		TargetFunc:  nil,
		ScoreConfig: &ScoreConfig{},
	}

	return &EstimateArgs{
		MustNeedArgs: mustNeedArgs,
		RuleConfig:   ruleConfig,
		RobotConfig:  robotConfig,
	}
}

// NewRobotEstimateHuArgs 初始化分析参数
func NewRobotEstimateHuArgs() *RobotEstimateHuArgs {
	return &RobotEstimateHuArgs{
		HandTiles:         make([]int32, 0, 34),   // 手牌，不含副露
		DiscardsTiles:     make([]int32, 0, 34),   // 所有用户的弃牌
		ResidueTiles:      make([]int32, 0, 34),   // 剩余牌
		KingIndex:         make([]int32, 0),       // 精牌的位置数组，无精空
		TreeHandCards:     make(lib.Cards, 0, 14), // 搜索树胡牌时的手牌
		NotAllowedOutCard: make(lib.Cards, 0, 14), // 不允许出的牌
	}
}

// 通过游戏参数构建分析参数
func (this *RobotEstimateHuArgs) Build(other *EstimateArgs) {
	this.HandCardsLen = int32(other.HandCards.Len())
	this.HandTiles = CardToTiles(other.HandCards)
	this.DiscardsTiles = CardToTiles(other.DisCards)
	this.ResidueTiles = GetReducedTiles(this.HandTiles, this.DiscardsTiles)
	this.KingIndex = CardsToIndex(other.KingCards)
	this.NotAllowedOutCard = other.NoOutCards
	this.MustCanOutCards = other.MustCanOutCards
}
