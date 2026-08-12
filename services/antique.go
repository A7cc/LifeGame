package services

import (
	"LifeGame/core"
	"LifeGame/utils"
	"fmt"

	"math/rand"
)

// 获取古董
func (a *App) GetAntique(level int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 计算用户的名声
	if l, _, _, _ := core.CalcReputationLevel(a.Userinfo.UFame); l < level {
		return M{
			"code": -1,
			"msg":  "你的名声不足，不能参加拍卖会！",
		}
	}
	// 判断用户参加拍卖次数是否达到上限
	if a.Userinfo.UOpportunity.OANum >= a.Gameinfo.GMaxHoldNum.MARoundNum {
		return M{
			"code": -1,
			"msg":  "参加拍卖行次数已达到上限，不能在拍卖",
		}
	}
	// 初始化对应等级的古董信息
	antiqueDataList := []core.Antique{}
	// 获取古董信息，根据古董等级获取
	for _, v := range core.GetAntiques() {
		if v.ALevel == level {
			antiqueDataList = append(antiqueDataList, v)
		}
	}
	if len(antiqueDataList) == 0 {
		return M{
			"code": -1,
			"msg":  "当前等级没有可拍卖的古董",
		}
	}
	// 随机获取一个古董
	antiqueData := antiqueDataList[rand.Intn(len(antiqueDataList))]
	// 初始化古董信息
	antiqueInfo := core.AntiqueInfo{
		// 古董id
		AIId: antiqueData.AId,
		// 古董名字
		AIName: antiqueData.AName,
		// 古董价格
		AIPrice: core.GetAntiqueMinPrice(antiqueData),
		// 古董等级
		AILevel: antiqueData.ALevel,
		// 古董年限
		AITime: 0,
		// 古董默认为未鉴定（1是真，2是假，其他都是没有鉴定）
		AIDisplay: 0,
		// 古董稀有度
		AIMaterial: antiqueData.AMaterial,
		// 古董图片
		AIImg: antiqueData.AImg,
		// 古董描述
		AIDesc: antiqueData.ADesc,
		// 古董完好程度
		AICondition: utils.RandomInRange(0, 10),
	}
	// 古董最高价格，最高价格年限是4年，超过4年就会存在非常大的风险为假品
	antiqueInfo.AIPrice_max = core.CalculateMaxPrice(antiqueInfo.AIPrice, antiqueData.ALevel, antiqueData.AMaterial, utils.RandomInRange(6, 8), 4)
	// 将古董给Gameinfo
	a.Gameinfo.GAntiqueInfo = antiqueInfo
	// 记录当前用户参加拍卖古董次数+1
	a.Userinfo.UOpportunity.OANum++
	return M{
		"code":           200,
		"currentAntique": antiqueInfo,
		"oanum":          a.Userinfo.UOpportunity.OANum,
	}
}

// 结束拍卖
func (a *App) AuctionEnd(Price, Id int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 计算用户持有的古董数量
	if len(a.Userinfo.UAntique) >= a.Gameinfo.GMaxHoldNum.MAHoldNum {
		return M{
			"code": -1,
			"msg":  "古董数量已达到上限，不能在拍卖",
		}
	}
	// 判断古董是否合法
	if Price < a.Gameinfo.GAntiqueInfo.AIPrice {
		return M{
			"code": -1,
			"msg":  "竞拍价格不能低于古董最低价格",
		}
	}
	if Id != a.Gameinfo.GAntiqueInfo.AIId || a.Gameinfo.GAntiqueInfo.AIId == 0 {
		return M{
			"code": -1,
			"msg":  "不是当前拍卖的古董",
		}
	}

	// 判断用户是否有足够的钱
	if a.Userinfo.UCash < Price {
		return M{
			"code": -1,
			"msg":  "用户余额不足",
		}
	}

	// 设置古董的竞拍价格
	a.Gameinfo.GAntiqueInfo.AIPrice = Price
	// 扣除用户钱
	a.Userinfo.UCash -= Price
	// 添加古董到用户持有的古董列表中
	a.Userinfo.UAntique = append(a.Userinfo.UAntique, a.Gameinfo.GAntiqueInfo)
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + 2)
	// 计算用户资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	// 清空Gameinfo.GAntiqueInfo
	a.Gameinfo.GAntiqueInfo = core.AntiqueInfo{}
	return M{
		"code":     200,
		"userinfo": a.userSnapshot(),
	}
}

