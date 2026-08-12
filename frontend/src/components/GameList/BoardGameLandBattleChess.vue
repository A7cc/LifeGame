<template>
  <el-dialog v-model="visible" title="军旗" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">
          标准军棋翻棋模式<br>
          <span class="rules">翻棋 → 移动 → 吃子 | 铁路直线/工兵转弯 | 行营保护</span>
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
            <div class="score-item">
              <span class="player-label player-red">🔴 你</span>
              <span class="score-value">{{ playerPieces }}</span>
            </div>
            <div class="turn-info">
              <span class="turn-indicator" :class="{ 'player-turn': currentPlayer === 1 }">
                {{ currentPlayer === 1 ? '你的回合' : 'AI回合' }}
              </span>
              <span class="move-count">第 {{ moveNumber }} 手</span>
            </div>
            <div class="score-item">
              <span class="player-label player-blue">🔵 AI</span>
              <span class="score-value">{{ aiPieces }}</span>
            </div>
          </div>
        </div>

        <div class="board-container">
          <svg class="board-svg" viewBox="0 0 300 520">
            <!-- 公路线 -->
            <g class="highways">
              <line v-for="h in highways" :key="'h'+h.start" :x1="h.x1" :y1="h.y1" :x2="h.x2" :y2="h.y2" stroke="#999" stroke-width="1"/>
            </g>
            <!-- 铁路线 -->
            <g class="railways">
              <line v-for="(r, i) in railways" :key="'r'+i" :x1="r.x1" :y1="r.y1" :x2="r.x2" :y2="r.y2" :stroke="r.color" stroke-width="4"/>
            </g>
            <!-- 兵站/行营/大本营 -->
            <g class="stations">
              <g v-for="s in stations" :key="s.loc" :transform="s.transform">
                <rect v-if="s.type === 'station'" x="0" y="0" :width="s.width" :height="s.height" class="station"/>
                <rect v-if="s.type === 'hq'" x="0" y="0" :width="s.width" :height="s.height" class="station-hq"/>
                <ellipse v-if="s.type === 'camp'" :cx="s.cx" :cy="s.cy" :rx="s.rx" :ry="s.ry" class="station-camp"/>
              </g>
            </g>
            <!-- 有效移动标记 -->
            <g class="valid-moves">
              <circle v-for="m in validMoves" :key="'m'+m" :cx="getPieceX(m)+PIECE_WIDTH/2" :cy="getPieceY(m)+PIECE_HEIGHT/2" r="6" class="valid-mark"/>
            </g>
            <!-- 选中标记 -->
            <g v-if="selectedLoc !== null" class="selection">
              <rect :x="getPieceX(selectedLoc)" :y="getPieceY(selectedLoc)" :width="PIECE_WIDTH" :height="PIECE_HEIGHT" class="selected-mark"/>
            </g>
            <!-- 棋子（最上层） -->
            <g class="pieces">
              <g v-for="p in pieces" :key="p.loc" :transform="p.transform">
                <rect v-if="p.piece" :width="p.width" :height="p.height" :class="getPieceClass(p)" @click.stop="handleCellClick(p.loc)" style="cursor: pointer"/>
                <text v-if="p.piece && p.piece.revealed" :x="p.width/2" :y="p.height/2" class="piece-text" :class="p.piece.owner === 1 ? 'player' : 'ai'" style="pointer-events: none">
                  {{ p.piece.name }}
                </text>
                <text v-else-if="p.piece && !p.piece.revealed" :x="p.width/2" :y="p.height/2" class="piece-hidden" style="pointer-events: none">
                  ?
                </text>
              </g>
            </g>
          </svg>
        </div>

        <div class="game-legend">
          <span class="legend-item">司>军>师>旅>团>营>连>排>工</span>
          <span class="legend-item">炸弹同归于尽</span>
          <span class="legend-item">工兵挖雷</span>
        </div>

        <div class="action-buttons">
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
import { ref, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  config: { type: Object, default: () => ({ id: 'landbattlechess', name: '军旗', entryCost: 300 }) }
})
const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => { visible.value = newVal; if (newVal) resetLocal() })
watch(visible, (newVal) => { emit('update:modelValue', newVal) })

