<template>
  <el-dialog v-model="visible" title="国际象棋" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">与AI进行一场智慧的国际象棋对弈</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 500 }}元</span>
          <span class="info-item">🏆 奖励：最高1000元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始对弈</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-header">
          <div class="player-info">
            <span class="player-name">🤔 你</span>
            <span class="player-captured" v-if="capturedPieces.player.length > 0">
              {{ capturedPieces.player.join(' ') }}
            </span>
          </div>
          <div class="game-info-inline">
            <span class="turn-indicator">{{ currentPlayer === 'player' ? '你的回合' : 'AI思考中...' }}</span>
            <span class="move-count">第 {{ moveNumber }} 步</span>
          </div>
          <div class="player-info">
            <span class="player-name">🤖 AI</span>
            <span class="player-captured" v-if="capturedPieces.ai.length > 0">
              {{ capturedPieces.ai.join(' ') }}
            </span>
          </div>
        </div>

        <div class="chessboard-container">
          <div class="chessboard">
            <div
              v-for="(row, rowIndex) in board"
              :key="rowIndex"
              class="chess-row"
            >
              <div
                v-for="(cell, colIndex) in row"
                :key="colIndex"
                class="chess-cell"
                :class="{
                  'black-cell': (rowIndex + colIndex) % 2 === 1,
                  'white-cell': (rowIndex + colIndex) % 2 === 0,
                  'selected': selectedCell && selectedCell.row === rowIndex && selectedCell.col === colIndex,
                  'valid-move': isValidMove(rowIndex, colIndex),
                  'last-move': lastMove && ((lastMove.from.row === rowIndex && lastMove.from.col === colIndex) ||
                    (lastMove.to.row === rowIndex && lastMove.to.col === colIndex))
                }"
                @click="handleCellClick(rowIndex, colIndex)"
              >
                <span class="chess-piece" :class="cell.color">{{ cell.piece || '' }}</span>
                <span v-if="isValidMove(rowIndex, colIndex)" class="move-indicator">●</span>
              </div>
            </div>
          </div>
        </div>

        <div class="game-status">
          <div v-if="lastMoveText" class="last-move">{{ lastMoveText }}</div>
          <div class="status-text">剩余步数：{{ movesLeft }}</div>
        </div>

        <div class="action-buttons">
          <el-button @click="undoMove" :disabled="!canUndo || processing">↩️ 悔棋</el-button>
          <el-button @click="offerDraw" :disabled="processing">🤝 提议和棋</el-button>
          <el-button type="primary" @click="makeAIMove" :disabled="currentPlayer !== 'ai' || processing">跳过AI思考</el-button>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">总步数：{{ moveNumber }}</div>
          <div class="stat-item">吃子数：{{ capturedPieces.player.length }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'chess', name: '国际象棋', entryCost: 500 })
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

const INITIAL_BOARD = [
  ['♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'],
  ['♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'],
  ['', '', '', '', '', '', '', ''],
  ['', '', '', '', '', '', '', ''],
  ['', '', '', '', '', '', '', ''],
  ['', '', '', '', '', '', '', ''],
  ['♙', '♙', '♙', '♙', '♙', '♙', '♙', '♙'],
  ['♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖']
]

const PIECES = {
  '♜': { name: '车', color: 'black' }, '♞': { name: '马', color: 'black' },
  '♝': { name: '象', color: 'black' }, '♛': { name: '后', color: 'black' },
  '♚': { name: '王', color: 'black' }, '♟': { name: '兵', color: 'black' },
  '♖': { name: '车', color: 'white' }, '♘': { name: '马', color: 'white' },
  '♗': { name: '象', color: 'white' }, '♕': { name: '后', color: 'white' },
  '♔': { name: '王', color: 'white' }, '♙': { name: '兵', color: 'white' }
}

// 本地游戏状态
const movesLeft = ref(10)
const moveNumber = ref(0)
const currentPlayer = ref('player')
const board = ref([])
const selectedCell = ref(null)
const validMoves = ref([])
const lastMove = ref(null)
const lastMoveText = ref('')
const capturedPieces = ref({ player: [], ai: [] })
const canUndo = ref(false)
const moveHistory = ref([])
const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

const resetLocal = () => {
  movesLeft.value = 10
  moveNumber.value = 0
  currentPlayer.value = 'player'
  selectedCell.value = null
  validMoves.value = []
  lastMove.value = null
  lastMoveText.value = ''
  capturedPieces.value = { player: [], ai: [] }
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
    // 初始化棋盘
    board.value = INITIAL_BOARD.map(row =>
      row.map(cell => {
        if (!cell) return { piece: '', color: '' }
        const pieceInfo = PIECES[cell]
        return { piece: cell, color: pieceInfo.color, name: pieceInfo.name }
      })
    )
    ElMessage.info('对弈开始！你执白棋先行')
  }
}

