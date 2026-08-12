<template>
  <div class="product-panel" data-testid="page-profile">
    <div class="panel-header">
      <div class="panel-title">角色档案</div>
    </div>

    <div class="profile-main">
      <!-- 顶部 HUD 区域 -->
      <div class="hud">
        <!-- 左：头像 -->
        <div class="hud-avatar">
          <span class="hud-emoji">{{ userInfo.usex ? '👨' : '👩' }}</span>
          <div class="hud-identity">
            <div class="hud-name">{{ userInfo.uname }}</div>
            <div class="hud-sub">{{ userInfo.uage }}岁 · {{ fameLevelText }}</div>
          </div>
        </div>
        <!-- 中：进度条 -->
        <div class="hud-bars">
          <div class="bar-item">
            <span class="bar-label">🛡️ 免疫</span>
            <div class="bar-track"><div class="bar-fill hp" :style="{ width: (userInfo.uimmunity / gameStore.gameInfo.gmaxholdnum.miholdnum * 100) + '%' }"></div></div>
            <span class="bar-val">{{ userInfo.uimmunity }} / {{ gameStore.gameInfo.gmaxholdnum.miholdnum }}</span>
          </div>
          <div class="bar-item">
            <span class="bar-label">⭐ 名声</span>
            <div class="bar-track"><div class="bar-fill fame" :style="{ width: (userInfo.ufame / gameStore.gameInfo.gmaxholdnum.mfaholdnum * 100) + '%' }"></div></div>
            <span class="bar-val">{{ userInfo.ufame }} / {{ gameStore.gameInfo.gmaxholdnum.mfaholdnum }}</span>
          </div>
        </div>
        <!-- 右：机会次数 -->
        <div class="hud-opportunity">
          <div class="opp-grid">
            <div class="opp-tag">
              <span class="opp-icon">💼</span>
              <span class="opp-num">{{ userInfo.uopportunity.ownum }}/{{ gameStore.gameInfo.gmaxholdnum.mwroundnum }}</span>
            </div>
            <div class="opp-tag">
              <span class="opp-icon">🎮</span>
              <span class="opp-num">{{ userInfo.uopportunity.ognum }}/{{ gameStore.gameInfo.gmaxholdnum.mgroundnum }}</span>
            </div>
            <div class="opp-tag">
              <span class="opp-icon">💑</span>
              <span class="opp-num">{{ userInfo.uopportunity.omnum }}/{{ gameStore.gameInfo.gmaxholdnum.mmroundnum }}</span>
            </div>
            <div class="opp-tag">
              <span class="opp-icon">🛍️</span>
              <span class="opp-num">{{ userInfo.uopportunity.osnum }}/{{ gameStore.gameInfo.gmaxholdnum.msroundnum }}</span>
            </div>
            <div class="opp-tag">
              <span class="opp-icon">🏺</span>
              <span class="opp-num">{{ userInfo.uopportunity.oanum }}/{{ gameStore.gameInfo.gmaxholdnum.maroundnum }}</span>
            </div>
          </div>
        </div>
      </div>
      <!-- 疾病情况 -->
      <div class="disease-row" :class="{ healthy: Object.keys(userInfo.udiseases || {}).length === 0 && userInfo.uimmunity >= 10, critical: userInfo.uimmunity < 10 }">
        <div class="disease-title">🏃 身体状况</div>
        <div class="disease-content" v-if="Object.keys(userInfo.udiseases || {}).length === 0">
          <span class="health-icon">{{ userInfo.uimmunity < 10 ? '🚑' : '💚' }}</span>
          <span class="health-text" :class="{ critical: userInfo.uimmunity < 10 }">{{ healthStateText }}</span>
        </div>
        <div class="disease-scroll" v-else>
          <span v-if="userInfo.uimmunity < 10" class="critical-health-text">🚑 低免疫第{{ userInfo.ucriticalhealthyears || 1 }}年</span>
          <div v-for="(d, did) in userInfo.udiseases" :key="did" class="disease-card">
            <span class="dc-name">{{ d.udname }}</span>
            <span class="dc-sev">{{ getSeverity(d.udseverity) }}</span>
          </div>
        </div>
      </div>
      <!-- 资产明细 -->
      <div class="finance-row">
        <div class="finance-cell total">
          <span class="fc-icon">💰</span>
          <div class="fc-text">
            <span class="fc-label">总资产</span>
            <span class="fc-val">{{ formatPrice(userInfo.uassets) }}</span>
          </div>
        </div>
        <div class="finance-cell">
          <span class="fc-icon">💵</span>
          <div class="fc-text">
            <span class="fc-label">现金</span>
            <span class="fc-val">{{ formatPrice(userInfo.ucash) }}</span>
          </div>
        </div>
        <div class="finance-cell">
          <span class="fc-icon">🏦</span>
          <div class="fc-text">
            <span class="fc-label">存款</span>
            <span class="fc-val">{{ formatPrice(userInfo.ubank) }}</span>
          </div>
        </div>
        <div class="finance-cell">
          <span class="fc-icon">📋</span>
          <div class="fc-text">
            <span class="fc-label">贷款</span>
            <span class="fc-val" :class="{ danger: userInfo.uloan > 0 }">{{ formatPrice(userInfo.uloan) }}</span>
          </div>
        </div>
        <div class="finance-cell" v-if="userInfo.uloanoverdue > 0">
          <span class="fc-icon">⚠️</span>
          <div class="fc-text">
            <span class="fc-label">逾期</span>
            <span class="fc-val danger">{{ userInfo.uloanoverdue }} 年</span>
          </div>
        </div>
      </div>

      <!-- 持有物品 -->
      <div class="holdings">
        <div class="holdings-tabs folder">
          <span v-for="tab in tabs" :key="tab.key" class="tab-btn" :class="{ active: activeTab === tab.key }" @click="activeTab = tab.key">{{ tab.icon }} {{ tab.label }}<small v-if="tab.count">{{ tab.count }}</small></span>
        </div>
        <div class="tab-content">
          <!-- 没有物品 -->
          <div class="section-empty" v-if="totalHoldings <= 0">
            <span class="section-empty-icon">📦</span>
            <span class="section-empty-text">暂无持有物品</span>
          </div>
          <!-- 有物品 -->
          <!-- 全部 -->
          <template v-if="activeTab === 'all'">
            <div v-for="group in allGroups" :key="group.key" class="group">
              <div class="group-title">{{ group.icon }} {{ group.label }}</div>
              <div class="card-grid">
                <div v-for="item in group.items" :key="item.id" class="item-card">
                  <div class="card-left">
                    <el-image v-if="item.img" :src="item.img" fit="cover" class="card-img" :preview-src-list="[item.img]" preview-teleported />
                    <span v-else class="card-icon">{{ item.icon }}</span>
                  </div>
                  <div class="card-right">
                    <div class="card-name">{{ item.name }}</div>
                    <div class="card-tag" v-if="item.tag">{{ item.tag }}</div>
                    <div class="card-price" v-if="item.price" :class="item.priceClass">{{ item.price }}</div>
                  </div>
                </div>
              </div>
            </div>
          </template>
          <!-- 单分类 -->
          <template v-else>
            <div v-for="group in allGroups.filter(g => g.key === activeTab)" :key="group.key" class="group">
              <div class="card-grid">
                <div v-for="item in group.items" :key="item.id" class="item-card">
                  <div class="card-left">
                    <el-image v-if="item.img" :src="item.img" fit="cover" class="card-img" :preview-src-list="[item.img]" preview-teleported />
                    <span v-else class="card-icon">{{ item.icon }}</span>
                  </div>
                  <div class="card-right">
                    <div class="card-name">{{ item.name }}</div>
                    <div class="card-tag" v-if="item.tag">{{ item.tag }}</div>
                    <div class="card-price" v-if="item.price" :class="item.priceClass">{{ item.price }}</div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useGameStore } from '@/src/stores/game'
