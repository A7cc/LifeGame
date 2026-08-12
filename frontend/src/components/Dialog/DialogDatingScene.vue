<template>
  <el-dialog
    v-model="visible"
    :title="`与${dating?.dname || ''}约会`"
    width="520px"
    align-center
    class="dating-scene-dialog"
  >
    <div v-if="dating" class="dating-content">
      <div class="property-info">
        <div class="property-details">
          <div class="detail-name">{{ dating.dname }}</div>
          <span class="detail-meta">{{ dating.dnationality }} · {{ dating.doccup }} · {{ dating.dage }}岁 · <span class="affinity-value">{{ dating.daffinitylevel }}</span></span>
          <span class="detail-desc">{{ dating.ddesc }}</span>
        </div>
        <div class="detail-action">
          <span class="cost-value">💰{{ parseInt(dating.dcost).toLocaleString() }}</span>
          <el-button size="small" type="warning" @click="handleConfirm" :disabled="selectedLocation === -1">💑 约会</el-button>
        </div>
      </div>

      <div class="property-scene-section">
        <div class="section-header">
          <div>
            <div class="section-title">🎯 猜猜 {{ dating.dname }} 喜欢的约会环境</div>
            <div class="section-hint">偏好场景成功率高、收益适中；非偏好场景成功率低，但成功回报更高</div>
          </div>
        </div>

        <div class="scene-list" data-testid="dating-scene-options">
          <button
            v-for="(location, index) in sceneOptions"
            :key="location"
            type="button"
            class="scene-item"
            :class="{ selected: selectedLocation === index }"
            :data-testid="`dating-location-${location}`"
            @click="selectedLocation = index"
          >
            <span class="scene-name">{{ location }}</span>
            <span v-if="selectedLocation === index" class="scene-check">✓</span>
          </button>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { createDatingSceneOptions } from '@/src/utils/datingScenes'

const props = defineProps({
  modelValue: Boolean,
  dating: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = ref(props.modelValue)
const sceneOptions = ref([])
const selectedLocation = ref(-1)

const refreshSceneOptions = () => {
  sceneOptions.value = createDatingSceneOptions(props.dating?.dlocations)
  selectedLocation.value = -1
}

watch(() => props.modelValue, val => {
  visible.value = val
  if (val) refreshSceneOptions()
}, { immediate: true })
watch(visible, val => emit('update:modelValue', val))
watch(() => props.dating?.did, () => {
  if (visible.value) refreshSceneOptions()
})

// 确认约会
const handleConfirm = () => {
  if (selectedLocation.value === -1) return
  emit('confirm', {
    dating: props.dating,
    locationIndex: selectedLocation.value,
    location: sceneOptions.value[selectedLocation.value]
  })
  visible.value = false
}
</script>

<style scoped>
.property-info {
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 10px;
  margin-bottom: 8px;
}

.property-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
}

.detail-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--font-color);
}

.detail-meta {
  font-size: 10px;
  color: var(--font-secondary);
}

.detail-desc {
  font-size: 10px;
  color: var(--font-secondary);
}

.detail-action {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  flex-shrink: 0;
}

.cost-value {
  font-size: 11px;
  font-weight: 500;
  color: var(--warning-color);
}

.affinity-value {
  font-size: 10px;
  font-weight: 600;
  color: var(--font-purple);
}

.property-scene-section {
  padding: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--font-secondary);
}

.section-hint {
  margin-top: 2px;
  font-size: 11px;
  color: var(--font-secondary);
}

.scene-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
  gap: 5px;
}

.scene-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 5px;
  min-height: 48px;
  padding: 8px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  color: var(--font-color);
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.scene-item.selected {
  border-color: var(--success-color);
  box-shadow: inset 0 0 0 1px var(--success-color);
}

.scene-name {
  flex: 1;
  font-size: 13px;
  color: inherit;
  font-weight: 600;
}

.scene-check {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  color: var(--success-color);
  font-weight: bold;
  z-index: 1;
}
</style>
