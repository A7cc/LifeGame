<template>
  <el-dialog v-model="visible" title="围棋" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">
          在9x9棋盘上进行围棋对弈<br>
          <span class="rules">围地吃子，以目数决胜</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 500 }}元</span>
          <span class="info-item">🏆 奖励：1000元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始对弈</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <div class="score-board">
            <div class="score-item">
              <span class="player-label">⚫ 你</span>
              <span class="score-value">{{ playerCaptures }}</span>
            </div>
            <div class="turn-info">
              <span class="turn-indicator">{{ currentPlayer === 1 ? '你的回合' : 'AI思考中...' }}</span>
              <span class="move-count">第 {{ moveNumber }} 手</span>
            </div>
            <div class="score-item">
              <span class="player-label">⚪ AI</span>
              <span class="score-value">{{ aiCaptures }}</span>
            </div>
          </div>
        </div>

        <div class="board-container">
          <div class="go-board">
            <svg :width="boardSize" :height="boardSize" class="board-svg">
              <!-- 绘制网格 -->
              <g stroke="#8b6914" stroke-width="1">
                <line v-for="i in 9" :key="'h-' + i"
                  :x1="padding" :y1="padding + (i - 1) * cellSize"
                  :x2="boardSize - padding" :y2="padding + (i - 1) * cellSize" />
                <line v-for="i in 9" :key="'v-' + i"
                  :x1="padding + (i - 1) * cellSize" :y1="padding"
                  :x2="padding + (i - 1) * cellSize" :y2="boardSize - padding" />
              </g>

              <!-- 星位 -->
              <g fill="#8b6914">
                <circle v-for="star in starPoints" :key="star.x + '-' + star.y"
                  :cx="padding + star.x * cellSize"
                  :cy="padding + star.y * cellSize"
                  r="4" />
              </g>

              <!-- 棋子和点击区域 -->
              <g>
                <template v-for="(row, rowIndex) in board" :key="'rect-row-' + rowIndex">
                  <rect v-for="(_, colIndex) in row"
                    :key="'rect-' + rowIndex + '-' + colIndex"
                    :x="padding + colIndex * cellSize - cellSize / 2"
                    :y="padding + rowIndex * cellSize - cellSize / 2"
                    :width="cellSize"
                    :height="cellSize"
                    fill="transparent"
                    :class="{
                      'last-move': lastMove && lastMove.row === rowIndex && lastMove.col === colIndex
                    }"
                    @click="handleCellClick(rowIndex, colIndex)"
                    style="cursor: pointer" />
                </template>

                <!-- 棋子 -->
                <template v-for="(row, rowIndex) in board" :key="'stone-row-' + rowIndex">
                  <circle v-for="(cell, colIndex) in row"
                    :key="'stone-' + rowIndex + '-' + colIndex"
                    v-if="cell !== 0"
                    :cx="padding + colIndex * cellSize"
                    :cy="padding + rowIndex * cellSize"
                    :r="stoneRadius"
                    :class="cell === 1 ? 'stone-black' : 'stone-white'"
                    stroke="rgba(0,0,0,0.2)"
                    stroke-width="1" />
                </template>
              </g>
            </svg>
          </div>
        </div>

        <div class="action-buttons">
          <el-button @click="passMove">⏭️ 跳过(Pass)</el-button>
          <el-button @click="resign">🏳️ 认输</el-button>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">你的提子：{{ playerCaptures }}</div>
          <div class="stat-item">AI提子：{{ aiCaptures }}</div>
          <div class="stat-item">总手数：{{ moveNumber }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'go', name: '围棋', entryCost: 500 })
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

// 使用小游戏基础逻辑
const { gameStarted, gameEnded, processing, startGame, endGame, reset } = useMiniGameBase(props.config)

const BOARD_SIZE = 9

const boardSize = ref(360)
const padding = ref(20)
const cellSize = ref(40)
const stoneRadius = ref(17)

const starPoints = [
  { x: 2, y: 2 }, { x: 6, y: 2 },
  { x: 4, y: 4 },
  { x: 2, y: 6 }, { x: 6, y: 6 }
]

const currentPlayer = ref(1) // 1: black (player), 2: white (AI)
const moveNumber = ref(0)
const consecutivePasses = ref(0)
const board = ref([])
const lastMove = ref(null)
const playerCaptures = ref(0)
const aiCaptures = ref(0)

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const resetLocal = () => {
  currentPlayer.value = 1
  moveNumber.value = 0
  consecutivePasses.value = 0
  board.value = []
  lastMove.value = null
  playerCaptures.value = 0
  aiCaptures.value = 0
  resultIcon.value = ''
  resultTitle.value = ''
  resultReward.value = ''
  reset()
}