// 操作古董
func (a *App) OperationAntique(id, functype int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 判断古董是否合法
	okTmp := -1
	topInfo := ""
	UAntiqueTmp := a.Userinfo.UAntique
	for i, v := range UAntiqueTmp {
		if v.AIId == id {
			okTmp = i
			break
		}
	}
	if okTmp <= -1 {
		return M{
			"code": -1,
			"msg":  "古董不存在",
		}
	}
	// 获取古董信息
	antiqueTmp := UAntiqueTmp[okTmp]
	// 判断操作
	switch functype {
	case 1:
		// 判断古董是否已经鉴定
		if antiqueTmp.AIDisplay == 1 || antiqueTmp.AIDisplay == 2 {
			return M{
				"code": -1,
				"msg":  antiqueTmp.AIName + "已经鉴定过了",
			}
		}
		appraisalCost := max(100, antiqueTmp.AIPrice/100)
		if a.Userinfo.UCash < appraisalCost {
			return M{"code": -1, "msg": fmt.Sprintf("鉴定需要 %d 元，现金不足", appraisalCost)}
		}
		a.Userinfo.UCash -= appraisalCost
		// 随机鉴定结果，鉴定后持有时间不增加。
		antiqueTmp.AIDisplay = core.CalcFakeProbWithDifficulty(antiqueTmp, core.GetDifficultyConfig(a.Gameinfo.GDifficulty).AntiqueFake)
		a.Userinfo.UAntique[okTmp].AIDisplay = antiqueTmp.AIDisplay
		// 返回内容
		if antiqueTmp.AIDisplay == 2 {
			a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - 5)
			topInfo = antiqueTmp.AIName + "鉴定结果为 假"
		} else {
			topInfo = antiqueTmp.AIName + "鉴定结果为 真"
		}
	case 2:
		// 判断古董是否超过5
		if antiqueTmp.AICondition >= 5 {
			return M{
				"code": -1,
				"msg":  antiqueTmp.AIName + "已经不能修复了",
			}
		}
		repairCost := max(200, antiqueTmp.AIPrice*(5-antiqueTmp.AICondition)/200)
		if a.Userinfo.UCash < repairCost {
			return M{"code": -1, "msg": fmt.Sprintf("修复需要 %d 元，现金不足", repairCost)}
		}
		a.Userinfo.UCash -= repairCost
		a.Userinfo.UAntique[okTmp].AICondition = rand.Intn(6) + 5
		antiqueTmp.AICondition = a.Userinfo.UAntique[okTmp].AICondition
		topInfo = antiqueTmp.AIName + "修复结果为 " + fmt.Sprint(antiqueTmp.AICondition)
	case 3:
		if antiqueTmp.AIDisplay != 1 && antiqueTmp.AIDisplay != 2 {
			antiqueTmp.AIDisplay = core.CalcFakeProbWithDifficulty(antiqueTmp, core.GetDifficultyConfig(a.Gameinfo.GDifficulty).AntiqueFake)
		}
		// 存放古董的临时变量
		priceTmp := 0
		// 确定古董真假
		if antiqueTmp.AIDisplay == 1 {
			// 计算古董价格
			priceTmp = core.CalculateMaxPrice(antiqueTmp.AIPrice, antiqueTmp.AILevel, antiqueTmp.AIMaterial, antiqueTmp.AICondition, antiqueTmp.AITime)
		} else {
			a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - 5)
			// 古董是假的
			priceTmp = antiqueTmp.AIPrice / 10
			// 判断古董价格是否小于1
			if priceTmp < 1 {
				priceTmp = 1
			}
		}
		topInfo = antiqueTmp.AIName + "出售的价格为 " + fmt.Sprint(priceTmp) + " 元"
		// 修改古董的持有
		a.Userinfo.UAntique = append(a.Userinfo.UAntique[:okTmp], a.Userinfo.UAntique[okTmp+1:]...)
		// 给用户加钱
		a.Userinfo.UCash += priceTmp
		a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	default:
		return M{
			"code": -1,
			"msg":  "出现未知操作类型",
		}
	}
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"userinfo": a.userSnapshot(),
		"topinfo":  topInfo,
	}
}
