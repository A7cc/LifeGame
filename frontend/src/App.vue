<template>
  <div class="wrapper">
    <div v-if="!startupChecked || retryingStartup" class="startup-state" data-testid="startup-loading">
      <div class="startup-card">
        <div class="startup-icon">⏳</div>
        <h2>正在初始化游戏</h2>
        <p>{{ startupStage || '正在连接后端服务…' }}</p>
      </div>
    </div>

    <div v-else-if="startupError" class="startup-state" data-testid="startup-error">
      <div class="startup-card startup-error-card">
        <div class="startup-icon">⚠️</div>
        <h2>游戏初始化失败</h2>
        <p class="startup-stage">失败阶段：{{ startupStage || '未知阶段' }}</p>
        <p class="startup-error-message">{{ startupError }}</p>
        <el-button data-testid="retry-startup" type="primary" :loading="retryingStartup" @click="retryStartup">重新尝试</el-button>
      </div>
    </div>

    <!-- ==================== 起始页面 ==================== -->
    <div v-else-if="!started" class="start-screen" data-testid="start-screen">
      <!-- 背景装饰：浮动圆形和闪烁星星 -->
      <div class="background-decoration">
        <div class="decoration-circle circle-1"></div>
        <div class="decoration-circle circle-2"></div>
        <div class="decoration-circle circle-3"></div>
        <div class="decoration-star star-1">✨</div>
        <div class="decoration-star star-2">🌟</div>
        <div class="decoration-star star-3">💫</div>
      </div>

      <div class="start-content">
        <!-- 标题区域 -->
        <div class="title-section">
          <h1 class="start-title">🎮 LifeGame</h1>
          <p class="start-subtitle">开启你的人生模拟之旅</p>
        </div>

        <!-- 名字输入区域 -->
        <div class="input-section">
          <div class="input-wrapper">
            <span class="input-icon">👤</span>
            <input data-testid="player-name" v-model="playerName" class="name-input" placeholder="请输入你的名字" />
          </div>
        </div>

        <!-- 性别选择 -->
        <div class="section-label">选择性别</div>
        <div class="gender-select">
          <div
            data-testid="gender-male"
            role="button"
            tabindex="0"
            aria-label="选择男生"
            :aria-pressed="gender === true"
            class="gender-option"
            :class="{ selected: gender === true }"
            @click="gender = true"
            @keydown.enter.space.prevent="gender = true"
          >
            <span class="gender-icon">👨</span>
            <span class="gender-text">男生</span>
          </div>
          <div
            data-testid="gender-female"
            role="button"
            tabindex="0"
            aria-label="选择女生"
            :aria-pressed="gender === false"
            class="gender-option"
            :class="{ selected: gender === false }"
            @click="gender = false"
            @keydown.enter.space.prevent="gender = false"
          >
            <span class="gender-icon">👩</span>
            <span class="gender-text">女生</span>
          </div>
        </div>

        <!-- 难度选择 -->
        <div class="section-label">选择难度</div>
        <div class="difficulty-select">
          <div
            v-for="diff in difficulties"
            :key="diff.level"
            :data-testid="`difficulty-${diff.level}`"
            role="button"
            tabindex="0"
            :aria-label="`选择${diff.name}难度`"
            :aria-pressed="difficulty === diff.level"
            class="difficulty-option"
            :class="{ selected: difficulty === diff.level }"
            @click="difficulty = diff.level"
            @keydown.enter.space.prevent="difficulty = diff.level"
          >
            <div class="diff-icon">{{ diff.icon }}</div>
            <div class="diff-name">{{ diff.name }}</div>
            <div class="diff-desc">{{ diff.description }}</div>
            <div class="diff-money">💰 {{ formatMoney(diff.initmoney) }}</div>
          </div>
        </div>

        <!-- 开始游戏按钮 -->
        <div class="start-buttons">
          <el-button data-testid="start-game" type="primary" size="large" class="start-btn" @click="startGame">
            <span>开始游戏</span>
            <span class="btn-arrow">→</span>
          </el-button>
          <el-button data-testid="open-load-dialog" type="success" size="large" class="load-btn" @click="openLoadDialog">
            <span>📂 加载存档</span>
          </el-button>
        </div>
      </div>
    </div>

    <!-- ==================== 游戏主页面 ==================== -->
    <div v-else data-testid="game-screen">
      <GameMain @exit="exitGame" />
    </div>

    <DialogLoadGame
      v-model="showLoadDialog"
      :saves="saveList"
      :loading="loadingSaves"
      @load="loadSave"
    />

    <DialogGameEvaluation
      v-model="showEvaluation"
      :evaluation="evaluationData"
      @restart="confirmExit"
    />
  </div>
