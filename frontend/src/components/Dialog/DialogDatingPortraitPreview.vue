<template>
  <el-dialog
    v-model="visible"
    aria-label="约会对象写真"
    width="720px"
    align-center
    destroy-on-close
    :show-close="false"
    class="dating-portrait-preview-dialog"
    data-testid="dating-portrait-preview"
  >
    <div class="portrait-preview-shell">
      <div class="portrait-preview-stage" :style="{ '--scene-accent': scene.accent }">
        <img :src="scene.image" class="portrait-backdrop" alt="" />
        <div class="portrait-vignette"></div>
        <div class="portrait-overlay-actions" data-testid="dating-portrait-overlay-actions">
          <el-button
            circle
            size="small"
            type="primary"
            title="关闭写真"
            aria-label="关闭写真"
            data-testid="close-dating-portrait-preview"
            @click="visible = false"
          >×</el-button>
        </div>

        <div class="portrait-copy">
          <div class="portrait-kicker">{{ scene.icon }} {{ scene.label }}</div>
          <div class="portrait-name">{{ dating?.dname || '约会对象' }}</div>
          <div class="portrait-meta">
            {{ dating?.dnationality || '未知国籍' }} · {{ dating?.doccup || '未知职业' }} · {{ dating?.dage || '?' }}岁
          </div>
          <div v-if="dating?.ddesc" class="portrait-description">{{ dating.ddesc }}</div>
        </div>

        <div class="portrait-scene-hint" data-testid="dating-scene-hint">
          {{ scene.icon }} {{ scene.occupation }}
        </div>

        <div
          class="portrait-relationship"
          :class="{ established: (dating?.daffinitylevel || '陌生人') !== '陌生人' }"
          data-testid="dating-relationship-status"
        >
          <span class="portrait-relationship-dot"></span>
          <strong>{{ dating?.daffinitylevel || '陌生人' }}</strong>
        </div>

        <div class="portrait-character-wrap">
          <img
            v-if="dating?.dimage"
            :src="dating.dimage"
            :alt="dating.dname"
            class="portrait-character"
            data-testid="dating-portrait-character"
          />
          <div v-else class="portrait-fallback">{{ dating?.dsex ? '🕺' : '💃' }}</div>
        </div>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { getDatingCareerScene } from '@/src/utils/datingCareerScenes'

