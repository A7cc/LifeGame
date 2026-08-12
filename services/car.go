package services

import (
	"LifeGame/core"
	"fmt"
)

// 购买车辆
func (a *App) BuyCar(carId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	targetCar := a.Gameinfo.GCarInfo[carId]
	// 检查车辆是否存在
	if targetCar == (core.CarInfo{}) {
		return M{
			"code": -1,
			"msg":  "车辆不存在",
		}
	}

	// 检查是否已拥有此车辆
	if a.Userinfo.UCar[carId] {
		return M{
			"code": -1,
			"msg":  "你已经拥有这辆车了",
		}
	}

	// 检查资金
	currentPrice := targetCar.CIPrice
	if a.Userinfo.UCash < currentPrice {
		return M{
			"code": -1,
			"msg":  fmt.Sprintf("现金不足！需要 %d 元", currentPrice),
		}
	}

	// 扣钱并添加车辆
	a.Userinfo.UCash -= currentPrice
	// 计算车辆带来的名声和免疫力
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + targetCar.CIFame)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + targetCar.CIHealth)
	// 添加车辆
	a.Userinfo.UCar[carId] = true
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功购买 %s！", targetCar.CIName),
		"userinfo": a.userSnapshot(),
	}
}

// 出售车辆
func (a *App) SellCar(carId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	targetCar := a.Gameinfo.GCarInfo[carId]
	// 检查车辆是否存在
	if targetCar == (core.CarInfo{}) {
		return M{
			"code": -1,
			"msg":  "车辆不存在",
		}
	}

	// 检查是否已拥有此车辆
	if !a.Userinfo.UCar[carId] {
		return M{
			"code": -1,
			"msg":  "你还没有这辆车",
		}
	}

	// 计算出售价格（当前市价的 70%）
	sellPrice := targetCar.CIPrice * 7 / 10

	// 加钱并移除车辆
	a.Userinfo.UCash += sellPrice
	// 计算车辆带来的名声和免疫力
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - targetCar.CIFame)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity - targetCar.CIHealth)
	// 移除车辆
	delete(a.Userinfo.UCar, carId)
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":     200,
		"msg":      fmt.Sprintf("成功出售 %s！获得 %d 元", targetCar.CIName, sellPrice),
		"userinfo": a.userSnapshot(),
	}
}
