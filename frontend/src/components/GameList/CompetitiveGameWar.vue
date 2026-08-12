<template>
  <el-dialog v-model="visible" title="小小大战争" width="720px" @close="handleClose">
    <div class="game-dialog">
      <!-- 游戏介绍 -->
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">指挥你的军队，击败敌人基地！</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 300 }}元</span>
        </div>
        <div class="intro-rules">
          <div class="rule-item">🎖️ 点击单位选中，绿色区域移动</div>
          <div class="rule-item">⚔️ 移动后可攻击红色区域的敌人</div>
          <div class="rule-item">🏭 占领建筑可生产更多单位</div>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%">开始战斗</el-button>
      </div>

      <!-- 游戏进行中 -->
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <div class="player-info">
            <span class="player-name">🎖️ 玩家</span>
            <span class="player-gold">💰 {{ gold[1] }}</span>
          </div>
          <div class="game-info-inline">
            <span class="turn-indicator" :class="{ enemy: currentPlayer === 2 }">
              {{ currentPlayer === 1 ? '你的回合' : 'AI行动中...' }}
            </span>
            <span class="move-count">回合 {{ turn }}</span>
          </div>
          <div class="player-info">
            <span class="player-name">🤖 AI</span>
            <span class="player-gold">💰 {{ gold[2] }}</span>
          </div>
        </div>

        <div class="map-container">
          <canvas ref="canvasRef" :width="MAP_W * TILE" :height="MAP_H * TILE" @click="handleMapClick" />
        </div>

        <div class="game-status">
          <div v-if="selectedInfo" class="selected-info">{{ selectedInfo }}</div>
          <div v-else class="status-text">点击单位选择</div>
        </div>

        <div class="action-buttons">
          <el-button @click="showBuildMenu = true" :disabled="currentPlayer !== 1">🏭 生产</el-button>
          <el-button @click="clearSelection" :disabled="!selected">❌ 取消</el-button>
          <el-button type="primary" @click="endTurn" :disabled="currentPlayer !== 1">⏭️ 结束回合</el-button>
        </div>

        <!-- 生产菜单 -->
        <el-dialog v-model="showBuildMenu" title="生产单位" width="300px" append-to-body>
          <div class="build-options">
            <div
              v-for="(unit, key) in UNIT_TYPES"
              :key="key"
              class="build-option"
              :class="{ disabled: gold[1] < unit.cost }"
              @click="produceUnit(key)"
            >
              <span class="unit-emoji">{{ unit.emoji }}</span>
              <span class="unit-name">{{ unit.name }}</span>
              <span class="unit-cost">{{ unit.cost }}💰</span>
            </div>
          </div>
        </el-dialog>
      </div>

      <!-- 游戏结束 -->
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">总回合：{{ turn }}</div>
          <div class="stat-item">消灭单位：{{ killedUnits }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  config: {
    type: Object,
    default: () => ({ id: 'war', name: '小小大战争', entryCost: 300 })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) resetGameState()
})
watch(visible, (val) => emit('update:modelValue', val))

// 游戏配置
const TILE = 48
const MAP_W = 12
const MAP_H = 8

const UNIT_TYPES = {
  SOLDIER: { name: '士兵', hp: 100, atk: 30, def: 10, move: 3, range: 1, cost: 50, emoji: '🎖️' },
  ARCHER: { name: '弓箭手', hp: 80, atk: 40, def: 5, move: 2, range: 3, cost: 60, emoji: '🏹' },
  TANK: { name: '坦克', hp: 200, atk: 60, def: 30, move: 2, range: 2, cost: 100, emoji: '🛡️' }
}

const COLORS = { 1: '#4a90d9', 2: '#d94a4a' }
const TERRAIN = { GRASS: '#5a8a5a', FOREST: '#2d5a2d', MOUNTAIN: '#6a6a5a', WATER: '#3a5a8a' }

// 游戏状态
const canvasRef = ref(null)
let ctx = null

