<template>
  <el-dialog
    v-model="visible"
    aria-label="房车场景"
    width="720px"
    align-center
    destroy-on-close
    :show-close="false"
    class="asset-purchase-moment-dialog"
    data-testid="asset-purchase-moment"
  >
    <div class="purchase-moment-shell">
      <div
        :key="playKey"
        class="purchase-stage"
        :class="[`purchase-${kind}`, { 'asset-showcase-stage': isPreview }]"
        data-testid="asset-purchase-stage"
      >
        <img :src="meta.backdrop" class="purchase-backdrop" alt="" />
        <div class="purchase-vignette"></div>
        <div class="purchase-overlay-actions" data-testid="asset-overlay-actions">
          <el-button
            circle
            size="small"
            type="primary"
            :title="isPreview ? '关闭鉴赏' : `收下${kind === 'car' ? '钥匙' : '新家'}`"
            :aria-label="isPreview ? '关闭鉴赏' : `收下${kind === 'car' ? '钥匙' : '新家'}`"
            data-testid="close-purchase-moment"
            @click="visible = false"
          >{{ isPreview ? '×' : '✓' }}</el-button>
          <el-button
            circle
            size="small"
            title="重播"
            aria-label="重播"
            data-testid="replay-purchase-moment"
            @click="restart"
          >↻</el-button>
        </div>
        <img
          v-if="itemImage"
          :src="itemImage"
          :alt="itemName"
          class="purchase-item-cutout"
          :class="`purchase-item-${kind}`"
          data-testid="asset-purchase-item"
        />

        <div v-if="!isPreview" class="celebration-layer" aria-hidden="true">
          <span
            v-for="index in 9"
            :key="index"
            :style="{
              left: `${5 + index * 11}%`,
              animationDelay: `${index * -0.42}s`,
              animationDuration: `${3.1 + (index % 3) * 0.45}s`,
            }"
          >{{ index % 2 ? '◆' : '✦' }}</span>
        </div>

        <div class="purchase-copy">
          <div class="purchase-kicker">{{ kicker }}</div>
          <div class="purchase-title">{{ headline }}</div>
          <div class="purchase-subtitle">{{ itemName }}</div>
          <div class="showcase-bonuses" data-testid="asset-showcase-bonuses">
            <span>{{ isPreview ? '参考价' : '成交价' }} {{ formatPrice(itemPrice) }}</span>
            <span v-if="itemHealth > 0">💚 健康 +{{ itemHealth }}</span>
            <span v-if="itemFame > 0">⭐ 名声 +{{ itemFame }}</span>
          </div>
        </div>

        <div v-if="!isPreview" class="purchase-sequence" data-testid="asset-purchase-sequence">
          <div
            v-for="(beat, index) in beats"
            :key="beat.text"
            class="purchase-beat"
            :style="{ animationDelay: `${0.55 + index * 2.2}s` }"
            :data-testid="`asset-purchase-beat-${index + 1}`"
          >
            <span>{{ beat.icon }}</span>
            <div>
              <strong>{{ beat.title }}</strong>
              <small>{{ beat.text }}</small>
            </div>
          </div>
        </div>

        <div v-else class="showcase-specs" data-testid="asset-showcase-specs">
          <span>{{ meta.showcaseIcon }} {{ meta.showcaseType }}</span>
          <span>🖼️ 场景化收藏展示</span>
          <span>✦ {{ meta.showcaseLine }}</span>
        </div>

        <div
          class="showcase-ownership"
          :class="{ owned: !isPreview || owned, 'purchase-complete': !isPreview }"
          :data-testid="isPreview ? 'asset-ownership-status' : 'asset-purchase-status'"
        >
          <span class="showcase-ownership-dot"></span>
          {{ isPreview ? (owned ? '已购买' : '未购买') : '已购买' }}
        </div>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { formatPrice } from '@/src/utils/format'

const props = defineProps({
  modelValue: Boolean,
  kind: {
    type: String,
    default: 'car',
    validator: value => ['car', 'house'].includes(value),
  },
  mode: {
    type: String,
    default: 'purchase',
    validator: value => ['purchase', 'preview'].includes(value),
  },
  item: {
    type: Object,
    default: () => ({}),
  },
  owned: Boolean,
})

const emit = defineEmits(['update:modelValue'])
const visible = ref(props.modelValue)
const playKey = ref(0)
const isPreview = computed(() => props.mode === 'preview')

const itemName = computed(() => props.kind === 'car'
  ? props.item?.ciname || '新座驾'
  : props.item?.hiname || '新居')
