package business

import (
	"context"
	"fmt"
	"kxmj.common/redis_cache"
	"kxmj.core.api/config"
	"kxmj.core.api/internal/dto"
	"math/rand"
	"regexp"
)

func SendSms(ctx context.Context, parameter *dto.SendSmsParameter) (*dto.SendSmsResult, error) {
	var code string
	if config.Default.RunMode == "develop" {
		code = "123456"
	} else {
		code = fmt.Sprintf("%d", rand.Intn(900000)+100000)
	}

	// todo 第三方短信发送接口实现

	ttl, err := redis_cache.GetCache().GetSmsCache().Set(ctx, parameter.TelNumber, parameter.Type, code)
	if err != nil {
		return nil, err
	}

	return &dto.SendSmsResult{TTL: ttl}, nil
}

func CheckSms(ctx context.Context, parameter *dto.CheckSmsParameter) (bool, error) {
	code, err := redis_cache.GetCache().GetSmsCache().Get(ctx, parameter.TelNumber, parameter.Type)
	if err != nil {
		return false, err
	}
	return code == parameter.Code, nil
}

func PhoneRegex(number string, phone string) (bool, error) {
	var pattern string
	switch number {
	case "86":
		pattern = "^1[3-9]\\d{9}$"
		break
	case "91": // 印度
		pattern = "^[6789]\\d{9}$"
		break
	case "852": // 香港手机号
		pattern = "^[5679]\\d{7}$"
		break
	case "853": // 澳门
		pattern = "^[6]\\d{7}$"
		break
	case "855": // 柬埔寨
		pattern = "^(\\d{8})|(\\d{9})$"
		break
	case "856": // 老挝
		pattern = "^(20)\\d{8}$"
		break
	case "886": // 台湾省
		pattern = "^(9)\\d{8}$"
		break
	case "60": // 马来西亚
		pattern = "^(\\d{9})|(\\d{10})$"
		break
	case "62": // 印度尼西亚
		pattern = "^([89]\\d{10})|([89]\\d{11})$"
		break
	case "63": // 菲律宾
		pattern = "^\\d{10}$"
		break
	case "65": // 新加坡
		pattern = "^[89]\\d{7}$"
		break
	case "66": // 泰国
		pattern = "^\\d{9}$"
		break
	case "84": // 越南
		pattern = "^\\d{9}$"
		break
	case "95": // 缅甸
		pattern = "^(9)\\d{9}$"
		break
	case "55": // 巴西
		pattern = "^\\d{11}$"
		break
	default:
		pattern = ""
		break
	}

	return regexp.MatchString(pattern, phone)
}
