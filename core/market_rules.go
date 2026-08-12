package core

import (
	"LifeGame/utils"
	"errors"
	"math"
	"math/rand"
	"strings"
)

// 获取商品当前价格
func getNowPrice(item Item, display, inter, limitOk bool) (int, string) {
	// 先创建PriceTmp设置正常价格
	PriceTmp := item.IPrice
	// 常规行情控制在基准价 ±10%，事件行情控制在 60%~160%，
	// 避免种子中的极端上下限造成十倍乃至百倍套利。
	PriceTmpUp := int(math.Ceil(float64(PriceTmp) * 1.10))
	PriceTmpDown := int(math.Floor(float64(PriceTmp) * 0.90))
	eventMax := int(math.Ceil(float64(PriceTmp) * 1.60))
	eventMin := int(math.Floor(float64(PriceTmp) * 0.60))
	if eventMax < item.IPrice_max {
		item.IPrice_max = eventMax
	}
	if eventMin > item.IPrice_min {
		item.IPrice_min = eventMin
	}
	// 如果价格浮动范围超过最大最小值，则设置最大最小值
	if PriceTmpUp >= item.IPrice_max {
		PriceTmpUp = item.IPrice_max
	}
	if PriceTmpDown <= item.IPrice_min {
		PriceTmpDown = item.IPrice_min
	}
	// 随机获取效果里的信息，平价、上涨、下跌
	// 计算购买的国家风险权重
	var num int
	if inter {
		// 如果是国外，增加下跌的概率
		// 用权重控制随机生成的范围，使得num == 2（下跌）的概率更大
		weight := rand.Intn(100) // 生成一个0到100的随机数
		if weight < 55 {         // 55% 的概率为下跌（国际市场风险更高）
			num = 2
		} else if weight < 70 { // 15% 的概率为上涨
			num = 1
		} else { // 30% 的概率为平价
			num = 0
		}
	} else {
		// 国内的正常随机选择
		if len(item.IEffects) > 0 {
			num = rand.Intn(len(item.IEffects))
		}
	}
	// 如果Display为false或者LimitOk为false，则不做大的涨跌价格浮动
	if !display || !limitOk {
		num = 0
	}
	// 根据随机数选择价格变化方式，平价存在20%上下浮动
	switch num {
	case 1:
		PriceTmp = utils.RandomInRange(PriceTmpUp, item.IPrice_max)
	case 2:
		PriceTmp = utils.RandomInRange(item.IPrice_min, PriceTmpDown)
	default:
		PriceTmp = utils.RandomInRange(PriceTmpDown, PriceTmpUp)
	}
	return PriceTmp, item.IEffects[num]
}

// 刷新物资
func refreshItem(gitemTmp map[int]ItemInfo, itemInfoData []Item, inter bool) (map[int]ItemInfo, error) {
	// 生成随机索引列表用于显示可购买的物资
	itemCount := len(itemInfoData)
	if itemCount == 0 {
		return gitemTmp, errors.New("市场物资为空")
	}

	showMarketNum := ShowMarketNum
	if itemCount < showMarketNum {
		showMarketNum = itemCount
	}

	// 生成随机索引
	randomIndices := make(map[int]struct{})
	for len(randomIndices) < showMarketNum {
		idx := rand.Intn(itemCount)
		randomIndices[idx] = struct{}{}
	}

	// 用于设置出现涨幅的物资数目限制, 最多为物资数量的1/2
	ItemLimit := int(showMarketNum / 2)
	// 刷新国内物资
	for idx, item := range itemInfoData {
		// 判断是否在随机列表中，如果在则显示
		_, DisplayTmp := randomIndices[idx]
		if DisplayTmp || len(gitemTmp) == 0 {
			// 设置是否允许涨幅
			ItemLimitOk := false
			if ItemLimit > 0 {
				ItemLimitOk = true
			}
			// 设置价格
			PriceTmp, EffectTmp := getNowPrice(item, DisplayTmp, inter, ItemLimitOk)
			// 只限制真正上涨的商品，跌价事件不占用上涨名额。
			if strings.Contains(EffectTmp, "📈") || strings.Contains(EffectTmp, "🔥") {
				ItemLimit--
			}
			// 赋值
			gitemTmp[item.IId] = ItemInfo{
				IIName:    item.IName,
				IIPrice:   PriceTmp,
				IIEffect:  EffectTmp,
				IIDisplay: DisplayTmp,
			}
		} else {
			// 赋值
			gitemTmp[item.IId] = ItemInfo{
				IIName:  item.IName,
				IIPrice: gitemTmp[item.IId].IIPrice,
				// 这里是值是从0开始的
				IIEffect:  itemInfoData[idx].IEffects[0],
				IIDisplay: false,
			}
		}
	}
	return gitemTmp, nil
}

