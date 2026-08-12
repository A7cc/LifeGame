<template>
  <el-dialog v-model="visible" title="黑白棋" width="440px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">
          策略棋类游戏，通过翻转对方棋子获胜<br>
          <span class="rules">夹住对方棋子即可翻转</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 300 }}元</span>
          <span class="info-item">🏆 奖励：600元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <div class="score-board">
            <div class="score-item" :class="{ active: currentPlayer === 1 }">
              <span class="disc-icon"></span>
              <span class="player-label">你</span>
              <span class="score-value">{{ blackCount }}</span>
            </div>
            <div class="turn-info">
              <span class="turn-indicator">{{ currentPlayer === 1 ? '你的回合' : 'AI思考中...' }}</span>
            </div>
            <div class="score-item" :class="{ active: currentPlayer === 2 }">
              <span class="disc-icon white"></span>
              <span class="player-label">AI</span>
              <span class="score-value">{{ whiteCount }}</span>
            </div>
          </div>
        </div>

        <div class="board-container">
          <div class="othello-board">
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
                  'valid-move': currentPlayer === 1 && isValidMove(rowIndex, colIndex),
                  'last-move': lastMove && lastMove.row === rowIndex && lastMove.col === colIndex
                }"
                @click="handleCellClick(rowIndex, colIndex)"
              >
                <div v-if="cell !== 0" class="disc" :class="cell === 1 ? 'black' : 'white'"></div>
                <div v-else-if="currentPlayer === 1 && isValidMove(rowIndex, colIndex)" class="hint"></div>
              </div>
            </div>
          </div>
        </div>

        <div class="game-status">
          <div v-if="!hasValidMove(1)" class="no-move-msg">你无处可下，跳过回合</div>
        </div>

        <div class="action-buttons">
          <el-button @click="resign">🏳️ 认输</el-button>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">⚫ 你的棋子：{{ finalBlack }}</div>
          <div class="stat-item">⚪ AI棋子：{{ finalWhite }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'othello', name: '黑白棋', entryCost: 300 })
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

const BOARD_SIZE = 8

const currentPlayer = ref(1) // 1: black (player), 2: white (AI)
const board = ref([])
const lastMove = ref(null)
const consecutivePasses = ref(0)

const blackCount = computed(() => countPieces(1))
const whiteCount = computed(() => countPieces(2))

const finalBlack = ref(0)
const finalWhite = ref(0)

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const resetLocal = () => {
  currentPlayer.value = 1
  board.value = []
  lastMove.value = null
  consecutivePasses.value = 0
  finalBlack.value = 0
  finalWhite.value = 0
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
    // 初始化棋盘
    board.value = Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(0))

    // 放置初始棋子
    const mid = BOARD_SIZE / 2
    board.value[mid - 1][mid - 1] = 2 // 白
    board.value[mid - 1][mid] = 1     // 黑
    board.value[mid][mid - 1] = 1     // 黑
    board.value[mid][mid] = 2         // 白

    ElMessage.info('游戏开始！你执黑棋先行')
  }
}

const countPieces = (player) => {
  let count = 0
  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === player) count++
    }
  }
  return count
}

const handleCellClick = (row, col) => {
  if (currentPlayer.value !== 1 || gameEnded.value) return
  if (board.value[row][col] !== 0) return
  if (!isValidMove(row, col)) return

  makeMove(row, col, 1)
}

const isValidMove = (row, col) => {
  if (board.value[row][col] !== 0) return false

  const directions = [
    [0, 1], [0, -1], [1, 0], [-1, 0],
    [1, 1], [1, -1], [-1, 1], [-1, -1]
  ]

  for (const [dr, dc] of directions) {
    if (wouldFlip(row, col, dr, dc)) {
      return true
    }
  }

  return false
}

const wouldFlip = (row, col, dr, dc) => {
  const player = currentPlayer.value
  const opponent = player === 1 ? 2 : 1
  let r = row + dr
  let c = col + dc
  let foundOpponent = false

  while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE) {
    if (board.value[r][c] === opponent) {
      foundOpponent = true
      r += dr
      c += dc
    } else if (board.value[r][c] === player) {
      return foundOpponent
    } else {
      return false
    }
  }

  return false
}

