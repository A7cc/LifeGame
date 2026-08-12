<template>
  <el-dialog
    v-model="visible"
    aria-label="约会互动场景"
    width="760px"
    align-center
    destroy-on-close
    :show-close="false"
    class="spouse-moment-dialog"
    data-testid="spouse-moment-dialog"
  >
    <div class="moment-shell">
      <div :key="playKey" class="moment-stage" :class="[`moment-${kind}`, `moment-${resolvedVariant}`]">
        <img :src="backdrop" class="moment-backdrop" alt="" />
        <div class="moment-vignette"></div>
        <div class="moment-overlay-actions" data-testid="moment-overlay-actions">
          <el-button
            circle
            size="small"
            type="primary"
            title="关闭短片"
            aria-label="关闭短片"
            data-testid="close-spouse-moment"
            @click="visible = false"
          >×</el-button>
          <el-button
            circle
            size="small"
            title="重播短片"
            aria-label="重播短片"
            data-testid="replay-spouse-moment"
            @click="restart"
          >↻</el-button>
          <el-button
            v-for="action in interactive && kind === 'date' ? interactionActions : []"
            :key="action.value"
            circle
            size="small"
            :type="action.type"
            :disabled="interactionPending || !action.enabled"
            :loading="interactionPending && activeAction === action.value"
            :class="{ 'outfit-action-anchor': action.value === 'outfit' }"
            :title="action.enabled ? `${action.label}：${action.description}` : action.reason"
            :aria-label="action.enabled ? action.label : `${action.label}（${action.reason}）`"
            :data-testid="`date-interaction-${action.value}`"
            @click="selectInteraction(action.value)"
          >{{ action.icon }}</el-button>
        </div>

        <div v-if="outfitChooserOpen" class="outfit-picker" data-testid="dating-outfit-picker">
          <div class="outfit-buttons">
            <el-button
              v-for="outfitOption in availableOutfits"
              :key="outfitOption.key"
              circle
              size="small"
              :type="outfitVariant === outfitOption.key ? 'primary' : ''"
              :plain="outfitVariant !== outfitOption.key"
              :disabled="interactionPending"
              :title="outfitOption.label"
              :aria-label="outfitOption.label"
              :data-testid="`dating-outfit-${outfitOption.key}`"
              @click="selectOutfit(outfitOption)"
            >{{ outfitOption.icon }}</el-button>
          </div>
        </div>

        <div v-if="kind === 'bath'" class="steam-layer" aria-hidden="true">
          <span
            v-for="index in 6"
            :key="index"
            :style="{ left: `${70 + index * 58}px`, animationDelay: `${index * -0.45}s` }"
          ></span>
        </div>
        <div v-else-if="isTenseMoment" class="tension-layer" aria-hidden="true">
          <span
            v-for="index in 5"
            :key="index"
            :style="{ top: `${38 + index * 54}px`, left: `${42 + index * 105}px`, animationDelay: `${index * -0.35}s` }"
          >💢</span>
        </div>
        <div v-else class="sparkle-layer" aria-hidden="true">
          <span
            v-for="index in 7"
            :key="index"
            :style="{
              top: `${35 + index * 37}px`,
              left: `${45 + index * 73}px`,
              fontSize: `${10 + index}px`,
              animationDelay: `${index * -0.5}s`,
            }"
          >♥</span>
        </div>

        <div class="moment-copy">
          <div class="moment-kicker">{{ momentKicker }}</div>
          <div class="moment-title">{{ headline }}</div>
          <div class="moment-subtitle">{{ subtitle }}</div>
        </div>

        <div class="spouse-frame" data-testid="dating-scene-character">
          <img v-if="portraitSource" :src="portraitSource" :alt="spouse?.dname" class="spouse-portrait" />
          <div v-else class="spouse-fallback spouse-portrait">{{ spouse?.dsex ? '🕺' : '💃' }}</div>
          <span
            v-for="(beat, index) in interactionBeats"
            :key="`${beat.icon}-${index}`"
            class="portrait-action"
            :style="{ animationDelay: `${0.55 + index * 2.8}s` }"
            aria-hidden="true"
          >{{ beat.icon }}</span>
        </div>

        <div class="interaction-sequence" data-testid="moment-interaction-sequence">
          <div
            v-for="(beat, index) in interactionBeats"
            :key="beat.text"
            class="interaction-beat"
            :style="{ animationDelay: `${0.55 + index * 2.8}s` }"
            :data-testid="`moment-interaction-${index + 1}`"
          >
            <span class="interaction-icon">{{ beat.icon }}</span>
            <span>{{ beat.text }}</span>
          </div>
        </div>

        <div
          class="moment-relationship"
          :class="{ established: relationshipEstablished }"
          data-testid="moment-relationship-status"
        >
          <span class="moment-relationship-dot"></span>
          <strong>{{ relationshipStatus || '陌生人' }}</strong>
        </div>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { getDatingScene } from '@/src/utils/datingScenes'
