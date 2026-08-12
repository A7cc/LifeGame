<template>
  <el-dialog v-model="visible" title="掷骰子游戏" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">掷骰子比大小，简单有趣</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 100 }}元</span>
          <span class="info-item">🏆 奖励：最高200元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="dice-display">
          <div class="dice-player">
            <div class="dice-label">你</div>
            <div class="dice-value">{{ playerValue || '?' }}</div>
          </div>
          <div class="dice-vs">VS</div>
          <div class="dice-computer">
            <div class="dice-label">电脑</div>
            <div class="dice-value">{{ computerValue || '?' }}</div>
          </div>
        </div>
        <el-button type="primary" @click="play" style="width: 100%; margin-top: 15px;">掷骰子</el-button>
        <div v-if="result" class="game-result-text">
          <div :class="{ win: result === 'win', lose: result === 'lose' }">{{ resultText }}</div>
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

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'dice', name: '掷骰子', entryCost: 100 })
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
const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)

// 本地游戏状态
const playerValue = ref('')
const computerValue = ref('')
const result = ref('')
const resultText = ref('')
const reward = ref(0)
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')

// 重置本地状态
const resetLocal = () => {
  playerValue.value = ''
  computerValue.value = ''
  result.value = ''
  resultText.value = ''
  reward.value = 0
  resultIcon.value = ''
  resultTitle.value = ''
  resultDetail.value = ''
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

const play = async () => {
	const player = Number(startData.value?.round?.player || 0)
	const computer = Number(startData.value?.round?.computer || 0)
  playerValue.value = player
  computerValue.value = computer

  let winCount = 0
  let customResultText = ''

  if (player > computer) {
    result.value = 'win'
    resultText.value = '你赢了！'
    reward.value = 200
    winCount = 1
    customResultText = `你赢了！获得200元`
    ElMessage.success('你赢了！获得200元！')
  } else if (player < computer) {
    result.value = 'lose'
    resultText.value = '你输了！'
    winCount = 0
    customResultText = '你输了！'
    ElMessage.error('你输了！')
  } else {
    result.value = 'draw'
    resultText.value = '平局！返还报名费'
    winCount = 2 // 后端约定：2 表示平局并返还报名费
    customResultText = '平局！返还报名费'
    ElMessage.info('平局！返还报名费')
  }

  // 结算
	const gameResult = await endGame(999, customResultText, { player, computer, reward: reward.value })

  if (gameResult) {
    reward.value = gameResult.payout
    resultIcon.value = gameResult.win ? '🎉' : '💔'
    resultTitle.value = gameResult.win ? '获胜！' : '失败！'
    resultDetail.value = `${resultText.value}；${gameResult.resultText}`

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

.lose {
  color: var(--error-color);
  font-weight: 600;
}

.dice-display {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 20px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin: 15px 0;
}

.dice-player, .dice-computer {
  text-align: center;
}

.dice-label {
  font-size: 12px;
  color: var(--font-light);
  margin-bottom: 5px;
}

.dice-value {
  font-size: 36px;
  font-weight: 600;
}

.dice-vs {
  font-size: 14px;
  color: var(--font-light);
  font-weight: 600;
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
