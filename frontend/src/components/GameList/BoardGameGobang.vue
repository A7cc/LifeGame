<template>
  <el-dialog v-model="visible" title="五子棋" width="450px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">在15x15的棋盘上连成五子获胜</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 300 }}元</span>
          <span class="info-item">🏆 奖励：600元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <span class="turn-indicator">{{ currentPlayer === 1 ? '你的回合 (⚫)' : 'AI思考中 (⚪)' }}</span>
          <span class="move-count">第 {{ moveNumber }} 手</span>
        </div>

        <div class="board-container">
          <div class="gomoku-board">
            <div
              v-for="(row, rowIndex) in board"
              :key="rowIndex"
              class="board-row"
            >
              <div
                v-for="(cell, colIndex) in row"
                :key="colIndex"
                class="board-cell"
                :class="{
                  'selected': selectedCell && selectedCell.row === rowIndex && selectedCell.col === colIndex,
                  'last-move': lastMove && lastMove.row === rowIndex && lastMove.col === colIndex,
                  'win-cell': isWinCell(rowIndex, colIndex)
                }"
                @click="handleCellClick(rowIndex, colIndex)"
              >
                <div v-if="cell !== 0" class="stone" :class="cell === 1 ? 'black' : 'white'"></div>
              </div>
            </div>
          </div>
        </div>

        <div class="action-buttons">
          <el-button @click="undoMove" :disabled="!canUndo">↩️ 悔棋</el-button>
          <el-button @click="resign">🏳️ 认输</el-button>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
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
    default: () => ({ id: 'gobang', name: '五子棋', entryCost: 300 })
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

const BOARD_SIZE = 15

const currentPlayer = ref(1) // 1: player, 2: AI
const moveNumber = ref(0)
const board = ref([])
const selectedCell = ref(null)
const lastMove = ref(null)
const winCells = ref([])
const canUndo = ref(false)
const moveHistory = ref([])

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const resetLocal = () => {
  currentPlayer.value = 1
  moveNumber.value = 0
  board.value = []
  selectedCell.value = null
  lastMove.value = null
  winCells.value = []
  canUndo.value = false
  moveHistory.value = []
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
    ElMessage.info('游戏开始！你执黑棋先行')
  }
}

const handleCellClick = (row, col) => {
  if (currentPlayer.value !== 1 || gameEnded.value) return
  if (board.value[row][col] !== 0) return

  makeMove(row, col, 1)
}

const makeMove = (row, col, player) => {
  // 保存历史
  moveHistory.value.push({
    board: JSON.parse(JSON.stringify(board.value)),
    currentPlayer: currentPlayer.value,
    moveNumber: moveNumber.value
  })
  canUndo.value = true

  // 落子
  board.value[row][col] = player
  lastMove.value = { row, col }
  moveNumber.value++

  // 检查胜利
  const winResult = checkWin(row, col, player)
  if (winResult) {
    winCells.value = winResult
    finishGame(player)
    return
  }

  // 检查平局
  if (moveNumber.value >= BOARD_SIZE * BOARD_SIZE) {
    finishGame(0)
    return
  }

  // 切换玩家
  currentPlayer.value = player === 1 ? 2 : 1
  if (currentPlayer.value === 2) {
    setTimeout(() => aiMove(), 500)
  }
}

const checkWin = (row, col, player) => {
  const directions = [
    [[0, 1], [0, -1]],   // 水平
    [[1, 0], [-1, 0]],   // 垂直
    [[1, 1], [-1, -1]], // 对角线
    [[1, -1], [-1, 1]]  // 反对角线
  ]

  for (const [dir1, dir2] of directions) {
    const cells = [[row, col]]

    // 向两个方向延伸
    for (const [dr, dc] of [dir1, dir2]) {
      let r = row + dr
      let c = col + dc
      while (
        r >= 0 && r < BOARD_SIZE &&
        c >= 0 && c < BOARD_SIZE &&
        board.value[r][c] === player
      ) {
        cells.push([r, c])
        r += dr
        c += dc
      }
    }

    if (cells.length >= 5) {
      return cells
    }
  }

  return null
}

const isWinCell = (row, col) => {
  return winCells.value.some(([r, c]) => r === row && c === col)
}

const aiMove = () => {
  // 简单AI策略：防守优先，进攻其次
  let bestScore = -Infinity
  let bestMoves = []

  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 0) {
        const score = evaluatePosition(r, c)
        if (score > bestScore) {
          bestScore = score
          bestMoves = [[r, c]]
        } else if (score === bestScore) {
          bestMoves.push([r, c])
        }
      }
    }
  }

  if (bestMoves.length > 0) {
    const [row, col] = bestMoves[Math.floor(Math.random() * bestMoves.length)]
    makeMove(row, col, 2)
  }
}

