<template>
  <el-dialog v-model="visible" title="老虎机" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">转动老虎机，赢取大奖！</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 200 }}元</span>
          <span class="info-item">🏆 奖励：最高1000元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="slot-machine">
          <div class="slot-reel">{{ reels[0] }}</div>
          <div class="slot-reel">{{ reels[1] }}</div>
          <div class="slot-reel">{{ reels[2] }}</div>
        </div>
        <el-button type="primary" @click="play" style="width: 100%; margin-top: 15px;" :disabled="spinning">
          {{ spinning ? '转动中...' : '开始' }}
        </el-button>
        <div v-if="result" class="game-result-text">
          <div :class="{ win: result === 'win' }">{{ resultText }}</div>
          <div v-if="reward > 0">获得 {{ reward }} 元！</div>
        </div>
      </div>
      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'
import { useCleanupTasks } from '@/src/composables/useCleanupTasks'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'slot', name: '老虎机', entryCost: 200 })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)
const { cleanup, clearManagedTimer, setManagedInterval, setManagedTimeout } = useCleanupTasks()

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetLocal()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

// 使用小游戏基础逻辑
const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)

// 本地游戏状态
const reels = ref(['🎰', '🎰', '🎰'])
const spinning = ref(false)
const result = ref('')
const resultText = ref('')
const reward = ref(0)
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')

const symbols = ['🍒', '🍋', '🍊', '🍇', '⭐', '💎', '7️⃣']
let spinInterval = null

// 重置本地状态
const resetLocal = () => {
  reels.value = ['🎰', '🎰', '🎰']
  spinning.value = false
  result.value = ''
  resultText.value = ''
  reward.value = 0
  resultIcon.value = ''
  resultTitle.value = ''
  resultDetail.value = ''
  if (spinInterval) {
    clearManagedTimer(spinInterval)
    spinInterval = null
  }
  reset()
}

const handleClose = () => {
  visible.value = false
  resetLocal()
}

// 开始游戏
const handleStartGame = async () => {
  await startGame()
}

const play = () => {
  spinning.value = true
  result.value = ''

  // 转动动画
  spinInterval = setManagedInterval(() => {
    reels.value = reels.value.map(() =>
      symbols[Math.floor(Math.random() * symbols.length)]
    )
  }, 100)

  // 2秒后停止
  setManagedTimeout(async () => {
    clearManagedTimer(spinInterval)
    spinInterval = null

		const results = startData.value?.round?.reels || ['🍒', '🍋', '🍊']
    reels.value = results
    spinning.value = false

    let winCount = 0
    let customResultText = ''

    // 判断奖励
    if (results[0] === results[1] && results[1] === results[2]) {
      // 三个相同
      const rewardMap = { '💎': 1000, '7️⃣': 500, '⭐': 300, '🍒': 200, '🍋': 150, '🍊': 150, '🍇': 150 }
      const outcomeMap = { '💎': 6, '7️⃣': 5, '⭐': 4, '🍒': 3, '🍋': 2, '🍊': 2, '🍇': 2 }
      const r = rewardMap[results[0]] || 100
      result.value = 'win'
      resultText.value = `大奖！三个${results[0]}！`
      reward.value = r
      winCount = outcomeMap[results[0]] || 2
      customResultText = `大奖！三个${results[0]}！获得${r}元`
      ElMessage.success(`大奖！获得${r}元！`)
    } else if (results[0] === results[1] || results[1] === results[2] || results[0] === results[2]) {
      // 两个相同
      const r = 50
      result.value = 'small'
      resultText.value = '小奖！两个相同！'
      reward.value = r
      winCount = 1
      customResultText = `小奖！两个相同！获得${r}元`
      ElMessage.success(`小奖！获得${r}元！`)
    } else {
      result.value = 'lose'
      resultText.value = '再接再厉！'
      winCount = 0
      customResultText = '没有中奖...'
      ElMessage.error('没有中奖...')
    }

    // 结算
		const gameResult = await endGame(999, customResultText, { reels: results.join(' '), reward: reward.value })

    if (gameResult) {
      reward.value = gameResult.payout
      resultIcon.value = gameResult.win ? '🎉' : '💔'
      resultTitle.value = gameResult.win ? '中奖！' : '未中奖'
      resultDetail.value = `${resultText.value}；${gameResult.resultText}`

      emit('complete', gameResult)
    }
  }, 2000)
}

defineExpose({
  cleanup,
})
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

.game-info-text {
  font-size: 14px;
  color: var(--font-secondary);
  text-align: center;
  margin-bottom: 15px;
}

.game-result-text {
  margin-top: 15px;
  padding: 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  text-align: center;
}

.win {
  color: var(--success-color);
  font-weight: 600;
}

.slot-machine {
  display: flex;
  justify-content: center;
  gap: 15px;
  padding: 20px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin: 15px 0;
}

.slot-reel {
  width: 60px;
  height: 80px;
  background: #fff;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
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
  margin-bottom: 8px;
}

.result-detail {
  font-size: 14px;
  color: var(--font-secondary);
}
</style>
