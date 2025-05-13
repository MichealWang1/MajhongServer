package err_code

const (
	InvalidOpen = 10001 // 无效操作
	GetItemNUll = 10002 // 领取物品为空
)

var Message = map[int]string{
	InvalidOpen: "无效操作",
	GetItemNUll: "领取物品为空",
}
