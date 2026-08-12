<template>
  <el-dialog v-model="visible" title="斗兽棋" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">
          策略棋类游戏，兽类相克，渡河制胜<br>
          <span class="rules">象>狮>虎>豹>狼>狗>猫>鼠，鼠可吃象</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 400 }}元</span>
          <span class="info-item">🏆 奖励：800元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <span class="turn-indicator">{{ currentPlayer === 1 ? '你的回合' : 'AI思考中...' }}</span>
          <span class="move-count">第 {{ moveNumber }} 手</span>
        </div>

        <div class="board-container">
          <div class="jungle-board">
            <div
              v-for="(row, rowIndex) in board"
              :key="rowIndex"
              class="board-row"
            >
              <div
                v-for="(cell, colIndex) in row"
                :key="colIndex"
                class="board-cell"
                :class="[
                  cell.type === 'water' ? 'water' : 'land',
                  cell.type === 'trap' ? 'trap' : '',
                  cell.type === 'den' ? 'den' : '',
                  cell.type === 'trap' && cell.trapColor === 1 ? 'trap-red' : '',
                  cell.type === 'trap' && cell.trapColor === 2 ? 'trap-blue' : '',
                  cell.type === 'den' && cell.denColor === 1 ? 'den-red' : '',
                  cell.type === 'den' && cell.denColor === 2 ? 'den-blue' : '',
                  selectedCell && selectedCell.row === rowIndex && selectedCell.col === colIndex ? 'selected' : '',
                  isValidMove(rowIndex, colIndex) ? 'valid-move' : '',
                  lastMove && ((lastMove.from.row === rowIndex && lastMove.from.col === colIndex) ||
                    (lastMove.to.row === rowIndex && lastMove.to.col === colIndex)) ? 'last-move' : ''
                ]"
                @click="handleCellClick(rowIndex, colIndex)"
              >
                <span v-if="cell.piece" class="animal" :class="cell.piece.owner === 1 ? 'player' : 'ai'">
                  {{ cell.piece.emoji }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="game-legend">
          <span class="legend-item">🌊 河流</span>
          <span class="legend-item">🪤 陷阱</span>
          <span class="legend-item">🏠 兽穴</span>
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
    default: () => ({ id: 'jungle', name: '斗兽棋', entryCost: 400 })
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

// 棋盘配置 9x7
const ANIMALS = {
  elephant: { rank: 8, emoji: '🐘' },
  lion: { rank: 7, emoji: '🦁' },
  tiger: { rank: 6, emoji: '🐯' },
  leopard: { rank: 5, emoji: '🐆' },
  wolf: { rank: 4, emoji: '🐺' },
  dog: { rank: 3, emoji: '🐕' },
  cat: { rank: 2, emoji: '🐱' },
  rat: { rank: 1, emoji: '🐀' }
}

// 初始布局
const INITIAL_LAYOUT = [
  { animal: 'lion', pos: [0, 0] },
  { animal: 'tiger', pos: [0, 6] },
  { animal: 'dog', pos: [1, 1] },
  { animal: 'cat', pos: [1, 5] },
  { animal: 'rat', pos: [2, 0] },
  { animal: 'leopard', pos: [2, 2] },
  { animal: 'wolf', pos: [2, 4] },
  { animal: 'elephant', pos: [2, 6] }
]

// 河流位置
const WATER_CELLS = [
  [3, 1], [3, 2], [4, 1], [4, 2], [5, 1], [5, 2],
  [3, 4], [3, 5], [4, 4], [4, 5], [5, 4], [5, 5]
]

// 陷阱位置
const TRAPS = {
  1: [[0, 2], [0, 4], [1, 3]], // 玩家（红色）
  2: [[8, 2], [8, 4], [7, 3]]  // AI（蓝色）
}

// 兽穴位置
const DENS = {
  1: [0, 3], // 玩家
  2: [8, 3]  // AI
}

const currentPlayer = ref(1)
const moveNumber = ref(0)
const board = ref([])
const selectedCell = ref(null)
const validMoves = ref([])
const lastMove = ref(null)
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
  validMoves.value = []
  lastMove.value = null
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
    initBoard()
    ElMessage.info('游戏开始！你的兽穴在上方')
  }
}

