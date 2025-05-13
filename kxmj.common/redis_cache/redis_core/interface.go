package redis_core

import "context"

type EventParams struct {
	Action  string                 // 更新事件类型
	Data    interface{}            // 更新表数据
	Updates map[string]interface{} // 当Action等于Update类型时记录已修改字段
}

type ICache interface {
	EventHandler(ctx context.Context, e *EventParams) // 事件回调
	GetTableTemplate() interface{}                    // 数据库表模板实例
}
