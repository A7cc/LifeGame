package services

import (
	"LifeGame/core"
	"LifeGame/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type miniGameWagerRule struct {
	Min     int
	Max     int
	Step    int
	Choices map[string]float64
}

// 赔率表示获胜时返还的总金额倍数（包含下注本金）。这些规则只存在于后端，
// 前端展示的数字不会参与最终结算。
var miniGameWagerRules = map[string]miniGameWagerRule{
	"horseracing": {
		Min: 100, Max: 5000, Step: 100,
		Choices: map[string]float64{"1": 1.7, "2": 2.2, "3": 4.8, "4": 8.8, "5": 14},
	},
	"roulette": {
		Min: 100, Max: 2000, Step: 100,
		Choices: map[string]float64{"red": 2, "black": 2, "even": 2, "odd": 2, "1-18": 2, "19-36": 2},
	},
	"baccarat": {
		Min: 200, Max: 5000, Step: 200,
		Choices: map[string]float64{"player": 2, "banker": 1.95, "tie": 9},
	},
	"blackjack": {
		Min: 200, Max: 6000, Step: 200,
		Choices: map[string]float64{"": 2},
	},
}

var allowedLotteryQuantities = map[int]bool{1: true, 5: true, 10: true, 20: true}

// GetMiniGameConfigs 获取所有小游戏配置信息。
func (a *App) GetMiniGameConfigs() H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}

	games := core.GetMiniGames()
	configs := make([]map[string]interface{}, 0, len(games))
	for _, game := range games {
		difficulty, ok := game.MGDifficulty[0]
		if !ok {
			continue
		}

		rewardText := ""
		if difficulty.SMGReward["totalmoney"] > 0 {
			rewardText = fmt.Sprintf("最高%d元", difficulty.SMGReward["totalmoney"])
		} else if difficulty.SMGReward["money"] > 0 {
			rewardText = fmt.Sprintf("每次%d元", difficulty.SMGReward["money"])
		}

		_, needWager := miniGameWagerRules[game.MGName]
		config := map[string]interface{}{
			"id":        game.MGName,
			"name":      game.MGCName,
			"icon":      game.MGIcon,
			"desc":      game.MGDesc,
			"category":  game.MGType,
			"type":      game.MGType,
			"entryCost": difficulty.SMGNeed,
			"reward":    rewardText,
			"needBet":   needWager || game.MGName == "lottery",
		}
		if rule, ok := miniGameWagerRules[game.MGName]; ok {
			config["wagerMin"] = rule.Min
			config["wagerMax"] = rule.Max
			config["wagerStep"] = rule.Step
		}
		configs = append(configs, config)
	}

	return M{"code": 200, "msg": "获取成功", "configs": configs}
}

// StartMiniGame 保留基础调用形式，普通小游戏和打工均从这里开始。
func (a *App) StartMiniGame(gname string, glevel int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	return a.startMiniGame(gname, glevel, 0, "", 1)
}

// StartMiniGameWithOptions 由后端校验并扣除博彩下注或彩票总价。
func (a *App) StartMiniGameWithOptions(gname string, glevel, wager int, choice string, quantity int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	return a.startMiniGame(gname, glevel, wager, choice, quantity)
}