// 使用小游戏基础逻辑
const { gameStarted, gameEnded, processing, startGame, endGame, reset } = useMiniGameBase(props.config)

// ==================== 棋盘常量 ====================
// 位置编号：loc = row * 5 + col，左下角为0，右上角为59
// 棋盘尺寸参数
const CELL_WIDTH = 48
const CELL_HEIGHT = 32
const PADDING_X = 30
const PADDING_Y = 30
const GAP_Y = 10 // 中间山界间隔
const PIECE_WIDTH = 36
const PIECE_HEIGHT = 24

// 棋盘邻接图（POSITION_GRAPH）
const POSITION_GRAPH = new Array(60)
POSITION_GRAPH[0] = [5, 1]
POSITION_GRAPH[1] = [6, 2, 0]
POSITION_GRAPH[2] = [7, 3, 1]
POSITION_GRAPH[3] = [8, 4, 2]
POSITION_GRAPH[4] = [9, 3]
POSITION_GRAPH[5] = [10, 11, 6, 0]
POSITION_GRAPH[6] = [11, 7, 1, 5]
POSITION_GRAPH[7] = [12, 13, 8, 2, 6, 11]
POSITION_GRAPH[8] = [13, 9, 3, 7]
POSITION_GRAPH[9] = [14, 4, 8, 13]
POSITION_GRAPH[10] = [15, 11, 5]
POSITION_GRAPH[11] = [16, 17, 12, 7, 6, 5, 10, 15]
POSITION_GRAPH[12] = [17, 13, 7, 11]
POSITION_GRAPH[13] = [18, 19, 14, 9, 8, 7, 12, 17]
POSITION_GRAPH[14] = [19, 9, 13]
POSITION_GRAPH[15] = [20, 21, 16, 11, 10]
POSITION_GRAPH[16] = [21, 17, 11, 15]
POSITION_GRAPH[17] = [22, 23, 18, 13, 12, 11, 16, 21]
POSITION_GRAPH[18] = [23, 19, 13, 17]
POSITION_GRAPH[19] = [24, 14, 13, 18, 23]
POSITION_GRAPH[20] = [25, 21, 15]
POSITION_GRAPH[21] = [26, 27, 22, 17, 16, 15, 20, 25]
POSITION_GRAPH[22] = [27, 23, 17, 21]
POSITION_GRAPH[23] = [28, 29, 24, 19, 18, 17, 22, 27]
POSITION_GRAPH[24] = [29, 19, 23]
POSITION_GRAPH[25] = [30, 26, 21, 20]
POSITION_GRAPH[26] = [27, 21, 25]
POSITION_GRAPH[27] = [32, 28, 23, 22, 21, 26]
POSITION_GRAPH[28] = [29, 23, 27]
POSITION_GRAPH[29] = [34, 24, 23, 28]
POSITION_GRAPH[30] = [35, 36, 31, 25]
POSITION_GRAPH[31] = [36, 32, 30]
POSITION_GRAPH[32] = [37, 38, 33, 27, 31, 36]
POSITION_GRAPH[33] = [38, 34, 32]
POSITION_GRAPH[34] = [39, 29, 33, 38]
POSITION_GRAPH[35] = [40, 36, 30]
POSITION_GRAPH[36] = [41, 42, 37, 32, 31, 30, 35, 40]
POSITION_GRAPH[37] = [42, 38, 32, 36]
POSITION_GRAPH[38] = [43, 44, 39, 34, 33, 32, 37, 42]
POSITION_GRAPH[39] = [44, 34, 38]
POSITION_GRAPH[40] = [45, 46, 41, 36, 35]
POSITION_GRAPH[41] = [46, 42, 36, 40]
POSITION_GRAPH[42] = [47, 48, 43, 38, 37, 36, 41, 46]
POSITION_GRAPH[43] = [48, 44, 38, 42]
POSITION_GRAPH[44] = [49, 39, 38, 43, 48]
POSITION_GRAPH[45] = [50, 46, 40]
POSITION_GRAPH[46] = [51, 52, 47, 42, 41, 40, 45, 50]
POSITION_GRAPH[47] = [52, 48, 42, 46]
POSITION_GRAPH[48] = [53, 54, 49, 44, 43, 42, 47, 52]
POSITION_GRAPH[49] = [54, 44, 48]
POSITION_GRAPH[50] = [55, 51, 46, 45]
POSITION_GRAPH[51] = [56, 52, 46, 50]
POSITION_GRAPH[52] = [57, 53, 48, 47, 46, 51]
POSITION_GRAPH[53] = [58, 54, 48, 52]
POSITION_GRAPH[54] = [59, 49, 48, 53]
POSITION_GRAPH[55] = [56, 50]
POSITION_GRAPH[56] = [57, 51, 55]
POSITION_GRAPH[57] = [58, 52, 56]
POSITION_GRAPH[58] = [59, 53, 57]
POSITION_GRAPH[59] = [54, 58]

