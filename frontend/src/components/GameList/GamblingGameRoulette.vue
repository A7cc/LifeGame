<template>
  <el-dialog v-model="visible" title="轮盘赌" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">命运之轮转动，赌上你的运气与财富</div>
        <div class="warning-box">
          <span class="warning-icon">⚠️</span>
          <span class="warning-text">庄家优势5.26%，长期必输！</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 500 }}元</span>
          <span class="info-item">🏆 奖励：最高3000元</span>
        </div>

        <!-- 下注选项 -->
        <div class="bet-options">
          <div class="bet-title">选择下注类型：</div>
          <div class="option-list">
            <div v-for="option in betOptions" :key="option.id" class="bet-option" :class="{ selected: selectedBet === option.id }" @click="selectBet(option.id)">
              <div class="option-color" :style="{ background: option.color }"></div>
              <div class="option-info">
                <div class="option-name">{{ option.name }}</div>
                <div class="option-odds">赔率：{{ option.odds }}x</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 下注金额 -->
        <div class="bet-section">
          <div class="bet-label">下注金额：💰 {{ betAmount }} 元</div>
          <el-slider v-model="betAmount" :min="100" :max="2000" :step="100" show-input :disabled="!selectedBet"></el-slider>
        </div>

        <el-button type="primary" @click="handleStartGame" :disabled="!selectedBet" style="width: 100%">旋转轮盘</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <div class="roulette-wheel" :style="{ transform: `rotate(${wheelRotation}deg)` }">
          <div class="wheel-center">🎡</div>
        </div>
        <div class="wheel-pointer">▼</div>
        <div class="game-status">{{ gameStatus }}</div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
        <div class="result-stats">
          <div class="stats-item">本次下注：💰 {{ lastBet }}</div>
          <div class="stats-item" :class="{ win: totalWinnings > 0 }">{{ totalWinnings >= 0 ? '赢取' : '输掉' }}：💰 {{ Math.abs(totalWinnings) }}</div>
        </div>
        <div v-if="lossStreak > 0" class="loss-warning">
          已连续输掉 {{ lossStreak }} 局，赌博只会让你越陷越深！
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
import { useCleanupTasks } from '@/src/composables/useCleanupTasks'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'roulette', name: '轮盘赌', entryCost: 500, needBet: true })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)
const { cleanup, clearManagedTimer, setManagedInterval, setManagedTimeout } = useCleanupTasks()

let spinTimer = null

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const selectedBet = ref(null)
const betAmount = ref(500)
const wheelRotation = ref(0)
const gameStatus = ref('')
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultReward = ref('')
const lastBet = ref(0)
const totalWinnings = ref(0)
const lossStreak = ref(0)

// 下注选项 - 庄家优势设计
const betOptions = [
  { id: 'red', name: '红色', color: '#e74c3c', odds: 2, winChance: 18/38 },
  { id: 'black', name: '黑色', color: '#2c3e50', odds: 2, winChance: 18/38 },
  { id: 'even', name: '偶数', color: '#3498db', odds: 2, winChance: 18/38 },
  { id: 'odd', name: '奇数', color: '#9b59b6', odds: 2, winChance: 18/38 },
  { id: '1-18', name: '小(1-18)', color: '#1abc9c', odds: 2, winChance: 18/38 },
  { id: '19-36', name: '大(19-36)', color: '#e67e22', odds: 2, winChance: 18/38 },
]

