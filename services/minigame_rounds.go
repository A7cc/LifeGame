package services

import (
	"LifeGame/core"
	cryptorand "crypto/rand"
	"errors"
	"math/big"
	mathrand "math/rand"
	"sort"
	"strconv"
)

var errMiniGameNotResolved = errors.New("请先完成本局二十一点操作")

var serverAuthoritativeMiniGames = map[string]bool{
	"rps":         true,
	"guess":       true,
	"dice":        true,
	"slot":        true,
	"poker":       true,
	"horseracing": true,
	"roulette":    true,
	"baccarat":    true,
	"blackjack":   true,
	"lottery":     true,
}

func secureRandomInt(max int) int {
	if max <= 1 {
		return 0
	}
	value, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(max)))
	if err == nil {
		return int(value.Int64())
	}
	// 系统随机源异常时仍保持游戏可用；math/rand 的全局源本身是并发安全的。
	return mathrand.Intn(max)
}

func (a *App) prepareAuthoritativeRound(session *core.MiniGameSession) {
	session.MGSRound = make(map[string]interface{})

	switch session.MGSName {
	case "rps":
		choices := []string{"rock", "scissors", "paper"}
		computer := choices[secureRandomInt(len(choices))]
		session.MGSRound["computerChoice"] = computer
		if (session.MGSChoice == "rock" && computer == "scissors") ||
			(session.MGSChoice == "scissors" && computer == "paper") ||
			(session.MGSChoice == "paper" && computer == "rock") {
			session.MGSOutcome = 1
		}
	case "guess":
		session.MGSSecret = secureRandomInt(100) + 1
	case "dice":
		player := secureRandomInt(6) + 1
		computer := secureRandomInt(6) + 1
		session.MGSRound["player"] = player
		session.MGSRound["computer"] = computer
		switch {
		case player > computer:
			session.MGSOutcome = 1
		case player == computer:
			session.MGSOutcome = 2
		}
	case "slot":
		symbols := []string{"🍒", "🍋", "🍊", "🍇", "⭐", "💎", "7️⃣"}
		reels := []string{
			symbols[secureRandomInt(len(symbols))],
			symbols[secureRandomInt(len(symbols))],
			symbols[secureRandomInt(len(symbols))],
		}
		session.MGSRound["reels"] = reels
		if reels[0] == reels[1] && reels[1] == reels[2] {
			outcomes := map[string]int{"💎": 6, "7️⃣": 5, "⭐": 4, "🍒": 3, "🍋": 2, "🍊": 2, "🍇": 2}
			session.MGSOutcome = outcomes[reels[0]]
		} else if reels[0] == reels[1] || reels[1] == reels[2] || reels[0] == reels[2] {
			session.MGSOutcome = 1
		}
	case "poker":
		winner := secureRandomInt(100)
		switch {
		case winner < 35:
			session.MGSOutcome = 1
			session.MGSRound["winner"] = "player"
		case winner < 70:
			session.MGSRound["winner"] = "opponent2"
		default:
			session.MGSRound["winner"] = "opponent3"
		}
	case "horseracing":
		order := weightedHorseOrder()
		session.MGSRound["finishOrder"] = order
		selected := 0
		for i, horse := range order {
			if session.MGSChoice == string(rune('0'+horse)) {
				selected = i + 1
				break
			}
		}
		session.MGSRound["rank"] = selected
		if selected == 1 {
			session.MGSOutcome = 1
		} else if selected > 1 && selected <= 3 {
			session.MGSOutcome = 2
		}
	case "roulette":
		// 0..36 加 37（代表 00），使用美式轮盘的 38 个格子。
		number := secureRandomInt(38)
		session.MGSRound["number"] = number
		if number == 37 {
			session.MGSRound["displayNumber"] = "00"
		} else {
			session.MGSRound["displayNumber"] = number
		}
		if rouletteChoiceWins(session.MGSChoice, number) {
			session.MGSOutcome = 1
		}
	case "baccarat":
		prepareBaccaratRound(session)
	case "blackjack":
		prepareBlackjackRound(session)
	case "lottery":
		if session.MGSPayout > 0 {
			session.MGSOutcome = 1
		}
	}

	if session.MGSName != "guess" && session.MGSName != "blackjack" {
		session.MGSRound["outcome"] = session.MGSOutcome
	}
}

func weightedHorseOrder() []int {
	horses := []int{1, 2, 3, 4, 5}
	weights := map[int]int{1: 400, 2: 300, 3: 150, 4: 90, 5: 60}
	order := make([]int, 0, len(horses))
	for len(horses) > 0 {
		total := 0
		for _, horse := range horses {
			total += weights[horse]
		}
		roll := secureRandomInt(total)
		selected := 0
		for i, horse := range horses {
			roll -= weights[horse]
			if roll < 0 {
				selected = i
				break
			}
		}
		order = append(order, horses[selected])
		horses = append(horses[:selected], horses[selected+1:]...)
	}
	return order
}