import { formatPrice } from '@/src/utils/format'

const gameStore = useGameStore()
const userInfo = computed(() => gameStore.userInfo)
const activeTab = ref('all')

const healthStateText = computed(() => {
  if (userInfo.value.uimmunity < 10) {
    const years = Number(userInfo.value.ucriticalhealthyears || 1)
    return years >= 2 ? '最后抢救期，请立即去医院' : `低免疫危机第${years}年，仍有治疗机会`
  }
  return '身体健康，无疾病'
})

// 计算严重程度
const getSeverity = (severity) => {
  if (severity === 1) return '🟢'
  if (severity === 2) return '🟡'
  if (severity === 3) return '🟠'
  if (severity === 4) return '🔴'
  return '💔'
}

const fameLevelText = computed(() => {
  const level = gameStore.calcReputationLevel()
  const texts = ['普通', '中等', '高级', '豪华', '私人']
  return level >= 0 ? `${texts[level]}` : '老赖'
})

const ownedHouses = computed(() => {
  const houses = userInfo.value.uhouse || {}
  const houseInfo = gameStore.gameInfo.ghouseinfo || {}
  const result = []
  for (const [id, owned] of Object.entries(houses)) {
    if (owned && houseInfo[id]) result.push({ ...houseInfo[id], hiid: id })
  }
  return result
})

