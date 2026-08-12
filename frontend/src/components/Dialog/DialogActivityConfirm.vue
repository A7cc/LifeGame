<template>
  <el-dialog
    data-testid="activity-confirm-dialog"
    v-model="visible"
    :title="`参与${activity?.name || ''}`"
    width="300px"
    align-center
    class="activity-confirm-dialog"
  >
    <div v-if="activity" class="activity-content">
      <!-- 活动信息 -->
      <div class="activity-info">
        <div class="activity-icon-large">{{ activity.icon }}</div>
        <div class="activity-details">
          <div class="activity-name">{{ activity.name }}</div>
          <div class="activity-desc">{{ activity.desc }}</div>
        </div>
        <div class="activity-action">
          <span class="cost-value">💰 {{ activity.entryCost?.toLocaleString() }}</span>
          <el-button data-testid="confirm-activity" type="primary" size="small" @click="handleConfirm">确认参与</el-button>
        </div>
      </div>

      <!-- 效果列表 -->
      <div class="effect-section">
        <div class="section-header">
          <div class="section-title">🎯 活动效果</div>
          <div class="section-controls">
            <el-tag size="small">效果浮动 ±30%</el-tag>
          </div>
        </div>
        <div class="effect-list">
          <div class="effect-row">
            <span class="effect-name">健康</span>
            <span class="effect-value">+{{ activity.healthGain }}</span>
          </div>
          <div v-if="activity.fameGain > 0" class="effect-row">
            <span class="effect-name">名声</span>
            <span class="effect-value">+{{ activity.fameGain }}</span>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  activity: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

// 对话框可见性
const visible = ref(props.modelValue)
watch(() => props.modelValue, val => visible.value = val)
watch(visible, val => emit('update:modelValue', val))

// 确认
const handleConfirm = () => {
  emit('confirm', props.activity)
  visible.value = false
}
</script>

<style scoped>
/* 活动信息 */
.activity-info {
  display: flex;
  gap: 10px;
  border-radius: 10px;
  margin-bottom: 5px;
}

.activity-icon-large {
  font-size: 30px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.activity-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
}

.activity-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
}

.activity-desc {
  font-size: 10px;
  color: var(--font-secondary);
  line-height: 1.4;
}

.activity-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
}

.cost-value {
  font-size: 11px;
  font-weight: 500;
  color: var(--warning-color);
}

/* 效果区域 */
.effect-section {
  padding: 0 4px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-secondary);
}

.effect-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 5px;
}

.effect-row {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 5px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 3px;
  cursor: pointer;
  position: relative;
}

.effect-name {
  flex: 1;
  font-size: 12px;
  color: var(--font-light);
  font-weight: 500;
}

.effect-value {
  font-size: 16px;
  color: var(--success-color);
  font-weight: bold;
}
</style>
