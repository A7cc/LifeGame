<template>
  <el-dialog v-model="visible" title="赛马博彩" width="600px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">挑选你心目中的冠军马匹，下注赢取丰厚奖励！</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 1000 }}元</span>
          <span class="info-item">🏆 奖励：最高5000元</span>
        </div>

        <!-- 马匹列表 -->
        <div class="horses-list">
          <div v-for="horse in horses" :key="horse.id" class="horse-item" :class="{ selected: selectedHorse === horse.id }" @click="selectHorse(horse.id)">
            <div class="horse-icon">{{ horse.icon }}</div>
            <div class="horse-info">
              <div class="horse-name">{{ horse.name }}</div>
              <div class="horse-odds">赔率：{{ horse.odds }}x</div>
              <div class="horse-stats">{{ horse.speed }}速度 | {{ horse.stamina }}耐力</div>
            </div>
          </div>
        </div>

        <!-- 下注金额 -->
        <div class="bet-section">
          <div class="bet-label">下注金额：💰 {{ betAmount }} 元</div>
          <el-slider v-model="betAmount" :min="100" :max="5000" :step="100" show-input :disabled="!selectedHorse"></el-slider>
        </div>

        <el-button type="primary" @click="handleStartGame" :disabled="!selectedHorse" style="width: 100%">开始比赛</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <div class="race-track">
          <div class="track-lanes">
            <div v-for="horse in racingHorses" :key="horse.id" class="race-lane">
              <div class="lane-number">{{ horse.lane }}</div>
              <div class="lane-track">
                <div class="horse-position" :style="{ left: horse.position + '%' }">
                  <span class="racing-horse">{{ horse.icon }}</span>
                </div>
                <div class="finish-line">🏁</div>
              </div>
            </div>
          </div>
        </div>

        <div class="race-status">{{ raceStatus }}</div>
        <div class="race-progress">比赛进行中...</div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
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
    default: () => ({ id: 'horse-racing', name: '赛马博彩', entryCost: 1000, needBet: true })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)
const { cleanup, clearManagedTimer, setManagedInterval, setManagedTimeout } = useCleanupTasks()

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const selectedHorse = ref(null)
const betAmount = ref(1000)

const raceStatus = ref('')
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultReward = ref('')

// 可选的马匹
const horses = ref([
  { id: 1, name: '闪电', icon: '🐎', odds: 1.7, speed: 95, stamina: 85, baseSpeed: 2.5 },
  { id: 2, name: '烈火', icon: '🦄', odds: 2.2, speed: 90, stamina: 90, baseSpeed: 2.3 },
  { id: 3, name: '风暴', icon: '🐴', odds: 4.8, speed: 85, stamina: 95, baseSpeed: 2.1 },
  { id: 4, name: '黑马', icon: '🐃', odds: 8.8, speed: 80, stamina: 80, baseSpeed: 1.9 },
  { id: 5, name: '幸运', icon: '🦓', odds: 14.0, speed: 75, stamina: 75, baseSpeed: 1.7 },
])

// 比赛中的马匹
const racingHorses = ref([])
let raceInterval = null

