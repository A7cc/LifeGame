package core

import (
	"errors"
	"math"
)

// 购买物资（国内buyDomesticGoods/国外buyForeignGoods）
func (u *User) BuyDomesticGoods(gameinfo *Game, itemid, snum int) error {
	// 设置国内
	GItemInfoTmp := gameinfo.GItemInsInfo
	UItemInfoTmp := u.UItemins
	UItemInfoMaxCountTmp := gameinfo.GMaxHoldNum.MDHoldNum
	// 判断是否超出购买数量
	for _, uIInfoTmp := range UItemInfoTmp {
		UItemInfoMaxCountTmp = UItemInfoMaxCountTmp - uIInfoTmp.UIINum
	}
	if (UItemInfoMaxCountTmp - snum) < 0 {
		return errors.New("购买商品数量超出最大物资持有量！")
	}
	// 判断是否物品编号是否在用户购买的物品列表中
	for i, iteminfo := range GItemInfoTmp {
		if i == itemid {
			if iteminfo.IIDisplay {
				// 用户买东西
				icashtmp := int(iteminfo.IIPrice) * snum
				// 检测用户资产是否足够
				if u.UCash < icashtmp || icashtmp <= 0 {
					return errors.New("资产不足，无法购买！")
				}
				// 检测用户是否拥有该物品，如果没有则添加，如果有则增加数量
				if uIITmp, exists := UItemInfoTmp[itemid]; !exists {
					uIITmp.UIICostPrice = iteminfo.IIPrice
					uIITmp.UIINum = snum
					UItemInfoTmp[itemid] = uIITmp
				} else {
					uIITmp.UIICostPrice = (uIITmp.UIICostPrice*uIITmp.UIINum + iteminfo.IIPrice*snum) / (uIITmp.UIINum + snum)
					uIITmp.UIINum += snum
					UItemInfoTmp[itemid] = uIITmp
				}
				// 扣除用户资产
				u.UCash -= icashtmp
				return nil
			}
			return errors.New("物品不存在！")
		}
	}
	return errors.New("购买失败")

}
func (u *User) BuyForeignGoods(gameinfo *Game, itemid, snum int) error {
	// 判断是否是国外
	GItemInfoTmp := gameinfo.GItemOutInfo
	UItemInfoTmp := u.UItemout
	UItemInfoMaxCountTmp := gameinfo.GMaxHoldNum.MFHoldNum

	// 判断是否超出购买数量
	for _, uIInfoTmp := range UItemInfoTmp {
		UItemInfoMaxCountTmp = UItemInfoMaxCountTmp - uIInfoTmp.UIINum
	}
	if (UItemInfoMaxCountTmp - snum) < 0 {
		return errors.New("购买商品数量超出最大物资持有量！")
	}
	// 判断是否物品编号是否在用户购买的物品列表中
	for i, iteminfo := range GItemInfoTmp {
		if i == itemid {
			if iteminfo.IIDisplay {
				// 用户买东西
				icashtmp := int(iteminfo.IIPrice) * snum
				// 检测用户资产是否足够
				if u.UCash < icashtmp || icashtmp <= 0 {
					return errors.New("资产不足，无法购买！")
				}
				// 检测用户是否拥有该物品，如果没有则添加，如果有则增加数量
				if uIITmp, exists := UItemInfoTmp[itemid]; !exists {
					uIITmp.UIICostPrice = iteminfo.IIPrice
					uIITmp.UIINum = snum
					UItemInfoTmp[itemid] = uIITmp
				} else {
					uIITmp.UIICostPrice = (uIITmp.UIICostPrice*uIITmp.UIINum + iteminfo.IIPrice*snum) / (uIITmp.UIINum + snum)
					uIITmp.UIINum += snum
					UItemInfoTmp[itemid] = uIITmp
				}
				// 扣除用户资产
				u.UCash -= icashtmp
				return nil
			}
			return errors.New("物品不存在！")
		}
	}
	return errors.New("购买失败")
}

