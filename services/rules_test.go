package services

import (
	"LifeGame/core"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func newRulesTestApp() *App {
	game := &core.Game{
		GItemInsInfo: map[int]core.ItemInfo{},
		GItemOutInfo: map[int]core.ItemInfo{},
		GCompanyInfo: map[int]core.CompanyInfo{},
		GStockInfo:   []core.StockInfo{},
		GHouseInfo:   map[int]core.HouseInfo{},
		GCarInfo:     map[int]core.CarInfo{},
		GMaxHoldNum: core.MaxHoldNum{
			MMRoundNum: 1,
			MSRoundNum: 1,
			MGRoundNum: 10,
			MWRoundNum: 10,
		},
	}
	user := core.NewUser("测试", true, 100, 18, 1000, nil)
	return &App{Gameinfo: game, Userinfo: user}
}

func responseCode(t *testing.T, response any) int {
	t.Helper()
	switch result := response.(type) {
	case M:
		code, ok := result["code"].(int)
		if !ok {
			t.Fatalf("response code type = %T, want int", result["code"])
		}
		return code
	case BasicResponse:
		return result.Code
	case StartupResponse:
		return result.Code
	case GameStateResponse:
		return result.Code
	case StockUpdateResponse:
		return result.Code
	case EvaluationResponse:
		return result.Code
	case DatingListResponse:
		return result.Code
	case DatingSceneResponse:
		return result.Code
	case DatingResultResponse:
		return result.Code
	case DatingRelationshipResponse:
		return result.Code
	case DatingInteractionResponse:
		return result.Code
	case SpouseInteractionResponse:
		return result.Code
	case SaveGameResponse:
		return result.Code
	case LoadGameResponse:
		return result.Code
	case ListSavesResponse:
		return result.Code
	default:
		t.Fatalf("response type = %T has no test adapter", response)
		return 0
	}
}

func responseMap(t *testing.T, response H) M {
	t.Helper()
	result, ok := response.(M)
	if !ok {
		t.Fatalf("response type = %T, want services.M", response)
	}
	return result
}

func useMiniGamesForTest(t *testing.T, games []core.MiniGame) {
	t.Helper()
	previous := core.CachedMiniGames
	core.CachedMiniGames = games
	t.Cleanup(func() { core.CachedMiniGames = previous })
}

func TestRepayLoanRequiresPrincipalAndInterest(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UCash = 0
	app.Userinfo.UBank = 110
	app.Userinfo.ULoan = 100

	if code := responseCode(t, app.RepayLoan(100)); code == 200 {
		t.Fatal("RepayLoan() accepted a principal-only payment")
	}
	if app.Userinfo.UBank != 110 || app.Userinfo.ULoan != 100 {
		t.Fatalf("failed repayment mutated balances: bank=%d loan=%d", app.Userinfo.UBank, app.Userinfo.ULoan)
	}

	if code := responseCode(t, app.RepayLoan(110)); code != 200 {
		t.Fatalf("RepayLoan() code = %d, want 200", code)
	}
	if app.Userinfo.UBank != 0 || app.Userinfo.ULoan != 0 {
		t.Fatalf("settled balances: bank=%d loan=%d, want 0/0", app.Userinfo.UBank, app.Userinfo.ULoan)
	}
}

func TestBankTasksUseNetFlowInsteadOfRecycledMoney(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UCash = 1_000

	if code := responseCode(t, app.OperationMoney("deposit", 500)); code != 200 {
		t.Fatalf("deposit code = %d", code)
	}
	if code := responseCode(t, app.OperationMoney("withdraw", 500)); code != 200 {
		t.Fatalf("withdraw code = %d", code)
	}
	if code := responseCode(t, app.OperationMoney("deposit", 500)); code != 200 {
		t.Fatalf("repeated deposit code = %d", code)
	}

	stats := app.Gameinfo.GBankTaskStats
	if stats.DepositAmount != 500 || stats.DepositCount != 1 {
		t.Fatalf("recycled deposits counted as progress: %#v", stats)
	}
	if stats.WithdrawAmount != 0 || stats.WithdrawCount != 0 {
		t.Fatalf("withdrawing a prior deposit counted as net withdrawal: %#v", stats)
	}
}

func TestDoEntertainmentPersistsStateAndEnforcesLimit(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UImmunity = 50

	if code := responseCode(t, app.DoEntertainment("film")); code != 200 {
		t.Fatalf("DoEntertainment() code = %d, want 200", code)
	}
	if app.Userinfo.UCash != 900 {
		t.Fatalf("cash = %d, want 900", app.Userinfo.UCash)
	}
	if app.Userinfo.UOpportunity.OSNum != 1 {
		t.Fatalf("OSNum = %d, want 1", app.Userinfo.UOpportunity.OSNum)
	}
	if app.Userinfo.UImmunity < 52 || app.Userinfo.UImmunity > 54 {
		t.Fatalf("immunity = %d, want a value in [52, 54]", app.Userinfo.UImmunity)
	}
	if app.Userinfo.UAssets != 900 {
		t.Fatalf("assets = %d, want 900", app.Userinfo.UAssets)
	}
	if record := app.Userinfo.UMiniGameRecords["film"]; record.MGRType != "activity" || record.PlayCount != 1 {
		t.Fatalf("entertainment record = %#v, want one activity play", record)
	}

	if code := responseCode(t, app.DoEntertainment("film")); code == 200 {
		t.Fatal("DoEntertainment() ignored the annual limit")
	}
	if app.Userinfo.UCash != 900 {
		t.Fatalf("rejected entertainment mutated cash to %d", app.Userinfo.UCash)
	}
}

func TestDoDatingEnforcesAnnualLimit(t *testing.T) {
	app := newRulesTestApp()
	app.Gameinfo.GDatingInfo = []core.DatingInfo{{DId: 1, DName: "测试对象", DCost: 100, DLocations: []string{"公园"}}}
	app.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "测试对象"}

	if code := responseCode(t, app.DoDating(1, "公园")); code != 200 {
		t.Fatalf("DoDating() code = %d, want 200", code)
	}
	cashAfterFirstDate := app.Userinfo.UCash
	if app.Userinfo.UOpportunity.OMNum != 1 {
		t.Fatalf("OMNum = %d, want 1", app.Userinfo.UOpportunity.OMNum)
	}

	if code := responseCode(t, app.DoDating(1, "公园")); code == 200 {
		t.Fatal("DoDating() ignored the annual limit")
	}
	if app.Userinfo.UCash != cashAfterFirstDate {
		t.Fatalf("rejected date mutated cash to %d", app.Userinfo.UCash)
	}
}

