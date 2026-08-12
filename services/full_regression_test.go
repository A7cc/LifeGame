package services

import (
	"LifeGame/core"
	"context"
	"math"
	"sort"
	"testing"
	"time"
)

func TestFullBackendFunctionalRegression(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	app := NewApp()
	app.Startup(context.Background())
	t.Cleanup(func() { app.Shutdown(context.Background()) })

	status := app.GetStartupStatus().Status
	if !status.Ready {
		t.Fatalf("application did not start: %#v", status)
	}

	resetGame := func(t *testing.T) {
		t.Helper()
		if code := responseCode(t, app.InitGame("全功能回归", true, core.DifficultyEasy)); code != 200 {
			t.Fatalf("InitGame() code = %d, want 200", code)
		}
		app.Userinfo.UCash = 200_000_000
		app.Userinfo.UFame = core.MaxFame
		app.Userinfo.UImmunity = 80
		app.Gameinfo.GMaxHoldNum.MGRoundNum = 1_000
		app.Gameinfo.GMaxHoldNum.MWRoundNum = 1_000
		app.Gameinfo.GMaxHoldNum.MMRoundNum = 1_000
		app.Gameinfo.GMaxHoldNum.MSRoundNum = 1_000
		app.Gameinfo.GMaxHoldNum.MARoundNum = 1_000
		app.Userinfo.UAssets = core.CalculateUserAssets(app.Userinfo, app.Gameinfo)
	}

	t.Run("seeded startup and game initialization", func(t *testing.T) {
		resetGame(t)
		if len(app.Gameinfo.GItemInsInfo) == 0 || len(app.Gameinfo.GItemOutInfo) == 0 ||
			len(app.Gameinfo.GCompanyInfo) == 0 || len(app.Gameinfo.GStockInfo) == 0 ||
			len(app.Gameinfo.GHouseInfo) == 0 || len(app.Gameinfo.GCarInfo) == 0 ||
			len(app.Gameinfo.GDatingInfo) == 0 {
			t.Fatalf("game initialization missed seeded data: %#v", app.Gameinfo)
		}
		if len(core.GetHospitals()) == 0 || len(core.GetTreats()) == 0 || len(core.GetAntiques()) == 0 {
			t.Fatal("reference data caches were not populated")
		}
		configs := responseMap(t, app.GetMiniGameConfigs())
		if got := len(configs["configs"].([]map[string]interface{})); got != len(core.GetMiniGames()) {
			t.Fatalf("mini-game config count = %d, want %d", got, len(core.GetMiniGames()))
		}
	})

	t.Run("domestic foreign company and stock trading", func(t *testing.T) {
		resetGame(t)
		domesticID := firstDisplayedItemID(t, app.Gameinfo.GItemInsInfo)
		foreignID := firstDisplayedItemID(t, app.Gameinfo.GItemOutInfo)
		if code := responseCode(t, app.BuyItem(domesticID, 1, "国内市场")); code != 200 {
			t.Fatalf("buy domestic code = %d", code)
		}
		if code := responseCode(t, app.SellItem(domesticID, 1, "国内市场")); code != 200 {
			t.Fatalf("sell domestic code = %d", code)
		}
		if code := responseCode(t, app.BuyItem(foreignID, 1, "国外市场")); code != 200 {
			t.Fatalf("buy foreign code = %d", code)
		}
		if code := responseCode(t, app.SellItem(foreignID, 1, "国外市场")); code != 200 {
			t.Fatalf("sell foreign code = %d", code)
		}

		companyID := firstActiveCompanyID(t, app.Gameinfo.GCompanyInfo)
		if code := responseCode(t, app.BuyItem(companyID, 1_000, "创业")); code != 200 {
			t.Fatalf("buy company code = %d", code)
		}
		holding := app.Userinfo.UCompany[companyID]
		holding.UCompanyHoldTime = app.Gameinfo.GCompanyInfo[companyID].CITime
		app.Userinfo.UCompany[companyID] = holding
		if code := responseCode(t, app.SellItem(companyID, 1_000, "创业")); code != 200 {
			t.Fatalf("sell company code = %d", code)
		}

		stockID := firstTradableStockID(t, app.Gameinfo.GStockInfo)
		if code := responseCode(t, app.BuyItem(stockID, 10, "股票")); code != 200 {
			t.Fatalf("buy stock code = %d", code)
		}
		if code := responseCode(t, app.SellItem(stockID, 10, "股票")); code != 200 {
			t.Fatalf("sell stock code = %d", code)
		}
		stockUpdate := app.UpdateStock()
		if stockUpdate.Version != 1 || stockUpdate.Epoch == "" {
			t.Fatalf("unexpected stock update clock: %#v", stockUpdate)
		}
		if code := responseCode(t, app.BuyItem(domesticID, 0, "国内市场")); code == 200 {
			t.Fatal("zero-quantity trade was accepted")
		}
	})

	t.Run("banking and bank tasks", func(t *testing.T) {
		resetGame(t)
		if code := responseCode(t, app.OperationMoney("deposit", 10_000)); code != 200 {
			t.Fatalf("deposit code = %d", code)
		}
		if code := responseCode(t, app.OperationMoney("withdraw", 1_000)); code != 200 {
			t.Fatalf("withdraw code = %d", code)
		}
		if code := responseCode(t, app.ApplyLoan(1_000)); code != 200 {
			t.Fatalf("loan code = %d", code)
		}
		due := 1_000 + int(math.Ceil(1_000*core.LoanInterestRate))
		if code := responseCode(t, app.RepayLoan(due)); code != 200 {
			t.Fatalf("repay code = %d", code)
		}

		tasks := responseMap(t, app.GetBankTaskList())["tasklist"].([]map[string]interface{})
		if len(tasks) == 0 {
			t.Fatal("bank task list is empty")
		}
		stats := &app.Gameinfo.GBankTaskStats
		stats.DepositAmount, stats.DepositCount = math.MaxInt32, math.MaxInt32
		stats.WithdrawAmount, stats.WithdrawCount = math.MaxInt32, math.MaxInt32
		stats.LoanAmount, stats.WorkCount = math.MaxInt32, math.MaxInt32
		taskID := tasks[0]["taskid"].(int)
		if code := responseCode(t, app.ClaimTaskReward(taskID)); code != 200 {
			t.Fatalf("claim bank task code = %d", code)
		}
	})

	t.Run("housing cars entertainment and dating", func(t *testing.T) {
		resetGame(t)
		houseID := cheapestHouseID(t, app.Gameinfo.GHouseInfo)
		if code := responseCode(t, app.BuyHouse(houseID)); code != 200 {
			t.Fatalf("buy house code = %d", code)
		}
		if code := responseCode(t, app.SellHouse(houseID)); code != 200 {
			t.Fatalf("sell house code = %d", code)
		}
		carID := cheapestCarID(t, app.Gameinfo.GCarInfo)
		if code := responseCode(t, app.BuyCar(carID)); code != 200 {
			t.Fatalf("buy car code = %d", code)
		}
		if code := responseCode(t, app.SellCar(carID)); code != 200 {
			t.Fatalf("sell car code = %d", code)
		}

		activities := responseMap(t, app.GetEntertainmentActivities())["activities"].([]EntertainmentActivity)
		if len(activities) == 0 {
			t.Fatal("entertainment activity list is empty")
		}
		if code := responseCode(t, app.DoEntertainment(activities[0].ID)); code != 200 {
			t.Fatalf("entertainment code = %d", code)
		}

		datings := app.GetDatingInfo().DatingList
		if len(datings) == 0 {
			t.Fatal("dating list is empty")
		}
		datingID := datings[0].DId
		if _, exists := app.Userinfo.UDating[datingID]; !exists {
			app.Userinfo.UDating[datingID] = core.UserDatingInfo{DDatingId: datingID, DName: datings[0].DName}
		}
		if code := responseCode(t, app.DoDating(datingID, datings[0].DLocations[0])); code != 200 {
			t.Fatalf("dating code = %d", code)
		}
	})

	t.Run("hospital treatment and emergency", func(t *testing.T) {
		resetGame(t)
		if code := responseCode(t, app.GetHospitalInfo()); code != 200 {
			t.Fatalf("hospital info code = %d", code)
		}
		hospitalType, treatmentID := compatibleTreatment(t)
		app.Userinfo.UImmunity = 79
		if code := responseCode(t, app.BuyTreatment(hospitalType, treatmentID)); code != 200 {
			t.Fatalf("buy treatment code = %d", code)
		}
		serviceHospital, serviceName := firstSpecialTreatment(t)
		if code := responseCode(t, app.SpecialTreatment(serviceName, serviceHospital)); code != 200 {
			t.Fatalf("special treatment code = %d", code)
		}

		app.Userinfo.UImmunity = 0
		app.Userinfo.UDiseases[999] = core.UDiseaseInfo{
			UDName: "回归测试重症", UDType: "red", UDSeverity: 5,
		}
		if code := responseCode(t, app.EmergencyTreatment()); code != 200 {
			t.Fatalf("emergency treatment code = %d", code)
		}
		if core.GetHealthEmergencyStatus(app.Userinfo).Required {
			t.Fatal("emergency state remained active after treatment")
		}
	})

	t.Run("antique auction appraisal and sale", func(t *testing.T) {
		resetGame(t)
		auction := responseMap(t, app.GetAntique(0))
		if responseCode(t, auction) != 200 {
			t.Fatalf("get antique response = %#v", auction)
		}
		antique := auction["currentAntique"].(core.AntiqueInfo)
		if code := responseCode(t, app.AuctionEnd(antique.AIPrice, antique.AIId)); code != 200 {
			t.Fatalf("auction end code = %d", code)
		}
		if code := responseCode(t, app.OperationAntique(antique.AIId, 1)); code != 200 {
			t.Fatalf("antique appraisal code = %d", code)
		}
		if code := responseCode(t, app.OperationAntique(antique.AIId, 3)); code != 200 {
			t.Fatalf("antique sale code = %d", code)
		}
	})

	t.Run("all mini games start cancel and settle", func(t *testing.T) {
		resetGame(t)
		games := core.GetMiniGames()
		if len(games) != 23 {
			t.Fatalf("seeded mini-game count = %d, want 23", len(games))
		}
		for _, game := range games {
			game := game
			t.Run(game.MGName, func(t *testing.T) {
				started := startConfiguredMiniGame(t, app, game)
				sessionID, ok := started["sessionid"].(string)
				if !ok || sessionID == "" {
					t.Fatalf("start response has no session ID: %#v", started)
				}
				if responseMap(t, app.CancelMiniGame(sessionID))["cancelled"] != true {
					t.Fatal("active mini-game was not cancelled")
				}

				startConfiguredMiniGame(t, app, game)
				app.MiniGameSession.MGSStartTime = time.Now().Add(-5 * time.Minute).Unix()
				if game.MGName == "blackjack" && !app.MiniGameSession.MGSResolved {
					if code := responseCode(t, app.MiniGameAction("stand")); code != 200 {
						t.Fatalf("blackjack stand code = %d", code)
					}
				}
				submitted := validMiniGameSubmission(app.MiniGameSession)
				result := responseMap(t, app.EndMiniGame(submitted))
				if _, ok := result["code"].(int); !ok {
					t.Fatalf("settlement has invalid response: %#v", result)
				}
				if app.MiniGameSession != nil {
					t.Fatal("settlement did not clear the mini-game session")
				}
			})
		}
	})

	t.Run("all dating candidates have achievable conditions", func(t *testing.T) {
		resetGame(t)
		user := app.Userinfo
		user.UAge = 100
		user.UCash = 100_000_000
		user.UBank = 100_000_000
		user.UFame = 100
		user.UImmunity = 100
		user.UStockProfit = 100_000_000
		for id := 1; id <= 50; id++ {
			user.UCar[id] = true
			user.UHouse[id] = true
		}
		user.UCompany[1] = core.UCompanyInfo{UCompanyNum: 1}
		user.UItemins[1] = core.UItemInfo{UIINum: 100}
		user.UAntique = []core.AntiqueInfo{
			{AIDisplay: 1, AIMaterial: 4},
			{AIDisplay: 1, AIMaterial: 5},
		}
		user.UDating[999] = core.UserDatingInfo{DCount: 100}
		for _, game := range core.GetMiniGames() {
			user.UMiniGameRecords[game.MGName] = core.MiniGameRecord{MGRType: game.MGType, PlayCount: 100, WinCount: 100}
		}
		for _, activity := range entertainmentActivities {
			user.UMiniGameRecords[activity.ID] = core.MiniGameRecord{MGRType: "activity", PlayCount: 100, WinCount: 100}
		}

		candidates := app.Gameinfo.GDatingInfo
		if len(candidates) != 100 {
			t.Fatalf("dating candidate count = %d, want 100", len(candidates))
		}
		femaleCount, maleCount := 0, 0
		for _, candidate := range candidates {
			if candidate.DSex {
				maleCount++
			} else {
				femaleCount++
			}
			reachable := false
			for age := core.UserAgeInit; age <= core.UserAgeMax; age++ {
				user.UAge = age
				if core.CheckDatingUnlock(user, candidate) {
					reachable = true
					break
				}
			}
			if !reachable {
				t.Errorf("dating candidate %d (%s) still has an unreachable condition: %#v", candidate.DId, candidate.DName, candidate.DMeetConditions)
			}
		}
		if femaleCount != 50 || maleCount != 50 {
			t.Fatalf("dating sides = female %d, male %d; want 50/50", femaleCount, maleCount)
		}
	})

	t.Run("all dating locations have server-side scene rules", func(t *testing.T) {
		resetGame(t)
		locations := make(map[string]struct{})
		for _, dating := range app.Gameinfo.GDatingInfo {
			for _, location := range dating.DLocations {
				locations[location] = struct{}{}
				if _, ok := datingSceneForLocation(location); !ok {
					t.Fatalf("dating location %q has no server-side scene rule", location)
				}
			}
		}
		if len(locations) != 108 {
			t.Fatalf("seeded dating locations = %d, want 108", len(locations))
		}
	})

	t.Run("save list load and delete", func(t *testing.T) {
		resetGame(t)
		savedCash := app.Userinfo.UCash
		save := app.SaveGame("完整回归存档")
		if save.Code != 200 {
			t.Fatalf("save response = %#v", save)
		}
		if save.SaveVersion != currentSaveVersion {
			t.Fatalf("saved version = %d, want %d", save.SaveVersion, currentSaveVersion)
		}
		saveID := int(save.SaveID)
		listed := app.ListSaves().Saves
		if len(listed) != 1 || listed[0].ID != saveID {
			t.Fatalf("save list = %#v, want save %d", listed, saveID)
		}
		if code := responseCode(t, app.OperationMoney("deposit", 1_234)); code != 200 {
			t.Fatalf("mutation before load code = %d", code)
		}
		if response := app.LoadGame(saveID); response.Code != 200 {
			t.Fatalf("load response = %#v", response)
		}
		if app.Userinfo.UCash != savedCash {
			t.Fatalf("loaded cash = %d, want %d", app.Userinfo.UCash, savedCash)
		}
		if response := app.DeleteSave(saveID); response.Code != 200 {
			t.Fatalf("delete save response = %#v", response)
		}
		if saves := app.ListSaves().Saves; len(saves) != 0 {
			t.Fatalf("deleted save still listed: %#v", saves)
		}
	})

	t.Run("year progression and game evaluation", func(t *testing.T) {
		resetGame(t)
		age := app.Userinfo.UAge
		if code := responseCode(t, app.NextTime()); code != 200 {
			t.Fatalf("next year code = %d", code)
		}
		if app.Userinfo.UAge != age+1 {
			t.Fatalf("age after NextTime = %d, want %d", app.Userinfo.UAge, age+1)
		}
		ended := app.EndGame()
		if ended.Code != 200 || ended.Evaluation == nil {
			t.Fatalf("end game response = %#v", ended)
		}
		if app.Gameinfo != nil || app.Userinfo != nil {
			t.Fatal("EndGame() did not clear active state")
		}
	})
}