const initBoard = () => {
  board.value = []

  for (let r = 0; r < 9; r++) {
    const row = []
    for (let c = 0; c < 7; c++) {
      const cell = { type: 'land', piece: null }

      // 标记河流
      if (WATER_CELLS.some(([wr, wc]) => wr === r && wc === c)) {
        cell.type = 'water'
      }

      // 标记陷阱
      for (const [player, traps] of Object.entries(TRAPS)) {
        if (traps.some(([tr, tc]) => tr === r && tc === c)) {
          cell.type = 'trap'
          cell.trapColor = parseInt(player)
        }
      }

      // 标记兽穴
      for (const [player, den] of Object.entries(DENS)) {
        if (den[0] === r && den[1] === c) {
          cell.type = 'den'
          cell.denColor = parseInt(player)
        }
      }

      row.push(cell)
    }
    board.value.push(row)
  }

  // 放置玩家的棋子（上方）
  INITIAL_LAYOUT.forEach(({ animal, pos }) => {
    board.value[pos[0]][pos[1]].piece = {
      type: animal,
      ...ANIMALS[animal],
      owner: 1
    }
  })

  // 放置AI的棋子（下方，镜像）
  INITIAL_LAYOUT.forEach(({ animal, pos }) => {
    const mirroredRow = 8 - pos[0]
    const mirroredCol = 6 - pos[1]
    board.value[mirroredRow][mirroredCol].piece = {
      type: animal,
      ...ANIMALS[animal],
      owner: 2
    }
  })
}

const handleCellClick = (row, col) => {
  if (currentPlayer.value !== 1 || gameEnded.value) return

  const cell = board.value[row][col]

  if (selectedCell.value) {
    const isValid = validMoves.value.some(m => m.row === row && m.col === col)
    if (isValid) {
      makeMove(selectedCell.value.row, selectedCell.value.col, row, col)
      selectedCell.value = null
      validMoves.value = []
      return
    }
  }

  if (cell.piece && cell.piece.owner === 1) {
    selectedCell.value = { row, col }
    calculateValidMoves(row, col)
  } else {
    selectedCell.value = null
    validMoves.value = []
  }
}

const calculateValidMoves = (row, col) => {
  validMoves.value = []
  const piece = board.value[row][col].piece

  const directions = [[0, 1], [0, -1], [1, 0], [-1, 0]]

  for (const [dr, dc] of directions) {
    let nr = row + dr
    let nc = col + dc

    // 检查边界
    if (nr < 0 || nr >= 9 || nc < 0 || nc >= 7) continue

    const targetCell = board.value[nr][nc]

    // 不能进入自己的兽穴
    if (targetCell.type === 'den' && targetCell.denColor === piece.owner) continue

    // 处理河流
    if (targetCell.type === 'water') {
      // 只有老鼠能进河
      if (piece.type === 'rat') {
        if (!targetCell.piece || (targetCell.piece.owner !== piece.owner)) {
          validMoves.value.push({ row: nr, col: nc })
        }
      }
      // 狮子和老虎可以跳过河
      else if (piece.type === 'lion' || piece.type === 'tiger') {
        // 检查河中是否有老鼠
        let jumpRow = nr
        let jumpCol = nc
        let blocked = false

        while (
          jumpRow >= 0 && jumpRow < 9 &&
          jumpCol >= 0 && jumpCol < 7 &&
          board.value[jumpRow][jumpCol].type === 'water'
        ) {
          if (board.value[jumpRow][jumpCol].piece) {
            blocked = true
            break
          }
          jumpRow += dr
          jumpCol += dc
        }

        if (!blocked && jumpRow >= 0 && jumpRow < 9 && jumpCol >= 0 && jumpCol < 7) {
          const landCell = board.value[jumpRow][jumpCol]
          if (!landCell.piece || (canCapture(piece, landCell.piece, row, col, jumpRow, jumpCol))) {
            validMoves.value.push({ row: jumpRow, col: jumpCol })
          }
        }
      }
      continue
    }

    // 普通移动
    if (!targetCell.piece) {
      // 从河里出来的老鼠不能吃岸上的象
      if (piece.type === 'rat' && board.value[row][col].type === 'water') {
        validMoves.value.push({ row: nr, col: nc })
      } else if (!targetCell.piece) {
        validMoves.value.push({ row: nr, col: nc })
      }
    } else if (targetCell.piece.owner !== piece.owner) {
      if (canCapture(piece, targetCell.piece, row, col, nr, nc)) {
        validMoves.value.push({ row: nr, col: nc })
      }
    }
  }
}