func TestDatingLocationIsValidatedAndReturnsSceneOutcome(t *testing.T) {
	app := newRulesTestApp()
	app.Gameinfo.GDatingInfo = []core.DatingInfo{{
		DId: 1, DName: "测试对象", DCost: 100, DLocations: []string{"公园"},
	}}
	app.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "测试对象"}
	beforeCash := app.Userinfo.UCash

	if result := app.DoDating(1, "不存在的场景"); result.Code == 200 {
		t.Fatal("DoDating() accepted an unknown location")
	}
	if app.Userinfo.UCash != beforeCash || app.Userinfo.UOpportunity.OMNum != 0 {
		t.Fatal("invalid dating location mutated user state")
	}

	result := app.DoDating(1, "赌场")
	if result.Code != 200 || result.Scene == nil {
		t.Fatalf("DoDating() non-preferred scene result = %#v", result)
	}
	if result.Scene.Preferred || result.Scene.SuccessRate > 0.20 || result.Scene.RewardTier != "high-risk" || result.Scene.RewardMultiplier != 2 ||
		!strings.Contains(result.Scene.Event, "不是对方特别喜欢") {
		t.Fatalf("non-preferred scene did not use high-risk rule = %#v", result.Scene)
	}

	preferredApp := newRulesTestApp()
	preferredApp.Gameinfo.GDatingInfo = app.Gameinfo.GDatingInfo
	preferredApp.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "测试对象"}
	preferredResult := preferredApp.DoDating(1, "公园")
	if preferredResult.Code != 200 || preferredResult.Scene == nil {
		t.Fatalf("DoDating() preferred scene result = %#v", preferredResult)
	}
	if preferredResult.Scene.Location != "公园" || preferredResult.Scene.Category != "nature" ||
		preferredResult.Scene.Label != "自然漫游" || !preferredResult.Scene.Preferred ||
		preferredResult.Scene.SuccessRate < 0.82 || preferredResult.Scene.RewardTier != "steady" || preferredResult.Scene.RewardMultiplier != 1 ||
		!strings.Contains(preferredResult.Scene.Event, "选中了对方喜欢") {
		t.Fatalf("preferred scene did not use steady rule = %#v", preferredResult.Scene)
	}
	steady := datingSceneChoiceRuleFor(0.60, true)
	risky := datingSceneChoiceRuleFor(0.60, false)
	if steady.SuccessRate <= risky.SuccessRate || risky.AffinityMin <= steady.AffinityMin+steady.AffinityRange-1 {
		t.Fatalf("dating scene risk/reward is not meaningful: steady=%#v risky=%#v", steady, risky)
	}
}

func TestDatingSceneCategoriesHaveDifferentGameplayEffects(t *testing.T) {
	expected := map[string]struct {
		location string
		fame     int
		health   int
		affinity int
	}{
		"dining":   {location: "餐厅", affinity: 2},
		"culture":  {location: "剧院", fame: 2},
		"nature":   {location: "公园", health: 2},
		"activity": {location: "游乐园", health: 1, affinity: 1},
		"study":    {location: "图书馆", fame: 1, affinity: 1},
		"wellness": {location: "温泉", health: 3},
		"luxury":   {location: "游艇", fame: 1, affinity: 2},
	}

	if len(datingLocationCategories) != 108 {
		t.Fatalf("dating location rules = %d, want 108", len(datingLocationCategories))
	}
	for category, want := range expected {
		rule, ok := datingSceneForLocation(want.location)
		if !ok || rule.Category != category || rule.FameBonus != want.fame ||
			rule.HealthBonus != want.health || rule.AffinityBonus != want.affinity {
			t.Fatalf("scene rule %q = %#v, want %#v", category, rule, want)
		}
	}
}

func TestDatingMomentsUnlockByRelationshipAndAdultStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		adults     bool
		roll       float64
		wantMoment string
		wantChange int
	}{
		{name: "new relationship chats", status: core.DatingStatusFriend, adults: true, roll: 0.50, wantMoment: datingMomentChat, wantChange: 1},
		{name: "new relationship can argue", status: core.DatingStatusFriend, adults: true, roll: 0.95, wantMoment: datingMomentArgument, wantChange: -4},
		{name: "dating stage does not unlock kissing", status: core.DatingStatusDating, adults: true, roll: 0.79, wantMoment: datingMomentChat, wantChange: 1},
		{name: "lover unlocks kissing", status: core.DatingStatusLover, adults: true, roll: 0.60, wantMoment: datingMomentKiss, wantChange: 2},
		{name: "lover can argue", status: core.DatingStatusExclusive, adults: true, roll: 0.90, wantMoment: datingMomentArgument, wantChange: -4},
		{name: "adult sweetheart unlocks intimacy", status: core.DatingStatusSweetheart, adults: true, roll: 0.70, wantMoment: datingMomentIntimacy, wantChange: 3},
		{name: "married relationship unlocks intimacy", status: core.DatingStatusMarried, adults: true, roll: 0.84, wantMoment: datingMomentIntimacy, wantChange: 3},
		{name: "intimacy is blocked unless both are adults", status: core.DatingStatusSweetheart, adults: false, roll: 0.70, wantMoment: datingMomentChat, wantChange: 1},
		{name: "sweetheart can argue", status: core.DatingStatusSweetheart, adults: true, roll: 0.90, wantMoment: datingMomentArgument, wantChange: -4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := datingMomentForRoll(tt.status, tt.adults, tt.roll)
			if got.Kind != tt.wantMoment || got.AffinityChange != tt.wantChange {
				t.Fatalf("datingMomentForRoll(%q, %v, %.2f) = %#v", tt.status, tt.adults, tt.roll, got)
			}
		})
	}
}