// 行营位置（邻接数为8）
const CAMP_LOCS = [11, 13, 17, 21, 23, 36, 38, 42, 46, 48]
// 大本营位置
const HQ_LOCS = [1, 3, 56, 58]

// 铁路顶点列表
const RAIL_VERTEXES = [5, 6, 7, 8, 9, 10, 14, 15, 19, 20, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 39, 40, 44, 45, 49, 50, 51, 52, 53, 54]

// 铁路邻接矩阵
const RAIL_EDGES = [
  [0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,1,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,1,0,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,1],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,1,0,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1,0],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,1],
  [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,1,0]
]

// ==================== 棋子定义 ====================
const PIECE_TYPES = [
  { name: '司令', rank: 9, count: 1 },
  { name: '军长', rank: 8, count: 1 },
  { name: '师长', rank: 7, count: 2 },
  { name: '旅长', rank: 6, count: 2 },
  { name: '团长', rank: 5, count: 2 },
  { name: '营长', rank: 4, count: 2 },
  { name: '连长', rank: 3, count: 3 },
  { name: '排长', rank: 2, count: 3 },
  { name: '工兵', rank: 1, count: 3, canDefuseMine: true },
  { name: '炸弹', rank: -1, count: 2, isBomb: true },
  { name: '地雷', rank: -2, count: 3, isMine: true, immobile: true },
  { name: '军旗', rank: 0, count: 1, isFlag: true, immobile: true }
]

// ==================== 游戏状态 ====================
const currentPlayer = ref(1)
const moveNumber = ref(0)
const board = ref(new Array(60).fill(null))
const selectedLoc = ref(null)
const validMoves = ref([])
const lastMove = ref(null)

const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

// ==================== 棋盘绘制数据 ====================
const highways = computed(() => initHighways())
const railways = computed(() => initRailways())
const stations = computed(() => initStations())
const pieces = computed(() => initPieces())

const playerPieces = computed(() => board.value.filter(p => p && p.owner === 1).length)
const aiPieces = computed(() => board.value.filter(p => p && p.owner === 2).length)

// 位置转像素坐标
const locToPixel = (loc) => {
  const x = loc % 5
  const y = Math.floor(loc / 5)
  const px = PADDING_X + x * CELL_WIDTH
  // 上半部分(行5-11)和下半部分(行0-4)之间有山界
  let py
  if (y <= 5) {
    py = PADDING_Y + (11 - y) * CELL_HEIGHT + GAP_Y
  } else {
    py = PADDING_Y + (11 - y) * CELL_HEIGHT
  }
  return { x: px, y: py }
}

