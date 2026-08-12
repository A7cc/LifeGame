package services

import (
	"LifeGame/core"
	"LifeGame/utils"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

const spouseBathCost = 200

func (a *App) datingRandomFloat() float64 {
	if a.randomRoll != nil {
		roll := a.randomRoll()
		return min(1, max(0, roll))
	}
	return rand.Float64()
}

func hasOtherCommittedPartner(user *core.User, datingID int) bool {
	if user == nil {
		return false
	}
	for id, relationship := range user.UDating {
		if id == datingID {
			continue
		}
		status := core.NormalizeDatingStatus(relationship.DStatus)
		if status == core.DatingStatusSweetheart || status == core.DatingStatusMarried {
			return true
		}
	}
	return false
}

func isCompatibleDating(user *core.User, dating core.DatingInfo) bool {
	return user != nil && user.USex != dating.DSex
}

// normalizeDatingRelationships removes relationships that conflict with the
// configured opposite-sex matching rule and enforces current relationship rules.
func normalizeDatingRelationships(user *core.User, game *core.Game) {
	if user == nil || game == nil {
		return
	}
	if len(user.UDating) == 0 {
		user.UMarriedDatingID = 0
		return
	}
	candidates := make(map[int]core.DatingInfo, len(game.GDatingInfo))
	for _, dating := range game.GDatingInfo {
		candidates[dating.DId] = dating
	}
	for datingID := range user.UDating {
		dating, exists := candidates[datingID]
		if !exists || !isCompatibleDating(user, dating) {
			delete(user.UDating, datingID)
		}
	}
	normalizeDatingState(user)
}

func normalizeDatingState(user *core.User) {
	if user == nil {
		return
	}
	if user.UMarriedDatingID != 0 {
		if relationship, exists := user.UDating[user.UMarriedDatingID]; exists {
			relationship.DStatus = core.DatingStatusMarried
			user.UDating[user.UMarriedDatingID] = relationship
		} else {
			user.UMarriedDatingID = 0
		}
	}
	for datingID, relationship := range user.UDating {
		status := core.NormalizeDatingStatus(relationship.DStatus)
		if status == core.DatingStatusMarried {
			if user.UMarriedDatingID == 0 {
				user.UMarriedDatingID = datingID
			} else if user.UMarriedDatingID != datingID {
				status = core.DatingStatusSweetheart
			}
		}
		relationship.DStatus = status
		user.UDating[datingID] = relationship
	}
	if user.UMarriedDatingID != 0 {
		return
	}
	committedID := 0
	for datingID, relationship := range user.UDating {
		if relationship.DStatus == core.DatingStatusSweetheart && (committedID == 0 || datingID < committedID) {
			committedID = datingID
		}
	}
	for datingID, relationship := range user.UDating {
		if datingID != committedID && relationship.DStatus == core.DatingStatusSweetheart {
			relationship.DStatus = core.DatingStatusExclusive
			user.UDating[datingID] = relationship
		}
	}
}

func newDatingRelationship(dating core.DatingInfo) core.UserDatingInfo {
	return core.UserDatingInfo{
		DDatingId: dating.DId,
		DName:     dating.DName,
		DStatus:   core.DatingStatusStranger,
	}
}

func datingRelationshipStatus(user *core.User, datingID int, relationship core.UserDatingInfo) string {
	if user.UMarriedDatingID == datingID {
		return core.DatingStatusMarried
	}
	status := core.NormalizeDatingStatus(relationship.DStatus)
	if status == core.DatingStatusFormer {
		return status
	}
	return core.GetDatingStatus(relationship.DAffinity, relationship.DCount, hasOtherCommittedPartner(user, datingID))
}

func findCompatibleDatingByID(user *core.User, game *core.Game, datingID int) (core.DatingInfo, bool) {
	if game == nil {
		return core.DatingInfo{}, false
	}
	for _, dating := range game.GDatingInfo {
		if dating.DId == datingID && isCompatibleDating(user, dating) {
			return dating, true
		}
	}
	return core.DatingInfo{}, false
}

// 获取约会信息列表
func (a *App) GetDatingInfo() DatingListResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingListResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if a.Userinfo.UDating == nil {
		a.Userinfo.UDating = make(map[int]core.UserDatingInfo)
	}
	normalizeDatingRelationships(a.Userinfo, a.Gameinfo)
	// 获取匹配的约会对象列表并更新解锁状态
	datingList := []core.DatingInfo{}
	meetingScenes := map[string]struct{}{}

	for _, dating := range a.Gameinfo.GDatingInfo {
		if !isCompatibleDating(a.Userinfo, dating) {
			continue
		}
		if dating.DMeetScene != "" {
			meetingScenes[dating.DMeetScene] = struct{}{}
		}
		// 获取用户约会信息（如果有）
		if userInfo, exists := a.Userinfo.UDating[dating.DId]; exists {
			dating.DAffinityLevel = datingRelationshipStatus(a.Userinfo, dating.DId, userInfo)
			userInfo.DStatus = dating.DAffinityLevel
			a.Userinfo.UDating[dating.DId] = userInfo
			// 更新解锁状态
			dating.DUnlocked = true
		} else {
			dating.DAffinityLevel = core.DatingStatusStranger
			// 配置了相遇场景的对象只能由 VisitDatingScene 解锁。
			unlocked := dating.DMeetScene == "" && core.CheckDatingUnlock(a.Userinfo, dating)
			// 如果解锁，则创建用户约会信息
			if unlocked {
				a.Userinfo.UDating[dating.DId] = newDatingRelationship(dating)
			}
			// 更新解锁状态
			dating.DUnlocked = unlocked
		}

		datingList = append(datingList, dating)
	}
	scenes := make([]string, 0, len(meetingScenes))
	for scene := range meetingScenes {
		scenes = append(scenes, scene)
	}
	sort.Strings(scenes)

	return DatingListResponse{
		Code:          200,
		Msg:           "获取成功",
		DatingList:    datingList,
		MeetingScenes: scenes,
		Userinfo:      a.userSnapshot(),
	}
}

