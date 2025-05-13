package utils

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(IsValidDigit("1111"))
	fmt.Println(IsValidDigit("0000"))
	fmt.Println(IsValidDigit("aaaa"))
	fmt.Println(IsValidDigit(""))
}