const getPieceX = (loc) => locToPixel(loc).x - PIECE_WIDTH / 2
const getPieceY = (loc) => locToPixel(loc).y - PIECE_HEIGHT / 2

const initHighways = () => {
  const highwayS = [0, 1, 2, 3, 4, 55, 56, 57, 58, 59, 0, 10, 15, 20, 35, 40, 45, 55, 5, 9, 7, 7, 15, 19, 30, 34, 32, 32, 40, 44]
  const highwayE = [5, 26, 27, 28, 9, 50, 31, 32, 33, 54, 4, 14, 19, 24, 39, 44, 49, 59, 29, 25, 15, 19, 27, 27, 54, 50, 40, 44, 52, 52]
  return highwayS.map((s, i) => {
    const p1 = locToPixel(s)
    const p2 = locToPixel(highwayE[i])
    return { start: s, x1: p1.x, y1: p1.y, x2: p2.x, y2: p2.y }
  })
}

const initRailways = () => {
  const railwayS = [5, 50, 54, 9, 25, 30, 27]
  const railwayE = [50, 54, 9, 5, 29, 34, 32]
  const result = []
  railwayS.forEach((s, i) => {
    const p1 = locToPixel(s)
    const p2 = locToPixel(railwayE[i])
    const len = Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2)
    const num = Math.floor(len / 8)
    for (let j = 0; j < num; j++) {
      const t1 = j / num
      const t2 = (j + 1) / num
      result.push({
        x1: p1.x + (p2.x - p1.x) * t1,
        y1: p1.y + (p2.y - p1.y) * t1,
        x2: p1.x + (p2.x - p1.x) * t2,
        y2: p1.y + (p2.y - p1.y) * t2,
        color: j % 2 === 0 ? '#111' : '#EAC611'
      })
    }
  })
  return result
}

const initStations = () => {
  return Array.from({ length: 60 }, (_, loc) => {
    const p = locToPixel(loc)
    const w = CELL_WIDTH - 6
    const h = CELL_HEIGHT - 4
    let type = 'station'
    if (HQ_LOCS.includes(loc)) type = 'hq'
    else if (CAMP_LOCS.includes(loc)) type = 'camp'
    return {
      loc,
      type,
      transform: `translate(${p.x - w/2}, ${p.y - h/2})`,
      width: w,
      height: h,
      cx: w / 2,
      cy: h / 2,
      rx: w / 2 + 3,
      ry: h / 2 + 3
    }
  })
}

const initPieces = () => {
  return Array.from({ length: 60 }, (_, loc) => {
    const p = locToPixel(loc)
    const w = PIECE_WIDTH
    const h = PIECE_HEIGHT
    return {
      loc,
      piece: board.value[loc],
      transform: `translate(${p.x - w/2}, ${p.y - h/2})`,
      width: w,
      height: h
    }
  })
}

const getPieceClass = (p) => {
  if (!p.piece || !p.piece.revealed) return 'piece-back'
  return p.piece.owner === 1 ? 'piece-red' : 'piece-blue'
}

// ==================== 初始化 ====================
const resetLocal = () => {
  currentPlayer.value = 1
  moveNumber.value = 0
  board.value = new Array(60).fill(null)
  selectedLoc.value = null
  validMoves.value = []
  lastMove.value = null
  resultIcon.value = ''
  resultTitle.value = ''
  resultReward.value = ''
  reset()
}

const handleClose = () => { visible.value = false; resetLocal() }

// 开始游戏
const handleStartGame = async () => {
  const success = await startGame()
  if (success) {
    initBoard()
    ElMessage.info('游戏开始！点击棋子翻开')
  }
}