const ownedCars = computed(() => {
  const cars = userInfo.value.ucar || {}
  const carInfo = gameStore.gameInfo.gcarinfo || {}
  const result = []
  for (const [id, owned] of Object.entries(cars)) {
    if (owned && carInfo[id]) result.push({ ...carInfo[id], ciid: id })
  }
  return result
})

const stockList = computed(() => {
  const stocks = userInfo.value.ustock || {}
  const result = []
  for (const [id, stock] of Object.entries(stocks)) {
    if (stock && stock.usnum > 0) result.push({ id, name: stock.usname || `股票#${id}`, num: stock.usnum, profit: stock.usprofit || 0 })
  }
  return result
})

const companyList = computed(() => {
  const companies = userInfo.value.ucompany || {}
  const companyInfo = gameStore.gameInfo.gcompanyinfo || {}
  const result = []
  for (const [id, company] of Object.entries(companies)) {
    const info = companyInfo[id]
    result.push({ id, name: info?.ciname || `公司#${id}`, num: company.ucompanynum || 0, costprice: company.ucompanycostprice || 0 })
  }
  return result
})

const datingList = computed(() => {
  const dating = userInfo.value.udating || {}
  const result = []
  for (const [id, d] of Object.entries(dating)) {
    result.push({ id, name: d.dname || `对象#${id}`, level: d.dstatus || '陌生人', count: d.dcount || 0 })
  }
  return result
})

const getAffinityIcon = (level) => {
  const icons = { '陌生人': '👋', '朋友': '🙂', '暧昧中': '💕', '交往中': '💑', '恋人': '💚', '专属恋人': '💖', '爱人': '👫', '已婚': '💍', '前任': '💔' }
  return icons[level] || '💕'
}

