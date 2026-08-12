<template>
  <div class="product-panel" data-testid="page-carshop">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title-block">
        <div class="panel-title">汽车专卖店</div>
        <div class="panel-subtitle">挑选座驾，购买成功后完成新车交付</div>
      </div>
      <div class="panel-controls">
        <el-tag type="success" size="small">车库 {{ Object.keys(userCars || {}).length }} 辆</el-tag>
        <el-button size="small" type="primary" @click="showTaxiDialog = true">🚕 出车接单</el-button>
      </div>
    </div>

    <!-- 中间信息区域 -->
    <div class="panel-main">
      <!-- 左侧：在售车型 -->
      <div class="panel-main-left">
        <div class="section-header">
          <span class="section-title">在售车型</span>
          <div class="section-controls">
            <el-tag size="small">{{ Object.keys(gameStore.gameInfo.gcarinfo || {}).length - Object.keys(gameStore.userInfo.ucar || {}).length }}/{{ Object.keys(gameStore.gameInfo.gcarinfo || {}).length }}款</el-tag>
          </div>
        </div>
        <div class="section-content">
          <div v-if="Object.keys(gameStore.gameInfo.gcarinfo || {}).length === 0" class="section-empty">
            <span class="section-empty-icon">🚗</span>
            <span class="section-empty-text">暂无在售车型</span>
          </div>
          <div class="property-grid">
            <div v-for="propertyinfo in gameStore.gameInfo.gcarinfo" :key="propertyinfo.ciid" :data-testid="`car-${propertyinfo.ciid}`" class="property-card" :class="{ owned: userInfo.ucar[propertyinfo.ciid], 'cannot-afford': userInfo.ucash < propertyinfo.ciprice }">
              <div class="card-header">
                <div class="card-avatar" @click.stop>
                  <el-image :src="propertyinfo.ciimg" fit="contain" class="avatar-image" @click.stop="openAssetPreview(propertyinfo)">
                    <template #error>
                      <span class="avatar-loading">🚗</span>
                    </template>
                  </el-image>
                </div>
                <div class="property-heading">
                  <span class="property-name">{{ propertyinfo.ciname }}</span>
                </div>
              </div>
              <div class="card-body">
                <div class="property-footer-row">
                  <div class="property-info">
                    <div class="property-bonus">
                      <span v-if="propertyinfo.cihealth > 0">💚 {{ propertyinfo.cihealth }} 健康</span>
                      <span v-if="propertyinfo.cifame > 0">⭐ {{ propertyinfo.cifame }} 名声</span>
                    </div>
                    <span class="property-price">💰 {{ formatPrice(propertyinfo.ciprice) }}</span>
                  </div>
                  <el-button :data-testid="`buy-car-${propertyinfo.ciid}`" v-if="!userInfo.ucar[propertyinfo.ciid]" type="primary" size="small" :disabled="userInfo.ucash < propertyinfo.ciprice" @click="handleBuyClick(propertyinfo)">
                    {{ userInfo.ucash < propertyinfo.ciprice ? '资金不足' : '购买' }}
                  </el-button>
                  <div v-else class="owned-badge">已拥有</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧面板 -->
      <div class="panel-main-right">
        <!-- 我的车库 -->
        <div class="garage-box">
          <div class="section-header">
            <span class="section-title">我的车库</span>
            <div class="section-controls">
              <el-tag type="success" size="small">{{ Object.keys(userCars || {}).length }}辆</el-tag>
            </div>
          </div>
          <div class="garage-content">
            <div v-if="Object.keys(userCars || {}).length === 0" class="section-empty">
              <span class="section-empty-icon">🚗</span>
              <span class="section-empty-text">暂无车辆</span>
            </div>
            <div v-else class="owned-cars-list">
              <div
                v-for="(ok, carId) in userCars" :key="carId" :data-testid="`owned-car-${carId}`" class="owned-property-card">
                <div class="owned-property-left">
                  <div class="owned-property-image">
                    <el-image
                      :src="getCarData(carId, 'img')"
                      fit="contain"
                      class="owned-thumbnail"
                      :data-testid="`owned-car-image-${carId}`"
                      @click.stop="openAssetPreview(getCarInfo(carId))"
                    >
                      <template #error><span class="owned-image-fallback">🚗</span></template>
                    </el-image>
                  </div>
                  <div class="property-details">
                    <span class="property-name">{{ getCarData(carId) }}</span>
                    <div class="property-meta-info">
                      <span class="property-sell-price">💰 {{ formatPrice(getCarData(carId, 'price')) }}</span>
                      <span class="property-bonus-inline">
                        <span v-if="getCarData(carId, 'health') > 0">💚 {{ getCarData(carId, 'health') }}</span>
                        <span v-if="getCarData(carId, 'fame') > 0">⭐ {{ getCarData(carId, 'fame') }}</span>
                      </span>
                    </div>
                  </div>
                </div>
                <el-button :data-testid="`sell-car-${carId}`" type="danger" size="small" @click="sellCar(carId)">出售</el-button>
              </div>
            </div>
          </div>
        </div>

        <LogPanel title="📋 交易记录" :items="runLogInfo" empty-icon="🚗" empty-text="暂无记录" @clear="clearLogs" />
      </div>
    </div>
  </div>

  <!-- 出车接单游戏 -->
  <OtherGameTaxiDriver v-model="showTaxiDialog" @complete="handleTaxiComplete" />
  <DialogAssetPurchaseMoment v-model="purchaseMomentVisible" kind="car" :item="purchasedCar" />
  <DialogAssetPurchaseMoment
    v-model="assetPreviewVisible"
    mode="preview"
    kind="car"
    :item="previewedCar"
    :owned="Boolean(userInfo.ucar?.[previewedCar?.ciid])"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useLogList } from '@/src/composables/useLogList'
