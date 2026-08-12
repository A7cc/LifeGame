package services

import (
	"LifeGame/core"
	"path"
	"strings"
)

const (
	datingInteractionChat     = "chat"
	datingInteractionCaress   = "caress"
	datingInteractionKiss     = "kiss"
	datingInteractionOutfit   = "outfit"
	datingInteractionIntimacy = "intimacy"
)

type datingInteractionSession struct {
	DatingID int
}

type datingInteractionRule struct {
	Action        string
	Label         string
	MinStatus     string
	AffinityBonus int
	ArgumentRisk  float64
	AdultOnly     bool
	SuccessEvent  string
}

var datingInteractionRules = map[string]datingInteractionRule{
	datingInteractionChat: {
		Action: datingInteractionChat, Label: "聊天", MinStatus: core.DatingStatusStranger,
		AffinityBonus: 1, ArgumentRisk: 0.08,
		SuccessEvent: "你认真听完对方的想法，也分享了自己的心事",
	},
	datingInteractionCaress: {
		Action: datingInteractionCaress, Label: "温柔抚摸", MinStatus: core.DatingStatusAmbiguous,
		AffinityBonus: 2, ArgumentRisk: 0.12,
		SuccessEvent: "你先确认对方是否舒服，然后轻轻抚摸了对方的脸颊",
	},
	datingInteractionKiss: {
		Action: datingInteractionKiss, Label: "亲吻", MinStatus: core.DatingStatusLover,
		AffinityBonus: 3, ArgumentRisk: 0.10,
		SuccessEvent: "在确认彼此的心意后，你们交换了一个温柔的吻",
	},
	datingInteractionOutfit: {
		Action: datingInteractionOutfit, Label: "邀请换装", MinStatus: core.DatingStatusFriend,
		AffinityBonus: 1, ArgumentRisk: 0.08,
		SuccessEvent: "你们一起挑选了新的约会造型，还为彼此的审美打了分",
	},
	datingInteractionIntimacy: {
		Action: datingInteractionIntimacy, Label: "亲密过夜", MinStatus: core.DatingStatusSweetheart,
		AffinityBonus: 5, ArgumentRisk: 0.10, AdultOnly: true,
		SuccessEvent: "两位成年人再次确认了彼此的意愿和边界，共度了一个私密的夜晚",
	},
}

type datingOutfitVariant struct {
	Key       string
	Label     string
	MinStatus string
	AdultOnly bool
}

var datingOutfitVariants = []datingOutfitVariant{
	{Key: "career", Label: "职业装", MinStatus: core.DatingStatusFriend},
	{Key: "homewear", Label: "居家装", MinStatus: core.DatingStatusFriend},
	{Key: "qipao", Label: "旗袍/国风", MinStatus: core.DatingStatusAmbiguous},
	{Key: "cosplay", Label: "Cosplay", MinStatus: core.DatingStatusDating},
	{Key: "swimwear", Label: "泳装", MinStatus: core.DatingStatusDating},
	{Key: "sleepwear", Label: "睡衣", MinStatus: core.DatingStatusLover},
	{Key: "romantic", Label: "情趣睡衣", MinStatus: core.DatingStatusExclusive, AdultOnly: true},
}

var datingOutfitDirectories = map[string]bool{
	"sleepwear": true,
	"romantic":  true,
	"swimwear":  true,
	"cosplay":   true,
	"qipao":     true,
	"homewear":  true,
}

func datingOutfitAsset(image, outfit string) string {
	image = strings.TrimSpace(image)
	if image == "" || outfit == "career" || !datingOutfitDirectories[outfit] {
		return image
	}
	directory, filename := path.Split(image)
	if filename == "" {
		return image
	}
	return path.Join(directory, outfit, filename)
}

func findDatingOutfitVariant(key string) (datingOutfitVariant, bool) {
	for _, outfit := range datingOutfitVariants {
		if outfit.Key == key {
			return outfit, true
		}
	}
	return datingOutfitVariant{}, false
}

func datingStatusRank(status string) int {
	switch core.NormalizeDatingStatus(status) {
	case core.DatingStatusFriend:
		return 1
	case core.DatingStatusAmbiguous:
		return 2
	case core.DatingStatusDating:
		return 3
	case core.DatingStatusLover:
		return 4
	case core.DatingStatusExclusive:
		return 5
	case core.DatingStatusSweetheart:
		return 6
	case core.DatingStatusMarried:
		return 7
	default:
		return 0
	}
}