const makeMove = (row, col, player) => {
  board.value[row][col] = player
  lastMove.value = { row, col }

  const directions = [
    [0, 1], [0, -1], [1, 0], [-1, 0],
    [1, 1], [1, -1], [-1, 1], [-1, -1]
  ]

  for (const [dr, dc] of directions) {
    if (wouldFlipAt(row, col, dr, dc, player)) {
      flipDiscs(row, col, dr, dc, player)
    }
  }

  consecutivePasses.value = 0
  currentPlayer.value = player === 1 ? 2 : 1

  // 检查游戏是否结束
  checkGameEnd()

  // AI移动
  if (currentPlayer.value === 2 && !gameEnded.value) {
    setTimeout(() => aiMove(), 700)
  }
}

const wouldFlipAt = (row, col, dr, dc, player) => {
  const opponent = player === 1 ? 2 : 1
  let r = row + dr
  let c = col + dc
  let foundOpponent = false

  while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE) {
    if (board.value[r][c] === opponent) {
      foundOpponent = true
      r += dr
      c += dc
    } else if (board.value[r][c] === player) {
      return foundOpponent
    } else {
      return false
    }
  }

  return false
}

const flipDiscs = (row, col, dr, dc, player) => {
  const opponent = player === 1 ? 2 : 1
  let r = row + dr
  let c = col + dc

  while (board.value[r][c] === opponent) {
    board.value[r][c] = player
    r += dr
    c += dc
  }
}

const hasValidMove = (player) => {
  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 0) {
        const directions = [
          [0, 1], [0, -1], [1, 0], [-1, 0],
          [1, 1], [1, -1], [-1, 1], [-1, -1]
        ]
        for (const [dr, dc] of directions) {
          if (wouldFlipForPlayer(r, c, dr, dc, player)) {
            return true
          }
        }
      }
    }
  }
  return false
}

const wouldFlipForPlayer = (row, col, dr, dc, player) => {
  const opponent = player === 1 ? 2 : 1
  let r = row + dr
  let c = col + dc
  let foundOpponent = false

  while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE) {
    if (board.value[r][c] === opponent) {
      foundOpponent = true
      r += dr
      c += dc
    } else if (board.value[r][c] === player) {
      return foundOpponent
    } else {
      return false
    }
  }

  return false
}

const checkGameEnd = () => {
  // 检查是否有空位
  let hasEmpty = false
  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 0) {
        hasEmpty = true
        break
      }
    }
  }

  if (!hasEmpty) {
    finishGame()
    return
  }

  // 检查当前玩家是否有合法移动
  if (!hasValidMove(currentPlayer.value)) {
    consecutivePasses.value++

    if (consecutivePasses.value >= 2) {
      finishGame()
      return
    }

    const nextPlayer = currentPlayer.value === 1 ? 2 : 1
    if (!hasValidMove(nextPlayer)) {
      finishGame()
      return
    }

    ElMessage.info(`${currentPlayer.value === 1 ? '你' : 'AI'}无处可下，跳过回合`)
    currentPlayer.value = nextPlayer

    if (currentPlayer.value === 2 && !gameEnded.value) {
      setTimeout(() => aiMove(), 700)
    }
  }
}

const aiMove = () => {
  let bestMove = null
  let bestScore = -Infinity

  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      if (board.value[r][c] === 0 && isValidMoveForAI(r, c)) {
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
    consecutivePasses.value++
    checkGameEnd()
  }
}

const isValidMoveForAI = (row, col) => {
  const directions = [
    [0, 1], [0, -1], [1, 0], [-1, 0],
    [1, 1], [1, -1], [-1, 1], [-1, -1]
  ]

  for (const [dr, dc] of directions) {
    if (wouldFlipForPlayer(row, col, dr, dc, 2)) {
      return true
    }
  }
  return false
}

