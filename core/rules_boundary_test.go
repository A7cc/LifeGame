package core

import (
	"math"
	"strings"
	"testing"
)

func TestScalarRulesClampBoundaries(t *testing.T) {
	previousMaxFame, previousMaxImmunity := MaxFame, MaxImmunity
	MaxFame, MaxImmunity = 150, 100
	t.Cleanup(func() {
		MaxFame, MaxImmunity = previousMaxFame, previousMaxImmunity
	})

	for _, test := range []struct {
		name string
		got  int
		want int
	}{
		{name: "immunity lower", got: CalcImmunity(-1), want: 0},
		{name: "immunity upper", got: CalcImmunity(101), want: 100},
		{name: "fame lower", got: CalcFame(-1), want: 0},
		{name: "fame upper", got: CalcFame(151), want: 150},
		{name: "disease lower", got: CalcDisease(0), want: 1},
		{name: "disease upper", got: CalcDisease(9), want: 5},
	} {
		if test.got != test.want {
			t.Errorf("%s = %d, want %d", test.name, test.got, test.want)
		}
	}
}

func TestReputationAndDifficultyBoundaries(t *testing.T) {
	levels := []struct {
		reputation int
		level      int
		name       string
		stockLimit int
	}{
		{-1, -1, "老赖", 0},
		{0, 0, "普通", 100_000},
		{70, 0, "普通", 100_000},
		{71, 1, "中等", 500_000},
		{111, 2, "高级", 1_000_000},
		{131, 3, "豪华", 5_000_000},
		{146, 4, "私人", 10_000_000},
		{151, -1, "老赖", 0},
	}
	for _, test := range levels {
		level, name, limit, _ := CalcReputationLevel(test.reputation)
		if level != test.level || name != test.name || limit != test.stockLimit {
			t.Errorf("CalcReputationLevel(%d) = %d/%s/%d", test.reputation, level, name, limit)
		}
	}

	if got := GetDifficultyConfig(999); got.Level != DifficultyNormal {
		t.Fatalf("unknown difficulty level = %d, want normal", got.Level)
	}
}

func TestDatingConditionsAndRelationshipBoundaries(t *testing.T) {
	user := NewUser("规则测试", true, 100, 35, 20_000, nil)
	user.UFame = 80
	user.UBank = 10_000
	user.UImmunity = 75
	user.UStockProfit = 5_000
	user.UHouse[1] = true
	user.UCar[1] = true
	user.UCompany[1] = UCompanyInfo{UCompanyNum: 1}
	user.UItemins[7] = UItemInfo{UIINum: 2}
	user.UItemout[7] = UItemInfo{UIINum: 3}
	user.UAntique = []AntiqueInfo{{AIDisplay: 1, AIMaterial: 4}}
	user.UDating[1] = UserDatingInfo{DCount: 6}
	user.UMiniGameRecords["chess"] = MiniGameRecord{MGRType: "board", PlayCount: 4, WinCount: 2}
	user.UMiniGameRecords["bank"] = MiniGameRecord{MGRType: "work", PlayCount: 3}
	user.UMiniGameRecords["lottery"] = MiniGameRecord{MGRType: "gambling", WinCount: 1}

	conditions := []MeetCondition{
		{CType: "fame", CValue: 80},
		{CType: "cash", CValue: 20_000},
		{CType: "bank", CValue: 10_000},
		{CType: "house", CValue: 1},
		{CType: "car", CValue: 1},
		{CType: "play_game", CTarget: "chess", CValue: 4},
		{CType: "win_game", CTarget: "chess", CValue: 2},
		{CType: "play_game", CValue: 7},
		{CType: "win_game", CValue: 3},
		{CType: "work_count", CValue: 3},
		{CType: "age", CValue: 35},
		{CType: "date_count", CValue: 6},
		{CType: "immunity", CValue: 75},
		{CType: "antique_rare", CValue: 1},
		{CType: "lottery_win", CValue: 1},
		{CType: "item_own", CTarget: "7", CValue: 5},
		{CType: "stock_profit", CValue: 5_000},
		{CType: "company_founder", CValue: 1},
		{CType: "random", CValue: 100},
	}
	for _, condition := range conditions {
		if !CheckCondition(user, condition) {
			t.Errorf("condition should pass: %#v", condition)
		}
	}
	if CheckCondition(user, MeetCondition{CType: "unknown", CValue: 0}) {
		t.Fatal("unknown condition passed")
	}
	if CheckCondition(user, MeetCondition{CType: "random", CValue: 0}) {
		t.Fatal("zero-percent random condition passed")
	}
	if !CheckDatingUnlock(user, DatingInfo{DMeetConditions: conditions}) {
		t.Fatal("candidate with satisfied conditions remained locked")
	}

	statuses := []struct {
		affinity int
		count    int
		already  bool
		want     string
	}{
		{90, 20, false, "爱人"},
		{90, 20, true, "专属恋人"},
		{50, 10, false, "恋人"},
		{30, 0, false, "交往中"},
		{20, 5, false, "暧昧中"},
		{10, 0, false, "朋友"},
		{-1, 10, false, "前任"},
		{0, 0, false, "陌生人"},
	}
	for _, test := range statuses {
		if got := GetDatingStatus(test.affinity, test.count, test.already); got != test.want {
			t.Errorf("GetDatingStatus(%d,%d,%v) = %q, want %q", test.affinity, test.count, test.already, got, test.want)
		}
	}

	if got := CalculateDatingSuccessRate(DatingInfo{DDesc: "善解人意"}, UserDatingInfo{DAffinity: 1_000}); got != 0.95 {
		t.Errorf("high dating rate = %.2f, want 0.95", got)
	}
	if got := CalculateDatingSuccessRate(DatingInfo{DDesc: "神秘"}, UserDatingInfo{DAffinity: -1_000}); got != 0.05 {
		t.Errorf("low dating rate = %.2f, want 0.05", got)
	}
}

