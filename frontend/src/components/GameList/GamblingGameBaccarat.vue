<template>
  <el-dialog v-model="visible" title="百家乐" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">百家乐 - 亚洲最受欢迎的博彩游戏</div>
        <div class="warning-box">
          <span class="warning-icon">⚠️</span>
          <span class="warning-text">庄家优势1.06%，看似很小但足以让你破产！</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 1000 }}元</span>
          <span class="info-item">🏆 奖励：最高6000元</span>
        </div>

        <!-- 下注选项 -->
        <div class="bet-options">
          <div class="bet-title">选择下注方：</div>
          <div class="option-list">
            <div class="bet-option" :class="{ selected: selectedBet === 'player' }" @click="selectBet('player')">
              <div class="option-icon">👤</div>
              <div class="option-info">
                <div class="option-name">闲家</div>
                <div class="option-odds">1:1 赔率</div>
              </div>
            </div>
            <div class="bet-option" :class="{ selected: selectedBet === 'banker' }" @click="selectBet('banker')">
              <div class="option-icon">🏦</div>
              <div class="option-info">
                <div class="option-name">庄家</div>
                <div class="option-odds">1:0.95 赔率（抽水5%）</div>
              </div>
            </div>
            <div class="bet-option" :class="{ selected: selectedBet === 'tie' }" @click="selectBet('tie')">
              <div class="option-icon">🤝</div>
              <div class="option-info">
                <div class="option-name">和局</div>
                <div class="option-odds">1:8 赔率</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 下注金额 -->
        <div class="bet-section">
          <div class="bet-label">下注金额：💰 {{ betAmount }} 元</div>
          <el-slider v-model="betAmount" :min="200" :max="5000" :step="200" show-input :disabled="!selectedBet"></el-slider>
        </div>

        <el-button type="primary" @click="handleStartGame" :disabled="!selectedBet" style="width: 100%">发牌</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <div class="baccarat-table">
          <!-- 庄家区域 -->
          <div class="hand-area banker-area">
            <div class="area-title">🏦 庄家</div>
            <div class="cards-area">
              <span v-for="(card, index) in bankerCards" :key="index" class="baccarat-card" :class="{ revealed: cardRevealed.banker >= index }">
                {{ cardRevealed.banker >= index ? card : '🂠' }}
              </span>
            </div>
            <div class="hand-score" v-if="cardRevealed.banker >= 2">
              点数：{{ bankerScore }}
            </div>
          </div>

          <!-- 游戏状态 -->
          <div class="game-center">
            <div class="status-text">{{ gameStatus }}</div>
          </div>

          <!-- 闲家区域 -->
          <div class="hand-area player-area">
            <div class="area-title">👤 闲家</div>
            <div class="cards-area">
              <span v-for="(card, index) in playerCards" :key="index" class="baccarat-card" :class="{ revealed: cardRevealed.player >= index }">
                {{ cardRevealed.player >= index ? card : '🂠' }}
              </span>
            </div>
            <div class="hand-score" v-if="cardRevealed.player >= 2">
              点数：{{ playerScore }}
            </div>
          </div>
        </div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
        <div class="result-stats">
          <div class="stats-row">庄家：{{ bankerCards.join(' ') }} = {{ bankerScore }}点</div>
          <div class="stats-row">闲家：{{ playerCards.join(' ') }} = {{ playerScore }}点</div>
          <div class="stats-row result-winnings" :class="{ win: totalWinnings > 0 }">
            {{ totalWinnings >= 0 ? '赢取' : '输掉' }}：💰 {{ Math.abs(totalWinnings) }}
          </div>
        </div>
        <div v-if="totalLost > 10000" class="loss-warning">
          💀 警告：你已累计损失 {{ totalLost }} 元！这就是赌博的下场！
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'baccarat', name: '百家乐', entryCost: 1000, needBet: true })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const selectedBet = ref(null)
const betAmount = ref(1000)
const gameStatus = ref('')
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultReward = ref('')
const totalWinnings = ref(0)
const totalLost = ref(0)
const roundsPlayed = ref(0)

const playerCards = ref([])
const bankerCards = ref([])
const cardRevealed = ref({ player: 0, banker: 0 })