import LogPanel from '@/src/components/Common/LogPanel.vue'
import DialogAssetPurchaseMoment from '@/src/components/Dialog/DialogAssetPurchaseMoment.vue'
import { BuyCar, SellCar } from "@/wailsjs/go/services/App.js"
import OtherGameTaxiDriver from '@/src/components/GameList/WorkGameTaxiDriver.vue'
import { formatPrice } from '@/src/utils/format'

const gameStore = useGameStore()
const userInfo = computed(() => gameStore.userInfo)

const { addLog, clearLogs, logs: runLogInfo } = useLogList({ green: '购买', red: '出售', blue: '出车' })

// 用户车辆
const userCars = computed(() => userInfo.value.ucar || [])

// 出租车游戏
const showTaxiDialog = ref(false)
const purchaseMomentVisible = ref(false)
const purchasedCar = ref(null)
const assetPreviewVisible = ref(false)
const previewedCar = ref(null)

const getCarInfo = carId => gameStore.gameInfo.gcarinfo?.[carId] || null

const openAssetPreview = car => {
  if (!car?.ciimg) return
  previewedCar.value = { ...car }
  assetPreviewVisible.value = true
}

// 获取车辆名称
const getCarData = (carId, type="name") => {
  const car = getCarInfo(carId)
  if (type === "name") {
    return car ? car.ciname : '未知车辆'
  } else if (type === "health") {
    return car ? car.cihealth : 0
  } else if (type === "fame") {
    return car ? car.cifame : 0
  } else if (type === "img") {
    return car ? car.ciimg : ''
  } else if (type === "price") {
    // 出售价格 = 70% 原价
    if (!car) return 0
    return Math.floor(car.ciprice * 0.7)
  } else {
    return '未知车辆信息'
  }
}

// 点击购买按钮
const handleBuyClick = async (car) => {
  if (userInfo.value.ucash < car.ciprice) {
    ElMessage.warning('资金不足')
    return
  }

  try {
    const result = await BuyCar(car.ciid)
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      ElMessage.success(`成功购买 ${car.ciname}！`)

      // 添加交易记录
      addLog(`${car.ciname}`, 'green')
      purchasedCar.value = { ...car }
      purchaseMomentVisible.value = true
    } else {
      ElMessage.error(result.msg || '购买失败')
    }
  } catch (err) {
    ElMessage.error('购买失败')
  }
}

