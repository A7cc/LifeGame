package core

import "testing"

func TestCalculateUserAssetsSubtractsLoan(t *testing.T) {
	game := &Game{
		GItemInsInfo: map[int]ItemInfo{},
		GItemOutInfo: map[int]ItemInfo{},
		GCompanyInfo: map[int]CompanyInfo{},
		GStockInfo:   []StockInfo{},
		GHouseInfo:   map[int]HouseInfo{},
		GCarInfo:     map[int]CarInfo{},
	}
	user := &User{
		UCash:    500,
		UBank:    200,
		ULoan:    200,
		UItemins: map[int]UItemInfo{},
		UItemout: map[int]UItemInfo{},
		UCompany: map[int]UCompanyInfo{},
		UStock:   map[int]UserStockInfo{},
		UHouse:   map[int]bool{},
		UCar:     map[int]bool{},
	}

	if got, want := CalculateUserAssets(user, game), 500; got != want {
		t.Fatalf("CalculateUserAssets() = %d, want %d", got, want)
	}
}

func TestProcessLoanAnnualFailsOnEleventhOverdue(t *testing.T) {
	user := &User{ULoan: 100, ULoanOverdue: 10}

	if _, ok := ProcessLoanAnnual(user, &Game{}); ok {
		t.Fatal("ProcessLoanAnnual() allowed an eleventh overdue")
	}
	if user.ULoanOverdue != 11 {
		t.Fatalf("ULoanOverdue = %d, want 11", user.ULoanOverdue)
	}
	if user.ULoan != 100 {
		t.Fatalf("bankrupt loan was mutated to %d", user.ULoan)
	}
}

func TestCalcImmunityUsesConfiguredMaximum(t *testing.T) {
	previousMax := MaxImmunity
	MaxImmunity = 120
	t.Cleanup(func() { MaxImmunity = previousMax })

	if got := CalcImmunity(150); got != 120 {
		t.Fatalf("CalcImmunity() = %d, want 120", got)
	}
}
