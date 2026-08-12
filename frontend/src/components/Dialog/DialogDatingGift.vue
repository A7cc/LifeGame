<template>
  <el-dialog
    v-model="visible"
    data-testid="dating-gift-dialog"
    :title="`送给${dating?.dname || ''}礼物`"
    width="520px"
    align-center
    class="dating-gift-dialog"
  >
    <div v-if="dating" class="dating-content">
      <div class="property-info">
        <div class="property-details">
          <div class="detail-name">{{ dating.dname }}</div>
          <span class="detail-meta">
            {{ dating.dnationality }} · {{ dating.doccup }} · {{ dating.dage }}岁 ·
            <span class="affinity-value">{{ dating.daffinitylevel }}</span>
          </span>
          <span class="detail-desc">{{ dating.ddesc }}</span>
        </div>
        <div class="detail-action">
          <span class="cost-value">{{ selectedOption ? `💰${giftCost.toLocaleString()}` : '选择礼物' }}</span>
          <el-button
            data-testid="confirm-dating-gift"
            size="small"
            type="warning"
            :loading="pending"
            :disabled="selectedGift === -1 || pending"
            @click="handleConfirm"
          >🎁 送礼</el-button>
        </div>
      </div>

      <div class="property-gift-section">
        <div class="section-header">
          <div>
            <div class="section-title">🎯 选择 {{ dating.dname }} 喜欢的礼物</div>
            <div class="section-hint">偏好礼物更容易送出、收益适中；非偏好礼物很难送出，但成功回报更高</div>
          </div>
        </div>

        <div v-loading="loading" class="gift-list" data-testid="dating-gift-options">
          <button
            v-for="(gift, index) in options"
            :key="gift.name"
            type="button"
            class="gift-item"
            :class="{ selected: selectedGift === index }"
            :data-testid="`gift-option-${dating.did}-${index}`"
            :disabled="pending"
            @click="selectedGift = index"
          >
            <span class="gift-name">{{ gift.name }}</span>
            <span v-if="selectedGift === index" class="gift-check">✓</span>
          </button>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  dating: { type: Object, default: null },
  pending: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  options: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = ref(props.modelValue)
const selectedGift = ref(-1)
const selectedOption = computed(() => props.options[selectedGift.value] || null)
const giftCost = computed(() => Number(selectedOption.value?.cost || 0))

watch(() => props.modelValue, value => {
  visible.value = value
  if (value) selectedGift.value = -1
}, { immediate: true })

watch(visible, value => emit('update:modelValue', value))

watch(() => props.dating?.did, () => {
  if (visible.value) selectedGift.value = -1
})

watch(() => props.options, () => {
  selectedGift.value = -1
})

const handleConfirm = () => {
  if (selectedGift.value === -1 || props.pending) return
  emit('confirm', selectedOption.value)
}
</script>

<style scoped>
.property-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  border-radius: 10px;
}

.property-details {
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  min-width: 0;
}

.detail-name {
  color: var(--font-color);
  font-size: 16px;
  font-weight: 600;
}

.detail-meta,
.detail-desc {
  overflow: hidden;
  color: var(--font-secondary);
  font-size: 10px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-action {
  display: flex;
  flex-shrink: 0;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.cost-value {
  color: var(--warning-color);
  font-size: 11px;
  font-weight: 500;
}

.affinity-value {
  color: var(--font-purple);
  font-size: 10px;
  font-weight: 600;
}

.property-gift-section {
  padding: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.section-title {
  color: var(--font-secondary);
  font-size: 13px;
  font-weight: 600;
}

.section-hint {
  margin-top: 2px;
  color: var(--font-secondary);
  font-size: 11px;
}

.gift-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
  gap: 5px;
}

.gift-item {
  position: relative;
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  gap: 5px;
  padding: 8px 12px;
  color: var(--font-color);
  font: inherit;
  text-align: left;
  cursor: pointer;
  border: 1px solid var(--border-color);
  border-radius: 5px;
  background: var(--panel-color);
}

.gift-item.selected {
  border-color: var(--success-color);
  box-shadow: inset 0 0 0 1px var(--success-color);
}

.gift-item:disabled {
  cursor: wait;
  opacity: 0.65;
}

.gift-name {
  width: 100%;
  color: inherit;
  font-size: 13px;
  font-weight: 600;
}

.gift-check {
  position: absolute;
  top: 50%;
  right: 6px;
  z-index: 1;
  color: var(--success-color);
  font-size: 16px;
  font-weight: bold;
  transform: translateY(-50%);
}
</style>
