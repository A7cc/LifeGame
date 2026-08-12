package services

import (
	"LifeGame/core"
	"fmt"
)

// 购买房屋
func (a *App) BuyHouse(houseId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	targetHouse := a.Gameinfo.GHouseInfo[houseId]
	// 检查房屋是否存在
	if targetHouse == (core.HouseInfo{}) {
		return M{
			"code": -1,
			"msg":  "房屋不存在",
		}
	}
	// 检查是否已经拥有此房屋
	if a.Userinfo.UHouse[houseId] {
		return M{
			"code": -1,
			"msg":  "你已经拥有这套房屋了",
		}
	}

	// 检查资金
	currentPrice := targetHouse.HIPrice
	if a.Userinfo.UCash < currentPrice {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("现金不足！需要 %d 元", currentPrice),
		}
	}

	// 扣钱并添加房屋
	a.Userinfo.UCash -= currentPrice
	// 计算房屋带来的名声和免疫力
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + targetHouse.HIFame)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + targetHouse.HIHealth)
	// 添加房屋
	a.Userinfo.UHouse[houseId] = true
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功购买 %s！", targetHouse.HIName),
		"userinfo": a.userSnapshot(),
	}
}

// 出售房屋
func (a *App) SellHouse(houseId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	targetHouse := a.Gameinfo.GHouseInfo[houseId]
	// 检查房屋是否存在
	if targetHouse == (core.HouseInfo{}) {
		return M{
			"code": -1,
			"msg":  "房屋不存在",
		}
	}
	// 检查是否已经拥有此房屋
	if !a.Userinfo.UHouse[houseId] {
		return M{
			"code": -1,
			"msg":  "你还没有这套房屋",
		}
	}

	// 计算出售价格（当前市价的 80%）
	sellPrice := targetHouse.HIPrice * 8 / 10

	// 加钱并移除房屋
	a.Userinfo.UCash += sellPrice
	// 计算房屋带来的名声和免疫力
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - targetHouse.HIFame)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity - targetHouse.HIHealth)
	// 移除房屋
	delete(a.Userinfo.UHouse, houseId)
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功出售 %s！获得 %d 元", targetHouse.HIName, sellPrice),
		"userinfo": a.userSnapshot(),
	}
}
