package kxmj_report

import "encoding/json"

type OrderTaskWin struct {
	OrderId   int64  `json:"order_id" redis:"order_id" gorm:"column:order_id;primary_key"` // 订单ID
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                // 用户ID
	Items     string `json:"items" redis:"items" gorm:"column:items"`                      // 领取物品:[{"id":102002,"count":"2180000000"},{"id":102002,"count":"2180000000"}]
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`       // 创建时间
}

func (o *OrderTaskWin) TableName() string {
	return "order_task_win"
}

func (o *OrderTaskWin) Schema() string {
	return "kxmj_report"
}

func (o *OrderTaskWin) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OrderTaskWin) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