</template>

<script setup>
import { defineAsyncComponent, ref, onBeforeUnmount, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useThemeStore } from '@/src/stores/theme'
import { activateBackgroundMusic, deactivateBackgroundMusic } from '@/src/composables/useBackgroundMusic'
import DialogGameEvaluation from './components/Dialog/DialogGameEvaluation.vue'
import DialogLoadGame from './components/Dialog/DialogLoadGame.vue'
import { GetStartupStatus, RetryStartup, InitGame, EndGame, ListSaves, LoadGame } from "@/wailsjs/go/services/App.js"

const GameMain = defineAsyncComponent(() => import('./components/GameMain.vue'))

// ==================== 页面状态 ====================
const started = ref(false)           // 是否进入游戏界面
const playerName = ref('')           // 用户名字
const gender = ref(true)             // 性别：true=男, false=女
const difficulty = ref(1)            // 难度：0=简单, 1=普通, 2=困难
const gameStore = useGameStore()     // 游戏数据存储
const themeStore = useThemeStore()   // 主题状态存储
const startupChecked = ref(false)
const retryingStartup = ref(false)
const startupStage = ref('')
const startupError = ref('')

const applyStartupStatus = (status = {}) => {
  startupStage.value = status.stage || ''
  startupError.value = status.ready ? '' : (status.error || '后端尚未完成初始化')
  startupChecked.value = true
}

const checkStartup = async () => {
  try {
    const result = await GetStartupStatus()
    applyStartupStatus(result?.status)
  } catch (error) {
    startupStage.value = '连接后端'
    startupError.value = error?.message || String(error)
    startupChecked.value = true
  }
}

const retryStartup = async () => {
  retryingStartup.value = true
  try {
    const result = await RetryStartup()
    applyStartupStatus(result?.status)
  } catch (error) {
    startupStage.value = '重新初始化'
    startupError.value = error?.message || String(error)
    startupChecked.value = true
  } finally {
    retryingStartup.value = false
  }
}

// ==================== 初始化主题 ====================
onMounted(async () => {
  themeStore.initTheme()
  activateBackgroundMusic()
  await checkStartup()
})

onBeforeUnmount(() => {
  themeStore.disposeThemeListener()
  deactivateBackgroundMusic()
})

// ==================== 游戏评价 ====================
const showEvaluation = ref(false)    // 是否显示评价对话框
const evaluationData = ref(null)     // 评价数据

// ==================== 加载存档 ====================
const showLoadDialog = ref(false)    // 是否显示加载存档对话框
const saveList = ref([])             // 存档列表
const loadingSaves = ref(false)      // 加载存档列表状态

// 打开加载存档对话框
const openLoadDialog = async () => {
  showLoadDialog.value = true
  loadingSaves.value = true
  try {
    const result = await ListSaves()
    if (result.code === 200) {
      saveList.value = result.saves || []
    } else {
      ElMessage.error(result.msg || '获取存档列表失败')
    }
  } catch (err) {
    console.error('获取存档列表失败:', err)
    ElMessage.error('获取存档列表失败')
  } finally {
    loadingSaves.value = false
  }
}

// 加载存档
const loadSave = async (saveId) => {
  try {
    const result = await LoadGame(saveId)
    if (result.code === 200) {
      gameStore.applyGameData(result)
      playerName.value = result.userinfo?.uname || ''
      gender.value = result.userinfo?.usex ?? true
      showLoadDialog.value = false
      started.value = true
      ElMessage.success('存档加载成功')
    } else {
      ElMessage.error(result.msg || '加载存档失败')
    }
  } catch (err) {
    console.error('加载存档失败:', err)
    ElMessage.error('加载存档失败')
  }
}

