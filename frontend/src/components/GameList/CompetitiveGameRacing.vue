<template>
  <el-dialog v-model="visible" title="赛车竞速" width="520px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">速度与激情的较量，谁是赛道之王</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 800 }}元</span>
        </div>
        <div class="car-select">
          <div class="select-title">选择你的赛车：</div>
          <div class="car-list">
            <div v-for="car in cars" :key="car.id" class="car-item" :class="{ selected: selectedCar === car.id }" @click="selectedCar = car.id">
              <span class="car-icon">{{ car.icon }}</span>
              <span class="car-name">{{ car.name }}</span>
              <span class="car-type">{{ car.type }}</span>
              <div class="car-stats">
                <div class="stat-bar">
                  <span class="stat-label">速度</span>
                  <div class="stat-bar-bg">
                    <div class="stat-fill" :style="{ width: car.speed + '%' }"></div>
                  </div>
                </div>
                <div class="stat-bar">
                  <span class="stat-label">加速</span>
                  <div class="stat-bar-bg">
                    <div class="stat-fill" :style="{ width: car.accel + '%' }"></div>
                  </div>
                </div>
                <div class="stat-bar">
                  <span class="stat-label">操控</span>
                  <div class="stat-bar-bg">
                    <div class="stat-fill" :style="{ width: car.handle + '%' }"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%; margin-top: 15px" :disabled="!selectedCar">开始比赛</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <!-- 顶部状态栏 -->
        <div class="top-bar">
          <div class="position-display">
            <span class="position-label">当前排名</span>
            <span class="position-number">{{ getPositionText(playerPosition) }}</span>
          </div>
          <div class="lap-display">
            <span class="lap-label">圈数</span>
            <span class="lap-number">{{ lap }}/3</span>
          </div>
          <div class="time-display">
            <span class="time-label">用时</span>
            <span class="time-number">{{ formatTime(elapsedTime) }}</span>
          </div>
          <div class="speed-display">
            <span class="speed-label">速度</span>
            <span class="speed-number">{{ Math.round(currentSpeed) }} km/h</span>
          </div>
        </div>

        <!-- 赛道 -->
        <div class="racing-track">
          <!-- 赛道背景 -->
          <div class="track-bg">
            <!-- 弯道提示 -->
            <div v-if="upcomingTurn" class="turn-warning" :class="upcomingTurn.type">
              {{ upcomingTurn.text }}
            </div>

            <!-- 赛道路面 -->
            <div class="track-road">
              <div class="road-lines"></div>
              <div class="road-markings">
                <div v-for="i in 20" :key="i" class="road-mark" :style="{ left: (i * 5) + '%' }"></div>
              </div>
            </div>

            <!-- 玩家赛车 -->
            <div class="racer player" :style="{ left: playerLane + '%', bottom: playerProgress + '%' }">
              <div class="car-sprite player-car" :style="{ transform: `rotate(${carTilt}deg)` }">
                {{ selectedCarObj?.icon }}
              </div>
              <div class="speed-lines" v-if="isBoosting"></div>
              <div class="player-label">你</div>
            </div>

            <!-- 对手赛车 -->
            <div v-for="opponent in opponents" :key="opponent.id" class="racer opponent"
                 :style="{ left: opponent.lane + '%', bottom: opponent.progress + '%' }">
              <div class="car-sprite opponent-car">
                {{ opponent.icon }}
              </div>
              <div class="opponent-label">{{ opponent.name }}</div>
            </div>

            <!-- 道具/障碍物 -->
            <div v-for="item in trackItems" :key="item.id" class="track-item" :class="item.type"
                 :style="{ left: item.lane + '%', bottom: item.progress + '%' }">
              {{ item.icon }}
            </div>

            <!-- 加速效果 -->
            <div class="boost-effect" v-if="isBoosting">
              <div class="boost-particle" v-for="i in 5" :key="i"></div>
            </div>
          </div>

          <!-- 迷你地图 -->
          <div class="minimap">
            <div class="minimap-track">
              <div class="minimap-progress" :style="{ strokeDashoffset: (100 - playerProgress / 3) + '%' }"></div>
            </div>
            <div class="minimap-players">
              <div class="minimap-dot player-dot" :style="{ bottom: (playerProgress / 3) + '%' }"></div>
              <div v-for="opp in opponents" :key="opp.id" class="minimap-dot opponent-dot"
                   :style="{ bottom: (opp.progress / 3) + '%' }"></div>
            </div>
          </div>
        </div>

        <!-- 操作面板 -->
        <div class="control-panel">
          <!-- 能量槽 -->
          <div class="energy-bar">
            <span class="energy-label">氮气</span>
            <div class="energy-bg">
              <div class="energy-fill" :style="{ width: boostEnergy + '%' }"></div>
            </div>
            <span class="energy-text">{{ boostEnergy }}%</span>
          </div>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <el-button @click="changeLane(-1)" :disabled="processing || playerLane <= 15" size="small">
              ⬅️ 左变道
            </el-button>
            <el-button @click="accelerate" :disabled="processing" type="primary" size="small">
              🚀 加速
            </el-button>
            <el-button @click="useBoost" :disabled="processing || boostEnergy < 20" type="warning" size="small">
              ⚡ 氮气加速
            </el-button>
            <el-button @click="brake" :disabled="processing" type="info" size="small">
              🛑 刹车
            </el-button>
            <el-button @click="changeLane(1)" :disabled="processing || playerLane >= 85" size="small">
              右变道 ➡️
            </el-button>
          </div>
        </div>

        <!-- 实时排名 -->
        <div class="leaderboard">
          <div class="leaderboard-title">实时排名</div>
          <div class="leaderboard-list">
            <div v-for="(racer, index) in sortedRacers" :key="racer.id"
                 class="leaderboard-item" :class="{ player: racer.isPlayer }">
              <span class="rank">{{ getPositionText(index + 1) }}</span>
              <span class="racer-name">{{ racer.name }}</span>
              <span class="racer-progress">{{ Math.round(racer.progress / 3) }}%</span>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">最终排名：{{ getPositionText(finalPosition) }}</div>
          <div class="stat-item">完成时间：{{ formatTime(finalTime) }}</div>
          <div class="stat-item">最高速度：{{ maxSpeedReached }} km/h</div>
          <div class="stat-item">使用加速：{{ boostCount }} 次</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'racing', name: '赛车竞速', entryCost: 400 })
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

