package err_code

// UnKnowError 从9001开始
const (
	UnKnowError = 0 // 未知错误

)

var Message = map[int]string{
	UnKnowError: "未知错误",
}
