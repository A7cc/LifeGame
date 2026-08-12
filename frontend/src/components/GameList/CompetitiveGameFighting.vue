<template>
  <el-dialog v-model="visible" title="格斗竞技" width="450px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">近身格斗对决，展示你的战斗技巧</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 300 }}元</span>
        </div>
        <div class="fighter-select">
          <div class="select-title">选择你的格斗家：</div>
          <div class="fighter-list">
            <div v-for="fighter in fighters" :key="fighter.id" class="fighter-item" :class="{ selected: selectedFighter === fighter.id }" @click="selectedFighter = fighter.id">
              <span class="fighter-icon">{{ fighter.icon }}</span>
              <span class="fighter-name">{{ fighter.name }}</span>
              <span class="fighter-style">{{ fighter.style }}</span>
            </div>
          </div>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%; margin-top: 15px" :disabled="!selectedFighter">开始格斗</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="fighting-arena">
          <div class="health-bars">
            <div class="health-bar">
              <span class="fighter-label">你</span>
              <div class="bar-bg">
                <div class="bar-fill player" :style="{ width: playerHealth + '%' }"></div>
              </div>
              <span class="health-text">{{ playerHealth }}%</span>
            </div>
            <div class="health-bar">
              <span class="health-text">{{ enemyHealth }}%</span>
              <div class="bar-bg">
                <div class="bar-fill enemy" :style="{ width: enemyHealth + '%' }"></div>
              </div>
              <span class="fighter-label">对手</span>
            </div>
          </div>
          <div class="arena-display">
            <div class="fighter-sprite player-fighter">{{ selectedFighterObj?.icon }}</div>
            <div class="vs-text">VS</div>
            <div class="fighter-sprite enemy-fighter">🥷</div>
          </div>
        </div>
        <div class="action-buttons">
          <el-button @click="attack" :disabled="processing">👊 攻击</el-button>
          <el-button @click="defend" :disabled="processing">🛡️ 防御</el-button>
          <el-button @click="special" :disabled="processing || !specialReady">⚡ 必杀</el-button>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">剩余血量: {{ playerHealth }}%</div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'fighting', name: '格斗竞技', entryCost: 300 })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const { gameStarted, gameEnded, processing, startGame: startMiniGame, endGame: endMiniGame, reset } = useMiniGameBase(props.config)

const fighters = [
  { id: 1, name: '拳击手', icon: '🥊', style: '力量型' },
  { id: 2, name: '武术家', icon: '🥋', style: '敏捷型' },
  { id: 3, name: '摔跤手', icon: '💪', style: '防御型' },
]

const selectedFighter = ref(null)

const playerHealth = ref(100)
const enemyHealth = ref(100)
const specialReady = ref(false)

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

let turnCount = 0

const selectedFighterObj = computed(() => fighters.find(f => f.id === selectedFighter.value))

