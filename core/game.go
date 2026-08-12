package core

import (
	"LifeGame/utils"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// 用于存放游戏的数据结构
// 游戏信息
type Game struct {
	// 游戏id
	GId int `json:"gid"`
	// 游戏名字
	GName string `json:"gname"`
	// 游戏时间，用于设置用户最大岁数
	GTime int `json:"gtime"`
	// 当前游戏难度
	GDifficulty int `json:"gdifficulty"`
	// 国内全部物资
	GItemInsInfo map[int]ItemInfo `json:"giteminsinfo"`
	// 国外全部物资
	GItemOutInfo map[int]ItemInfo `json:"gitemoutinfo"`
	// 公司信息
	GCompanyInfo map[int]CompanyInfo `json:"gcompanyinfo"`
	// 古董信息
	GAntiqueInfo AntiqueInfo `json:"gantiqueinfo"`
	// 游戏可以持有物品最大数量
	GMaxHoldNum MaxHoldNum `json:"gmaxholdnum"`
	// 股票数据信息
	GStockInfo []StockInfo `json:"gstockinfo"`
	// 股票新闻
	GStockNews []string `json:"gstocknews"`
	// 当前游戏年度已经生成的股票行情次数。
	GStockUpdateCount int `json:"gstockupdatecount"`
	// 银行任务统计
	GBankTaskStats BankTaskStats `json:"gbanktaskstats"`
	// 约会对象信息
	GDatingInfo []DatingInfo `json:"gdatinginfo"`
	// 房子信息
	GHouseInfo map[int]HouseInfo `json:"ghouseinfo"`
	// 车子信息
	GCarInfo map[int]CarInfo `json:"gcarinfo"`
}

// 初始化Game
// gname: 游戏名字，gtime: 游戏时间，ShowItemNumTemp：显示物资数量，ItemIns: 国内物资列表，ItemOut: 国外物资列表
func NewGame(gname string, gtime, ShowItemNumTemp int, ItemIns, ItemOut []Item) (*Game, error) {
	// 如果传入的物资列表为空，从缓存获取
	if len(ItemIns) == 0 {
		ItemIns = GetItemIns()
	}
	if len(ItemOut) == 0 {
		ItemOut = GetItemOut()
	}

	game := &Game{
		// 游戏id
		GId: 1,
		// 游戏名字
		GName: gname,
		// 游戏时间，用于设置用户最大岁数
		GTime: gtime,
		// 当前游戏难度
		GDifficulty: DifficultyNormal,
		// 当前古董信息
		GAntiqueInfo: AntiqueInfo{},
		// 国内全部物资
		GItemInsInfo: make(map[int]ItemInfo),
		// 国外全部物资
		GItemOutInfo: make(map[int]ItemInfo),
		// 所有公司信息
		GCompanyInfo: make(map[int]CompanyInfo),
		// 游戏中一些属性最大数量
		GMaxHoldNum: MaxHoldNum{
			// 最大持有国内物资数量
			MDHoldNum: 100,
			// 最大持有国外物资数量
			MFHoldNum: 100,
			// 公司最大数
			MCHoldNum: CompanyNum,
			// 最大持有古董数量
			MAHoldNum: UserAntiqueNum,
			// 最大名声
			MFAHoldNum: MaxFame,
			// 最大免疫力
			MIHoldNum: MaxImmunity,
			// 每回合打工最大次数
			MWRoundNum: 10,
			// 每回合玩游戏最大次数
			MGRoundNum: 10,
			// 每回合约会最大次数
			MMRoundNum: 10,
			// 每回合逛街最大次数
			MSRoundNum: 10,
			// 每回合参加拍卖会最大次数
			MARoundNum: 10,
		},
		// 银行任务统计
		GBankTaskStats: BankTaskStats{
			DepositAmount:  0,
			DepositCount:   0,
			WithdrawAmount: 0,
			WithdrawCount:  0,
			LoanAmount:     0,
			WorkCount:      0,
			CompletedTasks: []int{},
			ClaimedTasks:   []int{},
			CurrentTasks:   utils.RandomSample(GetBankTasks(), 3),
		},
		// 所有约会对象信息
		GDatingInfo: GetDatings(),
		// 股票新闻
		GStockNews:        []string{},
		GStockUpdateCount: 0,
		// 房子信息
		GHouseInfo: map[int]HouseInfo{},
		// 车子信息
		GCarInfo: map[int]CarInfo{},
	}

	// ==============生成房子信息=================
	for _, house := range GetHouses() {
		// 计算当前价格
		currentPrice := CalculateHousePrice(house)
		// 添加到房子信息列表
		game.GHouseInfo[house.HId] = HouseInfo{
			HIId:     house.HId,
			HIName:   house.HName,
			HIPrice:  currentPrice,
			HIHealth: CalculateLifestyleHealth(house.HHealth),
			HIFame:   CalculateLifestyleFame(house.HFame),
			HIImg:    house.HImg,
		}
	}
	// ==============生成车子信息=================
	for _, car := range GetCars() {
		// 计算当前价格
		currentPrice := CalculateCarPrice(car)
		// 添加到车子信息列表
		game.GCarInfo[car.CId] = CarInfo{
			CIId:     car.CId,
			CIName:   car.CName,
			CIPrice:  currentPrice,
			CIHealth: CalculateLifestyleHealth(car.CHealth),
			CIFame:   CalculateLifestyleFame(car.CFame),
			CIImg:    car.CImg,
		}
	}
	// ==============生成国内外物资=================
	err := game.GetRandItem()
	if err != nil {
		return nil, err
	}
	// ==============生成创业公司信息=================
	// 初始化公司信息
	companies := GetCompanies()
	for _, company := range companies {
		game.GCompanyInfo[company.CId] = CompanyInfo{
			CIName:   company.CName,
			CIPrice:  GenerateStartPrice(company.CPrice, company.CRisk),
			CIRisk:   company.CRisk,
			CIProfit: company.CProfit,
			CITime:   company.CTime,
			CIStatus: true,
		}
	}
	// ==============生成股票=================
	// 初始化股票信息
	stocks := GetStocks()
	for _, sd := range stocks {
		// 设置一个随机价格20%涨幅
		priceTmp := int(float64(sd.SPrice) * (1 + 0.20*(rand.Float64()*2-1)))
		stockTmp := StockInfo{
			SIId:      sd.SId,
			SIName:    sd.SName,
			SIPrice:   priceTmp,
			SIRisk:    sd.SRisk,
			SIHistory: []int{priceTmp},
			SIKlineHistory: [][4]int{
				// [开盘, 收盘, 最低, 最高]
				{priceTmp, priceTmp, priceTmp, priceTmp},
			},
			SIStatus: "",
		}
		game.GStockInfo = append(game.GStockInfo, stockTmp)
	}
	return game, nil
}

// 随机选择当前需要展示的物资
func (g *Game) GetRandItem() (err error) {
	// 刷新国内物资
	g.GItemInsInfo, err = refreshItem(g.GItemInsInfo, GetItemIns(), false)
	if err != nil {
		return err
	}
	// 刷新国外物资，国外物资由于涨价值大，所以设置了涨价更难的概率
	g.GItemOutInfo, err = refreshItem(g.GItemOutInfo, GetItemOut(), true)
	return err
}

// 随机设置创业信息
func (g *Game) GetRandCompany() []int {
	bankruptcy := []int{}
	for i, cinfo := range g.GCompanyInfo {
		// 如果公司是新建的，则设置初始价格
		if !cinfo.CIStatus {
			cinfo.CIPrice = GenerateStartPrice(cinfo.CIPrice, cinfo.CIRisk)
			cinfo.CIStatus = true
		}
		// ✨ 黑马事件（高风险低收益，5% 概率翻盘）
		if cinfo.CIRisk >= 7 && cinfo.CIProfit <= 4 && rand.Float64() < 0.05 {
			cinfo.CIProfit += rand.Intn(3) + 2
			if cinfo.CIProfit > 10 {
				cinfo.CIProfit = 10
			}
			boost := float64(cinfo.CIPrice) * (0.1 + rand.Float64()*0.2)
			cinfo.CIPrice = int(math.Round(float64(cinfo.CIPrice) + boost))
		}

		// === 价格波动 ===
		price := float64(cinfo.CIPrice)
		riskRatio := float64(cinfo.CIRisk) / 10.0
		profitRatio := float64(cinfo.CIProfit) / 10.0
		BaseBankruptRate := 0.001                // 最低破产率
		MaxVolatilityRate := 0.30                // 最大波动
		RevertStrength := 0.03                   // 回归趋势修正因子
		cpTmp := float64(GetCompanyBasePrice(i)) // 每个商品基准

		volatilityRange := riskRatio * MaxVolatilityRate
		rnd := rand.Float64()*2 - 1
		profitBias := (profitRatio - 0.5) * 0.4
		rnd += profitBias

		// 限制范围 [-1, 1]
		if rnd > 1 {
			rnd = 1
		} else if rnd < -1 {
			rnd = -1
		}

		priceChange := price * volatilityRange * rnd
		deviation := (cpTmp - price) / cpTmp
		priceChange += deviation * cpTmp * RevertStrength

		price += priceChange
		if price < 1 {
			price = 1
		}
		cinfo.CIPrice = int(math.Round(price))

		// === 平衡破产机制（温和修正）===
		base := BaseBankruptRate + math.Pow(riskRatio, 1.5)*0.15
		profitOffset := math.Sqrt(profitRatio) * 0.08

		unstablePenalty := 0.0
		if cinfo.CIProfit > cinfo.CIRisk {
			diff := float64(cinfo.CIProfit - cinfo.CIRisk)
			unstablePenalty = diff * 0.01
		}

		startupPenalty := 0.0
		if cinfo.CITime <= 2 {
			startupPenalty = 0.01
		}

		chance := (base - profitOffset + unstablePenalty + startupPenalty) * GetDifficultyConfig(g.GDifficulty).BankruptRate

		// 限制破产率在合理区间
		if chance < BaseBankruptRate {
			chance = BaseBankruptRate
		}
		if chance > 0.5 {
			chance = 0.5
		}
		// 判断是否破产
		if rand.Float64() < chance {
			cinfo.CIStatus = false
			bankruptcy = append(bankruptcy, i)
			// 最后设置数据
			g.GCompanyInfo[i] = cinfo
			continue
		} else {
			// 最后设置数据
			g.GCompanyInfo[i] = cinfo
		}
	}
	return bankruptcy
}

// 计算公司收益
func (g *Game) CalculateCompanyProfit(userinfo *User) (announceCompany []string) {
	for i, uc := range userinfo.UCompany {
		// 判断公司是否可以开始收益
		if uc.UCompanyHoldTime < g.GCompanyInfo[i].CITime {
			announceCompany = append(announceCompany, "⏳ "+uc.UCompanyName+"正在融资中")
			continue
		}
		profitTmp := float64(uc.UCompanyNum*g.GCompanyInfo[i].CIPrice*g.GCompanyInfo[i].CIProfit) * 0.01
		userinfo.UCash += max(int(math.Round(profitTmp)), 1)
		announceCompany = append(announceCompany, "🧧 "+g.GCompanyInfo[i].CIName+"去年获得分红"+strconv.Itoa(max(int(math.Round(profitTmp)), 1))+"元")
	}
	return announceCompany
}

// 更新股票价格信息
func (g *Game) UpdateStockPrices() {
	gstockinfoTmp := g.GStockInfo
	stockNews := GetStockNews()
	// 遍历股票信息
	for i, gsinfo := range gstockinfoTmp {
		// 计算长度
		historyLen := len(gsinfo.SIHistory)
		// 上一次收盘价（也就是当前价格）
		lastClose := gsinfo.SIPrice
		if historyLen >= 1 {
			lastClose = gsinfo.SIHistory[historyLen-1]
		}
		// 前一天的收盘价，如果没有（第一天），就用 lastClose 代替
		yesterdayClose := lastClose
		if historyLen > 1 {
			yesterdayClose = gsinfo.SIHistory[historyLen-2]
		}
		// 今日开盘价，等于上一次收盘价
		open := lastClose
		// 随机波动范围（模拟市场波动）
		// 非线性风险因子（1.35 的幂次）：风险越高，变化越大
		riskMultiplier := math.Pow(1.35, float64(gsinfo.SIRisk)) // 1~5 => 1.35~5.25
		// 基础波动幅度（可调节），结合正态波动（均值0，波动围绕0）
		baseVolatility := 3.0 // 放大基础波动
		change := rand.NormFloat64() * baseVolatility * riskMultiplier
		// 是否触发新闻影响（30% 几率）
		if rand.Float64() < 0.3 && len(stockNews) > 0 {
			eventTmp := stockNews[rand.Intn(len(stockNews))]
			eventTmp = strings.Replace(eventTmp, "XXXX", gsinfo.SIName, -1)

			// 查找关键词并影响价格
			for keyword, effect := range StockNewsKeywordEffect {
				if strings.Contains(eventTmp, keyword) {
					if effect > 0 && change > 0 {
						eventTmp = "📈 " + eventTmp
						change = change + effect
					} else if effect < 0 && change < 0 {
						eventTmp = "📉 " + eventTmp
						change = change + effect
					} else {
						eventTmp = ""
					}
					if eventTmp != "" {
						// 插入到头部
						g.GStockNews = append([]string{eventTmp}, g.GStockNews...)
					}
					if len(g.GStockNews) > 30 {
						g.GStockNews = g.GStockNews[:30]
					}
					break
				}
			}
		}

		// 计算新的价格（收盘价）
		newPrice := int(math.Round(float64(lastClose) + change))
		// 计算涨停价（昨收×1.3）和跌停价（昨收×0.8）
		maxPrice := int(math.Round(float64(yesterdayClose) * 1.3))
		minPrice := int(math.Round(float64(yesterdayClose) * 0.8))

		// 限制新价格不超过涨跌停范围，并设置状态
		if newPrice > maxPrice {
			newPrice = maxPrice
			g.GStockInfo[i].SIStatus = "涨停"
		} else if newPrice < minPrice {
			newPrice = minPrice
			g.GStockInfo[i].SIStatus = "跌停"
		} else {
			g.GStockInfo[i].SIStatus = ""
		}

		// 模拟今日最高价和最低价（在开盘和收盘价基础上加减浮动）
		close := newPrice
		high := int(math.Max(float64(open), float64(close))) + rand.Intn(gsinfo.SIRisk+1) // 风险高，高点可能更高
		low := int(math.Min(float64(open), float64(close))) - rand.Intn(gsinfo.SIRisk+1)  // 风险高，低点可能更低

		// 更新当前价格，并确保价格不为负（最低为 1）
		g.GStockInfo[i].SIPrice = int(math.Max(1, float64(newPrice)))

		// 将新价格加入价格历史，最多保留 30 条
		g.GStockInfo[i].SIHistory = append(gsinfo.SIHistory, newPrice)
		// 是移除 klineHistory 数组中的第一项（最旧的数据）
		if len(g.GStockInfo[i].SIHistory) > 30 {
			g.GStockInfo[i].SIHistory = g.GStockInfo[i].SIHistory[1:]
		}

		// 将新的 K 线数据加入历史，结构为 [开盘, 收盘, 最低, 最高]，最多保留 30 条
		kline := [4]int{open, close, low, high}
		g.GStockInfo[i].SIKlineHistory = append(g.GStockInfo[i].SIKlineHistory, kline)
		if len(g.GStockInfo[i].SIKlineHistory) > 30 {
			g.GStockInfo[i].SIKlineHistory = g.GStockInfo[i].SIKlineHistory[1:]
		}
	}
}

// 公告处理
func (g *Game) UpdateAnnounce(announce Announce) Announce {
	// 获取公告信息
	for _, v := range g.GItemInsInfo {
		if !strings.Contains(v.IIEffect, "⚖️") {
			announce.AnnounceIns = append(announce.AnnounceIns, v.IIEffect)
		}
	}
	// 判断是否有公告
	if len(announce.AnnounceIns) == 0 {
		announce.AnnounceIns = append(announce.AnnounceIns, "当前市场行情没有发生任何变化")
	}
	// 获取国外公告信息
	for _, v := range g.GItemOutInfo {
		if !strings.Contains(v.IIEffect, "⚖️") {
			announce.AnnounceOut = append(announce.AnnounceOut, v.IIEffect)
		}
	}
	// 判断是否有公告
	if len(announce.AnnounceOut) == 0 {
		announce.AnnounceOut = append(announce.AnnounceOut, "当前市场行情没有发生任何变化")
	}
	// 获取用户持有的公司信息
	if len(announce.AnnounceCompany) == 0 {
		announce.AnnounceCompany = append(announce.AnnounceCompany, "您还没有成立公司")
	}
	// 游戏公告
	if len(announce.AnnounceGame) == 0 {
		announce.AnnounceGame = append(announce.AnnounceGame, "没有什么有价值的公告")
	}
	// 健康公告
	if len(announce.AnnounceHealthy) == 0 {
		announce.AnnounceHealthy = append(announce.AnnounceHealthy, "你很健康，没有什么疾病")
	}
	return announce
}
