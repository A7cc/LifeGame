<template>
  <div class="bank-work">
    <!-- 头部 -->
    <div class="section-header">
      <span class="section-title">银行柜台打工</span>
      <div class="section-controls">
        <span v-if="isRunning" class="work-score">完成度：<span :class="getScoreClass()">{{ ((score / scoreTarget) * 100).toFixed(0) }}</span>%</span>
        <span v-if="isRunning" class="work-score">机会：<span style="color: var(--error-color); font-weight: bold;">{{ failCount }}</span>次</span>
        <el-button-group size="small">
          <el-button :type="worklevel === 0 ? 'primary' : ''" :disabled="isRunning" @click="changeLevel(0)">简单</el-button>
          <el-button :type="worklevel === 1 ? 'primary' : ''" :disabled="isRunning" @click="changeLevel(1)">中等</el-button>
          <el-button :type="worklevel === 2 ? 'primary' : ''" :disabled="isRunning" @click="changeLevel(2)">困难</el-button>
        </el-button-group>
        <el-button size="small" type="danger" :disabled="!isRunning" @click="endMiniGame">结束</el-button>
      </div>
    </div>

    <div class="bank-work-panel" v-if="isRunning">
      <!-- 打工内容部分 -->
      <div class="work-panel-body">
        <template v-for="(customer, index) in customers" :key="customer.id">
          <!-- 单个打工线选项 -->
          <div class="work-panel-row" :class="{ active: selectedCustomer === index }" @click="selectCustomer(index)">
            <div class="work-row-name">{{ customer.name }}</div>
            <div class="work-row-time">⏱ {{ customer.timeLeft }}s</div>
            <div class="work-row-other" :class="{ active: i < customer.correctSteps.length }" v-for="(step, i) in customer.steps" :key="i">
              {{ step }}
            </div>
          </div>
        </template>
      </div>
      <!-- 打工游戏按钮区域 -->
      <div class="work-panel-button">
        <el-button
          v-for="but in businessbuttonList2"
          :key="but"
          size="small"
          :type="getButtonType(but)"
          @click="handleClick(but)"
        >
          {{ but }}
        </el-button>
      </div>
    </div>

    <div class="bank-work-panel" v-else>
      <!-- 打工开始页面 -->
      <div class="work-start">
        <el-button type="success" size="large" @click="startWork">开始处理银行业务</el-button>
        <div class="work-tips">
          <p>💡 提示：按正确的顺序处理客户业务</p>
          <p>⏱ 越快完成，奖励越高</p>
          <p>💼 目前打工进度：{{ gameStore.userInfo.uopportunity.ownum }}/{{ gameStore.gameInfo.gmaxholdnum.mwroundnum }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount, watch } from 'vue'
import { useGameStore } from '@/src/stores/game'
import { useCleanupTasks } from '@/src/composables/useCleanupTasks'
import { ElMessage } from 'element-plus'
import { CancelMiniGame, StartMiniGame, EndMiniGame } from "@/wailsjs/go/services/App.js"

// 获取游戏信息
const gameStore = useGameStore()
const { cleanup, clearManagedTimer, setManagedInterval, setManagedTimeout } = useCleanupTasks()

// 打工状态
const isRunning = ref(false)
const activeSessionID = ref('')
// 工作难度等级
const worklevel = ref(0)
// 客户信息
const customers = ref([])
// 存放时间
const timers = []
// 选择客户
const selectedCustomer = ref(0)
// 存放分数
const score = ref(0)
// 目标分数
const scoreTarget = ref(30)
// 失败次数
const failCount = ref(5)
// 业务按钮
const businessbuttonList = ['挂号', '填资料', '审资质', '确认金额', '身份验证', '选卡类型', '签署协议', '存钱', '取钱', '开户', '放款']
const businessbuttonList2 = ref(businessbuttonList)
// 银行业务
const businessList = {
  '💰 存钱': ['挂号', '填资料', '确认金额', '存钱'],
  '💸 取钱': ['挂号', '填资料', '身份验证', '取钱'],
  '🤵 开户': ['挂号', '审资质', '选卡类型', '开户'],
  '🏦 贷款': ['挂号', '审资质', '签署协议', '放款'],
}