func TestDatingInteractionRulesHaveRiskAndReward(t *testing.T) {
	for action, rule := range datingInteractionRules {
		outcome, change := datingInteractionOutcomeForRoll(rule, rule.ArgumentRisk/2)
		if outcome != datingMomentArgument || change != -3 {
			t.Fatalf("%s argument outcome = %q, %d", action, outcome, change)
		}
		outcome, change = datingInteractionOutcomeForRoll(rule, rule.ArgumentRisk)
		if outcome != action || change != rule.AffinityBonus {
			t.Fatalf("%s success outcome = %q, %d", action, outcome, change)
		}
	}
}

func TestDatingInteractionIsRelationshipGatedAndSingleUse(t *testing.T) {
	app := newRulesTestApp()
	app.Gameinfo.GDatingInfo = []core.DatingInfo{{DId: 1, DName: "互动对象", DAge: 25, DSex: false}}
	app.Userinfo.UDating[1] = core.UserDatingInfo{
		DDatingId: 1, DName: "互动对象", DAffinity: 15, DCount: 3, DStatus: core.DatingStatusFriend,
	}
	app.pendingDatingInteraction = &datingInteractionSession{DatingID: 1}

	if result := app.DoDatingInteraction(1, "unknown", ""); result.Code == 200 {
		t.Fatal("unknown interaction was accepted")
	}
	if app.pendingDatingInteraction == nil {
		t.Fatal("invalid interaction consumed the pending session")
	}
	if result := app.DoDatingInteraction(1, datingInteractionKiss, ""); result.Code == 200 {
		t.Fatal("friend relationship bypassed the kiss gate")
	}
	if app.pendingDatingInteraction == nil {
		t.Fatal("locked interaction consumed the pending session")
	}

	relationship := app.Userinfo.UDating[1]
	relationship.DAffinity = 60
	relationship.DCount = 12
	relationship.DStatus = core.DatingStatusLover
	app.Userinfo.UDating[1] = relationship
	result := app.DoDatingInteraction(1, datingInteractionKiss, "")
	if result.Code != 200 || result.Datinginfo == nil || result.Userinfo == nil {
		t.Fatalf("unlocked interaction = %#v", result)
	}
	if app.pendingDatingInteraction != nil {
		t.Fatal("successful interaction did not consume its session")
	}
	if repeated := app.DoDatingInteraction(1, datingInteractionChat, ""); repeated.Code == 200 {
		t.Fatal("the same successful date settled more than one interaction")
	}
}

func TestDatingIntimacyRequiresBothAdults(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UAge = 17
	app.Gameinfo.GDatingInfo = []core.DatingInfo{{DId: 1, DName: "互动对象", DAge: 25, DSex: false}}
	app.Userinfo.UDating[1] = core.UserDatingInfo{
		DDatingId: 1, DName: "互动对象", DAffinity: 95, DCount: 25, DStatus: core.DatingStatusSweetheart,
	}
	app.pendingDatingInteraction = &datingInteractionSession{DatingID: 1}

	if result := app.DoDatingInteraction(1, datingInteractionIntimacy, ""); result.Code == 200 {
		t.Fatal("intimacy accepted an underage player")
	}
	if app.pendingDatingInteraction == nil {
		t.Fatal("adult validation failure consumed the pending session")
	}
	app.Userinfo.UAge = 18
	if result := app.DoDatingInteraction(1, datingInteractionIntimacy, ""); result.Code != 200 {
		t.Fatalf("adult intimacy was rejected: %#v", result)
	}
}

func TestDatingOutfitsAreExplicitAndRelationshipGated(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		outfit     string
		wantOK     bool
		wantImage  string
		playerAge  int
		partnerAge int
	}{
		{name: "friend can choose homewear", status: core.DatingStatusFriend, outfit: "homewear", wantOK: true, wantImage: "/images/datinginfo/dating-partner/female/homewear/07.webp", playerAge: 25, partnerAge: 25},
		{name: "friend cannot choose sleepwear", status: core.DatingStatusFriend, outfit: "sleepwear", playerAge: 25, partnerAge: 25},
		{name: "ambiguous unlocks qipao", status: core.DatingStatusAmbiguous, outfit: "qipao", wantOK: true, wantImage: "/images/datinginfo/dating-partner/female/qipao/07.webp", playerAge: 25, partnerAge: 25},
		{name: "dating unlocks swimwear", status: core.DatingStatusDating, outfit: "swimwear", wantOK: true, wantImage: "/images/datinginfo/dating-partner/female/swimwear/07.webp", playerAge: 25, partnerAge: 25},
		{name: "lover unlocks sleepwear", status: core.DatingStatusLover, outfit: "sleepwear", wantOK: true, wantImage: "/images/datinginfo/dating-partner/female/sleepwear/07.webp", playerAge: 25, partnerAge: 25},
		{name: "exclusive unlocks romantic outfit", status: core.DatingStatusExclusive, outfit: "romantic", wantOK: true, wantImage: "/images/datinginfo/dating-partner/female/romantic/07.webp", playerAge: 25, partnerAge: 25},
		{name: "romantic outfit requires adults", status: core.DatingStatusExclusive, outfit: "romantic", playerAge: 17, partnerAge: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newRulesTestApp()
			app.randomRoll = func() float64 { return 0.5 }
			app.Userinfo.UAge = tt.playerAge
			app.Gameinfo.GDatingInfo = []core.DatingInfo{{
				DId: 1, DName: "换装对象", DAge: tt.partnerAge, DSex: false,
				DImage: "/images/datinginfo/dating-partner/female/07.webp",
			}}
			app.Userinfo.UDating[1] = core.UserDatingInfo{
				DDatingId: 1, DName: "换装对象", DAffinity: 80, DCount: 16, DStatus: tt.status,
			}
			app.pendingDatingInteraction = &datingInteractionSession{DatingID: 1}

			result := app.DoDatingInteraction(1, datingInteractionOutfit, tt.outfit)
			if gotOK := result.Code == 200; gotOK != tt.wantOK {
				t.Fatalf("DoDatingInteraction(%q, %q) = %#v, want success %v", tt.status, tt.outfit, result, tt.wantOK)
			}
			if !tt.wantOK {
				if app.pendingDatingInteraction == nil {
					t.Fatal("locked outfit consumed the pending interaction")
				}
				return
			}
			if result.OutfitVariant != tt.outfit || result.OutfitImage != tt.wantImage {
				t.Fatalf("explicit outfit result = %#v", result)
			}
		})
	}
}