import { getDatingGiftScene } from '@/src/utils/datingGifts'
import { availableDatingOutfits, datingRelationshipRank } from '@/src/utils/datingOutfits'

const props = defineProps({
  modelValue: Boolean,
  spouse: {
    type: Object,
    default: () => ({}),
  },
  kind: {
    type: String,
    default: 'date',
    validator: value => ['date', 'bath', 'gift'].includes(value),
  },
  location: {
    type: String,
    default: '',
  },
  variant: {
    type: String,
    default: '',
    validator: value => ['', 'date', 'chat', 'caress', 'kiss', 'outfit', 'intimacy', 'argument', 'favorite', 'disliked', 'missed', 'rejected', 'risky-success'].includes(value),
  },
  characterImage: {
    type: String,
    default: '',
  },
  outfit: {
    type: String,
    default: '',
  },
  outfitVariant: {
    type: String,
    default: 'career',
  },
  interactive: Boolean,
  interactionPending: Boolean,
  interactionCompleted: Boolean,
  playerAge: {
    type: Number,
    default: 0,
  },
  gift: {
    type: String,
    default: '',
  },
  event: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['update:modelValue', 'interact'])
const visible = ref(props.modelValue)
const playKey = ref(0)
const activeAction = ref('')
const outfitChooserOpen = ref(false)

const dateScene = computed(() => getDatingScene(props.location))
const giftScene = computed(() => getDatingGiftScene(props.gift))
const portraitSource = computed(() => props.characterImage || props.spouse?.dimage || '')
const backdrop = computed(() => {
  if (props.kind === 'bath') return '/images/datinginfo/dating-moments/bath-spa.webp'
  if (props.kind === 'gift') return giftScene.value.image
  return dateScene.value.image
})
const kissStatuses = new Set(['恋人', '专属恋人', '爱人', '已婚'])
const relationshipStatus = computed(() => props.spouse?.daffinitylevel || '')
const relationshipEstablished = computed(() => Boolean(relationshipStatus.value) && relationshipStatus.value !== '陌生人')
const resolvedVariant = computed(() => {
  if (props.kind === 'bath') return 'bath'
  if (props.kind === 'gift') {
    if (props.variant === 'risky-success') return 'gift-risky-success'
    if (['disliked', 'missed', 'rejected'].includes(props.variant)) return 'gift-rejected'
    return 'gift-favorite'
  }
  if (props.variant) return props.variant
  return kissStatuses.has(relationshipStatus.value) ? 'kiss' : 'chat'
})

const momentMeta = computed(() => ({
  bath: {
    kicker: 'MARRIED LIFE', label: '婚后时光',
    headline: `和${props.spouse?.dname || '爱人'}一起放松`,
  },
  date: {
    kicker: 'DATE MEMORY', label: '约会时光',
    headline: `和${props.spouse?.dname || '爱人'}共度浪漫时光`,
  },
  chat: {
    kicker: 'HEART TO HEART', label: '真心聊天',
    headline: '一次更深的了解',
  },
  caress: {
    kicker: 'GENTLE TOUCH', label: '温柔抚摸',
    headline: '小心靠近彼此的瞬间',
  },
  kiss: {
    kicker: 'FIRST KISS', label: '温柔亲吻',
    headline: '心意靠近的瞬间',
  },
  intimacy: {
    kicker: 'PRIVATE MOMENT', label: '亲密过夜',
    headline: '只属于你们的夜晚',
  },
  outfit: {
    kicker: 'NEW LOOK', label: '换装互动',
    headline: '为这次约会换上新造型',
  },
  argument: {
    kicker: 'RELATIONSHIP TEST', label: '意见争执',
    headline: '气氛突然紧张起来',
  },
  'gift-favorite': {
    kicker: 'PERFECT GIFT', label: '正中喜好',
    headline: `${props.spouse?.dname || '对方'}收到了一份惊喜`,
  },
  'gift-disliked': {
    kicker: 'WRONG GIFT', label: '没有猜中',
    headline: '这份礼物似乎选错了',
  },
	'gift-risky-success': {
		kicker: 'HIGH RISK · HIGH REWARD', label: '意外打动',
		headline: '大胆的选择带来了意外惊喜',
	},
	'gift-rejected': {
		kicker: 'GIFT NOT ACCEPTED', label: '这次没有送出',
		headline: '这份心意没有在今天抵达',
	},
})[resolvedVariant.value])

const momentKicker = computed(() => momentMeta.value.kicker)
const headline = computed(() => momentMeta.value.headline)
const subtitle = computed(() => {
  if (props.kind === 'bath') return '蒸汽、烛光，还有两个人的安静片刻'
  if (props.kind === 'gift') return `${giftScene.value.icon} ${props.gift || '神秘礼物'} · ${momentMeta.value.label}`
  return `${dateScene.value.icon} ${props.location || '浪漫约会'} · ${momentMeta.value.label}`
})
const isTenseMoment = computed(() => ['argument', 'gift-disliked', 'gift-rejected'].includes(resolvedVariant.value))

const spousePronoun = computed(() => props.spouse?.dsex ? '他' : '她')
const availableOutfits = computed(() => availableDatingOutfits(
  relationshipStatus.value,
  props.playerAge,
  props.spouse?.dage,
))
const interactionActions = computed(() => {
  const rank = datingRelationshipRank(relationshipStatus.value)
  const bothAdults = props.playerAge >= 18 && Number(props.spouse?.dage || 0) >= 18
  const actions = [
    { value: 'chat', icon: '💬', label: '聊天', min: 0, type: 'primary', description: '安全的交流，好感额外+1' },
    { value: 'outfit', icon: '👗', label: '指定换装', min: 1, type: 'success', description: '选择当前关系阶段已解锁的服装，好感额外+1' },
    { value: 'caress', icon: '🫶', label: '抚摸', min: 2, type: 'warning', description: '确认对方意愿后温柔靠近，好感额外+2' },
    { value: 'kiss', icon: '💋', label: '亲吻', min: 4, type: 'danger', description: '恋人阶段解锁，好感额外+3' },
    { value: 'intimacy', icon: '🌙', label: '亲密过夜', min: 6, type: 'danger', description: '双方成年且爱人/已婚解锁，好感额外+5' },
  ]
  return actions.map(action => {
    const adultBlocked = action.value === 'intimacy' && !bothAdults
    const statusBlocked = rank < action.min
    return {
      ...action,
      enabled: !adultBlocked && !statusBlocked,
      reason: adultBlocked
        ? '亲密过夜仅限双方成年'
        : statusBlocked ? `当前关系阶段不足（${relationshipStatus.value || '陌生人'}）` : '',
    }
  })
})

const selectInteraction = action => {
  if (action === 'outfit') {
    outfitChooserOpen.value = !outfitChooserOpen.value
    return
  }
  outfitChooserOpen.value = false
  activeAction.value = action
  emit('interact', action, '')
}

const selectOutfit = outfitOption => {
  activeAction.value = 'outfit'
  outfitChooserOpen.value = false
  emit('interact', 'outfit', outfitOption.key)
}

const interactionBeats = computed(() => {
  const name = props.spouse?.dname || '对方'
  const pronoun = spousePronoun.value
  if (props.kind === 'bath') {
    return [
      { icon: '💬', text: `你和${name}聊起今天的小事，蒸汽里满是笑声` },
      { icon: '🫶', text: `你轻轻替${pronoun}擦去发梢的水珠` },
      { icon: '💕', text: `${pronoun}靠近你，你们在烛光中温柔依偎` },
    ]
  }
  if (props.kind === 'gift') {
    if (resolvedVariant.value === 'gift-rejected' || resolvedVariant.value === 'gift-disliked') {
      return [
        { icon: '🎁', text: `你把${props.gift || '礼物'}递给了${name}` },
        { icon: '😕', text: `${pronoun}看了看礼物，表情显得有些迟疑` },
        { icon: '💭', text: props.event || `${name}礼貌地表示这次不能收下` },
      ]
    }
		if (resolvedVariant.value === 'gift-risky-success') {
			return [
				{ icon: '🎲', text: `你冒险选择了${props.gift || '礼物'}，把它递给${name}` },
				{ icon: '✨', text: `这并不是${pronoun}平时偏爱的类型，却意外带来了新鲜感` },
				{ icon: '💖', text: props.event || `${name}被这次大胆的心意打动了` },
			]
		}
    return [
      { icon: '🎁', text: `你把精心挑选的${props.gift || '礼物'}递给了${name}` },
      { icon: '✨', text: `${pronoun}拆开礼物，眼睛一下亮了起来` },
      { icon: '💕', text: props.event || `${name}很开心你记住了自己的喜好` },
    ]
  }

  const beats = {
    chat: [
      { icon: '💬', text: `你和${name}聊起最近的生活` },
      { icon: '🎧', text: `你认真倾听${pronoun}的想法，也分享了自己的心事` },
      { icon: '😊', text: `话题越来越自然，${name}露出了放松的笑容` },
    ],
    caress: [
      { icon: '🫶', text: `你先询问${name}是否愿意更靠近一些` },
      { icon: '🤝', text: `${pronoun}给出了肯定的回应，彼此都放松了下来` },
      { icon: '💕', text: `你轻轻抚摸${pronoun}的脸颊和手背，${pronoun}微笑着靠近你` },
    ],
    kiss: [
      { icon: '👀', text: `你们的目光一次次相遇，距离也慢慢靠近` },
      { icon: '🫶', text: `你轻声确认${pronoun}的心意，${pronoun}微笑着点了点头` },
      { icon: '💋', text: `你轻轻抚摸${pronoun}的脸颊，你们交换了一个温柔的吻` },
    ],
    intimacy: [
      { icon: '💬', text: `你和${name}坦诚聊起彼此的感受和期待` },
      { icon: '🤝', text: `两位成年人再次确认了彼此的意愿和边界` },
      { icon: '🌙', text: `灯光慢慢暗下，你们共度了一个亲密而私密的夜晚` },
    ],
    outfit: [
      { icon: '👗', text: `你邀请${name}一起挑选今天的约会造型` },
      { icon: '🪞', text: `${pronoun}在镜子前试了几种风格，也询问了你的意见` },
      { icon: '✨', text: `${props.outfit || '新的约会造型'}让你们眼前一亮` },
    ],
    argument: [
      { icon: '❗', text: `一个话题引发了分歧，你和${name}都没能立刻说服对方` },
      { icon: '💢', text: `语气变得急促，原本轻松的气氛也紧张起来` },
      { icon: '🌫️', text: `你们最后停下争论，决定先给彼此一点冷静的时间` },
    ],
  }
  return beats[resolvedVariant.value] || beats.chat
})

const restart = () => {
  playKey.value += 1
}

watch(() => props.modelValue, value => {
  visible.value = value
  if (value) {
    outfitChooserOpen.value = false
    restart()
  }
})
watch(() => props.variant, () => {
  activeAction.value = ''
  if (visible.value) restart()
})
watch(visible, value => emit('update:modelValue', value))
</script>

<style scoped>
:global(.spouse-moment-dialog.el-dialog) {
  padding: 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

:global(.spouse-moment-dialog .el-dialog__header) {
  display: none;
}

:global(.spouse-moment-dialog .el-dialog__body) {
  padding: 0;
  background: transparent;
}

.moment-shell {
  display: flex;
  flex-direction: column;
}

.moment-stage {
  position: relative;
  aspect-ratio: 16 / 9;
  overflow: hidden;
  border: 4px solid #d5a43a;
  border-radius: 14px;
  background: #17131c;
  box-shadow:
    0 0 0 1px #6b450f,
    0 0 0 3px #f1cc72,
    0 0 0 4px #593807,
    0 0 22px rgb(226 176 65 / 42%),
    0 24px 58px rgb(3 7 18 / 64%);
  isolation: isolate;
}

.moment-stage::after {
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

.moment-backdrop {
  position: absolute;
  inset: -4%;
  width: 108%;
  height: 108%;
  object-fit: cover;
  animation: moment-pan 9s ease-out both;
}

.moment-vignette {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(10, 8, 15, 0.14) 30%, rgba(10, 8, 15, 0.76) 100%),
    linear-gradient(0deg, rgba(10, 8, 15, 0.7), transparent 42%);
}

.moment-overlay-actions {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 9;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.moment-overlay-actions :deep(.el-button),
.outfit-buttons :deep(.el-button) {
  width: 28px;
  min-width: 28px;
  height: 28px;
  margin-left: 0;
  padding: 0;
  border-color: rgb(255 255 255 / 34%);
  background: rgb(17 13 28 / 72%);
  color: #fff;
  box-shadow: 0 4px 14px rgb(0 0 0 / 30%);
  font-size: 14px;
  backdrop-filter: blur(8px);
}

.moment-overlay-actions :deep(.el-button--primary),
.outfit-buttons :deep(.el-button--primary) {
  border-color: rgb(255 185 203 / 74%);
  background: rgb(178 72 116 / 82%);
}

.moment-argument .moment-vignette,
.moment-gift-disliked .moment-vignette,
.moment-gift-rejected .moment-vignette {
  background:
    linear-gradient(90deg, rgba(36, 8, 13, 0.2) 25%, rgba(48, 9, 14, 0.8) 100%),
    linear-gradient(0deg, rgba(25, 5, 10, 0.82), transparent 48%);
}

.moment-intimacy .moment-vignette {
  animation: intimacy-fade 9s ease-in both;
}

.moment-copy {
  position: absolute;
  z-index: 3;
  left: 34px;
  top: 30px;
  max-width: 51%;
  color: #fff;
  text-shadow: 0 2px 12px rgba(0, 0, 0, 0.8);
  animation: copy-enter 1.1s 0.35s ease-out both;
}

.moment-kicker {
  margin-bottom: 5px;
  color: #ffd6a1;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.2em;
}

.moment-title {
  font-size: 25px;
  font-weight: 700;
}

.moment-subtitle {
  margin-top: 7px;
  color: rgba(255, 255, 255, 0.88);
  font-size: 13px;
}

.spouse-frame {
  position: absolute;
  z-index: 2;
  top: 12px;
  right: -24px;
  bottom: 0;
  width: 54%;
  pointer-events: none;
  animation: spouse-enter 1.2s 0.9s cubic-bezier(0.2, 0.8, 0.2, 1) both, spouse-float 3s 2.1s ease-in-out infinite;
}

.spouse-frame img,
.spouse-fallback {
  display: flex;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  object-fit: contain;
  object-position: center bottom;
  background: transparent;
  filter: drop-shadow(0 14px 16px rgba(0, 0, 0, 0.36));
  font-size: 86px;
}

.spouse-portrait {
  transform-origin: center bottom;
  animation: portrait-interact 9s ease-in-out both;
}

.moment-argument .spouse-portrait,
.moment-gift-disliked .spouse-portrait,
.moment-gift-rejected .spouse-portrait {
  animation: portrait-argument 9s ease-in-out both;
}

.moment-outfit .spouse-frame {
  filter: drop-shadow(0 0 18px rgba(204, 130, 255, 0.34));
}

.moment-outfit .spouse-portrait {
  animation: outfit-reveal 9s ease-in-out both;
}

.portrait-action {
  position: absolute;
  top: 66px;
  left: 22px;
  display: flex;
  width: 46px;
  height: 46px;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.48);
  border-radius: 50%;
  background: rgba(24, 18, 31, 0.82);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.32);
  font-size: 23px;
  opacity: 0;
  animation: action-pop 2.55s ease both;
}

.interaction-sequence {
  position: absolute;
  z-index: 3;
  bottom: 42px;
  left: 34px;
  width: min(390px, 55%);
  height: 64px;
}

.interaction-beat {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  color: #fff;
  font-size: 13px;
  line-height: 1.45;
  opacity: 0;
  text-shadow: 0 2px 9px rgba(0, 0, 0, 0.82);
  animation: beat-show 2.55s ease both;
}

.interaction-icon {
  flex: 0 0 auto;
  font-size: 24px;
}

.moment-argument .interaction-beat,
.moment-gift-disliked .interaction-beat,
.moment-gift-rejected .interaction-beat {
  color: #ffd2d2;
}

.moment-relationship {
  position: absolute;
  right: 24px;
  bottom: 20px;
  z-index: 4;
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.82);
  font-size: 11px;
  font-weight: 700;
  text-shadow: 0 2px 9px rgba(0, 0, 0, 0.82);
}

