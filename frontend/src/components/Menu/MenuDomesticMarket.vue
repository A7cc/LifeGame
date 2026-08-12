<template>
  <MarketTradePanel
    test-id="page-domestic"
    title="国内市场"
    region="国内市场"
    :market-items="gameStore.gameInfo.giteminsinfo || {}"
    :user-items="gameStore.userItemInsData"
    :user-item-count="gameStore.userItemInsCount"
    :user-item-limit="gameStore.gameInfo.gmaxholdnum.mdholdnum"
  />
</template>

<script setup>
import { useGameStore } from '@/src/stores/game'
import MarketTradePanel from '@/src/components/Menu/MarketTradePanel.vue'

const gameStore = useGameStore()
</script>
