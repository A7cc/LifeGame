package core

import (
	"math"
	"math/rand"
)

// 获取古董在拍卖行的起始最低价
func GetAntiqueMinPrice(antique Antique) int {
	rarityFactor := 1.0 + float64(antique.AMaterial)*0.1
	levelFactor := 1.0 + float64(antique.ALevel)*0.2
	basePrice := float64(antique.APrice) * rarityFactor * levelFactor

	// 添加浮动 ±10%
	fluctuation := 0.9 + rand.Float64()*0.2
	minPrice := basePrice * fluctuation

	return int(minPrice)
}

// 获取古董在拍卖行的最高价
func CalculateMaxPrice(aiPrice, aiLevel, aiMaterial, aiCondition, aitime int) int {
	// 收藏收益在第5年达到高点，之后因保管风险逐年回落。
	rarityFactor := 1.0 + float64(aiMaterial)*0.02
	conditionFactor := 0.85 + float64(aiCondition)*0.015
	levelFactor := 1.0 + float64(aiLevel)*0.05
	randomFactor := 0.95 + rand.Float64()*0.10
	holdYears := min(aitime, 5)
	timeFactor := 1.0 + float64(holdYears)*0.03
	if aitime > 5 {
		timeFactor -= float64(aitime-5) * 0.04
		if timeFactor < 0.8 {
			timeFactor = 0.8
		}
	}
	priceMax := float64(aiPrice) * rarityFactor * conditionFactor * levelFactor * randomFactor * timeFactor
	// 真品也可能因行情不佳小幅折价，但不再保证三倍暴利。
	minPrice := float64(aiPrice) * 0.8
	if priceMax < minPrice {
		priceMax = minPrice
	}

	return int(priceMax)
}

// CalcFakeProbWithDifficulty 让古董真伪风险真正受游戏难度影响。
func CalcFakeProbWithDifficulty(info AntiqueInfo, difficultyMultiplier float64) int {
	ageRisk := math.Abs(float64(info.AITime)-5) * 0.035
	fakeChance := 0.32 + float64(info.AIMaterial-5)*0.015 - float64(info.AILevel)*0.04 + ageRisk
	fakeChance *= difficultyMultiplier
	if fakeChance < 0.1 {
		fakeChance = 0.1
	}
	if fakeChance > 0.75 {
		fakeChance = 0.75
	}
	if rand.Float64() < fakeChance {
		return 2
	}
	return 1
}
