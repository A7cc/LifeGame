package services

import (
	"LifeGame/core"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (a *App) InitGame(name string, sex bool, difficulty int) GameStateResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireStartup(); errResp != nil {
		return gameStateError(responseMessage(errResp))
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return gameStateError("请输入名字")
	}
	if utf8.RuneCountInString(name) > 10 {
		return gameStateError("名字长度不能超过10个字符")
	}

	var err error
	a.Gameinfo, err = core.NewGame(core.GetGameName(), core.UserAgeMax, core.ShowMarketNum, nil, nil)
	if err != nil {
		return gameStateError(err.Error())
	}

	// 设置难度
	difficultyConfig := core.GetDifficultyConfig(difficulty)
	a.Gameinfo.GDifficulty = difficultyConfig.Level
	// 初始化玩家
	a.Userinfo = core.NewUser(name, sex, core.MaxItemNum, core.UserAgeInit, difficultyConfig.InitMoney, core.GetCompanies())
	a.MiniGameSession = nil
	a.pendingDatingInteraction = nil
	a.resetStockClock()
	announce := a.Gameinfo.UpdateAnnounce(core.Announce{})
	a.setAnnounce(announce)
	if _, _, err := a.Userinfo.RefreshAndValidateUserState(a.Gameinfo); err != nil {
		return gameStateError(err.Error())
	}

	return GameStateResponse{
		Code:         200,
		Msg:          "初始化成功",
		Gameinfo:     a.gameSnapshot(),
		Userinfo:     a.userSnapshot(),
		Announce:     &announce,
		Difficulty:   &difficultyConfig,
		StockEpoch:   a.stockEpoch,
		StockVersion: a.stockVersion,
	}
}