const resetGameState = () => {
  reset()
  selectedFighter.value = null
  playerHealth.value = 100
  enemyHealth.value = 100
  specialReady.value = false
  turnCount = 0
  resultIcon.value = ''
  resultTitle.value = ''
  resultReward.value = ''
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const handleStartGame = async () => {
  const success = await startMiniGame()
  if (!success) return
  ElMessage.info('格斗开始！')
}

const enemyAction = () => {
  const actions = ['attack', 'attack', 'defend', 'special']
  return actions[Math.floor(Math.random() * actions.length)]
}

const processTurn = (playerAction) => {
  const enemy = enemyAction()
  let playerDamage = 0
  let enemyDamage = 0

  // 玩家行动
  if (playerAction === 'attack') {
    enemyDamage = Math.floor(Math.random() * 15) + 10
  } else if (playerAction === 'special' && specialReady.value) {
    enemyDamage = Math.floor(Math.random() * 25) + 20
    specialReady.value = false
  }

  // 敌人行动
  if (enemy === 'attack') {
    playerDamage = Math.floor(Math.random() * 15) + 10
  } else if (enemy === 'special' && turnCount > 2) {
    playerDamage = Math.floor(Math.random() * 20) + 15
  }

  // 防御减伤
  if (playerAction === 'defend') {
    playerDamage = Math.floor(playerDamage * 0.3)
  }
  if (enemy === 'defend') {
    enemyDamage = Math.floor(enemyDamage * 0.3)
  }

  // 应用伤害
  playerHealth.value = Math.max(0, playerHealth.value - playerDamage)
  enemyHealth.value = Math.max(0, enemyHealth.value - enemyDamage)

  turnCount++
  if (turnCount >= 3) {
    specialReady.value = true
  }

  return { playerDamage, enemyDamage }
}

const attack = () => {
  processing.value = true
  const { playerDamage, enemyDamage } = processTurn('attack')
  setTimeout(() => {
    ElMessage.info(`你造成 ${enemyDamage} 伤害，受到 ${playerDamage} 伤害`)
    processing.value = false
    checkEnd()
  }, 500)
}

const defend = () => {
  processing.value = true
  const { playerDamage, enemyDamage } = processTurn('defend')
  setTimeout(() => {
    ElMessage.info(`防御姿态！受到 ${playerDamage} 伤害`)
    processing.value = false
    checkEnd()
  }, 500)
}

const special = () => {
  if (!specialReady.value) return
  processing.value = true
  const { playerDamage, enemyDamage } = processTurn('special')
  setTimeout(() => {
    ElMessage.success(`必杀技！造成 ${enemyDamage} 伤害，受到 ${playerDamage} 伤害`)
    processing.value = false
    checkEnd()
  }, 600)
}

const checkEnd = () => {
  if (playerHealth.value <= 0 || enemyHealth.value <= 0) {
    handleEndGame()
  }
}

const handleEndGame = async () => {
  const isWin = playerHealth.value > 0
  const winCount = isWin ? 1 : 0

  const resultText = isWin
    ? `KO获胜，剩余血量 ${playerHealth.value}%`
    : '被KO'

  const detail = {
    playerHealth: playerHealth.value,
    enemyHealth: enemyHealth.value,
    turns: turnCount
  }

  const gameResult = await endMiniGame(winCount, resultText, detail)

  resultIcon.value = isWin ? '🏆' : '😵'
  resultTitle.value = isWin ? 'KO！胜利！' : '被KO！'
  resultReward.value = gameResult
    ? `获得 ${gameResult.cashChange} 元，名声 +${gameResult.fameChange}`
    : ''

  emit('complete', {
    game: '格斗竞技',
    result: resultText,
    ...gameResult
  })
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

.fighter-select {
  margin-bottom: 15px;
}

.select-title {
  font-size: 13px;
  color: var(--font-secondary);
  margin-bottom: 10px;
}

.fighter-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.fighter-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 12px;
  background: var(--panel-color);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  width: 70px;
}

.fighter-item:hover {
  border-color: var(--el-color-primary);
  background: #fffbeb;
}

.fighter-item.selected {
  border-color: var(--el-color-primary);
  background: #fef3c7;
}

.fighter-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.fighter-name {
  font-size: 11px;
  color: var(--font-color);
  font-weight: 600;
}

.fighter-style {
  font-size: 10px;
  color: var(--font-light);
}

.game-playing {
  text-align: center;
}

.fighting-arena {
  background: linear-gradient(135deg, #744210 0%, #4a3728 100%);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 15px;
}

.health-bars {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
}

.health-bar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.fighter-label {
  font-size: 12px;
  color: #e2e8f0;
  width: 30px;
}

.bar-bg {
  flex: 1;
  height: 16px;
  background: #1a202c;
  border-radius: 8px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  transition: width 0.3s;
}

.bar-fill.player {
  background: linear-gradient(90deg, var(--success-color) 0%, #85ce61 100%);
}

.bar-fill.enemy {
  background: linear-gradient(90deg, var(--error-color) 0%, #f78989 100%);
}

.health-text {
  font-size: 12px;
  color: #e2e8f0;
  width: 40px;
  text-align: right;
}

.arena-display {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 20px 0;
}

.fighter-sprite {
  font-size: 48px;
}

.vs-text {
  font-size: 24px;
  font-weight: 600;
  color: #ffd700;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
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
}
</style>