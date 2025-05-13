package utils

import "strconv"

// UniqueInt32 去重
func UniqueInt32(list []int32) []int32 {
	var result []int32
	m := map[int32]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueUint32 去重
func UniqueUint32(list []uint32) []uint32 {
	var result []uint32
	m := map[uint32]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueUint8 去重
func UniqueUint8(list []uint8) []uint8 {
	var result []uint8
	m := map[uint8]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueInt 去重
func UniqueInt(list []int) []int {
	var result []int
	m := map[int]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueUint 去重
func UniqueUint(list []uint) []uint {
	var result []uint
	m := map[uint]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueInt64 去重
func UniqueInt64(list []int64) []int64 {
	var result []int64
	m := map[int64]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// UniqueString 去重
func UniqueString(list []string) []string {
	var result []string
	m := map[string]struct{}{}
	for _, v := range list {
		m[v] = struct{}{}
	}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

func ContainsUint32(elems []uint32, elem uint32) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func ContainsString(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func StringSliceToUint32(elems []string) []uint32 {
	var result []uint32
	for _, e := range elems {

		k, _ := strconv.ParseUint(e, 0, 64)

		result = append(result, uint32(k))
	}
	return result
}
