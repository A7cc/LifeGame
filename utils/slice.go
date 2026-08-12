package utils

import "math/rand"

// RandomSample 返回 source 的随机样本，不修改原切片。
func RandomSample[T any](source []T, count int) []T {
	if count <= 0 || len(source) == 0 {
		return []T{}
	}

	if len(source) <= count {
		result := make([]T, len(source))
		copy(result, source)
		return result
	}

	// 使用 Fisher-Yates 洗牌算法
	shuffled := make([]T, len(source))
	copy(shuffled, source)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}