// VisitDatingScene 主动前往场景，并只对该场景的候选人执行认识判定。
func (a *App) VisitDatingScene(scene string) DatingSceneResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingSceneResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	scene = strings.TrimSpace(scene)
	if scene == "" {
		return DatingSceneResponse{Code: -1, Msg: "请选择要前往的场景"}
	}

	validScene := false
	met := make([]core.DatingInfo, 0)
	for _, dating := range a.Gameinfo.GDatingInfo {
		if !isCompatibleDating(a.Userinfo, dating) || dating.DMeetScene != scene {
			continue
		}
		validScene = true
		if _, known := a.Userinfo.UDating[dating.DId]; known {
			continue
		}
		if core.CheckDatingUnlock(a.Userinfo, dating) {
			a.Userinfo.UDating[dating.DId] = newDatingRelationship(dating)
			dating.DUnlocked = true
			dating.DAffinityLevel = core.DatingStatusStranger
			met = append(met, dating)
		}
	}
	if !validScene {
		return DatingSceneResponse{Code: -1, Msg: "该场景没有可认识的约会对象"}
	}
	message := "这次没有遇到合适的人，可满足条件后再次前往"
	if len(met) > 0 {
		names := make([]string, 0, len(met))
		for _, dating := range met {
			names = append(names, dating.DName)
		}
		message = "你在" + scene + "认识了" + strings.Join(names, "、")
	}
	return DatingSceneResponse{
		Code:     200,
		Msg:      message,
		Scene:    scene,
		Met:      met,
		Userinfo: a.userSnapshot(),
	}
}

