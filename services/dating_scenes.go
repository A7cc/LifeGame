package services

import (
	"fmt"
	"math/rand"
)

type datingSceneRule struct {
	Category      string
	Label         string
	Effect        string
	FameBonus     int
	HealthBonus   int
	AffinityBonus int
	SuccessEvents []string
	FailureEvents []string
}

type datingSceneChoiceRule struct {
	SuccessRate      float64
	RewardMultiplier int
	AffinityMin      int
	AffinityRange    int
	FameBase         int
	HealthMin        int
	HealthRange      int
	FailureMin       int
	FailureRange     int
	FailureHealth    int
	RewardTier       string
}

// datingSceneChoiceRuleFor turns location preference into a clear
// steady-versus-high-risk choice while retaining a small influence from the
// relationship's base success rate.
func datingSceneChoiceRuleFor(baseRate float64, preferred bool) datingSceneChoiceRule {
	if preferred {
		return datingSceneChoiceRule{
			SuccessRate:      min(0.95, max(0.82, baseRate+0.25)),
			RewardMultiplier: 1,
			AffinityMin:      4,
			AffinityRange:    4,
			FameBase:         1,
			HealthMin:        1,
			HealthRange:      2,
			FailureMin:       1,
			FailureRange:     1,
			FailureHealth:    -1,
			RewardTier:       "steady",
		}
	}
	return datingSceneChoiceRule{
		SuccessRate:      min(0.20, max(0.08, baseRate-0.50)),
		RewardMultiplier: 2,
		AffinityMin:      10,
		AffinityRange:    6,
		FameBase:         2,
		HealthMin:        3,
		HealthRange:      3,
		FailureMin:       2,
		FailureRange:     3,
		FailureHealth:    -2,
		RewardTier:       "high-risk",
	}
}

var datingSceneRules = map[string]datingSceneRule{
	"dining": {
		Category: "dining", Label: "餐饮时光", Effect: "成功时好感额外+2",
		AffinityBonus: 2,
		SuccessEvents: []string{
			"你们在%s分享了喜欢的味道，聊天也越来越自然",
			"%s的灯光和美食让这次约会格外温暖",
		},
		FailureEvents: []string{
			"%s有些嘈杂，你们没能找到舒服的聊天节奏",
			"这次在%s的口味没有合拍，但也更了解了彼此",
		},
	},
	"culture": {
		Category: "culture", Label: "艺术约会", Effect: "成功时名声额外+2",
		FameBonus: 2,
		SuccessEvents: []string{
			"你们在%s发现了相近的审美，交流变得格外投契",
			"%s的艺术氛围让你们留下了一段有质感的回忆",
		},
		FailureEvents: []string{
			"你们对%s的内容看法不同，气氛一时有些拘谨",
			"%s的安排不太合拍，这次交流没有完全展开",
		},
	},
	"nature": {
		Category: "nature", Label: "自然漫游", Effect: "成功时健康额外+2",
		HealthBonus: 2,
		SuccessEvents: []string{
			"你们在%s并肩散步，风景让彼此都放松了下来",
			"%s的清新空气和晚霞让这次约会轻松又浪漫",
		},
		FailureEvents: []string{
			"%s的天气不太配合，你们只好提前结束行程",
			"这次%s之行有些匆忙，没能好好享受风景",
		},
	},
	"activity": {
		Category: "activity", Label: "活力同行", Effect: "成功时健康和好感额外+1",
		HealthBonus: 1, AffinityBonus: 1,
		SuccessEvents: []string{
			"你们在%s配合默契，笑声让距离迅速拉近",
			"%s的新鲜体验激发了共同话题，你们玩得很尽兴",
		},
		FailureEvents: []string{
			"%s的活动节奏有点不合拍，你们都有些手忙脚乱",
			"这次在%s发挥不佳，但至少留下了有趣的插曲",
		},
	},
	"study": {
		Category: "study", Label: "知性相伴", Effect: "成功时名声和好感额外+1",
		FameBonus: 1, AffinityBonus: 1,
		SuccessEvents: []string{
			"你们在%s认真交流，彼此的想法带来了新的启发",
			"%s安静的氛围让你们谈到了更深入的话题",
		},
		FailureEvents: []string{
			"%s的话题略显严肃，你们的交流没有预想中轻松",
			"这次在%s有些冷场，需要更多时间了解彼此",
		},
	},
	"wellness": {
		Category: "wellness", Label: "疗愈放松", Effect: "成功时健康额外+3",
		HealthBonus: 3,
		SuccessEvents: []string{
			"你们在%s放慢脚步，安静陪伴也显得格外亲密",
			"%s的舒缓体验让身心放松，彼此相处更加自在",
		},
		FailureEvents: []string{
			"你们还不太习惯%s的安静节奏，气氛稍显尴尬",
			"这次%s体验没有完全放松下来，下次可以换个方式",
		},
	},
	"luxury": {
		Category: "luxury", Label: "浪漫礼遇", Effect: "成功时名声+1、好感额外+2",
		FameBonus: 1, AffinityBonus: 2,
		SuccessEvents: []string{
			"你在%s的精心安排让对方感受到满满的重视",
			"%s的浪漫氛围让这次约会成为难忘的特别时刻",
		},
		FailureEvents: []string{
			"%s虽然精致，但安排显得有些刻意，彼此没有放松下来",
			"这次%s行程期待过高，实际相处反而有些拘束",
		},
	},
}

