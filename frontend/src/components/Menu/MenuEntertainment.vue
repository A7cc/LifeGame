<template>
  <div class="product-panel" data-testid="page-entertainment">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">娱乐中心</div>
    </div>

    <!-- 中间内容区域 -->
    <div class="entertainment-main">
      <!-- 左侧：游戏分类导航 -->
      <div class="gamenav">
        <div class="section-header">
          <span class="section-title">游戏分类</span>
        </div>

        <div class="game-categories">
          <div v-for="category in gameCategories" :key="category.id" :data-testid="`entertainment-category-${category.id}`" class="category-item" :class="{ active: currentCategory === category.id }" @click="switchCategory(category.id)">
            <div class="category-icon">{{ category.icon }}</div>
            <div class="category-info">
              <div class="category-name">{{ category.name }}</div>
              <div class="category-count">{{ category.count }}项</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧面板 -->
      <div class="right-panel">
        <!-- 游戏主区域 -->
        <div class="gamearena">
          <div class="section-header">
            <span class="section-title">{{ currentCategoryName }}</span>
          </div>

          <!-- 娱乐活动列表 -->
          <div v-if="currentCategory === 'activities'" class="minigames-list">
            <div v-for="activity in entertainmentGames" :key="activity.id" :data-testid="`entertainment-activity-${activity.id}`" class="minigames-row" :class="{ locked: userInfo.ucash < activity.entryCost }" @click="confirmActivity(activity)">
              <div class="minigames-icon">{{ activity.icon }}</div>
              <div class="minigames-info">
                <div class="minigames-name">{{ activity.name }}</div>
                <div class="minigames-desc">{{ activity.desc }}</div>
              </div>
              <div class="minigames-reward">
                <span class="minigames-cost">💰 {{ activity.entryCost.toLocaleString() }}</span>
                <div class="minigames-tag">
                  <span v-if="activity.healthGain > 0" class="minigames-health">💚 +{{ activity.healthGain }}</span>
                  <span v-if="activity.fameGain > 0" class="minigames-fame">⭐ +{{ activity.fameGain }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 休闲小游戏区域 -->
          <div v-else-if="currentCategory === 'casualgames'" class="minigames-list">
            <div v-for="game in casualGames" :key="game.id" class="minigames-row" @click="openGame(game)">
              <div class="minigames-icon">{{ game.icon }}</div>
              <div class="minigames-info">
                <div class="minigames-name">{{ game.name }}</div>
                <div class="minigames-desc">{{ game.desc }}</div>
              </div>
              <div class="minigames-reward">
                <span class="minigames-cost">💰 {{ game.entryCost.toLocaleString() }}</span>
                <span class="minigames-reward-text">{{ game.reward }}</span>
              </div>
            </div>
          </div>

          <!-- 棋牌游戏区域 -->
          <div v-else-if="currentCategory === 'boardgames'" class="minigames-list">
            <div v-for="game in boardGames" :key="game.id" class="minigames-row" @click="openGame(game)">
              <div class="minigames-icon">{{ game.icon }}</div>
              <div class="minigames-info">
                <div class="minigames-name">{{ game.name }}</div>
                <div class="minigames-desc">{{ game.desc }}</div>
              </div>
              <div class="minigames-reward">
                <span class="minigames-cost">💰 {{ game.entryCost.toLocaleString() }}</span>
                <span class="minigames-reward-text">{{ game.reward }}</span>
              </div>
            </div>
          </div>

          <!-- 竞技游戏区域 -->
          <div v-else-if="currentCategory === 'competitive'" class="minigames-list">
            <div v-for="game in competitiveGames" :key="game.id" class="minigames-row" @click="openGame(game)">
              <div class="minigames-icon">{{ game.icon }}</div>
              <div class="minigames-info">
                <div class="minigames-name">{{ game.name }}</div>
                <div class="minigames-desc">{{ game.desc }}</div>
              </div>
              <div class="minigames-reward">
                <span class="minigames-cost">💰 {{ game.entryCost.toLocaleString() }}</span>
                <span class="minigames-reward-text">{{ game.reward }}</span>
              </div>
            </div>
          </div>

          <!-- 博彩游戏区域 -->
          <div v-else-if="currentCategory === 'gambling'" class="minigames-list">
            <div v-for="game in gamblingGames" :key="game.id" class="minigames-row" @click="openGame(game)">
              <div class="minigames-icon">{{ game.icon }}</div>
              <div class="minigames-info">
                <div class="minigames-name">{{ game.name }}</div>
                <div class="minigames-desc">{{ game.desc }}</div>
              </div>
              <div class="minigames-reward">
                <span class="minigames-cost">💰 {{ game.entryCost.toLocaleString() }}元报名费</span>
                <span v-if="game.needBet" class="minigames-bet-tip">+ 下注金额</span>
                <span class="minigames-reward-text">{{ game.reward }}</span>
              </div>
            </div>
          </div>
        </div>

        <LogPanel title="📋 日志记录" :items="runLogInfo" empty-icon="🚴" empty-text="暂无活动记录" @clear="clearLogs" />
      </div>
    </div>
  </div>

  <!-- 确认参与活动对话框 -->
  <DialogActivityConfirm
    v-model="showActivityConfirm"
    :activity="selectedActivity"
    @confirm="doActivity"
  />

  <!-- 游戏组件 - 传入 config prop -->
  <GamblingGameRockPaperScissors v-if="showGamblingGameRockPaperScissors" v-model="showGamblingGameRockPaperScissors" :config="getGameConfig('rps')" @complete="onGameComplete" />
  <CasualGameGuessNumber v-if="showCasualGameGuessNumber" v-model="showCasualGameGuessNumber" :config="getGameConfig('guess')" @complete="onGameComplete" />
  <CasualGameDice v-if="showCasualGameDice" v-model="showCasualGameDice" :config="getGameConfig('dice')" @complete="onGameComplete" />
  <CasualGameSlotMachine v-if="showCasualGameSlotMachine" v-model="showCasualGameSlotMachine" :config="getGameConfig('slot')" @complete="onGameComplete" />
  <BoardGameChess v-if="showBoardGameChess" v-model="showBoardGameChess" :config="getGameConfig('chess')" @complete="onGameComplete" />
  <CompetitiveGameFPS v-if="showCompetitiveGameFPS" v-model="showCompetitiveGameFPS" :config="getGameConfig('fps')" @complete="onGameComplete" />
  <CompetitiveGameMOBA v-if="showCompetitiveGameMOBA" v-model="showCompetitiveGameMOBA" :config="getGameConfig('moba')" @complete="onGameComplete" />
  <CompetitiveGameRacing v-if="showCompetitiveGameRacing" v-model="showCompetitiveGameRacing" :config="getGameConfig('racing')" @complete="onGameComplete" />
  <CompetitiveGameFighting v-if="showCompetitiveGameFighting" v-model="showCompetitiveGameFighting" :config="getGameConfig('fighting')" @complete="onGameComplete" />
  <CompetitiveGameWar v-if="showCompetitiveGameWar" v-model="showCompetitiveGameWar" :config="getGameConfig('war')" @complete="onGameComplete" />
  <GamblingGamePoker v-if="showGamblingGamePoker" v-model="showGamblingGamePoker" :config="getGameConfig('poker')" @complete="onGameComplete" />
  <GamblingGameHorseRacing v-if="showGamblingGameHorseRacing" v-model="showGamblingGameHorseRacing" :config="getGameConfig('horseracing')" @complete="onGameComplete" />
  <GamblingGameRoulette v-if="showGamblingGameRoulette" v-model="showGamblingGameRoulette" :config="getGameConfig('roulette')" @complete="onGameComplete" />
  <GamblingGameBaccarat v-if="showGamblingGameBaccarat" v-model="showGamblingGameBaccarat" :config="getGameConfig('baccarat')" @complete="onGameComplete" />
  <GamblingGameBlackjack v-if="showGamblingGameBlackjack" v-model="showGamblingGameBlackjack" :config="getGameConfig('blackjack')" @complete="onGameComplete" />
  <GamblingGameLottery v-if="showGamblingGameLottery" v-model="showGamblingGameLottery" :config="getGameConfig('lottery')" @complete="onGameComplete" />
  <BoardGameGobang v-if="showBoardGameGobang" v-model="showBoardGameGobang" :config="getGameConfig('gobang')" @complete="onGameComplete" />
  <BoardGameJungle v-if="showBoardGameJungle" v-model="showBoardGameJungle" :config="getGameConfig('jungle')" @complete="onGameComplete" />
  <BoardGameGo v-if="showBoardGameGo" v-model="showBoardGameGo" :config="getGameConfig('go')" @complete="onGameComplete" />
  <BoardGameOthello v-if="showBoardGameOthello" v-model="showBoardGameOthello" :config="getGameConfig('othello')" @complete="onGameComplete" />
  <BoardGameLandBattleChess v-if="showBoardGameLandBattleChess" v-model="showBoardGameLandBattleChess" :config="getGameConfig('landbattle')" @complete="onGameComplete" />
</template>

<script setup>
import { defineAsyncComponent, ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useLogList } from '@/src/composables/useLogList'
import LogPanel from '@/src/components/Common/LogPanel.vue'
import DialogActivityConfirm from '@/src/components/Dialog/DialogActivityConfirm.vue'
import { DoEntertainment, GetEntertainmentActivities, GetMiniGameConfigs } from "@/wailsjs/go/services/App.js"

// 休闲游戏组件
const CasualGameGuessNumber = defineAsyncComponent(() => import('@/src/components/GameList/CasualGameGuessNumber.vue'))
const CasualGameDice = defineAsyncComponent(() => import('@/src/components/GameList/CasualGameDice.vue'))
const CasualGameSlotMachine = defineAsyncComponent(() => import('@/src/components/GameList/CasualGameSlotMachine.vue'))
// 棋牌游戏组件
const BoardGameGobang = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameGobang.vue'))
const BoardGameJungle = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameJungle.vue'))
const BoardGameGo = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameGo.vue'))
const BoardGameOthello = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameOthello.vue'))
const BoardGameLandBattleChess = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameLandBattleChess.vue'))
// 竞技游戏组件
const BoardGameChess = defineAsyncComponent(() => import('@/src/components/GameList/BoardGameChess.vue'))
const CompetitiveGameFPS = defineAsyncComponent(() => import('@/src/components/GameList/CompetitiveGameFPS.vue'))
const CompetitiveGameMOBA = defineAsyncComponent(() => import('@/src/components/GameList/CompetitiveGameMOBA.vue'))
const CompetitiveGameRacing = defineAsyncComponent(() => import('@/src/components/GameList/CompetitiveGameRacing.vue'))
const CompetitiveGameFighting = defineAsyncComponent(() => import('@/src/components/GameList/CompetitiveGameFighting.vue'))
const CompetitiveGameWar = defineAsyncComponent(() => import('@/src/components/GameList/CompetitiveGameWar.vue'))
// 博彩游戏组件
const GamblingGameRockPaperScissors = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameRockPaperScissors.vue'))
const GamblingGamePoker = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGamePoker.vue'))
const GamblingGameHorseRacing = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameHorseRacing.vue'))
const GamblingGameRoulette = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameRoulette.vue'))
const GamblingGameBaccarat = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameBaccarat.vue'))
const GamblingGameBlackjack = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameBlackjack.vue'))
const GamblingGameLottery = defineAsyncComponent(() => import('@/src/components/GameList/GamblingGameLottery.vue'))


// 游戏主要信息
const gameStore = useGameStore()
// 用户个人信息
const userInfo = computed(() => gameStore.userInfo)

const { addLog, clearLogs, logs: runLogInfo } = useLogList({ green: '日志', red: '主要' })

// 当前选中的分类
const currentCategory = ref('activities')

// 游戏分类
const gameCategories = computed(() => [
  { id: 'activities', name: '娱乐活动', icon: '🎬', count: entertainmentGames.value.length },
  { id: 'casualgames', name: '休闲游戏', icon: '🎮', count: casualGames.value.length },
  { id: 'boardgames', name: '棋牌游戏', icon: '♟️', count: boardGames.value.length },
  { id: 'competitive', name: '竞技游戏', icon: '🏆', count: competitiveGames.value.length },
  { id: 'gambling', name: '博彩游戏', icon: '🎰', count: gamblingGames.value.length },
])

// 当前分类名称
const currentCategoryName = computed(() => {
  const category = gameCategories.value.find(c => c.id === currentCategory.value)
  return category?.name || '娱乐活动'
})

// 切换分类
const switchCategory = (categoryId) => {
  currentCategory.value = categoryId
}

// 娱乐活动列表
const entertainmentGames = ref([])

// 确认参与活动对话框
const showActivityConfirm = ref(false)
const selectedActivity = ref(null)

// 确认参与活动
const confirmActivity = (activity) => {
  if (userInfo.value.ucash < activity.entryCost) {
    ElMessage.warning('资金不足')
    return
  }
  selectedActivity.value = activity
  showActivityConfirm.value = true
}

// 参与活动
const doActivity = async (activity) => {
  try {
    const result = await DoEntertainment(activity.id)
    if (result?.code !== 200) {
      ElMessage.error(result?.msg || '参与活动失败')
      return
    }

    gameStore.applyUserInfo(result.userinfo)
    const healthChange = Number(result.healthchange || 0)
    const fameChange = Number(result.famechange || 0)
    const messages = []
    if (healthChange !== 0) messages.push(`免疫力${healthChange > 0 ? '+' : ''}${healthChange}`)
    if (fameChange !== 0) messages.push(`名声${fameChange > 0 ? '+' : ''}${fameChange}`)
    ElMessage.success(`参与成功！${messages.join('，')}`)
    addLog(`你参加了${activity.name}，免疫力${healthChange > 0 ? '+' : ''}${healthChange}，名声${fameChange > 0 ? '+' : ''}${fameChange}`)
  } catch (error) {
    console.error('参与娱乐活动失败', error)
    ElMessage.error('参与活动失败')
  }
}

// ==================== 小游戏相关 ====================
// 游戏配置从后端动态获取
const miniGameConfigs = ref([])
const casualGames = ref([])
const boardGames = ref([])
const competitiveGames = ref([])
const gamblingGames = ref([])

const loadEntertainmentActivities = async () => {
  try {
    const result = await GetEntertainmentActivities()
    if (result?.code === 200) {
      entertainmentGames.value = result.activities || []
    } else {
      ElMessage.error(result?.msg || '获取娱乐活动失败')
    }
  } catch (error) {
    console.error('获取娱乐活动失败', error)
    ElMessage.error('获取娱乐活动失败')
  }
}

// 加载小游戏配置
const loadMiniGameConfigs = async () => {
  try {
    const result = await GetMiniGameConfigs()
    if (result?.code === 200 && result?.configs) {
      miniGameConfigs.value = result.configs

      // 按分类分组
      casualGames.value = result.configs.filter(c => c.category === 'casual')
      boardGames.value = result.configs.filter(c => c.category === 'board')
      competitiveGames.value = result.configs.filter(c => c.category === 'competitive')
      gamblingGames.value = result.configs.filter(c => c.category === 'gambling')
    }
  } catch (error) {
    console.error('获取小游戏配置失败', error)
  }
}

// 在组件初始化时加载配置
loadMiniGameConfigs()
loadEntertainmentActivities()

// 游戏对话框显示状态
// 休闲游戏开关
const showCasualGameGuessNumber = ref(false)
const showCasualGameDice = ref(false)
const showCasualGameSlotMachine = ref(false)
// 棋牌游戏开关
const showBoardGameGobang = ref(false)
const showBoardGameJungle = ref(false)
const showBoardGameGo = ref(false)
const showBoardGameOthello = ref(false)
const showBoardGameLandBattleChess = ref(false)
// 竞技游戏开关
const showBoardGameChess = ref(false)
const showCompetitiveGameFPS = ref(false)
const showCompetitiveGameMOBA = ref(false)
const showCompetitiveGameRacing = ref(false)
const showCompetitiveGameFighting = ref(false)
const showCompetitiveGameWar = ref(false)
// 博彩游戏开关
const showGamblingGameRockPaperScissors = ref(false)
const showGamblingGamePoker = ref(false)
const showGamblingGameHorseRacing = ref(false)
const showGamblingGameRoulette = ref(false)
const showGamblingGameBaccarat = ref(false)
const showGamblingGameBlackjack = ref(false)
const showGamblingGameLottery = ref(false)

const gameDialogRefs = {
  rps: showGamblingGameRockPaperScissors,
  guess: showCasualGameGuessNumber,
  dice: showCasualGameDice,
  slot: showCasualGameSlotMachine,
  gobang: showBoardGameGobang,
  jungle: showBoardGameJungle,
  go: showBoardGameGo,
  othello: showBoardGameOthello,
  landbattle: showBoardGameLandBattleChess,
  chess: showBoardGameChess,
  fps: showCompetitiveGameFPS,
  moba: showCompetitiveGameMOBA,
  racing: showCompetitiveGameRacing,
  fighting: showCompetitiveGameFighting,
  war: showCompetitiveGameWar,
  poker: showGamblingGamePoker,
  horseracing: showGamblingGameHorseRacing,
  roulette: showGamblingGameRoulette,
  baccarat: showGamblingGameBaccarat,
  blackjack: showGamblingGameBlackjack,
  lottery: showGamblingGameLottery,
}

// 获取游戏配置（用于传给小游戏组件）
const getGameConfig = (gameId) => {
  return miniGameConfigs.value.find(g => g.id === gameId) || { id: gameId, name: gameId, entryCost: 0 }
}

// 打开游戏（简化版：只打开对话框，不调用 StartMiniGame）
const openGame = (game) => {
  const dialogRef = gameDialogRefs[game.id]
  if (!dialogRef) {
    ElMessage.error('小游戏组件不存在')
    return
  }
  // 直接打开对话框，小游戏组件内部会调用 StartMiniGame
  dialogRef.value = true
}

// 游戏完成回调（简化版：只记录日志）
const onGameComplete = (result) => {
  if (!result) return
  const emoji = result.win ? '🎉' : '💔'
  let logText = `${emoji} ${result.gameName}：${result.resultText}`
  if (result.cashChange !== undefined && result.cashChange !== 0) {
    logText += ` | 净变化：${result.cashChange >= 0 ? '+' : ''}${result.cashChange}元`
  }
  addLog(logText)
}
</script>

<style scoped>
/* ==================== 中间内容区域 ==================== */
.entertainment-main {
  display: flex;
  gap: 5px;
  margin-top: 5px;
  flex: 1;
  min-height: 0;
}

/* ==================== 左侧：游戏分类导航 ==================== */
.gamenav {
  flex: 0 0 170px;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  background: var(--panel-color);
  padding: 10px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ==================== 右侧面板 ==================== */
.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
  min-width: 0;
}

/* ==================== 游戏主区域 ==================== */
.gamearena {
  flex: 1;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  background: var(--panel-color);
  padding: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0;
}

/* ==================== 游戏分类导航 ==================== */
.game-categories {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding-right: 4px;
  /* 隐藏滚动条 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.game-categories::-webkit-scrollbar {
  display: none;
}

/* 游戏分类项 */
.category-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 10px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.category-item:hover {
  background-color: var(--select-color);         /* 根据主题变化 */
}

.category-item.active {
  background-color: var(--select-color);         /* 根据主题变化 */
}

/* 游戏分类项->图标 */
.category-icon {
  font-size: 24px;
  flex-shrink: 0;
}
/* 游戏分类项->信息块 */
.category-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}
/* 游戏分类项->名字 */
.category-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--font-color);
}
/* 游戏分类项->游戏数目 */
.category-count {
  font-size: 11px;
  color: var(--font-secondary);
}

