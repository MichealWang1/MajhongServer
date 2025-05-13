package ai

import (
	lib "kxmj.common/mahjong"
	"sort"
	"strings"
)

// 牌相关处理

var CardToTilesMap = map[lib.Card]int32{
	0x01: 0, 0x02: 1, 0x03: 2, 0x04: 3, 0x05: 4, 0x06: 5, 0x07: 6, 0x08: 7, 0x09: 8, // 一万到九万
	0x11: 9, 0x12: 10, 0x13: 11, 0x14: 12, 0x15: 13, 0x16: 14, 0x17: 15, 0x18: 16, 0x19: 17, // 一条到九条
	0x21: 18, 0x22: 19, 0x23: 20, 0x24: 21, 0x25: 22, 0x26: 23, 0x27: 24, 0x28: 25, 0x29: 26, // 一筒到九筒
	0x31: 27, 0x32: 28, 0x33: 29, 0x34: 30, 0x35: 31, 0x36: 32, 0x37: 33, //东南西北中发白
}

var TileToCardMap = map[int32]lib.Card{
	0: 0x01, 1: 0x02, 2: 0x03, 3: 0x04, 4: 0x05, 5: 0x06, 6: 0x07, 7: 0x08, 8: 0x09, // 一万到九万
	9: 0x11, 10: 0x12, 11: 0x13, 12: 0x14, 13: 0x15, 14: 0x16, 15: 0x17, 16: 0x18, 17: 0x19, // 一条到九条
	18: 0x21, 19: 0x22, 20: 0x23, 21: 0x24, 22: 0x25, 23: 0x26, 24: 0x27, 25: 0x28, 26: 0x29, // 一筒到九筒
	27: 0x31, 28: 0x32, 29: 0x33, 30: 0x34, 31: 0x35, 32: 0x36, 33: 0x37, //东南西北中发白
}

const (
	MaxWTTTile int32 = 26
)

// CardToTiles
func CardToTiles(cards lib.Cards) []int32 {
	tiles := make([]int32, 34)
	for _, card := range cards {
		tiles[CardToTilesMap[card]]++
	}
	return tiles
}

// TileToCards
func TileToCards(tiles []int32) lib.Cards {
	res := make(lib.Cards, 0, 14)
	for t, count := range tiles {
		card := TileToCardMap[int32(t)]
		res = append(res, card.Repeat(int(count))...)
	}
	return res
}

// []{tile, tile}转换成cards
func Int32ToCards(in []int32) lib.Cards {
	res := make(lib.Cards, 0, 14)
	for _, v := range in {
		res = append(res, TileToCardMap[v])
	}
	return res
}

// 通过手牌和弃牌获取剩余牌tiles
func GetReducedTiles(handTiles, discardTiles []int32) []int32 {
	res := make([]int32, 34)
	for i := 0; i < len(res); i++ {
		res[i] = 4 - handTiles[i] - discardTiles[i]
	}
	return res
}

// CountInt32Map 计算map[牌]数量 里数量
func CountInt32Map(m map[int32]int32) int32 {
	count := int32(0)
	for _, num := range m {
		count += num
	}
	return count
}

// TilesToKey 根据Tiles构造唯一key值
func TilesToKey(tiles []int32) string {
	var builder strings.Builder
	builder.Grow(len(tiles))
	for _, i := range tiles {
		builder.WriteByte(byte('0' + i))
	}
	return builder.String()
}

// GetTilesValue 获取当前位置的牌值
func GetTilesValue(i int32) int32 {
	return i%9 + 1
}

// CountOfTiles 获取手中数量
func CountOfTiles(handTiles []int32) int32 {
	count := int32(0)
	for _, num := range handTiles {
		count += num
	}
	return count
}

func CardsToIndex(cards lib.Cards) []int32 {
	kingIndex := make([]int32, 0)
	for _, card := range cards {
		index := CardToTilesMap[card]
		kingIndex = append(kingIndex, index)
	}
	return kingIndex
}

func SortCombs(combs [][]int32) [][]int32 {
	res := make([][]int32, len(combs))
	for i := range combs {
		res[i] = make([]int32, len(combs[i]))
		copy(res[i], combs[i])
	}
	sort.SliceIsSorted(res, func(i, j int) bool {
		return len(res[i]) < len(res[j]) && res[i][0] < res[j][0]
	})
	return res
}

// tiles牌中花色数量顺序：下标0万1条2筒3字 对应 排行(不计算字)
func TilesToColorSort(tiles []int32) []int32 {
	colors := make(map[int32]int32)
	for i := int32(0); i < 3; i++ {
		colors[i] = 0
	}
	for i := 0; i <= int(MaxWTTTile); i++ {
		if tiles[i] == 0 {
			continue
		}
		colors[int32(i/9)]++
	}
	res := make([]int32, 4, 4)
	for i := 0; i < 3; i++ {
		tmpColor := int32(-1)
		maxColorCount := int32(-1)
		for k, v := range colors {
			if tmpColor == -1 || maxColorCount < v {
				tmpColor = k
				maxColorCount = v
			}
		}
		if tmpColor == -1 {
			break
		}
		res[tmpColor] = int32(i)
		delete(colors, tmpColor)
	}
	res[int32(3)] = 3
	return res
}

// 是否是万
func TileIsWan(i int) bool {
	return i <= 8
}

// 是否是条
func TileIsTiao(i int) bool {
	return i >= 9 && i <= 17
}

// 是否是筒
func TileIsTong(i int) bool {
	return i >= 18 && i <= 26
}

// TileIsWTT 是否是万条筒
func TileIsWTT(i int) bool {
	return i < 27
}

// TileIsDNXB 是否是东南西北
func TileIsDNXB(i int) bool {
	return i >= 27 && i <= 30
}

// TileIsZFB 是否是中发白
func TileIsZFB(i int) bool {
	return i >= 30 && i <= 33
}

// 是否是子
func TileIsZi(i int) bool {
	return i >= 27 && i <= 33
}