// buyCompany 购买公司
func (u *User) BuyCompany(gameinfo *Game, itemid, snum int) error {
	companyInfo, exists := gameinfo.GCompanyInfo[itemid]
	if !exists {
		return errors.New("公司不存在！")
	}

	// 判断当前公司是否有被购买
	isNew := u.UCompany[itemid].UCompanyNum <= 0
	// 创业公司数量同时受全局上限和名声等级控制。
	_, _, _, reputationCompanyLimit := CalcReputationLevel(u.UFame)
	companyLimit := min(gameinfo.GMaxHoldNum.MCHoldNum, reputationCompanyLimit)
	if isNew && len(u.UCompany)+1 > companyLimit {
		return errors.New("创业数量已达到上限！")
	}
	// 首次买入必须大于等于1000
	if isNew && snum < 1000 {
		return errors.New("最少购买1000手！")
	}
	// 判断公司是否破产
	if !companyInfo.CIStatus {
		return errors.New("该公司已经破产！不可购买！")
	}
	// 计算当前购买的钱
	pTmp := companyInfo.CIPrice * snum
	// 检测用户资产是否足够
	if u.UCash < pTmp || pTmp <= 0 {
		return errors.New("资产不足，无法购买！")
	}
	// 扣钱
	u.UCash -= pTmp
	// 计算原先的总价和当前购买的钱
	ucTmp := u.UCompany[itemid]
	// 设置公司名字
	ucTmp.UCompanyName = companyInfo.CIName
	pTmp += ucTmp.UCompanyCostPrice * ucTmp.UCompanyNum
	// 增加数量
	ucTmp.UCompanyNum += snum
	// 小数点忽略不计
	ucTmp.UCompanyCostPrice = pTmp / ucTmp.UCompanyNum
	// 更新用户公司信息
	u.UCompany[itemid] = ucTmp
	return nil
}

// buyStock 购买股票
func (u *User) BuyStock(gameinfo *Game, itemid, snum int) error {
	// 计算用户的名声，如果名声小于出售的股票，则无法出售
	if _, _, sp, _ := CalcReputationLevel(u.UFame); sp < u.UStockProfit {
		return errors.New("今天盈利已经超出当前名气的金额，无法再购买了！")
	}
	// 最少购买10股票
	if snum < 10 {
		return errors.New("最少购买10！")
	}
	// 判断是否是当前股票
	ggsTmp := StockInfo{}
	found := false
	for _, v := range gameinfo.GStockInfo {
		if v.SIId == itemid {
			ggsTmp = v
			found = true
			break
		}
	}
	if !found {
		return errors.New("股票不存在！")
	}
	// 计算当前购买的钱并收取双边手续费。
	pTmp := ggsTmp.SIPrice * snum
	fee := int(math.Ceil(float64(pTmp) * StockTransactionFeeRate))
	totalCost := pTmp + fee
	// 检测用户资产是否足够
	if u.UCash < totalCost || pTmp <= 0 {
		return errors.New("资产不足，无法购买！")
	}
	// 跌停
	if ggsTmp.SIStatus == "跌停" {
		return errors.New("当前为跌停状态，买入存在风险")
	}
	// 判断是否拥有该股票
	// 扣钱
	u.UCash -= totalCost
	// Go 语言为了保证 map 的一致性和安全性，不允许你直接修改 map 中结构体的字段
	usTmp := u.UStock[itemid]
	// 计算数量和成本价
	if usTmp.USNum == 0 {
		usTmp.USName = ggsTmp.SIName
		usTmp.USPrice_init = totalCost / snum
		usTmp.USNum = snum
	} else {
		pTmp = totalCost + usTmp.USPrice_init*usTmp.USNum
		usTmp.USNum += snum
		usTmp.USPrice_init = pTmp / usTmp.USNum
	}
	// 更新用户股票信息
	u.UStock[itemid] = usTmp
	return nil
}

// 出售物资（国内/国外）
func (u *User) SellDomesticGoods(gameinfo *Game, itemid, snum int) (int, error) {
	GItemInfoTmp := gameinfo.GItemInsInfo
	UItemInfoTmp := u.UItemins
	// 判断是否物品编号是否在用户购买的物品列表中
	if uIITmp, ok := UItemInfoTmp[itemid]; ok {
		// 循环获取物品信息
		for i, iteminfo := range GItemInfoTmp {
			if i == itemid {
				// 这个判断可以写在一起，但是想的是都i=itemid了，可以退出了
				if iteminfo.IIDisplay && uIITmp.UIINum >= snum {
					// 判断物品编号是否在当前的物品列表中
					icashtmp := 0
					// 出售物品
					// 检测用户是否拥有该物品，如果没跳出，如果有则计算数量
					uIITmp.UIINum -= snum
					gross := int(iteminfo.IIPrice) * snum
					icashtmp = gross - int(math.Ceil(float64(gross)*MarketSellFeeRate))
					// 扣除用户资产
					u.UCash += icashtmp
					// 如果数量为0，则删除物资记录
					if uIITmp.UIINum == 0 {
						delete(UItemInfoTmp, itemid)
					} else {
						UItemInfoTmp[itemid] = uIITmp
					}
					return icashtmp, nil
				} else if uIITmp.UIINum < snum {
					return 0, errors.New("出售的数量超过了拥有的数量！")
				}
				return 0, errors.New("该物品不可出售")
			}
		}
	}
	return 0, errors.New("出售失败")
}

