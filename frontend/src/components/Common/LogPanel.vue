<template>
  <div class="log-panel">
    <div class="section-header">
      <span class="section-title">{{ title }}</span>
      <div class="section-controls">
        <el-button size="small" text @click="$emit('clear')">清空</el-button>
      </div>
    </div>
    <div class="section-content">
      <div v-if="items.length === 0" class="section-empty">
        <span class="section-empty-icon">{{ emptyIcon }}</span>
        <span class="section-empty-text">{{ emptyText }}</span>
      </div>
      <div v-else class="log-list">
        <div v-for="item in items" :key="item.id" class="log-item">
          <span v-if="item.time" class="log-time">{{ item.time }}</span>
          <span v-if="item.action" class="log-action" :class="item.type">{{ item.action }}</span>
          <span class="log-item-text">{{ item.detail }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  items: {
    type: Array,
    default: () => [],
  },
  emptyIcon: {
    type: String,
    default: '📋',
  },
  emptyText: {
    type: String,
    default: '暂无记录',
  },
})

defineEmits(['clear'])
</script>