const cars = [
  { id: 1, name: '闪电', icon: '⚡', type: '极速型', speed: 95, accel: 85, handle: 60 },
  { id: 2, name: '风暴', icon: '🌀', type: '平衡型', speed: 80, accel: 80, handle: 80 },
  { id: 3, name: '坦克', icon: '💪', type: '稳定型', speed: 65, accel: 70, handle: 95 }
]

const selectedCar = ref(null)

// 赛车状态
const playerProgress = ref(0)
const playerLane = ref(50)
const playerPosition = ref(4)
const currentSpeed = ref(0)
const maxSpeedReached = ref(0)
const boostEnergy = ref(100)
const boostCount = ref(0)
const isBoosting = ref(false)
const carTilt = ref(0)

// 比赛状态
const lap = ref(1)
const elapsedTime = ref(0)
const finalTime = ref(0)
const finalPosition = ref(4)

// 对手
const opponents = ref([])
const trackItems = ref([])

// 弯道提示
const upcomingTurn = ref(null)

// 结果
const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

let gameTimer = null
let itemId = 0

const selectedCarObj = computed(() => cars.find(c => c.id === selectedCar.value))

const sortedRacers = computed(() => {
  const allRacers = [
    { id: 0, name: '你', progress: playerProgress.value, isPlayer: true },
    ...opponents.value.map(o => ({ id: o.id, name: o.name, progress: o.progress, isPlayer: false }))
  ]
  return allRacers.sort((a, b) => b.progress - a.progress)
})

