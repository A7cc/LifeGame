package core

import "LifeGame/utils"

// 设置一些常量（从配置文件读取，这里作为默认值）
const (
	// 商品上涨浮动，一般是1+20%
	DefaultUpPrice float64 = 1.2
	// 商品下跌浮动，一般是1-20%
	DefaultDownPrice float64 = 0.8
	// 显示物资数量
	DefaultShowItemNum int = 10
	// 设置公司初始值上下浮动
	DefaultCompanyFluct float64 = 0.02
	// 用户可以创业的数量
	DefaultCompanyNum int = 5
	// 用户最大持有古董数
	DefaultUserAntiqueNum int = 10
	// 市场物品显示的数目
	DefaultShowMarketNum int = 10
	// 初始年龄
	DefaultUserAgeInit int = 18
	// 最大年龄
	DefaultUserAgeMax int = 100
	// 贷款利率（年利率10%）
	DefaultLoanInterestRate = 0.10
	// 逾期惩罚（每次逾期增加20%利息）
	DefaultOverduePenalty = 0.20
	// 国内和国外最大物资持有量
	DefaultMaxItemNum int = 100
	// 最大名声
	DefaultMaxFame int = 150
	// 最大免疫力
	DefaultMaxImmunity int = 100
	// 免疫力最低阈值
	DefaultImmunityThreshold int = 10
	// 免疫力低于安全线后允许玩家完整抢救的年度数；再持续一年才会死亡。
	DefaultCriticalHealthGraceYears int = 2
	// 免疫力归零时保留的最低值，避免一次年度结算直接结束游戏。
	DefaultMinimumSurvivableImmunity int = 1
	// 进入紧急治疗后最低可达到的免疫力
	DefaultImmunityEmergencyThreshold int = 60
	// 进入最低免疫时基础费用
	DefaultImmunityEmergencyBaseCost int = 15000
	// 进入严重疾病基础费用
	DefaultDiseaseSeriousBaseCost int = 20000
	// 每次增加1增加费用
	DefaultDiseaseSeriousCostPerStep int = 3000
	// 每个游戏年度最多生成的股票行情次数，避免现实时间无限刷新套利。
	DefaultMaxStockUpdatesPerYear int = 20
	// 商品卖出手续费，抑制零成本反复交易。
	DefaultMarketSellFeeRate float64 = 0.05
	// 股票双边手续费。
	DefaultStockTransactionFeeRate float64 = 0.005
	// 银行存款年利率。
	DefaultDepositInterestRate float64 = 0.02
)

// 缓存的游戏数据（从数据库加载）
var (
	CachedItemIns   []Item
	CachedItemOut   []Item
	CachedCompanies []Company
	CachedAntiques  []Antique
	CachedMiniGames []MiniGame
	CachedStocks    []Stock
	CachedStockNews []string
	CachedHouses    []House
	CachedCars      []Car
	CachedDatings   []DatingInfo
	CachedBankTasks []BankTask
	CachedDiseases  []DiseaseInfo
	CachedTreats    []TreatInfo
	CachedHospitals []HospitalInfo
)

// 运行时配置值，启动时由当前配置文件刷新。
var (
	UpPrice                    = DefaultUpPrice
	DownPrice                  = DefaultDownPrice
	ShowItemNum                = DefaultShowItemNum
	CompanyFluct               = DefaultCompanyFluct
	CompanyNum                 = DefaultCompanyNum
	UserAntiqueNum             = DefaultUserAntiqueNum
	ShowMarketNum              = DefaultShowMarketNum
	UserAgeInit                = DefaultUserAgeInit
	UserAgeMax                 = DefaultUserAgeMax
	LoanInterestRate           = DefaultLoanInterestRate
	OverduePenalty             = DefaultOverduePenalty
	MaxItemNum                 = DefaultMaxItemNum
	MaxFame                    = DefaultMaxFame
	MaxImmunity                = DefaultMaxImmunity
	ImmunityThreshold          = DefaultImmunityThreshold
	ImmunityEmergencyThreshold = DefaultImmunityEmergencyThreshold
	ImmunityEmergencyBaseCost  = DefaultImmunityEmergencyBaseCost
	DiseaseSeriousBaseCost     = DefaultDiseaseSeriousBaseCost
	DiseaseSeriousCostPerStep  = DefaultDiseaseSeriousCostPerStep
	MaxStockUpdatesPerYear     = DefaultMaxStockUpdatesPerYear
	MarketSellFeeRate          = DefaultMarketSellFeeRate
	StockTransactionFeeRate    = DefaultStockTransactionFeeRate
	DepositInterestRate        = DefaultDepositInterestRate
)