// 难度设置
function difficulty(level) {
  if (level === 2) {
    scoreTarget.value = 20
    failCount.value = 10
  } else if (level === 1) {
    scoreTarget.value = 40
    failCount.value = 8
  } else {
    scoreTarget.value = 50
    failCount.value = 5
  }
}

// 洗牌算法，随机打乱顺序，增加难度
function shuffle(arr) {
  return arr
    .map(item => ({ item, sort: Math.random() }))
    .sort((a, b) => a.sort - b.sort)
    .map(({ item }) => item)
}

// 获取按钮类型
function getButtonType(but) {
  const len = but.length
  if (len === 2) return 'primary'
  if (len === 3) return 'success'
  if (len === 4) return 'danger'
  return 'warning'
}

// 获取分数颜色
function getScoreClass() {
  const percent = (score.value / scoreTarget.value) * 100
  if (percent >= 80) return 'work-score-excellent'
  if (percent >= 50) return 'work-score-good'
  return 'work-score-normal'
}

// 打工难度等级
function changeLevel(level) {
  // 设置当前等级
  worklevel.value = level
}

// 生成客户
function generateCustomer() {
  // 随机生成客户
  const keys = Object.keys(businessList)
  // 存放生成的客户
  const randomKey = keys[Math.floor(Math.random() * keys.length)]
  const time = worklevel.value === 0 ? 30 : 25
  return {
    // 客户要处理的业务
    name: randomKey,
    // 步骤
    steps: businessList[randomKey],
    // 倒计时
    timeLeft: time,
    // 存放步骤结果
    correctSteps: [],
    // 客户业务是否在进行中
    finished: false,
    // 客户业务是被玩家真正的处理完成是否成功
    success: false,
  }
}
// 开始打工
const startWork = async () => {
  const data = await StartMiniGame('bank', worklevel.value)
  if (data.code != 200) {
    // 弹窗提示
    ElMessage.error('打工失败！'+data.msg)
    return
  }
  activeSessionID.value = String(data.sessionid || '')
  // 更新用户信息
  gameStore.applyUserInfo(data.userinfo)
  ElMessage.success('开始打工！')
  // 设置打工按钮不可见状态
  isRunning.value = true
  // 分数
  score.value = 0
  // 设置难度
  difficulty(worklevel.value)
  // 当等级是3时，打乱按钮顺序
  if (worklevel.value >= 3) {
    businessbuttonList2.value = shuffle(businessbuttonList)
  }
  // 根据难度生成对应的客户
  customers.value = Array.from({ length: 3 + worklevel.value }, () => generateCustomer())
  // 清除之前的定时器
  clearAllTimers()
  // 循环为多个客户设置倒计时定时器
  customers.value.forEach((customer, i) => startTimer(customer, i))
}

// 选择处理的客户
function selectCustomer(index) {
  selectedCustomer.value = index
}

// 点击按钮
function handleClick(step) {
  const customerIndex = selectedCustomer.value
  // 获取当前客户
  const customer = customers.value[customerIndex]
  // 如果客户已经完成（不管是成功还是失败），就不再处理
  if (customer.finished) return

  const expected = customer.steps[customer.correctSteps.length]
  // 判断当前步骤是否相等
  if (step === expected) {
    // 如果正确，就添加到正确步骤数组中
    customer.correctSteps.push(step)
    // 如果步骤数组的长度大于等于步骤的长度，说明客户已经处理完
    if (customer.correctSteps.length >= customer.steps.length) {
      // 客户处理完成
      customer.finished = true
      customer.success = true
      // 增加分数
      if (customer.timeLeft) {
        if (worklevel.value === 0 && customer.timeLeft >= 20) {
          score.value = score.value + 2
        } else if ((worklevel.value === 1 || worklevel.value === 2) && customer.timeLeft >= 10) {
          score.value = score.value + 2
        }
      }
      score.value++
      // 清除当前的定时器
      clearManagedTimer(timers[customerIndex])
      // 失败次数过多
      if (failCount.value <= 0) {
        endMiniGame()
        return
      }
      // 判断是否达成目标分数
      if (score.value >= scoreTarget.value) {
        endMiniGame()
        return
      }
      // 停止0.5秒，然后用新的客户替换掉当前的
      setManagedTimeout(() => {
        // 用新的客户替换掉当前的
        customers.value[customerIndex] = generateCustomer()
        // 为新客户设置倒计时定时器
        startTimer(customers.value[customerIndex], customerIndex)
      }, 500)
    }
  }
}