// 执行约会
func (a *App) DoDating(datingId int, location string) DatingResultResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingResultResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if a.Userinfo.UOpportunity.OMNum >= a.Gameinfo.GMaxHoldNum.MMRoundNum {
		return DatingResultResponse{Code: -1, Msg: "本年度约会次数已用完，请明年再来"}
	}
	dating := core.DatingInfo{}
	found := false
	// 检查该约会对象是否存在
	for _, d := range a.Gameinfo.GDatingInfo {
		if d.DId == datingId {
			dating = d
			found = true
			break
		}
	}
	if !found {
		return DatingResultResponse{Code: -1, Msg: "不存在该约会对象"}
	}
	if !isCompatibleDating(a.Userinfo, dating) {
		return DatingResultResponse{Code: -1, Msg: "该约会对象与当前角色性别不匹配"}
	}

	// 获取或创建用户约会信息，如果不存在则给用户添加
	userDating, exists := a.Userinfo.UDating[datingId]
	if !exists {
		return DatingResultResponse{Code: -1, Msg: "并不认识该约会对象"}
	}
	if core.NormalizeDatingStatus(userDating.DStatus) == core.DatingStatusFormer {
		return DatingResultResponse{Code: -1, Msg: "你们已经分手，暂时不能继续约会"}
	}
	location = strings.TrimSpace(location)
	if location == "" {
		return DatingResultResponse{Code: -1, Msg: "请选择约会地点"}
	}
	preferredLocation := false
	for _, candidate := range dating.DLocations {
		if candidate == location {
			preferredLocation = true
			break
		}
	}
	sceneRule, sceneConfigured := datingSceneForLocation(location)
	if !sceneConfigured {
		return DatingResultResponse{Code: -1, Msg: "该约会地点尚未配置场景"}
	}
	// 检查资金
	if a.Userinfo.UCash < dating.DCost {
		return DatingResultResponse{Code: -1, Msg: "资金不足"}
	}
	// 开始新约会时作废上一次未使用的短片互动。
	a.pendingDatingInteraction = nil
	// 扣除约会费用
	a.Userinfo.UCash -= dating.DCost
	a.Userinfo.UOpportunity.OMNum++

	// 喜欢的场景是高成功率、普通收益的稳妥路线；非偏好场景是
	// 低成功率、高收益的冒险路线。
	baseSuccessRate := core.CalculateDatingSuccessRate(dating, userDating)
	choiceRule := datingSceneChoiceRuleFor(baseSuccessRate, preferredLocation)
	successRate := choiceRule.SuccessRate
	preferenceRateChange := successRate - baseSuccessRate
	success := a.datingRandomFloat() < successRate

	// 计算属性变化
	// 名声
	fameChange := 0
	// 免疫力
	healthChange := 0
	// 好感度
	affinityChange := 0
	moment := datingMomentOutcome{}

	if success {
		// 成功收益由路线决定：冒险路线的基础收益和场景效果都更高。
		affinityChange = choiceRule.AffinityMin + rand.Intn(choiceRule.AffinityRange)
		fameChange = choiceRule.FameBase
		healthChange = choiceRule.HealthMin + rand.Intn(choiceRule.HealthRange)

		// 根据职业加成不同属性
		switch {
		case utils.ContainsAny(dating.DOccup, []string{"健身教练", "医生", "护士", "瑜伽教练", "奥运冠军", "花艺师", "厨师", "导游", "中医师", "心理医生"}):
			healthChange += 1
		case utils.ContainsAny(dating.DOccup, []string{"名媛千金", "上流名媛", "青年投资人", "商界精英", "艺术家", "律师", "作家", "歌手", "演员", "企业家", "博士后", "消防员", "珠宝设计师", "大学教授", "模特", "服装设计师", "科学家", "导演", "飞行员", "小提琴家", "电台主持人", "金融分析师"}):
			fameChange += 1
		}
		fameChange += sceneRule.FameBonus * choiceRule.RewardMultiplier
		healthChange += sceneRule.HealthBonus * choiceRule.RewardMultiplier
		affinityChange += sceneRule.AffinityBonus * choiceRule.RewardMultiplier

		projectedStatus := core.GetDatingStatus(
			userDating.DAffinity+affinityChange,
			userDating.DCount+1,
			hasOtherCommittedPartner(a.Userinfo, datingId),
		)
		if a.Userinfo.UMarriedDatingID == datingId {
			projectedStatus = core.DatingStatusMarried
		}
		moment = datingMomentForRoll(projectedStatus, a.Userinfo.UAge >= 18 && dating.DAge >= 18, rand.Float64())
		affinityChange += moment.AffinityChange
	} else {
		// 冒险路线失败的代价也略高。
		affinityChange = -(choiceRule.FailureMin + rand.Intn(choiceRule.FailureRange))
		healthChange = choiceRule.FailureHealth
	}

	// 更新属性
	// 好感度
	userDating.DAffinity += affinityChange
	if userDating.DAffinity < -10 {
		userDating.DAffinity = -10
	}
	if userDating.DAffinity > 100 {
		userDating.DAffinity = 100
	}
	// 约会次数
	userDating.DCount++
	// 关系状态
	if a.Userinfo.UMarriedDatingID == datingId {
		userDating.DStatus = core.DatingStatusMarried
	} else {
		userDating.DStatus = core.GetDatingStatus(userDating.DAffinity, userDating.DCount, hasOtherCommittedPartner(a.Userinfo, datingId))
	}

	// 应用属性变化
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + fameChange)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + healthChange)

	// 保存用户约会信息
	a.Userinfo.UDating[datingId] = userDating

	// 计算总资产
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	if success {
		a.pendingDatingInteraction = &datingInteractionSession{DatingID: datingId}
	}

	sceneFameChange := 0
	sceneHealthChange := 0
	sceneAffinityChange := 0
	if success {
		sceneFameChange = sceneRule.FameBonus * choiceRule.RewardMultiplier
		sceneHealthChange = sceneRule.HealthBonus * choiceRule.RewardMultiplier
		sceneAffinityChange = sceneRule.AffinityBonus * choiceRule.RewardMultiplier
	}
	sceneEffect := "偏好场景：成功率高，成功收益中等；" + sceneRule.Effect
	if !preferredLocation {
		sceneEffect = "非偏好场景：成功率低，成功时基础收益更高且场景效果翻倍"
	}

	return DatingResultResponse{
		Code:           200,
		Msg:            "约会完成",
		Success:        success,
		FameChange:     fameChange,
		HealthChange:   healthChange,
		AffinityChange: affinityChange,
		Scene: &DatingSceneOutcome{
			Location:             location,
			Category:             sceneRule.Category,
			Label:                sceneRule.Label,
			Event:                datingSceneEvent(sceneRule, location, success, preferredLocation),
			Effect:               sceneEffect,
			Preferred:            preferredLocation,
			PreferenceRateChange: preferenceRateChange,
			SuccessRate:          successRate,
			RewardTier:           choiceRule.RewardTier,
			RewardMultiplier:     choiceRule.RewardMultiplier,
			Moment:               moment.Kind,
			MomentLabel:          moment.Label,
			MomentEvent:          moment.Event,
			MomentAffinityChange: moment.AffinityChange,
			FameChange:           sceneFameChange,
			HealthChange:         sceneHealthChange,
			AffinityChange:       sceneAffinityChange,
		},
		Userinfo:   a.userSnapshot(),
		Datinginfo: &userDating,
	}
}

