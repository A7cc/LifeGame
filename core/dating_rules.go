package core

import (
	"math/rand"
	"strconv"
)

const (
	DatingStatusStranger   = "陌生人"
	DatingStatusFriend     = "朋友"
	DatingStatusAmbiguous  = "暧昧中"
	DatingStatusDating     = "交往中"
	DatingStatusLover      = "恋人"
	DatingStatusExclusive  = "专属恋人"
	DatingStatusSweetheart = "爱人"
	DatingStatusMarried    = "已婚"
	DatingStatusFormer     = "前任"
)

// 获取用户约会次数
func getTotalDates(userinfo *User) int {
	totalDates := 0
	for _, dating := range userinfo.UDating {
		totalDates += dating.DCount
	}
	return totalDates
}

// 获取用户拥有的商品数量
func getItemOwnCount(userinfo *User, target string) int {
	// 计算指定商品数量
	if target != "" {
		itemId, err := strconv.Atoi(target)
		if err != nil {
			return 0
		}
		return userinfo.UItemins[itemId].UIINum + userinfo.UItemout[itemId].UIINum
	}
	// 计算所有商品数量
	total := 0
	for _, item := range userinfo.UItemins {
		total += item.UIINum
	}
	for _, item := range userinfo.UItemout {
		total += item.UIINum
	}
	return total
}

// 获取用户拥有的古董数量
func getRareAntiqueCount(userinfo *User) int {
	total := 0
	for _, antique := range userinfo.UAntique {
		// 古董必须有为真，且为珍品
		if antique.AIDisplay == 1 && antique.AIMaterial >= 4 {
			total++
		}
	}
	return total
}

// 获取用户游玩小游戏的信息
func getMiniGameRecord(userinfo *User, target string) MiniGameRecord {
	if userinfo.UMiniGameRecords == nil {
		return MiniGameRecord{}
	}
	if target == "" {
		total := MiniGameRecord{}
		for _, record := range userinfo.UMiniGameRecords {
			total.PlayCount += record.PlayCount
			total.WinCount += record.WinCount
		}
		return total
	}
	return userinfo.UMiniGameRecords[target]
}

// 检查约会对象解锁条件
func CheckDatingUnlock(userinfo *User, dating DatingInfo) bool {
	oknum := 0
	// 检查每个条件是否满足
	for _, condition := range dating.DMeetConditions {
		unlocked := false
		if condition.CType == "random" {
			// 随机邂逅在同一游戏年度内保持稳定，避免反复打开页面刷结果；
			// 年龄变化后会得到新的判定。
			roll := (dating.DId*37 + userinfo.UAge*17 + userinfo.Uid*13) % 100
			unlocked = roll < condition.CValue
		} else {
			unlocked = CheckCondition(userinfo, condition)
		}
		if unlocked {
			oknum++
		}
	}
	if oknum >= len(dating.DMeetConditions) {
		return true
	} else {
		return false
	}
}

// 检查单个条件
func CheckCondition(userinfo *User, condition MeetCondition) bool {
	switch condition.CType {
	case "fame":
		// 名声要求
		if userinfo.UFame >= condition.CValue {
			return true
		}
		return false

	case "cash":
		// 现金要求
		if userinfo.UCash >= condition.CValue {
			return true
		}
		return false

	case "bank":
		// 存款要求
		if userinfo.UBank >= condition.CValue {
			return true
		}
		return false

	case "house":
		// 拥有房子是否满足
		if len(userinfo.UHouse) >= condition.CValue {
			return true
		}
		return false

	case "car":
		// 拥有车子是否满足
		if len(userinfo.UCar) >= condition.CValue {
			return true
		}
		return false

	case "play_game":
		// 玩过游戏
		if getMiniGameRecord(userinfo, condition.CTarget).PlayCount >= condition.CValue {
			return true
		}
		return false

	case "win_game":
		// 游戏获胜次数
		if getMiniGameRecord(userinfo, condition.CTarget).WinCount >= condition.CValue {
			return true
		}
		return false

	case "work_count":
		// 打工次数
		total := 0
		for _, record := range userinfo.UMiniGameRecords {
			if record.MGRType == "work" {
				total += record.PlayCount
			}
		}
		if total >= condition.CValue {
			return true
		}
		return false

	case "age":
		// 年龄要求
		if userinfo.UAge >= condition.CValue {
			return true
		}
		return false

	case "random":
		// 随机遇见
		if rand.Float64() < float64(condition.CValue)/100.0 {
			return true
		}
		return false

	case "date_count":
		// 约会次数要求
		if getTotalDates(userinfo) >= condition.CValue {
			return true
		}
		return false

	case "immunity":
		// 免疫力要求
		if userinfo.UImmunity >= condition.CValue {
			return true
		}
		return false

	case "antique_rare":
		// 拥有稀有古董
		if getRareAntiqueCount(userinfo) >= condition.CValue {
			return true
		}
		return false

	case "lottery_win":
		// 彩票中奖
		if getMiniGameRecord(userinfo, "lottery").WinCount >= condition.CValue {
			return true
		}
		return false

	case "item_own":
		// 拥有物品
		if getItemOwnCount(userinfo, condition.CTarget) >= condition.CValue {
			return true
		}
		return false

	case "stock_profit":
		// 当年股票盈利
		if userinfo.UStockProfit >= condition.CValue {
			return true
		}
		return false

	case "company_founder":
		// 创业者
		if len(userinfo.UCompany) > 0 {
			return true
		}
		return false

	default:
		return false
	}
}

