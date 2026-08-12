<template>
  <MarketTradePanel
    test-id="page-foreign"
    title="国外市场"
    region="国外市场"
    :market-items="gameStore.gameInfo.gitemoutinfo || {}"
    :user-items="gameStore.userItemOutData"
    :user-item-count="gameStore.userItemOutCount"
    :user-item-limit="gameStore.gameInfo.gmaxholdnum.mfholdnum"
  />
</template>

<script setup>
import { useGameStore } from '@/src/stores/game'
import MarketTradePanel from '@/src/components/Menu/MarketTradePanel.vue'

const gameStore = useGameStore()
</script>
