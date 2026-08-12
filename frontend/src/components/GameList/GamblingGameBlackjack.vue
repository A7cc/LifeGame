<template>
  <el-dialog v-model="visible" title="二十一点" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">追求21点，超越庄家而不爆牌</div>
        <div class="warning-box">
          <span class="warning-icon">⚠️</span>
          <span class="warning-text">庄家优势0.5%-2%，看似公平实则必输！</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 800 }}元</span>
          <span class="info-item">🏆 奖励：最高4000元</span>
        </div>

        <!-- 下注金额 -->
        <div class="bet-section">
          <div class="bet-label">下注金额：💰 {{ betAmount }} 元</div>
          <el-slider v-model="betAmount" :min="200" :max="3000" :step="200" show-input></el-slider>
        </div>

        <el-button type="primary" @click="handleStartGame" style="width: 100%">开始发牌</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <div class="blackjack-table">
          <!-- 庄家区域 -->
          <div class="hand-area dealer-area">
            <div class="area-title">🎰 庄家</div>
            <div class="cards-area">
              <span v-for="(card, index) in dealerCards" :key="index" class="bj-card" :class="{ hidden: index === 1 && !showDealerCard }">
                {{ (index === 1 && !showDealerCard) ? '🂠' : card }}
              </span>
            </div>
            <div class="hand-score" v-if="showDealerCard">
              点数：{{ dealerScore }}
            </div>
          </div>

          <!-- 游戏状态 -->
          <div class="game-center">
            <div class="status-text">{{ gameStatus }}</div>
          </div>

          <!-- 玩家区域 -->
          <div class="hand-area player-area">
            <div class="area-title">👤 你</div>
            <div class="cards-area">
              <span v-for="(card, index) in playerCards" :key="index" class="bj-card">
                {{ card }}
              </span>
            </div>
            <div class="hand-score">
              点数：{{ playerScore }}
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons" v-if="!playerBusted">
          <el-button @click="hit" :disabled="processing" type="primary">要牌</el-button>
          <el-button @click="stand" :disabled="processing">停牌</el-button>
          <el-button @click="doubleDown" :disabled="processing || playerCards.length !== 2" type="warning">加倍</el-button>
        </div>

        <div v-if="playerBusted" class="bust-message">💥 爆牌了！</div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
        <div class="result-stats">
          <div class="stats-row">你的牌：{{ playerCards.join(' ') }} = {{ playerScore }}</div>
          <div class="stats-row">庄家牌：{{ dealerCards.join(' ') }} = {{ dealerScore }}</div>
          <div class="stats-row result-winnings" :class="{ win: totalWinnings > 0 }">
            {{ totalWinnings >= 0 ? '赢取' : '输掉' }}：💰 {{ Math.abs(totalWinnings) }}
          </div>
        </div>
        <div v-if="consecutiveLosses >= 3" class="loss-warning">
          已连续输 {{ consecutiveLosses }} 局！这就是赌博陷阱！
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'blackjack', name: '二十一点', entryCost: 800, needBet: true })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, addWager, action, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const betAmount = ref(800)
const gameStatus = ref('')
const showDealerCard = ref(false)
const playerBusted = ref(false)
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultReward = ref('')
const totalWinnings = ref(0)
const consecutiveLosses = ref(0)
const roundsPlayed = ref(0)

const playerCards = ref([])
const dealerCards = ref([])