// 开始计时
function startTimer(customer, index) {
  // 为每个客户设置一个 每秒执行一次的定时器（保存到 timers 数组中）,每个定时器会倒计时并检查客户是否已处理完
  timers[index] = setManagedInterval(() => {
    // 如果这个客户已经完成（不管是成功还是失败），就不继续倒计时，直接停止定时器
    if (customer.finished) {
      clearManagedTimer(timers[index])
      return
    }
    // 每秒让 timeLeft 减 1，相当于倒计时
    customer.timeLeft--
    // 如果客户时间用完，还没处理，就判定为失败
    if (customer.timeLeft <= 0) {
      customer.finished = true
      customer.success = false
      // 清除当前的定时器
      clearManagedTimer(timers[index])
      // 减少失败次数
      failCount.value--
      // ✅ 失败次数过多
      if (failCount.value <= 0) {
        endMiniGame()
        return
      }
      // 用新的客户替换掉当前的
      customers.value[index] = generateCustomer()
      // 为新客户设置倒计时定时器
      startTimer(customers.value[index], index)
    }
  }, 1000)
}

// 结束打工
const endMiniGame = async () => {
  // 设置打工按钮可见状态
  isRunning.value = false
  // 清除之前的定时器
  clearAllTimers()
  // 清空客户
  customers.value = []
  // 判断是否达到失败次数
  if (failCount.value <= 0) {
    void cancelActiveSession()
    ElMessage.error('打工结束！你出现的错误太多被辞退了')
    return
  }
  // 判断是否达到目标分数
  if (score.value < scoreTarget.value) {
    void cancelActiveSession()
    ElMessage.error('打工失败！未达到目标分数')
    return
  }
  // 获取工资
  const data = await EndMiniGame(score.value)
  activeSessionID.value = ''
  // 清空得分
  score.value = 0
  if (data.code != 200) {
    ElMessage.error(`打工结束！${data.msg || '老板跑路了！'}`)
    return
  }
  // 更新用户信息
  gameStore.applyUserInfo(data.userinfo)
  // 刷新任务列表（通过emit通知父组件）
  emit('refreshTasks')
  // 弹窗提示
  ElMessage.success(`打工结束！${data.msg}`)
}

const cancelActiveSession = async () => {
  const sessionID = activeSessionID.value
  activeSessionID.value = ''
  if (!sessionID) return
  try {
    await CancelMiniGame(sessionID)
  } catch (error) {
    console.error('取消银行打工失败', error)
  }
}

// 清除之前所有通过 setInterval 创建的定时器
function clearAllTimers() {
  timers.forEach((timer) => clearManagedTimer(timer))
  timers.length = 0
  cleanup()
}

// 监听用户年龄变化，当新年到来时强制结束打工
watch(() => gameStore.userInfo?.uage, (newAge, oldAge) => {
  // 当年龄增加时（新年），重新加载任务列表
  if (newAge && oldAge && newAge > oldAge) {
    // 如果正在打工，强制结束
    if (isRunning.value) {
      // 清除定时器
      clearAllTimers()
      // 重置打工状态
      isRunning.value = false
      customers.value = []
      score.value = 0
      void cancelActiveSession()
      // 提示用户
      ElMessage.warning('⏰ 工作时间结束，未完成的工作无报酬')
    }
  }
})

// 卸载组件后清空定时器
onBeforeUnmount(() => {
  // 清除之前的定时器
  clearAllTimers()
  // 清空客户
  customers.value = []
  void cancelActiveSession()
})

// emit事件
const emit = defineEmits(['refreshTasks'])
</script>