func TestEmergencyTreatmentCostsAndCuresOnlySevereDisease(t *testing.T) {
	user := &User{
		UImmunity: 8,
		UDiseases: map[int]UDiseaseInfo{
			1: {UDName: "轻症", UDType: "green", UDSeverity: 2},
			2: {UDName: "危重症", UDType: "red", UDSeverity: 5},
		},
	}
	status := GetHealthEmergencyStatus(user)
	wantCost := ImmunityEmergencyBaseCost + 2*DiseaseSeriousCostPerStep + DiseaseSeriousBaseCost + 3*DiseaseSeriousCostPerStep
	if !status.Required || status.Cost != wantCost || len(status.SevereDiseases) != 1 {
		t.Fatalf("emergency status = %#v, want cost %d and one severe disease", status, wantCost)
	}

	cured := ApplyHealthEmergencyTreatment(user)
	if len(cured) != 1 || cured[0].UDName != "危重症" {
		t.Fatalf("cured diseases = %#v", cured)
	}
	if _, exists := user.UDiseases[1]; !exists {
		t.Fatal("treatment removed a non-severe disease")
	}
	if user.UImmunity != DefaultImmunityEmergencyThreshold {
		t.Fatalf("immunity after emergency = %d, want %d", user.UImmunity, DefaultImmunityEmergencyThreshold)
	}
}

func TestDiseaseDamageHasAnnualCap(t *testing.T) {
	user := &User{UDiseases: map[int]UDiseaseInfo{
		1: {UHealthImpact: -25, UDSeverity: 5},
		2: {UHealthImpact: -20, UDSeverity: 5},
	}}
	for _, difficulty := range []int{DifficultyEasy, DifficultyNormal, DifficultyHard} {
		maxLoss := GetMaxAnnualImmunityLoss(difficulty)
		for i := 0; i < 100; i++ {
			change, _ := GetImmunityEvent(difficulty, user)
			if change < -maxLoss {
				t.Fatalf("difficulty %d immunity change = %d, cap = -%d", difficulty, change, maxLoss)
			}
		}
	}
}

func TestCriticalHealthProvidesTwoFullRescueYears(t *testing.T) {
	user := &User{UImmunity: 0}
	gameOver, left := AdvanceCriticalHealthYear(user)
	if gameOver || left != 1 || user.UImmunity != DefaultMinimumSurvivableImmunity {
		t.Fatalf("first critical year = over %v, left %d, immunity %d", gameOver, left, user.UImmunity)
	}
	gameOver, left = AdvanceCriticalHealthYear(user)
	if gameOver || left != 0 {
		t.Fatalf("second critical year = over %v, left %d", gameOver, left)
	}
	gameOver, _ = AdvanceCriticalHealthYear(user)
	if !gameOver {
		t.Fatal("third consecutive critical year did not end the game")
	}

	user.UImmunity = ImmunityThreshold
	ResetCriticalHealthIfRecovered(user)
	if user.UCriticalHealthYears != 0 {
		t.Fatalf("recovered critical years = %d, want 0", user.UCriticalHealthYears)
	}
}