var datingLocationCategories = makeDatingLocationCategories(map[string][]string{
	"dining": {
		"餐厅", "茶馆", "厨房", "咖啡馆", "咖啡厅", "美食节", "奶茶店",
		"商务餐厅", "烧烤", "西餐厅", "夜市",
	},
	"culture": {
		"博物馆", "电台", "电影节", "电影院", "歌剧院", "画廊", "画展", "剧场",
		"剧院", "录音棚", "排练厅", "片场", "琴房", "摄影棚", "首映礼",
		"舞蹈室", "舞台", "艺术馆", "艺术区", "音乐会", "音乐节", "音乐厅",
		"影院", "展览", "KTV", "T台", "时装周",
	},
	"nature": {
		"茶园", "公园", "观景台", "海边", "海岛", "户外", "花店", "花市", "花园",
		"景点", "街拍", "森林", "医院花园", "植物园",
	},
	"activity": {
		"超市", "电竞馆", "赌场", "高尔夫", "家居城", "健身房", "街区", "俱乐部",
		"马戏团", "买手店", "赛事", "商场", "市集", "体育馆", "网红地", "训练场",
		"夜场", "游乐园", "游戏展", "展会", "直播间",
	},
	"study": {
		"大学", "工作室", "讲座", "科技馆", "科技园", "律所", "实验室", "书店",
		"书房", "图书馆", "校园", "医院", "幼儿园", "诊所", "证券所",
	},
	"wellness": {"静修中心", "温泉", "养生馆", "瑜伽馆", "中医馆"},
	"luxury": {
		"度假村", "豪华酒店", "会所", "婚礼场地", "机场", "酒店", "酒会", "拍卖会",
		"沙龙", "时尚派对", "私人会所", "私人庄园", "游艇", "珠宝店", "VIP厅",
	},
})

func makeDatingLocationCategories(groups map[string][]string) map[string]string {
	locations := make(map[string]string)
	for category, names := range groups {
		for _, name := range names {
			locations[name] = category
		}
	}
	return locations
}

func datingSceneForLocation(location string) (datingSceneRule, bool) {
	category, exists := datingLocationCategories[location]
	if !exists {
		return datingSceneRule{}, false
	}
	return datingSceneRules[category], true
}

func datingSceneEvent(rule datingSceneRule, location string, success bool, preferred bool) string {
	events := rule.FailureEvents
	if success {
		events = rule.SuccessEvents
	}
	prefix := fmt.Sprintf("你选择的%s不是对方特别喜欢的环境，这次冒险没有成功。", location)
	if success {
		prefix = fmt.Sprintf("你选择的%s不是对方特别喜欢的环境，但这次冒险带来了意外惊喜！", location)
	}
	if preferred {
		prefix = fmt.Sprintf("你选中了对方喜欢的%s，稳妥的安排让约会更容易成功。", location)
	}
	if len(events) == 0 {
		return prefix + fmt.Sprintf("你们在%s度过了一段时光", location)
	}
	return prefix + fmt.Sprintf(events[rand.Intn(len(events))], location)
}
