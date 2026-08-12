<template>
  <el-dialog v-model="visible" title="猜拳游戏" width="400px" @close="handleClose">
    <div class="game-dialog">
      <div class="game-info-text">选择你的出拳：</div>
      <div class="rps-choices">
        <div class="rps-choice" @click="play('rock')">✊ 石头</div>
        <div class="rps-choice" @click="play('scissors')">✋ 剪刀</div>
        <div class="rps-choice" @click="play('paper')">✌️ 布</div>
      </div>
      <div v-if="result" class="game-result-text">
        <div>你出：{{ playerChoice }}</div>
        <div>电脑出：{{ computerChoice }}</div>
        <div class="result-text" :class="{ win: result === 'win', lose: result === 'lose' }">
          {{ resultText }}
        </div>
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
    default: () => ({ id: 'rps', name: '猜拳游戏', entryCost: 50, needBet: false })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) reset()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const result = ref('')
const playerChoice = ref('')
const computerChoice = ref('')
const resultText = ref('')

const choices = ['rock', 'scissors', 'paper']
const choiceNames = { rock: '石头', scissors: '剪刀', paper: '布' }

const handleClose = () => {
  visible.value = false
  reset()
  result.value = ''
  playerChoice.value = ''
  computerChoice.value = ''
  resultText.value = ''
}

const handleStartGame = async (choice) => {
	const success = await startGame({ choice })
  if (!success) return
  gameStarted.value = true
}

const play = async (choice) => {
	if (!gameStarted.value) {
		await handleStartGame(choice)
	}
	if (!gameStarted.value) return

	const computer = startData.value?.round?.computerChoice
	if (!choices.includes(computer)) return
  playerChoice.value = choiceNames[choice]
  computerChoice.value = choiceNames[computer]

  let winCount = 0
  if (choice === computer) {
    result.value = 'draw'
    resultText.value = '平局！'
    winCount = 0
  } else if (
    (choice === 'rock' && computer === 'scissors') ||
    (choice === 'scissors' && computer === 'paper') ||
    (choice === 'paper' && computer === 'rock')
  ) {
    result.value = 'win'
    resultText.value = '你赢了！'
    winCount = 1
    ElMessage.success('你赢了！')
  } else {
    result.value = 'lose'
    resultText.value = '你输了！'
    winCount = 0
    ElMessage.error('你输了！')
  }

	const gameResult = await endGame(999, `${playerChoice.value} vs ${computerChoice.value}: ${resultText.value}`, {
    playerChoice: choice,
    computerChoice: computer
  })

  emit('complete', gameResult)
}
</script>

<style scoped>
.game-dialog {
  padding: 10px 0;
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

.result-text {
  font-size: 16px;
  font-weight: 600;
  margin-top: 8px;
}

.result-text.win {
  color: var(--success-color);
}

.result-text.lose {
  color: var(--error-color);
}

.rps-choices {
  display: flex;
  gap: 10px;
  justify-content: center;
}

.rps-choice {
  flex: 1;
  padding: 15px;
  background: var(--panel-color);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 18px;
}

.rps-choice:hover {
  border-color: var(--el-color-primary);
  background: #fffbeb;
  transform: scale(1.05);
}
</style>
