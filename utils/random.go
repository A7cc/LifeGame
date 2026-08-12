package utils

import "math/rand"

// RandomInRange 返回闭区间 [min, max] 的随机整数。
func RandomInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}

// RandFloat 返回闭区间 [min, max] 的随机浮点数。
func RandFloat(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	return min + rand.Float64()*(max-min)
}