func (a *App) startMiniGame(gname string, glevel, wager int, choice string, quantity int) H {
	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}

	if gname == "" {
		return M{"code": -1, "msg": "小游戏名称不能为空"}
	}
	mgame, found := core.FindMiniGameConfig(gname)
	if !found {
		return M{"code": -1, "msg": "小游戏不存在"}
	}
	difficulty, ok := mgame.MGDifficulty[glevel]
	if !ok {
		return M{"code": -1, "msg": "难度不存在"}
	}

	cost, tickets, payout, err := prepareMiniGameEconomy(mgame, difficulty, wager, choice, quantity)
	if err != nil {
		return M{"code": -1, "msg": err.Error()}
	}
	if a.Userinfo.UCash < cost {
		return M{"code": -1, "msg": fmt.Sprintf("资金不足，需要 %d 元", cost)}
	}
	// 只有新一局通过全部校验后才替换旧会话，避免无效请求抹掉待结算游戏。
	a.MiniGameSession = nil

	switch mgame.MGType {
	case "front", "back", "casual", "board", "competitive", "gambling":
		if a.Userinfo.UOpportunity.OGNum >= a.Gameinfo.GMaxHoldNum.MGRoundNum {
			return M{"code": -1, "msg": "小游戏次数已用完，请明年再来"}
		}
		a.Userinfo.UCash -= cost
		a.Userinfo.UOpportunity.OGNum++
	case "work":
		if a.Userinfo.UOpportunity.OWNum >= a.Gameinfo.GMaxHoldNum.MWRoundNum {
			return M{"code": -1, "msg": "您当前已超出工作次数，请休息一下！"}
		}
		a.Userinfo.UOpportunity.OWNum++
		a.recordBankTaskStats("work", 1)
	default:
		return M{"code": -1, "msg": "小游戏类型错误"}
	}

	if a.Userinfo.UMiniGameRecords == nil {
		a.Userinfo.UMiniGameRecords = make(map[string]core.MiniGameRecord)
	}
	record := a.Userinfo.UMiniGameRecords[gname]
	record.MGRType = mgame.MGType
	record.PlayCount++
	a.Userinfo.UMiniGameRecords[gname] = record

	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	session := &core.MiniGameSession{
		MGSID:        newRuntimeID(),
		MGSName:      gname,
		MGSSubInfo:   difficulty,
		MGSLevel:     glevel,
		MGSType:      mgame.MGType,
		MGSStartTime: time.Now().Unix(),
		MGSWager:     wager,
		MGSChoice:    choice,
		MGSQuantity:  quantity,
		MGSCost:      cost,
		MGSPayout:    payout,
		MGSTickets:   tickets,
	}
	a.prepareAuthoritativeRound(session)
	a.MiniGameSession = session

	response := M{
		"code":      200,
		"msg":       "开始",
		"cost":      cost,
		"sessionid": session.MGSID,
		"userinfo":  a.userSnapshot(),
	}
	round := publicMiniGameRound(session)
	if len(tickets) > 0 {
		if round == nil {
			round = make(map[string]interface{})
		}
		round["tickets"] = tickets
	}
	if len(round) > 0 {
		response["round"] = round
	}
	return response
}

// CancelMiniGame 显式放弃尚未结算的小游戏。已扣除的报名费和下注不退还，
// 避免玩家看到随机局面后反复取消。sessionID 不匹配时按幂等成功处理，
// 从而保证旧窗口的延迟请求不会取消刚开始的新游戏。
func (a *App) CancelMiniGame(sessionID string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	session := a.MiniGameSession
	if session == nil || sessionID == "" || session.MGSID != sessionID {
		return M{"code": 200, "msg": "没有需要取消的小游戏", "cancelled": false}
	}

	a.MiniGameSession = nil
	return M{
		"code":      200,
		"msg":       "小游戏已取消，已支付费用不予退还",
		"cancelled": true,
		"cost":      session.MGSCost,
	}
}

func prepareMiniGameEconomy(mgame core.MiniGame, difficulty core.SubMiniGame, wager int, choice string, quantity int) (int, []int, int, error) {
	if mgame.MGName == "lottery" {
		if wager != 0 || choice != "" || !allowedLotteryQuantities[quantity] {
			return 0, nil, 0, fmt.Errorf("彩票购买数量无效")
		}
		tickets := generateLotteryTickets(quantity)
		payout := 0
		for _, prize := range tickets {
			payout += prize
		}
		return difficulty.SMGNeed * quantity, tickets, payout, nil
	}
	if mgame.MGName == "rps" {
		if wager != 0 || quantity != 1 || (choice != "rock" && choice != "paper" && choice != "scissors") {
			return 0, nil, 0, fmt.Errorf("猜拳选项无效")
		}
		return difficulty.SMGNeed, nil, 0, nil
	}

	if rule, ok := miniGameWagerRules[mgame.MGName]; ok {
		if wager < rule.Min || wager > rule.Max || wager%rule.Step != 0 {
			return 0, nil, 0, fmt.Errorf("下注金额必须在 %d 至 %d 元之间，且为 %d 的倍数", rule.Min, rule.Max, rule.Step)
		}
		if _, ok := rule.Choices[choice]; !ok {
			return 0, nil, 0, fmt.Errorf("下注选项无效")
		}
		if quantity != 1 {
			return 0, nil, 0, fmt.Errorf("购买数量无效")
		}
		return difficulty.SMGNeed + wager, nil, 0, nil
	}

	if wager != 0 || choice != "" || quantity != 1 {
		return 0, nil, 0, fmt.Errorf("该小游戏不支持下注参数")
	}
	return difficulty.SMGNeed, nil, 0, nil
}

