<template>
  <el-dialog
    data-testid="load-dialog"
    :model-value="modelValue"
    title="加载存档"
    width="400px"
    align-center
    class="load-dialog"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <div class="load-content" v-loading="loading">
      <div v-if="saves.length > 0" class="save-grid">
        <button
          v-for="save in saves"
          :key="save.id"
          type="button"
          :data-testid="`load-save-${save.id}`"
          class="save-card"
          @click="emit('load', save.id)"
        >
          <span class="save-card-name">{{ save.name }}</span>
          <span class="save-card-info">
            <span>{{ save.game_year }}岁</span>
            <span>{{ save.created_at }}</span>
          </span>
        </button>
      </div>
      <div v-else class="no-saves-tip">暂无存档，请先开始新游戏</div>
    </div>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
defineProps({
  modelValue: { type: Boolean, default: false },
  saves: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'load'])
</script>

<style scoped>
.load-dialog :deep(.el-dialog) {
  border-radius: 16px;
}

.load-content {
  min-height: 150px;
}

.save-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.save-card {
  display: flex;
  flex-direction: column;
  width: 100%;
  padding: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #dcdfe6;
  text-align: left;
  font: inherit;
}

.save-card:hover,
.save-card:focus-visible {
  background: linear-gradient(135deg, #ecf5ff 0%, #d9ecff 100%);
  border-color: #409eff;
  transform: translateY(-2px);
  outline: none;
}

.save-card-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.save-card-info {
  font-size: 12px;
  color: #909399;
  display: flex;
  gap: 12px;
}

.no-saves-tip {
  text-align: center;
  color: #909399;
  padding: 40px;
  font-size: 14px;
}
</style>
