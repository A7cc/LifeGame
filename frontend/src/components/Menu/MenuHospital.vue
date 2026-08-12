<template>
  <div class="product-panel" data-testid="page-hospital">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">医院健康中心</div>
      <div class="panel-controls">
        <el-tag size="small">{{ userInfo.uimmunity }}/{{ gameStore.gameInfo.gmaxholdnum.miholdnum }}</el-tag>
      </div>
    </div>

    <div class="hospital-main">
      <!-- 左侧：健康状态和疾病 -->
      <div class="health-panel">
        <div class="section-header">
          <span class="section-title">身体状态</span>
          <div v-if="emergencyStatus?.required"  class="section-controls">
            <el-button type="danger" size="small" :loading="emergencyLoading" @click="handleEmergencyTreatment">接受急诊</el-button>
          </div>
        </div>
        <div class="section-content">
          <div v-if="Object.keys(userInfo.udiseases || {}).length === 0" class="section-empty">
            <span class="section-empty-icon">💪</span>
            <span class="section-empty-text">{{ immunityAdvice }}</span>
          </div>
          <!-- 存在严重疾病 -->
          <div v-else-if="emergencyStatus?.required" class="emergency-panel">
            <div class="emergency-main">
              <div class="emergency-title">⚠️ 必须立即急诊</div>
              <div class="emergency-reasons">
                {{ emergencyStatus.reasons?.join('；') || '身体状态危险' }}
              </div>
              <div class="emergency-cost">💰 {{ emergencyStatus.cost }}</div>
            </div>
          </div>

          <!-- 当前疾病 -->
          <div v-else class="log-list">
            <div v-for="(d, did) in userInfo.udiseases" :key="did" class="log-item">
              <span class="log-action" :class="d.udtype">{{ d.udname }}</span>
              <span class="log-item-text">{{ d.usymptoms }}</span>
              <span>{{ getSeverity(d.udseverity) }}</span>
            </div>
          </div>

        </div>
      </div>

      <!-- 右侧：医院卡片 -->
      <div class="hospital-cards">
        <div class="section-header">
          <span class="section-title">医疗机构</span>
        </div>
        <div class="cards-grid">
          <div v-for="card in hospitalCards" :key="card.htype" :data-testid="`hospital-${card.htype}`" class="hospital-card" @click="selectHospital(card)">
            <div class="card-icon">{{ card.hicon }}</div>
            <div class="card-name">{{ card.hname }}</div>
            <div class="card-desc">{{ card.hdescription }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部：游戏区和日志 -->
    <div class="hospital-bottom">
      <!-- 左侧：游戏区 -->
      <div class="game-panel">
        <div class="section-header">
          <span class="section-title">{{ currentHospital?.hicon }} {{ currentHospital?.hname || '选择医院' }}</span>
          <el-button v-if="isPlaying" size="small" type="danger" @click="endGame">结束</el-button>
          <el-button v-if="currentHospital && !isPlaying" size="small" type="primary" @click="currentHospital = null">返回</el-button>
        </div>

        <div class="game-content">
          <!-- 未选择医院 -->
          <div v-if="!currentHospital" class="game-empty">
            <span>👆 请选择一家医疗机构</span>
          </div>

          <!-- 服务选择 -->
          <div v-else-if="!isPlaying" class="service-select">
            <div class="service-header">
              <span class="service-title">选择服务</span>
            </div>
            <div class="service-list">
              <div v-for="svc in currentHospital.hservices" :key="svc.hsid" :data-testid="`hospital-service-${svc.hsid}`" class="service-item" @click="startService(svc)">
                <span class="svc-name">{{ svc.hsname }}</span>
                <span class="svc-price" v-if="svc.hsprice > 0">💰 {{ svc.hsprice }}</span>
                <span class="svc-desc">{{ svc.hsdesc }}</span>
              </div>
            </div>
          </div>

          <!-- 买药游戏 -->
          <div v-else-if="gameType === 'medicine'" class="medicine-game">
            <div class="symptom-box">
              <div class="symptom-title">你的症状：</div>
              <div class="symptom-text">{{ currentSymptom }}</div>
            </div>
            <div class="medicine-list">
              <div v-for="(t, tid) in availableTreatments" :key="tid" :data-testid="`medicine-${tid}`" class="medicine-item" @click="buyTreatment(t, parseInt(tid))">
                <div class="med-name">{{ t.tname }}</div>
                <div class="med-desc">{{ t.tdesc }}</div>
                <div class="med-price">💰 {{ t.tprice }}</div>
              </div>
            </div>
          </div>

          <!-- 打针 -->
          <div v-else-if="gameType === 'injection'" class="surgery-game">
            <div class="surgery-confirm">
              <div class="surgery-warning">⚠️ 打针风险提示</div>
              <div class="surgery-info">
                <p>💰 费用：价格中等</p>
                <p>🏥 可治愈：大部分疾病</p>
              </div>
              <el-button type="danger" size="small" @click="doSurgery('打针')" :loading="surgeryLoading">
                确认打针
              </el-button>
            </div>
          </div>

          <!-- 针灸 -->
          <div v-else-if="gameType === 'acupuncture'" class="surgery-game">
            <div class="surgery-confirm">
              <div class="surgery-warning">⚠️ 针灸风险提示</div>
              <div class="surgery-info">
                <p>💰 费用：价格中等</p>
                <p>🏥 可治愈：扭伤、腰肌劳损</p>
              </div>
              <el-button type="danger" size="small" @click="doSurgery('针灸')" :loading="surgeryLoading">
                确认针灸
              </el-button>
            </div>
          </div>

          <!-- 手术 -->
          <div v-else-if="gameType === 'surgery'" class="surgery-game">
            <div class="surgery-confirm">
              <div class="surgery-warning">⚠️ 手术风险提示</div>
              <div class="surgery-info">
                <p>💰 费用：价格昂贵</p>
                <p>🏥 可治愈：重大疾病</p>
              </div>
              <el-button type="danger" size="small" @click="doSurgery('手术')" :loading="surgeryLoading">
                确认手术
              </el-button>
            </div>
          </div>
          <!-- 脱胎换骨 -->
          <div v-else-if="gameType === 'bereborn'" class="surgery-game">
            <div class="surgery-confirm">
              <div class="surgery-warning">⚠️ 脱胎换骨风险提示</div>
              <div class="surgery-info">
                <p>💰 费用：价格昂贵</p>
                <p>🏥 可治愈：专治压力疾病</p>
              </div>
              <el-button type="danger" size="large" @click="doSurgery('脱胎换骨')" :loading="surgeryLoading">
                确认脱胎换骨
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <LogPanel title="📋 治疗记录" :items="runLogInfo" empty-icon="📋" empty-text="暂无治疗记录" @clear="clearLogs" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useLogList } from '@/src/composables/useLogList'
import LogPanel from '@/src/components/Common/LogPanel.vue'
import { GetHospitalInfo, BuyTreatment, SpecialTreatment, EmergencyTreatment } from "@/wailsjs/go/services/App.js"