const initBoard = () => {
  const allPieces = []
  for (const type of PIECE_TYPES) {
    for (let i = 0; i < type.count; i++) {
      allPieces.push(createPiece(type, 1))
      allPieces.push(createPiece(type, 2))
    }
  }
  shuffle(allPieces)

  // 行营不能放棋子
  const validLocs = Array.from({ length: 60 }, (_, i) => i).filter(loc => !CAMP_LOCS.includes(loc))
  shuffle(validLocs)

  board.value = new Array(60).fill(null)
  for (let i = 0; i < allPieces.length && i < validLocs.length; i++) {
    board.value[validLocs[i]] = allPieces[i]
  }
}

const createPiece = (type, owner) => ({
  name: type.name,
  rank: type.rank,
  owner,
  revealed: false,
  canDefuseMine: type.canDefuseMine || false,
  isBomb: type.isBomb || false,
  isMine: type.isMine || false,
  isFlag: type.isFlag || false,
  immobile: type.immobile || false
})

const shuffle = (arr) => { for (let i = arr.length - 1; i > 0; i--) { const j = Math.floor(Math.random() * (i + 1)); [arr[i], arr[j]] = [arr[j], arr[i]] } }

// ==================== 交互处理 ====================
const handleCellClick = (loc) => {
  if (gameEnded.value || currentPlayer.value !== 1) return
  const piece = board.value[loc]

  // 如果已选中棋子，尝试移动
  if (selectedLoc.value !== null) {
    if (validMoves.value.includes(loc)) {
      makeMove(selectedLoc.value, loc)
    }
    selectedLoc.value = null
    validMoves.value = []
    return
  }

  // 翻开未翻棋子
  if (piece && !piece.revealed) {
    revealPiece(loc)
    return
  }

  // 选择己方已翻开的可移动棋子
  if (piece && piece.revealed && piece.owner === 1 && !piece.immobile) {
    selectedLoc.value = loc
    validMoves.value = getValidMoves(loc)
  }
}

const revealPiece = (loc) => {
  board.value[loc].revealed = true
  lastMove.value = { to: loc }
  moveNumber.value++
  switchTurn()
}

// ==================== 移动规则 ====================
const getValidMoves = (loc) => {
  const piece = board.value[loc]
  if (!piece || !piece.revealed || piece.immobile) return []
  const moves = []

  // 普通移动（相邻位置）
  for (const next of POSITION_GRAPH[loc]) {
    if (canMoveTo(next, piece)) moves.push(next)
  }

  // 铁路移动
  if (RAIL_VERTEXES.includes(loc)) {
    if (piece.canDefuseMine) {
      // 工兵：BFS转弯
      moves.push(...getEngineerRailMoves(loc, piece))
    } else {
      // 普通棋子：直线
      moves.push(...getNormalRailMoves(loc, piece))
    }
  }

  return [...new Set(moves)]
}

const canMoveTo = (loc, piece) => {
  const target = board.value[loc]
  // 空格可以移动
  if (!target) return true
  // 不能进入己方棋子
  if (target.owner === piece.owner) return false
  // 行营中有棋子时不能进入（行营保护）
  if (CAMP_LOCS.includes(loc)) return false
  // 可以攻击敌方棋子（包括未翻开的，翻开后再判断胜负）
  return true
}

const getNormalRailMoves = (loc, piece) => {
  const moves = []
  const sx = loc % 5, sy = Math.floor(loc / 5)
  const directions = [[0, 1], [0, -1], [1, 0], [-1, 0]]

  for (const [dx, dy] of directions) {
    let nx = sx + dx, ny = sy + dy
    while (nx >= 0 && nx < 5 && ny >= 0 && ny < 12) {
      const nloc = ny * 5 + nx
      if (!RAIL_VERTEXES.includes(nloc)) break
      const target = board.value[nloc]
      if (target) {
        if (!target.revealed || target.owner === piece.owner || CAMP_LOCS.includes(nloc)) break
        moves.push(nloc)
        break
      }
      moves.push(nloc)
      nx += dx
      ny += dy
    }
  }
  return moves
}