const handleCellClick = (row, col) => {
  if (currentPlayer.value !== 'player' || processing.value) return

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

  if (cell.piece && cell.color === 'white') {
    selectedCell.value = { row, col }
    calculateValidMoves(row, col, cell)
  } else {
    selectedCell.value = null
    validMoves.value = []
  }
}

const calculateValidMoves = (row, col, piece) => {
  validMoves.value = []
  const moves = []

  switch (piece.name) {
    case '兵':
      if (row > 0 && !board.value[row - 1][col].piece) {
        moves.push({ row: row - 1, col })
        if (row === 6 && !board.value[row - 2][col].piece) {
          moves.push({ row: row - 2, col })
        }
      }
      if (row > 0 && col > 0 && board.value[row - 1][col - 1].piece && board.value[row - 1][col - 1].color === 'black') {
        moves.push({ row: row - 1, col: col - 1 })
      }
      if (row > 0 && col < 7 && board.value[row - 1][col + 1].piece && board.value[row - 1][col + 1].color === 'black') {
        moves.push({ row: row - 1, col: col + 1 })
      }
      break
    case '车':
      addLineMoves(row, col, [[0, 1], [0, -1], [1, 0], [-1, 0]], moves)
      break
    case '马':
      const knightMoves = [[-2, -1], [-2, 1], [-1, -2], [-1, 2], [1, -2], [1, 2], [2, -1], [2, 1]]
      knightMoves.forEach(([dr, dc]) => {
        const nr = row + dr, nc = col + dc
        if (nr >= 0 && nr < 8 && nc >= 0 && nc < 8) {
          if (!board.value[nr][nc].piece || board.value[nr][nc].color !== piece.color) {
            moves.push({ row: nr, col: nc })
          }
        }
      })
      break
    case '象':
      addLineMoves(row, col, [[1, 1], [1, -1], [-1, 1], [-1, -1]], moves)
      break
    case '后':
      addLineMoves(row, col, [[0, 1], [0, -1], [1, 0], [-1, 0], [1, 1], [1, -1], [-1, 1], [-1, -1]], moves)
      break
    case '王':
      const kingMoves = [[-1, -1], [-1, 0], [-1, 1], [0, -1], [0, 1], [1, -1], [1, 0], [1, 1]]
      kingMoves.forEach(([dr, dc]) => {
        const nr = row + dr, nc = col + dc
        if (nr >= 0 && nr < 8 && nc >= 0 && nc < 8) {
          if (!board.value[nr][nc].piece || board.value[nr][nc].color !== piece.color) {
            moves.push({ row: nr, col: nc })
          }
        }
      })
      break
  }

  validMoves.value = moves
}

const addLineMoves = (row, col, directions, moves) => {
  directions.forEach(([dr, dc]) => {
    let nr = row + dr, nc = col + dc
    while (nr >= 0 && nr < 8 && nc >= 0 && nc < 8) {
      if (board.value[nr][nc].piece) {
        if (board.value[nr][nc].color !== board.value[row][col].color) {
          moves.push({ row: nr, col: nc })
        }
        break
      }
      moves.push({ row: nr, col: nc })
      nr += dr
      nc += dc
    }
  })
}

const isValidMove = (row, col) => {
  return validMoves.value.some(m => m.row === row && m.col === col)
}

const makeMove = (fromRow, fromCol, toRow, toCol) => {
  const piece = board.value[fromRow][fromCol]
  const targetPiece = board.value[toRow][toCol]

  moveHistory.value.push({
    board: JSON.parse(JSON.stringify(board.value)),
    capturedPieces: JSON.parse(JSON.stringify(capturedPieces.value)),
    movesLeft: movesLeft.value
  })
  canUndo.value = true

  if (targetPiece.piece) {
    if (piece.color === 'white') {
      capturedPieces.value.player.push(targetPiece.piece)
    } else {
      capturedPieces.value.ai.push(targetPiece.piece)
    }
  }

  board.value[toRow][toCol] = piece
  board.value[fromRow][fromCol] = { piece: '', color: '', name: '' }

  lastMove.value = { from: { row: fromRow, col: fromCol }, to: { row: toRow, col: toCol } }
  lastMoveText.value = `${piece.piece} ${String.fromCharCode(97 + fromCol)}${8 - fromRow} → ${String.fromCharCode(97 + toCol)}${8 - toRow}`

  moveNumber.value++
  movesLeft.value--

  if (movesLeft.value <= 0) {
    finishGame()
    return
  }

  currentPlayer.value = currentPlayer.value === 'player' ? 'ai' : 'player'

  if (currentPlayer.value === 'ai') {
    setTimeout(() => aiMove(), 1000)
  }
}