const canCapture = (attacker, defender, fromRow, fromCol, toRow, toCol) => {
  const attackerInWater = board.value[fromRow][fromCol].type === 'water'
  const defenderInWater = board.value[toRow][toCol].type === 'water'

  // 水里的老鼠不能吃岸上的动物
  if (attacker.type === 'rat' && attackerInWater && !defenderInWater) {
    return false
  }

  // 岸上的动物不能吃水里的老鼠
  if (!attackerInWater && defenderInWater) {
    return false
  }

  // 陷阱中的动物可以被任何动物吃
  const defenderCell = board.value[toRow][toCol]
  if (defenderCell.type === 'trap' && defenderCell.trapColor !== defender.owner) {
    return true
  }

  // 老鼠吃象
  if (attacker.type === 'rat' && defender.type === 'elephant') {
    return true
  }

  // 象不能吃老鼠
  if (attacker.type === 'elephant' && defender.type === 'rat') {
    return false
  }

  // 等级高或相同才能吃
  return attacker.rank >= defender.rank
}

const isValidMove = (row, col) => {
  return validMoves.value.some(m => m.row === row && m.col === col)
}

const makeMove = (fromRow, fromCol, toRow, toCol) => {
  const piece = board.value[fromRow][fromCol].piece

  // 保存历史
  moveHistory.value.push({
    board: JSON.parse(JSON.stringify(board.value)),
    currentPlayer: currentPlayer.value,
    moveNumber: moveNumber.value
  })
  canUndo.value = true

  // 移动棋子
  board.value[toRow][toCol].piece = piece
  board.value[fromRow][fromCol].piece = null

  lastMove.value = { from: { row: fromRow, col: fromCol }, to: { row: toRow, col: toCol } }
  moveNumber.value++

  // 检查是否进入对方兽穴
  const targetCell = board.value[toRow][toCol]
  if (targetCell.type === 'den' && targetCell.denColor !== piece.owner) {
    finishGame(piece.owner)
    return
  }

  // 检查是否吃掉对方所有棋子
  if (checkAllPiecesCaptured(piece.owner)) {
    finishGame(piece.owner)
    return
  }

  currentPlayer.value = currentPlayer.value === 1 ? 2 : 1
  if (currentPlayer.value === 2) {
    setTimeout(() => aiMove(), 600)
  }
}

const checkAllPiecesCaptured = (player) => {
  const opponent = player === 1 ? 2 : 1
  for (let r = 0; r < 9; r++) {
    for (let c = 0; c < 7; c++) {
      if (board.value[r][c].piece && board.value[r][c].piece.owner === opponent) {
        return false
      }
    }
  }
  return true
}

const aiMove = () => {
  const moves = []

  // 收集所有可能的移动
  for (let r = 0; r < 9; r++) {
    for (let c = 0; c < 7; c++) {
      const cell = board.value[r][c]
      if (cell.piece && cell.piece.owner === 2) {
        calculateValidMoves(r, c)
        if (validMoves.value.length > 0) {
          validMoves.value.forEach(move => {
            moves.push({
              from: { row: r, col: c },
              to: move,
              score: evaluateMove(r, c, move.row, move.col)
            })
          })
        }
      }
    }
  }

  if (moves.length === 0) {
    finishGame(1) // 玩家获胜
    return
  }

  // 按分数排序，选择最高分的移动
  moves.sort((a, b) => b.score - a.score)
  const bestMoves = moves.filter(m => m.score === moves[0].score)
  const selected = bestMoves[Math.floor(Math.random() * bestMoves.length)]

  makeMove(selected.from.row, selected.from.col, selected.to.row, selected.to.col)
}

