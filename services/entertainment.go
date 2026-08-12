package services

import (
	"LifeGame/core"
	"math"
	"math/rand"
)

// EntertainmentActivity 是由后端维护的娱乐活动配置，避免前后端各自维护经济数值。
type EntertainmentActivity struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Desc       string `json:"desc"`
	EntryCost  int    `json:"entryCost"`
	HealthGain int    `json:"healthGain"`
	FameGain   int    `json:"fameGain"`
}

var entertainmentActivities = []EntertainmentActivity{
	{ID: "film", Name: "看电影", Icon: "🎬", Desc: "观看最新电影，放松心情", EntryCost: 100, HealthGain: 3},
	{ID: "keepfit", Name: "健身", Icon: "🏃", Desc: "到健身房锻炼身体", EntryCost: 500, HealthGain: 5, FameGain: 1},
	{ID: "tour", Name: "旅游", Icon: "✈️", Desc: "短途旅行，开阔视野", EntryCost: 5000, HealthGain: 10, FameGain: 2},
	{ID: "concert", Name: "音乐会", Icon: "🎵", Desc: "欣赏高雅音乐会", EntryCost: 8000, HealthGain: 8, FameGain: 3},
	{ID: "food", Name: "美食", Icon: "🍽️", Desc: "品尝精致美食", EntryCost: 2000, HealthGain: 4, FameGain: 1},
	{ID: "spa", Name: "SPA", Icon: "💆", Desc: "享受SPA按摩放松", EntryCost: 3000, HealthGain: 7, FameGain: 1},
	{ID: "golf", Name: "高尔夫", Icon: "⛳", Desc: "打高尔夫球，商务社交", EntryCost: 15000, HealthGain: 6, FameGain: 4},
	{ID: "party", Name: "派对", Icon: "🎉", Desc: "参加高端社交派对", EntryCost: 20000, HealthGain: 5, FameGain: 5},
}

func findEntertainmentActivity(id string) (EntertainmentActivity, bool) {
	for _, activity := range entertainmentActivities {
		if activity.ID == id {
			return activity, true
		}
	}
	return EntertainmentActivity{}, false
}

// GetEntertainmentActivities 返回后端权威的娱乐活动配置。
func (a *App) GetEntertainmentActivities() H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	return M{
		"code":       200,
		"activities": entertainmentActivities,
	}
}

// DoEntertainment 在后端完成扣款、次数与属性结算，保证存档和界面状态一致。
func (a *App) DoEntertainment(activityID string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}

	activity, ok := findEntertainmentActivity(activityID)
	if !ok {
		return M{"code": -1, "msg": "娱乐活动不存在"}
	}
	if a.Userinfo.UOpportunity.OSNum >= a.Gameinfo.GMaxHoldNum.MSRoundNum {
		return M{"code": -1, "msg": "本年度娱乐次数已用完，请明年再来"}
	}
	if a.Userinfo.UCash < activity.EntryCost {
		return M{"code": -1, "msg": "资金不足"}
	}

	// 同一年度连续娱乐收益递减，避免十次重复派对直接刷满声望。
	repeatMultiplier := math.Max(0.25, 1.0-float64(a.Userinfo.UOpportunity.OSNum)*0.15)
	healthMultiplier := (0.7 + rand.Float64()*0.6) * repeatMultiplier
	fameMultiplier := 1.0
	if activity.FameGain > 0 {
		fameMultiplier = (0.7 + rand.Float64()*0.6) * repeatMultiplier
	}
	healthChange := int(float64(activity.HealthGain)*healthMultiplier + 0.5)
	fameChange := int(float64(activity.FameGain)*fameMultiplier + 0.5)

	a.Userinfo.UCash -= activity.EntryCost
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + healthChange)
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + fameChange)
	a.Userinfo.UOpportunity.OSNum++
	if a.Userinfo.UMiniGameRecords == nil {
		a.Userinfo.UMiniGameRecords = make(map[string]core.MiniGameRecord)
	}
	record := a.Userinfo.UMiniGameRecords[activity.ID]
	record.MGRType = "activity"
	record.PlayCount++
	a.Userinfo.UMiniGameRecords[activity.ID] = record
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return M{
		"code":         200,
		"msg":          "参与成功",
		"healthchange": healthChange,
		"famechange":   fameChange,
		"userinfo":     a.userSnapshot(),
	}
}