// 统一的物品分组
const allGroups = computed(() => {
  const groups = []
  if (ownedHouses.value.length) groups.push({
    key: 'house', icon: '🏠', label: '房产', items: ownedHouses.value.map(h => ({
      id: h.hiid, img: h.hiimg, icon: '🏠', name: h.hiname,
      tag: [h.hihealth && `💚+${h.hihealth}`, h.hifame && `⭐+${h.hifame}`].filter(Boolean).join(' ') || null,
      price: formatPrice(h.hiprice)
    }))
  })
  if (ownedCars.value.length) groups.push({
    key: 'car', icon: '🚗', label: '汽车', items: ownedCars.value.map(c => ({
      id: c.ciid, img: c.ciimg, icon: '🚗', name: c.ciname,
      tag: [c.cihealth && `💚+${c.cihealth}`, c.cifame && `⭐+${c.cifame}`].filter(Boolean).join(' ') || null,
      price: formatPrice(c.ciprice)
    }))
  })
  if (userInfo.value.uantique?.length) groups.push({
    key: 'antique', icon: '🖼️', label: '古董', items: userInfo.value.uantique.map(a => ({
      id: a.aiid, img: a.aiimg, icon: '🖼️', name: a.ainame,
      tag: a.airarity || null, price: formatPrice(a.aiprice || 0)
    }))
  })
  if (gameStore.userItemInsCount) groups.push({
    key: 'ins', icon: '📦', label: '国内商品', items: gameStore.userItemInsData.map(i => ({
      id: i.id, icon: '📦', name: i.iiname, tag: `x${i.num}`, price: formatPrice(i.iiprice)
    }))
  })
  if (gameStore.userItemOutCount) groups.push({
    key: 'out', icon: '🌍', label: '国外商品', items: gameStore.userItemOutData.map(i => ({
      id: i.id, icon: '🌍', name: i.iiname, tag: `x${i.num}`, price: formatPrice(i.iiprice)
    }))
  })
  if (stockList.value.length) groups.push({
    key: 'stock', icon: '📊', label: '股票', items: stockList.value.map(s => ({
      id: s.id, icon: '📊', name: s.name, tag: `${s.num}股`,
      price: (s.profit > 0 ? '+' : '') + formatPrice(s.profit),
      priceClass: s.profit < 0 ? 'danger' : s.profit > 0 ? 'success' : ''
    }))
  })
  if (companyList.value.length) groups.push({
    key: 'company', icon: '🚀', label: '公司', items: companyList.value.map(c => ({
      id: c.id, icon: '🚀', name: c.name, tag: `${c.num}期`, price: formatPrice(c.costprice)
    }))
  })
  if (datingList.value.length) groups.push({
    key: 'dating', icon: '💖', label: '约会', items: datingList.value.map(d => ({
      id: d.id, icon: getAffinityIcon(d.level), name: d.name, tag: `${d.level} · ${d.count}次`
    }))
  })
  return groups
})

const tabs = computed(() => {
  const t = [{ key: 'all', icon: '📋', label: '全部', count: 0 }]
  let total = 0
  for (const g of allGroups.value) {
    t.push({ key: g.key, icon: g.icon, label: g.label, count: g.items.length })
    total += g.items.length
  }
  t[0].count = total
  return t
})

const totalHoldings = computed(() => tabs.value[0]?.count || 0)
</script>

<style scoped>
.profile-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-top: 5px;
  min-height: 0;
}


/* ==================== HUD 顶部 ==================== */
.hud {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 5px 14px;
  background: var(--gradient-primary);
  border-radius: 12px;
}

.hud-avatar {
  display: flex;
  align-items: center;
  gap: 10px;
}

.hud-emoji {
  width: 48px;
  height: 48px;
  background: rgba(255,255,255,0.15);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.hud-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--font-white);
}

.hud-sub {
  font-size: 11px;
  color: rgba(255,255,255,0.6);
  margin-top: 1px;
}