const props = defineProps({
  modelValue: Boolean,
  dating: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue'])
const visible = ref(props.modelValue)
const scene = computed(() => getDatingCareerScene(props.dating))

watch(() => props.modelValue, value => {
  visible.value = value
})
watch(visible, value => emit('update:modelValue', value))
</script>

<style scoped>
:global(.dating-portrait-preview-dialog.el-dialog) {
  padding: 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

:global(.dating-portrait-preview-dialog .el-dialog__header) {
  display: none;
}

:global(.dating-portrait-preview-dialog .el-dialog__body) {
  padding: 0;
  background: transparent;
}

.portrait-preview-shell {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.portrait-preview-stage {
  position: relative;
  overflow: hidden;
  aspect-ratio: 16 / 9;
  border: 4px solid #d5a43a;
  border-radius: 14px;
  background: #111827;
  box-shadow:
    0 0 0 1px #6b450f,
    0 0 0 3px #f1cc72,
    0 0 0 4px #593807,
    0 0 22px rgb(226 176 65 / 42%),
    0 24px 58px rgb(3 7 18 / 64%);
  isolation: isolate;
}

.portrait-preview-stage::after {
  position: absolute;
  inset: 5px;
  z-index: 6;
  border: 1px solid rgb(255 235 169 / 88%);
  border-radius: 9px;
  box-shadow:
    0 0 0 1px rgb(91 55 7 / 78%),
    inset 0 0 8px rgb(35 20 2 / 62%),
    0 0 7px rgb(245 204 111 / 28%);
  content: '';
  pointer-events: none;
}

.portrait-backdrop {
  position: absolute;
  inset: -3%;
  width: 106%;
  height: 106%;
  object-fit: cover;
  animation: portrait-scene-pan 9s ease-out both;
}

.portrait-vignette {
  position: absolute;
  inset: 0;
  z-index: 1;
  background:
    linear-gradient(90deg, rgb(5 11 24 / 78%) 0%, rgb(5 11 24 / 38%) 48%, transparent 72%),
    linear-gradient(0deg, rgb(5 11 24 / 68%), transparent 48%);
}

.portrait-overlay-actions {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 7;
  display: flex;
}

.portrait-overlay-actions :deep(.el-button) {
  width: 28px;
  min-width: 28px;
  height: 28px;
  margin-left: 0;
  padding: 0;
  border-color: color-mix(in srgb, var(--scene-accent) 66%, #fff);
  background: rgb(7 15 30 / 70%);
  color: #fff;
  box-shadow: 0 4px 14px rgb(0 0 0 / 28%);
  font-size: 15px;
  backdrop-filter: blur(7px);
}

.portrait-copy {
  position: absolute;
  top: 30px;
  left: 30px;
  z-index: 3;
  width: 46%;
  color: #fff;
  text-shadow: 0 2px 13px rgb(0 0 0 / 78%);
  animation: portrait-copy-enter 0.85s 0.2s ease-out both;
}

.portrait-kicker {
  margin-bottom: 6px;
  color: color-mix(in srgb, var(--scene-accent) 62%, #fff);
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.12em;
}

.portrait-name {
  font-size: 30px;
  font-weight: 760;
}

.portrait-meta {
  margin-top: 5px;
  color: rgb(255 255 255 / 88%);
  font-size: 13px;
}

.portrait-description {
  margin-top: 14px;
  color: rgb(255 255 255 / 86%);
  font-size: 11px;
  line-height: 1.55;
  text-shadow: 0 2px 9px rgb(0 0 0 / 86%);
}

.portrait-character-wrap {
  position: absolute;
  right: -1%;
  bottom: 0;
  z-index: 2;
  width: 58%;
  height: 98%;
  pointer-events: none;
  animation: portrait-character-enter 1s 0.35s cubic-bezier(0.2, 0.8, 0.2, 1) both;
}

.portrait-character,
.portrait-fallback {
  display: flex;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  object-fit: contain;
  object-position: center bottom;
  filter: drop-shadow(0 17px 18px rgb(0 0 0 / 45%));
  font-size: 100px;
}

.portrait-scene-hint {
  position: absolute;
  bottom: 18px;
  left: 28px;
  z-index: 4;
  color: rgb(255 255 255 / 82%);
  font-size: 11px;
  text-shadow: 0 2px 9px rgb(0 0 0 / 82%);
  pointer-events: none;
  animation: portrait-hint-lifecycle 4.6s ease both;
}

.portrait-relationship {
  position: absolute;
  right: 24px;
  bottom: 20px;
  z-index: 4;
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgb(255 255 255 / 82%);
  font-size: 11px;
  font-weight: 700;
  text-shadow: 0 2px 9px rgb(0 0 0 / 82%);
}

.portrait-relationship strong {
  font: inherit;
}

.portrait-relationship-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #94a3b8;
  box-shadow: 0 0 8px rgb(148 163 184 / 70%);
}

.portrait-relationship.established {
  color: #d1fae5;
}

.portrait-relationship.established .portrait-relationship-dot {
  background: #6ee7b7;
  box-shadow: 0 0 9px rgb(110 231 183 / 82%);
}

@keyframes portrait-scene-pan {
  from { transform: scale(1.01) translate3d(-0.8%, 0.3%, 0); }
  to { transform: scale(1.08) translate3d(1%, -0.5%, 0); }
}

@keyframes portrait-copy-enter {
  from { opacity: 0; transform: translateY(14px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes portrait-character-enter {
  from { opacity: 0; transform: translateX(10%) scale(0.96); }
  to { opacity: 1; transform: translateX(0) scale(1); }
}

@keyframes portrait-hint-lifecycle {
  0% { opacity: 0; transform: translateY(6px); }
  12%, 70% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(3px); }
}

@media (prefers-reduced-motion: reduce) {
  .portrait-preview-stage * {
    animation-duration: 1ms !important;
    animation-delay: 0ms !important;
  }
}
</style>
