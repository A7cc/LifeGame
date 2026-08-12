package services

import (
	"LifeGame/core"
	"fmt"
	"strings"
)

// 获取医院信息
func (a *App) GetHospitalInfo() H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 判断当前是否必须先进行急诊
	emergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty)

	return M{
		"code":          200,
		"hospitalcards": core.GetHospitals(),
		"treatments":    core.GetTreats(),
		"emergency":     emergency,
	}
}

// BuyTreatment 买药
func (a *App) BuyTreatment(hospitalType string, treatmentId int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 判断当前是否必须先进行急诊
	if emergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty); emergency.Required {
		return M{
			"code":      -1,
			"msg":       fmt.Sprintf("当前必须先接受急诊治疗，费用 %d 元", emergency.Cost),
			"emergency": emergency,
		}
	}
	if len(a.Userinfo.UDiseases) == 0 && a.Userinfo.UImmunity >= 100 {
		return M{"code": -1, "msg": "当前免疫力已满，不需要自行服药"}
	}
	hospitalExists := false
	for _, h := range core.GetHospitals() {
		if h.HType == hospitalType {
			hospitalExists = true
			break
		}
	}
	if !hospitalExists {
		return M{"code": -1, "msg": "医院不存在"}
	}

	// 找到治疗方式
	var treatment *core.TreatInfo
	for _, t := range core.GetTreats() {
		if t.TId == treatmentId {
			treatmentInfo := t
			treatment = &treatmentInfo
			break
		}
	}

	if treatment == nil {
		return M{"code": -1, "msg": "药品不存在"}
	}
	if treatment.TSource != hospitalType && treatment.TSource != "pharmacy" {
		return M{"code": -1, "msg": "该医院没有这个药品"}
	}

	// 检查资金
	if a.Userinfo.UCash < treatment.TPrice {
		return M{"code": -1, "msg": "现金不足"}
	}

	// 扣钱
	a.Userinfo.UCash -= treatment.TPrice
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	// 恢复免疫力
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treatment.THeal)

	// 处理副作用
	if treatment.TSideEffect < 0 {
		a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treatment.TSideEffect)
	}

	// 检查是否治愈疾病
	curedDiseases := []core.UDiseaseInfo{}
	udiseaseTmp := a.Userinfo.UDiseases
	for i, d := range udiseaseTmp {
		for _, t := range d.UTreatments {
			// 如果治疗方式在疾病的治疗方式中，则治愈
			if strings.Contains(t, treatment.TName) {
				if d.UDSeverity > 4 {
					continue
				} else {
					curedDiseases = append(curedDiseases, d)
					delete(a.Userinfo.UDiseases, i)
				}
			}
		}
	}
	core.ResetCriticalHealthIfRecovered(a.Userinfo)

	return M{
		"code":          200,
		"msg":           "购买成功！",
		"userinfo":      a.userSnapshot(),
		"curedDiseases": curedDiseases,
	}
}

// EmergencyTreatment 急诊治疗：严重疾病或免疫力过低时必须先支付高价处理。
func (a *App) EmergencyTreatment() H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}

	// 判断当前是否必须先进行急诊
	emergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty)
	if !emergency.Required {
		return M{
			"code":      -1,
			"msg":       "当前不需要急诊治疗",
			"emergency": emergency,
		}
	}
	if a.Userinfo.UCash < emergency.Cost {
		return M{
			"code":      -1,
			"msg":       fmt.Sprintf("急诊费用需要 %d 元，现金不足", emergency.Cost),
			"emergency": emergency,
		}
	}

	// 支付费用
	a.Userinfo.UCash -= emergency.Cost
	// 执行急诊治疗，清除严重疾病并把免疫力拉回安全线。
	curedDiseases := core.ApplyHealthEmergencyTreatment(a.Userinfo)
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	// 再次判断当前是否必须先进行急诊
	nextEmergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty)

	msg := fmt.Sprintf("急诊治疗完成，花费 %d 元", emergency.Cost)
	if len(curedDiseases) > 0 {
		names := make([]string, 0, len(curedDiseases))
		for _, disease := range curedDiseases {
			names = append(names, disease.UDName)
		}
		msg += "，处理了" + strings.Join(names, "、")
	}

	return M{
		"code":          200,
		"msg":           msg,
		"cost":          emergency.Cost,
		"userinfo":      a.userSnapshot(),
		"emergency":     nextEmergency,
		"curedDiseases": curedDiseases,
	}
}

