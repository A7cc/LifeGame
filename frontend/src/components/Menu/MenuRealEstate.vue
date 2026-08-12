<template>
  <div class="product-panel" data-testid="page-realestate">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title-block">
        <div class="panel-title">售楼中心</div>
        <div class="panel-subtitle">选择理想住所，购买成功后完成新居交付</div>
      </div>
      <div class="panel-controls">
        <el-tag type="success" size="small">名下 {{ Object.keys(userHouses || {}).length }} 套</el-tag>
      </div>
    </div>

    <!-- 中间信息区域 -->
    <div class="panel-main">
      <!-- 左侧：在售房屋 -->
      <div class="panel-main-left">
        <div class="section-header">
          <span class="section-title">在售房屋</span>
          <div class="section-controls">
            <el-tag size="small">{{ Object.keys(gameStore.gameInfo.ghouseinfo || {}).length - Object.keys(gameStore.userInfo.uhouse || {}).length }}/{{ Object.keys(gameStore.gameInfo.ghouseinfo || {}).length }}款</el-tag>
          </div>
        </div>
        <div class="section-content">
          <div v-if="Object.keys(gameStore.gameInfo.ghouseinfo || {}).length === 0" class="section-empty">
            <span class="section-empty-icon">🏠</span>
            <span class="section-empty-text">暂无房源</span>
          </div>
          <div class="property-grid">
            <div v-for="propertyinfo in gameStore.gameInfo.ghouseinfo" :key="propertyinfo.hiid" :data-testid="`house-${propertyinfo.hiid}`" class="property-card" :class="{ owned: userInfo.uhouse[propertyinfo.hiid], 'cannot-afford': userInfo.ucash < propertyinfo.hiprice }">
              <div class="card-header">
                <div class="card-avatar" @click.stop>
                  <el-image v-if="propertyinfo.hiimg" :src="propertyinfo.hiimg" fit="contain" class="avatar-image" @click.stop="openAssetPreview(propertyinfo)">
                    <template #error>
                      <span class="avatar-loading">🏠</span>
                    </template>
                  </el-image>
                </div>
                <div class="property-heading">
                  <span class="property-name">{{ propertyinfo.hiname }}</span>
                </div>
              </div>
              <div class="card-body">
                <div class="property-footer-row">
                  <div class="property-info">
                    <div class="property-bonus">
                      <span v-if="propertyinfo.hihealth > 0">💚 {{ propertyinfo.hihealth }} 健康</span>
                      <span v-if="propertyinfo.hifame > 0">⭐ {{ propertyinfo.hifame }} 名声</span>
                    </div>
                    <span class="property-price">💰 {{ formatPrice(propertyinfo.hiprice) }}</span>
                  </div>
                  <el-button :data-testid="`buy-house-${propertyinfo.hiid}`" v-if="!userInfo.uhouse[propertyinfo.hiid]" type="primary" size="small" :disabled="userInfo.ucash < propertyinfo.hiprice" @click="handleBuyClick(propertyinfo)">
                    {{ userInfo.ucash < propertyinfo.hiprice ? '资金不足' : '购买' }}
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
        <!-- 我的房产 -->
        <div class="garage-box">
          <div class="section-header">
            <span class="section-title">我的房产</span>
            <div class="section-controls">
              <el-tag type="success" size="small">{{ Object.keys(userHouses || {}).length }}套</el-tag>
            </div>
          </div>
          <div class="garage-content">
            <div v-if="Object.keys(userHouses || {}).length === 0" class="section-empty">
              <span class="section-empty-icon">🏠</span>
              <span class="section-empty-text">暂无房产</span>
            </div>
            <div v-else class="owned-houses-list">
            <div v-for="(ok, houseId) in userHouses" :key="houseId" :data-testid="`owned-house-${houseId}`" class="owned-property-card">
              <div class="owned-property-left">
                <div class="owned-property-image">
                  <el-image
                    :src="getHouseData(houseId, 'img')"
                    fit="contain"
                    class="owned-thumbnail"
                    :data-testid="`owned-house-image-${houseId}`"
                    @click.stop="openAssetPreview(getHouseInfo(houseId))"
                  >
                    <template #error><span class="owned-image-fallback">🏠</span></template>
                  </el-image>
                </div>
                <div class="property-details">
                  <span class="property-name">{{ getHouseData(houseId) }}</span>
                  <div class="property-meta-info">
                    <span class="property-sell-price">💰 {{ formatPrice(getHouseData(houseId, 'price')) }}</span>
                    <span class="property-bonus-inline">
                      <span v-if="getHouseData(houseId, 'health') > 0">💚 {{ getHouseData(houseId, 'health') }}</span>
                      <span v-if="getHouseData(houseId, 'fame') > 0">⭐ {{ getHouseData(houseId, 'fame') }}</span>
                    </span>
                  </div>
                </div>
              </div>
              <el-button :data-testid="`sell-house-${houseId}`" type="danger" size="small" @click="sellHouse(houseId)">出售</el-button>
            </div>
          </div>
        </div>
      </div>

      <LogPanel title="📋 交易记录" :items="runLogInfo" empty-icon="🏠" empty-text="暂无记录" @clear="clearLogs" />
      </div>
    </div>
  </div>
  <DialogAssetPurchaseMoment v-model="purchaseMomentVisible" kind="house" :item="purchasedHouse" />
  <DialogAssetPurchaseMoment
    v-model="assetPreviewVisible"
    mode="preview"
    kind="house"
    :item="previewedHouse"
    :owned="Boolean(userInfo.uhouse?.[previewedHouse?.hiid])"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useLogList } from '@/src/composables/useLogList'
