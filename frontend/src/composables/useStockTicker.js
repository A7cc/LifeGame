import { onBeforeUnmount, onMounted } from 'vue'
import { UpdateStock } from '@/wailsjs/go/services/App.js'
import { useGameStore } from '@/src/stores/game'

// 行情轮询挂在 GameMain 上，而不是股市路由页面上。这样玩家浏览银行、医院等
// 页面时市场仍会持续变化，同时用 updating 防止慢请求发生重叠。
export function useStockTicker(intervalMs = 5000) {
  const gameStore = useGameStore()
  let timer = null
  let updating = false

  const refresh = async () => {
    if (updating || gameStore.stockMarketClosed || !Array.isArray(gameStore.gameInfo?.gstockinfo) || gameStore.gameInfo.gstockinfo.length === 0) {
      return
    }

    updating = true
    try {
      const result = await UpdateStock()
      if (result?.code === 200) {
        gameStore.applyStockUpdate(result)
      }
    } catch (error) {
      console.error('更新股票行情失败', error)
    } finally {
      updating = false
    }
  }

  onMounted(() => {
    timer = window.setInterval(refresh, intervalMs)
  })

  onBeforeUnmount(() => {
    if (timer !== null) {
      window.clearInterval(timer)
      timer = null
    }
  })

  return { refreshStockMarket: refresh }
}