// SpecialTreatment 特殊治疗
func (a *App) SpecialTreatment(treatType, hospitalType string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	// 判断当前是否必须先进行急诊
	if emergency := core.GetHealthEmergencyStatus(a.Userinfo, a.Gameinfo.GDifficulty); emergency.Required {
		return M{
			"code":      -1,
			"msg":       fmt.Sprintf("当前必须先接受急诊治疗，费用 %d 元", emergency.Cost),
			"emergency": emergency,
		}
	}
	var treat *core.HospitalServiceInfo
	for _, h := range core.GetHospitals() {
		if h.HType == hospitalType {
			for _, hs := range h.HServices {
				if hs.HSName == treatType {
					serviceInfo := hs
					treat = &serviceInfo
					break
				}
			}
		}
	}
	if treat == nil {
		return M{"code": -1, "msg": "你的这个治疗方式不存在"}
	}

	// 检查资金
	if a.Userinfo.UCash < treat.HSPrice {
		return M{"code": -1, "msg": "现金不足"}
	}

	// 扣钱
	a.Userinfo.UCash -= treat.HSPrice
	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	curedDiseases := ""
	// 根据游戏得分计算治疗效果
	switch treat.HSName {
	case "打针":
		udiseaseTmp := a.Userinfo.UDiseases
		for i, d := range udiseaseTmp {
			for _, t := range d.UTreatments {
				// 如果治疗方式在疾病的治疗方式中，则治愈
				if strings.Contains(t, treat.HSName) {
					d.UDSeverity -= 4
					curedDiseases += d.UDName + "，"
					if d.UDSeverity <= 0 {
						delete(a.Userinfo.UDiseases, i)
					} else {
						a.Userinfo.UDiseases[i] = d
					}
				}
			}
		}
		if curedDiseases == "" {
			curedDiseases = "打针治疗成功！"
		} else {
			curedDiseases = "打针治疗成功！治疗了" + strings.TrimRight(curedDiseases, "，")
		}
		// 增加免疫力
		a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treat.HSImmunity)
	case "针灸":
		curedDiseases = "针灸治疗完成！"
		// 增加免疫力
		a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treat.HSImmunity)
	case "手术":
		// 治愈所有需要手术的疾病
		udiseaseTmp := a.Userinfo.UDiseases
		for i, d := range udiseaseTmp {
			for _, t := range d.UTreatments {
				// 如果治疗方式在疾病的治疗方式中，则治愈
				if strings.Contains(t, treat.HSName) {
					d.UDSeverity -= 4
					curedDiseases += d.UDName + "，"
					if d.UDSeverity <= 0 {
						delete(a.Userinfo.UDiseases, i)
					} else {
						a.Userinfo.UDiseases[i] = d
					}
				}
			}
		}
		if curedDiseases == "" {
			curedDiseases = "手术治疗成功！"
		} else {
			curedDiseases = "手术治疗成功！治疗了" + strings.TrimRight(curedDiseases, "，")
		}
		// 增加免疫力
		a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treat.HSImmunity)
	case "脱胎换骨":
		// 治愈所有需要手术的疾病
		udiseaseTmp := a.Userinfo.UDiseases
		for i, d := range udiseaseTmp {
			for _, t := range d.UTreatments {
				// 如果治疗方式在疾病的治疗方式中，则治愈
				if strings.Contains(t, treat.HSName) {
					d.UDSeverity -= 4
					curedDiseases += d.UDName + "，"
					if d.UDSeverity <= 0 {
						delete(a.Userinfo.UDiseases, i)
					} else {
						a.Userinfo.UDiseases[i] = d
					}
				}
			}
		}
		if curedDiseases == "" {
			curedDiseases = "脱胎换骨治疗成功！"
		} else {
			curedDiseases = "脱胎换骨治疗成功！治疗了" + strings.TrimRight(curedDiseases, "，")
		}
		// 增加免疫力
		a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + treat.HSImmunity)
	default:
		return M{"code": -1, "msg": "你的这个治疗方式不存在"}
	}
	core.ResetCriticalHealthIfRecovered(a.Userinfo)

	return M{
		"code":     200,
		"msg":      curedDiseases,
		"userinfo": a.userSnapshot(),
	}
}