var rouletteRedNumbers = map[int]bool{
	1: true, 3: true, 5: true, 7: true, 9: true, 12: true, 14: true, 16: true, 18: true,
	19: true, 21: true, 23: true, 25: true, 27: true, 30: true, 32: true, 34: true, 36: true,
}

func rouletteChoiceWins(choice string, number int) bool {
	if number == 0 || number == 37 {
		return false
	}
	switch choice {
	case "red":
		return rouletteRedNumbers[number]
	case "black":
		return !rouletteRedNumbers[number]
	case "even":
		return number%2 == 0
	case "odd":
		return number%2 == 1
	case "1-18":
		return number <= 18
	case "19-36":
		return number >= 19
	default:
		return false
	}
}

var cardSuits = []string{"♠", "♥", "♦", "♣"}
var cardRanks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func shuffledDeck() []string {
	deck := make([]string, 0, len(cardSuits)*len(cardRanks))
	for _, suit := range cardSuits {
		for _, rank := range cardRanks {
			deck = append(deck, suit+rank)
		}
	}
	for i := len(deck) - 1; i > 0; i-- {
		j := secureRandomInt(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	return deck
}

func drawCard(deck *[]string) string {
	last := len(*deck) - 1
	card := (*deck)[last]
	*deck = (*deck)[:last]
	return card
}

func cardRank(card string) string {
	runes := []rune(card)
	if len(runes) < 2 {
		return ""
	}
	return string(runes[1:])
}

func baccaratValue(card string) int {
	switch cardRank(card) {
	case "A":
		return 1
	case "10", "J", "Q", "K":
		return 0
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	default:
		return 0
	}
}

func baccaratScore(cards []string) int {
	total := 0
	for _, card := range cards {
		total += baccaratValue(card)
	}
	return total % 10
}

func prepareBaccaratRound(session *core.MiniGameSession) {
	deck := shuffledDeck()
	player := []string{drawCard(&deck), drawCard(&deck)}
	banker := []string{drawCard(&deck), drawCard(&deck)}
	playerScore := baccaratScore(player)
	bankerScore := baccaratScore(banker)

	if playerScore < 8 && bankerScore < 8 {
		playerThirdValue := -1
		if playerScore <= 5 {
			card := drawCard(&deck)
			player = append(player, card)
			playerThirdValue = baccaratValue(card)
			playerScore = baccaratScore(player)
		}
		drawBanker := false
		if playerThirdValue < 0 {
			drawBanker = bankerScore <= 5
		} else {
			switch bankerScore {
			case 0, 1, 2:
				drawBanker = true
			case 3:
				drawBanker = playerThirdValue != 8
			case 4:
				drawBanker = playerThirdValue >= 2 && playerThirdValue <= 7
			case 5:
				drawBanker = playerThirdValue >= 4 && playerThirdValue <= 7
			case 6:
				drawBanker = playerThirdValue == 6 || playerThirdValue == 7
			}
		}
		if drawBanker {
			banker = append(banker, drawCard(&deck))
			bankerScore = baccaratScore(banker)
		}
	}

	winner := "tie"
	if playerScore > bankerScore {
		winner = "player"
	} else if bankerScore > playerScore {
		winner = "banker"
	}
	session.MGSRound["playerCards"] = player
	session.MGSRound["bankerCards"] = banker
	session.MGSRound["playerScore"] = playerScore
	session.MGSRound["bankerScore"] = bankerScore
	session.MGSRound["winner"] = winner
	if winner == session.MGSChoice {
		session.MGSOutcome = 1
	} else if winner == "tie" {
		session.MGSOutcome = 2
	}
}

func blackjackScore(cards []string) int {
	total := 0
	aces := 0
	for _, card := range cards {
		switch cardRank(card) {
		case "A":
			total += 11
			aces++
		case "10", "K", "Q", "J":
			total += 10
		default:
			value, _ := strconv.Atoi(cardRank(card))
			total += value
		}
	}
	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}
	return total
}

func prepareBlackjackRound(session *core.MiniGameSession) {
	session.MGSDeck = shuffledDeck()
	session.MGSPlayerCards = []string{drawCard(&session.MGSDeck), drawCard(&session.MGSDeck)}
	session.MGSDealerCards = []string{drawCard(&session.MGSDeck), drawCard(&session.MGSDeck)}
	playerNatural := blackjackScore(session.MGSPlayerCards) == 21
	dealerNatural := blackjackScore(session.MGSDealerCards) == 21
	if playerNatural || dealerNatural {
		session.MGSResolved = true
		switch {
		case playerNatural && dealerNatural:
			session.MGSOutcome = 3
		case playerNatural:
			session.MGSOutcome = 2
		default:
			session.MGSOutcome = 0
		}
	}
	updateBlackjackPublicRound(session)
}

func updateBlackjackPublicRound(session *core.MiniGameSession) {
	dealerCards := append([]string(nil), session.MGSDealerCards...)
	dealerScore := 0
	if !session.MGSResolved && len(dealerCards) > 1 {
		dealerCards[1] = "🂠"
	} else {
		dealerScore = blackjackScore(session.MGSDealerCards)
	}
	session.MGSRound = map[string]interface{}{
		"playerCards": append([]string(nil), session.MGSPlayerCards...),
		"dealerCards": dealerCards,
		"playerScore": blackjackScore(session.MGSPlayerCards),
		"dealerScore": dealerScore,
		"resolved":    session.MGSResolved,
	}
	if session.MGSResolved {
		session.MGSRound["outcome"] = session.MGSOutcome
	}
}

func resolveBlackjack(session *core.MiniGameSession) {
	for blackjackScore(session.MGSDealerCards) < 17 {
		session.MGSDealerCards = append(session.MGSDealerCards, drawCard(&session.MGSDeck))
	}
	playerScore := blackjackScore(session.MGSPlayerCards)
	dealerScore := blackjackScore(session.MGSDealerCards)
	playerNatural := len(session.MGSPlayerCards) == 2 && playerScore == 21
	switch {
	case playerScore > 21:
		session.MGSOutcome = 0
	case dealerScore > 21 || playerScore > dealerScore:
		if playerNatural {
			session.MGSOutcome = 2
		} else {
			session.MGSOutcome = 1
		}
	case playerScore == dealerScore:
		session.MGSOutcome = 3
	default:
		session.MGSOutcome = 0
	}
	session.MGSResolved = true
	updateBlackjackPublicRound(session)
}

// MiniGameAction 执行必须由后端掌握隐藏状态的游戏动作。
func (a *App) MiniGameAction(action string) H {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	if errResp := a.requireGame(); errResp != nil {
		return errResp
	}
	session := a.MiniGameSession
	if session == nil || session.MGSName != "blackjack" {
		return M{"code": -1, "msg": "当前没有可操作的二十一点游戏"}
	}
	if session.MGSResolved {
		return M{"code": -1, "msg": "本局已经结束", "round": publicMiniGameRound(session)}
	}

	switch action {
	case "hit":
		session.MGSPlayerCards = append(session.MGSPlayerCards, drawCard(&session.MGSDeck))
		if blackjackScore(session.MGSPlayerCards) > 21 {
			session.MGSOutcome = 0
			session.MGSResolved = true
		}
		updateBlackjackPublicRound(session)
	case "stand":
		resolveBlackjack(session)
	default:
		return M{"code": -1, "msg": "不支持的游戏动作"}
	}

	return M{"code": 200, "msg": "操作成功", "round": publicMiniGameRound(session)}
}

func authoritativeOutcome(session *core.MiniGameSession, submitted int) (int, error) {
	if !serverAuthoritativeMiniGames[session.MGSName] {
		return submitted, nil
	}
	if session.MGSName == "guess" {
		if submitted < 1 || submitted > 100 {
			return 0, errors.New("猜测数字必须在 1 到 100 之间")
		}
		difference := submitted - session.MGSSecret
		if difference < 0 {
			difference = -difference
		}
		session.MGSOutcome = 0
		if difference == 0 {
			session.MGSOutcome = 1
		} else if difference <= 10 {
			session.MGSOutcome = 2
		}
		session.MGSRound = map[string]interface{}{"answer": session.MGSSecret, "guess": submitted, "outcome": session.MGSOutcome}
	}
	if session.MGSName == "blackjack" && !session.MGSResolved {
		return 0, errMiniGameNotResolved
	}
	return session.MGSOutcome, nil
}

func publicMiniGameRound(session *core.MiniGameSession) map[string]interface{} {
	if session == nil || session.MGSRound == nil {
		return nil
	}
	result := make(map[string]interface{}, len(session.MGSRound))
	keys := make([]string, 0, len(session.MGSRound))
	for key := range session.MGSRound {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		switch value := session.MGSRound[key].(type) {
		case []string:
			result[key] = append([]string(nil), value...)
		case []int:
			result[key] = append([]int(nil), value...)
		default:
			result[key] = value
		}
	}
	return result
}