func TestDatingMatchesOppositeSexAndRejectsBypass(t *testing.T) {
	candidates := []core.DatingInfo{
		{DId: 1, DName: "女性对象", DSex: false, DCost: 100, DLocations: []string{"公园"}},
		{DId: 51, DName: "男性对象", DSex: true, DCost: 100, DLocations: []string{"公园"}},
	}

	maleApp := newRulesTestApp()
	maleApp.Gameinfo.GDatingInfo = candidates
	maleList := maleApp.GetDatingInfo().DatingList
	if len(maleList) != 1 || maleList[0].DSex {
		t.Fatalf("male player dating list = %#v, want one female candidate", maleList)
	}
	maleApp.Userinfo.UDating[51] = core.UserDatingInfo{DDatingId: 51, DName: "男性对象"}
	if code := responseCode(t, maleApp.DoDating(51, "公园")); code == 200 {
		t.Fatal("male player dated a male candidate through a direct backend call")
	}

	femaleApp := newRulesTestApp()
	femaleApp.Userinfo.USex = false
	femaleApp.Gameinfo.GDatingInfo = candidates
	femaleApp.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "女性对象"}
	femaleList := femaleApp.GetDatingInfo().DatingList
	if len(femaleList) != 1 || !femaleList[0].DSex {
		t.Fatalf("female player dating list = %#v, want one male candidate", femaleList)
	}
	if _, exists := femaleApp.Userinfo.UDating[1]; exists {
		t.Fatal("female player's incompatible relationship was not removed")
	}
	femaleApp.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "女性对象"}
	if code := responseCode(t, femaleApp.DoDating(1, "公园")); code == 200 {
		t.Fatal("female player dated a female candidate through a direct backend call")
	}
	if code := responseCode(t, femaleApp.DoDating(51, "公园")); code != 200 {
		t.Fatalf("female player could not date the matched male candidate: code = %d", code)
	}
}

func TestDatingSceneRequiresVisitBeforeMeeting(t *testing.T) {
	app := newRulesTestApp()
	app.Gameinfo.GDatingInfo = []core.DatingInfo{{
		DId:        1,
		DName:      "场景对象",
		DSex:       false,
		DCost:      100,
		DMeetScene: "医院",
		DMeetConditions: []core.MeetCondition{
			{CType: "age", CValue: 18},
			{CType: "random", CValue: 100},
		},
	}}

	listed := app.GetDatingInfo()
	datings := listed.DatingList
	if datings[0].DUnlocked {
		t.Fatal("scene candidate unlocked by merely opening the dating page")
	}
	if _, exists := app.Userinfo.UDating[1]; exists {
		t.Fatal("scene candidate was added before visiting the scene")
	}
	if code := responseCode(t, app.VisitDatingScene("公园")); code == 200 {
		t.Fatal("unknown meeting scene was accepted")
	}
	visited := app.VisitDatingScene("医院")
	if responseCode(t, visited) != 200 || len(visited.Met) != 1 {
		t.Fatalf("VisitDatingScene() = %#v", visited)
	}
	if _, exists := app.Userinfo.UDating[1]; !exists {
		t.Fatal("candidate was not added after a successful scene visit")
	}
}

