<template>
  <div class="product-panel" data-testid="page-stock">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">股市</div>
      <div class="panel-controls">
        <el-tag :type="gameStore.stockMarketClosed ? 'info' : 'success'" size="small">
          {{ gameStore.stockMarketClosed ? '本年已收盘' : `本年剩余 ${gameStore.stockRemaining} 次行情` }}
        </el-tag>
        <div class="panel-btn-group">
          <el-radio-group v-model="chartType" size="small">
            <el-radio-button label="line">分时图</el-radio-button>
            <el-radio-button label="kline">日K图</el-radio-button>
          </el-radio-group>
        </div>
        <div class="panel-btn-group">
          <el-select data-testid="stock-select" v-model="selectedStockId" size="small" placeholder="选择股票" style="width: 150px">
            <el-option v-for="s in gameStore.gameInfo.gstockinfo" :key="s.siid" :label="s.siname" :value="s.siid" />
          </el-select>
        </div>
      </div>
    </div>

    <!-- 中间信息区域 -->
    <div class="main-content">
      <div class="chart-content">
        <div ref="chartRef" style="height: 100%; width: 100%"></div>
      </div>
      <div class="log-panel">
        <div class="section-header">
          <span class="section-title">📋 市场动态</span>
          <div class="section-controls">
              <el-button size="small" text @click="runLogInfo = []">清空</el-button>
            </div>
        </div>
        <div class="section-content">
          <div v-if="runLogInfo.length === 0" class="section-empty">
            <span class="section-empty-icon">📊</span>
            <span class="section-empty-text">暂无记录</span>
          </div>
          <div v-else class="log-list">
            <div v-for="loginfo in runLogInfo.slice(0, 30)" class="log-item">
              <span class="log-item-icon">{{ loginfo.includes('📈') ? '📈' : loginfo.includes('📉') ? '📉' : '📑' }}</span>
              <span class="log-item-text">{{ loginfo.replace(/📈\s?|📉\s?/g, '') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部信息区域 -->
    <div class="bottom-info">
      <!-- 交易记录 -->
      <div class="log">
        <div class="section-header">
          <span class="section-title">交易记录</span>
        </div>

        <div class="log-table">
          <el-table :data="tradeHistory.slice().reverse()" size="small" border height="100%">
            <el-table-column prop="time" label="时间"/>
            <el-table-column prop="type" label="操作" width="70" />
            <el-table-column prop="stock" label="股票" width="70" />
            <el-table-column prop="amount" label="数量" width="60" />
            <el-table-column prop="price" label="价格" width="70" />
          </el-table>
        </div>
      </div>

      <!-- 用户操作区域 -->
      <div class="uactions">
        <div class="section-header">
          <span class="section-title">账户信息</span>
        </div>

        <el-descriptions border :column="2" size="small">
          <el-descriptions-item label="当前价格" width="90" >
            ￥{{ selectedStock?.siprice }}
            <span v-if="selectedStock?.sistatus === '涨停'" style="color: red;">（涨停）</span>
              <span v-else-if="selectedStock?.sistatus === '跌停'" style="color: blue;">（跌停）</span>
          </el-descriptions-item>
          <el-descriptions-item label="涨跌幅" width="70" >
            <span :style="{ color: percentageChange > 0 ? 'red' : percentageChange < 0 ? 'green' : '#333' }">
              {{ percentageChange }}%
            </span>
        </el-descriptions-item>
          <el-descriptions-item label="持有股票" width="90" ><span data-testid="stock-holding">{{ gameStore.userInfo.ustock[selectedStock?.siid]?.usnum || 0 }} 股</span></el-descriptions-item>
          <el-descriptions-item label="持仓成本" width="70" >￥{{ gameStore.userInfo.ustock[selectedStock?.siid]?.usprice_init || 0 }}</el-descriptions-item>
          <el-descriptions-item label="今年盈亏" width="90" >
            <span :style="{ color: gameStore.userInfo.ustockprofit < 0 ? 'green' :  gameStore.userInfo.ustockprofit > 0 ? 'red' : '' }">
              ￥{{  gameStore.userInfo.ustockprofit || 0 }}
            </span></el-descriptions-item>
          <el-descriptions-item label="累计盈亏" width="70" >
            <span :style="{ color: gameStore.userInfo.ustock[selectedStock?.siid]?.usprofit < 0 ? 'green' : gameStore.userInfo.ustock[selectedStock?.siid]?.usprofit > 0 ? 'red' : '' }">
              ￥{{ gameStore.userInfo.ustock[selectedStock?.siid]?.usprofit || 0 }}
            </span>
          </el-descriptions-item>
        </el-descriptions>

        <div class="trade-buttons">
          <el-button data-testid="stock-buy" type="primary" class="full-btn" size="small" @click="tradeStock('购买')">买入</el-button>
          <el-input-number data-testid="stock-amount" size="small" v-model="amount" :min="10" :max="9999" />
          <el-button data-testid="stock-sell" type="danger" class="full-btn" size="small" @click="tradeStock('出售')">卖出</el-button>
        </div>
        <div class="trade-buttons">
          <el-button data-testid="stock-buy-all" type="primary" class="full-btn" size="small" @click="tradeStock('全仓购买')">全仓买入</el-button>
          <el-button data-testid="stock-sell-all" type="danger" class="full-btn" size="small" @click="tradeStock('全仓出售')">全仓卖出</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useGameStore } from '@/src/stores/game'
import { ElMessage } from 'element-plus'
import { BuyItem, SellItem } from "@/wailsjs/go/services/App.js";

// ECharts 按需引入，避免股市页把完整图表库打进业务 chunk。
import { init, use } from 'echarts/core'
import { LineChart, CandlestickChart } from 'echarts/charts'
import { GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, CandlestickChart, GridComponent, CanvasRenderer])
// 获取游戏数据
const gameStore = useGameStore()

