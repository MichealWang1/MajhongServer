package mq

import (
	"encoding/json"
	"fmt"
	"kxmj.common/entities/kxmj_core"
	"kxmj.common/log"
	"reflect"
	"strings"
)

// UpdateStatistics 更新用户游戏统计数据
func UpdateStatistics(data *kxmj_core.UserGameStatistics) error {
	defer func() {
		if err := recover(); err != nil {
			log.Sugar().Error(fmt.Sprintf("panic err:%v", err))
		}
	}()

	err := Default().Publish(GameLoggerQueue, data)

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("publish mq err:%v", err))
		return err
	}

	log.Sugar().Info(data)
	return nil
}

// AddGameLogger 新增数据库记录，data必须是数据库对应的实体指针，该方法会把记录发布到消息队列中，再由logger服务订阅入库
func AddGameLogger(data interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			log.Sugar().Error(fmt.Sprintf("panic err:%v", err))
		}
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Marshal err:%v", err))
		return err
	}

	// 反射获取table name
	reflectValue := reflect.ValueOf(data)
	method := reflectValue.MethodByName("TableName")
	values := method.Call(nil)
	tableName := values[0].String()

	// 反射获取过滤字段auto_increment
	filterKey := ""
	reflectType := reflect.TypeOf(data).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		if strings.Contains(reflectType.Field(i).Tag.Get("gorm"), "auto_increment") {
			filterKey = reflectType.Field(i).Tag.Get("json")
			break
		}
	}

	err = Default().Publish(GameLoggerQueue, &Record{
		TableName: tableName,
		Data:      string(jsonData),
		FilterKey: filterKey,
	})

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("publish mq err:%v", err))
		return err
	}

	log.Sugar().Info(string(jsonData))
	return nil
}
