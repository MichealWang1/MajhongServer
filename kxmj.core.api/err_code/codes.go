package err_code

const (
	PasswordError = 2001 // 密码错误

)

var Message = map[int]string{
	PasswordError: "密码错误",
}