// UpdateConfigValues 从配置更新常量值
func UpdateConfigValues() {
	if AppConfig != nil {
		UpPrice = AppConfig.Market.UpPrice
		DownPrice = AppConfig.Market.DownPrice
		ShowItemNum = AppConfig.Market.ShowItemNum
		ShowMarketNum = AppConfig.Market.ShowMarketNum
		MaxItemNum = AppConfig.Market.MaxItemNum
		MaxFame = AppConfig.Market.MaxFame
		MaxImmunity = AppConfig.Market.MaxImmunity
		CompanyFluct = AppConfig.Company.Fluct
		CompanyNum = AppConfig.Company.MaxNum
		UserAntiqueNum = AppConfig.Game.AntiqueMaxNum
		LoanInterestRate = AppConfig.Loan.InterestRate
		OverduePenalty = AppConfig.Loan.OverduePenalty
		UserAgeInit = AppConfig.Game.AgeInit
		UserAgeMax = AppConfig.Game.AgeMax
	}
}

// =========== 市场物资信息及其相关获取方法 ===========
// GetItemIns 获取国内物资列表（从缓存读取）
func GetItemIns() []Item {
	if len(CachedItemIns) > 0 {
		return CachedItemIns
	}
	return []Item{}
}

// GetItemOut 获取国外物资列表（从缓存读取）
func GetItemOut() []Item {
	if len(CachedItemOut) > 0 {
		return CachedItemOut
	}
	return []Item{}
}

// =========== 公司信息及其相关获取方法 ===========

// GetCompanies 获取公司列表（从缓存读取）
func GetCompanies() []Company {
	if len(CachedCompanies) > 0 {
		return CachedCompanies
	}
	return []Company{}
}

// GetCompanyBasePrice 获取公司基础价格（根据公司ID）
func GetCompanyBasePrice(companyId int) int {
	companies := GetCompanies()
	for _, c := range companies {
		if c.CId == companyId {
			return c.CPrice
		}
	}
	return 100 // 默认值
}

// =========== 股票信息及其相关获取方法 ===========

// GetStocks 获取股票列表（从缓存读取）
func GetStocks() []Stock {
	if len(CachedStocks) > 0 {
		return CachedStocks
	}
	return []Stock{}
}

// GetStockNews 获取股票新闻列表（从缓存读取）
func GetStockNews() []string {
	if len(CachedStockNews) > 0 {
		return CachedStockNews
	}
	return []string{}
}

// 股票新闻关键词与影响幅度
var StockNewsKeywordEffect = map[string]float64{
	"上升": utils.RandFloat(1, 3),
	"上涨": utils.RandFloat(2, 4),
	"上扬": utils.RandFloat(2, 4),
	"反弹": utils.RandFloat(1, 2),
	"增长": utils.RandFloat(1, 3),
	"大涨": utils.RandFloat(5, 8),
	"飙升": utils.RandFloat(6, 10),
	"下跌": -utils.RandFloat(2, 4),
	"下降": -utils.RandFloat(1, 3),
	"下挫": -utils.RandFloat(3, 5),
	"回调": -utils.RandFloat(1, 2),
	"暴跌": -utils.RandFloat(6, 10),
}

// =========== 古董信息及其相关获取方法 ===========
// GetAntiques 获取古董列表（从缓存读取）
func GetAntiques() []Antique {
	if len(CachedAntiques) > 0 {
		return CachedAntiques
	}
	return []Antique{}
}

// =========== 小游戏信息及其相关获取方法 ===========
// 获取小游戏信息
func GetMiniGames() []MiniGame {
	if len(CachedMiniGames) > 0 {
		return CachedMiniGames
	}
	return []MiniGame{}
}

// =========== 银行信息及其相关获取方法 ===========
// GetBankTasks 获取银行任务（从缓存读取）
func GetBankTasks() []BankTask {
	if len(CachedBankTasks) > 0 {
		return CachedBankTasks
	}
	return nil
}

// =========== 约会对象信息及其相关获取方法 ===========
// GetDatings 获取约会对象信息（从缓存读取）
func GetDatings() []DatingInfo {
	if len(CachedDatings) > 0 {
		return CachedDatings
	}
	return nil
}

// =========== 房产信息及其相关获取方法 ===========
// GetHouses 获取房产列表（从缓存读取）
func GetHouses() []House {
	if len(CachedHouses) > 0 {
		return CachedHouses
	}
	return []House{}
}

// =========== 车辆信息及其相关获取方法 ===========
// GetCars 获取车辆列表（从缓存读取）
func GetCars() []Car {
	if len(CachedCars) > 0 {
		return CachedCars
	}
	return []Car{}
}

// ========== 疾病信息及其相关获取方法 ===========
// GetDiseases 获取疾病列表（从缓存读取）
func GetDiseases() []DiseaseInfo {
	if len(CachedDiseases) > 0 {
		return CachedDiseases
	}
	return []DiseaseInfo{}
}

// ========== 治疗信息及其相关获取方法 ===========
// GetTreats 获取治疗列表（从缓存读取）
func GetTreats() []TreatInfo {
	if len(CachedTreats) > 0 {
		return CachedTreats
	}
	return []TreatInfo{}
}

// ========== 医院信息及其相关获取方法 ===========
// GetHospitals 获取医院列表（从缓存读取）
func GetHospitals() []HospitalInfo {
	if len(CachedHospitals) > 0 {
		return CachedHospitals
	}
	return []HospitalInfo{}
}
