package services

import (
	"errors"
	"strconv"
)

/*
* BuyItem 购买物资
* @Param itemid 物资ID
* @Param snum 购买数量
* @Param region 购买的商品类型
 */
func (a *App) BuyItem(itemid, snum int, region string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	if snum <= 0 {
		return M{
			"code": -1,
			"msg":  "购买失败：购买数量必须大于0！",
		}
	}
	err := errors.New("购买失败：未知商品类型！")
	switch region {
	case "国内市场":
		err = a.Userinfo.BuyDomesticGoods(a.Gameinfo, itemid, snum)
	case "国外市场":
		err = a.Userinfo.BuyForeignGoods(a.Gameinfo, itemid, snum)
	case "创业":
		err = a.Userinfo.BuyCompany(a.Gameinfo, itemid, snum)
	case "股票":
		err = a.Userinfo.BuyStock(a.Gameinfo, itemid, snum)
	}
	if err != nil {
		return M{
			"code": -1,
			"msg":  err.Error(),
		}
	}
	// 计算用户资产
	_, _, err = a.Userinfo.RefreshAndValidateUserState(a.Gameinfo)
	if err != nil {
		return M{
			"code": -1,
			"msg":  err.Error(),
		}
	} else {
		return M{
			"code":     200,
			"userinfo": a.userSnapshot(),
		}
	}
}

// ===========================================
/*
* SellItem 出售物资
* @Param itemid 物资ID
* @Param snum 出售数量
* @Param region 出售的商品类型
 */
func (a *App) SellItem(itemid, snum int, region string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	if snum <= 0 {
		return M{
			"code": -1,
			"msg":  "出售失败：出售数量必须大于0！",
		}
	}

	err := errors.New("出售失败：未知商品类型！")
	money := 0
	// 根据商品类型，执行出售操作
	switch region {
	case "国内市场":
		money, err = a.Userinfo.SellDomesticGoods(a.Gameinfo, itemid, snum)
	case "国外市场":
		money, err = a.Userinfo.SellForeignGoods(a.Gameinfo, itemid, snum)
	case "创业":
		money, err = a.Userinfo.SellCompany(a.Gameinfo, itemid, snum)
	case "股票":
		money, err = a.Userinfo.SellStock(a.Gameinfo, itemid, snum)
	}
	if err != nil {
		return M{
			"code": -1,
			"msg":  err.Error(),
		}
	}

	// 计算用户资产
	_, _, err = a.Userinfo.RefreshAndValidateUserState(a.Gameinfo)
	if err != nil {
		return M{
			"code": -1,
			"msg":  err.Error(),
		}
	} else {
		return M{
			"code":     200,
			"msg":      "，获得" + strconv.Itoa(money) + "元！",
			"userinfo": a.userSnapshot(),
		}
	}
}