func generateLotteryTickets(quantity int) []int {
	tickets := make([]int, quantity)
	for i := range tickets {
		roll := float64(secureRandomInt(10000)) / 100
		switch {
		case roll < 50:
			tickets[i] = 0
		case roll < 80:
			tickets[i] = 100
		case roll < 95:
			tickets[i] = 200
		case roll < 99:
			tickets[i] = 500
		case roll < 99.9:
			tickets[i] = 1000
		case roll < 99.99:
			tickets[i] = 5000
		default:
			tickets[i] = 15000
		}
	}
	return tickets
}

// AddMiniGameWager 支持二十一点加倍，并立即由后端扣除新增下注。
func (a *App) AddMiniGameWager(amount int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	session := a.MiniGameSession
	if session == nil || session.MGSName != "blackjack" {
		return M{"code": -1, "msg": "当前游戏不支持追加下注"}
	}
	if session.MGSResolved || len(session.MGSPlayerCards) != 2 {
		return M{"code": -1, "msg": "当前牌局不能加倍"}
	}
	rule := miniGameWagerRules["blackjack"]
	newWager := session.MGSWager + amount
	if amount <= 0 || newWager > rule.Max || amount%rule.Step != 0 {
		return M{"code": -1, "msg": "追加下注金额无效"}
	}
	if a.Userinfo.UCash < amount {
		return M{"code": -1, "msg": "资金不足，无法加倍"}
	}
	a.Userinfo.UCash -= amount
	session.MGSWager = newWager
	session.MGSCost += amount
	// 二十一点加倍只能再拿一张牌，然后由后端自动完成庄家回合。
	session.MGSPlayerCards = append(session.MGSPlayerCards, drawCard(&session.MGSDeck))
	resolveBlackjack(session)
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	return M{
		"code":     200,
		"msg":      "加倍成功",
		"wager":    newWager,
		"round":    publicMiniGameRound(session),
		"userinfo": a.userSnapshot(),
	}
}

