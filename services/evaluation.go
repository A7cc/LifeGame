package services

import (
	"LifeGame/core"
	"fmt"
)

type GameEvaluation struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	Score             int    `json:"score"`
	WealthScore       int    `json:"wealthscore"`
	HealthScore       int    `json:"healthscore"`
	FameScore         int    `json:"famescore"`
	AgeScore          int    `json:"agescore"`
	CareerScore       int    `json:"careerscore"`
	RelationshipScore int    `json:"relationshipscore"`
	CollectionScore   int    `json:"collectionscore"`
}

func calculateEvaluation(userinfo *core.User) GameEvaluation {
	eval := GameEvaluation{}

	eval.WealthScore = min(25, max(0, userinfo.UAssets/1000000))

	eval.HealthScore = userinfo.UImmunity * 20 / 100

	eval.FameScore = min(15, max(0, userinfo.UFame*15/core.MaxFame))

	playedWork, gameWins := 0, 0
	for _, record := range userinfo.UMiniGameRecords {
		if record.MGRType == "work" {
			playedWork += record.PlayCount
		}
		if record.MGRType != "work" && record.MGRType != "activity" {
			gameWins += record.WinCount
		}
	}
	eval.CareerScore = min(4, playedWork/20) + min(3, len(userinfo.UCompany)) + min(3, gameWins/20)

	totalDates, totalGifts, bestRelationship := 0, 0, 0
	for _, relationship := range userinfo.UDating {
		totalDates += relationship.DCount
		totalGifts += relationship.DGiftCount
		base := 0
		switch core.NormalizeDatingStatus(relationship.DStatus) {
		case core.DatingStatusDating:
			base = 1
		case core.DatingStatusLover:
			base = 2
		case core.DatingStatusExclusive:
			base = 3
		case core.DatingStatusSweetheart:
			base = 4
		case core.DatingStatusMarried:
			base = 5
		}
		bestRelationship = max(bestRelationship, base)
	}
	eval.RelationshipScore = min(10, bestRelationship+min(3, totalDates/10)+min(2, totalGifts/5))

	trueAntiques := 0
	for _, antique := range userinfo.UAntique {
		if antique.AIDisplay == 1 {
			trueAntiques++
		}
	}
	eval.CollectionScore = min(2, trueAntiques/2) + min(1, len(userinfo.UHouse)/5) + min(1, len(userinfo.UCar)/5)
	if len(userinfo.UItemins)+len(userinfo.UItemout) > 0 {
		eval.CollectionScore++
	}

	if userinfo.UAge > core.UserAgeInit {
		eval.AgeScore = (userinfo.UAge - core.UserAgeInit) * 15 / (core.UserAgeMax - core.UserAgeInit)
	}
	if eval.AgeScore > 15 {
		eval.AgeScore = 15
	}

	eval.Score = eval.WealthScore + eval.HealthScore + eval.FameScore + eval.AgeScore + eval.CareerScore + eval.RelationshipScore + eval.CollectionScore
	eval.Title = getEvaluationTitle(eval.Score)
	eval.Description = getEvaluationDescription(eval, userinfo)

	return eval
}

func getEvaluationTitle(score int) string {
	switch {
	case score >= 90:
		return "👑 传奇人生"
	case score >= 75:
		return "🌟 辉煌人生"
	case score >= 60:
		return "👍 成功人生"
	case score >= 40:
		return "😊 平凡人生"
	case score >= 20:
		return "😔 艰难人生"
	default:
		return "💀 惨淡人生"
	}
}

func getEvaluationDescription(eval GameEvaluation, userinfo *core.User) string {
	var desc string
	desc += fmt.Sprintf("你的人生以 %d 岁结束\n", userinfo.UAge)
	desc += fmt.Sprintf("最终资产：%d 元\n", userinfo.UAssets)
	desc += fmt.Sprintf("最终免疫力：%d\n", userinfo.UImmunity)
	_, name, _, _ := core.CalcReputationLevel(userinfo.UFame)
	desc += fmt.Sprintf("最终名声：%s\n", name)

	if eval.WealthScore < 10 {
		desc += "\n💰 财富方面还有很大提升空间"
	}
	if eval.HealthScore < 15 {
		desc += "\n💚 健康是人生的基石，记得多去医院检查"
	}
	if eval.FameScore < 10 {
		desc += "\n⭐ 名声能带来更多机会，尝试提升社会地位"
	}
	if eval.AgeScore < 15 {
		desc += "\n⏰ 长寿是人生的终极目标，保重身体最重要"
	}
	if eval.CareerScore < 5 {
		desc += "\n💼 工作、创业和游戏成就还能进一步积累"
	}
	if eval.RelationshipScore < 5 {
		desc += "\n💖 稳定关系和共同经历也是人生的重要部分"
	}
	if eval.CollectionScore < 3 {
		desc += "\n🏡 收藏、房产和车辆可以丰富人生轨迹"
	}
	return desc
}

func (a *App) EndGame() EvaluationResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return EvaluationResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	evaluation := calculateEvaluation(a.Userinfo)
	a.Gameinfo = nil
	a.Userinfo = nil
	a.Announce = nil
	a.pendingDatingInteraction = nil
	a.stockEpoch = ""
	a.stockVersion = 0
	a.stockUpdatedAt = 0

	return EvaluationResponse{Code: 200, Msg: "游戏结束", Evaluation: &evaluation}
}