func firstDisplayedItemID(t *testing.T, items map[int]core.ItemInfo) int {
	t.Helper()
	ids := make([]int, 0, len(items))
	for id, item := range items {
		if item.IIDisplay && item.IIPrice > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		t.Fatal("no displayed market item")
	}
	sort.Ints(ids)
	return ids[0]
}

func firstActiveCompanyID(t *testing.T, companies map[int]core.CompanyInfo) int {
	t.Helper()
	ids := make([]int, 0, len(companies))
	for id, company := range companies {
		if company.CIStatus && company.CIPrice > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		t.Fatal("no active company")
	}
	sort.Ints(ids)
	return ids[0]
}

func firstTradableStockID(t *testing.T, stocks []core.StockInfo) int {
	t.Helper()
	for _, stock := range stocks {
		if stock.SIPrice > 0 && stock.SIStatus != "涨停" && stock.SIStatus != "跌停" {
			return stock.SIId
		}
	}
	t.Fatal("no tradable stock")
	return 0
}

func cheapestHouseID(t *testing.T, houses map[int]core.HouseInfo) int {
	t.Helper()
	id, price := 0, math.MaxInt
	for candidate, house := range houses {
		if house.HIPrice > 0 && house.HIPrice < price {
			id, price = candidate, house.HIPrice
		}
	}
	if id == 0 {
		t.Fatal("no purchasable house")
	}
	return id
}