func TestDatingMarriageIsExclusiveAndSupportsBreakupAndDivorce(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UCash = 100_000
	app.Gameinfo.GMaxHoldNum.MMRoundNum = 10
	app.Gameinfo.GDatingInfo = []core.DatingInfo{
		{DId: 1, DName: "对象一", DSex: false, DCost: 100, DLocations: []string{"公园"}},
		{DId: 2, DName: "对象二", DSex: false, DCost: 100, DLocations: []string{"公园"}},
		{DId: 3, DName: "对象三", DSex: false, DCost: 100, DLocations: []string{"公园"}},
	}
	app.Userinfo.UDating[1] = core.UserDatingInfo{DDatingId: 1, DName: "对象一", DAffinity: 90, DCount: 20, DStatus: core.DatingStatusSweetheart}
	app.Userinfo.UDating[2] = core.UserDatingInfo{DDatingId: 2, DName: "对象二", DAffinity: 95, DCount: 25, DStatus: core.DatingStatusSweetheart}
	app.Userinfo.UDating[3] = core.UserDatingInfo{DDatingId: 3, DName: "对象三", DAffinity: 40, DCount: 8, DStatus: core.DatingStatusDating}

	if code := responseCode(t, app.MarryDating(1)); code != 200 {
		t.Fatalf("MarryDating(1) code = %d", code)
	}
	if app.Userinfo.UMarriedDatingID != 1 || app.Userinfo.UDating[1].DStatus != core.DatingStatusMarried {
		t.Fatalf("marriage state = spouse %d, relationship %#v", app.Userinfo.UMarriedDatingID, app.Userinfo.UDating[1])
	}
	if code := responseCode(t, app.MarryDating(2)); code == 200 {
		t.Fatal("second simultaneous marriage succeeded")
	}
	beforeCash := app.Userinfo.UCash
	bathed := app.BatheWithSpouse(1)
	if bathed.Code != 200 || bathed.Interaction != "bath" || bathed.AffinityChange != 2 {
		t.Fatalf("BatheWithSpouse(1) = %#v", bathed)
	}
	if app.Userinfo.UCash != beforeCash-spouseBathCost || app.Userinfo.UOpportunity.OMNum != 1 {
		t.Fatalf("bath state: cash=%d interactions=%d", app.Userinfo.UCash, app.Userinfo.UOpportunity.OMNum)
	}
	if code := responseCode(t, app.BatheWithSpouse(2)); code == 200 {
		t.Fatal("bath interaction with a non-spouse succeeded")
	}
	if code := responseCode(t, app.DoDating(1, "公园")); code != 200 {
		t.Fatalf("married spouse could not continue dating: code = %d", code)
	}
	if app.Userinfo.UDating[1].DStatus != core.DatingStatusMarried {
		t.Fatalf("dating changed married status to %q", app.Userinfo.UDating[1].DStatus)
	}
	if code := responseCode(t, app.BreakUpDating(1)); code == 200 {
		t.Fatal("married relationship was ended without divorce")
	}
	if code := responseCode(t, app.BreakUpDating(3)); code != 200 {
		t.Fatalf("BreakUpDating(3) code = %d", code)
	}
	if code := responseCode(t, app.DoDating(3, "公园")); code == 200 {
		t.Fatal("dating a former partner succeeded")
	}
	if code := responseCode(t, app.DivorceDating(1)); code != 200 {
		t.Fatalf("DivorceDating(1) code = %d", code)
	}
	if app.Userinfo.UMarriedDatingID != 0 || app.Userinfo.UDating[1].DStatus != core.DatingStatusFormer {
		t.Fatalf("divorce state = spouse %d, relationship %#v", app.Userinfo.UMarriedDatingID, app.Userinfo.UDating[1])
	}
	if code := responseCode(t, app.MarryDating(2)); code != 200 {
		t.Fatalf("marriage after divorce code = %d", code)
	}
}

func TestDatingGiftConsumesOneActionAndUsesPreference(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UCash = 10_000
	app.Gameinfo.GMaxHoldNum.MMRoundNum = 3
	app.Gameinfo.GDatingInfo = []core.DatingInfo{
		{DId: 1, DName: "礼物对象", DSex: false, DCost: 600, DGifts: []string{"鲜花", "书"}},
		{DId: 2, DName: "其他对象", DSex: false, DCost: 600, DGifts: []string{"石头"}},
	}
	app.Userinfo.UDating[1] = core.UserDatingInfo{
		DDatingId: 1, DName: "礼物对象", DAffinity: 50, DCount: 5, DStatus: core.DatingStatusDating,
	}

	if code := responseCode(t, app.GiveDatingGift(1, "不喜欢的礼物")); code == 200 {
		t.Fatal("an unlisted gift was accepted")
	}
	if app.Userinfo.UCash != 10_000 || app.Userinfo.UOpportunity.OMNum != 0 {
		t.Fatal("failed gift mutated player state")
	}
	preferredRule, ok := datingGiftRuleFor(app.Gameinfo, "鲜花")
	if !ok {
		t.Fatal("preferred gift did not receive a rule")
	}
	app.randomRoll = func() float64 { return 0 }
	giftResult := app.GiveDatingGift(1, "鲜花")
	if giftResult.Code != 200 {
		t.Fatalf("preferred gift result = %#v", giftResult)
	}
	if giftResult.Event == "" || giftResult.Gift != "鲜花" || giftResult.GiftCost != preferredRule.Cost || giftResult.AffinityChange != preferredRule.PreferredEffect || !giftResult.Preferred || !giftResult.Success || giftResult.Outcome != "favorite" || giftResult.SuccessRate != preferredGiftSuccessRate {
		t.Fatalf("gift event result = %#v", giftResult)
	}
	relationship := app.Userinfo.UDating[1]
	if app.Userinfo.UCash != 10_000-preferredRule.Cost || app.Userinfo.UOpportunity.OMNum != 1 || relationship.DGiftCount != 1 || relationship.DAffinity != 50+preferredRule.PreferredEffect {
		t.Fatalf("gift state = cash %d, actions %d, relationship %#v", app.Userinfo.UCash, app.Userinfo.UOpportunity.OMNum, relationship)
	}

	dislikedRule, ok := datingGiftRuleFor(app.Gameinfo, "石头")
	if !ok {
		t.Fatal("disliked gift did not receive a rule")
	}
	dislikedResult := app.GiveDatingGift(1, "石头")
	if dislikedResult.Code != 200 || dislikedResult.Preferred || !dislikedResult.Success || dislikedResult.Outcome != "risky-success" || dislikedResult.AffinityChange != dislikedRule.RiskyEffect || dislikedResult.SuccessRate != riskyGiftSuccessRate {
		t.Fatalf("disliked gift result = %#v", dislikedResult)
	}
	relationship = app.Userinfo.UDating[1]
	if relationship.DGiftCount != 2 || relationship.DAffinity != 50+preferredRule.PreferredEffect+dislikedRule.RiskyEffect || dislikedRule.RiskyEffect <= preferredRule.PreferredEffect {
		t.Fatalf("disliked gift state = %#v", relationship)
	}

	app.randomRoll = func() float64 { return 1 }
	rejectedResult := app.GiveDatingGift(1, "石头")
	if rejectedResult.Code != 200 || rejectedResult.Success || rejectedResult.Outcome != "rejected" || rejectedResult.AffinityChange != -2 {
		t.Fatalf("rejected gift result = %#v", rejectedResult)
	}
}

