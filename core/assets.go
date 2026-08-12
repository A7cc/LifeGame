package core

// CalculateUserAssets 计算用户总资产（不修改 UAssets 字段，仅返回计算值）
// 用于获取用户准确的资产总值，避免数据不一致
func CalculateUserAssets(userinfo *User, gameinfo *Game) int {
	// 初始化资产值
	assetstmp := 0
	// 1. 遍历用户国内物资
	for i, giinfo := range gameinfo.GItemInsInfo {
		if j, exists := userinfo.UItemins[i]; exists {
			if giinfo.IIDisplay {
				assetstmp = assetstmp + giinfo.IIPrice*j.UIINum
			} else {
				assetstmp = assetstmp + j.UIICostPrice*j.UIINum
			}
		}
	}

	// 2. 遍历用户国外物资
	for i, goinfo := range gameinfo.GItemOutInfo {
		if j, exists := userinfo.UItemout[i]; exists {
			if goinfo.IIDisplay {
				assetstmp = assetstmp + goinfo.IIPrice*j.UIINum
			} else {
				assetstmp = assetstmp + j.UIICostPrice*j.UIINum
			}
		}
	}

	// 3. 添加古董价值
	for _, ua := range userinfo.UAntique {
		assetstmp = assetstmp + ua.AIPrice
	}

	// 4. 添加公司当前市值
	for companyID, uc := range userinfo.UCompany {
		price := uc.UCompanyCostPrice
		if company, exists := gameinfo.GCompanyInfo[companyID]; exists {
			price = company.CIPrice
		}
		assetstmp = assetstmp + price*uc.UCompanyNum
	}

	// 5. 添加股票当前市值
	for stockId, us := range userinfo.UStock {
		var currentPrice int
		for _, gs := range gameinfo.GStockInfo {
			if gs.SIId == stockId {
				currentPrice = gs.SIPrice
				break
			}
		}
		assetstmp = assetstmp + currentPrice*us.USNum
	}

	// 6. 添加房屋价值（按基础价格计算）
	for _, h := range gameinfo.GHouseInfo {
		if userinfo.UHouse[h.HIId] {
			assetstmp = assetstmp + h.HIPrice
		}
	}

	// 7. 添加车辆价值（按基础价格计算）
	for _, c := range gameinfo.GCarInfo {
		if userinfo.UCar[c.CIId] {
			assetstmp = assetstmp + c.CIPrice
		}
	}

	// 8. 加上用户银行存款和现金，并扣除尚未偿还的贷款。
	// 贷款到账只会改变资产构成，不应凭空增加净资产。
	assetstmp = assetstmp + userinfo.UCash + userinfo.UBank - userinfo.ULoan

	return assetstmp
}