// ==================== 难度配置 ====================
const difficulties = ref([
  {
    level: 0,
    name: '简单',
    description: '健康波动小，破产概率低，适合新手',
    initmoney: 1000000,
    icon: '🌱'
  },
  {
    level: 1,
    name: '普通',
    description: '标准平衡，适合有经验的玩家',
    initmoney: 300000,
    icon: '⚖️'
  },
  {
    level: 2,
    name: '困难',
    description: '健康波动大，破产概率高，极具挑战性',
    initmoney: 50000,
    icon: '🔥'
  },
])

// ==================== 工具函数 ====================
/**
 * 格式化金额显示
 * @param {number} money - 金额数值
 * @returns {string} 格式化后的金额字符串
 */
const formatMoney = (money) => {
  if (money >= 10000) {
    return (money / 10000).toFixed(0) + '万元'
  }
  return money.toLocaleString() + '元'
}

// ==================== 游戏流程 ====================
/**
 * 开始游戏，初始化用户数据
 */
const startGame = async () => {
  if (!playerName.value) {
    ElMessage.error('请输入你的名字')
    return
  }
  try {
    // 调用后端初始化游戏
    const gamedata = await InitGame(playerName.value, gender.value, difficulty.value)
    if (gamedata.code == 200) {
      // 保存游戏数据到store
      gameStore.applyGameData(gamedata)
      // 进入游戏界面
      started.value = true
    } else {
      ElMessage.error(gamedata.msg || '初始化失败')
      confirmExit()
    }
  } catch (err) {
    console.error('调用 InitGame 异常：', err)
    ElMessage.error('网络或服务器错误')
    confirmExit()
  }
}

/**
 * 结束游戏，显示评价
 */
const exitGame = async () => {
  const gamedata = await EndGame()
  if (gamedata.code == 200) {
    // 保存评价数据并显示对话框
    evaluationData.value = gamedata.evaluation
    showEvaluation.value = true
  } else {
    ElMessage.error(gamedata.msg || '游戏结束失败')
    confirmExit()
  }
}

/**
 * 重置游戏数据到初始状态
 */
const resetGameData = () => {
  // 重置游戏信息
  gameStore.resetGameState()
  // 返回起始页面
  started.value = false
}

/**
 * 确认退出并重新开始
 */
const confirmExit = () => {
  showEvaluation.value = false
  resetGameData()
}
</script>

<style scoped>
/* ==================== 全局容器 ==================== */
.wrapper {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #071225 0%, #0b1f3a 50%, #123a63 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  animation: fadeIn 0.6s ease-in-out;
  position: relative;
  overflow: hidden;
}

.startup-state {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  box-sizing: border-box;
}

.startup-card {
  width: min(520px, 100%);
  padding: 36px;
  border-radius: 18px;
  text-align: center;
  color: #f5f7fa;
  background: rgba(20, 24, 48, 0.92);
  border: 1px solid rgba(255, 255, 255, 0.16);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
}

.startup-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.startup-stage {
  color: #f0c78a;
}

.startup-error-message {
  margin: 16px 0 24px;
  padding: 12px;
  border-radius: 8px;
  color: #ffc9c9;
  background: rgba(245, 108, 108, 0.12);
  overflow-wrap: anywhere;
}

/* ==================== 背景装饰 ==================== */
.background-decoration {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s infinite ease-in-out;
}