// 生成当前公司的价格
func GenerateStartPrice(cPrice int, cRisk int) int {
	// 计算最大浮动百分比，比如CRisk=5时最大±10%
	maxFluctuation := float64(cRisk) * CompanyFluct

	// 随机浮动值：[-maxFluctuation, +maxFluctuation]
	randomFluctuation := (rand.Float64()*2 - 1) * maxFluctuation

	// 初始价格 = 基础价格 * (1 + 随机浮动)
	startPrice := float64(cPrice) * (1 + randomFluctuation)

	// 保证价格最少为1
	if startPrice < 1 {
		startPrice = 1
	}

	return int(startPrice)
}

// 计算房屋当前价格（根据市场波动）
func CalculateHousePrice(house House) int {
	// 价格在基础价格 ±20% 范围内波动
	priceRange := house.HPrice / 5 // 20%
	minPrice := house.HPrice - priceRange
	maxPrice := house.HPrice + priceRange

	// 确保不超过设定范围
	if minPrice < house.HPrice_min {
		minPrice = house.HPrice_min
	}
	if maxPrice > house.HPrice_max {
		maxPrice = house.HPrice_max
	}

	// 确保价格范围有效
	if minPrice >= maxPrice {
		minPrice = house.HPrice_min
		maxPrice = house.HPrice_max
		// 如果仍然无效，返回基础价格
		if minPrice >= maxPrice {
			return house.HPrice
		}
	}

	// 随机生成当前价格
	return utils.RandomInRange(minPrice, maxPrice)
}

// 计算车辆当前价格（根据市场波动）
func CalculateCarPrice(car Car) int {
	// 价格在基础价格 ±25% 范围内波动（车辆贬值更快）
	priceRange := car.CPrice / 4
	minPrice := car.CPrice - priceRange
	maxPrice := car.CPrice + priceRange

	// 确保不超过设定范围
	if minPrice < car.CPrice_min {
		minPrice = car.CPrice_min
	}
	if maxPrice > car.CPrice_max {
		maxPrice = car.CPrice_max
	}

	// 确保价格范围有效
	if minPrice >= maxPrice {
		minPrice = car.CPrice_min
		maxPrice = car.CPrice_max
		// 如果仍然无效，返回基础价格
		if minPrice >= maxPrice {
			return car.CPrice
		}
	}

	// 随机生成当前价格
	return utils.RandomInRange(minPrice, maxPrice)
}

// CalculateLifestyleFame 将房车的永久声望压缩到可持续区间，
// 防止一次购买直接跳过全部声望成长。
func CalculateLifestyleFame(raw int) int {
	if raw <= 0 {
		return 0
	}
	return int(math.Ceil(float64(raw) / 3.0))
}

// CalculateLifestyleHealth 房车仍提供生活品质收益，但不再轻易把免疫力堆满。
func CalculateLifestyleHealth(raw int) int {
	if raw <= 1 {
		return max(raw, 0)
	}
	return int(math.Ceil(float64(raw) * 0.6))
}