const resetGameState = () => {
  reset()
  selectedCar.value = null
  playerProgress.value = 0
  playerLane.value = 50
  playerPosition.value = 4
  currentSpeed.value = 0
  maxSpeedReached.value = 0
  boostEnergy.value = 100
  boostCount.value = 0
  isBoosting.value = false
  carTilt.value = 0
  lap.value = 1
  elapsedTime.value = 0
  finalTime.value = 0
  finalPosition.value = 4
  opponents.value = []
  trackItems.value = []
  upcomingTurn.value = null
  resultIcon.value = ''
  resultTitle.value = ''
  resultReward.value = ''
  if (gameTimer) clearInterval(gameTimer)
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const handleStartGame = async () => {
  const success = await startMiniGame()
  if (!success) return

  // 初始化对手
  opponents.value = [
    { id: 1, name: '对手A', icon: '🚗', progress: 0, lane: 20 },
    { id: 2, name: '对手B', icon: '🚕', progress: 0, lane: 40 },
    { id: 3, name: '对手C', icon: '🚙', progress: 0, lane: 60 },
    { id: 4, name: '对手D', icon: '🚌', progress: 0, lane: 80 }
  ]

  // 游戏主循环
  gameTimer = setInterval(() => {
    elapsedTime.value++

    // 更新速度
    const car = selectedCarObj.value
    const naturalDecel = 0.5
    const targetSpeed = car.speed * (0.8 + Math.random() * 0.4)

    if (currentSpeed.value < targetSpeed) {
      currentSpeed.value = Math.min(targetSpeed, currentSpeed.value + car.accel * 0.2)
    } else {
      currentSpeed.value = Math.max(0, currentSpeed.value - naturalDecel)
    }

    maxSpeedReached.value = Math.max(maxSpeedReached.value, currentSpeed.value)

    // 更新玩家进度
    const speedBonus = isBoosting.value ? 1.5 : 1
    playerProgress.value += (currentSpeed.value / 100) * 0.3 * speedBonus

    // 更新对手
    opponents.value.forEach(opp => {
      const oppSpeed = 50 + Math.random() * 40
      opp.progress += (oppSpeed / 100) * 0.28
      opp.lane = Math.max(15, Math.min(85, opp.lane + (Math.random() - 0.5) * 5))
    })

    // 生成道具/障碍物
    if (Math.random() < 0.1) {
      spawnItem()
    }

    // 弯道提示
    if (Math.floor(elapsedTime.value) % 10 === 0 && Math.floor(elapsedTime.value) > 0) {
      showTurnWarning()
    }

    // 更新排名
    updatePosition()

    // 更新圈数
    if (playerProgress.value >= lap.value * 100 && lap.value < 3) {
      lap.value++
    }

    // 能量恢复
    if (boostEnergy.value < 100 && !isBoosting.value) {
      boostEnergy.value += 0.3
    }

    // 检查游戏结束
    if (playerProgress.value >= 300) {
      finalTime.value = elapsedTime.value
      clearInterval(gameTimer)
      handleEndGame(true)
    } else if (opponents.value.some(o => o.progress >= 300)) {
      clearInterval(gameTimer)
      handleEndGame(false)
    }
  }, 100)
}

const accelerate = () => {
  if (processing.value) return

  processing.value = true
  const car = selectedCarObj.value
  currentSpeed.value = Math.min(car.speed * 1.2, currentSpeed.value + car.accel * 0.5)

  setTimeout(() => {
    processing.value = false
  }, 200)
}

const brake = () => {
  if (processing.value) return

  processing.value = true
  currentSpeed.value = Math.max(0, currentSpeed.value - 30)

  setTimeout(() => {
    processing.value = false
  }, 200)
}

const useBoost = () => {
  if (processing.value || boostEnergy.value < 20) return

  isBoosting.value = true
  boostEnergy.value -= 20
  boostCount.value++

  setTimeout(() => {
    isBoosting.value = false
  }, 2000)
}

const changeLane = (direction) => {
  if (processing.value) return

  const newLane = playerLane.value + (direction * 15)
  if (newLane >= 15 && newLane <= 85) {
    playerLane.value = newLane
    carTilt.value = direction * 10

    setTimeout(() => {
      carTilt.value = 0
    }, 300)
  }

  // 检查是否有道具
  checkItemCollision()
}

const spawnItem = () => {
  const types = [
    { type: 'boost', icon: '⚡', effect: 'speed' },
    { type: 'coin', icon: '💰', effect: 'gold' },
    { type: 'obstacle', icon: '🚧', effect: 'slow' }
  ]

  const item = types[Math.floor(Math.random() * types.length)]
  trackItems.value.push({
    id: itemId++,
    ...item,
    lane: 15 + Math.random() * 70,
    progress: playerProgress.value + 30
  })

  // 道具消失
  setTimeout(() => {
    const index = trackItems.value.findIndex(i => i.id === itemId - 1)
    if (index !== -1) trackItems.value.splice(index, 1)
  }, 5000)
}

const checkItemCollision = () => {
  trackItems.value.forEach((item, index) => {
    if (Math.abs(item.lane - playerLane.value) < 10 && Math.abs(item.progress - playerProgress.value) < 5) {
      // 碰撞
      switch (item.effect) {
        case 'speed':
          currentSpeed.value = Math.min(200, currentSpeed.value + 30)
          break
        case 'gold':
          // 什么也不做，只是视觉效果
          break
        case 'slow':
          currentSpeed.value = Math.max(0, currentSpeed.value - 20)
          break
      }
      trackItems.value.splice(index, 1)
    }
  })
}

const showTurnWarning = () => {
  const turns = [
    { text: '左弯道', type: 'left' },
    { text: '右弯道', type: 'right' },
    { text: '急转弯！', type: 'sharp' }
  ]

  const turn = turns[Math.floor(Math.random() * turns.length)]
  upcomingTurn.value = turn

  setTimeout(() => {
    upcomingTurn.value = null
  }, 3000)
}

const updatePosition = () => {
  const allProgress = [
    playerProgress.value,
    ...opponents.value.map(o => o.progress)
  ].sort((a, b) => b - a)

  playerPosition.value = allProgress.indexOf(playerProgress.value) + 1
}

const getPositionText = (pos) => {
  const positions = ['🥇 第1', '🥈 第2', '🥉 第3', '第4', '第5']
  return positions[pos - 1] || `${pos}位`
}

const formatTime = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  const ms = Math.floor((seconds % 1) * 100)
  return `${mins}:${secs.toString().padStart(2, '0')}.${ms.toString().padStart(2, '0')}`
}