// 存储所有的交易历史记录
const tradeHistory = ref([])
// 存储随机的新闻事件，这些事件可能影响股票价格
const runLogInfo = ref([])
// 绑定图表容器 DOM
const chartRef = ref(null)
// 保存 ECharts 实例，或者说数据
const chart = ref(null)
// 当前图表类型，默认是“分时图”
const chartType = ref('line')
// 避免短时间内频繁触发 chart.resize() 导致卡顿或抖动
let observer = null
// 设置数量
const amount = ref(1)
// 统一图表渲染入口
const renderChart = () => {
  if (!chart.value && chartRef.value) {
    chart.value = init(chartRef.value)
  }
  if (!chart.value) return // 防止报错
  const stockdata = selectedStock.value
  if (!stockdata?.sihistory) return
  // 通过选择chartType来决定渲染哪种图表
  if (chartType.value === 'line') {
    // 直线图
    renderLineChart(stockdata)
  } else if (chartType.value === 'kline') {
    // K线图
    renderKLineChart(stockdata)
  }
}

// 渲染分时图函数
const renderLineChart = (stockdata) => {
  const history = stockdata.sihistory.slice(-30) // 只取最近30条
  chart.value.setOption({
    grid: {
      top: 20,     // ⭐ 减少顶部空白
      left: 40,
      right: 20,
      bottom: 20   // 可保留底部用于刻度
    },
    xAxis: {
      type: 'category',
      data: history.map((_, i) => i + 1),
    },
    yAxis: { type: 'value' },
    series: [{
      data: history,
      type: 'line',
      // areaStyle是开启“面积图”模式，会在折线图下方填充一块颜色区域。
      areaStyle: {
        // color: 'rgba(0, 136, 212, 0.2)'  // 可加颜色和透明度
      },
    }]
  })
}

// 渲染日K线图函数
const renderKLineChart = (stockdata) => {
  const klineHistory = stockdata.siklinehistory.slice(-30) // 只取最近30条
  chart.value.setOption({
    xAxis: {
      type: 'category',
      data: klineHistory.map((_, i) => `Day ${i + 1}`),
    },
    yAxis: { type: 'value' },
    series: [{
      type: 'candlestick',
      data: klineHistory,
      itemStyle: {
        color: '#ec0000',
        color0: '#00da3c',
        borderColor: '#8A0000',
        borderColor0: '#008F28'
      },
    }]
  })
}

// 用于保存当前选中的股票的名字，默认是科技股
const selectedStockId = ref(1)
// 使用 computed 是为了在依赖的值（比如 stock.selectedStockId）发生变化时，自动重新计算 selectedStock 的值
const selectedStock = computed(() => {
  return gameStore.gameInfo.gstockinfo.find(s => s.siid === selectedStockId.value)
})
// 用于计算当前选中股票的涨跌幅
const percentageChange = computed(() => {
  // 从当前选中的股票对象中，获取其价格历史记录数组（sihistory）
  const h = selectedStock.value?.sihistory
  if (!h || h.length < 2) return 0
  const change = ((h[h.length - 1] - h[h.length - 2]) / h[h.length - 2]) * 100
  return change.toFixed(2)
})

