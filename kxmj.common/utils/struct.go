package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	//TypeOf会返回目标数据的类型，比如int/float/struct/指针等
	typ := reflect.TypeOf(data)
	//ValueOf返回目标数据的的值
	val := reflect.ValueOf(data)

	if val.Elem().Kind() != reflect.Struct {
		return result, errors.New(fmt.Sprintf("Invalid struct %v", typ))
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		field := typ.Elem().Field(i) //字段的数据类型
		tag := field.Tag
		key := tag.Get("redis")
		value := val.Elem().Field(i)

		switch value.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
			result[key] = value.Int()
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
			result[key] = value.Uint()
		case reflect.Float32, reflect.Float64:
			result[key] = value.Float()
		case reflect.String:
			result[key] = value.String()
		case reflect.Bool:
			result[key] = value.Bool()
		default:
			result[key] = value.Interface()
		}
	}

	return result, nil
}