// 出售物资（国内/国外）
func (u *User) SellForeignGoods(gameinfo *Game, itemid, snum int) (int, error) {
	GItemInfoTmp := gameinfo.GItemOutInfo
	UItemInfoTmp := u.UItemout
	// 判断是否物品编号是否在用户购买的物品列表中
	if uIITmp, ok := UItemInfoTmp[itemid]; ok {
		// 循环获取物品信息
		for i, iteminfo := range GItemInfoTmp {
			if i == itemid {
				// 这个判断可以写在一起，但是想的是都i=itemid了，可以退出了
				if iteminfo.IIDisplay && uIITmp.UIINum >= snum {
					// 判断物品编号是否在当前的物品列表中
					icashtmp := 0
					// 出售物品
					// 检测用户是否拥有该物品，如果没跳出，如果有则计算数量
					uIITmp.UIINum -= snum
					gross := int(iteminfo.IIPrice) * snum
					icashtmp = gross - int(math.Ceil(float64(gross)*MarketSellFeeRate))
					// 扣除用户资产
					u.UCash += icashtmp
					// 如果数量为0，则删除物资记录
					if uIITmp.UIINum == 0 {
						delete(UItemInfoTmp, itemid)
					} else {
						UItemInfoTmp[itemid] = uIITmp
					}
					return icashtmp, nil
				} else if uIITmp.UIINum < snum {
					return 0, errors.New("出售的数量超过了拥有的数量！")
				}
				return 0, errors.New("该物品不可出售")
			}
		}
	}
	return 0, errors.New("出售失败")
}

// sellCompany 出售公司
func (u *User) SellCompany(gameinfo *Game, itemid, snum int) (int, error) {
	companyInfo, exists := gameinfo.GCompanyInfo[itemid]
	if !exists {
		return 0, errors.New("公司不存在！")
	}
	// 判断用户是否持有该公司
	if u.UCompany[itemid].UCompanyNum <= 0 {
		return 0, errors.New("你还没有持有该公司！")
	}

	// 判断用户要出售的数量是否超过拥有的数量
	if u.UCompany[itemid].UCompanyNum < snum {
		return 0, errors.New("出售的数量超过了拥有的数量！")
	}
	// 判断是否能够出售
	if u.UCompany[itemid].UCompanyHoldTime < companyInfo.CITime {
		return 0, errors.New("该公司正在融资中，不能进行出售交易！")
	}
	// 按当前公司市价结算，并收取10%的退出费用。
	ucTmp := u.UCompany[itemid]
	pTmp := companyInfo.CIPrice * snum * 9 / 10
	ucTmp.UCompanyNum -= snum
	// 计算总的盈亏
	u.UCash += pTmp
	// 如果数量为0，则删除
	if ucTmp.UCompanyNum == 0 {
		delete(u.UCompany, itemid)
	} else {
		u.UCompany[itemid] = ucTmp
	}
	return pTmp, nil
}

// sellStock 出售股票
func (u *User) SellStock(gameinfo *Game, itemid, snum int) (int, error) {
	// 计算用户的名声，如果名声小于出售的股票，则无法出售
	if _, _, sp, _ := CalcReputationLevel(u.UFame); sp < u.UStockProfit {
		return 0, errors.New("今天盈利已经超出当前名气的金额，无法再出售了！")
	}
	// 判断是否是当前股票
	ggsTmp := StockInfo{}
	found := false
	for _, v := range gameinfo.GStockInfo {
		if v.SIId == itemid {
			ggsTmp = v
			found = true
			break
		}
	}
	if !found {
		return 0, errors.New("股票不存在！")
	}
	// 用来存放用户股票信息
	usTmp := u.UStock[itemid]
	if usTmp.USNum <= 0 {
		return 0, errors.New("你还没有持有该股票！")
	}
	// 判断出售的数量超过拥有的数量
	if snum > usTmp.USNum {
		return 0, errors.New("出售的数量超过了拥有的数量！")
	}
	// 涨停
	if ggsTmp.SIStatus == "涨停" {
		return 0, errors.New("涨停状态下无法卖出")
	}
	// 计算当前出售的钱、手续费以及本次净利润，并在成交前校验盈利上限。
	gross := ggsTmp.SIPrice * snum
	fee := int(math.Ceil(float64(gross) * StockTransactionFeeRate))
	pTmp := gross - fee
	tradeProfit := pTmp - usTmp.USPrice_init*snum
	_, _, profitLimit, _ := CalcReputationLevel(u.UFame)
	if tradeProfit > 0 && u.UStockProfit+tradeProfit > profitLimit {
		return 0, errors.New("本次卖出将超过当前名气允许的年度股票盈利上限！")
	}
	// 加钱
	u.UCash += pTmp
	// 增加数量
	usTmp.USNum -= snum
	// 计算总的盈亏
	usTmp.USProfit += tradeProfit
	// 计算当天盈亏
	u.UStockProfit += tradeProfit
	// 如果数量为0，则删除股票记录
	if usTmp.USNum == 0 {
		delete(u.UStock, itemid)
	} else {
		// 更新用户股票信息
		u.UStock[itemid] = usTmp
	}
	return pTmp, nil
}
