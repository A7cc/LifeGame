package services

import (
	"LifeGame/core"
	"fmt"
	"math"
)

// 存钱取钱操作
func (a *App) OperationMoney(ope string, count int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	if count <= 0 {
		return M{
			"code": -1,
			"msg":  "操作金额必须大于0",
		}
	}

	msg := ""
	switch ope {
	case "deposit":
		// 存款
		if count > a.Userinfo.UCash {
			return M{
				"code": -1,
				"msg":  "用户余额不足",
			}
		}
		a.Userinfo.UCash -= count
		a.Userinfo.UBank += count
		msg = "存款"
	case "withdraw":
		// 取款
		if count > a.Userinfo.UBank {
			return M{
				"code": -1,
				"msg":  "银行余额不足",
			}
		}
		a.Userinfo.UBank -= count
		a.Userinfo.UCash += count
		msg = "取款"
	default:
		return M{
			"code": -1,
			"msg":  "操作类型错误",
		}
	}
	// 记录任务统计
	a.recordBankTaskStats(ope, count)
	// 计算用户资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("已%s %d 元", msg, count),
		"userinfo": a.userSnapshot(),
	}
}

// 申请贷款
func (a *App) ApplyLoan(amount int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 检查是否已有贷款
	if a.Userinfo.ULoan > 0 {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("你已有贷款 %d 元未还清，无法再次申请", a.Userinfo.ULoan),
		}
	}
	// 检查贷款金额是否合法
	if amount <= 0 {
		return M{
			"code": -1,
			"msg":  "贷款金额必须大于0",
		}
	}
	// 获取最大贷款额度
	maxLoan := core.GetMaxLoanAmount(a.Userinfo, a.Gameinfo)

	if amount > maxLoan {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("超过最大贷款额度！根据你的资产，最多可贷 %d 元", maxLoan),
		}
	}

	// 发放贷款
	a.Userinfo.ULoan += amount
	a.Userinfo.UBank += amount
	// 记录任务统计
	a.recordBankTaskStats("loan", amount)
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功贷款 %d 元，年利率10%%，记得按时还款", amount),
		"userinfo": a.userSnapshot(),
	}
}

// 获取银行任务列表
func (a *App) GetBankTaskList() H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 获取当前年度的任务
	currentTasks := a.Gameinfo.GBankTaskStats.CurrentTasks
	taskStatusList := make([]map[string]interface{}, 0)

	for _, task := range currentTasks {
		isCompleted := core.CheckBankTaskProgress(task, a.Gameinfo.GBankTaskStats)
		isClaimed := false
		for _, claimedId := range a.Gameinfo.GBankTaskStats.ClaimedTasks {
			if claimedId == task.TaskId {
				isClaimed = true
				break
			}
		}
		taskStatusList = append(taskStatusList, map[string]interface{}{
			"taskid":    task.TaskId,
			"taskname":  task.TaskName,
			"taskdesc":  task.TaskDesc,
			"tasktype":  task.TaskType,
			"target":    task.TargetValue,
			"reward":    task.Reward,
			"completed": isCompleted,
			"claimed":   isClaimed,
			"current":   core.GetTaskCurrentValue(task.TaskType, a.Gameinfo.GBankTaskStats),
		})
	}

	return M{
		"code":     200,
		"tasklist": taskStatusList,
	}
}

// 领取任务奖励
func (a *App) ClaimTaskReward(taskId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	var targetTask *core.BankTask

	// 在当前年度任务中查找
	for i := range a.Gameinfo.GBankTaskStats.CurrentTasks {
		if a.Gameinfo.GBankTaskStats.CurrentTasks[i].TaskId == taskId {
			targetTask = &a.Gameinfo.GBankTaskStats.CurrentTasks[i]
			break
		}
	}

	if targetTask == nil {
		return M{
			"code": -1,
			"msg":  "任务不存在",
		}
	}

	// 检查任务是否已完成
	if !core.CheckBankTaskProgress(*targetTask, a.Gameinfo.GBankTaskStats) {
		return M{
			"code": -1,
			"msg":  "任务未完成，无法领取奖励",
		}
	}

	// 检查是否已经领取过
	for _, claimedId := range a.Gameinfo.GBankTaskStats.ClaimedTasks {
		if claimedId == taskId {
			return M{
				"code": -1,
				"msg":  "奖励已领取",
			}
		}
	}

	// 发放奖励
	a.Userinfo.UCash += targetTask.Reward
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	a.Gameinfo.GBankTaskStats.ClaimedTasks = append(a.Gameinfo.GBankTaskStats.ClaimedTasks, taskId)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功领取任务奖励：%d 元！", targetTask.Reward),
		"reward":   targetTask.Reward,
		"userinfo": a.userSnapshot(),
	}
}

// 记录任务统计
func (a *App) recordBankTaskStats(taskType string, amount int) {
	switch taskType {
	case "deposit":
		a.Gameinfo.GBankTaskStats.NetBankFlow += amount
		if a.Gameinfo.GBankTaskStats.NetBankFlow > a.Gameinfo.GBankTaskStats.DepositAmount {
			a.Gameinfo.GBankTaskStats.DepositAmount = a.Gameinfo.GBankTaskStats.NetBankFlow
			a.Gameinfo.GBankTaskStats.DepositCount++
		}
	case "withdraw":
		a.Gameinfo.GBankTaskStats.NetBankFlow -= amount
		netWithdrawal := -a.Gameinfo.GBankTaskStats.NetBankFlow
		if netWithdrawal > a.Gameinfo.GBankTaskStats.WithdrawAmount {
			a.Gameinfo.GBankTaskStats.WithdrawAmount = netWithdrawal
			a.Gameinfo.GBankTaskStats.WithdrawCount++
		}
	case "loan":
		a.Gameinfo.GBankTaskStats.LoanAmount += amount
	case "work":
		a.Gameinfo.GBankTaskStats.WorkCount++
	}
}

// 偿还贷款
func (a *App) RepayLoan(amount int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 检查是否有贷款
	if a.Userinfo.ULoan <= 0 {
		return M{
			"code": -1,
			"msg":  "你没有需要偿还的贷款",
		}
	}

	// 计算应还金额（本金+利息）
	interest := int(math.Ceil(float64(a.Userinfo.ULoan) * core.LoanInterestRate))
	totalDue := a.Userinfo.ULoan + interest

	// 当前产品只支持一次性结清，避免只归还本金从而绕过利息。
	if amount <= 0 {
		return M{
			"code": -1,
			"msg":  "还款金额必须大于0",
		}
	}

	if amount < totalDue {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("当前仅支持一次性结清，需要 %d 元（本金%d + 利息%d）", totalDue, a.Userinfo.ULoan, interest),
		}
	}

	// 检查存款是否足够
	if a.Userinfo.UBank < totalDue {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("存款不足！需要 %d 元（本金%d + 利息%d）", totalDue, a.Userinfo.ULoan, interest),
		}
	}

	// 一次性扣除本金和利息并结清贷款。
	a.Userinfo.UBank -= totalDue
	a.Userinfo.ULoan = 0
	a.Userinfo.ULoanOverdue = 0

	// 更新资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	msg := fmt.Sprintf("成功偿还本金 %d 元及利息 %d 元，贷款已全部还清！", totalDue-interest, interest)

	return M{
		"code":     200,
		"msg":      msg,
		"userinfo": a.userSnapshot(),
	}
}