import LogPanel from '@/src/components/Common/LogPanel.vue'
import DialogAssetPurchaseMoment from '@/src/components/Dialog/DialogAssetPurchaseMoment.vue'
import { BuyHouse, SellHouse } from "@/wailsjs/go/services/App.js"
import { formatPrice } from '@/src/utils/format'

const gameStore = useGameStore()
const userInfo = computed(() => gameStore.userInfo)

const { addLog, clearLogs, logs: runLogInfo } = useLogList({ green: '购买', red: '出售' })

// 用户房产
const userHouses = computed(() => userInfo.value.uhouse || [])
const purchaseMomentVisible = ref(false)
const purchasedHouse = ref(null)
const assetPreviewVisible = ref(false)
const previewedHouse = ref(null)

const getHouseInfo = houseId => gameStore.gameInfo.ghouseinfo?.[houseId] || null

const openAssetPreview = house => {
  if (!house?.hiimg) return
  previewedHouse.value = { ...house }
  assetPreviewVisible.value = true
}

// 获取房屋名称\健康\名声\价格
const getHouseData = (houseId, type="name") => {
  const house = getHouseInfo(houseId)
  if (type === "name") {
    return house ? house.hiname : '未知房屋'
  } else if (type === "health") {
    return house ? house.hihealth : 0
  } else if (type === "fame") {
    return house ? house.hifame : 0
  } else if (type === "img") {
    return house ? house.hiimg : ''
  } else if (type === "price") {
    // 出售价格 = 80% 原价
    if (!house) return 0
    return Math.floor(house.hiprice * 0.8)
  } else {
    return '未知房屋信息'
  }
}

// 点击购买按钮
const handleBuyClick = async (house) => {
  if (userInfo.value.ucash < house.hiprice) {
    ElMessage.warning('资金不足')
    return
  }

  try {
    const result = await BuyHouse(house.hiid)
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)
      ElMessage.success(`成功购买 ${house.hiname}！`)

      // 添加交易记录
      addLog(`${house.hiname}`, 'green')
      purchasedHouse.value = { ...house }
      purchaseMomentVisible.value = true
    } else {
      ElMessage.error(result.msg || '购买失败')
    }
  } catch (err) {
    ElMessage.error('购买失败')
  }
}

// 出售房屋
const sellHouse = async (houseId) => {
  try {
    const result = await SellHouse(Number(houseId))
    if (result.code === 200) {
      gameStore.applyUserInfo(result.userinfo)

      // 添加交易记录
      addLog(`${getHouseData(houseId)}`, 'red')

      ElMessage.success(result.msg)
    } else {
      ElMessage.error(result.msg || '出售失败')
    }
  } catch (err) {
    ElMessage.error('出售错误')
  }
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

/* ==================== 左侧：在售房屋 ==================== */
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

/* ==================== 房屋卡片 ==================== */
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
  color: var(--warning-color);
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

/* ==================== 我的房产 ==================== */
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

.owned-houses-list {
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

/* ==================== 交易日志记录 ==================== */
.log-panel {
  height: 140px;
}
</style>
