<template>
  <el-dialog v-model="visible" title="猜数字游戏" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">猜一个1-100之间的数字</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ config?.entryCost || 10 }}元</span>
          <span class="info-item">🏆 奖励：最高500元</span>
        </div>
        <el-button type="primary" @click="handleStartGame" :disabled="processing" style="width: 100%">开始游戏</el-button>
      </div>
      <div v-else-if="!gameEnded" class="game-playing">
        <div class="game-info-text">猜一个1-100之间的数字</div>
        <el-input-number v-model="guessInput" :min="1" :max="100" placeholder="输入数字" style="width: 100%; margin-bottom: 10px;" />
        <el-button type="primary" @click="play" style="width: 100%;">猜！</el-button>
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

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'guess', name: '猜数字', entryCost: 10 })
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

// 本地游戏状态
const guessInput = ref(50)
const result = ref('')
const resultText = ref('')
const reward = ref(0)
const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')

// 重置本地状态
const resetLocal = () => {
  guessInput.value = 50
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
	const guess = guessInput.value
	const gameResult = await endGame(guess, null, { guess })
	if (!gameResult) return

	const answer = Number(gameResult.round?.answer || 0)
	if (gameResult.outcome === 1) {
		result.value = 'win'
		resultText.value = '猜对了！'
		ElMessage.success('恭喜！获得500元！')
	} else if (gameResult.outcome === 2) {
		result.value = 'close'
		resultText.value = `很接近！正确答案是${answer}`
		ElMessage.success('差一点点！返还50元')
	} else {
		result.value = 'lose'
		resultText.value = `猜错了！正确答案是${answer}`
		ElMessage.error(`猜错了！正确答案是${answer}`)
	}

	reward.value = gameResult.payout
	resultIcon.value = gameResult.win ? '🎉' : '💔'
	resultTitle.value = gameResult.win ? '获胜！' : '失败！'
	resultDetail.value = `${resultText.value}；${gameResult.resultText}`
	emit('complete', gameResult)
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