const aiMove = () => {
  processing.value = true

  const blackPieces = []
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      const cell = board.value[r][c]
      if (cell.piece && cell.color === 'black') {
        blackPieces.push({ row: r, col: c, piece: cell })
      }
    }
  }

  if (blackPieces.length === 0) {
    processing.value = false
    finishGame()
    return
  }

  const randomPiece = blackPieces[Math.floor(Math.random() * blackPieces.length)]
  calculateValidMoves(randomPiece.row, randomPiece.col, randomPiece.piece)

  if (validMoves.value.length === 0) {
    processing.value = false
    aiMove()
    return
  }

  const randomMove = validMoves.value[Math.floor(Math.random() * validMoves.value.length)]
  makeMove(randomPiece.row, randomPiece.col, randomMove.row, randomMove.col)

  processing.value = false
}

const makeAIMove = () => {
  if (currentPlayer.value === 'ai' && !processing.value) {
    aiMove()
  }
}

const undoMove = () => {
  if (!canUndo.value || moveHistory.value.length === 0) return

  const lastState = moveHistory.value.pop()
  board.value = lastState.board
  capturedPieces.value = lastState.capturedPieces
  movesLeft.value = lastState.movesLeft
  moveNumber.value--
  currentPlayer.value = 'player'
  canUndo.value = moveHistory.value.length > 0

  ElMessage.info('已悔棋')
}

const offerDraw = async () => {
  const accept = Math.random() < 0.3
  if (accept) {
    const gameResult = await endGame(0, '和棋', { moves: moveNumber.value })
    if (gameResult) {
      resultIcon.value = '🤝'
      resultTitle.value = '和棋'
      resultReward.value = gameResult.resultText
      emit('complete', gameResult)
    }
  } else {
    ElMessage.info('AI拒绝和棋提议')
  }
}

const finishGame = async () => {
  const playerScore = capturedPieces.value.player.length * 2
  const aiScore = capturedPieces.value.ai.length * 2
  const isWin = playerScore > aiScore || (playerScore === aiScore && Math.random() < 0.5)

  const winCount = isWin ? 1 : 0
  const customResultText = isWin
    ? `获胜，吃子 ${capturedPieces.value.player.length} 个`
    : '惜败'

  const gameResult = await endGame(winCount, customResultText, {
    moves: moveNumber.value,
    captured: capturedPieces.value.player.length
  })

  if (gameResult) {
    resultIcon.value = isWin ? '🏆' : '😢'
    resultTitle.value = isWin ? '恭喜获胜！' : '惜败！'
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

.player-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.player-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
}

.player-captured {
  font-size: 14px;
  color: var(--font-secondary);
}

.game-info-inline {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.turn-indicator {
  font-size: 11px;
  font-weight: 600;
  color: var(--el-color-primary);
}

.move-count {
  font-size: 10px;
  color: var(--font-secondary);
}

.chessboard-container {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.chessboard {
  display: grid;
  grid-template-columns: repeat(8, 38px);
  grid-template-rows: repeat(8, 38px);
  border: 3px solid #5c4033;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.chess-cell {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  cursor: pointer;
  position: relative;
  box-sizing: border-box;
  overflow: hidden;
  flex-shrink: 0;
}

.white-cell {
  background: #f0d9b5;
}

.black-cell {
  background: #b58863;
}

.chess-cell:hover {
  background: rgba(255, 255, 0, 0.2);
}

.chess-cell.selected {
  background: rgba(255, 255, 0, 0.3) !important;
}

.chess-cell.valid-move {
  background: rgba(var(--success-color-rgb), 0.25) !important;
}

.chess-cell.last-move {
  background: rgba(255, 255, 0, 0.15) !important;
}

.chess-piece {
  user-select: none;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
  line-height: 1;
}

.move-indicator {
  position: absolute;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(var(--success-color-rgb), 0.5);
  pointer-events: none;
  z-index: 1;
}

.game-status {
  margin-bottom: 12px;
  min-height: 40px;
}

.last-move {
  font-size: 12px;
  color: var(--font-secondary);
  margin-bottom: 4px;
}

.status-text {
  font-size: 12px;
  color: var(--font-color);
  font-weight: 500;
}

.action-buttons {
  display: flex;
  gap: 8px;
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
  font-size: 14px;
  color: var(--font-secondary);
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