const getEngineerRailMoves = (loc, piece) => {
  const moves = []
  const visited = new Set([loc])
  const queue = [loc]

  while (queue.length > 0) {
    const cur = queue.shift()
    const curIdx = RAIL_VERTEXES.indexOf(cur)
    for (let i = 0; i < RAIL_EDGES[curIdx].length; i++) {
      if (RAIL_EDGES[curIdx][i] === 0) continue
      const nloc = RAIL_VERTEXES[i]
      if (visited.has(nloc)) continue
      visited.add(nloc)
      const target = board.value[nloc]
      if (target) {
        if (!target.revealed || target.owner === piece.owner || CAMP_LOCS.includes(nloc)) continue
        moves.push(nloc)
        continue
      }
      moves.push(nloc)
      queue.push(nloc)
    }
  }
  return moves
}

// ==================== 战斗逻辑 ====================
const makeMove = (from, to) => {
  const attacker = board.value[from]
  const defender = board.value[to]

  if (defender) {
    // 如果目标棋子未翻开，先翻开
    if (!defender.revealed) {
      defender.revealed = true
    }

    const result = resolveBattle(attacker, defender)
    if (result.bothDie) {
      board.value[from] = null
      board.value[to] = null
    } else if (result.attackerWins) {
      board.value[to] = attacker
      board.value[from] = null
    } else {
      board.value[from] = null
    }
    if (defender.isFlag) { finishGame(attacker.owner); return }
  } else {
    board.value[to] = attacker
    board.value[from] = null
  }

  lastMove.value = { from, to }
  moveNumber.value++
  if (checkGameEnd()) return
  switchTurn()
}

const resolveBattle = (attacker, defender) => {
  if (attacker.isBomb || defender.isBomb) return { bothDie: true }
  if (defender.isMine) return attacker.canDefuseMine ? { attackerWins: true } : { attackerWins: false }
  if (defender.isFlag) return { attackerWins: true }
  if (attacker.rank > defender.rank) return { attackerWins: true }
  if (attacker.rank === defender.rank) return { bothDie: true }
  return { attackerWins: false }
}

const checkGameEnd = () => {
  for (const player of [1, 2]) {
    if (!hasMovablePieces(player) && !hasHiddenPieces(player)) {
      finishGame(player === 1 ? 2 : 1)
      return true
    }
  }
  return false
}

const hasMovablePieces = (owner) => board.value.some(p => p && p.owner === owner && p.revealed && !p.immobile)
const hasHiddenPieces = (owner) => board.value.some(p => p && p.owner === owner && !p.revealed)

const switchTurn = () => {
  currentPlayer.value = currentPlayer.value === 1 ? 2 : 1
  if (currentPlayer.value === 2) setTimeout(aiMove, 500)
}

// ==================== AI ====================
const aiMove = () => {
  if (gameEnded.value) return
  const capture = findBestCaptureMove()
  if (capture) { makeMove(capture.from, capture.to); return }
  const move = findBestMove()
  if (move) { makeMove(move.from, move.to); return }
  const hidden = board.value.findIndex(p => p && !p.revealed)
  if (hidden >= 0) { revealPiece(hidden); return }
  finishGame(1)
}

const findBestCaptureMove = () => {
  let best = null, bestScore = 0
  for (let loc = 0; loc < 60; loc++) {
    const piece = board.value[loc]
    if (piece && piece.owner === 2 && piece.revealed && !piece.immobile) {
      for (const to of getValidMoves(loc)) {
        const target = board.value[to]
        if (target && target.owner === 1) {
          const score = target.isFlag ? 1000 : target.rank * 10 + 50
          if (score > bestScore) { bestScore = score; best = { from: loc, to } }
        }
      }
    }
  }
  return best
}

const findBestMove = () => {
  const moves = []
  for (let loc = 0; loc < 60; loc++) {
    const piece = board.value[loc]
    if (piece && piece.owner === 2 && piece.revealed && !piece.immobile) {
      for (const to of getValidMoves(loc)) {
        if (!board.value[to]) moves.push({ from: loc, to })
      }
    }
  }
  return moves.length > 0 ? moves[Math.floor(Math.random() * moves.length)] : null
}