const gameStore = useGameStore()
const userInfo = computed(() => gameStore.userInfo)

// 医院数据
const hospitalCards = ref([])
const emergencyStatus = ref(null)
const emergencyLoading = ref(false)
const currentHospital = ref(null)
const isPlaying = ref(false)
const gameType = ref('')

// 药品游戏
const availableTreatments = ref(null)
const currentSymptom = ref('')
// 治疗
const surgeryLoading = ref(false)

const { addLog, clearLogs, logs: runLogInfo } = useLogList()

// 根据免疫力计算健康状态
const immunityAdvice = computed(() => {
  const h = userInfo.value.uimmunity
  if (h >= 90) return "💪 免疫力非常强！身体很健康！"
  if (h >= 80) return "💚 免疫力很强，状态优秀"
  if (h >= 70) return "💛 免疫力良好，继续保持"
  if (h >= 50) return "🧡 免疫力一般，注意休息"
  if (h >= 30) return "🧠 免疫力较低，建议就医检查"
  if (h >= 10) return "❗ 免疫力较低，建议就医检查"
  return "💔 免疫力很差！必须立即就医！"
})

// 计算严重程度
const getSeverity = (severity) => {
  if (severity === 1) return '🟢'
  if (severity === 2) return '🟡'
  if (severity === 3) return '🟠'
  if (severity === 4) return '🔴'
  return '💔'
}

