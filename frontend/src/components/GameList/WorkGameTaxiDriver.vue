<template>
  <el-dialog v-model="dialogVisible" title="出车接单" width="480px" :close-on-click-modal="false" @close="handleClose">
    <div class="taxi-game">
      <!-- 游戏区域 -->
      <div class="game-layout" v-if="isRunning && !showResult">
        <!-- 左侧：地图 -->
        <div class="map-section">
          <div class="map" :style="{ gridTemplateColumns: `repeat(${cols}, 20px)` }">
            <div v-for="r in rows" :key="'r-' + r" class="row">
              <div v-for="c in cols" :key="`c-${r}-${c}`" class="cell" @click="clickCell(r-1, c-1)">
                <span v-if="carPos[0] === r-1 && carPos[1] === c-1">🚕</span>
                <span v-else-if="startPos[0] === r-1 && startPos[1] === c-1">🟢</span>
                <span v-else-if="endPos[0] === r-1 && endPos[1] === c-1">🔴</span>
                <span v-else-if="getMap(r-1, c-1) === 1">🏢</span>
                <span v-else-if="getMap(r-1, c-1) === 2">🚥</span>
                <span v-else-if="getMap(r-1, c-1) === 3">🚧</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：信息面板 -->
        <div class="info-panel">
          <!-- 状态信息 -->
          <div class="info-section">
            <div class="info-title">📊 游戏状态</div>
            <div class="info-item">
              <span class="label">⏱️ 已用时间</span>
              <span class="value">{{ timeUsed }}秒</span>
            </div>
            <div class="info-item">
              <span class="label">⛽ 剩余油量</span>
              <span class="value">{{ fuel }}</span>
            </div>
            <div class="info-item">
              <span class="label">📍 当前位置</span>
              <span class="value">({{ carPos[0] + 1 }},{{ carPos[1] + 1 }})</span>
            </div>
          </div>

          <!-- 乘客需求 -->
          <div class="info-section" v-if="currentRequest">
            <div class="info-title">🎯 乘客需求</div>
            <div class="request-name">{{ currentRequest.title }}</div>
            <div class="request-desc">{{ currentRequest.effect }}</div>
          </div>

          <!-- 方向控制 -->
          <div class="control-section">
            <div class="info-title">🎮 方向控制</div>
            <div class="d-pad">
              <button class="btn" @click="move('up')">⬆️</button>
              <div class="d-row">
                <button class="btn" @click="move('left')">⬅️</button>
                <button class="btn" @click="move('down')">⬇️</button>
                <button class="btn" @click="move('right')">➡️</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 开始界面 -->
      <div class="start-screen" v-if="!isRunning && !showResult">
        <div class="icon">🚕</div>
        <div class="title">准备好接单了吗？</div>
        <div class="difficulty-select">
          <div class="diff-title">选择难度</div>
          <div class="diff-buttons">
            <button
              v-for="(config, key) in difficultyConfig"
              :key="key"
              class="diff-btn"
              :class="{ active: worklevel == key }"
              @click="worklevel = key"
            >
              {{ config.name }}
            </button>
          </div>
        </div>
        <div class="tips">
          <div class="tip">💡 从绿点到红点接送乘客</div>
          <div class="tip">⛽ 管理好油量</div>
          <div class="tip">⭐ 高评分获得更多收入</div>
        </div>
      </div>

      <!-- 结果界面 -->
      <div class="result-screen" v-if="showResult">
        <div class="stars">
          <span v-for="i in 5" :key="i" :class="{ on: i <= rating }">⭐</span>
        </div>
        <div class="text">{{ feedback }}</div>
        <div class="money">{{ earnings }}</div>
      </div>
    </div>
    <template #footer>
      <el-button @click="handleClose">{{ showResult ? '关闭' : '放弃' }}</el-button>
      <el-button v-if="!isRunning && !showResult" type="primary" @click="startGame">开始接单</el-button>
      <el-button v-if="isRunning && !showResult" type="danger" @click="endGame">结束</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { CancelMiniGame, StartMiniGame, EndMiniGame } from "@/wailsjs/go/services/App.js"

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue', 'complete'])