func cheapestCarID(t *testing.T, cars map[int]core.CarInfo) int {
	t.Helper()
	id, price := 0, math.MaxInt
	for candidate, car := range cars {
		if car.CIPrice > 0 && car.CIPrice < price {
			id, price = candidate, car.CIPrice
		}
	}
	if id == 0 {
		t.Fatal("no purchasable car")
	}
	return id
}

func compatibleTreatment(t *testing.T) (string, int) {
	t.Helper()
	for _, hospital := range core.GetHospitals() {
		for _, treatment := range core.GetTreats() {
			if treatment.TSource == hospital.HType || treatment.TSource == "pharmacy" {
				return hospital.HType, treatment.TId
			}
		}
	}
	t.Fatal("no compatible hospital treatment")
	return "", 0
}

func firstSpecialTreatment(t *testing.T) (string, string) {
	t.Helper()
	for _, hospital := range core.GetHospitals() {
		for _, service := range hospital.HServices {
			switch service.HSName {
			case "打针", "针灸", "手术", "脱胎换骨":
				return hospital.HType, service.HSName
			}
		}
	}
	t.Fatal("no supported special treatment")
	return "", ""
}

func startConfiguredMiniGame(t *testing.T, app *App, game core.MiniGame) M {
	t.Helper()
	level := 0
	if _, ok := game.MGDifficulty[level]; !ok {
		levels := make([]int, 0, len(game.MGDifficulty))
		for candidate := range game.MGDifficulty {
			levels = append(levels, candidate)
		}
		sort.Ints(levels)
		if len(levels) == 0 {
			t.Fatalf("mini-game %s has no difficulty", game.MGName)
		}
		level = levels[0]
	}

	var response H
	switch game.MGName {
	case "rps":
		response = app.StartMiniGameWithOptions(game.MGName, level, 0, "rock", 1)
	case "lottery":
		response = app.StartMiniGameWithOptions(game.MGName, level, 0, "", 1)
	default:
		if rule, ok := miniGameWagerRules[game.MGName]; ok {
			choices := make([]string, 0, len(rule.Choices))
			for choice := range rule.Choices {
				choices = append(choices, choice)
			}
			sort.Strings(choices)
			response = app.StartMiniGameWithOptions(game.MGName, level, rule.Min, choices[0], 1)
		} else {
			response = app.StartMiniGame(game.MGName, level)
		}
	}
	result := responseMap(t, response)
	if responseCode(t, result) != 200 {
		t.Fatalf("start %s response = %#v", game.MGName, result)
	}
	return result
}

func validMiniGameSubmission(session *core.MiniGameSession) int {
	if session.MGSName == "guess" {
		return session.MGSSecret
	}
	difficulty := session.MGSSubInfo
	if difficulty.SMGTarget == 0 {
		return 1
	}
	if difficulty.SMGReward["money"] != 0 {
		return 1
	}
	return difficulty.SMGTarget
}