func datingInteractionOutcomeForRoll(rule datingInteractionRule, roll float64) (outcome string, affinityChange int) {
	if roll < rule.ArgumentRisk {
		return datingMomentArgument, -3
	}
	return rule.Action, rule.AffinityBonus
}

// DoDatingInteraction 结算成功约会短片中的一次玩家互动。
func (a *App) DoDatingInteraction(datingID int, action, outfitKey string) DatingInteractionResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingInteractionResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	action = strings.TrimSpace(strings.ToLower(action))
	rule, exists := datingInteractionRules[action]
	if !exists {
		return DatingInteractionResponse{Code: -1, Msg: "不支持该约会互动"}
	}
	if a.pendingDatingInteraction == nil || a.pendingDatingInteraction.DatingID != datingID {
		return DatingInteractionResponse{Code: -1, Msg: "当前没有待选择的成功约会互动"}
	}
	dating, exists := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID)
	if !exists {
		return DatingInteractionResponse{Code: -1, Msg: "约会对象资料不存在或不匹配"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists || core.NormalizeDatingStatus(relationship.DStatus) == core.DatingStatusFormer {
		return DatingInteractionResponse{Code: -1, Msg: "当前关系不能进行该互动"}
	}
	if datingStatusRank(relationship.DStatus) < datingStatusRank(rule.MinStatus) {
		return DatingInteractionResponse{Code: -1, Msg: "关系阶段不足，暂未解锁" + rule.Label}
	}
	if rule.AdultOnly && (a.Userinfo.UAge < 18 || dating.DAge < 18) {
		return DatingInteractionResponse{Code: -1, Msg: "亲密过夜仅限双方成年且达到爱人或已婚阶段"}
	}

	var selectedOutfit datingOutfitVariant
	if action == datingInteractionOutfit {
		outfitKey = strings.TrimSpace(strings.ToLower(outfitKey))
		var outfitExists bool
		selectedOutfit, outfitExists = findDatingOutfitVariant(outfitKey)
		if !outfitExists {
			return DatingInteractionResponse{Code: -1, Msg: "请选择有效的换装造型"}
		}
		if datingStatusRank(relationship.DStatus) < datingStatusRank(selectedOutfit.MinStatus) {
			return DatingInteractionResponse{Code: -1, Msg: "当前关系阶段还不能更换" + selectedOutfit.Label}
		}
		if selectedOutfit.AdultOnly && (a.Userinfo.UAge < 18 || dating.DAge < 18) {
			return DatingInteractionResponse{Code: -1, Msg: "情趣睡衣仅在双方成年后解锁"}
		}
	}

	// 通过所有校验后立即消耗本次互动，防止重复点击刷好感。
	a.pendingDatingInteraction = nil
	outcome, requestedChange := datingInteractionOutcomeForRoll(rule, a.datingRandomFloat())
	previousAffinity := relationship.DAffinity
	relationship.DAffinity = min(100, max(-10, relationship.DAffinity+requestedChange))
	if a.Userinfo.UMarriedDatingID == datingID {
		relationship.DStatus = core.DatingStatusMarried
	} else {
		relationship.DStatus = core.GetDatingStatus(relationship.DAffinity, relationship.DCount, hasOtherCommittedPartner(a.Userinfo, datingID))
	}
	a.Userinfo.UDating[datingID] = relationship

	response := DatingInteractionResponse{
		Code:           200,
		Interaction:    action,
		Outcome:        outcome,
		Label:          rule.Label,
		AffinityChange: relationship.DAffinity - previousAffinity,
		Datinginfo:     &relationship,
		Userinfo:       a.userSnapshot(),
	}
	if outcome == datingMomentArgument {
		response.Msg = "互动时出现了分歧，短片将播放争执场景"
		response.Event = "你们对这次互动的节奏理解不同，在争执后决定先冷静一下"
		return response
	}

	response.Msg = rule.Label + "互动成功"
	response.Event = rule.SuccessEvent
	if action == datingInteractionOutfit {
		response.Label = "更换" + selectedOutfit.Label
		response.Msg = selectedOutfit.Label + "换装成功"
		response.Event = "你们一起选定了" + selectedOutfit.Label + "，这次造型正好适合当前的关系阶段"
		response.Outfit = selectedOutfit.Label
		response.OutfitVariant = selectedOutfit.Key
		response.OutfitImage = datingOutfitAsset(dating.DImage, selectedOutfit.Key)
	}
	return response
}
