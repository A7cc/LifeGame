<template>
  <el-dialog v-model="visible" title="彩票刮刮乐" width="500px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">试试你的运气，也许能中大奖！</div>
        <div class="warning-box">
          <span class="warning-icon">⚠️</span>
          <span class="warning-text">彩票返奖率仅50%，中大奖概率极低！</span>
        </div>
        <div class="game-info">
          <span class="info-item">💰 单张价格：{{ ticketPrice }}元</span>
          <span class="info-item">🏆 头奖：10000元</span>
          <span class="info-item">📊 中奖率：35%</span>
        </div>

        <!-- 购买数量 -->
        <div class="purchase-section">
          <div class="purchase-label">购买数量：{{ ticketCount }} 张</div>
          <div class="ticket-buttons">
            <el-button @click="ticketCount = 1" :type="ticketCount === 1 ? 'primary' : ''" size="small">1张</el-button>
            <el-button @click="ticketCount = 5" :type="ticketCount === 5 ? 'primary' : ''" size="small">5张</el-button>
            <el-button @click="ticketCount = 10" :type="ticketCount === 10 ? 'primary' : ''" size="small">10张</el-button>
            <el-button @click="ticketCount = 20" :type="ticketCount === 20 ? 'primary' : ''" size="small">20张</el-button>
          </div>
          <div class="total-cost">总花费：💰 {{ totalCost }} 元</div>
        </div>

        <el-button type="primary" @click="handleStartGame" style="width: 100%">购买并刮奖</el-button>

        <!-- 统计信息 -->
        <div v-if="totalSpent > 0" class="stats-section">
          <div class="stats-title">本次会话统计</div>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-label">购买张数</div>
              <div class="stat-value">{{ totalTickets }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">总花费</div>
              <div class="stat-value loss">💰 {{ totalSpent }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">总收获</div>
              <div class="stat-value" :class="{ win: totalWon > 0 }">💰 {{ totalWon }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">净盈亏</div>
              <div class="stat-value" :class="{ win: netProfit >= 0, loss: netProfit < 0 }">
                {{ netProfit >= 0 ? '+' : '' }}{{ netProfit }}
              </div>
            </div>
          </div>
          <div v-if="netProfit < 0" class="loss-warning">
            💔 已亏损 {{ Math.abs(netProfit) }} 元 - 这就是彩票的真相！
          </div>
        </div>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <div class="scratch-area">
          <div class="tickets-grid">
            <div
              v-for="(ticket, index) in tickets"
              :key="index"
              class="lottery-ticket"
              :class="{ scratched: ticket.scratched }"
              @click="scratchTicket(index)"
            >
              <div class="ticket-cover" v-if="!ticket.scratched">
                <span class="scratch-text">点击刮奖</span>
                <span class="scratch-icon">🎫</span>
              </div>
              <div class="ticket-content" v-else>
                <div class="ticket-prize" :class="getPrizeClass(ticket.prize)">
                  {{ ticket.prize > 0 ? `💰 ${ticket.prize}` : '😢' }}
                </div>
                <div class="ticket-label">
                  {{ ticket.prize > 0 ? '中奖！' : '未中奖' }}
                </div>
              </div>
            </div>
          </div>
          <div class="scratch-progress">
            已刮 {{ scratchedCount }} / {{ ticketCount }} 张
          </div>
        </div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-detail">{{ resultDetail }}</div>

        <!-- 本轮结果 -->
        <div class="round-result">
          <div class="round-title">本轮结果</div>
          <div class="prize-breakdown">
            <div v-for="(count, prize) in prizeBreakdown" :key="prize" class="prize-item">
              <span class="prize-label">{{ prize === '0' ? '未中奖' : prize + '元' }}：</span>
              <span class="prize-count">{{ count }} 张</span>
            </div>
          </div>
          <div class="round-summary">
            <div class="summary-item">花费：💰 {{ totalCost }}</div>
            <div class="summary-item" :class="{ win: roundWinnings > 0 }">
              收入：💰 {{ roundWinnings }}
            </div>
            <div class="summary-item result-net" :class="{ win: roundNet >= 0, loss: roundNet < 0 }">
              净盈亏：{{ roundNet >= 0 ? '+' : '' }}{{ roundNet }}
            </div>
          </div>
        </div>

        <div class="result-message">{{ resultMessage }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'lottery', name: '彩票刮刮乐', entryCost: 100, needBet: true })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const { gameStarted, gameEnded, processing, startData, startGame, endGame, reset } = useMiniGameBase(props.config)
const visible = ref(props.modelValue)

const ticketPrice = computed(() => props.config?.entryCost || 50)
const totalCost = computed(() => ticketCount.value * ticketPrice.value)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetAll()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const JACKPOT = 10000

const ticketCount = ref(1)
const tickets = ref([])
const scratchedCount = ref(0)

const totalTickets = ref(0)
const totalSpent = ref(0)
const totalWon = ref(0)

const resultIcon = ref('')
const resultTitle = ref('')
const resultDetail = ref('')
const resultMessage = ref('')

const netProfit = computed(() => totalWon.value - totalSpent.value)
const roundWinnings = ref(0)
const roundNet = ref(0)
const prizeBreakdown = ref({})

const resetAll = () => {
  resetGameState()
  totalTickets.value = 0
  totalSpent.value = 0
  totalWon.value = 0
}

const resetGameState = () => {
  reset()
  tickets.value = []
  scratchedCount.value = 0
  prizeBreakdown.value = {}
}

const handleClose = () => {
  visible.value = false
  resetAll()
}

const handleStartGame = async () => {
  const success = await startGame({ quantity: ticketCount.value })
  if (!success) return

  totalTickets.value += ticketCount.value
  totalSpent.value += totalCost.value

  scratchedCount.value = 0
  prizeBreakdown.value = {}

  // 奖面由后端生成，前端只负责刮奖动画和展示。
  tickets.value = (startData.value?.round?.tickets || []).map(prize => ({
    prize,
    scratched: false
  }))
}

const scratchTicket = (index) => {
  if (tickets.value[index].scratched) return

  tickets.value[index].scratched = true
  scratchedCount.value++

  // 统计
  const prize = tickets.value[index].prize
  prizeBreakdown.value[prize] = (prizeBreakdown.value[prize] || 0) + 1

  // 检查是否全部刮完
  if (scratchedCount.value >= ticketCount.value) {
    setTimeout(() => showResults(), 500)
  }
}

const showResults = async () => {
  // 计算本轮结果
  let winnings = 0
  tickets.value.forEach(ticket => {
    if (ticket.prize > 0) {
      winnings += ticket.prize
    }
  })

  roundWinnings.value = winnings
  roundNet.value = winnings - totalCost.value
  totalWon.value += winnings

  // 计算胜负
  let winCount = winnings > 0 ? 1 : 0
  let customResultText = ''
  let detail = { tickets: ticketCount.value, cost: totalCost.value, winnings: winnings, prizeBreakdown: prizeBreakdown.value }

  // 设置结果
  if (winnings > totalCost.value) {
    resultIcon.value = '🎉'
    resultTitle.value = '运气不错！'
    resultDetail.value = `本轮赢了一些`
    resultMessage.value = '这次运气好，但不要以为每次都能赢！'
    customResultText = `赢了${winnings - totalCost.value}元`
  } else if (winnings > 0) {
    resultIcon.value = '😊'
    resultTitle.value = '小赚一点'
    resultDetail.value = `回了一点血`
    resultMessage.value = '虽有小赢，但长期买彩票只会让你越买越穷！'
    customResultText = `输了${totalCost.value - winnings}元`
  } else {
    resultIcon.value = '😢'
    resultTitle.value = '全部未中奖'
    resultDetail.value = '运气不佳'
    resultMessage.value = '这就是彩票的真相 - 庄家永远赚你的钱！'
    customResultText = `输了${totalCost.value}元`
    winCount = 0
  }

  // 检查是否中了大奖
  const hasJackpot = tickets.value.some(t => t.prize >= JACKPOT)
  if (hasJackpot) {
    resultIcon.value = '🏆'
    resultTitle.value = '中大奖了！'
    resultMessage.value = '极小概率事件！别以为下次还能中！'
    winCount = 1
    customResultText = `中了大奖${winnings}元`
  }

  const gameResult = await endGame(winCount, customResultText, detail)
  emit('complete', gameResult)
}

const getPrizeClass = (prize) => {
  if (prize >= 1000) return 'jackpot'
  if (prize >= 500) return 'big-prize'
  if (prize >= 50) return 'medium-prize'
  if (prize > 0) return 'small-prize'
  return 'no-prize'
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

.warning-box {
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.warning-icon {
  font-size: 18px;
}

.warning-text {
  font-size: 13px;
  color: #c00;
  font-weight: 600;
}

.game-info {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.info-item {
  font-size: 12px;
  padding: 6px 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.purchase-section {
  background: var(--panel-color);
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.purchase-label {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--font-color);
}

.ticket-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 12px;
}

.total-cost {
  font-size: 16px;
  font-weight: 600;
  color: var(--warning-color);
}

.stats-section {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  margin-top: 20px;
}

.stats-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--font-color);
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 10px;
}

.stat-item {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  padding: 10px;
  border-radius: 6px;
}

.stat-label {
  font-size: 11px;
  color: var(--font-light);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
}

.stat-value.win {
  color: var(--success-color);
}

.stat-value.loss {
  color: var(--error-color);
}

.loss-warning {
  background: #fee;
  color: #c00;
  padding: 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}

.game-playing {
  text-align: center;
}

.scratch-area {
  padding: 10px 0;
}

.tickets-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin-bottom: 15px;
}

.lottery-ticket {
  aspect-ratio: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  overflow: hidden;
  position: relative;
}

.lottery-ticket:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.ticket-cover {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--font-color);
}

.scratch-text {
  font-size: 11px;
  margin-bottom: 4px;
}

.scratch-icon {
  font-size: 20px;
}

.ticket-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--panel-color);
}