.circle-1 { width: 300px; height: 300px; background: #3b82f6; top: -100px; right: -100px; animation-delay: 0s; }
.circle-2 { width: 200px; height: 200px; background: #1d4ed8; bottom: -50px; left: -50px; animation-delay: 5s; }
.circle-3 { width: 150px; height: 150px; background: #f59e0b; top: 50%; left: 10%; animation-delay: 10s; }

.decoration-star {
  position: absolute;
  font-size: 24px;
  opacity: 0.6;
  animation: twinkle 3s infinite ease-in-out;
}

.star-1 { top: 15%; right: 20%; animation-delay: 0s; }
.star-2 { top: 25%; left: 15%; animation-delay: 1s; }
.star-3 { bottom: 30%; right: 10%; animation-delay: 2s; }

/* ==================== 起始页面 ==================== */
.start-screen {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  z-index: 1;
}

.start-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  animation: slideUp 0.8s ease-out;
}

/* ----- 标题区域 ----- */
.title-section {
  text-align: center;
  margin-bottom: 40px;
}

.start-title {
  font-size: 56px;
  font-weight: 700;
  background: linear-gradient(135deg, #ffe066 0%, #ff922b 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 10px;
  filter: drop-shadow(2px 2px 4px rgba(0, 0, 0, 0.3));
}

.start-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 2px;
}

/* ----- 输入区域 ----- */
.input-section { margin-bottom: 25px; }

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  font-size: 20px;
  opacity: 0.6;
}

.name-input {
  padding: 14px 20px 14px 50px;
  border-radius: 16px;
  font-size: 16px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  outline: none;
  width: 320px;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.name-input::placeholder { color: rgba(255, 255, 255, 0.5); }

.name-input:focus {
  border-color: #60a5fa;
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 0 20px rgb(37 99 235 / 35%);
}

.section-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 12px;
  letter-spacing: 1px;
  font-weight: 500;
}

/* ----- 性别选择 ----- */
.gender-select {
  display: flex;
  gap: 16px;
  margin-bottom: 30px;
}

.gender-option {
  padding: 16px 28px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.gender-option:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
}

.gender-option.selected {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  border-color: #60a5fa;
  box-shadow: 0 8px 25px rgb(37 99 235 / 42%);
}

.gender-icon { font-size: 24px; }

.gender-text {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.gender-option.selected .gender-text { color: #fff; }

/* ----- 难度选择 ----- */
.difficulty-select {
  display: flex;
  gap: 12px;
  margin-bottom: 35px;
}

.difficulty-option {
  padding: 20px 18px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
  cursor: pointer;
  min-width: 140px;
  text-align: center;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.difficulty-option:hover {
  background: rgba(255, 255, 255, 0.12);
  transform: translateY(-4px);
}

.difficulty-option.selected {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  border-color: #60a5fa;
  box-shadow: 0 8px 25px rgb(37 99 235 / 42%);
}

.difficulty-option.selected .diff-name,
.difficulty-option.selected .diff-desc,
.difficulty-option.selected .diff-money { color: #fff; }

.diff-icon { font-size: 32px; margin-bottom: 8px; }

.diff-name {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.9);
}

.diff-desc {
  font-size: 11px;
  opacity: 0.7;
  margin-bottom: 8px;
  line-height: 1.4;
  color: rgba(255, 255, 255, 0.6);
  min-height: 32px;
}

.diff-money {
  font-size: 13px;
  font-weight: 600;
  opacity: 0.9;
  color: rgba(255, 255, 255, 0.8);
}

/* ----- 开始按钮 ----- */
.start-buttons {
  display: flex;
  gap: 16px;
  align-items: center;
}

.start-btn {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  color: white;
  border-radius: 16px;
  font-size: 18px;
  font-weight: 600;
  padding: 14px 48px;
  border: none;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.start-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgb(37 99 235 / 42%);
}

.btn-arrow { transition: transform 0.3s ease; }

.start-btn:hover .btn-arrow { transform: translateX(4px); }

.load-btn {
  background: linear-gradient(135deg, #67c23a 0%, #5daf34 100%);
  color: white;
  border-radius: 16px;
  font-size: 18px;
  font-weight: 600;
  padding: 14px 32px;
  border: none;
  transition: all 0.3s ease;
}

.load-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(103, 194, 58, 0.4);
}

/* ==================== 动画效果 ==================== */

/* 页面淡入 */
@keyframes fadeIn {
  0% { opacity: 0; }
  100% { opacity: 1; }
}

/* 内容上滑 */
@keyframes slideUp {
  0% { opacity: 0; transform: translateY(30px); }
  100% { opacity: 1; transform: translateY(0); }
}

/* 圆形浮动 */
@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

/* 星星闪烁 */
@keyframes twinkle {
  0%, 100% { opacity: 0.3; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.2); }
}
</style>