const itemPrice = computed(() => Number(props.kind === 'car' ? props.item?.ciprice : props.item?.hiprice) || 0)
const itemHealth = computed(() => Number(props.kind === 'car' ? props.item?.cihealth : props.item?.hihealth) || 0)
const itemFame = computed(() => Number(props.kind === 'car' ? props.item?.cifame : props.item?.hifame) || 0)
const itemImage = computed(() => props.kind === 'car' ? props.item?.ciimg : props.item?.hiimg)

const meta = computed(() => props.kind === 'car' ? {
  backdrop: '/images/carinfo/car-moments/car-showroom.webp',
  kicker: 'NEW CAR DELIVERY',
  previewKicker: 'PRIVATE GARAGE COLLECTION',
  headline: '属于你的新座驾正式交付',
  showcaseHeadline: '在私人展厅里鉴赏这台座驾',
  showcaseIcon: '🏁',
  showcaseType: '车辆收藏',
  showcaseLine: '将车辆放入真实展厅，更直观地感受比例与轮廓',
} : {
  backdrop: '/images/houseinfo/house-moments/home-handover.webp',
  kicker: 'NEW HOME HANDOVER',
  previewKicker: 'LIFESTYLE RESIDENCE',
  headline: '属于你的新生活正式启程',
  showcaseHeadline: '在交付现场预览这处理想住所',
  showcaseIcon: '📍',
  showcaseType: '房产收藏',
  showcaseLine: '将房屋放入交付环境，更直观地感受建筑与生活氛围',
})

const headline = computed(() => isPreview.value ? meta.value.showcaseHeadline : meta.value.headline)
const kicker = computed(() => isPreview.value ? meta.value.previewKicker : meta.value.kicker)

const beats = computed(() => props.kind === 'car' ? [
  { icon: '🧾', title: '订单确认', text: `${itemName.value}的购车手续已经完成` },
  { icon: '🔑', title: '钥匙交接', text: '车辆检查完毕，销售顾问将钥匙交到你手中' },
  { icon: '🚗', title: '正式入库', text: `${itemName.value}已停入你的车库` },
] : [
  { icon: '🧾', title: '签约完成', text: `${itemName.value}的产权手续已经完成` },
  { icon: '🗝️', title: '钥匙交接', text: '验房完成，新居钥匙正式交到你手中' },
  { icon: '🏠', title: '正式入住', text: `${itemName.value}已加入你的房产` },
])

const restart = () => {
  playKey.value += 1
}

watch(() => props.modelValue, value => {
  visible.value = value
  if (value) restart()
})
watch(visible, value => emit('update:modelValue', value))
</script>

