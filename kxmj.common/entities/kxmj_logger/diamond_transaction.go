package kxmj_logger

import "encoding/json"

type DiamondTransaction struct {
	Id           int64  `json:"id" redis:"id" gorm:"column:id;primary_key"`                      // 日志ID
	OrderId      int64  `json:"order_id" redis:"order_id" gorm:"column:order_id"`                // 订单号
	UserId       uint32 `json:"user_id" redis:"user_id" gorm:"column:user_id"`                   // 用户ID
	Type         uint8  `json:"type" redis:"type" gorm:"column:type"`                            // 交易类型：1 增加；2 扣减；
	BusinessType uint8  `json:"business_type" redis:"business_type" gorm:"column:business_type"` // 业务类型：1 商城；2 任务；3 邮件；
	Count        string `json:"count" redis:"count" gorm:"column:count"`                         // 数量
	CreatedAt    uint32 `json:"created_at" redis:"created_at" gorm:"column:created_at"`          // 创建时间
}

func (d *DiamondTransaction) TableName() string {
	return "diamond_transaction"
}

func (d *DiamondTransaction) Schema() string {
	return "kxmj_logger"
}

func (d *DiamondTransaction) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DiamondTransaction) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