.moment-relationship strong {
  font: inherit;
}

.moment-relationship-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #94a3b8;
  box-shadow: 0 0 8px rgba(148, 163, 184, 0.7);
}

.moment-relationship.established {
  color: #d1fae5;
}

.moment-relationship.established .moment-relationship-dot {
  background: #6ee7b7;
  box-shadow: 0 0 9px rgba(110, 231, 183, 0.82);
}

.outfit-picker {
  position: absolute;
  top: 127px;
  right: 48px;
  z-index: 9;
  display: flex;
  width: max-content;
  max-width: 430px;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 12px;
  background: linear-gradient(145deg, rgb(31 20 42 / 92%), rgb(12 18 34 / 92%));
  color: #fff;
  box-shadow: 0 10px 30px rgb(0 0 0 / 34%);
  font-size: 11px;
  backdrop-filter: blur(10px);
  transform: translateY(-50%);
}

.outfit-picker::after {
  position: absolute;
  top: 50%;
  right: -6px;
  width: 10px;
  height: 10px;
  background: rgb(22 18 36 / 96%);
  content: '';
  transform: translateY(-50%) rotate(45deg);
}

.outfit-buttons {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 5px;
}

.steam-layer span {
  position: absolute;
  bottom: 95px;
  width: 36px;
  height: 90px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.13);
  filter: blur(10px);
  animation: steam-rise 3.2s calc(var(--steam-index) * -0.45s) ease-in infinite;
}

