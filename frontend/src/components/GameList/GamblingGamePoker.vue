<template>
  <el-dialog v-model="visible" title="德州扑克" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">心理博弈，筹码比拼，牌运与智慧的较量</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：2000元</span>
          <span class="info-item">🏆 奖励：5000元</span>
        </div>
        <div class="table-info">
          <div class="info-row">
            <span>🎴 单桌对战</span>
            <span>👤 3位玩家</span>
          </div>
          <div class="info-row">
            <span>💵 起始筹码：5000</span>
            <span>🎯 盲注：50/100</span>
          </div>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%">上桌</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="poker-table">
          <!-- 公共牌 -->
          <div class="community-cards">
            <div class="card-placeholder" v-for="(card, index) in 5" :key="index">
              <span v-if="communityCards[index]" class="poker-card">{{ communityCards[index] }}</span>
              <span v-else class="card-back">🂠</span>
            </div>
          </div>

          <!-- 对手信息 -->
          <div class="opponents">
            <div class="opponent">
              <span class="opponent-name">玩家2</span>
              <span class="opponent-chips">💰 {{ opponent2Chips }}</span>
              <span class="opponent-cards">{{ opponent2Cards > 0 ? '🂠🂠' : '💭' }}</span>
            </div>
            <div class="opponent">
              <span class="opponent-name">玩家3</span>
              <span class="opponent-chips">💰 {{ opponent3Chips }}</span>
              <span class="opponent-cards">{{ opponent3Cards > 0 ? '🂠🂠' : '💭' }}</span>
            </div>
          </div>

          <!-- 玩家信息 -->
          <div class="player-area">
            <div class="player-hand">
              <span class="poker-card" v-for="card in playerHand" :key="card">{{ card }}</span>
            </div>
            <div class="player-info">
              <span>你的筹码：💰 {{ playerChips }}</span>
              <span class="hand-rank" v-if="handRank">{{ handRank }}</span>
            </div>
          </div>

          <!-- 底池 -->
          <div class="pot-display">
            <span class="pot-label">底池</span>
            <span class="pot-amount">💰 {{ pot }}</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <el-button @click="fold" :disabled="processing" type="info">弃牌</el-button>
          <el-button @click="call" :disabled="processing">{{ currentBet > 0 ? '跟注' : '过牌' }}</el-button>
          <el-button @click="raise" :disabled="processing || playerChips < currentBet * 2">加注</el-button>
          <el-button @click="allIn" :disabled="processing || playerChips <= 0" type="danger">全押</el-button>
        </div>

        <div class="game-status">{{ gameStatus }}</div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
        <div class="result-chips">最终筹码：💰 {{ playerChips }}</div>
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
    default: () => ({ id: 'poker', name: '德州扑克', entryCost: 2000, needBet: false })
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

const STARTING_CHIPS = 5000

const suits = ['♠', '♥', '♦', '♣']
const ranks = ['A', 'K', 'Q', 'J', '10', '9', '8', '7', '6', '5', '4', '3', '2']

const playerChips = ref(STARTING_CHIPS)
const opponent2Chips = ref(STARTING_CHIPS)
const opponent3Chips = ref(STARTING_CHIPS)
const pot = ref(0)
const currentBet = ref(0)

const playerHand = ref([])
const communityCards = ref([])
const opponent2Cards = ref(0)
const opponent3Cards = ref(0)

const gameStatus = ref('')
const handRank = ref('')
const round = ref(0) // 0: preflop, 1: flop, 2: turn, 3: river

const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultReward = ref('')

const generateCard = () => {
  const suit = suits[Math.floor(Math.random() * suits.length)]
  const rank = ranks[Math.floor(Math.random() * ranks.length)]
  const color = (suit === '♥' || suit === '♦') ? 'red' : 'black'
  return { display: `${suit}${rank}`, color, rank, suit }
}