const { gameStarted, gameEnded, processing, startGame: startMiniGame, endGame: endMiniGame, reset } = useMiniGameBase(props.config)

const currentPlayer = ref(1)
const turn = ref(1)
const gold = reactive({ 1: 150, 2: 150 })
const map = ref([])
const units = ref([])
const buildings = ref([])
const selected = ref(null)
const moveTiles = ref([])
const attackTiles = ref([])
const showBuildMenu = ref(false)
const busy = ref(false)

const killedUnits = ref(0)
const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const selectedInfo = computed(() => {
  if (!selected.value) return ''
  if (selected.value.type === 'BASE' || selected.value.type === 'BARRACKS') {
    return `🏰 ${selected.value.type === 'BASE' ? '基地' : '兵营'} HP:${selected.value.hp}/${selected.value.maxHp}`
  }
  const c = UNIT_TYPES[selected.value.type]
  return `${c.emoji} ${c.name} HP:${selected.value.hp}/${selected.value.maxHp} ⚔️${c.atk}`
})

const resetGameState = () => {
  reset()
  currentPlayer.value = 1
  turn.value = 1
  gold[1] = 150
  gold[2] = 150
  map.value = []
  units.value = []
  buildings.value = []
  selected.value = null
  moveTiles.value = []
  attackTiles.value = []
  killedUnits.value = 0
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

  nextTick(() => {
    ctx = canvasRef.value.getContext('2d')
    initMap()
    ElMessage.info('战斗开始！点击蓝色单位行动')
  })
}

const initMap = () => {
  // 生成地图
  for (let y = 0; y < MAP_H; y++) {
    map.value[y] = []
    for (let x = 0; x < MAP_W; x++) {
      let t = 'GRASS'
      const r = Math.random()
      if (r < 0.08) t = 'FOREST'
      else if (r < 0.12) t = 'MOUNTAIN'
      else if (r < 0.15) t = 'WATER'
      if ((x < 3 && y < 3) || (x > MAP_W - 4 && y > MAP_H - 4)) t = 'GRASS'
      map.value[y][x] = { terrain: t, unit: null, building: null }
    }
  }

  // 玩家单位
  addUnit(1, 1, 'SOLDIER', 1)
  addUnit(2, 1, 'ARCHER', 1)
  addBuilding(1, 0, 'BASE', 1)
  addBuilding(2, 0, 'BARRACKS', 1)

  // AI单位
  addUnit(MAP_W - 2, MAP_H - 2, 'SOLDIER', 2)
  addUnit(MAP_W - 3, MAP_H - 2, 'ARCHER', 2)
  addBuilding(MAP_W - 2, MAP_H - 1, 'BASE', 2)
  addBuilding(MAP_W - 3, MAP_H - 1, 'BARRACKS', 2)

  render()
}

const addUnit = (x, y, type, p) => {
  const c = UNIT_TYPES[type]
  const u = { x, y, type, p, hp: c.hp, maxHp: c.hp, moved: false, attacked: false }
  units.value.push(u)
  map.value[y][x].unit = u
}

const addBuilding = (x, y, type, p) => {
  const b = { x, y, type, p, hp: 400, maxHp: 400 }
  buildings.value.push(b)
  map.value[y][x].building = b
}

const handleMapClick = (e) => {
  if (busy.value || currentPlayer.value !== 1) return

  const rect = canvasRef.value.getBoundingClientRect()
  const x = Math.floor((e.clientX - rect.left) / TILE)
  const y = Math.floor((e.clientY - rect.top) / TILE)
  if (x < 0 || x >= MAP_W || y < 0 || y >= MAP_H) return

  const tile = map.value[y][x]

  // 移动
  if (selected.value && moveTiles.value.some(t => t.x === x && t.y === y)) {
    moveUnit(selected.value, x, y)
    return
  }

  // 攻击
  if (selected.value && attackTiles.value.some(t => t.x === x && t.y === y)) {
    attack(selected.value, tile.unit)
    return
  }

  // 选择
  if (tile.unit && tile.unit.p === 1 && !tile.unit.moved) {
    selectUnit(tile.unit)
  } else if (tile.building && tile.building.p === 1) {
    selected.value = tile.building
    moveTiles.value = []
    attackTiles.value = []
    render()
  } else {
    clearSelection()
  }
}