func TestDatingGiftOptionsContainOnePreferenceAndDifferentPrices(t *testing.T) {
	app := newRulesTestApp()
	app.Gameinfo.GDatingInfo = []core.DatingInfo{
		{DId: 1, DName: "目标", DSex: false, DGifts: []string{"鲜花", "书"}},
		{DId: 2, DName: "甲", DSex: false, DGifts: []string{"石头", "相机"}},
		{DId: 3, DName: "乙", DSex: false, DGifts: []string{"钢笔", "香水"}},
	}
	options := createDatingGiftOptions(app.Gameinfo, app.Gameinfo.GDatingInfo[0], rand.New(rand.NewSource(7)))
	if len(options) != datingGiftOptionCount {
		t.Fatalf("gift option count = %d", len(options))
	}
	names := map[string]struct{}{}
	prices := map[int]struct{}{}
	effects := map[int]struct{}{}
	preferredCount := 0
	for _, option := range options {
		names[option.Name] = struct{}{}
		prices[option.Cost] = struct{}{}
		effects[option.PreferredEffect] = struct{}{}
		if isPreferredGift(app.Gameinfo.GDatingInfo[0], option.Name) {
			preferredCount++
		}
		if option.RiskyEffect <= option.PreferredEffect {
			t.Fatalf("risky gift reward is not higher: %#v", option)
		}
	}
	if len(names) != 3 || len(prices) != 3 || len(effects) != 3 || preferredCount != 1 {
		t.Fatalf("invalid gift choices: %#v", options)
	}
}

func TestEvaluationUsesAllCategoriesAndStaysWithinOneHundred(t *testing.T) {
	user := core.NewUser("满分人生", true, 100, core.UserAgeMax, 0, nil)
	user.UAssets = 25_000_000
	user.UImmunity = 100
	user.UFame = core.MaxFame
	user.UCompany = map[int]core.UCompanyInfo{1: {}, 2: {}, 3: {}}
	user.UMiniGameRecords = map[string]core.MiniGameRecord{
		"work": {MGRType: "work", PlayCount: 80},
		"game": {MGRType: "casual", WinCount: 60},
	}
	user.UDating[1] = core.UserDatingInfo{DStatus: core.DatingStatusMarried, DCount: 30, DGiftCount: 10}
	user.UAntique = []core.AntiqueInfo{{AIDisplay: 1}, {AIDisplay: 1}, {AIDisplay: 1}, {AIDisplay: 1}}
	user.UHouse = map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	user.UCar = map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	user.UItemins = map[int]core.UItemInfo{1: {UIINum: 1}}

	evaluation := calculateEvaluation(user)
	if evaluation.Score != 100 {
		t.Fatalf("complete-life score = %#v, want 100", evaluation)
	}

	user.UAssets = -1_000_000
	user.UFame = -100
	evaluation = calculateEvaluation(user)
	if evaluation.WealthScore != 0 || evaluation.FameScore != 0 || evaluation.Score < 0 || evaluation.Score > 100 {
		t.Fatalf("clamped evaluation = %#v", evaluation)
	}
}

func TestGetHospitalInfoRequiresGame(t *testing.T) {
	if code := responseCode(t, NewApp().GetHospitalInfo()); code == 200 {
		t.Fatal("GetHospitalInfo() succeeded without an active game")
	}
}

func TestSpecialTreatmentDeletesZeroSeverityDisease(t *testing.T) {
	previousHospitals := core.CachedHospitals
	core.CachedHospitals = []core.HospitalInfo{{
		HType: "test",
		HServices: []core.HospitalServiceInfo{{
			HSName: "打针",
		}},
	}}
	t.Cleanup(func() { core.CachedHospitals = previousHospitals })

	app := newRulesTestApp()
	app.Userinfo.UImmunity = 80
	app.Userinfo.UDiseases[1] = core.UDiseaseInfo{
		UDName:      "测试疾病",
		UTreatments: []string{"打针"},
		UDSeverity:  4,
	}

	if code := responseCode(t, app.SpecialTreatment("打针", "test")); code != 200 {
		t.Fatalf("SpecialTreatment() code = %d, want 200", code)
	}
	if len(app.Userinfo.UDiseases) != 0 {
		t.Fatalf("diseases = %#v, want none", app.Userinfo.UDiseases)
	}
}

func TestWageredMiniGameUsesBackendEconomy(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "roulette",
		MGType: "gambling",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 500, SMGReward: map[string]int{"totalmoney": 3000}},
		},
	}})

	app := newRulesTestApp()
	app.Userinfo.UCash = 10000
	if code := responseCode(t, app.StartMiniGameWithOptions("roulette", 0, 500, "red", 1)); code != 200 {
		t.Fatalf("StartMiniGameWithOptions() code = %d, want 200", code)
	}
	if app.Userinfo.UCash != 9000 {
		t.Fatalf("cash after entry and wager = %d, want 9000", app.Userinfo.UCash)
	}

	expectedOutcome := app.MiniGameSession.MGSOutcome
	expectedPayout := 0
	if expectedOutcome == 1 {
		expectedPayout = 1000
	}
	// 提交与后端相反的结果，验证前端参数不能改变结算。
	result := responseMap(t, app.EndMiniGame(1-expectedOutcome))
	if result["outcome"] != expectedOutcome || result["payout"] != expectedPayout {
		t.Fatalf("backend outcome/payout = %#v/%#v, want %d/%d", result["outcome"], result["payout"], expectedOutcome, expectedPayout)
	}
	if app.Userinfo.UCash != 9000+expectedPayout {
		t.Fatalf("cash after settlement = %d, want %d", app.Userinfo.UCash, 9000+expectedPayout)
	}
}

func TestWagerValidationDoesNotMutateUser(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "roulette",
		MGType: "gambling",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 500, SMGReward: map[string]int{"totalmoney": 3000}},
		},
	}})

	app := newRulesTestApp()
	before := app.Userinfo.UCash
	if code := responseCode(t, app.StartMiniGameWithOptions("roulette", 0, 550, "invalid", 1)); code == 200 {
		t.Fatal("StartMiniGameWithOptions() accepted an invalid wager")
	}
	if app.Userinfo.UCash != before || app.Userinfo.UOpportunity.OGNum != 0 {
		t.Fatalf("invalid wager mutated user: cash=%d games=%d", app.Userinfo.UCash, app.Userinfo.UOpportunity.OGNum)
	}
}

