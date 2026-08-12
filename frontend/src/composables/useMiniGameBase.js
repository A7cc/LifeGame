import { onBeforeUnmount, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { AddMiniGameWager, CancelMiniGame, MiniGameAction, StartMiniGameWithOptions, EndMiniGame } from '@/wailsjs/go/services/App.js'

/**
 * 小游戏结果结构
 * @typedef {Object} GameResult
 * @property {string} gameId - 游戏ID
 * @property {string} gameName - 游戏名称
 * @property {boolean} win - 是否获胜
 * @property {number} score - 分数/胜利次数
 * @property {string} resultText - 结果描述
 * @property {number} [cashChange] - 现金变化
 * @property {number} [fameChange] - 名声变化
 * @property {number} [payout] - 后端实际返还金额
 * @property {number} [netChange] - 扣除报名费和下注后的本局净变化
 * @property {any} [detail] - 游戏特定数据
 */

/**
 * 小游戏配置结构
 * @typedef {Object} MiniGameConfig
 * @property {string} id - 游戏ID
 * @property {string} name - 游戏名称
 * @property {number} entryCost - 报名费
 * @property {boolean} [needBet] - 是否需要额外下注
 */

/**
 * 小游戏基础逻辑 composable
 * 封装 StartMiniGame 和 EndMiniGame API 调用
 *
 * @param {MiniGameConfig} config - 游戏配置
 * @returns {Object} 游戏状态和方法
 */
export function useMiniGameBase(config) {
  const gameStore = useGameStore()

  // 游戏状态
  const gameStarted = ref(false)
  const gameEnded = ref(false)
  const processing = ref(false)
  const startData = ref(null)

  // 资金快照（用于计算净变化）
  const cashSnapshot = ref(0)
  const fameSnapshot = ref(0)

  /**
   * 开始游戏
   * 调用后端 StartMiniGame API 扣除报名费
   * @returns {Promise<boolean>} 是否成功开始
   */
  const startGame = async ({ wager = 0, choice = '', quantity = 1 } = {}) => {
    // 这里只做快速提示，最终金额和参数合法性由后端校验。
    const entryCost = config?.entryCost || 0
    const estimatedCost = config?.id === 'lottery'
      ? entryCost * quantity
      : entryCost + wager
    if (gameStore.userInfo.ucash < estimatedCost) {
      ElMessage.error('资金不足！')
      return false
    }

    processing.value = true
    cashSnapshot.value = Number(gameStore.userInfo?.ucash || 0)
    fameSnapshot.value = Number(gameStore.userInfo?.ufame || 0)

    try {
      const result = await StartMiniGameWithOptions(
        config.id,
        0,
        Number(wager) || 0,
        String(choice ?? ''),
        Number(quantity) || 1
      )

      if (result?.code !== 200) {
        ElMessage.error(result?.msg || '游戏开始失败')
        processing.value = false
        return false
      }

      // 更新用户信息（后端已扣除报名费）
      if (result?.userinfo) {
        gameStore.applyUserInfo(result.userinfo)
      }
      startData.value = result

      gameStarted.value = true
      gameEnded.value = false
      processing.value = false

      return true
    } catch (error) {
      console.error('游戏开始失败', error)
      ElMessage.error('游戏开始失败')
      processing.value = false
      return false
    }
  }

  /**
   * 结束游戏
   * 调用后端 EndMiniGame API 发放奖励
   * @param {number} winCount - 胜利次数/分数（1=获胜，0=失败）
   * @param {string} customResultText - 自定义结果文本
   * @param {any} detail - 游戏特定数据
   * @returns {Promise<GameResult|null>} 游戏结果
   */
  const endGame = async (winCount, customResultText = null, detail = null) => {
    if (!gameStarted.value) return null

    gameEnded.value = true
    processing.value = true

    try {
      const result = await EndMiniGame(Math.round(winCount))

      // 更新用户信息
      if (result?.userinfo) {
        gameStore.applyUserInfo(result.userinfo)
      }

      // 计算变化
      const currentCash = Number(gameStore.userInfo?.ucash || 0)
      const currentFame = Number(gameStore.userInfo?.ufame || 0)
      const cashChange = currentCash - cashSnapshot.value
      const fameChange = currentFame - fameSnapshot.value

      const isWin = typeof result?.win === 'boolean'
        ? result.win
        : result?.code === 200
      const resultText = isWin
        ? (result?.msg || customResultText || '获胜')
        : (result?.msg || '失败')

      /** @type {GameResult} */
      const gameResult = {
        gameId: config.id,
        gameName: config.name,
        win: isWin,
        score: Number(result?.outcome ?? winCount),
        outcome: Number(result?.outcome ?? winCount),
        resultText,
        cashChange,
        fameChange,
        payout: Number(result?.payout || 0),
        netChange: Number(result?.netchange ?? cashChange),
        round: result?.round || null,
        detail: result?.detail || detail
      }

      processing.value = false
      return gameResult
    } catch (error) {
      console.error('游戏结算失败', error)
      ElMessage.error('游戏结算失败')
      processing.value = false
      return null
    }
  }

  const addWager = async (amount) => {
    try {
      const result = await AddMiniGameWager(Number(amount) || 0)
      if (result?.code !== 200) {
        ElMessage.error(result?.msg || '追加下注失败')
        return null
      }
      if (result?.userinfo) {
        gameStore.applyUserInfo(result.userinfo)
      }
      return result
    } catch (error) {
      console.error('追加下注失败', error)
      ElMessage.error('追加下注失败')
      return null
    }
  }

  const action = async (name) => {
    try {
      const result = await MiniGameAction(String(name || ''))
      if (result?.code !== 200) {
        ElMessage.error(result?.msg || '游戏操作失败')
        return null
      }
      return result
    } catch (error) {
      console.error('游戏操作失败', error)
      ElMessage.error('游戏操作失败')
      return null
    }
  }

  /**
   * 放弃当前尚未结算的后端会话。会话标识由后端生成，延迟到达的旧请求
   * 不会误取消随后开始的新一局。
   */
  const cancelGame = async (sessionId = startData.value?.sessionid) => {
    if (!sessionId) return false

    try {
      const result = await CancelMiniGame(String(sessionId))
      return Boolean(result?.cancelled)
    } catch (error) {
      console.error('取消小游戏失败', error)
      return false
    }
  }

  /**
   * 重置游戏状态
   */
  const reset = () => {
    const sessionId = startData.value?.sessionid
    gameStarted.value = false
    gameEnded.value = false
    processing.value = false
    cashSnapshot.value = 0
    fameSnapshot.value = 0
    startData.value = null
    if (sessionId) void cancelGame(sessionId)
  }

  onBeforeUnmount(() => {
    const sessionId = startData.value?.sessionid
    if (sessionId) void cancelGame(sessionId)
  })

  return {
    gameStarted,
    gameEnded,
    processing,
    startData,
    startGame,
    addWager,
    action,
    endGame,
    cancelGame,
    reset,
    cashSnapshot,
    fameSnapshot
  }
}
