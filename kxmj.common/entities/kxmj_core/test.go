package kxmj_core

import "encoding/json"

type Test struct {
	Id      uint32 `json:"id" redis:"id" gorm:"column:id;primary_key;auto_increment"`
	Column1 string `json:"Column1" redis:"Column1" gorm:"column:Column1"`
}

func (t *Test) TableName() string {
	return "test"
}

func (t *Test) Schema() string {
	return "kxmj_core"
}

func (t *Test) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Test) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