const selectUnit = (u) => {
  selected.value = u
  calcRange(u)
  render()
}

const clearSelection = () => {
  selected.value = null
  moveTiles.value = []
  attackTiles.value = []
  render()
}

const calcRange = (u) => {
  moveTiles.value = []
  attackTiles.value = []
  const c = UNIT_TYPES[u.type]

  for (let y = 0; y < MAP_H; y++) {
    for (let x = 0; x < MAP_W; x++) {
      const d = Math.abs(x - u.x) + Math.abs(y - u.y)
      if (d > 0 && d <= c.move && !map.value[y][x].unit && map.value[y][x].terrain !== 'WATER') {
        moveTiles.value.push({ x, y })
      }
      if (d > 0 && d <= c.range) {
        const t = map.value[y][x].unit
        if (t && t.p !== u.p) attackTiles.value.push({ x, y })
      }
    }
  }
}

const moveUnit = (u, nx, ny) => {
  busy.value = true
  map.value[u.y][u.x].unit = null
  u.x = nx
  u.y = ny
  u.moved = true
  map.value[ny][nx].unit = u
  busy.value = false
  calcRange(u)
  render()
}

const attack = (attacker, target) => {
  if (attacker.attacked) return

  const c = UNIT_TYPES[attacker.type]
  const dmg = Math.max(1, c.atk - UNIT_TYPES[target.type].def)
  target.hp -= dmg
  attacker.attacked = true
  attacker.moved = true

  ElMessage.warning(`${UNIT_TYPES[target.type].emoji} 受到 ${dmg} 伤害！`)

  if (target.hp <= 0) {
    map.value[target.y][target.x].unit = null
    units.value = units.value.filter(u => u !== target)
    killedUnits.value++
    ElMessage.success(`${UNIT_TYPES[target.type].emoji} 被消灭！`)
    checkGameEnd()
  }
  clearSelection()
}

const produceUnit = (type) => {
  const c = UNIT_TYPES[type]
  if (gold[1] < c.cost) {
    ElMessage.error('金币不足！')
    return
  }

  const b = buildings.value.find(b => b.p === 1)
  if (!b) return

  const dirs = [[0, 1], [1, 0], [0, -1], [-1, 0]]
  let sx = -1, sy = -1

  for (const [dx, dy] of dirs) {
    const nx = b.x + dx, ny = b.y + dy
    if (nx >= 0 && nx < MAP_W && ny >= 0 && ny < MAP_H && !map.value[ny][nx].unit && map.value[ny][nx].terrain !== 'WATER') {
      sx = nx; sy = ny; break
    }
  }

  if (sx < 0) {
    ElMessage.warning('无空位！')
    return
  }

  gold[1] -= c.cost
  addUnit(sx, sy, type, 1)
  showBuildMenu.value = false
  ElMessage.success(`生产了 ${c.emoji} ${c.name}！`)
  render()
}

const endTurn = () => {
  units.value.forEach(u => {
    if (u.p === currentPlayer.value) { u.moved = false; u.attacked = false }
  })
  currentPlayer.value = currentPlayer.value === 1 ? 2 : 1
  if (currentPlayer.value === 1) turn.value++
  gold[currentPlayer.value] += 30

  clearSelection()

  if (currentPlayer.value === 2) {
    setTimeout(aiTurn, 800)
  }
}