// 出售车辆
const sellCar = async (carId) => {
  try {
    const result = await SellCar(Number(carId))
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)

      // 添加交易记录
      addLog(`${getCarData(carId)}`, 'red')

      ElMessage.success(result.msg)
    } else {
      ElMessage.error(result.msg || '出售失败')
    }
  } catch (err) {
    ElMessage.error('出售错误')
  }
}

// 出租车游戏完成处理
const handleTaxiComplete = (result) => {
  addLog(`出车接单 - ${result.rating}星 - +${result.earnings}元`, 'blue')
}
</script>

<style scoped>
.panel-title-block {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.panel-subtitle {
  color: var(--font-secondary);
  font-size: 10px;
}

/* ==================== 主内容区域 ==================== */
.panel-main {
  display: flex;
  gap: 5px;
  margin-top: 10px;
  flex: 1;
  min-height: 0;
}

/* ==================== 左侧：在售车型 ==================== */
.panel-main-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.section-content {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.section-content::-webkit-scrollbar {
  display: none;
}

.property-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 7px;
}

/* ==================== 车子卡片 ==================== */
.property-card {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.property-card {
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.property-card.owned {
  opacity: 0.82;
}

.property-card.cannot-afford {
  opacity: 0.62;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 10px;
  background: var(--panel-color);
  border-bottom: 1px solid var(--border-color);
}

.card-header .property-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.card-avatar {
  width: 45px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-image {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid var(--border-color);
}

.property-heading {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 3px;
}

.avatar-loading {
  /* 关键代码：开启 Flex 布局实现居中 */
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center;     /* 垂直居中 */
  
  /* 确保容器填满图片区域 */
  width: 100%;
  height: 100%;
  
  /* 背景色和字体样式（可选） */
  background-color: var(--panel-color);
  font-size: 24px;           /* Emoji 大小 */
  color: var(--border-color);
}

.property-price {
  font-size: 11px;
  font-weight: 600;
  color: var(--warning-color);
}

.card-body {
  padding: 8px 10px;
  flex: 1;
}

.property-footer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 5px;
}

.property-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.property-bonus {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  color: var(--success-color);
  font-weight: 500;
  font-size: 10px;
}

.owned-badge {
  color: var(--font-color);
  font-size: 12px;
  font-weight: 500;
}

/* ==================== 右侧面板 ==================== */
.panel-main-right {
  width: 280px;
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex-shrink: 0;
}

/* ==================== 我的车库 ==================== */
.garage-box {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  padding: 12px;
  min-height: 0;
}

.garage-content {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.garage-content::-webkit-scrollbar {
  display: none;
}

.owned-cars-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.owned-property-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 7px;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.owned-property-card:hover {
  border-color: var(--primary-color);
  background: var(--panel-main-color);
}

.owned-property-left {
  display: flex;
  align-items: center;
  min-width: 0;
  gap: 8px;
}

.owned-property-image {
  display: flex;
  width: 46px;
  height: 46px;
  flex: 0 0 46px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--panel-main-color);
  box-shadow: inset 0 0 0 1px rgb(255 255 255 / 4%);
}

.owned-thumbnail {
  width: 42px;
  height: 42px;
  cursor: zoom-in;
}

.owned-image-fallback {
  display: flex;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.property-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--font-color);
}

.property-details {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.property-sell-price {
  font-size: 10px;
  color: var(--warning-color);
}

.property-meta-info {
  display: flex;
  gap: 5px;
  font-size: 10px;
}

.property-bonus-inline {
  display: flex;
  gap: 5px;
  font-size: 10px;
  color: var(--success-color);
}

/* ==================== 交易记录框 ==================== */
.log-panel {
  height: 140px;
}
</style>