const resetGameState = () => {
  reset()
  playerChips.value = STARTING_CHIPS
  opponent2Chips.value = STARTING_CHIPS
  opponent3Chips.value = STARTING_CHIPS
  pot.value = 0
  currentBet.value = 0
  playerHand.value = []
  communityCards.value = []
  opponent2Cards.value = 0
  opponent3Cards.value = 0
  gameStatus.value = ''
  handRank.value = ''
  round.value = 0
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const handleStartGame = async () => {
  const success = await startGame()
  if (!success) return

  // 发手牌
  playerHand.value = [generateCard().display, generateCard().display]
  opponent2Cards.value = 2
  opponent3Cards.value = 2

  // 盲注
  playerChips.value -= 50
  opponent2Chips.value -= 100
  opponent3Chips.value -= 50
  pot.value = 200
  currentBet.value = 100

  gameStatus.value = 'Preflop - 轮到你行动'
  evaluateHand()
}

const evaluateHand = () => {
  const handStrength = Math.random()
  if (handStrength < 0.05) {
    handRank.value = '皇家同花顺！'
  } else if (handStrength < 0.1) {
    handRank.value = '同花顺'
  } else if (handStrength < 0.2) {
    handRank.value = '四条'
  } else if (handStrength < 0.3) {
    handRank.value = '葫芦'
  } else if (handStrength < 0.45) {
    handRank.value = '同花'
  } else if (handStrength < 0.6) {
    handRank.value = '顺子'
  } else if (handStrength < 0.75) {
    handRank.value = '三条'
  } else if (handStrength < 0.85) {
    handRank.value = '两对'
  } else {
    handRank.value = '一对'
  }
}

const nextRound = () => {
  round.value++
  if (round.value === 1) {
    // Flop
    communityCards.value = [
      generateCard().display,
      generateCard().display,
      generateCard().display
    ]
    gameStatus.value = 'Flop - 轮到你行动'
  } else if (round.value === 2) {
    // Turn
    communityCards.value.push(generateCard().display)
    gameStatus.value = 'Turn - 轮到你行动'
  } else if (round.value === 3) {
    // River
    communityCards.value.push(generateCard().display)
    gameStatus.value = 'River - 最后行动'
  } else {
    // Showdown
    showdown()
    return
  }
  evaluateHand()
  currentBet.value = 0
}

const showdown = async () => {
	const authoritativeWinner = startData.value?.round?.winner

  let winCount = 0
  let customResultText = ''
  let detail = {}

	if (authoritativeWinner === 'player') {
    // 玩家赢
    const winnings = pot.value
    playerChips.value += winnings
    resultIcon.value = '🏆'
    resultTitle.value = '你赢了！'
    resultDetail.value = `赢得底池 ${winnings} 筹码`
		resultReward.value = '等待后端结算...'
    winCount = 1
    customResultText = `获胜，筹码 ${playerChips.value}`
    detail = { finalChips: playerChips.value, pot: winnings }
	} else if (authoritativeWinner === 'opponent2') {
    // 玩家2赢
    opponent2Chips.value += pot.value
    resultIcon.value = '😢'
    resultTitle.value = '玩家2获胜'
    resultDetail.value = '你的牌不够大'
		resultReward.value = '等待后端结算...'
    winCount = 0
    customResultText = '失败'
    detail = { winner: '玩家2' }
  } else {
    // 玩家3赢
    opponent3Chips.value += pot.value
    resultIcon.value = '😢'
    resultTitle.value = '玩家3获胜'
    resultDetail.value = '运气不佳'
		resultReward.value = '等待后端结算...'
    winCount = 0
    customResultText = '失败'
    detail = { winner: '玩家3' }
  }

	const gameResult = await endGame(999, customResultText, detail)
  if (gameResult) {
    resultReward.value = `${gameResult.resultText}（本局净变化 ${gameResult.netChange} 元）`
  }
  emit('complete', gameResult)
}

const fold = () => {
  processing.value = true
  setTimeout(() => {
    // 对手行动
    if (Math.random() < 0.3) {
      opponent2Chips.value -= Math.floor(Math.random() * 200) + 50
      opponent3Chips.value -= Math.floor(Math.random() * 200) + 50
    }
    nextRound()
    processing.value = false
  }, 500)
}

const call = () => {
  processing.value = true
  setTimeout(() => {
    const callAmount = currentBet.value
    playerChips.value -= callAmount
    pot.value += callAmount

    // 对手行动
    const opp2Action = Math.random()
    const opp3Action = Math.random()

    if (opp2Action < 0.6) {
      const bet = Math.floor(Math.random() * 200) + 50
      opponent2Chips.value -= bet
      pot.value += bet
    }

    if (opp3Action < 0.6) {
      const bet = Math.floor(Math.random() * 200) + 50
      opponent3Chips.value -= bet
      pot.value += bet
    }

    nextRound()
    processing.value = false
  }, 500)
}

const raise = () => {
  processing.value = true
  setTimeout(() => {
    const raiseAmount = currentBet.value * 2
    playerChips.value -= raiseAmount
    pot.value += raiseAmount
    currentBet.value = raiseAmount

    // 对手跟注或弃牌
    if (Math.random() < 0.7) {
      opponent2Chips.value -= raiseAmount
      opponent3Chips.value -= raiseAmount
      pot.value += raiseAmount * 2
    }

    nextRound()
    processing.value = false
  }, 500)
}

const allIn = () => {
  processing.value = true
  setTimeout(() => {
    const allInAmount = playerChips.value
    pot.value += allInAmount
    playerChips.value = 0

    // 对手反应
    if (Math.random() < 0.5) {
      opponent2Chips.value -= Math.min(opponent2Chips.value, allInAmount)
      opponent3Chips.value -= Math.min(opponent3Chips.value, allInAmount)
      pot.value += Math.min(opponent2Chips.value, allInAmount)
      pot.value += Math.min(opponent3Chips.value, allInAmount)
    }

    // 直接进入摊牌
    round.value = 4
    showdown()
    processing.value = false
  }, 500)
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
  margin-bottom: 20px;
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

.table-info {
  background: #1a472a;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  color: #fff;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  margin-bottom: 8px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.game-playing {
  text-align: center;
}

.poker-table {
  background: #1a472a;
  border: 8px solid #8b4513;
  border-radius: 100px;
  padding: 20px;
  margin-bottom: 15px;
}

.community-cards {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 15px;
}

.card-placeholder {
  width: 40px;
  height: 56px;
  background: #fff;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.poker-card {
  font-size: 18px;
  font-weight: 600;
}

.poker-card:nth-child(odd) {
  color: var(--error-color);
}

.poker-card:nth-child(even) {
  color: var(--font-color);
}

.card-back {
  font-size: 24px;
}

.opponents {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 15px;
}

.opponent {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  color: #fff;
}

.opponent-name {
  font-size: 12px;
  display: block;
  margin-bottom: 4px;
}

.opponent-chips {
  font-size: 11px;
  display: block;
  margin-bottom: 4px;
}

.opponent-cards {
  font-size: 14px;
}

.player-area {
  margin-top: 15px;
}

.player-hand {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 10px;
}

.player-info {
  color: #fff;
  font-size: 13px;
}

.hand-rank {
  margin-left: 10px;
  color: #ffd700;
  font-weight: 600;
}

.pot-display {
  text-align: center;
  color: #fff;
  margin-top: 10px;
}

.pot-label {
  font-size: 12px;
  display: block;
}

.pot-amount {
  font-size: 16px;
  font-weight: 600;
  color: #ffd700;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 10px;
}

.game-status {
  font-size: 12px;
  color: var(--font-light);
  text-align: center;
}

.game-result {
  text-align: center;
  padding: 20px 0;
}

.result-icon {
  font-size: 64px;
  margin-bottom: 10px;
}

.result-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 8px;
}

.result-detail {
  font-size: 14px;
  color: var(--font-light);
  margin-bottom: 5px;
}

.result-chips {
  font-size: 14px;
  margin-bottom: 8px;
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