/* ==================== 小游戏列表 ==================== */
.minigames-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding-right: 4px;
  /* 隐藏滚动条 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.minigames-list::-webkit-scrollbar {
  display: none;
}

/* ==================== 小游戏样式 ==================== */
.minigames-row {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 5px 12px;
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.minigames-row:hover {
  background-color: var(--select-color);         /* 根据主题变化 */
}

.minigames-icon {
  font-size: 28px;
  flex-shrink: 0;
}

.minigames-info {
  flex: 1;
  min-width: 0;
}

.minigames-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--font-color);
  margin-bottom: 2px;
}

.minigames-desc {
  font-size: 10px;
  color: var(--font-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.minigames-reward {
  display: flex;
  flex-direction: column;
  gap: 5px;
  align-items: flex-end;
  flex-shrink: 0;
}

.minigames-cost {
  font-size: 10px;
  font-weight: 600;
  color: var(--warning-color);
}

.minigames-reward-text {
  font-size: 10px;
  color: var(--error-color);
  padding: 2px 6px;
  border-radius: 8px;
}

.minigames-tag {
  display: flex;
  gap: 5px;
}

.minigames-health {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 6px;
  font-weight: 500;
  color: var(--success-color);
}

.minigames-fame {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 6px;
  font-weight: 500;
  color: var(--warning-color);
}

.minigames-bet-tip {
  font-size: 10px;
  padding: 1px 5px;
  border-radius: 6px;
  background: #fee;
  color: #c00;
  font-weight: 500;
}

/* ==================== 日志记录框 ==================== */
.log-panel {
  height: 120px;
}
</style>