const handleClose = () => {
  visible.value = false
  resetLocal()
}

// 开始游戏
const handleStartGame = async () => {
  const success = await startGame()
  if (success) {
    // 初始化空棋盘
    board.value = Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0))
    ElMessage.info('对弈开始！你执黑棋先行')
  }
}

const handleCellClick = (row, col) => {
  if (currentPlayer.value !== 1 || gameEnded.value) return
  if (board.value[row][col] !== 0) return

  // 检查是否是合法落子（不自杀）
  if (!isLegalMove(row, col, 1)) {
    ElMessage.warning('此处不能落子（会被立即提吃）')
    return
  }

  makeMove(row, col, 1)
}

const isLegalMove = (row, col, player) => {
  // 临时落子
  const tempBoard = board.value.map(r => [...r])
  tempBoard[row][col] = player

  // 检查是否能提对方的子
  const opponent = player === 1 ? 2 : 1
  const captured = getCapturedStones(tempBoard, row, col, opponent)

  // 如果能提子，则是合法的
  if (captured.length > 0) return true

  // 检查自己是否会被提（自杀）
  const selfCaptured = getCapturedStones(tempBoard, row, col, player)
  return selfCaptured.length === 0
}

const makeMove = (row, col, player) => {
  const opponent = player === 1 ? 2 : 1

  // 先检查提子
  const captured = getCapturedStones(board.value, row, col, opponent)
  captured.forEach(([r, c]) => {
    board.value[r][c] = 0
    if (player === 1) {
      playerCaptures.value++
    } else {
      aiCaptures.value++
    }
  })

  // 落子
  board.value[row][col] = player
  lastMove.value = { row, col }
  moveNumber.value++
  consecutivePasses.value = 0

  currentPlayer.value = opponent

  if (currentPlayer.value === 2) {
    setTimeout(() => aiMove(), 800)
  }
}

const getCapturedStones = (boardState, row, col, color) => {
  const captured = []
  const checked = new Set()

  const directions = [[0, 1], [0, -1], [1, 0], [-1, 0]]

  // 检查四个方向的相邻同色棋子
  for (const [dr, dc] of directions) {
    const nr = row + dr
    const nc = col + dc

    if (nr >= 0 && nr < BOARD_SIZE && nc >= 0 && nc < BOARD_SIZE) {
      if (boardState[nr][nc] === color && !checked.has(`${nr},${nc}`)) {
        const group = getGroup(boardState, nr, nc, color)

        // 标记已检查
        group.forEach(([r, c]) => checked.add(`${r},${c}`))

        // 检查该组是否有气
        if (!hasLiberties(boardState, group)) {
          captured.push(...group)
        }
      }
    }
  }

  return captured
}

const getGroup = (boardState, row, col, color) => {
  const group = []
  const visited = new Set()
  const stack = [[row, col]]

  while (stack.length > 0) {
    const [r, c] = stack.pop()
    const key = `${r},${c}`

    if (visited.has(key)) continue
    if (r < 0 || r >= BOARD_SIZE || c < 0 || c >= BOARD_SIZE) continue
    if (boardState[r][c] !== color) continue

    visited.add(key)
    group.push([r, c])

    stack.push([r + 1, c], [r - 1, c], [r, c + 1], [r, c - 1])
  }

  return group
}

const hasLiberties = (boardState, group) => {
  const directions = [[0, 1], [0, -1], [1, 0], [-1, 0]]

  for (const [row, col] of group) {
    for (const [dr, dc] of directions) {
      const nr = row + dr
      const nc = col + dc

      if (nr >= 0 && nr < BOARD_SIZE && nc >= 0 && nc < BOARD_SIZE) {
        if (boardState[nr][nc] === 0) {
          return true
        }
      }
    }
  }

  return false
}

const passMove = () => {
  consecutivePasses.value++

  if (consecutivePasses.value >= 2) {
    endGameByScore()
    return
  }

  ElMessage.info('你选择跳过')
  currentPlayer.value = 2

  setTimeout(() => aiMove(), 600)
}