.ticket-prize {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.ticket-prize.jackpot {
  font-size: 14px;
  color: #ffd700;
}

.ticket-prize.big-prize {
  color: var(--warning-color);
}

.ticket-prize.medium-prize {
  color: var(--success-color);
}

.ticket-prize.small-prize {
  color: #409eff;
}

.ticket-prize.no-prize {
  color: var(--font-light);
}

.ticket-label {
  font-size: 10px;
  color: var(--font-light);
}

.scratch-progress {
  font-size: 14px;
  font-weight: 600;
  color: var(--font-color);
}

.game-result {
  text-align: center;
  padding: 20px 0;
}

.result-icon {
  font-size: 56px;
  margin-bottom: 10px;
}

.result-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.result-detail {
  font-size: 13px;
  color: var(--font-light);
  margin-bottom: 15px;
}

.round-result {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
}

.round-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--font-color);
}

.prize-breakdown {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
  margin-bottom: 10px;
}

.prize-item {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
}

.prize-label {
  color: var(--font-secondary);
}

.prize-count {
  font-weight: 600;
  color: #409eff;
}

.round-summary {
  border-top: 1px solid #e0e0e0;
  padding-top: 10px;
}

.summary-item {
  font-size: 13px;
  padding: 4px 0;
  color: var(--font-secondary);
}

.result-net {
  font-size: 16px;
  font-weight: 600;
  margin-top: 4px;
  padding-top: 8px;
  border-top: 1px dashed var(--border-color);
  color: var(--error-color);
}

.result-net.win {
  color: var(--success-color);
}

.result-message {
  font-size: 13px;
  color: var(--font-light);
  line-height: 1.5;
  padding: 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}
</style>