// 监听用户疾病变化，及时更新疾病列表
watch(() => [userInfo.value.uimmunity, userInfo.value.udiseases],
  () => {
    loadHospitalInfo()
  },
  { deep: true }
)

const selectHospital = (card) => {
  currentHospital.value = card
  isPlaying.value = false
  gameType.value = ''
}

const startService = async (svc) => {
  if (emergencyStatus.value?.required) {
    ElMessage.warning(`当前必须先接受急诊治疗，费用 ${emergencyStatus.value.cost} 元`)
    return
  }

  gameType.value = svc.hstype
  isPlaying.value = true

  if (svc.hstype === 'medicine') {
    startMedicineGame()
  }
}

// 买药游戏
const startMedicineGame = () => {
  // 生成症状描述
  const diseases = Object.values(userInfo.value.udiseases || {})
  currentSymptom.value = diseases.length > 0
    ? diseases.map(d => `${d.udname}：${d.usymptoms}`).join('；') : '身体不适，建议检查'

  // 获取可用药品
  GetHospitalInfo().then(res => {
    if (res.code === 200) {
      // 根据医院类型过滤药品
      const treatments = res.treatments || {}
      const result = {}
      Object.values(treatments).forEach(t => {
        if (t.tsource === currentHospital.value.htype || t.tsource === 'pharmacy') {
          result[t.tid] = t
        }
      })

      availableTreatments.value = result
      emergencyStatus.value = res.emergency || null
    }
  })
}

// 购买药品
const buyTreatment = async (treat, tid) => {
  if (emergencyStatus.value?.required) {
    ElMessage.warning(`当前必须先接受急诊治疗，费用 ${emergencyStatus.value.cost} 元`)
    return
  }

  if (userInfo.value.ucash < treat.tprice) {
    ElMessage.warning('现金不足')
    return
  }

  try {
    // 当前医院类型，药品id
    const result = await BuyTreatment(currentHospital.value.htype, parseInt(tid))
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      if (result.curedDiseases.length > 0) {
        addLog(`💊 购买${treat.tname}，治疗了${result.curedDiseases.length}种疾病`, 'green')
      }else{
        addLog(`💊 购买${treat.tname}，增加了免疫力`, 'red')
      }
      
      ElMessage.success(result.msg)
      loadHospitalInfo()
    } else {
      if (result.emergency) {
        emergencyStatus.value = result.emergency
      }
      ElMessage.error(result.msg)
    }
  } catch (e) {
    ElMessage.error('购买失败')
  }
}

// 治疗大型疾病
const doSurgery = async (treatmentType) => {
  if (emergencyStatus.value?.required) {
    ElMessage.warning(`当前必须先接受急诊治疗，费用 ${emergencyStatus.value.cost} 元`)
    return
  }

  surgeryLoading.value = true

  try {
    const result = await SpecialTreatment(treatmentType,currentHospital.value.htype)
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      addLog(`${result.msg}`, 'green')
      ElMessage.success(result.msg)
      loadHospitalInfo()
    } else {
      if (result.emergency) {
        emergencyStatus.value = result.emergency
      }
      ElMessage.error(result.msg)
    }
  } catch (e) {
    ElMessage.error(treatmentType+'失败')
  }

  surgeryLoading.value = false
  endGame()
}

const handleEmergencyTreatment = async () => {
  if (!emergencyStatus.value?.required) return

  emergencyLoading.value = true
  try {
    const result = await EmergencyTreatment()
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      emergencyStatus.value = result.emergency
      addLog(`${result.msg}`, 'green')
      ElMessage.success(result.msg)
      endGame()
      loadHospitalInfo()
    } else {
      if (result.emergency) {
        emergencyStatus.value = result.emergency
      }
      ElMessage.error(result.msg)
    }
  } catch (e) {
    ElMessage.error('急诊治疗失败')
  } finally {
    emergencyLoading.value = false
  }
}

const endGame = () => {
  isPlaying.value = false
  gameType.value = ''
}

const loadHospitalInfo = async () => {
  try {
    const result = await GetHospitalInfo()
    if (result.code === 200) {
      hospitalCards.value = result.hospitalcards || []
      emergencyStatus.value = result.emergency || null
    }
  } catch (e) {
    console.error('加载医院信息失败', e)
  }
}

