package item

// ExtraType 扩展属性
type ExtraType uint8

const (
	Attack ExtraType = 1 // 攻击属性
	Rise   ExtraType = 2 // 复活属性
)

type Extra struct {
	Type  ExtraType `json:"type"`  // 属性类型
	Value uint32    `json:"value"` // 属性值
}
