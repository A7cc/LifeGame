<template>
  <el-dialog
    data-testid="evaluation-dialog"
    :model-value="modelValue"
    width="600px"
    :close-on-click-modal="false"
    :show-close="false"
    class="evaluation-dialog"
    destroy-on-close
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #header>
      <div class="dialog-header">
        <span class="header-icon">🎮</span>
        <span class="header-title">游戏结束</span>
      </div>
    </template>

    <div v-if="evaluation" class="evaluation-content">
      <div class="top-section">
        <div class="score-section">
            <div class="score-circle">
            <div class="score-number">{{ evaluation.score }}</div>
            <div class="score-label">人生评分</div>
            </div>
            <div class="score-title">{{ evaluation.title }}</div>
            <el-button data-testid="restart-game" type="primary" class="restart-btn-small" @click="emit('restart')">
                <span>重新开始</span>
                <span class="btn-icon">🔄</span>
                </el-button>
        </div>

        <div class="evaluation-details">
            <div v-for="detail in scoreDetails" :key="detail.key" class="detail-item">
            <div class="detail-header">
                <span class="detail-icon">{{ detail.icon }}</span>
                <span class="detail-label">{{ detail.label }}</span>
                <span class="detail-score">{{ detail.score }}/{{ detail.max }}</span>
            </div>
            <div class="detail-progress">
                <div
                class="progress-bar"
                :style="{
                    width: Math.round((detail.score / detail.max) * 100) + '%',
                    background: getScoreGradient(detail.score, detail.max),
                }"
                ></div>
            </div>
            </div>
        </div>
      </div>

      <div class="evaluation-description">{{ evaluation.description }}</div>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  evaluation: { type: Object, default: null },
})

const emit = defineEmits(['update:modelValue', 'restart'])

const scoreDetails = computed(() => {
  const evaluation = props.evaluation || {}
  return [
    { key: 'wealth', icon: '💰', label: '财富评分', score: evaluation.wealthscore || 0, max: 25 },
    { key: 'health', icon: '💚', label: '健康评分', score: evaluation.healthscore || 0, max: 20 },
    { key: 'fame', icon: '⭐', label: '名声评分', score: evaluation.famescore || 0, max: 15 },
    { key: 'career', icon: '💼', label: '事业评分', score: evaluation.careerscore || 0, max: 10 },
    { key: 'relationship', icon: '💖', label: '关系评分', score: evaluation.relationshipscore || 0, max: 10 },
    { key: 'collection', icon: '🏡', label: '收藏评分', score: evaluation.collectionscore || 0, max: 5 },
  ]
})

const getScoreGradient = (score, max) => {
  const percent = (score / max) * 100
  if (percent >= 80) return 'linear-gradient(90deg, #67c23a, #85ce61)'
  if (percent >= 60) return 'linear-gradient(90deg, #409eff, #66b1ff)'
  if (percent >= 40) return 'linear-gradient(90deg, #e6a23c, #f0c78a)'
  return 'linear-gradient(90deg, #f56c6c, #f89898)'
}
</script>

<style scoped>
:global(.evaluation-dialog.el-dialog) {
  border: 2px solid var(--primary-color);
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 18px 54px rgb(var(--primary-color-rgb) / 24%);
}

:global(.evaluation-dialog .el-dialog__header),
:global(.evaluation-dialog .el-dialog__body),
:global(.evaluation-dialog .el-dialog__footer) {
  padding: 0;
  margin: 0;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 18px 24px;
  border-radius: 18px;
  background: var(--gradient-primary);
}

.header-icon {
  font-size: 26px;
}

.header-title {
  font-size: 19px;
  font-weight: 600;
  color: white;
}

.evaluation-content {
  padding: 18px 24px;
  background: var(--panel-main-color);
}

.top-section {
  display: flex;
  gap: 24px;
  align-items: center;
}

.score-section {
  flex-shrink: 0;
  text-align: center;
  width: 130px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.score-circle {
  width: 85px;
  height: 85px;
  border-radius: 50%;
  background: var(--gradient-primary);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: 0 auto 10px;
  box-shadow: 0 8px 30px rgb(var(--primary-color-rgb) / 30%);
}

.score-number {
  font-size: 36px;
  font-weight: 800;
  color: white;
  line-height: 1;
}

.score-label {
  font-size: 11px;
  color: rgb(255 255 255 / 90%);
}

.score-title {
  font-size: 15px;
  font-weight: 700;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 14px;
}

.restart-btn-small {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.restart-btn-small:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 18px rgba(102, 126, 234, 0.4);
}

.restart-btn-small .btn-icon {
  margin-left: 4px;
  font-size: 13px;
  transition: transform 0.5s ease;
}

.restart-btn-small:hover .btn-icon { transform: rotate(180deg); }

.evaluation-details {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.detail-item {
  padding: 12px 14px;
  background: var(--panel-color);
  border-radius: 14px;
  border: 1px solid var(--primary-color);
  box-shadow: 0 2px 8px var(--panel-shadow);
  transition: all 0.2s ease;
}

.detail-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.detail-icon {
  font-size: 18px;
}

.detail-label {
  flex: 1;
  font-size: 13px;
  font-weight: 600;
  color: var(--font-color);
}

.detail-score {
  font-size: 12px;
  font-weight: 700;
  color: var(--font-secondary);
}

.detail-progress {
  width: 100%;
  height: 7px;
  background: var(--panel-sub-color);
  border-radius: 10px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  border-radius: 10px;
  transition: width 0.8s ease-out;
}

.evaluation-description {
  padding: 14px 16px;
  background: var(--accent-color);
  border-radius: 14px;
  font-size: 13px;
  line-height: 1.7;
  color: var(--accent-strong-color);
  white-space: pre-line;
  margin-top: 14px;
  border: 1px solid var(--warning-color);
}
</style>