func TestLotteryTicketsAndPayoutComeFromBackend(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "lottery",
		MGType: "gambling",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 100, SMGReward: map[string]int{"totalmoney": 10000}},
		},
	}})

	app := newRulesTestApp()
	app.Userinfo.UCash = 10000
	result := responseMap(t, app.StartMiniGameWithOptions("lottery", 0, 0, "", 5))
	if result["cost"] != 500 || app.Userinfo.UCash != 9500 {
		t.Fatalf("lottery start cost=%#v cash=%d, want 500/9500", result["cost"], app.Userinfo.UCash)
	}
	if app.MiniGameSession == nil || len(app.MiniGameSession.MGSTickets) != 5 {
		t.Fatalf("backend tickets = %#v, want 5 tickets", app.MiniGameSession)
	}

	expectedPayout := app.MiniGameSession.MGSPayout
	settlement := responseMap(t, app.EndMiniGame(999))
	if settlement["payout"] != expectedPayout {
		t.Fatalf("lottery payout = %#v, want %d", settlement["payout"], expectedPayout)
	}
	if app.Userinfo.UCash != 9500+expectedPayout {
		t.Fatalf("cash = %d, want %d", app.Userinfo.UCash, 9500+expectedPayout)
	}
}

func TestBlackjackDoubleDownIsChargedByBackend(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "blackjack",
		MGType: "gambling",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 800, SMGReward: map[string]int{"totalmoney": 4000}},
		},
	}})

	app := newRulesTestApp()
	app.Userinfo.UCash = 10000
	if code := responseCode(t, app.StartMiniGameWithOptions("blackjack", 0, 800, "", 1)); code != 200 {
		t.Fatalf("start blackjack code = %d, want 200", code)
	}
	// 固定后端牌堆，确保加倍后玩家 21 点并获胜。
	app.MiniGameSession.MGSResolved = false
	app.MiniGameSession.MGSPlayerCards = []string{"♠10", "♥9"}
	app.MiniGameSession.MGSDealerCards = []string{"♣10", "♦8"}
	app.MiniGameSession.MGSDeck = []string{"♠2"}
	if code := responseCode(t, app.AddMiniGameWager(800)); code != 200 {
		t.Fatalf("AddMiniGameWager() code = %d, want 200", code)
	}
	if app.Userinfo.UCash != 7600 || app.MiniGameSession.MGSWager != 1600 {
		t.Fatalf("double down state: cash=%d wager=%d", app.Userinfo.UCash, app.MiniGameSession.MGSWager)
	}
	settlement := responseMap(t, app.EndMiniGame(0)) // 伪造“失败”也不能覆盖后端胜局。
	if settlement["outcome"] != 1 || settlement["payout"] != 3200 {
		t.Fatalf("blackjack outcome/payout = %#v/%#v, want 1/3200", settlement["outcome"], settlement["payout"])
	}
	if app.Userinfo.UCash != 10800 {
		t.Fatalf("cash after double-down win = %d, want 10800", app.Userinfo.UCash)
	}
}

func TestCountRewardRejectsClientOverflow(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "taxi",
		MGType: "work",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGTarget: 5, SMGReward: map[string]int{"money": 100}},
		},
	}})

	app := newRulesTestApp()
	before := app.Userinfo.UCash
	if code := responseCode(t, app.StartMiniGame("taxi", 0)); code != 200 {
		t.Fatalf("StartMiniGame() code = %d, want 200", code)
	}
	if code := responseCode(t, app.EndMiniGame(999)); code == 200 {
		t.Fatal("EndMiniGame() accepted a completion count above the configured target")
	}
	if app.Userinfo.UCash != before {
		t.Fatalf("invalid completion count paid money: cash=%d, want %d", app.Userinfo.UCash, before)
	}
}

func TestMiniGameDurationOnlyRejectsPrematurePositiveResults(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "duration-test",
		MGType: "board",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGMinRunTime: 60, SMGReward: map[string]int{"totalmoney": 1}},
		},
	}})

	app := newRulesTestApp()
	if code := responseCode(t, app.StartMiniGame("duration-test", 0)); code != 200 {
		t.Fatalf("first StartMiniGame() code = %d", code)
	}
	if result := responseMap(t, app.EndMiniGame(1)); responseCode(t, result) == 200 {
		t.Fatalf("premature positive result was accepted: %#v", result)
	}
	if app.MiniGameSession != nil {
		t.Fatal("rejected result left its mini-game session active")
	}

	if code := responseCode(t, app.StartMiniGame("duration-test", 0)); code != 200 {
		t.Fatalf("second StartMiniGame() code = %d", code)
	}
	result := responseMap(t, app.EndMiniGame(0))
	if responseCode(t, result) == 200 || result["msg"] != "你输了！" {
		t.Fatalf("immediate loss response = %#v, want a normal loss", result)
	}
	if app.MiniGameSession != nil {
		t.Fatal("immediate loss left its mini-game session active")
	}
}

func TestVariableCasualRewardIsSelectedByBackendTier(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "guess",
		MGType: "casual",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 10, SMGReward: map[string]int{"totalmoney": 500}},
		},
	}})

	app := newRulesTestApp()
	app.Userinfo.UCash = 1000
	if code := responseCode(t, app.StartMiniGame("guess", 0)); code != 200 {
		t.Fatalf("StartMiniGame() code = %d, want 200", code)
	}
	app.MiniGameSession.MGSSecret = 50
	result := responseMap(t, app.EndMiniGame(45))
	if result["payout"] != 30 {
		t.Fatalf("near-guess payout = %#v, want backend tier 30", result["payout"])
	}
	if app.Userinfo.UCash != 1020 {
		t.Fatalf("cash = %d, want 1020", app.Userinfo.UCash)
	}
}