.sparkle-layer span {
  position: absolute;
  color: rgba(255, 210, 185, 0.75);
  animation: sparkle-float 3.6s ease-in-out infinite;
}

.tension-layer span {
  position: absolute;
  color: rgba(255, 142, 142, 0.8);
  font-size: 20px;
  animation: tension-jitter 1.35s ease-in-out infinite;
}

@keyframes moment-pan {
  from { transform: scale(1.02) translate3d(-0.7%, 0.3%, 0); }
  to { transform: scale(1.1) translate3d(1.2%, -0.5%, 0); }
}

@keyframes copy-enter {
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes spouse-enter {
  from { opacity: 0; transform: translateX(42px) rotate(2deg); }
  to { opacity: 1; transform: translateX(0) rotate(0); }
}

@keyframes spouse-float {
  50% { transform: translateY(-5px); }
}

@keyframes portrait-interact {
  0%, 20% { transform: scale(1) translateX(0) rotate(0); }
  32%, 48% { transform: scale(1.025) translateX(-3px) rotate(-0.7deg); }
  62%, 82% { transform: scale(1.055) translateX(-6px) rotate(-1.2deg); }
  100% { transform: scale(1.02) translateX(-2px) rotate(0); }
}

@keyframes portrait-argument {
  0%, 22% { transform: translateX(0) rotate(0); }
  34%, 48% { transform: translateX(5px) rotate(0.8deg); }
  62%, 80% { transform: translateX(9px) rotate(1.2deg); filter: saturate(0.82); }
  100% { transform: translateX(5px) rotate(0); filter: saturate(0.72); }
}

@keyframes outfit-reveal {
  0%, 24% { opacity: 0.3; transform: scale(0.96); filter: saturate(0.55) brightness(0.72); }
  38%, 62% { opacity: 1; transform: scale(1.045); filter: saturate(1.2) brightness(1.08); }
  100% { opacity: 1; transform: scale(1.02); filter: saturate(1.06); }
}

@keyframes intimacy-fade {
  0%, 64% { background-color: rgba(8, 6, 13, 0); }
  100% { background-color: rgba(8, 6, 13, 0.88); }
}

@keyframes beat-show {
  0% { opacity: 0; transform: translateY(14px) scale(0.97); }
  14%, 78% { opacity: 1; transform: translateY(0) scale(1); }
  100% { opacity: 0; transform: translateY(-7px) scale(0.99); }
}

@keyframes action-pop {
  0% { opacity: 0; transform: translateY(10px) scale(0.55) rotate(-12deg); }
  18%, 74% { opacity: 1; transform: translateY(0) scale(1) rotate(0); }
  100% { opacity: 0; transform: translateY(-12px) scale(0.75) rotate(8deg); }
}

@keyframes steam-rise {
  0% { opacity: 0; transform: translateY(20px) scale(0.7); }
  45% { opacity: 0.65; }
  100% { opacity: 0; transform: translateY(-115px) scale(1.35); }
}

@keyframes sparkle-float {
  0%, 100% { opacity: 0.18; transform: translateY(9px) scale(0.8); }
  50% { opacity: 0.9; transform: translateY(-9px) scale(1.15); }
}

@keyframes tension-jitter {
  0%, 100% { opacity: 0.25; transform: translateX(-2px) rotate(-5deg); }
  45% { opacity: 0.9; transform: translateX(3px) rotate(5deg) scale(1.12); }
}

@media (prefers-reduced-motion: reduce) {
  .moment-stage * {
    animation-duration: 1ms !important;
    animation-delay: 0ms !important;
  }

  .interaction-beat:last-child,
  .portrait-action:last-child {
    opacity: 1;
    transform: none;
    animation: none !important;
  }
}

@media (max-width: 720px) {
  .outfit-picker {
    max-width: min(360px, calc(100% - 76px));
    flex-wrap: wrap;
    justify-content: flex-end;
  }
}
</style>
