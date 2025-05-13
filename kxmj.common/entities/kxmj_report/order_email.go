package kxmj_report

import "encoding/json"

type OrderEmail struct {
	OrderId   int64  `json:"order_id" redis:"order_id" gorm:"column:order_id;primary_key"` // 订单ID
	EmailId   uint32 `json:"email_id" redis:"email_id" gorm:"column:email_id"`             // 邮件ID
	EmailType uint8  `json:"email_type" redis:"email_type" gorm:"column:email_type"`       // 邮件类型：1 福利发放；2 系统通知
	UserId    uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                // 用户ID
	Status    uint8  `json:"status" redis:"status" gorm:"column:status"`                   // 邮件状态：1 未读；2 已读；3 已领取；
	DrawTime  uint32 `json:"draw_time" redis:"draw_time" gorm:"column:draw_time"`          // 领取时间
	CreatedAt uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`       // 创建时间
	UpdatedAt uint32 `json:"updated_at" redis:"updated_at" gorm:"column:updated_at"`       // 更新时间
}

func (o *OrderEmail) TableName() string {
	return "order_email"
}

func (o *OrderEmail) Schema() string {
	return "kxmj_report"
}

func (o *OrderEmail) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OrderEmail) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