const gameStore = useGameStore()
const dialogVisible = ref(props.modelValue)
const rows = 13
const cols = 13
const isRunning = ref(false)
const activeSessionID = ref('')

const cityMap = ref([])
const startPos = ref([0, 0])
const endPos = ref([0, 0])
const carPos = ref([0, 0])
const trafficLights = ref(new Set())
const fuel = ref(0)
const timeUsed = ref(0)
const rating = ref(0)
const feedback = ref('')
const showResult = ref(false)
const crashCount = ref(0)
const stepCount = ref(0)
const redLightCrashCount = ref(0)
const earnings = ref('')
const currentRequest = ref(null)
const worklevel = ref(0) // 0, 1, 2

let gameTimer = null
let lightTimer = null

const passengerRequests = [
  { title: '赶时间', effect: '10秒内到达', type: 'timeLimit', value: 10 },
  { title: '怕绕路', effect: '走最短路线', type: 'shortest', value: 1 },
  { title: '怕撞车', effect: '不能撞墙', type: 'noCrash', value: 0 },
  { title: '守规矩', effect: '不闯红灯', type: 'redLight', value: 0 },
  { title: '小费', effect: '快速完成奖励', type: 'bonus', value: 20 },
]

// 难度配置
const difficultyConfig = {
  0: { fuelBonus: 10, lightInterval: 5000, obstacleRate: 0.12, name: '简单' },
  1: { fuelBonus: 5, lightInterval: 3000, obstacleRate: 0.18, name: '中等' },
  2: { fuelBonus: 2, lightInterval: 2000, obstacleRate: 0.25, name: '困难' }
}

watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val) {
    // 打开弹窗时重置游戏状态
    resetGameState()
  }
})
watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
  if (!val) cleanup()
})

const getMap = (r, c) => cityMap.value[r]?.[c] ?? -1

const generateMap = () => {
  const config = difficultyConfig[worklevel.value]

  while (true) {
    const map = Array.from({ length: rows }, () => Array(cols).fill(0))
    let cond = [0, 1, 1, 0]
    if (Math.random() > 0.5) { cond[0] = 1; cond[2] = 0 }
    if (Math.random() > 0.5) { cond[1] = 1; cond[3] = 0 }

    for (let r = 0; r < rows; r++) {
      if (r % 2 === cond[0]) {
        for (let c = 0; c < cols; c++) {
          if (c % 2 === cond[1]) map[r][c] = 1
        }
      }
    }

    trafficLights.value.clear()
    for (let r = 0; r < rows; r++) {
      if (r % 2 === cond[2]) {
        for (let c = 0; c < cols; c++) {
          if (c % 2 === cond[3]) {
            map[r][c] = 2
            trafficLights.value.add(`${r}-${c}`)
          }
        }
      }
    }

    const maxObs = Math.floor(rows * cols * config.obstacleRate * Math.random())
    let added = 0
    while (added < maxObs) {
      const r = Math.floor(Math.random() * rows)
      const c = Math.floor(Math.random() * cols)
      if (map[r][c] === 0) { map[r][c] = 3; added++ }
    }

    let start, end
    do { start = [Math.floor(Math.random() * rows), Math.floor(Math.random() * cols)] }
    while (map[start[0]][start[1]] !== 0)
    do { end = [Math.floor(Math.random() * rows), Math.floor(Math.random() * cols)] }
    while (map[end[0]][end[1]] !== 0 || (end[0] === start[0] && end[1] === start[1]))

    const [shortest, ok] = findPath(map, start, end)
    if (ok) {
      cityMap.value = map
      startPos.value = start
      endPos.value = end
      carPos.value = [...start]
      fuel.value = shortest + config.fuelBonus
      return
    }
  }
}