// 获取约会状态，already表示是否已经有爱人了，不能有第二个
func GetDatingStatus(affinity, count int, already bool) string {
	// 爱人
	if count >= 20 && affinity >= 90 && !already {
		return DatingStatusSweetheart
	}
	// 专属
	if count >= 15 && affinity >= 70 {
		return DatingStatusExclusive
	}
	// 恋人
	if count >= 10 && affinity >= 50 {
		return DatingStatusLover
	}
	// 交往中
	if affinity >= 30 {
		return DatingStatusDating
	}
	// 暧昧中
	if count >= 5 && affinity >= 20 {
		return DatingStatusAmbiguous
	}
	// 朋友
	if affinity >= 10 {
		return DatingStatusFriend
	}
	// 前任
	if count >= 10 && affinity < 0 {
		return DatingStatusFormer
	}
	// 默认为陌生人
	return DatingStatusStranger
}

// NormalizeDatingStatus 将空值和非法状态收敛到当前关系状态集合。
func NormalizeDatingStatus(status string) string {
	switch status {
	case "", DatingStatusStranger:
		return DatingStatusStranger
	case DatingStatusFriend:
		return DatingStatusFriend
	case DatingStatusAmbiguous:
		return DatingStatusAmbiguous
	case DatingStatusDating:
		return DatingStatusDating
	case DatingStatusLover:
		return DatingStatusLover
	case DatingStatusExclusive:
		return DatingStatusExclusive
	case DatingStatusSweetheart:
		return DatingStatusSweetheart
	case DatingStatusMarried:
		return DatingStatusMarried
	case DatingStatusFormer:
		return DatingStatusFormer
	default:
		return DatingStatusStranger
	}
}

// 计算约会成功率
func CalculateDatingSuccessRate(dating DatingInfo, userDating UserDatingInfo) float64 {
	baseRate := 0.6 // 基础成功率60%

	// 好感度加成：好感度每增加10点，成功率增加5%
	affinityBonus := float64(userDating.DAffinity) / 10.0 * 0.05

	// 约会次数加成：约会次数每增加5次，成功率增加3%
	countBonus := float64(userDating.DCount) / 5.0 * 0.03
	successRate := baseRate + affinityBonus + countBonus

	// 直接使用约会消费档位表达人物难度，避免简介文案中的偶然关键词
	// 改变成功率。
	personalityBonus := 0.0
	switch {
	case dating.DCost <= 1000:
		personalityBonus = 0.10
	case dating.DCost <= 2000:
		personalityBonus = 0.05
	case dating.DCost <= 5000:
		personalityBonus = 0
	case dating.DCost <= 8000:
		personalityBonus = -0.05
	default:
		personalityBonus = -0.10
	}
	successRate = successRate + personalityBonus

	// 限制在5%-95%之间
	if successRate < 0.05 {
		successRate = 0.05
	}
	if successRate > 0.95 {
		successRate = 0.95
	}

	return successRate
}