func (a *App) NextTime() GameStateResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return gameStateError(responseMessage(errResp))
	}
	a.pendingDatingInteraction = nil
	announce := core.Announce{}

	err := a.Gameinfo.GetRandItem()
	if err != nil {
		return gameStateError(err.Error())
	}

	announce.AnnounceCompany = a.Gameinfo.CalculateCompanyProfit(a.Userinfo)

	bankruptcy := a.Gameinfo.GetRandCompany()
	for _, i := range bankruptcy {
		if a.Userinfo.UCompany[i].UCompanyNum > 0 {
			liquidationValue := a.Gameinfo.GCompanyInfo[i].CIPrice * a.Userinfo.UCompany[i].UCompanyNum
			costPrice := int(math.Round(float64(liquidationValue) * 0.5))
			a.Userinfo.UCash += costPrice
			delete(a.Userinfo.UCompany, i)
			announce.AnnounceCompany = append(announce.AnnounceCompany, "💀 "+a.Gameinfo.GCompanyInfo[i].CIName+"破产了！按照50%清算资产，获得"+strconv.Itoa(costPrice)+"元")
		} else {
			announce.AnnounceCompany = append(announce.AnnounceCompany, "💀 "+a.Gameinfo.GCompanyInfo[i].CIName+"破产了！")
		}
	}

	announce.AnnounceGame = a.settleAllStocks()

	if a.Gameinfo.GAntiqueInfo.AIId != 0 {
		a.Gameinfo.GAntiqueInfo = core.AntiqueInfo{}
		announce.AnnounceGame = append(announce.AnnounceGame, "⏰ 拍卖会时间结束，未完成的竞拍已取消")
	}

	if depositMsg := core.ProcessDepositAnnual(a.Userinfo); depositMsg != "" {
		announce.AnnounceGame = append(announce.AnnounceGame, depositMsg)
	}
	loanMsg, ok := core.ProcessLoanAnnual(a.Userinfo, a.Gameinfo)
	if !ok {
		return gameStateError("💔 贷款逾期超过10次，你破产了！")
	}
	if loanMsg != "" {
		announce.AnnounceGame = append(announce.AnnounceGame, loanMsg)
	}

	a.Userinfo.UAge += 1
	a.Userinfo.UOpportunity.OWNum = 0
	a.Userinfo.UOpportunity.OGNum = 0
	a.Userinfo.UOpportunity.OMNum = 0
	a.Userinfo.UOpportunity.OSNum = 0
	a.Userinfo.UOpportunity.OANum = 0
	a.Userinfo.UStockProfit = 0
	a.Gameinfo.GStockUpdateCount = 0
	// 清空游戏信息
	if a.MiniGameSession != nil {
		announce.AnnounceGame = append(announce.AnnounceGame, "🎮 游戏结束，未完成的游戏已取消")
	}
	a.MiniGameSession = nil
	core.ResetBankTaskStats(a.Gameinfo)

	for i, ucinfo := range a.Userinfo.UCompany {
		ucinfo.UCompanyHoldTime += 1
		a.Userinfo.UCompany[i] = ucinfo
	}

	for i := range a.Userinfo.UAntique {
		if a.Userinfo.UAntique[i].AITime >= 10 {
			continue
		}
		if a.Userinfo.UAntique[i].AIDisplay == 1 || a.Userinfo.UAntique[i].AIDisplay == 2 {
			continue
		}
		a.Userinfo.UAntique[i].AITime++
	}

	for _, house := range core.GetHouses() {
		currentPrice := core.CalculateHousePrice(house)
		a.Gameinfo.GHouseInfo[house.HId] = core.HouseInfo{
			HIId:     house.HId,
			HIName:   house.HName,
			HIPrice:  currentPrice,
			HIHealth: core.CalculateLifestyleHealth(house.HHealth),
			HIFame:   core.CalculateLifestyleFame(house.HFame),
			HIImg:    house.HImg,
		}
	}

	for _, car := range core.GetCars() {
		currentPrice := core.CalculateCarPrice(car)
		a.Gameinfo.GCarInfo[car.CId] = core.CarInfo{
			CIId:     car.CId,
			CIName:   car.CName,
			CIPrice:  currentPrice,
			CIHealth: core.CalculateLifestyleHealth(car.CHealth),
			CIFame:   core.CalculateLifestyleFame(car.CFame),
			CIImg:    car.CImg,
		}
	}

	currentAssets, prevAssets, err := a.Userinfo.RefreshAndValidateUserState(a.Gameinfo)
	if err != nil {
		return gameStateError(err.Error())
	}
	announce.AnnounceGame = append(announce.AnnounceGame, a.Userinfo.ApplyAssetFameMilestones(prevAssets, currentAssets)...)

	change, reason := core.GetImmunityEvent(a.Gameinfo.GDifficulty, a.Userinfo)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + change)
	announce.AnnounceHealthy = append(announce.AnnounceHealthy, reason)
	if _, _, err := a.Userinfo.RefreshAndValidateUserState(a.Gameinfo); err != nil {
		return gameStateError(err.Error())
	}
	criticalGameOver, rescueYearsLeft := core.AdvanceCriticalHealthYear(a.Userinfo)
	if criticalGameOver {
		return gameStateError(fmt.Sprintf("免疫力连续%d年低于%d，抢救期结束，游戏结束！", a.Userinfo.UCriticalHealthYears, core.ImmunityThreshold))
	}
	if a.Userinfo.UCriticalHealthYears > 0 {
		if rescueYearsLeft > 0 {
			announce.AnnounceHealthy = append(announce.AnnounceHealthy,
				fmt.Sprintf("🚨 免疫力已连续%d年低于%d，还有%d个年度可以就医恢复，不会立即死亡", a.Userinfo.UCriticalHealthYears, core.ImmunityThreshold, rescueYearsLeft))
		} else {
			announce.AnnounceHealthy = append(announce.AnnounceHealthy,
				fmt.Sprintf("🚑 免疫力已连续%d年低于%d，今年是最后抢救期，请立即去医院", a.Userinfo.UCriticalHealthYears, core.ImmunityThreshold))
		}
	}

	UDiseasesTmp := a.Userinfo.UDiseases
	for i, disease := range UDiseasesTmp {
		disease.UDTime += 1
		if disease.UUpgradeDays > 0 && disease.UDTime%disease.UUpgradeDays == 0 {
			disease.UDSeverity = core.CalcDisease(disease.UDSeverity + 1)
		}
		a.Userinfo.UDiseases[i] = disease
	}

	diseaseTmp := core.GenerateDisease(a.Userinfo)
	if diseaseTmp.DId != 0 {
		if disease, ok := a.Userinfo.UDiseases[diseaseTmp.DId]; ok {
			disease.UDSeverity = core.CalcDisease(disease.UDSeverity + 1)
			a.Userinfo.UDiseases[diseaseTmp.DId] = disease
			announce.AnnounceHealthy = append(announce.AnnounceHealthy, "🤒 你的"+diseaseTmp.DName+"加重了！")
		} else {
			a.Userinfo.UDiseases[diseaseTmp.DId] = core.UDiseaseInfo{
				UDName:        diseaseTmp.DName,
				UDType:        diseaseTmp.DType,
				USymptoms:     diseaseTmp.DSymptoms,
				UHealthImpact: diseaseTmp.DHealthImpact,
				UUpgradeDays:  diseaseTmp.DUpgradeDays,
				UTreatments:   diseaseTmp.DTreatments,
				// 新病不会直接达到5级；初始严重度由疾病类型决定。
				UDSeverity: core.GetInitialDiseaseSeverity(diseaseTmp.DType),
				UDTime:     0,
			}
			announce.AnnounceHealthy = append(announce.AnnounceHealthy, "🤒 你得了"+diseaseTmp.DName+"，症状为"+diseaseTmp.DSymptoms)
		}
	}

	diseaseTmp2, DescTmp2 := core.GenerateAssetDisease(a.Userinfo, currentAssets, prevAssets)
	if diseaseTmp2.DId != 0 {
		if disease, ok := a.Userinfo.UDiseases[diseaseTmp2.DId]; ok {
			disease.UDSeverity = core.CalcDisease(disease.UDSeverity + 1)
			a.Userinfo.UDiseases[diseaseTmp2.DId] = disease
			announce.AnnounceHealthy = append(announce.AnnounceHealthy, DescTmp2+"你的"+disease.UDName+"加重了！")
		} else {
			a.Userinfo.UDiseases[diseaseTmp2.DId] = core.UDiseaseInfo{
				UDName:        diseaseTmp2.DName,
				UDType:        diseaseTmp2.DType,
				USymptoms:     diseaseTmp2.DSymptoms,
				UHealthImpact: diseaseTmp2.DHealthImpact,
				UUpgradeDays:  diseaseTmp2.DUpgradeDays,
				UTreatments:   diseaseTmp2.DTreatments,
				UDSeverity:    core.GetInitialDiseaseSeverity(diseaseTmp2.DType),
				UDTime:        0,
			}
			announce.AnnounceHealthy = append(announce.AnnounceHealthy, DescTmp2)
		}
	}

	if emergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty); emergency.Required {
		announce.AnnounceHealthy = append(announce.AnnounceHealthy,
			fmt.Sprintf("🏥 当前需要急诊：%s，预计费用%d元", strings.Join(emergency.Reasons, "；"), emergency.Cost))
	}

	announce = a.Gameinfo.UpdateAnnounce(announce)
	a.setAnnounce(announce)

	return GameStateResponse{
		Code:         200,
		Msg:          "新年新气象",
		Gameinfo:     a.gameSnapshot(),
		Userinfo:     a.userSnapshot(),
		Announce:     &announce,
		StockEpoch:   a.stockEpoch,
		StockVersion: a.stockVersion,
	}
}

func (a *App) settleAllStocks() []string {
	var announces []string

	if len(a.Userinfo.UStock) == 0 {
		return announces
	}

	for siid, stock := range a.Userinfo.UStock {
		if stock.USNum <= 0 {
			continue
		}

		var currentPrice int
		var stockName string
		for _, gs := range a.Gameinfo.GStockInfo {
			if gs.SIId == siid {
				currentPrice = gs.SIPrice
				stockName = gs.SIName
				break
			}
		}

		if currentPrice <= 0 {
			continue
		}

		totalValue := currentPrice * stock.USNum
		profit := totalValue - stock.USPrice_init*stock.USNum

		a.Userinfo.UCash += totalValue

		if profit >= 0 {
			announces = append(announces, fmt.Sprintf("📈 股票年末结算：%s %d股，每股%d元，获利%d元", stockName, stock.USNum, currentPrice, profit))
		} else {
			announces = append(announces, fmt.Sprintf("📉 股票年末结算：%s %d股，每股%d元，亏损%d元", stockName, stock.USNum, currentPrice, -profit))
		}
	}

	a.Userinfo.UStock = make(map[int]core.UserStockInfo)

	return announces
}