// MarryDating 在达到爱人门槛后建立唯一婚姻关系。
func (a *App) MarryDating(datingID int) DatingRelationshipResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingRelationshipResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if a.Userinfo.UMarriedDatingID != 0 {
		if a.Userinfo.UMarriedDatingID == datingID {
			return DatingRelationshipResponse{Code: 200, Msg: "你们已经结婚", Userinfo: a.userSnapshot()}
		}
		return DatingRelationshipResponse{Code: -1, Msg: "你已经有配偶，不能重复结婚"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists {
		return DatingRelationshipResponse{Code: -1, Msg: "并不认识该约会对象"}
	}
	if core.NormalizeDatingStatus(relationship.DStatus) == core.DatingStatusFormer {
		return DatingRelationshipResponse{Code: -1, Msg: "已经分手的关系不能直接结婚"}
	}
	if relationship.DCount < 20 || relationship.DAffinity < 90 {
		return DatingRelationshipResponse{Code: -1, Msg: "结婚需要累计约会20次且好感度达到90"}
	}
	dating, found := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID)
	if !found {
		return DatingRelationshipResponse{Code: -1, Msg: "约会对象资料不存在"}
	}
	marriageCost := max(5000, dating.DCost*5)
	if a.Userinfo.UCash < marriageCost {
		return DatingRelationshipResponse{Code: -1, Msg: fmt.Sprintf("举办婚礼需要 %d 元，现金不足", marriageCost)}
	}
	a.Userinfo.UCash -= marriageCost
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + 5)
	relationship.DStatus = core.DatingStatusMarried
	a.Userinfo.UDating[datingID] = relationship
	a.Userinfo.UMarriedDatingID = datingID
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	return DatingRelationshipResponse{
		Code:       200,
		Msg:        fmt.Sprintf("结婚成功，婚礼花费 %d 元，名声+5", marriageCost),
		Datinginfo: &relationship,
		Userinfo:   a.userSnapshot(),
	}
}