<style scoped>
:global(.asset-purchase-moment-dialog.el-dialog) {
  padding: 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

:global(.asset-purchase-moment-dialog .el-dialog__header) {
  display: none;
}

:global(.asset-purchase-moment-dialog .el-dialog__body) {
  padding: 0;
  background: transparent;
}

.purchase-moment-shell {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.purchase-stage {
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

.purchase-stage::after {
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

.purchase-backdrop {
  position: absolute;
  inset: -4%;
  width: 108%;
  height: 108%;
  object-fit: cover;
  animation: purchase-pan 7s ease-out both;
}

.purchase-vignette {
  position: absolute;
  inset: 0;
  z-index: 1;
  background:
    linear-gradient(90deg, rgb(5 13 28 / 72%), transparent 58%),
    linear-gradient(0deg, rgb(5 10 24 / 76%), transparent 48%);
}

.purchase-overlay-actions {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 7;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.purchase-overlay-actions :deep(.el-button) {
  width: 28px;
  min-width: 28px;
  height: 28px;
  margin-left: 0;
  padding: 0;
  border-color: rgb(255 255 255 / 38%);
  background: rgb(7 17 36 / 68%);
  color: #fff;
  box-shadow: 0 4px 14px rgb(0 0 0 / 28%);
  font-size: 15px;
  backdrop-filter: blur(7px);
}

.purchase-overlay-actions :deep(.el-button--primary) {
  border-color: rgb(255 216 140 / 72%);
  background: rgb(47 107 208 / 82%);
}

.purchase-item-cutout {
  position: absolute;
  right: 3%;
  z-index: 2;
  object-fit: contain;
  object-position: center bottom;
  filter: drop-shadow(0 18px 18px rgb(0 0 0 / 48%));
  animation: purchase-item-enter 1.25s 0.35s cubic-bezier(0.2, 0.8, 0.2, 1) both,
    purchase-item-float 3.2s 1.6s ease-in-out infinite alternate;
}

.purchase-item-car {
  bottom: 5%;
  width: 56%;
  height: 55%;
}

.purchase-item-house {
  right: 2%;
  bottom: -2%;
  width: 55%;
  height: 76%;
}

.celebration-layer {
  position: absolute;
  inset: 0;
  z-index: 4;
  pointer-events: none;
}

.purchase-copy {
  position: absolute;
  top: 28px;
  left: 30px;
  z-index: 3;
  max-width: 55%;
  color: #fff;
  text-shadow: 0 2px 13px rgb(0 0 0 / 75%);
  animation: purchase-copy-enter 0.9s 0.25s ease-out both;
}

.purchase-kicker {
  margin-bottom: 5px;
  color: #ffd88c;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.18em;
}

.purchase-title {
  font-size: 24px;
  font-weight: 750;
}

.purchase-subtitle {
  margin-top: 6px;
  color: rgb(255 255 255 / 90%);
  font-size: 15px;
  font-weight: 650;
}

.showcase-bonuses {
  display: flex;
  flex-wrap: wrap;
  gap: 5px 12px;
  margin-top: 9px;
  color: rgb(255 255 255 / 86%);
  font-size: 10px;
  line-height: 1.4;
}

.celebration-layer span {
  position: absolute;
  top: -25px;
  color: #ffd166;
  font-size: 15px;
  opacity: 0;
  text-shadow: 0 0 8px rgb(255 209 102 / 80%);
  animation: confetti-fall 3.4s ease-in infinite;
}

.purchase-sequence {
  position: absolute;
  left: 28px;
  z-index: 3;
  bottom: 64px;
  width: min(300px, 45%);
  height: 66px;
}

.showcase-specs {
  position: absolute;
  bottom: 58px;
  left: 28px;
  z-index: 3;
  display: flex;
  width: min(330px, 48%);
  flex-direction: column;
  gap: 5px;
  color: rgb(255 255 255 / 88%);
  font-size: 10px;
  line-height: 1.4;
  text-shadow: 0 2px 9px rgb(0 0 0 / 82%);
  pointer-events: none;
  animation: showcase-hint-lifecycle 4.6s ease both;
}

.showcase-specs span:first-child {
  color: #ffd88c;
  font-size: 12px;
  font-weight: 700;
}

.purchase-beat {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 13px;
  color: #fff;
  opacity: 0;
  text-shadow: 0 2px 9px rgb(0 0 0 / 82%);
  animation: purchase-beat-show 2.05s ease both;
}

.purchase-beat > span { flex: 0 0 auto; font-size: 25px; }
.purchase-beat div { display: flex; flex-direction: column; gap: 2px; }
.purchase-beat strong { font-size: 12px; }
.purchase-beat small { color: rgb(255 255 255 / 78%); font-size: 10px; line-height: 1.35; }

.showcase-ownership {
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

.showcase-ownership-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #94a3b8;
  box-shadow: 0 0 8px rgb(148 163 184 / 70%);
}

.showcase-ownership.owned {
  color: #d1fae5;
}

.showcase-ownership.owned .showcase-ownership-dot {
  background: #6ee7b7;
  box-shadow: 0 0 9px rgb(110 231 183 / 82%);
}

.showcase-ownership.purchase-complete {
  opacity: 0;
  animation: purchase-status-enter 0.9s 6.1s ease both;
}

@keyframes purchase-pan {
  from { transform: scale(1.01) translate3d(-0.8%, 0.4%, 0); }
  to { transform: scale(1.09) translate3d(1.1%, -0.5%, 0); }
}

@keyframes purchase-copy-enter {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes showcase-hint-lifecycle {
  0% { opacity: 0; transform: translateY(6px); }
  12%, 70% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(3px); }
}

@keyframes purchase-item-enter {
  from { opacity: 0; transform: translate3d(15%, 5%, 0) scale(0.9); }
  to { opacity: 1; transform: translate3d(0, 0, 0) scale(1); }
}

@keyframes purchase-item-float {
  from { transform: translate3d(0, 0, 0) scale(1); }
  to { transform: translate3d(0, -1.5%, 0) scale(1.015); }
}

@keyframes confetti-fall {
  0% { opacity: 0; transform: translateY(-20px) rotate(0); }
  14% { opacity: 0.9; }
  100% { opacity: 0; transform: translateY(410px) rotate(420deg); }
}

@keyframes purchase-beat-show {
  0% { opacity: 0; transform: translateY(13px) scale(0.97); }
  15%, 78% { opacity: 1; transform: translateY(0) scale(1); }
  100% { opacity: 0; transform: translateY(-7px) scale(0.99); }
}

@keyframes purchase-status-enter {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

</style>