const evaluatePosition = (row, col) => {
  let score = Math.random() * 10

  // 角落位置优先级最高
  if ((row === 0 || row === 7) && (col === 0 || col === 7)) {
    score += 100
  }
  // 角落旁边的位置（星位）要避免
  else if (
    (row === 0 && (col === 1 || col === 6)) ||
    (row === 7 && (col === 1 || col === 6)) ||
    (col === 0 && (row === 1 || row === 6)) ||
    (col === 7 && (row === 1 || row === 6))
  ) {
    score -= 50
  }
  // 边缘位置比较好
  else if (row === 0 || row === 7 || col === 0 || col === 7) {
    score += 20
  }

  // 中心位置也不错
  if (row >= 2 && row <= 5 && col >= 2 && col <= 5) {
    score += 15
  }

  // 计算能翻转多少子
  const flipped = countFlipped(row, col, 2)
  score += flipped * 2

  return score
}

const countFlipped = (row, col, player) => {
  let count = 0
  const directions = [
    [0, 1], [0, -1], [1, 0], [-1, 0],
    [1, 1], [1, -1], [-1, 1], [-1, -1]
  ]

  for (const [dr, dc] of directions) {
    let r = row + dr
    let c = col + dc
    let potential = 0

    while (r >= 0 && r < BOARD_SIZE && c >= 0 && c < BOARD_SIZE) {
      const opponent = player === 1 ? 2 : 1
      if (board.value[r][c] === opponent) {
        potential++
        r += dr
        c += dc
      } else if (board.value[r][c] === player) {
        count += potential
        break
      } else {
        break
      }
    }
  }

  return count
}

const finishGame = async () => {
  finalBlack.value = blackCount.value
  finalWhite.value = whiteCount.value

  const winner = finalBlack.value > finalWhite.value ? 1 :
    finalWhite.value > finalBlack.value ? 2 : 0

  const winCount = winner === 1 ? 1 : 0
  const customResultText = winner === 1
    ? `获胜，${finalBlack.value}:${finalWhite.value}`
    : winner === 2
      ? `败北，${finalBlack.value}:${finalWhite.value}`
      : `平局，${finalBlack.value}:${finalWhite.value}`

  const gameResult = await endGame(winCount, customResultText, {
    black: finalBlack.value,
    white: finalWhite.value,
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

const resign = () => {
  ElMessageBox.confirm('确定要认输吗？', '认输', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    finalBlack.value = blackCount.value
    finalWhite.value = whiteCount.value

    const gameResult = await endGame(0, '认输', {
      black: finalBlack.value,
      white: finalWhite.value
    })

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
  color: var(--font-color);
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
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.3s;
}

.score-item.active {
  background: var(--el-color-primary);
  color: white;
}

.disc-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #1a1a1a;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.disc-icon.white {
  background: #f5f5f5;
  border: 1px solid #ccc;
}

.player-label {
  font-size: 12px;
  font-weight: 600;
}

.score-value {
  font-size: 20px;
  font-weight: 700;
}

.turn-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.turn-indicator {
  font-size: 11px;
  font-weight: 600;
  color: var(--el-color-primary);
}

.board-container {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.othello-board {
  background: #2d5a27;
  padding: 12px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.board-row {
  display: flex;
}

.board-cell {
  width: 38px;
  height: 38px;
  border: 1px solid #1e4d1a;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
}

.board-cell.valid-move {
  background: rgba(var(--success-color-rgb), 0.2);
}

.board-cell.valid-move:hover {
  background: rgba(var(--success-color-rgb), 0.4);
}

.board-cell.last-move {
  background: rgba(255, 255, 0, 0.2);
}

.disc {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
  transition: all 0.3s;
}

.disc.black {
  background: radial-gradient(circle at 30% 30%, #333, #000);
}

.disc.white {
  background: radial-gradient(circle at 30% 30%, #fff, #ddd);
}

.hint {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 0, 0, 0.5);
}

.game-status {
  margin-bottom: 10px;
  height: 24px;
}

.no-move-msg {
  font-size: 12px;
  color: var(--warning-color);
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
