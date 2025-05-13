package lib

import (
	"errors"
	"fmt"
	"sort"
)

type Cards []Card

func (c Cards) Len() int {
	return len(c)
}

func (c Cards) String() string {
	return fmt.Sprintf("%#v", c)
}
func (c Cards) ToInt32() []int32 {
	res := make([]int32, len(c))
	for i, card := range c {
		res[i] = int32(card)
	}
	return res
}

func (c Cards) ToUint32() []uint32 {
	res := make([]uint32, len(c))
	for i, card := range c {
		res[i] = uint32(card)
	}
	return res
}

func Int32ToCards(ints []int32) Cards {
	res := make(Cards, 0, len(ints))
	for _, v := range ints {
		res = append(res, Card(v))
	}
	return res
}

func Uint32ToCards(ints []uint32) Cards {
	res := make(Cards, 0, len(ints))
	for _, v := range ints {
		res = append(res, Card(v))
	}
	return res
}

// Sort 根据牌大小排序
func (c *Cards) Sort() {
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i] < (*c)[j]
	})
}

// SortByValue 根据牌值排序
func (c *Cards) SortByValue() {
	sort.Slice(*c, func(i, j int) bool {
		return (*c)[i].GetValue() < (*c)[j].GetValue()
	})
}

// In 是否包含这张牌
func (c Cards) In(card Card) bool {
	for _, card2 := range c {
		if card == card2 {
			return true
		}
	}
	return false
}

// IsContain 是否全包含
func (c Cards) IsContain(cards Cards) bool {
	oriCards := c.Copy()
	tarCards := cards.Copy()
	oriCards.Sort()
	tarCards.Sort()
	i, j := 0, 0
	for ; i < oriCards.Len() && j < tarCards.Len(); i++ {
		if oriCards[i] == tarCards[j] {
			j++
		}
	}
	return j >= tarCards.Len()
}

// Copy 拷贝牌
func (c Cards) Copy() Cards {
	res := make(Cards, len(c))
	copy(res, c)
	return res
}

// DeleteCard 删除单张牌
func (c Cards) DeleteCard(card Card) (Cards, error) {
	// 判断牌是否在里面
	if !c.In(card) {
		return c, errors.New(fmt.Sprintf("card:%#v not in cards:%#v", card, c))
	}
	cards := c.Copy()
	res := make(Cards, 0, cards.Len()-1)
	for i := cards.Len() - 1; i >= 0; i-- {
		if cards[i] == card {
			res = append(cards[:i], cards[i+1:]...)
			break
		}
	}
	return res, nil
}

// DeleteCards 删除一些牌
func (c Cards) DeleteCards(cards Cards) (Cards, error) {
	if !c.IsContain(cards) {
		return nil, errors.New(fmt.Sprintf("cards:%#v not in cards:%#v", cards, c))
	}
	var err error
	res := c.Copy()
	for _, card := range cards {
		res, err = res.DeleteCard(card)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// AddCard 添加单张牌
func (c Cards) AddCard(card Card) Cards {
	res := c.Copy()
	res = append(res, card)
	return res
}

// AddCards 添加一些牌
func (c Cards) AddCards(cards Cards) Cards {
	res := c.Copy()
	res = append(res, cards...)
	return res
}

// GetCount 获取某张牌的数量
func (c Cards) GetCount(card Card) int {
	count := 0
	for _, card2 := range c {
		if card == card2 {
			count++
		}
	}
	return count
}

// IsSameColor 是否是同一花色
func (c Cards) IsSameColor() bool {
	if c.Len() == 0 {
		return true
	}
	c1 := c[0].GetColor()
	for _, card := range c {
		if c1 != card.GetColor() {
			return false
		}
	}
	return true
}

// Intersection 两个cards的交集
func (c Cards) Intersection(cards Cards) Cards {
	oriCards := c.Copy()
	tarCards := cards.Copy()
	oriCards.Sort()
	tarCards.Sort()
	res := make(Cards, 0, oriCards.Len())
	for i, j := 0, 0; i < oriCards.Len() && j < tarCards.Len(); {
		if oriCards[i] == tarCards[j] {
			res = append(res, oriCards[i])
			i++
			j++
		} else if oriCards[i] < tarCards[j] {
			i++
		} else {
			j++
		}
	}
	return res
}

// ToMap 转换成Map，每张牌对应的张数
func (c Cards) ToMap() map[Card]int32 {
	res := make(map[Card]int32, c.Len())
	for _, card := range c {
		res[card]++
	}
	return res
}

// ToUnique 去重
func (c Cards) ToUnique() Cards {
	res := make(Cards, 0, c.Len())
	for card, _ := range c.ToMap() {
		res = append(res, card)
	}
	return res
}

// RemoveCards 去除赖子牌
func (c Cards) RemoveCards(kingCards Cards) Cards {
	res := make(Cards, 0, c.Len())
	for i := 0; i < c.Len(); i++ {
		if !kingCards.In(c[i]) {
			res = append(res, c[i])
		}
	}
	return res
}

// GetAllColorCards 获取某种花色的所有牌
func (c Cards) GetAllColorCards(color CardColor) Cards {
	res := make(Cards, 0, c.Len())
	for _, card := range c {
		if card.GetColor() == color {
			res = append(res, card)
		}
	}
	return res
}