// GiveDatingGift 消耗一次年度关系行动，把礼物偏好从展示数据变成可玩的成长路线。
func (a *App) GiveDatingGift(datingID int, gift string) DatingRelationshipResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingRelationshipResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	dating, found := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID)
	if !found {
		return DatingRelationshipResponse{Code: -1, Msg: "约会对象不存在"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists || core.NormalizeDatingStatus(relationship.DStatus) == core.DatingStatusFormer {
		return DatingRelationshipResponse{Code: -1, Msg: "当前关系不能赠送礼物"}
	}
	gift = strings.TrimSpace(gift)
	rule, validGift := datingGiftRuleFor(a.Gameinfo, gift)
	if !validGift {
		return DatingRelationshipResponse{Code: -1, Msg: "该礼物不在可选目录中"}
	}
	if a.Userinfo.UOpportunity.OMNum >= a.Gameinfo.GMaxHoldNum.MMRoundNum {
		return DatingRelationshipResponse{Code: -1, Msg: "本年度关系行动次数已用完，请明年再来"}
	}
	giftCost := rule.Cost
	if a.Userinfo.UCash < giftCost {
		return DatingRelationshipResponse{Code: -1, Msg: fmt.Sprintf("准备礼物需要 %d 元，现金不足", giftCost)}
	}
	preferred := isPreferredGift(dating, gift)
	success, requestedChange, successRate, outcome := resolveDatingGift(rule, preferred, a.datingRandomFloat())
	a.Userinfo.UCash -= giftCost
	a.Userinfo.UOpportunity.OMNum++
	event := datingGiftEvent(dating, gift, preferred, success)
	previousAffinity := relationship.DAffinity
	relationship.DGiftCount++
	relationship.DAffinity = min(100, max(-10, relationship.DAffinity+requestedChange))
	affinityChange := relationship.DAffinity - previousAffinity
	if a.Userinfo.UMarriedDatingID == datingID {
		relationship.DStatus = core.DatingStatusMarried
	} else {
		relationship.DStatus = core.GetDatingStatus(relationship.DAffinity, relationship.DCount, hasOtherCommittedPartner(a.Userinfo, datingID))
	}
	a.Userinfo.UDating[datingID] = relationship
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	resultText := fmt.Sprintf("礼物被婉拒，好感度%d", affinityChange)
	if success && preferred {
		resultText = fmt.Sprintf("稳妥送礼成功，好感度+%d", affinityChange)
	} else if success {
		resultText = fmt.Sprintf("冒险送礼成功，获得高回报，好感度+%d", affinityChange)
	} else if preferred {
		resultText = fmt.Sprintf("偏好礼物这次没有送出，好感度%d", affinityChange)
	}
	return DatingRelationshipResponse{
		Code: 200, Msg: fmt.Sprintf("赠送了%s，花费%d元，%s", gift, giftCost, resultText),
		Event: event, Gift: gift, GiftCost: giftCost, AffinityChange: affinityChange,
		Outcome: outcome, Preferred: preferred, Success: success, SuccessRate: successRate,
		Datinginfo: &relationship, Userinfo: a.userSnapshot(),
	}
}

// BatheWithSpouse 是婚后专属互动，与普通约会共用年度关系行动次数。
// 互动保持非露骨表达，只结算健康与好感变化，短片由前端播放。
func (a *App) BatheWithSpouse(datingID int) SpouseInteractionResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return SpouseInteractionResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if datingID == 0 || a.Userinfo.UMarriedDatingID != datingID {
		return SpouseInteractionResponse{Code: -1, Msg: "只能和当前配偶一起洗澡"}
	}
	if _, exists := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID); !exists {
		return SpouseInteractionResponse{Code: -1, Msg: "配偶资料不存在或不匹配"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists || core.NormalizeDatingStatus(relationship.DStatus) != core.DatingStatusMarried {
		return SpouseInteractionResponse{Code: -1, Msg: "婚姻关系数据不存在"}
	}
	if a.Userinfo.UOpportunity.OMNum >= a.Gameinfo.GMaxHoldNum.MMRoundNum {
		return SpouseInteractionResponse{Code: -1, Msg: "本年度关系互动次数已用完，请明年再来"}
	}
	if a.Userinfo.UCash < spouseBathCost {
		return SpouseInteractionResponse{Code: -1, Msg: "资金不足，一起洗澡需要200元"}
	}

	previousAffinity := relationship.DAffinity
	previousHealth := a.Userinfo.UImmunity
	relationship.DAffinity = min(100, relationship.DAffinity+2)
	relationship.DStatus = core.DatingStatusMarried
	a.Userinfo.UDating[datingID] = relationship
	a.Userinfo.UCash -= spouseBathCost
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + 3)
	a.Userinfo.UOpportunity.OMNum++
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	return SpouseInteractionResponse{
		Code:           200,
		Msg:            "你们一起泡了个舒服的澡，度过了温馨时光",
		Interaction:    "bath",
		Cost:           spouseBathCost,
		AffinityChange: relationship.DAffinity - previousAffinity,
		HealthChange:   a.Userinfo.UImmunity - previousHealth,
		Datinginfo:     &relationship,
		Userinfo:       a.userSnapshot(),
	}
}

