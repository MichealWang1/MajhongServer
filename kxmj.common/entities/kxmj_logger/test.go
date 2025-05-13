package kxmj_logger

import "encoding/json"

type Test struct {
	Id uint32 `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"`
}

func (t *Test) TableName() string {
	return "test"
}

func (t *Test) Schema() string {
	return "kxmj_logger"
}

func (t *Test) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Test) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
