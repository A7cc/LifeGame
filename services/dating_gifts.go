package services

import (
	"LifeGame/core"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

const datingGiftOptionCount = 3

const (
	preferredGiftSuccessRate = 0.88
	riskyGiftSuccessRate     = 0.15
)

type datingGiftRule struct {
	Cost            int
	PreferredEffect int
	RiskyEffect     int
}

func datingGiftCatalog(game *core.Game) []string {
	if game == nil {
		return nil
	}
	seen := make(map[string]struct{})
	for _, dating := range game.GDatingInfo {
		for _, gift := range dating.DGifts {
			gift = strings.TrimSpace(gift)
			if gift != "" {
				seen[gift] = struct{}{}
			}
		}
	}
	catalog := make([]string, 0, len(seen))
	for gift := range seen {
		catalog = append(catalog, gift)
	}
	sort.Strings(catalog)
	return catalog
}

// datingGiftRuleFor assigns every catalogued gift a stable, distinct price.
// Affinity effects scale in nine bands so the three generated choices can
// deliberately offer different risk/reward levels.
func datingGiftRuleFor(game *core.Game, gift string) (datingGiftRule, bool) {
	gift = strings.TrimSpace(gift)
	catalog := datingGiftCatalog(game)
	index := sort.SearchStrings(catalog, gift)
	if index >= len(catalog) || catalog[index] != gift {
		return datingGiftRule{}, false
	}
	band := 0
	if len(catalog) > 1 {
		band = index * 8 / (len(catalog) - 1)
	}
	return datingGiftRule{
		Cost:            200 + index*80,
		PreferredEffect: 2 + band,
		RiskyEffect:     5 + band*2,
	}, true
}

func isPreferredGift(dating core.DatingInfo, gift string) bool {
	gift = strings.TrimSpace(gift)
	for _, preferred := range dating.DGifts {
		if strings.TrimSpace(preferred) == gift {
			return true
		}
	}
	return false
}

func giftOption(game *core.Game, gift string) DatingGiftOption {
	rule, _ := datingGiftRuleFor(game, gift)
	return DatingGiftOption{
		Name:            gift,
		Cost:            rule.Cost,
		PreferredEffect: rule.PreferredEffect,
		RiskyEffect:     rule.RiskyEffect,
	}
}

func createDatingGiftOptions(game *core.Game, dating core.DatingInfo, random *rand.Rand) []DatingGiftOption {
	if game == nil || random == nil {
		return nil
	}
	preferred := make([]string, 0, len(dating.DGifts))
	preferredSet := make(map[string]struct{}, len(dating.DGifts))
	for _, gift := range dating.DGifts {
		gift = strings.TrimSpace(gift)
		if gift == "" {
			continue
		}
		if _, exists := preferredSet[gift]; exists {
			continue
		}
		if _, valid := datingGiftRuleFor(game, gift); !valid {
			continue
		}
		preferredSet[gift] = struct{}{}
		preferred = append(preferred, gift)
	}
	if len(preferred) == 0 {
		return nil
	}

	favorite := preferred[random.Intn(len(preferred))]
	favoriteRule, _ := datingGiftRuleFor(game, favorite)
	catalog := datingGiftCatalog(game)
	random.Shuffle(len(catalog), func(left, right int) {
		catalog[left], catalog[right] = catalog[right], catalog[left]
	})

	choices := []string{favorite}
	usedEffects := map[int]struct{}{favoriteRule.PreferredEffect: {}}
	for _, candidate := range catalog {
		if _, isPreferred := preferredSet[candidate]; isPreferred {
			continue
		}
		rule, _ := datingGiftRuleFor(game, candidate)
		if _, duplicateEffect := usedEffects[rule.PreferredEffect]; duplicateEffect {
			continue
		}
		choices = append(choices, candidate)
		usedEffects[rule.PreferredEffect] = struct{}{}
		if len(choices) == datingGiftOptionCount {
			break
		}
	}
	// 极端的小型测试数据可能没有三个不同效果档位，此时仍补齐不同礼物。
	if len(choices) < datingGiftOptionCount {
		for _, candidate := range catalog {
			if _, isPreferred := preferredSet[candidate]; isPreferred {
				continue
			}
			alreadyChosen := false
			for _, chosen := range choices {
				if chosen == candidate {
					alreadyChosen = true
					break
				}
			}
			if alreadyChosen {
				continue
			}
			choices = append(choices, candidate)
			if len(choices) == datingGiftOptionCount {
				break
			}
		}
	}
	random.Shuffle(len(choices), func(left, right int) {
		choices[left], choices[right] = choices[right], choices[left]
	})

	options := make([]DatingGiftOption, 0, len(choices))
	for _, gift := range choices {
		options = append(options, giftOption(game, gift))
	}
	return options
}

func resolveDatingGift(rule datingGiftRule, preferred bool, roll float64) (bool, int, float64, string) {
	if preferred {
		if roll < preferredGiftSuccessRate {
			return true, rule.PreferredEffect, preferredGiftSuccessRate, "favorite"
		}
		return false, -1, preferredGiftSuccessRate, "missed"
	}
	if roll < riskyGiftSuccessRate {
		return true, rule.RiskyEffect, riskyGiftSuccessRate, "risky-success"
	}
	return false, -2, riskyGiftSuccessRate, "rejected"
}

func datingGiftEvent(dating core.DatingInfo, gift string, preferred, success bool) string {
	if success && preferred {
		return fmt.Sprintf("%s开心地收下%s，原来你真的记住了这份喜好。", dating.DName, gift)
	}
	if success {
		return fmt.Sprintf("%s原本没有期待%s，却被这次大胆的选择意外打动。", dating.DName, gift)
	}
	if preferred {
		return fmt.Sprintf("%s虽然喜欢%s，但今天的时机不太合适，没有收下这份心意。", dating.DName, gift)
	}
	return fmt.Sprintf("%s觉得%s并不适合自己，礼貌地婉拒了这份礼物。", dating.DName, gift)
}

// GetDatingGiftOptions 每次返回三个不同价位的礼物，其中恰好一个是对象喜欢的。
func (a *App) GetDatingGiftOptions(datingID int) DatingGiftOptionsResponse {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return DatingGiftOptionsResponse{Code: -1, Msg: responseMessage(errResp)}
	}
	dating, found := findCompatibleDatingByID(a.Userinfo, a.Gameinfo, datingID)
	if !found {
		return DatingGiftOptionsResponse{Code: -1, Msg: "约会对象不存在"}
	}
	relationship, exists := a.Userinfo.UDating[datingID]
	if !exists || core.NormalizeDatingStatus(relationship.DStatus) == core.DatingStatusFormer {
		return DatingGiftOptionsResponse{Code: -1, Msg: "当前关系不能赠送礼物"}
	}
	random := rand.New(rand.NewSource(rand.Int63()))
	options := createDatingGiftOptions(a.Gameinfo, dating, random)
	if len(options) != datingGiftOptionCount {
		return DatingGiftOptionsResponse{Code: -1, Msg: "礼物选项配置不足"}
	}
	return DatingGiftOptionsResponse{Code: 200, Msg: "请选择礼物", Options: options}
}