onMounted(() => {
  loadHospitalInfo()
})
</script>

<style scoped>
/* 急诊面板 */
.emergency-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #fff5f5 0%, #ffe8e8 100%);
  border-radius: 8px;
  border: 2px solid var(--error-color);
  box-shadow: 0 2px 8px rgba(220, 53, 69, 0.15);
}

.emergency-main {
  text-align: center;
}

.emergency-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--error-color);
}

.emergency-reasons {
  margin: 6px 0px;
  font-size: 12px;
  color: #666;
  max-width: 240px;
  word-break: break-word;
}

.emergency-cost {
  font-size: 13px;
  font-weight: 600;
  color: var(--warning-color);
  white-space: nowrap;
}

.hospital-main {
  display: flex;
  gap: 5px;
  margin-top: 5px;
  height: 180px;
}

.health-panel, .hospital-cards {
  border: 1px solid var(--border-color);
  border-radius: 5px;
  background: var(--panel-color);
  padding: 10px;
  display: flex;
  flex-direction: column;
}

.health-panel {
  width: 300px;
  flex-shrink: 0;
}

.hospital-cards {
  flex: 1;
}

/* 医院卡片 */
.cards-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.hospital-card {
  background: var(--panel-sub-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.hospital-card:hover {
  background-color: var(--select-color);         /* 根据主题变化 */
}

.card-icon {
  font-size: 28px;
  margin-bottom: 6px;
}

.card-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-color);
  margin-bottom: 4px;
}

.card-desc {
  font-size: 11px;
  color: var(--font-secondary);
}

/* 底部区域 */
.hospital-bottom {
  flex: 1;
  display: flex;
  gap: 8px;
  margin-top: 8px;
  min-height: 0;
}

.game-panel {
  border: 1px solid var(--border-color);
  border-radius: 5px;
  background: var(--panel-color);
  padding: 10px;
  display: flex;
  flex-direction: column;
  flex: 2;
}

/* 游戏内容 */
.game-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.game-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--font-secondary);
  font-size: 14px;
}

.service-select {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.service-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--font-secondary);
}

.service-list {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.service-item {
  width: 130px;
  background: var(--panel-sub-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 10px 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.service-item:hover {
  border-color: var(--select-border-color);
  background: var(--select-color);
}

.svc-name {
  font-weight: 600;
  color: var(--font-color);
}

.svc-price {
  color: var(--warning-color);
  font-size: 12px;
  margin-left: 8px;
}

.svc-desc {
  display: block;
  font-size: 11px;
  color: var(--font-secondary);
  margin-top: 4px;
}

/* 买药游戏 */
.medicine-game {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
  overflow: hidden;
}

.symptom-box {
  background: #fff5f5;
  border: 1px solid var(--error-secondary-color);
  border-radius: 5px;
  padding: 10px;
}

.symptom-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--error-color);
  margin-bottom: 6px;
}

.symptom-text {
  font-size: 12px;
  color: var(--font-light);
}

.medicine-list {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 8px;
  overflow-y: auto;
  /* 移除滚动条预留空间 */
  scrollbar-width: none;              /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;           /* IE/Edge 浏览器隐藏滚动条 */
}

.medicine-list::-webkit-scrollbar {
  display: none;          /* Chrome/Safari 隐藏滚动条 */
}

.medicine-item {
  background: var(--panel-sub-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.medicine-item:hover {
  border-color: var(--select-border-color);
  background: var(--select-color);
}

.med-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
}

.med-desc {
  font-size: 11px;
  color: var(--font-secondary);
  margin: 4px 0;
}

.med-price {
  font-size: 11px;
  color: var(--warning-color);
}

/* 手术 */
.surgery-game {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.surgery-confirm {
  text-align: center;
  padding: 20px;
  background: var(--panel-sub-color);
  border-radius: 5px;
  border: 1px solid var(--border-color);
}

.surgery-warning {
  font-size: 16px;
  font-weight: 600;
  color: var(--warning-color);
  margin-bottom: 15px;
}

.surgery-info {
  text-align: left;
  margin-bottom: 20px;
}

.surgery-info p {
  margin: 8px 0;
  font-size: 13px;
  color: var(--font-secondary);
}

/* 日志 */
.log-panel {
  flex: 1;
}
</style>