const resetGameState = () => {
  reset()
  selectedHorse.value = null
  betAmount.value = 1000
  raceStatus.value = ''
  if (raceInterval) {
    clearManagedTimer(raceInterval)
    raceInterval = null
  }
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const selectHorse = (horseId) => {
  if (gameStarted.value) return
  selectedHorse.value = horseId
}

const handleStartGame = async () => {
  if (!selectedHorse.value) {
    ElMessage.warning('请先选择一匹马！')
    return
  }

  const success = await startGame({
    wager: betAmount.value,
    choice: String(selectedHorse.value)
  })
  if (!success) return

  // 初始化比赛马匹
	const finishOrder = startData.value?.round?.finishOrder || []
	racingHorses.value = horses.value.map((horse, index) => ({
		...horse,
		lane: index + 1,
		baseSpeed: 2.7 - Math.max(0, finishOrder.indexOf(horse.id)) * 0.25,
    position: 0,
    finished: false
  }))

  // 开始比赛
  startRace()
}

const startRace = () => {
  raceStatus.value = '比赛开始！'

  raceInterval = setManagedInterval(() => {
    let allFinished = true

    racingHorses.value.forEach(horse => {
      if (!horse.finished) {
        allFinished = false
        // 计算移动速度（带随机性）
        const randomFactor = 0.8 + Math.random() * 0.4
        const moveAmount = horse.baseSpeed * randomFactor * 0.8
        horse.position = Math.min(horse.position + moveAmount, 90)

        if (horse.position >= 90 && !horse.finished) {
          horse.finished = true
          horse.finishOrder = racingHorses.value.filter(h => h.finished).length
        }
      }
    })

    // 更新状态
    const finishedCount = racingHorses.value.filter(h => h.finished).length
    if (finishedCount > 0) {
      const leadingHorse = [...racingHorses.value]
        .filter(h => h.finished)
        .sort((a, b) => a.finishOrder - b.finishOrder)[0]
      raceStatus.value = `${leadingHorse.name} 率先冲线！`
    }

    // 检查比赛是否结束
    if (allFinished || racingHorses.value.filter(h => h.finished).length >= 3) {
      clearManagedTimer(raceInterval)
      raceInterval = null
      setManagedTimeout(() => endRace(), 500)
    }
  }, 200)
}

const endRace = async () => {
	// 最终名次由后端开局时生成；前端赛道移动只负责动画展示。
	const authoritativeOrder = startData.value?.round?.finishOrder || []
	const finalStandings = [...racingHorses.value].sort((a, b) =>
		authoritativeOrder.indexOf(a.id) - authoritativeOrder.indexOf(b.id)
	)

  const winner = finalStandings[0]
  const userHorse = racingHorses.value.find(h => h.id === selectedHorse.value)
  const userRank = finalStandings.findIndex(h => h.id === selectedHorse.value) + 1

  let winCount = 0
  let customResultText = ''
  let detail = { horse: userHorse.name, rank: userRank, betAmount: betAmount.value }

  if (userHorse.id === winner.id) {
    // 玩家下注的马获胜
    const winnings = Math.round(betAmount.value * userHorse.odds)
    resultIcon.value = '🏆'
    resultTitle.value = `${userHorse.name} 获得冠军！`
    resultDetail.value = `你赢取了 ${winnings} 元！赔率 ${userHorse.odds}x`
		resultReward.value = '等待后端结算...'
    winCount = 1
    customResultText = `获胜，${userHorse.name}第1名，赢取${winnings}元`
    detail.winnings = winnings
  } else if (userRank <= 3) {
    // 进入前三名，获得部分奖励
    const consolationPrize = Math.round(betAmount.value * 0.5)
    resultIcon.value = '🥈'
    resultTitle.value = `${userHorse.name} 第${userRank}名`
    resultDetail.value = `进入前三，获得安慰奖 ${consolationPrize} 元`
		resultReward.value = '等待后端结算...'
    winCount = 2 // 后端约定：2 表示进入前三名的安慰奖
    customResultText = `${userHorse.name}第${userRank}名，获得${consolationPrize}元`
    detail.winnings = consolationPrize
  } else {
    // 没有进入前三
    resultIcon.value = '😢'
    resultTitle.value = `${userHorse.name} 第${userRank}名`
    resultDetail.value = '很遗憾，没有进入前三名'
		resultReward.value = '等待后端结算...'
    winCount = 0
    customResultText = `${userHorse.name}第${userRank}名`
    detail.winnings = 0
  }

	const gameResult = await endGame(999, customResultText, detail)
  if (gameResult) {
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

.horses-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
  max-height: 280px;
  overflow-y: auto;
}

.horse-item {
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

.horse-item:hover {
  border-color: #409eff;
}

.horse-item.selected {
  border-color: var(--success-color);
}

.horse-icon {
  font-size: 32px;
  flex-shrink: 0;
}

.horse-info {
  flex: 1;
  text-align: left;
}

.horse-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 4px;
}

.horse-odds {
  font-size: 14px;
  color: var(--warning-color);
  font-weight: 600;
  margin-bottom: 2px;
}

.horse-stats {
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

.race-track {
  background: linear-gradient(90deg, #8b7355 0%, #c4a77d 50%, #8b7355 100%);
  border: 4px solid #5d4e37;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 15px;
}

.track-lanes {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.race-lane {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lane-number {
  width: 24px;
  height: 24px;
  background: var(--panel-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.lane-track {
  flex: 1;
  height: 40px;
  background: #a0d8ef;
  border-radius: 20px;
  position: relative;
  border: 2px dashed #fff;
  overflow: hidden;
}

.horse-position {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  transition: left 0.2s linear;
}

.racing-horse {
  font-size: 28px;
  filter: drop-shadow(2px 2px 2px rgba(0,0,0,0.3));
}

.finish-line {
  position: absolute;
  right: 5%;
  top: 50%;
  transform: translateY(-50%);
  font-size: 24px;
}

.race-status {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 8px;
}

.race-progress {
  font-size: 14px;
  color: var(--font-light);
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
  margin-bottom: 8px;
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
  font-weight: 600;
}
</style>