const handleEndGame = async (isPlayerWin) => {
  finalPosition.value = playerPosition.value

  const isWin = isPlayerWin && playerPosition.value <= 2
  const winCount = isWin ? 1 : 0

  const resultText = isWin
    ? `${getPositionText(finalPosition.value)}，用时 ${formatTime(finalTime.value)}`
    : getPositionText(finalPosition.value)

  const detail = {
    finalPosition: finalPosition.value,
    finalTime: finalTime.value,
    maxSpeed: maxSpeedReached.value,
    boostCount: boostCount.value
  }

  const gameResult = await endMiniGame(winCount, resultText, detail)

  resultIcon.value = isWin ? '🏁' : '🏳️'
  resultTitle.value = `比赛完成！${getPositionText(finalPosition.value)}`
  resultReward.value = gameResult
    ? `获得 ${gameResult.cashChange} 元，名声 +${gameResult.fameChange}`
    : ''

  emit('complete', {
    game: '赛车竞速',
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

.car-select {
  margin-bottom: 15px;
}

.select-title {
  font-size: 13px;
  color: var(--font-secondary);
  margin-bottom: 10px;
}

.car-list {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}

.car-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 12px;
  background: var(--panel-color);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  width: 85px;
}

.car-item:hover {
  border-color: var(--el-color-primary);
  background: #fffbeb;
}

.car-item.selected {
  border-color: var(--el-color-primary);
  background: #fef3c7;
}

.car-icon {
  font-size: 32px;
  margin-bottom: 6px;
}

.car-name {
  font-size: 12px;
  color: var(--font-color);
  font-weight: 600;
  margin-bottom: 2px;
}

.car-type {
  font-size: 10px;
  color: var(--font-light);
  margin-bottom: 8px;
}

.car-stats {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 100%;
}

.stat-bar {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-label {
  font-size: 9px;
  color: var(--font-secondary);
  min-width: 30px;
}

.stat-bar-bg {
  flex: 1;
  height: 6px;
  background: #e0e0e0;
  border-radius: 3px;
  overflow: hidden;
}

.stat-fill {
  height: 100%;
  background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%);
  transition: width 0.3s;
}

.game-playing {
  text-align: center;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 10px;
  background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
  border-radius: 8px;
  margin-bottom: 10px;
  color: #fff;
}

.position-display, .lap-display, .time-display, .speed-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.position-label, .lap-label, .time-label, .speed-label {
  font-size: 10px;
  opacity: 0.8;
}

.position-number {
  font-size: 14px;
  font-weight: 600;
}

.lap-number, .time-number, .speed-number {
  font-size: 12px;
  font-weight: 500;
  font-family: monospace;
}

.racing-track {
  position: relative;
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #374151 0%, #1f2937 100%);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 10px;
}

.track-bg {
  position: relative;
  width: 100%;
  height: 100%;
}

.turn-warning {
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  z-index: 100;
  animation: blink 0.5s infinite;
}

.turn-warning.left {
  background: rgba(59, 130, 246, 0.9);
  color: #fff;
}

.turn-warning.right {
  background: rgba(239, 68, 68, 0.9);
  color: #fff;
}

.turn-warning.sharp {
  background: rgba(245, 158, 11, 0.9);
  color: #fff;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.track-road {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, #4b5563 0%, #6b7280 50%, #4b5563 100%);
}

.road-lines {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 100%;
  background: repeating-linear-gradient(
    to bottom,
    #fff 0px,
    #fff 20px,
    transparent 20px,
    transparent 40px
  );
}

.road-markings {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.road-mark {
  position: absolute;
  width: 3px;
  height: 30px;
  background: rgba(255, 255, 255, 0.3);
}

.racer {
  position: absolute;
  transform: translate(-50%, 50%);
  transition: left 0.3s, bottom 0.1s;
}

.car-sprite {
  font-size: 36px;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.5));
  transition: transform 0.3s;
}

.player-car {
  z-index: 10;
}

.speed-lines {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  width: 40px;
  height: 20px;
  background: linear-gradient(to top, rgba(255, 255, 255, 0.6), transparent);
  animation: speedLines 0.2s infinite;
}

@keyframes speedLines {
  0% { height: 15px; }
  100% { height: 25px; }
}

.player-label, .opponent-label {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  white-space: nowrap;
  padding: 2px 6px;
  border-radius: 4px;
}

.player-label {
  background: rgba(var(--success-color-rgb), 0.9);
  color: #fff;
}

.opponent-label {
  background: rgba(var(--error-color-rgb), 0.9);
  color: #fff;
}

.track-item {
  position: absolute;
  transform: translate(-50%, 50%);
  font-size: 24px;
  z-index: 5;
  animation: float 2s infinite ease-in-out;
}

@keyframes float {
  0%, 100% { transform: translate(-50%, 50%) translateY(0); }
  50% { transform: translate(-50%, 50%) translateY(-5px); }
}

.boost-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.boost-particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: #fbbf24;
  border-radius: 50%;
  animation: boostParticle 0.5s infinite;
}