.hud-bars {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.bar-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.bar-label { font-size: 11px; flex-shrink: 0; color: rgba(255,255,255,0.8); }

.bar-track {
  flex: 1;
  height: 5px;
  background: rgba(255,255,255,0.2);
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.bar-fill.hp { background: linear-gradient(to right, var(--success-color), var(--success-secondary-color)); }
.bar-fill.fame { background: linear-gradient(to right, var(--warning-color), var(--warning-secondary-color)); }

.bar-val { font-size: 10px; color: rgba(255,255,255,0.9); min-width: 50px; text-align: right; }

/* 机会次数区域 */
.hud-opportunity {
  display: flex;
  align-items: center;
  padding-left: 12px;
  border-left: 1px solid rgba(255,255,255,0.15);
  margin-left: 8px;
}

.opp-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 6px;
}

.opp-tag {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 5px 8px;
  background: rgba(255,255,255,0.1);
  border-radius: 6px;
  min-width: 44px;
}

.opp-icon {
  font-size: 16px;
  line-height: 1;
}

.opp-num {
  font-size: 11px;
  color: rgba(255,255,255,0.95);
  font-weight: 600;
  margin-top: 2px;
}
.bar-fill.work { background: linear-gradient(to right, var(--info-color), var(--info-secondary-color)); }

.bar-val {
  font-size: 11px;
  font-weight: 600;
  color: rgba(255,255,255,0.8);
  min-width: 45px;
  text-align: right;
}

/* ==================== 资产明细 ==================== */
.finance-row {
  display: flex;
  gap: 6px;
  margin-top: 5px;
  margin-bottom: 5px;
}

.finance-cell {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.fc-icon { font-size: 16px; flex-shrink: 0; }

.fc-text { display: flex; flex-direction: column; min-width: 0; }

.fc-label {
  font-size: 10px;
  color: var(--font-secondary);
}

.fc-val {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-color);
  white-space: nowrap;
}

.fc-val.danger { color: var(--error-color); }

/* ==================== Tab 切换 ==================== */
.holdings {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

/* 文件夹标签式 */
.holdings-tabs.folder {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.holdings-tabs.folder .tab-btn {
  padding: 6px 12px 8px 12px;
  font-size: 11px;
  color: var(--font-secondary);
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
  background: var(--panel-log-color);
  border: 1px solid var(--border-color);
  border-bottom: none;
  border-radius: 6px 6px 0 0;
  position: relative;
  z-index: 1;
}

.holdings-tabs.folder .tab-btn:hover {
  background: rgba(64, 158, 255, 0.1);
  color: var(--primary-color);
}

.holdings-tabs.folder .tab-btn.active {
  background: var(--panel-color);
  color: var(--primary-color);
  font-weight: 500;
  z-index: 2;
  border-bottom-color: var(--panel-color);
  margin-bottom: -1px;
}

.holdings-tabs.folder .tab-btn small {
  font-weight: 400;
  color: var(--font-light);
  margin-left: 3px;
}

.tab-content {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 0 6px 6px 6px;
  padding: 8px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  scrollbar-width: none;
}

.tab-content::-webkit-scrollbar {
  display: none;
}

.empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--font-light);
  font-size: 13px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

/* ==================== 物品卡片 ==================== */
.group {
  margin-bottom: 3px;
}

.group-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--font-secondary);
  margin-bottom: 6px;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
  gap: 5px;
}

.item-card {
  background: var(--panel-log-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: border-color 0.2s;
}

.item-card:hover {
  border-color: var(--primary-color);
}

.card-left {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-img {
  width: 36px;
  height: 36px;
  border-radius: 6px;
}

.card-icon {
  font-size: 24px;
}

.card-right {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.card-name {
  font-size: 11px;
  font-weight: 500;
  color: var(--font-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-tag {
  font-size: 10px;
  color: var(--font-secondary);
}

.card-price {
  font-size: 11px;
  font-weight: 600;
  color: var(--warning-color);
}

.card-price.danger { color: var(--error-color); }
.card-price.success { color: var(--success-color); }

/* ==================== 疾病情况 ==================== */
.disease-row {
  margin-top: 5px;
  padding: 6px 10px;
  background: linear-gradient(135deg, rgba(255, 245, 245, 0.5) 0%, rgba(255, 232, 232, 0.5) 100%);
  border: 1px solid rgba(255, 199, 199, 0.5);
  border-radius: 8px;
}
.disease-row.healthy {
  background: linear-gradient(135deg, rgba(240, 249, 235, 0.5) 0%, rgba(225, 243, 216, 0.5) 100%);
  border: 1px solid rgba(194, 231, 176, 0.5);
}
.disease-row.critical {
  background: color-mix(in srgb, var(--error-color) 12%, var(--panel-color));
  border-color: color-mix(in srgb, var(--error-color) 48%, transparent);
}

.disease-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--error-color);
  margin-bottom: 4px;
}
.disease-row.healthy .disease-title {
  color: var(--success-color);
}

.disease-content {
  display: flex;
  align-items: center;
  gap: 4px;
}
.health-icon { font-size: 14px; }
.health-text { font-size: 11px; color: var(--success-color); }
.health-text.critical,
.critical-health-text { font-size: 11px; color: var(--error-color); font-weight: 600; }

.disease-scroll {
  display: flex;
  gap: 4px;
  overflow-x: auto;
  scrollbar-width: none;
}
.disease-scroll::-webkit-scrollbar { display: none; }

.disease-card {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 1px 3px;
  background: var(--panel-color);
  border-radius: 3px;
  border-left: 2px solid var(--warning-color);
  flex-shrink: 0;
}
.dc-name { font-size: 11px; color: var(--font-color); }
.dc-sev { font-size: 10px; color: var(--font-secondary); }
</style>