const evaluatePosition = (row, col) => {
  // 评估位置价值，防守和进攻都考虑
  let score = 0

  // 检查所有方向的连子情况
  const directions = [[0, 1], [1, 0], [1, 1], [1, -1]]

  for (const [dr, dc] of directions) {
    // 评估AI在此位置落子的价值（进攻）
    const aiScore = evaluateDirection(row, col, dr, dc, 2)
    // 评估玩家在此位置落子的价值（防守）
    const playerScore = evaluateDirection(row, col, dr, dc, 1)

    // 防守权重更高，因为阻止玩家胜利更重要
    score += aiScore + playerScore * 1.2
  }

  // 中心位置加分
  const centerDist = Math.abs(row - 7) + Math.abs(col - 7)
  score += (14 - centerDist) * 2

  return score
}

const evaluateDirection = (row, col, dr, dc, player) => {
  let count = 1
  let openEnds = 0

  // 正方向
  let r = row + dr
  let c = col + dc
  while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE && board.value[r][c] === player) {
    count++
    r += dr
    c += dc
  }
  if (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE && board.value[r][c] === 0) {
    openEnds++
  }

  // 反方向
  r = row - dr
  c = col - dc
  while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE && board.value[r][c] === player) {
    count++
    r -= dr
    c -= dc
  }
  if (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE && board.value[r][c] === 0) {
    openEnds++
  }

  // 评分
  if (count >= 5) return 100000
  if (count === 4 && openEnds === 2) return 10000
  if (count === 4 && openEnds === 1) return 1000
  if (count === 3 && openEnds === 2) return 500
  if (count === 3 && openEnds === 1) return 100
  if (count === 2 && openEnds === 2) return 50
  if (count === 2 && openEnds === 1) return 10

  return count
}

const undoMove = () => {
  if (!canUndo.value || moveHistory.value.length === 0) return

  // 悔两步（玩家和AI各一步）
  if (moveHistory.value.length >= 2) {
    moveHistory.value.pop() // AI的步
    const state = moveHistory.value.pop() // 玩家的步
    board.value = state.board
    currentPlayer.value = state.currentPlayer
    moveNumber.value = state.moveNumber
    lastMove.value = null
  } else if (moveHistory.value.length === 1) {
    const state = moveHistory.value.pop()
    board.value = state.board
    currentPlayer.value = state.currentPlayer
    moveNumber.value = state.moveNumber
    lastMove.value = null
  }

  canUndo.value = moveHistory.value.length > 0
  ElMessage.info('已悔棋')
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

const finishGame = async (winner) => {
  const winCount = winner === 1 ? 1 : 0
  const customResultText = winner === 1 ? '获胜' : winner === 2 ? '败北' : '平局'

  const gameResult = await endGame(winCount, customResultText, {
    moves: moveNumber.value,
    winner
  })

  if (gameResult) {
    if (winner === 1) {
      resultIcon.value = '🏆'
      resultTitle.value = '恭喜获胜！'
    } else if (winner === 2) {
      resultIcon.value = '😢'
      resultTitle.value = 'AI获胜！'
    } else {
      resultIcon.value = '🤝'
      resultTitle.value = '平局！'
    }
    resultReward.value = gameResult.resultText

    emit('complete', gameResult)
  }
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
  color: var(--font-color);
}

.game-playing {
  text-align: center;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 10px;
}

.turn-indicator {
  font-size: 12px;
  font-weight: 600;
  color: var(--el-color-primary);
}

.move-count {
  font-size: 11px;
  color: var(--font-secondary);
}

.board-container {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.gomoku-board {
  display: inline-grid;
  grid-template-columns: repeat(15, 24px);
  grid-template-rows: repeat(15, 24px);
  gap: 1px;
  background: #daa520;
  padding: 10px;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.board-row {
  display: contents;
}

.board-cell {
  width: 24px;
  height: 24px;
  background: #daa520;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}

.board-cell::before,
.board-cell::after {
  content: '';
  position: absolute;
  background: #8b6914;
}

.board-cell::before {
  width: 100%;
  height: 1px;
  top: 50%;
  left: 0;
}

.board-cell::after {
  width: 1px;
  height: 100%;
  left: 50%;
  top: 0;
}

.board-cell:hover {
  background: rgba(255, 255, 255, 0.2);
}

.board-cell.selected {
  background: rgba(255, 255, 0, 0.3);
}

.board-cell.last-move::before {
  background: var(--error-color);
  height: 2px;
}

.board-cell.win-cell {
  background: rgba(var(--success-color-rgb), 0.3);
}

.stone {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  position: relative;
  z-index: 1;
  box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.stone.black {
  background: radial-gradient(circle at 30% 30%, #444, #000);
}

.stone.white {
  background: radial-gradient(circle at 30% 30%, #fff, #ccc);
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
  color: var(--font-color);
  margin-bottom: 15px;
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
