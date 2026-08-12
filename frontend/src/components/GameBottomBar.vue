<template>
  <div class="status-left" data-testid="game-status">
    <span>净资产：{{ formatPrice(gameStore.userInfo?.uassets) }}</span>
    <span>现金：{{ formatPrice(gameStore.userInfo?.ucash) }}</span>
    <span>存款：{{ formatPrice(gameStore.userInfo?.ubank) }}</span>
    <span>名声：{{ gameStore.userInfo?.ufame }}</span>
    <span>状态：{{ healthStatus }}</span>
  </div>
  <div class="status-right">
    <!-- 明年按钮 -->
    <el-button data-testid="next-year" type="primary" @click="nextTime" size="small">明年</el-button>
  </div>
</template>

<script setup>
import { NextTime } from "@/wailsjs/go/services/App.js";
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { formatPrice } from '@/src/utils/format'

// 用pinia存放数据
const gameStore = useGameStore()

const healthStatus = computed(() => {
  const user = gameStore.userInfo || {}
  if (Number(user.uimmunity || 0) < 10) return `危急（第${Number(user.ucriticalhealthyears || 1)}年，请就医）`
  if (Object.values(user.udiseases || {}).some(disease => Number(disease?.udseverity || 0) >= 5)) return '重症（需急诊）'
  return Object.keys(user.udiseases || {}).length > 0 ? '生病' : '健康'
})

// 定义事件用于传递数据
const emit = defineEmits(['exit', 'updateGameData']);

// 处理第二天按钮点击事件
const nextTime = async  () => {
  try {
    const gamedata = await NextTime();
    if (gamedata.code == 200) {
      // 将gamedata传递给父组件
      emit('updateGameData', gamedata);
    }else{
      ElMessage.error(gamedata.msg);
      
      emit('exit');
    }
  } catch (err) {
    ElMessage.error(err || '初始化失败');
  }
}
</script>

<style scoped>
.status-left span,
.status-right span {
  margin-right: 20px; /* 每个状态之间间距 */
}

.status-right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.status-right .el-button {
  font-size: 14px;  /* 字体大小更小 */
  padding: 4px 12px; /* 按钮内边距更小 */
}
</style>
