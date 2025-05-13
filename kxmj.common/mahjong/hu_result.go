package lib

type HuResult struct {
	Kind  int64 // 是否能胡牌位(平胡、七对、十三幺等，能确定胡的权位)
	Right int64 // 胡牌的其他权位(清一色、一条龙、碰碰胡等，在胡基础上能加分的权位)
}

func NewHuResult() *HuResult {
	return &HuResult{
		Kind:  0,
		Right: 0,
	}
}

func (h *HuResult) Copy() *HuResult {
	return &HuResult{
		Kind:  h.Kind,
		Right: h.Right,
	}
}

// GetKind 获取能胡牌权位
func (h *HuResult) GetKind() int64 { return h.Kind }

// GetRight 获取额外权位
func (h *HuResult) GetRight() int64 { return h.Right }

// SetKind 设置胡牌权位
func (h *HuResult) SetKind(k int64) { h.Kind |= k }

// SetRight 设置额外权位
func (h *HuResult) SetRight(r int64) { h.Right |= r }

// RemoveKind 移除一种胡牌权位
func (h *HuResult) RemoveKind(k int64) { h.Kind &^= k }

// RemoveRight 移除一种额外权位
func (h *HuResult) RemoveRight(r int64) { h.Right &^= r }

// HasKind 是否有这种权位
func (h *HuResult) HasKind(k int64) bool { return h.Kind&k != 0 }

// HasRight 是否有这种权位
func (h *HuResult) HasRight(r int64) bool { return h.Right&r != 0 }

// ResetKind 重置胡牌权位
func (h *HuResult) ResetKind() { h.Kind = 0 }

// ResetRight 重置额外权位
func (h *HuResult) ResetRight() { h.Right = 0 }

// Remove 将权位置为空
func (h *HuResult) Remove() { h = nil }