const aiTurn = () => {
  const aiUnits = units.value.filter(u => u.p === 2 && !u.moved)

  if (aiUnits.length === 0) {
    // AI生产
    const types = Object.keys(UNIT_TYPES).filter(t => UNIT_TYPES[t].cost <= gold[2])
    if (types.length > 0) {
      const t = types[Math.floor(Math.random() * types.length)]
      const b = buildings.value.find(b => b.p === 2)
      if (b) {
        const dirs = [[0, 1], [1, 0], [0, -1], [-1, 0]]
        for (const [dx, dy] of dirs) {
          const nx = b.x + dx, ny = b.y + dy
          if (nx >= 0 && nx < MAP_W && ny >= 0 && ny < MAP_H && !map.value[ny][nx].unit && map.value[ny][nx].terrain !== 'WATER') {
            gold[2] -= UNIT_TYPES[t].cost
            addUnit(nx, ny, t, 2)
            break
          }
        }
      }
    }
    endTurn()
    return
  }

  const u = aiUnits[0]
  const enemy = units.value.find(e => e.p === 1)
  if (!enemy) { u.moved = true; setTimeout(aiTurn, 300); return }

  const d = Math.abs(u.x - enemy.x) + Math.abs(u.y - enemy.y)
  const c = UNIT_TYPES[u.type]

  if (d <= c.range) {
    selectUnit(u)
    setTimeout(() => {
      attack(u, enemy)
      setTimeout(aiTurn, 400)
    }, 300)
  } else {
    let dx = enemy.x - u.x, dy = enemy.y - u.y
    let nx = u.x + (dx !== 0 ? Math.sign(dx) : 0)
    let ny = u.y + (dy !== 0 && dx === 0 ? Math.sign(dy) : 0)

    if (nx >= 0 && nx < MAP_W && ny >= 0 && ny < MAP_H && !map.value[ny][nx].unit && map.value[ny][nx].terrain !== 'WATER') {
      map.value[u.y][u.x].unit = null
      u.x = nx; u.y = ny; u.moved = true
      map.value[ny][nx].unit = u
      render()
    } else {
      u.moved = true
    }
    setTimeout(aiTurn, 300)
  }
}

const checkGameEnd = () => {
  const playerBase = buildings.value.find(b => b.type === 'BASE' && b.p === 1)
  const aiBase = buildings.value.find(b => b.type === 'BASE' && b.p === 2)

  if (!aiBase || aiBase.hp <= 0) {
    handleEndGame(true)
  } else if (!playerBase || playerBase.hp <= 0) {
    handleEndGame(false)
  }
}

const handleEndGame = async (isWin) => {
  const winCount = isWin ? 1 : 0

  const resultText = isWin ? '胜利' : '战败'

  const detail = {
    turns: turn.value,
    killedUnits: killedUnits.value,
    remainingGold: gold[1]
  }

  const gameResult = await endMiniGame(winCount, resultText, detail)

  resultIcon.value = isWin ? '🏆' : '😢'
  resultTitle.value = isWin ? '胜利！' : '战败！'
  resultReward.value = gameResult
    ? `获得 ${gameResult.cashChange} 元，名声 +${gameResult.fameChange}`
    : ''

  emit('complete', {
    game: '小小大战争',
    result: resultText,
    ...gameResult
  })
}