const resetGameState = () => {
  reset()
  clearManagedTimer(spinTimer)
  spinTimer = null
  selectedBet.value = null
  betAmount.value = 500
  wheelRotation.value = 0
  gameStatus.value = ''
  lossStreak.value = 0
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const selectBet = (betId) => {
  if (gameStarted.value) return
  selectedBet.value = betId
}

const handleStartGame = async () => {
  if (!selectedBet.value) {
    ElMessage.warning('请先选择下注类型！')
    return
  }

  const success = await startGame({
    wager: betAmount.value,
    choice: selectedBet.value
  })
  if (!success) return

  lastBet.value = betAmount.value

  // 模拟轮盘旋转
  spinWheel()
}

const spinWheel = () => {
  gameStatus.value = '轮盘转动中...'
  const targetRotation = 360 * 5 + Math.random() * 360

  let current = 0
  spinTimer = setManagedInterval(() => {
    current += 15
    wheelRotation.value = current

    if (current >= targetRotation) {
      clearManagedTimer(spinTimer)
      spinTimer = null
      determineResult()
    }
  }, 20)
}

const determineResult = () => {
  const betOption = betOptions.find(b => b.id === selectedBet.value)

	const playerWins = Number(startData.value?.round?.outcome || 0) === 1
	const displayNumber = startData.value?.round?.displayNumber

	gameStatus.value = `${playerWins ? '🎉 恭喜中奖！' : '💔 很遗憾...'} 开奖号码：${displayNumber}`

  setManagedTimeout(() => {
    showResult(playerWins, betOption)
  }, 1000)
}

const showResult = async (playerWins, betOption) => {
  let winCount = 0
  let customResultText = ''
  let detail = { betType: selectedBet.value, betAmount: betAmount.value }

  if (playerWins) {
    const winnings = betAmount.value * betOption.odds
    totalWinnings.value = winnings - betAmount.value
    lossStreak.value = 0

    resultIcon.value = '🎰'
    resultTitle.value = '你赢了！'
    resultDetail.value = `幸运女神眷顾了你一次`
    resultReward.value = `赢取 ${winnings} 元（净赚 ${winnings - betAmount.value}）`
    winCount = 1
    customResultText = `获胜，赢取${winnings}元`
    detail.winnings = winnings
  } else {
    totalWinnings.value = -betAmount.value
    lossStreak.value++

    resultIcon.value = '💸'
    resultTitle.value = '你输了！'
    resultDetail.value = `庄家永远占据优势，这就是赌博的真相`
    resultReward.value = `损失 ${betAmount.value} 元`

    // 连续输的警告信息
    if (lossStreak.value >= 3) {
      resultDetail.value = `你已经连续输掉 ${lossStreak.value} 局。赌博只会让你越陷越深！`
    }

    winCount = 0
    customResultText = `失败，损失${betAmount.value}元`
    detail.winnings = 0
  }

	const gameResult = await endGame(999, customResultText, detail)
  if (gameResult) {
    totalWinnings.value = gameResult.netChange
    resultReward.value = `${gameResult.resultText}（本局净变化 ${gameResult.netChange} 元）`
  }
  emit('complete', gameResult)
}

defineExpose({
  cleanup,
})
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
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.bet-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
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

.option-color {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  flex-shrink: 0;
}

.option-info {
  flex: 1;
}

.option-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-color);
}

.option-odds {
  font-size: 11px;
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
  padding: 20px 0;
}

.roulette-wheel {
  width: 200px;
  height: 200px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: conic-gradient(
    #e74c3c 0deg 10deg, #2c3e50 10deg 20deg,
    #e74c3c 20deg 30deg, #2c3e50 30deg 40deg,
    #e74c3c 40deg 50deg, #2c3e50 50deg 60deg,
    #e74c3c 60deg 70deg, #2c3e50 70deg 80deg,
    #e74c3c 80deg 90deg, #2c3e50 90deg 100deg,
    #e74c3c 100deg 110deg, #2c3e50 110deg 120deg,
    #e74c3c 120deg 130deg, #2c3e50 130deg 140deg,
    #e74c3c 140deg 150deg, #2c3e50 150deg 160deg,
    #e74c3c 160deg 170deg, #2c3e50 170deg 180deg,
    #e74c3c 180deg 190deg, #2c3e50 190deg 200deg,
    #e74c3c 200deg 210deg, #2c3e50 210deg 220deg,
    #e74c3c 220deg 230deg, #2c3e50 230deg 240deg,
    #e74c3c 240deg 250deg, #2c3e50 250deg 260deg,
    #e74c3c 260deg 270deg, #2c3e50 270deg 280deg,
    #e74c3c 280deg 290deg, #2c3e50 290deg 300deg,
    #e74c3c 300deg 310deg, #2c3e50 310deg 320deg,
    #e74c3c 320deg 330deg, #2c3e50 330deg 340deg,
    #0a0 340deg 350deg, #2c3e50 350deg 360deg
  );
  border: 8px solid #5d4e37;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 4s cubic-bezier(0.17, 0.67, 0.12, 0.99);
  box-shadow: 0 0 20px rgba(0,0,0,0.3);
}

.wheel-center {
  width: 80px;
  height: 80px;
  background: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  border: 4px solid #5d4e37;
}

.wheel-pointer {
  font-size: 32px;
  color: #ffd700;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
  margin-bottom: 10px;
}

.game-status {
  font-size: 18px;
  font-weight: 600;
  color: var(--font-color);
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
  font-size: 13px;
  color: var(--font-light);
  margin-bottom: 15px;
  line-height: 1.5;
}

.result-stats {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 15px;
}

.stats-item {
  font-size: 14px;
  padding: 8px 15px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.stats-item.win {
  background: #f0f9ff;
  color: var(--success-color);
  font-weight: 600;
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