func TestZeroImmunityEntersRecoverableStateInsteadOfImmediateDeath(t *testing.T) {
	user := NewUser("健康测试", true, MaxItemNum, UserAgeInit, 1000, nil)
	user.UImmunity = 0
	game := &Game{
		GTime:        UserAgeMax,
		GItemInsInfo: map[int]ItemInfo{}, GItemOutInfo: map[int]ItemInfo{},
		GCompanyInfo: map[int]CompanyInfo{}, GHouseInfo: map[int]HouseInfo{}, GCarInfo: map[int]CarInfo{},
	}
	if _, _, err := user.RefreshAndValidateUserState(game); err != nil {
		t.Fatalf("zero immunity ended the game immediately: %v", err)
	}
	if user.UImmunity != DefaultMinimumSurvivableImmunity {
		t.Fatalf("survivable immunity = %d, want %d", user.UImmunity, DefaultMinimumSurvivableImmunity)
	}
}

func TestDiseaseGenerationIsLessFrequentAndNeverStartsCritical(t *testing.T) {
	youngHealthy := &User{UAge: 18, UImmunity: 100}
	olderUnhealthy := &User{UAge: 90, UImmunity: 0}
	if chance := CalculateDiseaseChance(youngHealthy); chance != 15 {
		t.Fatalf("young healthy disease chance = %d, want 15", chance)
	}
	if chance := CalculateDiseaseChance(olderUnhealthy); chance <= CalculateDiseaseChance(youngHealthy) || chance > 55 {
		t.Fatalf("older unhealthy disease chance = %d", chance)
	}
	for _, diseaseType := range []string{"green", "cyan", "yellow", "red"} {
		for i := 0; i < 100; i++ {
			if severity := GetInitialDiseaseSeverity(diseaseType); severity < 1 || severity > 3 {
				t.Fatalf("%s initial severity = %d", diseaseType, severity)
			}
		}
	}

	previous := CachedDiseases
	CachedDiseases = []DiseaseInfo{{DId: 1, DName: "已有疾病", DType: "green"}}
	t.Cleanup(func() { CachedDiseases = previous })
	olderUnhealthy.UDiseases = map[int]UDiseaseInfo{1: {UDName: "已有疾病"}}
	for i := 0; i < 100; i++ {
		if disease := GenerateDisease(olderUnhealthy); disease.DId != 0 {
			t.Fatalf("active disease was generated again: %#v", disease)
		}
	}
}

func TestMarketPricesStayWithinConfiguredBounds(t *testing.T) {
	previousUp, previousDown, previousShow := UpPrice, DownPrice, ShowMarketNum
	UpPrice, DownPrice, ShowMarketNum = 1.2, 0.8, 2
	t.Cleanup(func() { UpPrice, DownPrice, ShowMarketNum = previousUp, previousDown, previousShow })

	item := Item{
		IId: 1, IName: "测试商品", IPrice: 100, IPrice_min: 50, IPrice_max: 150,
		IEffects: map[int]string{0: "⚖️ 平稳", 1: "📈 上涨", 2: "📉 下跌"},
	}
	for i := 0; i < 200; i++ {
		price, _ := getNowPrice(item, true, i%2 == 0, true)
		if price < item.IPrice_min || price > item.IPrice_max {
			t.Fatalf("generated price %d outside [%d,%d]", price, item.IPrice_min, item.IPrice_max)
		}
	}
	if price := CalculateHousePrice(House{HPrice: 100, HPrice_min: 80, HPrice_max: 120}); price < 80 || price > 120 {
		t.Fatalf("house price = %d", price)
	}
	if price := CalculateCarPrice(Car{CPrice: 100, CPrice_min: 75, CPrice_max: 125}); price < 75 || price > 125 {
		t.Fatalf("car price = %d", price)
	}
	if _, err := refreshItem(nil, nil, false); err == nil {
		t.Fatal("empty market refresh succeeded")
	}
	items := []Item{item, {IId: 2, IName: "B", IPrice: 100, IPrice_min: 50, IPrice_max: 150, IEffects: item.IEffects}, {IId: 3, IName: "C", IPrice: 100, IPrice_min: 50, IPrice_max: 150, IEffects: item.IEffects}}
	refreshed, err := refreshItem(map[int]ItemInfo{}, items, false)
	if err != nil {
		t.Fatal(err)
	}
	displayed := 0
	for _, info := range refreshed {
		if info.IIDisplay {
			displayed++
		}
	}
	if displayed != ShowMarketNum {
		t.Fatalf("displayed items = %d, want %d", displayed, ShowMarketNum)
	}
}

