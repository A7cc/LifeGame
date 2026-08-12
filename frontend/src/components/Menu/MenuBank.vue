<template>
  <div class="product-panel" data-testid="page-bank">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">银行</div>
    </div>

    <!-- 资金信息 -->
    <div class="bank-main">
      <!-- 左侧：存取款和贷款 -->
      <div class="bankhouse">
        <div class="section-header">
          <span class="section-title">银行业务</span>
          <div class="section-controls">
            <!-- 还款按钮放在标题行 -->
            <el-button data-testid="bank-repay" v-if="userInfo.uloan > 0" size="small" type="danger" @click="handleRepayLoan">还款 {{ (userInfo.uloan + calculateInterest()).toLocaleString() }}元</el-button>
          </div>
        </div>

        <div class="bankhouse-content">
          <!-- 存钱 -->
          <div class="action-item">
            <span>存钱</span>
            <el-input-number data-testid="bank-deposit-amount" v-model="SaveCount" :min="0" :step="100" size="small"/>
            <el-button data-testid="bank-deposit" size="small" type="primary" @click="handleOperationMoney('deposit')">存</el-button>
          </div>
          <!-- 取钱 -->
          <div class="action-item">
            <span>取钱</span>
            <el-input-number data-testid="bank-withdraw-amount" v-model="DrawCount" :min="0" :step="100" size="small"/>
            <el-button data-testid="bank-withdraw" size="small" type="success" @click="handleOperationMoney('withdraw')">取</el-button>
          </div>
          <!-- 贷款 -->
          <div class="action-item">
            <span>贷款</span>
            <el-input-number data-testid="bank-loan-amount" v-model="LoanCount" :min="0" :step="1000" size="small"/>
            <el-button data-testid="bank-loan" size="small" type="warning" @click="handleApplyLoan" :disabled="userInfo.uloan > 0">贷</el-button>
          </div>
        </div>
      </div>

      <!-- 中间：任务信息 -->
      <div class="operationinfo">
        <div class="section-header">
          <span class="section-title">当前任务</span>
        </div>
        <div class="task-list" v-loading="loadingTasks">
          <div v-for="task in taskList" :key="task.taskid" class="task-item" :class="{ completed: task.completed, claimed: isTaskClaimed(task.taskid) }">
            <div class="task-header">
              <span class="task-name">{{ task.taskname }}</span>
              <el-tag v-if="task.completed && !isTaskClaimed(task.taskid)" type="success" size="small" @click="handleClaimReward(task.taskid, task.reward)">领取奖励</el-tag>
              <el-tag v-else :type="getTaskTagType(task)" size="small">{{ getTaskStatusText(task) }}</el-tag>
            </div>
            <div class="task-desc">{{ task.taskdesc }}</div>
            <div class="task-progress">
              <span class="progress-text">进度：{{ formatTaskProgress(task) }}</span>
              <span class="reward-text">奖励：{{ task.reward.toLocaleString() }}元</span>
            </div>

          </div>
          <div v-if="taskList.length === 0 && !loadingTasks" class="no-tasks">暂无可用任务</div>
        </div>
      </div>

      <!-- 右侧：资金情况 -->
      <div class="ubankinfo">
        <div class="section-header">
          <span class="section-title">当前资金情况</span>
        </div>

        <el-descriptions border :column="2" size="small">
          <el-descriptions-item label="现金">
            <span class="money-value" data-testid="bank-cash">{{ formatPrice(userInfo.ucash) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="存款">
            <span class="money-value" data-testid="bank-balance">{{ formatPrice(userInfo.ubank) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="贷款">
            <span class="loan-value" data-testid="bank-loan-balance">{{ formatPrice(userInfo.uloan) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="逾期次数">
            <span class="overdue-value">{{ userInfo.uloanoverdue }} 次</span>
          </el-descriptions-item>
          <el-descriptions-item label="净资产">
            <span class="money-value">{{ formatPrice(userInfo.uassets) }}</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

    <!-- 打工区域 -->
    <div class="bank-bottom">
      <SpecialGameBank @refreshTasks="loadTaskList" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useGameStore } from '@/src/stores/game'
import { ElMessage, ElMessageBox } from 'element-plus'
import { OperationMoney, ApplyLoan, RepayLoan, GetBankTaskList, ClaimTaskReward } from "@/wailsjs/go/services/App.js"
import { formatPrice } from '@/src/utils/format'
import SpecialGameBank from '@/src/components/GameList/WorkGameBank.vue'

// 获取游戏信息
const gameStore = useGameStore()
// 获取用户信息
const userInfo = computed(() => gameStore.userInfo)

// 任务相关
const taskList = ref([])
const loadingTasks = ref(false)

// 计算利息
const calculateInterest = () => {
  return Math.ceil(userInfo.value.uloan * 0.1)
}
// 存钱金额
const SaveCount = ref(100)
// 取钱金额
const DrawCount = ref(100)
// 贷款金额
const LoanCount = ref(100)

// 存钱取钱操作
const handleOperationMoney = async (opemsg) => {
  // 存放数量
  let count = 0
  if (opemsg === 'deposit') count = SaveCount.value
  else if (opemsg === 'withdraw') count = DrawCount.value
  else { ElMessage.error('操作类型错误'); return }

  // 判断金额是否大于0
  if (count <= 0) {
    ElMessage.error('请输入大于0的金额')
    return
  }

  const data = await OperationMoney(opemsg, count)
  if (data.code == 200) {
    // 赋值给用户
    gameStore.applyUserInfo(data.userinfo)
    // 刷新任务列表
    await loadTaskList()
    // 弹窗提示
    ElMessage.success(data.msg)
  } else {
    // 弹窗提示
    ElMessage.error(data.msg || `${opemsg}失败`)
  }
}
// 申请贷款
const handleApplyLoan = async () => {
  if (LoanCount.value <= 0) {
    ElMessage.error('请输入贷款金额')
    return
  }

  try {
    const result = await ApplyLoan(LoanCount.value)
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      // 刷新任务列表
      await loadTaskList()
      ElMessage.success(result.msg)
    } else {
      ElMessage.error(result.msg)
    }
  } catch (err) {
    ElMessage.error('申请贷款失败')
  }
}

// 还款
const handleRepayLoan = async () => {
  if (userInfo.value.uloan <= 0) {
    ElMessage.warning('没有需要偿还的贷款')
    return
  }

  const totalDue = userInfo.value.uloan + calculateInterest()
  const interest = calculateInterest()

  // 弹出确认对话框
  ElMessageBox.confirm(
    `应还总额：${totalDue.toLocaleString()} 元`,
    '偿还贷款确认',
    {
      confirmButtonText: '确认还款',
      cancelButtonText: '取消',
      message: `
<div style="line-height: 1.8;">
  <p style="margin: 8px 0;">💰 当前贷款：<strong>${userInfo.value.uloan.toLocaleString()}</strong> 元</p>
  <p style="margin: 8px 0;">📊 应付利息：<strong>${interest.toLocaleString()}</strong> 元（年利率10%）</p>
  <p style="margin: 8px 0;">💵 应还总额：<strong style="color: var(--warning-color); font-size: 16px;">${totalDue.toLocaleString()}</strong> 元</p>
  <p style="margin: 8px 0;">🏦 当前存款：<strong>${userInfo.value.ubank.toLocaleString()}</strong> 元</p>
</div>
      `,
      dangerouslyUseHTMLString: true,
    }
  ).then(async () => {
    // 检查存款是否足够
    if (userInfo.value.ubank < totalDue) {
      ElMessage.warning(`存款不足！需要 ${totalDue.toLocaleString()} 元`)
      return
    }

    try {
      const result = await RepayLoan(totalDue)
      if (result.code === 200) {
        gameStore.applyUserInfo(result.userinfo)
        // 刷新任务列表
        await loadTaskList()
        ElMessage.success(result.msg)
      } else {
        ElMessage.error(result.msg)
      }
    } catch (err) {
      ElMessage.error('还款失败')
    }
  }).catch(() => {
    // 用户取消
  })
}

// ============== 任务相关方法 ==============
// 加载任务列表
const loadTaskList = async () => {
  loadingTasks.value = true
  try {
    const result = await GetBankTaskList()
    if (result.code === 200) {
      taskList.value = result.tasklist
    } else {
      ElMessage.error('获取任务列表失败')
    }
  } catch (err) {
    console.error('加载任务列表失败:', err)
  } finally {
    loadingTasks.value = false
  }
}

// 格式化任务进度
const formatTaskProgress = (task) => {
  if (task.completed) {
    return '已完成'
  }
  const current = task.current || 0
  const target = task.target
  return `${current.toLocaleString()} / ${target.toLocaleString()}`
}

// 获取任务标签类型
const getTaskTagType = (task) => {
  if (isTaskClaimed(task.taskid)) return 'info'
  if (task.completed) return 'success'
  return 'warning'
}

// 获取任务状态文本
const getTaskStatusText = (task) => {
  if (isTaskClaimed(task.taskid)) return '已领取'
  if (task.completed) return '可领取'
  return '进行中'
}

// 检查任务是否已领取
const isTaskClaimed = (taskId) => {
  const task = taskList.value.find(t => t.taskid === taskId)
  return task?.claimed || false
}

// 领取任务奖励
const handleClaimReward = async (taskId, reward) => {
  try {
    const result = await ClaimTaskReward(taskId)
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      ElMessage.success(result.msg || `成功领取 ${reward.toLocaleString()} 元奖励！`)
      // 刷新任务列表
      await loadTaskList()
    } else {
      ElMessage.error(result.msg)
    }
  } catch (err) {
    console.error('领取奖励失败:', err)
    ElMessage.error('领取奖励失败')
  }
}

// 组件挂载时加载任务列表
onMounted(() => {
  loadTaskList()
})

// 监听用户年龄变化，当新年到来时重新加载任务列表
watch(() => gameStore.userInfo?.uage, (newAge, oldAge) => {
  // 当年龄增加时（新年），重新加载任务列表
  if (newAge && oldAge && newAge > oldAge) {
    loadTaskList()
  }
})
</script>

<style scoped>
/* 中间：资金部分 */
.bank-main {
  display: flex;
  gap: 5px;
  margin-top: 5px;
  height: 135px;
}

/* 资金部分的存钱取钱贷款、资产以及任务 */
.bankhouse, .operationinfo, .ubankinfo {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--panel-color);
  padding: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.bankhouse {
  flex: 1.2;
}
.operationinfo {
  flex: 1.3;
}
.ubankinfo {
  flex: 1.7;
}

/* 资金部分子内容 */
.bankhouse-content {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
/* 资金部分->存取钱贷款动作 */
.action-item {
  font-size: 15px;
  display: flex;
  justify-content: space-between; /* 左右两端对齐 */
  align-items: center;
  width: 100%;  /* 撑满每一行 */
  padding: 0 5px;
  box-sizing: border-box;
    color: var(--font-color);
}

/* 资金部分->用户资金信息颜色 */
.money-value {
  font-weight: 500;
  color: var(--success-color);
}
.loan-value {
  font-weight: 500;
  color: var(--warning-color);
}
.overdue-value {
  font-weight: 500;
  color: var(--error-color);
}

/* 银行底部（打工区域） */
.bank-bottom {
  height: calc(100% - 200px);
  display: flex;
  flex: 1;
  gap: 5px;
  margin-top: 5px;
}

/* ============== 任务列表样式 ============== */
.task-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding-right: 4px;
  /* 隐藏滚动条 */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE/Edge */
}

.task-list::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.task-item {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 5px 8px;
  background: var(--panel-color);
  transition: all 0.3s ease;
}

.task-item:hover {
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.task-item.completed {
  border-color: var(--success-color);
  background: var(--panel-color);
}

.task-item.claimed {
  opacity: 0.6;
  border-color: var(--font-light);
  background: var(--panel-log-color);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-name {
  font-weight: 500;
  font-size: 13px;
  color: var(--font-color);
}

.task-desc {
  font-size: 12px;
  color: var(--font-secondary);
  line-height: 1.4;
}

.task-progress {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.progress-text {
  color: var(--info-color);
}

.reward-text {
  color: var(--warning-color);
  font-weight: 500;
}

.no-tasks {
  text-align: center;
  color: var(--font-light);
  padding: 30px 0;
  font-size: 14px;
}
</style>