.boost-particle:nth-child(1) { left: 45%; bottom: 20%; animation-delay: 0s; }
.boost-particle:nth-child(2) { left: 55%; bottom: 25%; animation-delay: 0.1s; }
.boost-particle:nth-child(3) { left: 50%; bottom: 15%; animation-delay: 0.2s; }
.boost-particle:nth-child(4) { left: 48%; bottom: 22%; animation-delay: 0.3s; }
.boost-particle:nth-child(5) { left: 52%; bottom: 18%; animation-delay: 0.4s; }

@keyframes boostParticle {
  0% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(-20px); }
}

.minimap {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 80px;
  height: 80px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 8px;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.minimap-track {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid rgba(255, 255, 255, 0.3);
}

.minimap-progress {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid var(--success-color);
  border-style: dashed;
  transform: rotate(-90deg);
  transition: stroke-dashoffset 0.3s;
}

.minimap-players {
  position: absolute;
  top: 5px;
  left: 5px;
  right: 5px;
  bottom: 5px;
}

.minimap-dot {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  transition: bottom 0.3s;
}

.player-dot {
  background: var(--success-color);
  box-shadow: 0 0 4px var(--success-color);
}

.opponent-dot {
  background: var(--error-color);
}

.control-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 10px;
}

.energy-bar {
  display: flex;
  align-items: center;
  gap: 8px;
}

.energy-label {
  font-size: 11px;
  font-weight: 600;
  color: #f59e0b;
  min-width: 35px;
}

.energy-bg {
  flex: 1;
  height: 12px;
  background: #e5e7eb;
  border-radius: 6px;
  overflow: hidden;
}

.energy-fill {
  height: 100%;
  background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%);
  transition: width 0.3s;
}

.energy-text {
  font-size: 11px;
  font-weight: 600;
  color: #f59e0b;
  min-width: 35px;
  text-align: right;
}

.action-buttons {
  display: flex;
  gap: 5px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-buttons :deep(.el-button) {
  padding: 5px 8px;
  font-size: 12px;
}

.leaderboard {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 8px;
}

.leaderboard-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 8px;
}

.leaderboard-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.leaderboard-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 8px;
  background: #fff;
  border-radius: 4px;
  font-size: 12px;
}

.leaderboard-item.player {
  background: #fef3c7;
  border: 1px solid #f59e0b;
}

.rank {
  font-weight: 600;
  min-width: 60px;
}

.racer-name {
  flex: 1;
  text-align: left;
}

.racer-progress {
  font-weight: 600;
  color: #f59e0b;
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
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 15px;
  margin-bottom: 15px;
}

.stat-item {
  font-size: 13px;
  color: var(--font-secondary);
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
