<template>
  <el-dialog v-model="visible" title="FPS射击" width="550px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">快节奏射击对决，考验你的反应速度</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 200 }}元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%">开始比赛</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <!-- 顶部状态栏 -->
        <div class="top-bar">
          <div class="health-bar-container">
            <span class="bar-label">生命值</span>
            <div class="health-bar">
              <div class="health-fill" :style="{ width: playerHealth + '%', background: getHealthColor(playerHealth) }"></div>
            </div>
            <span class="health-text">{{ playerHealth }}%</span>
          </div>
          <div class="timer">{{ formatTime(timeLeft) }}</div>
          <div class="health-bar-container">
            <span class="health-text">{{ enemyHealth }}%</span>
            <div class="health-bar">
              <div class="health-fill enemy" :style="{ width: enemyHealth + '%', background: getHealthColor(enemyHealth) }"></div>
            </div>
            <span class="bar-label">对手</span>
          </div>
        </div>

        <!-- 射击场地 -->
        <div class="fps-arena" @click="arenaClick">
          <div class="arena-bg">
            <div class="grid-overlay"></div>

            <!-- 玩家 -->
            <div class="character player" :style="{ left: playerX + '%', top: playerY + '%' }">
              <div class="character-sprite">🔫</div>
              <div class="health-mini">{{ playerHealth }}%</div>
              <div class="character-name">你</div>
            </div>

            <!-- 对手 -->
            <div class="character enemy" :style="{ left: enemyX + '%', top: enemyY + '%' }">
              <div class="character-sprite">👾</div>
              <div class="health-mini">{{ enemyHealth }}%</div>
              <div class="character-name">对手</div>
            </div>

            <!-- 目标（射击场） -->
            <div
              v-for="target in targets"
              :key="target.id"
              class="target"
              :class="{ hit: target.hit, special: target.special }"
              :style="{ left: target.x + '%', top: target.y + '%' }"
              @click.stop="shootTarget(target)"
            >
              <span class="target-icon">{{ target.icon }}</span>
              <span class="target-points">+{{ target.points }}</span>
            </div>

            <!-- 弹孔效果 -->
            <div
              v-for="hole in bulletHoles"
              :key="hole.id"
              class="bullet-hole"
              :style="{ left: hole.x + '%', top: hole.y + '%' }"
            >
              💥
            </div>

            <!-- 浮动伤害数字 -->
            <div
              v-for="damage in damageNumbers"
              :key="damage.id"
              class="damage-number"
              :class="{ heal: damage.heal, enemy: damage.enemy }"
              :style="{ left: damage.x + '%', top: damage.y + '%' }"
            >
              {{ damage.heal ? '+' : '-' }}{{ damage.value }}
            </div>
          </div>
        </div>

        <!-- 统计面板 -->
        <div class="stats-panel">
          <div class="stat-item">
            <span class="stat-icon">🎯</span>
            <span class="stat-label">击杀</span>
            <span class="stat-value">{{ playerKills }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-icon">💀</span>
            <span class="stat-label">被杀</span>
            <span class="stat-value">{{ playerDeaths }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-icon">🔥</span>
            <span class="stat-label">连击</span>
            <span class="stat-value">{{ combo }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-icon">💯</span>
            <span class="stat-label">分数</span>
            <span class="stat-value">{{ score }}</span>
          </div>
        </div>

        <!-- 武器和操作 -->
        <div class="weapon-panel">
          <div class="weapon-info">
            <span class="weapon-icon">{{ currentWeapon.icon }}</span>
            <span class="weapon-name">{{ currentWeapon.name }}</span>
            <span class="weapon-ammo">{{ ammo }}/{{ maxAmmo }}</span>
          </div>
          <div class="action-buttons">
            <el-button @click="reload" :disabled="processing || ammo === maxAmmo" size="small">🔄 换弹</el-button>
            <el-button @click="useSpecial" :disabled="processing || !specialReady" type="warning" size="small">⚡ 必杀</el-button>
            <el-button @click="shoot" :disabled="processing" type="primary" size="small">🔫 射击</el-button>
          </div>
        </div>

        <!-- 击杀日志 -->
        <div class="kill-log">
          <div v-for="log in killLog" :key="log.id" class="log-entry" :class="log.type">
            {{ log.text }}
          </div>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-row">击杀数：{{ playerKills }}</div>
          <div class="stat-row">被击杀：{{ playerDeaths }}</div>
          <div class="stat-row">最高连击：{{ maxCombo }}</div>
          <div class="stat-row">最终得分：{{ score }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useCleanupTasks } from '@/src/composables/useCleanupTasks'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'fps', name: 'FPS射击', entryCost: 200 })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetLocal()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const { gameStarted, gameEnded, processing, startGame, endGame, reset } = useMiniGameBase(props.config)
const { cleanup, clearManagedTimer, setManagedInterval, setManagedTimeout } = useCleanupTasks()

const weapons = [
  { id: 1, name: '手枪', icon: '🔫', damage: 10, fireRate: 1 },
  { id: 2, name: '步枪', icon: '🔫', damage: 15, fireRate: 1.2 },
  { id: 3, name: '霰弹枪', icon: '💥', damage: 25, fireRate: 0.8 },
  { id: 4, name: '狙击枪', icon: '🎯', damage: 40, fireRate: 0.6 }
]

const playerHealth = ref(100)
const enemyHealth = ref(100)
const playerKills = ref(0)
const playerDeaths = ref(0)
const timeLeft = ref(30)
const score = ref(0)
const combo = ref(0)
const maxCombo = ref(0)
const ammo = ref(30)
const maxAmmo = ref(30)
const specialReady = ref(false)

const playerX = ref(25)
const playerY = ref(75)
const enemyX = ref(75)
const enemyY = ref(25)

const targets = ref([])
const bulletHoles = ref([])
const damageNumbers = ref([])
const killLog = ref([])
const currentWeapon = ref(weapons[0])

let targetId = 0
let holeId = 0
let damageId = 0
let logId = 0
let gameTimer = null

const resetLocal = () => {
  cleanup()
  gameTimer = null
  gameStarted.value = false
  gameEnded.value = false
  processing.value = false
  playerHealth.value = 100
  enemyHealth.value = 100
  playerKills.value = 0
  playerDeaths.value = 0
  timeLeft.value = 30
  score.value = 0
  combo.value = 0
  maxCombo.value = 0
  ammo.value = 30
  specialReady.value = false
  playerX.value = 25
  playerY.value = 75
  enemyX.value = 75
  enemyY.value = 25
  targets.value = []
  bulletHoles.value = []
  damageNumbers.value = []
  killLog.value = []
  currentWeapon.value = weapons[0]
}

const handleClose = () => {
  visible.value = false
  reset()
  resetLocal()
}

const handleStartGame = async () => {
  const success = await startGame()
  if (!success) return

  // 开始计时
  gameTimer = setManagedInterval(() => {
    timeLeft.value--

    // 对手行动
    enemyAction()

    // 生成目标
    if (Math.random() < 0.3) {
      spawnTarget()
    }

    if (timeLeft.value <= 0 || playerHealth.value <= 0 || enemyHealth.value <= 0) {
      clearManagedTimer(gameTimer)
      gameTimer = null
      handleEndGame()
    }
  }, 1000)
}

const spawnTarget = () => {
  const special = Math.random() < 0.2
  const spawnedTargetId = targetId++
  targets.value.push({
    id: spawnedTargetId,
    x: Math.random() * 80 + 10,
    y: Math.random() * 60 + 20,
    icon: special ? '⭐' : '🎯',
    points: special ? 50 : 10,
    special: special,
    hit: false
  })

  // 目标5秒后消失
  setManagedTimeout(() => {
    const index = targets.value.findIndex(t => t.id === spawnedTargetId)
    if (index !== -1 && !targets.value[index].hit) {
      targets.value.splice(index, 1)
      combo.value = 0
    }
  }, 5000)
}

const arenaClick = (e) => {
  if (processing.value) return

  const rect = e.currentTarget.getBoundingClientRect()
  const x = ((e.clientX - rect.left) / rect.width) * 100
  const y = ((e.clientY - rect.top) / rect.height) * 100

  const currentHoleId = holeId++
  bulletHoles.value.push({ id: currentHoleId, x, y })
  setManagedTimeout(() => {
    bulletHoles.value = bulletHoles.value.filter(hole => hole.id !== currentHoleId)
  }, 2000)
}

const showDamageNumber = (damageNumber) => {
  const currentDamageId = damageId++
  damageNumbers.value.push({ id: currentDamageId, ...damageNumber })
  setManagedTimeout(() => {
    damageNumbers.value = damageNumbers.value.filter(item => item.id !== currentDamageId)
  }, 1500)
}

const shootTarget = (target) => {
  if (processing.value || target.hit) return

  processing.value = true
  target.hit = true

  // 计算伤害
  const baseDamage = currentWeapon.value.damage
  const damage = baseDamage + (combo.value * 2)

  // 显示伤害
  showDamageNumber({
    x: target.x,
    y: target.y - 5,
    value: damage,
    enemy: false,
    heal: false
  })

  // 更新分数和连击
  score.value += target.points * (1 + combo.value * 0.1)
  combo.value++
  if (combo.value > maxCombo.value) maxCombo.value = combo.value

  // 对手受伤
  enemyHealth.value = Math.max(0, enemyHealth.value - damage)
  ammo.value--

  addLog(`你命中目标，造成 ${damage} 伤害！`, 'player')

  if (enemyHealth.value <= 0) {
    playerKills.value++
    enemyHealth.value = 100
    addLog('击杀对手！+100分', 'kill')
    score.value += 100
  }

  // 检查必杀
  if (combo.value >= 5) {
    specialReady.value = true
  }

  processing.value = false

  // 移除目标
  setManagedTimeout(() => {
    const index = targets.value.findIndex(t => t.id === target.id)
    if (index !== -1) targets.value.splice(index, 1)
  }, 300)
}

const shoot = () => {
  if (processing.value || ammo.value <= 0) {
    if (ammo.value <= 0) {
      ElMessage.warning('需要换弹！')
    }
    return
  }

  processing.value = true
  ammo.value--

  // 随机命中对手
  if (Math.random() < 0.4) {
    const damage = currentWeapon.value.damage
    enemyHealth.value = Math.max(0, enemyHealth.value - damage)

    showDamageNumber({
      x: enemyX.value,
      y: enemyY.value,
      value: damage,
      enemy: true,
      heal: false
    })

    combo.value++
    if (combo.value > maxCombo.value) maxCombo.value = combo.value

    addLog(`直接命中对手，造成 ${damage} 伤害！`, 'player')

    if (enemyHealth.value <= 0) {
      playerKills.value++
      enemyHealth.value = 100
      addLog('击杀对手！', 'kill')
    }
  } else {
    combo.value = 0
    addLog('未命中', 'miss')
  }

  processing.value = false
}

const reload = () => {
  if (processing.value || ammo.value === maxAmmo.value) return

  processing.value = true
  addLog('换弹中...', 'info')

  setManagedTimeout(() => {
    ammo.value = maxAmmo.value
    processing.value = false
    addLog('换弹完成', 'info')
  }, 1500)
}

const useSpecial = () => {
  if (!specialReady.value || processing.value) return

  specialReady.value = false
  combo.value = 0

  const damage = 50
  enemyHealth.value = Math.max(0, enemyHealth.value - damage)

  showDamageNumber({
    x: enemyX.value,
    y: enemyY.value,
    value: damage,
    enemy: true,
    heal: false
  })

  addLog(`必杀技！造成 ${damage} 伤害！`, 'special')

  if (enemyHealth.value <= 0) {
    playerKills.value++
    enemyHealth.value = 100
    addLog('必杀击杀！', 'kill')
  }
}

const enemyAction = () => {
  // 对手移动
  enemyX.value = Math.max(10, Math.min(90, enemyX.value + (Math.random() - 0.5) * 20))
  enemyY.value = Math.max(10, Math.min(90, enemyY.value + (Math.random() - 0.5) * 20))

  // 对手射击
  if (Math.random() < 0.4) {
    const damage = Math.floor(Math.random() * 15) + 5
    playerHealth.value = Math.max(0, playerHealth.value - damage)

    showDamageNumber({
      x: playerX.value,
      y: playerY.value,
      value: damage,
      enemy: false,
      heal: false
    })

    addLog(`被击中，受到 ${damage} 伤害`, 'enemy')

    if (playerHealth.value <= 0) {
      playerDeaths.value++
      playerHealth.value = 100
      combo.value = 0
      addLog('你被击败！重生中...', 'death')
    }
  }

  // 对手偶尔生成目标
  if (Math.random() < 0.15) {
    spawnTarget()
  }
}

const addLog = (text, type) => {
  killLog.value.unshift({
    id: logId++,
    text,
    type
  })
  if (killLog.value.length > 5) {
    killLog.value.pop()
  }
}

const formatTime = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const getHealthColor = (health) => {
  if (health > 60) return '#67c23a'
  if (health > 30) return '#e6a23c'
  return '#f56c6c'
}

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const handleEndGame = async () => {
  cleanup()
  gameTimer = null
  const isWin = playerKills.value > playerDeaths.value || score.value > 500
  const winCount = isWin ? 1 : 0

  const resultText = isWin
    ? `获胜，击杀 ${playerKills.value}`
    : '失败'

  const detail = {
    kills: playerKills.value,
    deaths: playerDeaths.value,
    score: score.value,
    maxCombo: maxCombo.value
  }

  const gameResult = await endGame(winCount, resultText, detail)

  resultIcon.value = isWin ? '🏆' : '💀'
  resultTitle.value = isWin ? '胜利！' : '失败！'
  resultReward.value = gameResult
    ? `获得 ${gameResult.cashChange} 元，名声 +${gameResult.fameChange}`
    : ''

  emit('complete', {
    game: 'FPS射击',
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

.game-playing {
  text-align: center;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border-radius: 8px;
  margin-bottom: 10px;
}

.health-bar-container {
  display: flex;
  align-items: center;
  gap: 6px;
}

.bar-label {
  font-size: 11px;
  color: #a0aec0;
}

.health-bar {
  width: 80px;
  height: 12px;
  background: #2d3748;
  border-radius: 6px;
  overflow: hidden;
}

.health-fill {
  height: 100%;
  transition: width 0.3s, background 0.3s;
}

.health-fill.enemy {
  background: linear-gradient(90deg, #f56c6c 0%, #f78989 100%);
}

.health-text {
  font-size: 11px;
  color: #e2e8f0;
  font-weight: 600;
  min-width: 35px;
  text-align: right;
}

.timer {
  font-size: 20px;
  font-weight: 600;
  color: #ffd700;
  font-family: monospace;
}

.fps-arena {
  position: relative;
  width: 100%;
  height: 250px;
  background: linear-gradient(135deg, #0f3460 0%, #16213e 100%);
  border-radius: 12px;
  overflow: hidden;
  cursor: crosshair;
  margin-bottom: 10px;
}

.arena-bg {
  position: relative;
  width: 100%;
  height: 100%;
}

.grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 25px 25px;
  pointer-events: none;
}

.character {
  position: absolute;
  transform: translate(-50%, -50%);
  text-align: center;
  transition: all 0.3s;
}

.character-sprite {
  font-size: 32px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.5));
}

.health-mini {
  font-size: 10px;
  color: #fff;
  background: rgba(0, 0, 0, 0.5);
  padding: 1px 4px;
  border-radius: 4px;
  margin-top: 2px;
}

.character-name {
  font-size: 10px;
  color: #a0aec0;
}

.target {
  position: absolute;
  transform: translate(-50%, -50%);
  cursor: pointer;
  animation: pulse 1s infinite;
  transition: all 0.2s;
}

.target:hover {
  transform: translate(-50%, -50%) scale(1.2);
}

.target.hit {
  opacity: 0;
  transform: translate(-50%, -50%) scale(0);
}

.target.special {
  z-index: 10;
}

.target-icon {
  font-size: 28px;
  display: block;
}

.target-points {
  font-size: 10px;
  color: #ffd700;
  font-weight: 600;
  display: block;
}

@keyframes pulse {
  0%, 100% { transform: translate(-50%, -50%) scale(1); }
  50% { transform: translate(-50%, -50%) scale(1.1); }
}

.bullet-hole {
  position: absolute;
  transform: translate(-50%, -50%);
  font-size: 16px;
  pointer-events: none;
  animation: fadeOut 2s forwards;
}

@keyframes fadeOut {
  0% { opacity: 1; }
  100% { opacity: 0; }
}

.damage-number {
  position: absolute;
  transform: translate(-50%, -50%);
  font-size: 16px;
  font-weight: 600;
  animation: floatUp 1.5s forwards;
  pointer-events: none;
}

.damage-number.heal {
  color: var(--success-color);
}

.damage-number.enemy {
  color: var(--error-color);
}

@keyframes floatUp {
  0% { opacity: 1; transform: translate(-50%, -50%); }
  100% { opacity: 0; transform: translate(-50%, -150%); }
}

.stats-panel {
  display: flex;
  justify-content: space-around;
  padding: 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 10px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-icon {
  font-size: 20px;
}

.stat-label {
  font-size: 11px;
  color: var(--font-light);
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
}

.weapon-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background: linear-gradient(135deg, #2d3748 0%, #1a202c 100%);
  border-radius: 8px;
  margin-bottom: 10px;
}

.weapon-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.weapon-icon {
  font-size: 24px;
}

.weapon-name {
  font-size: 13px;
  color: #e2e8f0;
  font-weight: 500;
}

.weapon-ammo {
  font-size: 12px;
  color: #ffd700;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.kill-log {
  max-height: 80px;
  overflow-y: auto;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 8px;
}

.log-entry {
  font-size: 12px;
  padding: 4px 8px;
  margin-bottom: 4px;
  border-radius: 4px;
  text-align: left;
}

.log-entry.player {
  background: #d1fae5;
  color: #065f46;
}

.log-entry.enemy {
  background: #fee2e2;
  color: #991b1b;
}

.log-entry.kill {
  background: #fef3c7;
  color: #92400e;
  font-weight: 600;
}

.log-entry.death {
  background: #e5e7eb;
  color: #374151;
}

.log-entry.miss {
  background: #f3f4f6;
  color: #6b7280;
}

.log-entry.special {
  background: #ddd6fe;
  color: #5b21b6;
  font-weight: 600;
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
  margin-bottom: 15px;
}

.result-stats {
  margin-bottom: 15px;
}

.stat-row {
  font-size: 14px;
  color: var(--font-secondary);
  margin-bottom: 5px;
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