<style scoped>
/* 打工区域 */
.bank-work {
  flex: 2; /* ✅ 占据一半宽度 */
  border: 1px solid var(--border-color); /* ✅ 添加实线边框 */
  border-radius: 5px; /* ✅ 添加圆角边框 */
  background: var(--panel-color);
  padding: 10px; /* ✅ 添加内边距 */
  display: flex; /* ✅ 使子元素水平排列 */
  flex-direction: column; /* ✅ 子元素垂直排列 */
  overflow: hidden; /* 隐藏溢出内容 */
}

/* 打工游戏内容部分 */
.bank-work-panel {
  width: 100%;
  flex: 1 1 auto;                   /* 自动填充可用空间 */
  min-height: 0;                   /* 避免被内容撑高父容器 */
  display: flex;
  flex-direction: column;   /* ✅ 子元素垂直排列 */
  overflow: auto;                /* ✅ 避免内容撑开 chart，高度超出时滚动 */
  box-sizing: border-box;
  gap: 5px;
}

.work-score {
  font-size: 13px;
  margin-right: 5px;
}

/* 完成度颜色变化 */
.work-score-excellent { color: var(--success-color); }
.work-score-good { color: var(--warning-color); }
.work-score-normal { color: var(--info-color); }

.work-start {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
}

/* 打工提示 */
.work-tips {
  margin-top: 20px;
  text-align: center;
  color: var(--font-secondary);
  font-size: 13px;
  margin: 5px 0;
}

/* 滚动区域只包含数据行，支持垂直滚动，适用于显示多个条目 */
.work-panel-body {
  flex: 1;                                       /* 占满父容器剩余高度 */
  overflow-y: auto;                              /* 垂直方向滚动 */
  display: flex;                                 /* 使用 Flex 布局 */
  flex-direction: column;                        /* 子元素垂直排列 */
  margin-top: 3px;                               /* 顶部留白，避免与固定头部重叠 */
  gap: 5px;                                     /* 子元素之间有5px的间距 */
  /* 移除滚动条预留空间 */
  scrollbar-width: none;                         /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;                      /* IE/Edge 浏览器隐藏滚动条 */
}

.work-panel-body::-webkit-scrollbar {
  display: none;        /* Chrome/Safari 隐藏滚动条 */
}

/* 打工游戏按钮区域 */
.work-panel-button {
  height: 60px;                                   /* 固定高度 */
  display: flex;                              /* Flex 布局 */
  align-items: center;                        /* 垂直居中 */
  justify-content: space-between;            /* 子元素两端对齐 */
  border-radius: 8px;                         /* 圆角 */
}

.work-panel-button .el-button {
  width: 100%; /* 按钮宽度占满父容器 */
  margin: 0px 4px;
  transition: transform 0.2s ease;
}

/* 单个打工线选项 */
.work-panel-row {
  gap: 1px;
  border: 1px solid var(--font-light);  /* 分割线 */
  display: flex;                              /* Flex 布局 */
  align-items: center;                        /* 垂直居中 */
  justify-content: space-between;            /* 子元素两端对齐 */
  padding: 5px 5px;                          /* 内边距 */
  font-size: 14px;                            /* 字号 */
  border-radius: 8px;                         /* 圆角 */
}
.work-panel-row.active {
  border: 1px solid var(--info-color);
}

/* 打工线名字、时间、其他信息 */
.work-row-time, .work-row-name, .work-row-other {
  flex: 1;                                    /* 平均分配宽度 */
  text-align: center;                         /* 文本居中 */
  padding: 5px 0px;                          /* 内边距 */
  border-radius: 6px;                         /* 圆角 */
}
/* 打工线名字 */
.work-row-name {  
  background-color: var(--select-color);                 /* 淡蓝背景 */
  border: 1px solid var(--select-border-color);                  /*边框更蓝一些 */
}
/* 打工线时间 */
.work-row-time {
  width: 50px;
  text-align: center;
  color: var(--warning-color);
  font-weight: bold;
}
.work-row-other.active {
  border: 1px solid var(--info-color);
  background-color: var(--info-secondary-color);                 /* 淡蓝背景 */
}
</style>