const evaluateMove = (fromRow, fromCol, toRow, toCol) => {
  let score = Math.random() * 10 // 基础随机分

  const piece = board.value[fromRow][fromCol].piece
  const targetCell = board.value[toRow][toCol]

  // 吃子得分
  if (targetCell.piece) {
    score += targetCell.piece.rank * 20
  }

  // 靠近对方兽穴得分
  const enemyDen = DENS[1]
  const distBefore = Math.abs(fromRow - enemyDen[0]) + Math.abs(fromCol - enemyDen[1])
  const distAfter = Math.abs(toRow - enemyDen[0]) + Math.abs(toCol - enemyDen[1])
  if (distAfter < distBefore) {
    score += (distBefore - distAfter) * 5
  }

  // 进入对方兽穴得分
  if (targetCell.type === 'den' && targetCell.denColor === 1) {
    score += 1000
  }

  // 保护高等级棋子
  if (piece.rank >= 6) {
    score -= 5
  }

  return score
}

const undoMove = () => {
  if (!canUndo.value || moveHistory.value.length === 0) return

  if (moveHistory.value.length >= 2) {
    moveHistory.value.pop()
    const state = moveHistory.value.pop()
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
      resultTitle.value = 'AI获胜！'
      resultReward.value = gameResult.resultText
      emit('complete', gameResult)
    }
  }).catch(() => {})
}

const finishGame = async (winner) => {
  const winCount = winner === 1 ? 1 : 0
  const customResultText = winner === 1 ? '获胜' : '败北'

  const gameResult = await endGame(winCount, customResultText, {
    moves: moveNumber.value,
    winner
  })

  if (gameResult) {
    if (winner === 1) {
      resultIcon.value = '🏆'
      resultTitle.value = '恭喜获胜！'
    } else {
      resultIcon.value = '😢'
      resultTitle.value = 'AI获胜！'
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
  color: var(--font-light);
}

.board-container {
  display: flex;
  justify-content: center;
  margin-bottom: 10px;
}

.jungle-board {
  display: grid;
  grid-template-columns: repeat(7, 38px);
  grid-template-rows: repeat(9, 38px);
  gap: 1px;
  background: #8b6914;
  padding: 6px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.board-row {
  display: contents;
}

.board-cell {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  font-size: 20px;
  box-sizing: border-box;
  overflow: hidden;
  flex-shrink: 0;
}

.board-cell.land {
  background: #daa520;
}

.board-cell.water {
  background: #4a90d9;
}

.board-cell.trap {
  background: #cd853f;
}

.board-cell.trap-red {
  background: #e8947a;
}

.board-cell.trap-blue {
  background: #7ab0e8;
}

.board-cell.den {
  background: #ffd700;
  border: 2px solid #8b6914;
}

.board-cell.den-red {
  background: #ffb6c1;
}

.board-cell.den-blue {
  background: #add8e6;
}

.board-cell:hover {
  filter: brightness(1.1);
}

.board-cell.selected {
  box-shadow: inset 0 0 0 2px rgba(255, 255, 0, 0.8);
}

.board-cell.valid-move {
  box-shadow: inset 0 0 0 2px rgba(var(--success-color-rgb), 0.6);
}

.board-cell.last-move {
  filter: brightness(1.2);
}

.animal {
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  line-height: 1;
}

.animal.player {
  filter: drop-shadow(0 0 2px #ff4444);
}

.animal.ai {
  filter: drop-shadow(0 0 2px #4444ff);
}

.game-legend {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 10px;
  font-size: 11px;
  color: var(--font-secondary);
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