func TestStateMutationsAreSerializedAndResponsesAreSnapshots(t *testing.T) {
	app := newRulesTestApp()
	app.Userinfo.UCash = 1000

	first := responseMap(t, app.OperationMoney("deposit", 100))
	snapshot, ok := first["userinfo"].(*core.User)
	if !ok {
		t.Fatalf("userinfo snapshot type = %T", first["userinfo"])
	}

	var wg sync.WaitGroup
	for i := 0; i < 90; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.OperationMoney("deposit", 10)
		}()
	}
	wg.Wait()

	if app.Userinfo.UCash != 0 || app.Userinfo.UBank != 1000 {
		t.Fatalf("serialized balances = cash %d, bank %d; want 0/1000", app.Userinfo.UCash, app.Userinfo.UBank)
	}
	if snapshot.UCash != 900 || snapshot.UBank != 100 {
		t.Fatalf("returned response was mutated later: cash %d, bank %d", snapshot.UCash, snapshot.UBank)
	}
}

func TestStockUpdatesReturnVersionedIncrementOnly(t *testing.T) {
	app := newRulesTestApp()
	app.stockEpoch = "test-epoch"

	first := app.UpdateStock()
	if first.Epoch != "test-epoch" || first.Version != 1 {
		t.Fatalf("first stock clock = %#v/%#v, want test-epoch/1", first.Epoch, first.Version)
	}
	if first.Stocks == nil {
		t.Fatal("stocks should be an initialized slice")
	}

	second := app.UpdateStock()
	if second.Version != 1 {
		t.Fatalf("throttled stock version = %#v, want 1", second.Version)
	}
	app.stockUpdatedAt = time.Now().Add(-stockUpdateMinInterval).UnixMilli()
	second = app.UpdateStock()
	if second.Version != 2 {
		t.Fatalf("second stock version = %#v, want 2", second.Version)
	}

	app.Gameinfo.GStockUpdateCount = core.MaxStockUpdatesPerYear - 1
	app.stockUpdatedAt = time.Now().Add(-stockUpdateMinInterval).UnixMilli()
	closed := app.UpdateStock()
	if !closed.MarketClosed || closed.Remaining != 0 {
		t.Fatalf("closing update = %#v", closed)
	}
	versionAtClose := closed.Version
	closed = app.UpdateStock()
	if !closed.MarketClosed || closed.Version != versionAtClose {
		t.Fatalf("closed market advanced again: %#v", closed)
	}
}

func TestCancelMiniGameUsesSessionIDAndDoesNotRefund(t *testing.T) {
	useMiniGamesForTest(t, []core.MiniGame{{
		MGName: "test-cancel",
		MGType: "casual",
		MGDifficulty: map[int]core.SubMiniGame{
			0: {SMGNeed: 100, SMGReward: map[string]int{"totalmoney": 500}},
		},
	}})

	app := newRulesTestApp()
	started := responseMap(t, app.StartMiniGame("test-cancel", 0))
	sessionID, ok := started["sessionid"].(string)
	if !ok || sessionID == "" {
		t.Fatalf("sessionid = %#v, want a non-empty string", started["sessionid"])
	}
	if app.Userinfo.UCash != 900 {
		t.Fatalf("cash after start = %d, want 900", app.Userinfo.UCash)
	}

	wrong := responseMap(t, app.CancelMiniGame("another-session"))
	if wrong["cancelled"] != false || app.MiniGameSession == nil {
		t.Fatal("a mismatched cancellation removed the active session")
	}
	if responseMap(t, app.CancelMiniGame(sessionID))["cancelled"] != true {
		t.Fatal("matching cancellation did not remove the active session")
	}
	if app.MiniGameSession != nil || app.Userinfo.UCash != 900 {
		t.Fatalf("cancelled state: session=%#v cash=%d, want nil/900", app.MiniGameSession, app.Userinfo.UCash)
	}

	newStart := responseMap(t, app.StartMiniGame("test-cancel", 0))
	newSessionID := newStart["sessionid"].(string)
	responseMap(t, app.CancelMiniGame(sessionID))
	if app.MiniGameSession == nil || app.MiniGameSession.MGSID != newSessionID {
		t.Fatal("a delayed cancellation from the previous game removed the new session")
	}
}

func TestDatingOutfitAssetUsesVariantDirectory(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		outfit   string
		expected string
	}{
		{name: "sleepwear", image: "/images/datinginfo/dating-partner/female/07.webp", outfit: "sleepwear", expected: "/images/datinginfo/dating-partner/female/sleepwear/07.webp"},
		{name: "romantic", image: "/images/datinginfo/dating-partner/female/07.webp", outfit: "romantic", expected: "/images/datinginfo/dating-partner/female/romantic/07.webp"},
		{name: "swimwear", image: "/images/datinginfo/dating-partner/male/42.webp", outfit: "swimwear", expected: "/images/datinginfo/dating-partner/male/swimwear/42.webp"},
		{name: "cosplay", image: "/images/datinginfo/dating-partner/male/42.webp", outfit: "cosplay", expected: "/images/datinginfo/dating-partner/male/cosplay/42.webp"},
		{name: "qipao", image: "/images/datinginfo/dating-partner/male/42.webp", outfit: "qipao", expected: "/images/datinginfo/dating-partner/male/qipao/42.webp"},
		{name: "homewear", image: "/images/datinginfo/dating-partner/male/42.webp", outfit: "homewear", expected: "/images/datinginfo/dating-partner/male/homewear/42.webp"},
		{name: "career", image: "/images/datinginfo/dating-partner/male/42.webp", outfit: "career", expected: "/images/datinginfo/dating-partner/male/42.webp"},
		{name: "unknown outfit", image: "/images/datinginfo/dating-partner/female/07.webp", outfit: "formal", expected: "/images/datinginfo/dating-partner/female/07.webp"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := datingOutfitAsset(test.image, test.outfit); actual != test.expected {
				t.Fatalf("datingOutfitAsset(%q, %q) = %q, want %q", test.image, test.outfit, actual, test.expected)
			}
		})
	}
	if len(datingOutfitVariants) != 7 {
		t.Fatalf("dating outfit variants = %d, want 7", len(datingOutfitVariants))
	}
}
