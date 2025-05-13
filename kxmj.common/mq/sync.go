package mq

import (
	"encoding/json"
	"fmt"
	"kxmj.common/log"
	"reflect"
	"strings"
)

type SyncType uint8

const (
	AddOrUpdate SyncType = 1 // 新增或者修改
	Delete      SyncType = 2 // 删除操作
)

type SyncEvent struct {
	Schema     string   `json:"schema"`     // 库名
	TableName  string   `json:"tableName"`  // 表名
	FilterKey  string   `json:"filterKey"`  // 过滤字段
	PrimaryKey string   `json:"primaryKey"` // 主键ID
	Type       SyncType `json:"type"`       // 同步方式
	Data       string   `json:"data"`       // 表数据
}

// SyncTable 同步更新表数据
func SyncTable(data interface{}, syncType SyncType) error {
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

	method = reflectValue.MethodByName("Schema")
	values = method.Call(nil)
	schema := values[0].String()

	// 反射获取过滤字段auto_increment
	filterKey := ""
	reflectType := reflect.TypeOf(data).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		if strings.Contains(reflectType.Field(i).Tag.Get("gorm"), "auto_increment") {
			filterKey = reflectType.Field(i).Tag.Get("json")
			break
		}
	}

	primaryKey := ""
	for i := 0; i < reflectType.NumField(); i++ {
		if strings.Contains(reflectType.Field(i).Tag.Get("gorm"), "primary_key") {
			primaryKey = reflectType.Field(i).Tag.Get("json")
			break
		}
	}

	err = Default().Publish(SyncDataQueue, &SyncEvent{
		Schema:     schema,
		TableName:  tableName,
		FilterKey:  filterKey,
		PrimaryKey: primaryKey,
		Type:       syncType,
		Data:       string(jsonData),
	})

	if err != nil {
		log.Sugar().Error(fmt.Sprintf("publish mq err:%v", err))
		return err
	}

	log.Sugar().Info(string(jsonData))
	return nil
}