func TestTradingRulesPreserveBalancesAndHoldings(t *testing.T) {
	game := &Game{
		GItemInsInfo: map[int]ItemInfo{1: {IIName: "商品", IIPrice: 100, IIDisplay: true}},
		GCompanyInfo: map[int]CompanyInfo{1: {CIName: "公司", CIPrice: 2, CITime: 2, CIStatus: true}},
		GStockInfo:   []StockInfo{{SIId: 1, SIName: "股票", SIPrice: 10}},
		GMaxHoldNum:  MaxHoldNum{MDHoldNum: 5, MCHoldNum: 2},
	}
	user := NewUser("交易测试", true, 5, 18, 10_000, nil)

	if err := user.BuyDomesticGoods(game, 1, 2); err != nil {
		t.Fatal(err)
	}
	if user.UCash != 9_800 || user.UItemins[1].UIINum != 2 {
		t.Fatalf("domestic buy state = cash %d, holding %#v", user.UCash, user.UItemins[1])
	}
	if _, err := user.SellDomesticGoods(game, 1, 3); err == nil {
		t.Fatal("oversell succeeded")
	}
	if proceeds, err := user.SellDomesticGoods(game, 1, 2); err != nil || proceeds != 190 || user.UCash != 9_990 {
		t.Fatalf("domestic sell = proceeds %d, cash %d, err %v", proceeds, user.UCash, err)
	}

	if err := user.BuyCompany(game, 1, 999); err == nil {
		t.Fatal("company purchase below initial minimum succeeded")
	}
	if err := user.BuyCompany(game, 1, 1_000); err != nil {
		t.Fatal(err)
	}
	if _, err := user.SellCompany(game, 1, 1_000); err == nil || !strings.Contains(err.Error(), "融资") {
		t.Fatalf("company sold before holding period: %v", err)
	}
	holding := user.UCompany[1]
	holding.UCompanyHoldTime = 2
	user.UCompany[1] = holding
	if _, err := user.SellCompany(game, 1, 1_000); err != nil {
		t.Fatal(err)
	}

	if err := user.BuyStock(game, 1, 9); err == nil {
		t.Fatal("stock purchase below lot size succeeded")
	}
	if err := user.BuyStock(game, 1, 10); err != nil {
		t.Fatal(err)
	}
	if proceeds, err := user.SellStock(game, 1, 10); err != nil || proceeds != 99 {
		t.Fatalf("stock sale = proceeds %d, err %v", proceeds, err)
	}
}

func TestAssetMilestonesApplyOnlyWhenCrossed(t *testing.T) {
	user := &User{UFame: 0}
	announcements := user.ApplyAssetFameMilestones(99_999, 10_000_000)
	if user.UFame != 6 || len(announcements) != 3 {
		t.Fatalf("milestones fame=%d announcements=%#v", user.UFame, announcements)
	}
	if repeated := user.ApplyAssetFameMilestones(10_000_000, math.MaxInt); len(repeated) != 1 || user.UFame != 11 {
		t.Fatalf("later milestone fame=%d announcements=%#v", user.UFame, repeated)
	}
}

func TestBalancedAnnualEconomyRules(t *testing.T) {
	user := &User{UBank: 10_000}
	message := ProcessDepositAnnual(user)
	if user.UBank != 10_200 || message == "" {
		t.Fatalf("deposit interest = bank %d, message %q; want 10200 and an announcement", user.UBank, message)
	}
	if fame := CalculateLifestyleFame(10); fame != 4 {
		t.Fatalf("scaled lifestyle fame = %d, want 4", fame)
	}
	if health := CalculateLifestyleHealth(10); health != 6 {
		t.Fatalf("scaled lifestyle health = %d, want 6", health)
	}

	for i := 0; i < 100; i++ {
		price := CalculateMaxPrice(100_000, 5, 5, 5, 5)
		if price < 80_000 || price > 250_000 {
			t.Fatalf("balanced antique resale price = %d, outside expected range", price)
		}
	}
}