const resetGameState = () => {
  reset()
  selectedBet.value = null
  betAmount.value = 1000
  gameStatus.value = ''
  playerCards.value = []
  bankerCards.value = []
  cardRevealed.value = { player: 0, banker: 0 }
  totalLost.value = 0
  roundsPlayed.value = 0
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const selectBet = (betId) => {
  if (gameStarted.value) return
  selectedBet.value = betId
}

const playerScore = ref(0)
const bankerScore = ref(0)

const handleStartGame = async () => {
  if (!selectedBet.value) {
    ElMessage.warning('请先选择下注方！')
    return
  }

  const success = await startGame({
    wager: betAmount.value,
    choice: selectedBet.value
  })
  if (!success) return

  roundsPlayed.value++
  gameEnded.value = false

	// 牌面和最终结果由后端生成，前端只负责逐张展示。
	playerCards.value = startData.value?.round?.playerCards || []
	bankerCards.value = startData.value?.round?.bankerCards || []
	playerScore.value = Number(startData.value?.round?.playerScore || 0)
	bankerScore.value = Number(startData.value?.round?.bankerScore || 0)

  // 逐步揭示牌
  revealCards()
}

const revealCards = () => {
  gameStatus.value = '发牌中...'

  const revealSequence = async () => {
    // 揭示闲家第一张
    await delay(500)
    cardRevealed.value.player = 0

    // 揭示庄家第一张
    await delay(500)
    cardRevealed.value.banker = 0

    // 揭示闲家第二张
    await delay(500)
    cardRevealed.value.player = 1

    // 揭示庄家第二张
    await delay(500)
    cardRevealed.value.banker = 1

		// 第三张牌规则已由后端执行，这里只继续揭示最终牌面。
		await delay(800)
		cardRevealed.value.player = playerCards.value.length - 1
		cardRevealed.value.banker = bankerCards.value.length - 1
		await delay(500)
		determineResult()
  }

  revealSequence()
}

const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

const determineResult = async () => {
  let winnings = -betAmount.value

	const winner = startData.value?.round?.winner
	const adjustedPlayerWin = winner === 'player'
	const adjustedBankerWin = winner === 'banker'
	const adjustedTie = winner === 'tie'

  let winCount = 0
  let customResultText = ''
  let detail = { betType: selectedBet.value, betAmount: betAmount.value, bankerScore: bankerScore.value, playerScore: playerScore.value }

  if (selectedBet.value === 'player') {
    if (adjustedPlayerWin) {
      winnings = betAmount.value * 1
      resultIcon.value = '🎉'
      resultTitle.value = '闲家胜！你赢了！'
      resultReward.value = `赢取 ${betAmount.value + winnings} 元`
      winCount = 1
      customResultText = `获胜，赢取${winnings}元`
      detail.winnings = winnings
    } else if (adjustedTie) {
      winnings = 0
      resultIcon.value = '🤝'
      resultTitle.value = '和局'
      resultReward.value = '退还下注'
      winCount = 2
      customResultText = '和局'
      detail.winnings = 0
    } else {
      totalLost.value += betAmount.value
      resultIcon.value = '💸'
      resultTitle.value = '庄家胜！你输了！'
      resultReward.value = `损失 ${betAmount.value} 元`
      winCount = 0
      customResultText = `失败，损失${betAmount.value}元`
      detail.winnings = 0
    }
  } else if (selectedBet.value === 'banker') {
    if (adjustedBankerWin) {
      winnings = Math.round(betAmount.value * 0.95)
      resultIcon.value = '🎉'
      resultTitle.value = '庄家胜！你赢了！'
      resultReward.value = `赢取 ${betAmount.value + winnings} 元（扣除5%佣金）`
      winCount = 1
      customResultText = `获胜，赢取${winnings}元`
      detail.winnings = winnings
    } else if (adjustedTie) {
      winnings = 0
      resultIcon.value = '🤝'
      resultTitle.value = '和局'
      resultReward.value = '退还下注'
      winCount = 2
      customResultText = '和局'
      detail.winnings = 0
    } else {
      totalLost.value += betAmount.value
      resultIcon.value = '💸'
      resultTitle.value = '闲家胜！你输了！'
      resultReward.value = `损失 ${betAmount.value} 元`
      winCount = 0
      customResultText = `失败，损失${betAmount.value}元`
      detail.winnings = 0
    }
  } else if (selectedBet.value === 'tie') {
    if (adjustedTie) {
      winnings = betAmount.value * 8
      resultIcon.value = '🎰'
      resultTitle.value = '和局！大奖！'
      resultReward.value = `赢取 ${betAmount.value + winnings} 元！`
      winCount = 1
      customResultText = `获胜，赢取${winnings}元`
      detail.winnings = winnings
    } else {
      totalLost.value += betAmount.value
      resultIcon.value = '💸'
		resultTitle.value = adjustedPlayerWin ? '闲家胜' : '庄家胜'
      resultDetail.value = '很遗憾，不是和局'
      resultReward.value = `损失 ${betAmount.value} 元`
      winCount = 0
      customResultText = `失败，损失${betAmount.value}元`
      detail.winnings = 0
    }
  }

  totalWinnings.value = winnings
  gameStatus.value = adjustedTie ? '和局！' : (adjustedPlayerWin ? '闲家胜！' : '庄家胜！')

  if (winnings <= 0 && roundsPlayed.value > 5) {
    resultDetail.value = `已玩 ${roundsPlayed.value} 局，累计损失 ${totalLost.value} 元。庄家优势虽小，但长期必输！`
  }

	const gameResult = await endGame(999, customResultText, detail)
  if (gameResult) {
    totalWinnings.value = gameResult.netChange
    resultReward.value = `${gameResult.resultText}（本局净变化 ${gameResult.netChange} 元）`
  }
  emit('complete', gameResult)
}
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

.bet-options {
  margin-bottom: 20px;
}

.bet-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  text-align: left;
}

.option-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bet-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--panel-color);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.bet-option:hover {
  border-color: #409eff;
}

.bet-option.selected {
  border-color: var(--success-color);
}

.option-icon {
  font-size: 28px;
}

.option-info {
  flex: 1;
  text-align: left;
}

.option-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 2px;
}

.option-odds {
  font-size: 12px;
  color: var(--font-light);
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

.baccarat-table {
  background: #1a472a;
  border-radius: 16px;
  padding: 20px;
  border: 6px solid #8b4513;
}

.hand-area {
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 10px;
}

.banker-area {
  background: rgba(255, 100, 100, 0.1);
  border: 1px solid rgba(255, 100, 100, 0.3);
}

.player-area {
  background: rgba(100, 100, 255, 0.1);
  border: 1px solid rgba(100, 100, 255, 0.3);
}

.area-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 10px;
  color: #fff;
}

.cards-area {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 8px;
}

.baccarat-card {
  font-size: 24px;
  opacity: 0.5;
  transition: opacity 0.3s;
}

.baccarat-card.revealed {
  opacity: 1;
}

.hand-score {
  color: #ffd700;
  font-size: 14px;
  font-weight: 600;
}

.game-center {
  padding: 10px;
  background: rgba(255, 215, 0, 0.1);
  border-radius: 8px;
  margin: 10px 0;
}

.status-text {
  font-size: 16px;
  font-weight: 600;
  color: #ffd700;
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
  padding: 12px;
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
