package codes

const (
	ReadyReqInfoError        = 50001 + iota // 准备请求信息错误
	PlayerIsReady                           // 玩家已经准备
	NowStateNotReady                        // 当前状态不能准备
	SwapReqInfoError                        // 换牌信息错误
	PlayerIsSwap                            // 玩家已经换牌
	NowStateNotSwap                         // 当前状态不能换牌
	ChooseMissReqInfoError                  // 选缺信息错误
	PlayerIsChooseMiss                      // 玩家已进行选缺
	NowStateNotChooseMiss                   // 当前状态不能选缺
	PlayerHandNotHaveTheCard                // 玩家手中没有这张牌
	PlayerHandHaveSwapCard                  // 玩家手里还有缺牌
	PlayerNotHaveOutCardAuth                // 该位子没有出牌权限
	OperateReqInfoError                     // 请求操作信息错误
	PlayerNotHaveOperateAuth                // 玩家没有这个操作权限
	NowStateNotOperate                      // 该状态不能操作
	PlayerNotExist                          // 玩家不存在
)

var gameMessage = map[int]string{
	ReadyReqInfoError:        "准备请求信息错误",
	PlayerIsReady:            "玩家已经准备",
	NowStateNotReady:         "当前状态不能准备",
	SwapReqInfoError:         "换牌信息错误",
	PlayerIsSwap:             "玩家已经换牌",
	NowStateNotSwap:          "当前状态不能换牌",
	ChooseMissReqInfoError:   "选缺信息错误",
	PlayerIsChooseMiss:       "玩家已进行选缺",
	NowStateNotChooseMiss:    "当前状态不能选缺",
	PlayerHandNotHaveTheCard: "玩家手中没有这张牌",
	PlayerHandHaveSwapCard:   "玩家手里还有缺牌",
	PlayerNotHaveOutCardAuth: "该位子没有出牌权限",
	OperateReqInfoError:      "请求操作信息错误",
	PlayerNotHaveOperateAuth: "玩家没有这个操作权限",
	NowStateNotOperate:       "该状态不能操作",
	PlayerNotExist:           "玩家不存在",
}