// EndMiniGame 的参数只作为技能表现或玩家输入使用；纯随机游戏的结果由后端会话决定。
func (a *App) EndMiniGame(submitted int) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	session := a.MiniGameSession
	if session == nil {
		return M{"code": -1, "msg": "小游戏未开始", "userinfo": a.userSnapshot()}
	}

	mgame, found := core.FindMiniGameConfig(session.MGSName)
	if !found {
		return M{"code": -1, "msg": "小游戏不存在", "userinfo": a.userSnapshot()}
	}
	// 时长校验只用于玩家主动提交的获胜/得分结果。0 代表认输或失败，
	// 应允许立即结束，否则关闭棋盘也会留下无法正常结算的会话。
	if submitted > 0 && !serverAuthoritativeMiniGames[session.MGSName] {
		if err := core.CheckAndClearExpiredSession(session); err != nil {
			a.MiniGameSession = nil
			return M{"code": -1, "msg": err.Error(), "userinfo": a.userSnapshot()}
		}
	}
	outcome, err := authoritativeOutcome(session, submitted)
	if err != nil {
		return M{"code": -1, "msg": err.Error(), "userinfo": a.userSnapshot()}
	}
	a.MiniGameSession = nil

	if response, handled := a.settleVariableCasualGame(session, outcome); handled {
		return enrichMiniGameSettlement(response, session)
	}
	if response, handled := a.settleWageredMiniGame(session, outcome); handled {
		return enrichMiniGameSettlement(response, session)
	}

	difficulty := session.MGSSubInfo
	wintype := "bool"
	if difficulty.SMGTarget != 0 {
		wintype = "total"
	}
	if difficulty.SMGTarget != 0 && difficulty.SMGReward["money"] != 0 {
		wintype = "count"
	}

	switch wintype {
	case "bool":
		if outcome != 1 {
			return enrichMiniGameSettlement(M{"code": -1, "msg": "你输了！", "win": false, "userinfo": a.userSnapshot()}, session)
		}
	case "total":
		if outcome < difficulty.SMGTarget {
			return M{"code": -1, "msg": "你未达到目标分数！", "win": false, "userinfo": a.userSnapshot()}
		}
	case "count":
		if outcome <= 0 || outcome > difficulty.SMGTarget {
			return M{"code": -1, "msg": "提交的完成次数无效", "win": false, "userinfo": a.userSnapshot()}
		}
	default:
		return M{"code": -1, "msg": "你的结果类型错误！", "userinfo": a.userSnapshot()}
	}

	if utils.RandomInRange(0, 100) < 5 && mgame.MGType == "work" {
		return M{"code": -1, "msg": "老板跑路了！", "win": false, "userinfo": a.userSnapshot()}
	}

	rewardInfo := make([]string, 0, len(difficulty.SMGReward))
	payout := 0
	for key, value := range difficulty.SMGReward {
		switch key {
		case "totalmoney":
			a.Userinfo.UCash += value
			payout += value
			rewardInfo = append(rewardInfo, fmt.Sprintf("获得 %d 元", value))
		case "money":
			money := value * outcome
			maximum := difficulty.SMGTarget * value
			if money > maximum {
				money = maximum
			}
			a.Userinfo.UCash += money
			payout += money
			rewardInfo = append(rewardInfo, fmt.Sprintf("获得 %d 元", money))
		case "fame":
			a.Userinfo.UFame = core.CalcFame(a.Userinfo.UFame + value)
			rewardInfo = append(rewardInfo, fmt.Sprintf("获得 %d 点声望", value))
		case "immunity":
			a.Userinfo.UImmunity = core.CalcImmunity(a.Userinfo.UImmunity + value)
			rewardInfo = append(rewardInfo, fmt.Sprintf("获得 %d 点免疫力", value))
		default:
			rewardInfo = append(rewardInfo, fmt.Sprintf("出现未知奖励类型：%s", key))
		}
	}
	if len(rewardInfo) == 0 {
		return M{"code": -1, "msg": "没有可发放的奖励", "userinfo": a.userSnapshot()}
	}

	a.recordMiniGameWin(mgame)
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	return enrichMiniGameSettlement(M{
		"code":      200,
		"msg":       strings.Join(rewardInfo, "，"),
		"win":       true,
		"payout":    payout,
		"netchange": payout - session.MGSCost,
		"userinfo":  a.userSnapshot(),
	}, session)
}

func enrichMiniGameSettlement(response H, session *core.MiniGameSession) H {
	result, ok := response.(M)
	if !ok || !serverAuthoritativeMiniGames[session.MGSName] {
		return response
	}
	result["outcome"] = session.MGSOutcome
	if round := publicMiniGameRound(session); len(round) > 0 {
		result["round"] = round
	}
	return result
}