// 渲染
const render = () => {
  if (!ctx) return

  ctx.fillStyle = '#3a5a3a'
  ctx.fillRect(0, 0, MAP_W * TILE, MAP_H * TILE)

  // 地形
  for (let y = 0; y < MAP_H; y++) {
    for (let x = 0; x < MAP_W; x++) {
      ctx.fillStyle = TERRAIN[map.value[y][x].terrain]
      ctx.fillRect(x * TILE, y * TILE, TILE, TILE)
      ctx.strokeStyle = 'rgba(0,0,0,0.2)'
      ctx.strokeRect(x * TILE, y * TILE, TILE, TILE)
    }
  }

  // 移动范围
  ctx.fillStyle = 'rgba(0,200,0,0.35)'
  moveTiles.value.forEach(t => ctx.fillRect(t.x * TILE, t.y * TILE, TILE, TILE))

  // 攻击范围
  ctx.fillStyle = 'rgba(200,0,0,0.4)'
  attackTiles.value.forEach(t => ctx.fillRect(t.x * TILE, t.y * TILE, TILE, TILE))

  // 建筑
  buildings.value.forEach(b => {
    const x = b.x * TILE, y = b.y * TILE
    ctx.fillStyle = COLORS[b.p]
    ctx.globalAlpha = 0.8
    ctx.fillRect(x + 3, y + 3, TILE - 6, TILE - 6)
    ctx.globalAlpha = 1
    ctx.font = '22px Arial'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(b.type === 'BASE' ? '🏰' : '🏭', x + TILE / 2, y + TILE / 2)

    // 血条
    ctx.fillStyle = '#333'
    ctx.fillRect(x + 4, y + 2, TILE - 8, 3)
    ctx.fillStyle = '#0f0'
    ctx.fillRect(x + 4, y + 2, (TILE - 8) * (b.hp / b.maxHp), 3)

    if (selected.value === b) {
      ctx.strokeStyle = '#ff0'
      ctx.lineWidth = 2
      ctx.strokeRect(x + 1, y + 1, TILE - 2, TILE - 2)
      ctx.lineWidth = 1
    }
  })

  // 单位
  units.value.forEach(u => {
    const px = u.x * TILE + TILE / 2
    const py = u.y * TILE + TILE / 2
    const c = UNIT_TYPES[u.type]

    if (selected.value === u) {
      ctx.strokeStyle = '#ff0'
      ctx.lineWidth = 2
      ctx.beginPath()
      ctx.arc(px, py, TILE / 2 - 4, 0, Math.PI * 2)
      ctx.stroke()
      ctx.lineWidth = 1
    }

    ctx.fillStyle = COLORS[u.p]
    ctx.globalAlpha = u.moved ? 0.5 : 1
    ctx.beginPath()
    ctx.arc(px, py, TILE / 2 - 10, 0, Math.PI * 2)
    ctx.fill()
    ctx.globalAlpha = 1

    ctx.font = '20px Arial'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(c.emoji, px, py)

    ctx.fillStyle = '#333'
    ctx.fillRect(px - 12, py - 20, 24, 3)
    ctx.fillStyle = u.hp > u.maxHp / 2 ? '#0f0' : u.hp > u.maxHp / 4 ? '#ff0' : '#f00'
    ctx.fillRect(px - 12, py - 20, 24 * (u.hp / u.maxHp), 3)
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
  margin-bottom: 15px;
}

.game-info {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 15px;
}

.info-item {
  font-size: 13px;
  padding: 6px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.intro-rules {
  background: #f5f7fa;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 15px;
  text-align: left;
}

.rule-item {
  font-size: 12px;
  color: var(--font-secondary);
  padding: 4px 0;
}

.game-playing {
  text-align: center;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 10px;
}

.player-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.player-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
}

.player-gold {
  font-size: 11px;
  color: var(--warning-color);
}

.game-info-inline {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.turn-indicator {
  font-size: 11px;
  font-weight: 600;
  color: #409eff;
}

.turn-indicator.enemy {
  color: var(--error-color);
}

.move-count {
  font-size: 10px;
  color: var(--font-light);
}

.map-container {
  display: flex;
  justify-content: center;
  margin-bottom: 10px;
}

.map-container canvas {
  border: 2px solid #5c4033;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.game-status {
  margin-bottom: 10px;
  min-height: 24px;
}

.selected-info {
  font-size: 12px;
  color: var(--font-color);
  font-weight: 500;
}

.status-text {
  font-size: 12px;
  color: var(--font-light);
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-buttons :deep(.el-button) {
  padding: 6px 12px;
  font-size: 12px;
}

.build-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.build-option {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.build-option:hover {
  background: #e9ecf0;
}

.build-option.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.unit-emoji {
  font-size: 18px;
  margin-right: 8px;
}

.unit-name {
  flex: 1;
  font-size: 13px;
  color: var(--font-color);
}

.unit-cost {
  font-size: 12px;
  color: var(--warning-color);
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
  margin-bottom: 12px;
}

.result-stats {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 12px;
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