const aiMove = () => {
  // 简单AI策略
  let bestMove = null
  let bestScore = -Infinity

  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 0 && isLegalMove(r, c, 2)) {
        const score = evaluatePosition(r, c)
        if (score > bestScore) {
          bestScore = score
          bestMove = { row: r, col: c }
        }
      }
    }
  }

  if (bestMove) {
    makeMove(bestMove.row, bestMove.col, 2)
  } else {
    // AI也跳过
    consecutivePasses.value++
    if (consecutivePasses.value >= 2) {
      endGameByScore()
    } else {
      ElMessage.info('AI选择跳过')
      currentPlayer.value = 1
    }
  }
}

const evaluatePosition = (row, col) => {
  let score = Math.random() * 20

  // 中心加分
  const centerDist = Math.abs(row - 4) + Math.abs(col - 4)
  score += (8 - centerDist) * 3

  // 星位加分
  if (starPoints.some(s => s.x === col && s.y === row)) {
    score += 15
  }

  // 靠近自己的子加分（连成一片）
  let adjacentSame = 0
  const directions = [[0, 1], [0, -1], [1, 0], [-1, 0]]
  for (const [dr, dc] of directions) {
    const nr = row + dr
    const nc = col + dc
    if (nr >= 0 && nr < BOARD_SIZE && nc >= 0 && nc < BOARD_SIZE) {
      if (board.value[nr][nc] === 2) {
        adjacentSame++
      } else if (board.value[nr][nc] === 1) {
        // 靠近对方子也加分（进攻）
        score += 8
      }
    }
  }
  score += adjacentSame * 10

  // 边角减分（不易做活）
  if (row === 0 || row === 8 || col === 0 || col === 8) {
    score -= 5
  }

  // 检查能否提子
  const tempBoard = board.value.map(r => [...r])
  tempBoard[row][col] = 2
  const captured = getCapturedStones(tempBoard, row, col, 1)
  score += captured.length * 50

  return score
}

const countTerritory = () => {
  // 简化版：只计算提子和盘面子数
  let playerStones = 0
  let aiStones = 0

  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 1) playerStones++
      else if (board.value[r][c] === 2) aiStones++
    }
  }

  return {
    player: playerStones + playerCaptures.value,
    ai: aiStones + aiCaptures.value
  }
}

const endGameByScore = async () => {
  const scores = countTerritory()

  // 贴目（AI执白，贴2.75目）
  const adjustedAi = scores.ai + 2.75
  const winner = scores.player > adjustedAi ? 1 : 2

  const winCount = winner === 1 ? 1 : 0
  const customResultText = winner === 1
    ? `获胜，目数 ${scores.player}`
    : '败北'

  const gameResult = await endGame(winCount, customResultText, {
    playerScore: scores.player,
    aiScore: Math.round(adjustedAi),
    moves: moveNumber.value
  })

  if (gameResult) {
    if (winner === 1) {
      resultIcon.value = '🏆'
      resultTitle.value = `恭喜获胜！(${scores.player} vs ${Math.round(adjustedAi)})`
    } else {
      resultIcon.value = '😢'
      resultTitle.value = `AI获胜！(${scores.player} vs ${Math.round(adjustedAi)})`
    }
    resultReward.value = gameResult.resultText

    emit('complete', gameResult)
  }
}

const resign = () => {
  ElMessageBox.confirm('确定要认输吗？', '认输', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const gameResult = await endGame(0, '认输', { moves: moveNumber.value })
    if (gameResult) {
      resultIcon.value = '😢'
      resultTitle.value = '认输'
      resultReward.value = gameResult.resultText
      emit('complete', gameResult)
    }
  }).catch(() => {})
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
  line-height: 1.6;
}

.rules {
  font-size: 12px;
  color: var(--font-light);
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

.game-header {
  margin-bottom: 10px;
}

.score-board {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.score-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.player-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
}

.score-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-color-primary);
}

.turn-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.turn-indicator {
  font-size: 11px;
  font-weight: 600;
  color: var(--el-color-primary);
}

.move-count {
  font-size: 10px;
  color: var(--font-light);
}

.board-container {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.go-board {
  background: #daa520;
  padding: padding;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.board-svg {
  display: block;
}

.stone-black {
  fill: #1a1a1a;
}

.stone-white {
  fill: #f5f5f5;
}

.last-move {
  fill: rgba(255, 0, 0, 0.1);
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
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 15px;
}

.result-stats {
  display: flex;
  justify-content: center;
  gap: 20px;
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
