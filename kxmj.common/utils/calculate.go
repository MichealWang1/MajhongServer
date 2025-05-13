package utils

import (
	"math/big"
	"regexp"
)

// IsValidDigit 检查一个字符串是否是一个合法数值型
func IsValidDigit(str string) bool {
	if len(str) <= 0 {
		str = "0"
	}

	ok, _ := regexp.Match("-[0-9]+(\\\\.[0-9]+)?|[0-9]+(\\\\.[0-9]+)?", []byte(str))
	return ok
}

// Add 计算s1+s2;如果s1或s2不能转换成整形则返回false
func Add(s1, s2 string) (*big.Int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return nil, false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	return new(big.Int).Add(n1, n2), true
}

// Sub  计算s1-s2;如果s1或s2不能转换成整形则返回false
func Sub(s1, s2 string) (*big.Int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	return new(big.Int).Sub(n1, n2), true
}

// Mul  计算s1*s2;如果s1或s2不能转换成整形则返回false
func Mul(s1, s2 string) (*big.Int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	return new(big.Int).Mul(n1, n2), true
}

// Quo  计算s1/s2;如果s1或s2不能转换成整形或者s2==“0”则返回false
func Quo(s1, s2 string) (*big.Int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	if n2.Sign() == 0 {
		return new(big.Int).SetInt64(0), false
	}
	return new(big.Int).Quo(n1, n2), true
}

// Mod  计算s1%s2;如果s1或s2不能转换成整形或者s2==“0”则返回false
func Mod(s1, s2 string) (*big.Int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return new(big.Int).SetInt64(0), false
	}
	if n2.Sign() == 0 {
		return new(big.Int).SetInt64(0), false
	}
	return new(big.Int).Mod(n1, n2), true
}

// Cmp 比较两个string大小,false:string转big.Int出错
//
//	-1 if s1 < s2
//	 0 if s1 == s2
//	+1 if s1 > s2
func Cmp(s1, s2 string) int {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return 0
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return 0
	}
	return n1.Cmp(n2)
}

// Zero 获取0值big.Int实例
func Zero() *big.Int {
	return new(big.Int).SetInt64(0)
}

// AddToString 计算s1+s2;如果s1或s2不能转换成整形则返回false
func AddToString(s1, s2 string) (string, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return "0", false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return "0", false
	}
	return new(big.Int).Add(n1, n2).String(), true
}

// SubToString  计算s1-s2;如果s1或s2不能转换成整形则返回false
func SubToString(s1, s2 string) (string, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return "0", false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return "0", false
	}
	return new(big.Int).Sub(n1, n2).String(), true
}

// MulToString  计算s1*s2;如果s1或s2不能转换成整形则返回false
func MulToString(s1, s2 string) (string, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return "0", false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return "0", false
	}
	return new(big.Int).Mul(n1, n2).String(), true
}

// QuoToString  计算s1/s2;如果s1或s2不能转换成整形或者s2==“0”则返回false
func QuoToString(s1, s2 string) (string, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return "0", false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return "0", false
	}
	if n2.Sign() == 0 {
		return "0", false
	}
	return new(big.Int).Quo(n1, n2).String(), true
}

// ModToString  计算s1%s2;如果s1或s2不能转换成整形或者s2==“0”则返回false
func ModToString(s1, s2 string) (string, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return "0", false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return "0", false
	}
	if n2.Sign() == 0 {
		return "0", false
	}
	return new(big.Int).Mod(n1, n2).String(), true
}

// CmpToString 比较两个string大小,false:string转big.Int出错
//
//	-1 if s1 < s2
//	 0 if s1 == s2
//	+1 if s1 > s2
func CmpToString(s1, s2 string) (int, bool) {
	if len(s1) <= 0 {
		s1 = "0"
	}

	if len(s2) <= 0 {
		s2 = "0"
	}

	n1, b := new(big.Int).SetString(s1, 10)
	if !b {
		return 0, false
	}
	n2, b := new(big.Int).SetString(s2, 10)
	if !b {
		return 0, false
	}
	return n1.Cmp(n2), true
}