// 少数休闲游戏有多个奖励档位。前端只提交档位编号，具体金额仍只由后端决定。
func (a *App) settleVariableCasualGame(session *core.MiniGameSession, outcome int) (H, bool) {
	payout := 0
	valid := true

	switch session.MGSName {
	case "guess":
		switch outcome {
		case 0:
		case 1:
			payout = session.MGSSubInfo.SMGReward["totalmoney"]
		case 2:
			payout = 30
		default:
			valid = false
		}
	case "dice":
		switch outcome {
		case 0:
		case 1:
			payout = session.MGSSubInfo.SMGReward["totalmoney"]
		case 2:
			payout = session.MGSSubInfo.SMGNeed
		default:
			valid = false
		}
	case "slot":
		payoutByOutcome := map[int]int{0: 0, 1: 300, 2: 2000, 3: 3000, 4: 4000, 5: 6000, 6: 8000}
		var ok bool
		payout, ok = payoutByOutcome[outcome]
		valid = ok
	default:
		return nil, false
	}

	if !valid {
		return M{"code": -1, "msg": "小游戏奖励档位无效", "userinfo": a.userSnapshot()}, true
	}
	if payout > 0 {
		a.Userinfo.UCash += payout
		mgame, _ := core.FindMiniGameConfig(session.MGSName)
		a.recordMiniGameWin(mgame)
	}
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)
	msg := "本局未获得奖励"
	code := -1
	if payout > 0 {
		msg = fmt.Sprintf("后端结算返还 %d 元", payout)
		code = 200
	}
	return M{
		"code":      code,
		"msg":       msg,
		"win":       payout > 0,
		"payout":    payout,
		"cost":      session.MGSCost,
		"netchange": payout - session.MGSCost,
		"userinfo":  a.userSnapshot(),
	}, true
}

func (a *App) settleWageredMiniGame(session *core.MiniGameSession, outcome int) (H, bool) {
	payout := 0
	win := false
	valid := true

	switch session.MGSName {
	case "lottery":
		payout = session.MGSPayout
		win = payout > 0
	case "horseracing":
		rule := miniGameWagerRules[session.MGSName]
		switch outcome {
		case 0:
		case 1:
			payout = int(math.Round(float64(session.MGSWager) * rule.Choices[session.MGSChoice]))
			win = true
		case 2:
			payout = session.MGSWager / 2
		default:
			valid = false
		}
	case "roulette":
		if outcome == 1 {
			rule := miniGameWagerRules[session.MGSName]
			payout = int(math.Round(float64(session.MGSWager) * rule.Choices[session.MGSChoice]))
			win = true
		} else if outcome != 0 {
			valid = false
		}
	case "baccarat":
		switch outcome {
		case 0:
		case 1:
			rule := miniGameWagerRules[session.MGSName]
			payout = int(math.Round(float64(session.MGSWager) * rule.Choices[session.MGSChoice]))
			win = true
		case 2:
			payout = session.MGSWager
		default:
			valid = false
		}
	case "blackjack":
		switch outcome {
		case 0:
		case 1:
			payout = session.MGSWager * 2
			win = true
		case 2:
			payout = int(math.Round(float64(session.MGSWager) * 2.5))
			win = true
		case 3:
			payout = session.MGSWager
		default:
			valid = false
		}
	default:
		return nil, false
	}

	if !valid {
		return M{"code": -1, "msg": "小游戏结果无效", "userinfo": a.userSnapshot()}, true
	}
	if payout > 0 {
		a.Userinfo.UCash += payout
	}
	if win {
		mgame, _ := core.FindMiniGameConfig(session.MGSName)
		a.recordMiniGameWin(mgame)
	}
	a.Userinfo.UAssets = core.CalculateUserAssets(a.Userinfo, a.Gameinfo)

	netChange := payout - session.MGSCost
	msg := "本局未获得奖励"
	if payout > 0 {
		msg = "后端结算返还 " + strconv.Itoa(payout) + " 元"
	}
	code := -1
	if payout > 0 {
		code = 200
	}
	return M{
		"code":      code,
		"msg":       msg,
		"win":       win,
		"payout":    payout,
		"wager":     session.MGSWager,
		"cost":      session.MGSCost,
		"netchange": netChange,
		"userinfo":  a.userSnapshot(),
	}, true
}

func (a *App) recordMiniGameWin(mgame core.MiniGame) {
	if a.Userinfo.UMiniGameRecords == nil {
		a.Userinfo.UMiniGameRecords = make(map[string]core.MiniGameRecord)
	}
	record := a.Userinfo.UMiniGameRecords[mgame.MGName]
	record.MGRType = mgame.MGType
	record.WinCount++
	a.Userinfo.UMiniGameRecords[mgame.MGName] = record
}
