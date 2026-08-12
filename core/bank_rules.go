package core

import (
	"LifeGame/utils"
	"fmt"
	"math"
)

// CheckBankTaskProgress 检查银行任务进度（返回任务是否完成）
func CheckBankTaskProgress(task BankTask, stats BankTaskStats) bool {
	// 根据任务类型检查进度
	switch task.TaskType {
	case "deposit":
		return stats.DepositAmount >= task.TargetValue
	case "depositcount":
		return stats.DepositCount >= task.TargetValue
	case "withdraw":
		return stats.WithdrawAmount >= task.TargetValue
	case "withdrawcount":
		return stats.WithdrawCount >= task.TargetValue
	case "loan":
		return stats.LoanAmount >= task.TargetValue
	case "work":
		return stats.WorkCount >= task.TargetValue
	default:
		return false
	}
}

// ResetBankTaskStats 重置银行任务统计（新年时调用）
func ResetBankTaskStats(gameinfo *Game) {
	// 获取所有可用任务
	allTasks := GetBankTasks()

	// 随机抽取3个任务
	selectedTasks := utils.RandomSample(allTasks, 3)

	gameinfo.GBankTaskStats = BankTaskStats{
		NetBankFlow:    0,
		DepositAmount:  0,
		DepositCount:   0,
		WithdrawAmount: 0,
		WithdrawCount:  0,
		LoanAmount:     0,
		WorkCount:      0,
		CompletedTasks: []int{},
		ClaimedTasks:   []int{},
		CurrentTasks:   selectedTasks,
	}
}

// ProcessDepositAnnual 为真实留在银行中的存款结算利息。
func ProcessDepositAnnual(userinfo *User) string {
	if userinfo == nil || userinfo.UBank <= 0 {
		return ""
	}
	interest := int(math.Floor(float64(userinfo.UBank) * DepositInterestRate))
	if interest <= 0 {
		return ""
	}
	userinfo.UBank += interest
	return fmt.Sprintf("🏦 银行存款获得年度利息 %d 元", interest)
}

// 获取任务当前值
func GetTaskCurrentValue(taskType string, stats BankTaskStats) int {
	switch taskType {
	case "deposit":
		return stats.DepositAmount
	case "depositcount":
		return stats.DepositCount
	case "withdraw":
		return stats.WithdrawAmount
	case "withdrawcount":
		return stats.WithdrawCount
	case "loan":
		return stats.LoanAmount
	case "work":
		return stats.WorkCount
	default:
		return 0
	}
}

// 获取最大贷款额度（综合考虑车、房、古董、名声等级、股票额度）
func GetMaxLoanAmount(userinfo *User, gameinfo *Game) int {
	// 计算名声等级
	level, _, stockLimit, _ := CalcReputationLevel(userinfo.UFame)

	// 老赖不能贷款
	if level < 0 {
		return 0
	}

	// 1. 计算车辆总价值（按市价70%计算）
	carTotalValue := 0
	for _, c := range gameinfo.GCarInfo {
		if userinfo.UCar[c.CIId] {
			carTotalValue += c.CIPrice * 7 / 10
		}
	}

	// 2. 计算房子总价值（按市价80%计算）
	houseTotalValue := 0
	for _, h := range gameinfo.GHouseInfo {
		if userinfo.UHouse[h.HIId] {
			houseTotalValue += h.HIPrice * 8 / 10
		}
	}

	// 3. 计算公司总价值（按市价80%计算）
	companyTotalValue := 0
	for _, c := range userinfo.UCompany {
		companyTotalValue += c.UCompanyCostPrice * c.UCompanyNum * 8 / 10
	}

	// 4. 获取最贵古董的价格
	maxAntiquePrice := 0
	for _, antique := range userinfo.UAntique {
		if antique.AIPrice > maxAntiquePrice {
			maxAntiquePrice = antique.AIPrice
		}
	}

	// 5. 名声等级对应的股票额度基础值
	baseAmount := stockLimit / 4

	// 6. 根据等级调整倍数
	var multiplier int
	switch level {
	case 0: // 普通
		multiplier = 1
	case 1: // 中等
		multiplier = 1
	case 2: // 高级
		multiplier = 2
	case 3: // 豪华
		multiplier = 3
	case 4: // 私人
		multiplier = 5
	default:
		multiplier = 0
	}

	// 7. 综合计算最大贷款额度
	// 基础额度（基于名声和股票额度）
	baseLoan := baseAmount * multiplier

	// 抵押物价值：车房公司总价值的50% + 最贵古董的80%
	collateralValue := (carTotalValue+houseTotalValue+companyTotalValue)/2 + maxAntiquePrice*80/100

	// 最大贷款额度 = 基础额度 + 抵押物价值
	maxLoan := baseLoan + collateralValue

	return maxLoan
}

// 处理年度贷款利息和逾期（在NextTime中调用）
func ProcessLoanAnnual(userinfo *User, gameinfo *Game) (string, bool) {
	// 检查是否有贷款
	if userinfo.ULoan <= 0 {
		return "", true
	}

	// 计算利息，贷款利率（年利率10%）
	interest := int(math.Ceil(float64(userinfo.ULoan) * LoanInterestRate))

	// 检查是否有足够现金还利息
	if userinfo.UBank >= interest {
		userinfo.UBank -= interest
		userinfo.UAssets = CalculateUserAssets(userinfo, gameinfo)
		return fmt.Sprintf("💰 银行自动从你的账户扣除贷款利息 %d 元", interest), true
	} else {
		// 存款不足，记录本次逾期；超过10次立即破产。
		userinfo.ULoanOverdue++
		if userinfo.ULoanOverdue > 10 {
			return "", false
		}

		// 逾期惩罚：贷款金额增加，逾期惩罚（每次逾期贷款增加20%利息）
		penalty := int(math.Ceil(float64(userinfo.ULoan) * OverduePenalty))
		userinfo.ULoan += penalty

		// 逾期影响名声
		userinfo.UFame = CalcFame(userinfo.UFame - 5)
		// 逾期影响免疫力
		userinfo.UImmunity = CalcImmunity(userinfo.UImmunity - 5)

		return fmt.Sprintf("⚠️ 贷款逾期！利息 %d 元无法偿还，贷款增加 %d 元，名声-5，免疫力-5", interest, penalty), true
	}
}