// ==================== 游戏结束 ====================
const resign = () => {
  ElMessageBox.confirm('确定要认输吗？', '认输', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
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
.game-dialog { padding: 10px 0; }
.game-intro { text-align: center; }
.intro-text { font-size: 14px; color: var(--font-secondary); margin-bottom: 15px; line-height: 1.6; }
.rules { font-size: 12px; color: var(--font-light); }
.game-info { display: flex; justify-content: center; gap: 20px; margin-bottom: 20px; }
.info-item { font-size: 13px; padding: 6px 12px; background: var(--panel-color); border: 1px solid var(--border-color); border-radius: 6px; color: var(--font-color); }
.game-playing { text-align: center; }
.game-header { margin-bottom: 12px; }
.score-board { display: flex; justify-content: space-between; align-items: center; padding: 12px 20px; background: linear-gradient(135deg, var(--panel-color) 0%, rgba(var(--el-color-primary-rgb), 0.1) 100%); border: 1px solid var(--border-color); border-radius: 12px; }
.score-item { display: flex; flex-direction: column; gap: 4px; align-items: center; }
.player-label { font-size: 13px; font-weight: 600; }
.player-red { color: #c0392b; }
.player-blue { color: #2980b9; }
.score-value { font-size: 22px; font-weight: 700; color: var(--el-color-primary); }
.turn-info { display: flex; flex-direction: column; gap: 4px; align-items: center; }
.turn-indicator { font-size: 12px; font-weight: 600; color: var(--font-secondary); padding: 4px 12px; background: var(--panel-color); border-radius: 20px; }
.turn-indicator.player-turn { color: var(--el-color-primary); background: rgba(var(--el-color-primary-rgb), 0.1); }
.move-count { font-size: 11px; color: var(--font-light); }

.board-container { display: flex; justify-content: center; margin-bottom: 12px; }
.board-svg { width: 300px; height: 520px; background: #f5f0e1; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.2); }

.station { stroke: #666; stroke-width: 1; fill: #fff; }
.station-hq { stroke: #333; stroke-width: 2; fill: #e8e8e8; }
.station-camp { stroke: #2d6b30; stroke-width: 2; fill: #d4edda; }

.piece-back { fill: #d4c8b4; stroke: #8b7355; stroke-width: 2; cursor: pointer; }
.piece-red { fill: #fff; stroke: #e74c3c; stroke-width: 2; cursor: pointer; }
.piece-blue { fill: #fff; stroke: #3498db; stroke-width: 2; cursor: pointer; }
.piece-text { font-size: 11px; font-weight: bold; text-anchor: middle; dominant-baseline: central; pointer-events: none; }
.piece-text.player { fill: #c0392b; }
.piece-text.ai { fill: #2980b9; }
.piece-hidden { font-size: 16px; font-weight: bold; fill: #5d4e37; text-anchor: middle; dominant-baseline: central; pointer-events: none; }
.selected-mark { fill: none; stroke: #ffd700; stroke-width: 3; }
.valid-mark { fill: rgba(76, 175, 80, 0.6); }

.game-legend { display: flex; justify-content: center; flex-wrap: wrap; gap: 10px; margin-bottom: 12px; font-size: 11px; color: var(--font-secondary); }
.legend-item { padding: 4px 10px; background: var(--panel-color); border-radius: 20px; border: 1px solid var(--border-color); }
.action-buttons { display: flex; gap: 10px; justify-content: center; }
.game-result { text-align: center; padding: 20px 0; }
.result-icon { font-size: 64px; margin-bottom: 10px; }
.result-title { font-size: 20px; font-weight: 600; color: var(--font-color); margin-bottom: 15px; }
.result-reward { font-size: 14px; color: var(--success-color); }
</style>