// 买卖股票
const tradeStock = async (type) => {
  const stocktmp = selectedStock.value
  if (!stocktmp) return

  // 检查全仓出售时是否持有股票
  if (type === '全仓出售') {
    const holding = gameStore.userInfo.ustock[stocktmp.siid]?.usnum || 0
    if (holding <= 0) {
      ElMessage.warning('你没有持有该股票')
      return
    }
  }

  try{
    let data = null
    // 判断当前是购买还是出售
    if (type === '购买') {
      data = await BuyItem(stocktmp.siid, amount.value, "股票")
    } else if (type === '全仓购买'){
      const buyAmount = Math.floor(gameStore.userInfo?.ucash / stocktmp.siprice)
      if (buyAmount < 10) {
        ElMessage.warning('资金不足以购买最少10股')
        return
      }
      data = await BuyItem(stocktmp.siid, buyAmount, "股票")
    } else if (type === '全仓出售'){
      data = await SellItem(stocktmp.siid, gameStore.userInfo.ustock[stocktmp.siid].usnum, "股票")
    } else {
      data = await SellItem(stocktmp.siid, amount.value, "股票")
    }
    if (data.code == 200) {
      gameStore.applyUserInfo(data.userinfo)
      // 加上数据信息
      tradeHistory.value.push({ type: type, stock: stocktmp.siname, amount: amount.value, price: stocktmp.siprice, time: new Date().toLocaleTimeString() })
      // 弹窗提示
      ElMessage.success(`已${type} ${amount.value} 个${stocktmp.siname}`)
    } else {
      ElMessage.error(data.msg || `${type}失败`);
    }
  }catch (err) {
      console.error('调用 BuyItem/SellItem 异常：', err)
      ElMessage.error(`${type}失败`);
  }
}

// 生命周期钩子
// 挂载后初始化图表并绑定窗口大小监听
onMounted(() => {
  // 默认选择第一个股票
  if (gameStore.gameInfo.gstockinfo.length > 0) {
    selectedStockId.value = gameStore.gameInfo.gstockinfo[0].siid
  }
  runLogInfo.value = gameStore.gameInfo.gstocknews || []
  // 用于监听一个 DOM 元素的尺寸变化（比如宽高变化），常用于响应式布局、图表重绘、组件自适应等场景
  renderChart()
  observer = new ResizeObserver(() => {
    chart.value?.resize()
  })
  if (chartRef.value) {
    observer.observe(chartRef.value)
  }
})
// 卸载前释放图表资源
onBeforeUnmount(() => {
  if (observer && chartRef.value) {
    observer.unobserve(chartRef.value)
  }
  chart.value?.dispose()
  chart.value = null
})
// 当 chartType.value 的值发生变化时，调用 renderChart 函数重新渲染图表
watch(() => chartType.value, renderChart)
// 当 stock.selectedStockId（选中的股票 ID）发生变化时，调用 renderChart，说明切换股票时需要更新图表
watch(() => selectedStockId.value, renderChart)
// 全局行情轮询替换 gameInfo 后，股市页只负责响应数据并重绘。
watch(() => gameStore.gameInfo.gstockinfo, renderChart, { flush: 'post' })
watch(() => gameStore.gameInfo.gstocknews, (news) => {
  runLogInfo.value = news || []
})
</script>

<style scoped>

.panel-btn-group {
  display: flex;
  align-items: center;
}

/* 图表 + 行情区域 */
.main-content {
  display: flex;  /* ✅ 使子元素水平排列 */
  flex: 1; /* ✅ 占据剩余空间 */
  gap: 5px; /* ✅ 子元素之间有5px的间距 */
  margin-top: 5px;   /* ✅ 子元素之间有5px的间距 */
  overflow: auto;       /* ✅ 避免内容撑开，高度超出时滚动 */
}
.chart-content {
  flex: 2; /* ✅ 占据一半宽度 */
  display: flex; /* ✅ 使子元素垂直排列 */
  justify-content: center; /* ✅ 子元素水平居中 */
  align-items: center; /* ✅ 子元素垂直居中 */
}
.log-panel {
  flex: 1; /* ✅ 占据一半宽度 */
}

/* 底部信息区：交易记录 + 用户操作 */
.bottom-info {
  display: flex;
  gap: 5px;
  margin-top: 5px;
  height: 220px;
}

.log, .uactions {
  flex: 1; /* ✅ 平分宽度 */
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--panel-color);
  display: flex;
  flex-direction: column;
  padding: 8px;
  overflow: hidden;
}

.log-table {
  flex: 1;
  overflow-y: auto;
}
/* 交易按钮 */
.trade-buttons {
  display: flex;
  gap: 5px; /* 按钮之间的间距 */
  margin-top: 10px;
}

.full-btn {
  flex: 1; /* 占据父容器剩余空间 */
}

.el-timeline {
  display: flex;
  flex: 1;
  flex-direction: column;
  height: 80%;
  /* 移除滚动条预留空间 */
  scrollbar-width: none;                         /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;                      /* IE/Edge 浏览器隐藏滚动条 */
}
.el-timeline::-webkit-scrollbar {
  display: none;                                 /* Chrome/Safari 隐藏滚动条 */
}
</style>