const resetGameState = () => {
  reset()
  betAmount.value = 800
  gameStatus.value = ''
  showDealerCard.value = false
  playerBusted.value = false
  playerCards.value = []
  dealerCards.value = []
  consecutiveLosses.value = 0
  roundsPlayed.value = 0
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const playerScore = ref(0)
const dealerScore = ref(0)

const applyRound = (round) => {
	if (!round) return
	playerCards.value = [...(round.playerCards || [])]
	dealerCards.value = [...(round.dealerCards || [])]
	playerScore.value = Number(round.playerScore || 0)
	dealerScore.value = Number(round.dealerScore || 0)
	showDealerCard.value = Boolean(round.resolved)
	playerBusted.value = playerScore.value > 21
}

const handleStartGame = async () => {
  const success = await startGame({ wager: betAmount.value })
  if (!success) return

  roundsPlayed.value++
  playerBusted.value = false
  showDealerCard.value = false

	applyRound(startData.value?.round)

	gameStatus.value = '选择：要牌 or 停牌？'
	if (startData.value?.round?.resolved) {
		setTimeout(() => endRound(), 500)
	}
}

const hit = async () => {
	processing.value = true
	const response = await action('hit')
	processing.value = false
	if (!response) return
	applyRound(response.round)
	if (response.round?.resolved) {
		gameStatus.value = '爆牌！'
		setTimeout(() => endRound(), 500)
	} else {
		gameStatus.value = playerScore.value === 21 ? '21点！请选择停牌' : '选择：要牌 or 停牌？'
	}
}

const stand = async () => {
	processing.value = true
	gameStatus.value = '庄家回合...'
	const response = await action('stand')
	processing.value = false
	if (!response) return
	applyRound(response.round)
	await delay(500)
	endRound()
}

const doubleDown = async () => {
  const doubleAmount = betAmount.value
	const response = await addWager(doubleAmount)
	if (!response) return
	betAmount.value += doubleAmount
	applyRound(response.round)
	gameStatus.value = '加倍完成，庄家结算中...'
	setTimeout(() => endRound(), 500)
}

const endRound = async () => {
	const detail = { playerCards: playerCards.value, dealerCards: dealerCards.value, playerScore: playerScore.value, dealerScore: dealerScore.value, betAmount: betAmount.value }
	const gameResult = await endGame(999, null, detail)
	if (!gameResult) return
	applyRound(gameResult.round)

	switch (gameResult.outcome) {
		case 2:
			resultIcon.value = '🃏'
			resultTitle.value = 'Blackjack！'
			consecutiveLosses.value = 0
			break
		case 1:
			resultIcon.value = '🎉'
			resultTitle.value = dealerScore.value > 21 ? '庄家爆牌！你赢了！' : '你赢了！'
			consecutiveLosses.value = 0
			break
		case 3:
			resultIcon.value = '🤝'
			resultTitle.value = '平局'
			consecutiveLosses.value = 0
			break
		default:
			resultIcon.value = playerScore.value > 21 ? '💥' : '😢'
			resultTitle.value = playerScore.value > 21 ? '爆牌！你输了！' : '庄家赢！'
			consecutiveLosses.value++
	}

	totalWinnings.value = gameResult.netChange
	resultReward.value = `${gameResult.resultText}（本局净变化 ${gameResult.netChange} 元）`

  if (consecutiveLosses.value >= 3) {
    resultDetail.value = `连续输 ${consecutiveLosses.value} 局。看似公平的游戏，庄家总有优势！`
  } else {
    resultDetail.value = ''
  }

	emit('complete', gameResult)
}

const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))
</script>

<style scoped>
.game-dialog {
  padding: 10px 0;
}

.game-intro {
  text-align: center;
}

.intro-text {
  font-size: 14px;
  color: var(--font-secondary);
  margin-bottom: 15px;
}

.warning-box {
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.warning-icon {
  font-size: 18px;
}

.warning-text {
  font-size: 13px;
  color: #c00;
  font-weight: 600;
}

.game-info {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.info-item {
  font-size: 13px;
  padding: 6px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.bet-section {
  background: var(--panel-color);
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.bet-label {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--font-color);
}

.game-playing {
  text-align: center;
}

.blackjack-table {
  background: #1a472a;
  border-radius: 12px;
  padding: 20px;
  border: 6px solid #8b4513;
  margin-bottom: 15px;
}

.hand-area {
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 10px;
}

.dealer-area {
  background: rgba(255, 100, 100, 0.1);
  border: 1px solid rgba(255, 100, 100, 0.3);
}

.player-area {
  background: rgba(100, 100, 255, 0.1);
  border: 1px solid rgba(100, 100, 255, 0.3);
}

.area-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #fff;
}

.cards-area {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 6px;
}

.bj-card {
  font-size: 22px;
  transition: all 0.3s;
}

.bj-card.hidden {
  opacity: 0.5;
  filter: blur(2px);
}

.hand-score {
  color: #ffd700;
  font-size: 13px;
  font-weight: 600;
}

.game-center {
  padding: 8px;
  background: rgba(255, 215, 0, 0.1);
  border-radius: 6px;
  margin: 8px 0;
}

.status-text {
  font-size: 14px;
  font-weight: 600;
  color: #ffd700;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-bottom: 10px;
}

.bust-message {
  font-size: 18px;
  font-weight: 600;
  color: var(--error-color);
  margin-bottom: 10px;
}

.game-result {
  text-align: center;
  padding: 20px 0;
}

.result-icon {
  font-size: 56px;
  margin-bottom: 10px;
}

.result-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.result-detail {
  font-size: 13px;
  color: var(--font-light);
  margin-bottom: 15px;
  line-height: 1.5;
}

.result-stats {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 15px;
}

.stats-row {
  font-size: 13px;
  padding: 4px 0;
  color: var(--font-secondary);
}

.result-winnings {
  font-size: 16px;
  font-weight: 600;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid var(--border-color);
  color: var(--error-color);
}

.result-winnings.win {
  color: var(--success-color);
}

.loss-warning {
  background: #fee;
  color: #c00;
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 600;
}

.result-reward {
  font-size: 14px;
  font-weight: 600;
}
</style>