const findPath = (map, start, end) => {
  const visited = Array.from({ length: rows }, () => Array(cols).fill(false))
  const queue = [[start[0], start[1], 0]]
  visited[start[0]][start[1]] = true
  const dirs = [[1,0],[-1,0],[0,1],[0,-1]]

  while (queue.length) {
    const [r, c, steps] = queue.shift()
    if (r === end[0] && c === end[1]) return [steps, true]

    for (const [dr, dc] of dirs) {
      const nr = r + dr, nc = c + dc
      if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] && (map[nr][nc] === 0 || map[nr][nc] === 2)) {
        visited[nr][nc] = true
        queue.push([nr, nc, steps + 1])
      }
    }
  }
  return [-1, false]
}

const updateLights = () => {
  trafficLights.value.forEach(coord => {
    const [r, c] = coord.split('-').map(Number)
    cityMap.value[r][c] = Math.random() > 0.5 ? 2 : 0
  })
}

// 重置游戏状态
const resetGameState = () => {
  isRunning.value = false
  showResult.value = false
  timeUsed.value = 0
  crashCount.value = 0
  stepCount.value = 0
  redLightCrashCount.value = 0
}

const startGame = async () => {
  generateMap()

  timeUsed.value = 0
  showResult.value = false
  crashCount.value = 0
  stepCount.value = 0
  redLightCrashCount.value = 0
  currentRequest.value = passengerRequests[Math.floor(Math.random() * passengerRequests.length)]
  const data = await StartMiniGame('taxi', Number(worklevel.value))
  if (data.code != 200) {
    // 弹窗提示
    ElMessage.error('打工失败！'+data.msg)
    return
  }
  activeSessionID.value = String(data.sessionid || '')
  // 更新用户信息
  gameStore.applyUserInfo(data.userinfo)
  ElMessage.success('开始打工！')
  // 设置打工按钮不可见状态
  isRunning.value = true
  // 设置难度
  const config = difficultyConfig[Number(worklevel.value)]

  cleanup()
  gameTimer = setInterval(() => timeUsed.value++, 1000)
  updateLights()
  lightTimer = setInterval(updateLights, config.lightInterval)
  ElMessage.success(`🚕 新订单已生成！难度：${config.name}`)
}

const endGame = () => {
  // 提前结束游戏，视为失败
  finish(false)
}

const move = (dir) => {
  if (!isRunning.value) return
  if (fuel.value <= 0) {
    ElMessage.error('⛽ 油量耗尽！')
    finish(false)
    return
  }

  const [r, c] = carPos.value
  let nr = r, nc = c
  if (dir === 'up') nr--
  else if (dir === 'down') nr++
  else if (dir === 'left') nc--
  else if (dir === 'right') nc++

  // 边界检查
  if (nr < 0 || nr >= rows || nc < 0 || nc >= cols) {
    ElMessage.warning('⛔ 已到边界')
    return
  }

  const cell = getMap(nr, nc)

  // 检查是否碰到建筑或施工
  if (cell === 1 || cell === 3) {
    crashCount.value++
    ElMessage.warning('⛔ 不能进入')
    return
  }

  // 检查红绿灯
  if (cell === 2) {
    redLightCrashCount.value++
    ElMessage.warning('🚦 闯红灯！')
  }

  fuel.value--
  stepCount.value++
  carPos.value = [nr, nc]

  if (nr === endPos.value[0] && nc === endPos.value[1]) {
    finish(true)
  }
}

const finish = async (success) => {
  isRunning.value = false
  showResult.value = true
  cleanup()

  if (!success) {
    void cancelActiveSession()
    rating.value = 1
    feedback.value = '订单失败！😡'
    earnings.value = ''
    return
  }

  const [shortestPath] = findPath(cityMap.value, startPos.value, endPos.value)
  if (shortestPath < 0) {
    void cancelActiveSession()
    rating.value = 1
    feedback.value = '地图错误！'
    earnings.value = ''
    return
  }

  const diff = stepCount.value - shortestPath
  let stars = 5

  if (currentRequest.value) {
    const req = currentRequest.value
    if (req.type === 'timeLimit' && timeUsed.value > req.value) stars -= 1
    if (req.type === 'shortest' && diff > req.value) stars -= 1
    if (req.type === 'noCrash' && crashCount.value > 0) stars -= 1
    if (req.type === 'redLight' && redLightCrashCount.value > 0) stars -= 1
  }

  stars = Math.max(1, Math.min(5, stars))
  rating.value = stars

  const feedbacks = { 5: '非常满意！', 4: '满意！', 3: '还行吧。', 2: '不太满意。', 1: '差评！😡' }
  feedback.value = feedbacks[stars]

  // 计算收入
  const data = await EndMiniGame(stars)
  activeSessionID.value = ''
  if (data.code != 200) {
    ElMessage.error('打工结束！'+data.msg || '打工结束！老板跑路了！')
    return
  }
  // 更新用户信息
  gameStore.applyUserInfo(data.userinfo)
  // 弹窗提示
  ElMessage.success(`打工结束！${data.msg}`)
  emit('complete', { game: '出车接单', rating: stars, earnings: data.msg })
}

