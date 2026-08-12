package services

import "LifeGame/core"

const (
	datingMomentChat     = "chat"
	datingMomentKiss     = "kiss"
	datingMomentIntimacy = "intimacy"
	datingMomentArgument = "argument"
)

type datingMomentOutcome struct {
	Kind           string
	Label          string
	Event          string
	AffinityChange int
}

var datingMoments = map[string]datingMomentOutcome{
	datingMomentChat: {
		Kind: datingMomentChat, Label: "真心聊天",
		Event:          "你们分享了彼此的近况，一次认真的倾听让关系更加自然",
		AffinityChange: 1,
	},
	datingMomentKiss: {
		Kind: datingMomentKiss, Label: "温柔亲吻",
		Event:          "在确认彼此的心意后，你们交换了一个温柔的吻",
		AffinityChange: 2,
	},
	datingMomentIntimacy: {
		Kind: datingMomentIntimacy, Label: "亲密过夜",
		Event:          "两位成年人再次确认彼此的意愿，共度了一个亲密而私密的夜晚",
		AffinityChange: 3,
	},
	datingMomentArgument: {
		Kind: datingMomentArgument, Label: "意见争执",
		Event:          "你们因为一个分歧发生了争执，尽管最后冷静了下来，气氛仍然有些紧张",
		AffinityChange: -4,
	},
}

// datingMomentForRoll 让关系阶段改变短片内容和风险，roll 应在 [0, 1) 内。
func datingMomentForRoll(status string, bothAdults bool, roll float64) datingMomentOutcome {
	status = core.NormalizeDatingStatus(status)
	if roll < 0 {
		roll = 0
	}
	if roll >= 1 {
		roll = 0.999999
	}

	switch status {
	case core.DatingStatusSweetheart, core.DatingStatusMarried:
		if roll < 0.30 {
			return datingMoments[datingMomentChat]
		}
		if roll < 0.65 {
			return datingMoments[datingMomentKiss]
		}
		if bothAdults && roll < 0.85 {
			return datingMoments[datingMomentIntimacy]
		}
		if !bothAdults && roll < 0.85 {
			return datingMoments[datingMomentChat]
		}
		return datingMoments[datingMomentArgument]
	case core.DatingStatusLover, core.DatingStatusExclusive:
		if roll < 0.45 {
			return datingMoments[datingMomentChat]
		}
		if roll < 0.85 {
			return datingMoments[datingMomentKiss]
		}
		return datingMoments[datingMomentArgument]
	case core.DatingStatusAmbiguous, core.DatingStatusDating:
		if roll < 0.80 {
			return datingMoments[datingMomentChat]
		}
		return datingMoments[datingMomentArgument]
	default:
		if roll < 0.90 {
			return datingMoments[datingMomentChat]
		}
		return datingMoments[datingMomentArgument]
	}
}