// BreakUpDating 结束未婚关系，保留历史为前任。婚姻必须通过离婚结束。
func (a *App) BreakUpDating(datingID int) DatingRelationshipResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingRelationshipResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if a.Userinfo.UMarriedDatingID == datingID {
		return DatingRelationshipResponse{Code: -1, Msg: "已婚关系请使用离婚"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists {
		return DatingRelationshipResponse{Code: -1, Msg: "并不认识该约会对象"}
	}
	status := core.NormalizeDatingStatus(relationship.DStatus)
	if status == core.DatingStatusFormer {
		return DatingRelationshipResponse{Code: -1, Msg: "你们已经分手"}
	}
	if status == core.DatingStatusStranger || status == core.DatingStatusFriend {
		return DatingRelationshipResponse{Code: -1, Msg: "当前还不是恋爱关系"}
	}
	if relationship.DAffinity >= 0 {
		relationship.DAffinity = -1
	}
	relationship.DStatus = core.DatingStatusFormer
	a.Userinfo.UDating[datingID] = relationship
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - 2)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity - 2)
	return DatingRelationshipResponse{
		Code:       200,
		Msg:        "已经分手，名声-2，免疫力-2",
		Datinginfo: &relationship,
		Userinfo:   a.userSnapshot(),
	}
}

// DivorceDating 结束当前唯一婚姻，并将关系保留为前任。
func (a *App) DivorceDating(datingID int) DatingRelationshipResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingRelationshipResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	if a.Userinfo.UMarriedDatingID != datingID {
		return DatingRelationshipResponse{Code: -1, Msg: "该对象不是你的配偶"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists {
		return DatingRelationshipResponse{Code: -1, Msg: "婚姻关系数据不存在"}
	}
	dating, _ := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID)
	divorceCost := max(3000, dating.DCost*3)
	paid := min(a.Userinfo.UCash, divorceCost)
	a.Userinfo.UCash -= paid
	relationship.DAffinity = -10
	relationship.DStatus = core.DatingStatusFormer
	a.Userinfo.UDating[datingID] = relationship
	a.Userinfo.UMarriedDatingID = 0
	a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame - 10)
	a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity - 5)
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	return DatingRelationshipResponse{
		Code:       200,
		Msg:        fmt.Sprintf("已经离婚，支付%d元，名声-10，免疫力-5", paid),
		Datinginfo: &relationship,
		Userinfo:   a.userSnapshot(),
	}
}