const clickCell = (r, c) => {
  const [cr, cc] = carPos.value
  if (r === cr && c === cc + 1) move('right')
  else if (r === cr && c === cc - 1) move('left')
  else if (r === cr - 1 && c === cc) move('up')
  else if (r === cr + 1 && c === cc) move('down')
}

const handleClose = () => {
  // 关闭弹窗，清理状态
  void cancelActiveSession()
  dialogVisible.value = false
}

const cancelActiveSession = async () => {
  const sessionID = activeSessionID.value
  activeSessionID.value = ''
  if (!sessionID) return
  try {
    await CancelMiniGame(sessionID)
  } catch (error) {
    console.error('取消出租车打工失败', error)
  }
}

const cleanup = () => {
  if (gameTimer) clearInterval(gameTimer)
  if (lightTimer) clearInterval(lightTimer)
  gameTimer = null
  lightTimer = null
}

onUnmounted(() => {
  cleanup()
  void cancelActiveSession()
})
</script>

<style scoped>
.taxi-game {
  padding: 5px 0;
}

.game-layout {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.map-section {
  flex-shrink: 0;
}

.map {
  display: grid;
  gap: 2px;
}

.row {
  display: contents;
}

.cell {
  width: 20px;
  height: 20px;
  border: 1px solid var(--border-color);
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  cursor: pointer;
  background: var(--panel-color);
  transition: all 0.15s;
}

.info-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 140px;
}

.info-section {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 8px;
}

.info-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 6px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 0;
  font-size: 11px;
}

.label {
  color: var(--font-secondary);
}

.value {
  font-weight: 600;
  color: var(--font-color);
}

.request-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--warning-color);
  margin-bottom: 3px;
}

.request-desc {
  font-size: 11px;
  color: var(--font-secondary);
  line-height: 1.3;
}

.control-section {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 8px;
}

.d-pad {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
}

.btn {
  width: 32px;
  height: 32px;
  border: none;
  background: var(--select-color);
  border-radius: 5px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.btn:hover {
  transform: scale(1.05);
}

.btn:active {
  transform: scale(0.95);
}

.d-row {
  display: flex;
  gap: 3px;
}

.start-screen {
  text-align: center;
  padding: 20px 15px;
}

.icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 12px;
}

.difficulty-select {
  margin-bottom: 12px;
}

.diff-title {
  font-size: 12px;
  color: var(--font-secondary);
  margin-bottom: 6px;
}

.diff-buttons {
  display: flex;
  justify-content: center;
  gap: 6px;
}

.diff-btn {
  padding: 5px 12px;
  border: 1px solid var(--border-color);
  background: var(--panel-color);
  border-radius: 5px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.diff-btn.active {
  background: var(--select-color);
  border-color: var(--select-border-color);
}

.tips {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: center;
}

.tip {
  font-size: 11px;
  color: var(--font-secondary);
  padding: 5px 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
}

.result-screen {
  text-align: center;
  padding: 25px;
}

.stars {
  font-size: 32px;
  margin-bottom: 12px;
}

.stars span {
  opacity: 0.25;
}

.stars span.on {
  opacity: 1;
}

.text {
  font-size: 16px;
  color: var(--font-secondary);
  margin-bottom: 15px;
}

.money {
  font-size: 26px;
  color: var(--success-color);
  font-weight: 600;
}
</style>
