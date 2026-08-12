<template>
  <div class="product-panel" data-testid="page-dating">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">约会中心</div>
      <div v-if="meetingScenes.length" class="meeting-scene-controls">
        <span class="meeting-scene-label">📍 场景邂逅</span>
        <el-select v-model="selectedMeetingScene" size="small" placeholder="选择场景" class="meeting-scene-select">
          <el-option v-for="scene in meetingScenes" :key="scene" :label="scene" :value="scene" />
        </el-select>
        <el-button
          data-testid="visit-dating-scene"
          size="small"
          type="success"
          :loading="visitingScene"
          :disabled="!selectedMeetingScene"
          @click="visitMeetingScene"
        >前往</el-button>
      </div>
    </div>

    <!-- 中间信息区域 -->
    <div class="panel-main">
      <!-- 左侧：异性约会对象列表 -->
      <div class="panel-main-left">
        <div class="section-header">
          <span class="section-title" data-testid="dating-list-title">{{ partnerLabel }}列表</span>
          <div class="section-controls">
            <el-tag size="small">{{ availablePropertys.filter(d => d.dunlocked).length }}/{{ availablePropertys.length }}</el-tag>
          </div>
        </div>

        <div class="section-content" v-loading="loading">
          <div v-if="availablePropertys.length === 0 && !loading" class="section-empty">
            <span class="section-empty-icon">🔒</span>
            <span class="section-empty-text">暂无{{ partnerLabel }}，提升名声解锁更多</span>
          </div>
          <div v-else class="property-grid">
            <div
              v-for="propertyinfo in availablePropertys"
              :key="propertyinfo.did"
              class="property-card"
              :class="{ locked: !propertyinfo.dunlocked }"
              :data-testid="`dating-card-${propertyinfo.did}`"
              :data-dating-sex="propertyinfo.dsex ? 'male' : 'female'"
            >
              <div class="card-header">
                <div class="card-avatar" @click.stop>
                  <el-image
                    :src="propertyinfo.dimage"
                    fit="contain"
                    class="avatar-image"
                    @click.stop="openDatingPortraitPreview(propertyinfo)"
                  >
                    <template #error>
                      <span class="avatar-loading">{{ propertyinfo.dsex ? '🕺' : '💃' }}</span>
                    </template>
                  </el-image>
                </div>
                <div class="card-heading">
                  <span class="property-name">{{ propertyinfo.dname }}</span>
                  <span class="card-meta">{{ propertyinfo.dnationality }} · {{ propertyinfo.doccup }} · {{ propertyinfo.dage }}岁</span>
                </div>
                <el-tag size="small" effect="plain" :type="getRelationshipTagType(propertyinfo.daffinitylevel)">
                  {{ propertyinfo.daffinitylevel }}
                </el-tag>
              </div>

              <div class="card-body">
                <div class="card-body-content">
                  <div class="card-description" :title="propertyinfo.ddesc">{{ propertyinfo.ddesc }}</div>
                <!-- 解锁条件 -->
                <div v-if="!propertyinfo.dunlocked" class="unlock-or-like">
                  <div class="unlock-title">🔓 解锁条件</div>
                  <div class="ul-content">
                    <div v-if="propertyinfo.dmeetscene" class="ul-item scene-condition">
                      <span class="ul-icon">📍</span>
                      <span class="ul-value">满足其他条件后，前往{{ propertyinfo.dmeetscene }}尝试认识</span>
                    </div>
                    <div
                      v-for="(condition, idx) in getUnlockConditions(propertyinfo)"
                      :key="idx"
                      class="ul-item"
                      :class="{ met: condition.met }"
                    >
                      <span class="ul-icon">{{ condition.met ? '✅' : '⭕' }}</span>
                      <span class="ul-value">{{ condition.text }}</span>
                    </div>
                  </div>
                </div>

                <!-- 约会对象喜好 -->
                <div v-else class="unlock-or-like">
                  <div class="pref-title">💝 喜好</div>
                  <div class="ul-content">
                    <div v-if="propertyinfo.dgifts && propertyinfo.dgifts.length > 0" class="ul-item">
                      <span class="ul-icon">礼物：</span>
                      <span class="ul-value">{{ propertyinfo.dgifts.join('、') }}</span>
                    </div>
                    <div v-if="propertyinfo.dlocations && propertyinfo.dlocations.length > 0" class="ul-item">
                      <span class="ul-icon">地点：</span>
                      <span class="ul-value">{{ propertyinfo.dlocations.join('、') }}</span>
                    </div>
                  </div>
                </div>
                </div>

                <div class="card-footer">
                  <div class="card-actions">
                    <el-button
                      v-if="propertyinfo.dunlocked && propertyinfo.daffinitylevel !== '前任' && propertyinfo.dgifts?.length"
                      :data-testid="`gift-dating-${propertyinfo.did}`"
                      size="small"
                      type="success"
                      plain
                      :loading="giftPending === propertyinfo.did"
                      :disabled="giftPending === propertyinfo.did"
                      @click="openGiftDialog(propertyinfo)"
                    >送礼</el-button>
                    <el-button
                      v-if="propertyinfo.dunlocked && propertyinfo.daffinitylevel !== '前任'"
                      :data-testid="`date-dating-${propertyinfo.did}`"
                      type="primary"
                      size="small"
                      :disabled="userInfo.ucash < propertyinfo.dcost"
                      @click="handleCardClick(propertyinfo)"
                    >{{ userInfo.ucash < propertyinfo.dcost ? '资金不足' : '约会' }}</el-button>
                    <el-button v-else-if="!propertyinfo.dunlocked" size="small" disabled>未解锁</el-button>
                    <el-button v-else size="small" disabled>已分手</el-button>
                    <el-button
                      v-if="propertyinfo.daffinitylevel === '爱人' && !userInfo.umarrieddatingid"
                      :data-testid="`marry-dating-${propertyinfo.did}`"
                      type="danger"
                      size="small"
                      @click="marryDating(propertyinfo)"
                    >结婚</el-button>
                    <el-button
                      v-if="canBreakUp(propertyinfo)"
                      :data-testid="`breakup-dating-${propertyinfo.did}`"
                      size="small"
                      @click="breakUpDating(propertyinfo)"
                    >分手</el-button>
                    <el-button
                      v-if="propertyinfo.daffinitylevel === '已婚'"
                      :data-testid="`bath-with-spouse-${propertyinfo.did}`"
                      type="primary"
                      plain
                      size="small"
                      :loading="spouseInteractionPending === propertyinfo.did"
                      :disabled="userInfo.ucash < 200"
                      @click="batheWithSpouse(propertyinfo)"
                    >洗澡</el-button>
                    <el-button
                      v-if="propertyinfo.daffinitylevel === '已婚'"
                      :data-testid="`divorce-dating-${propertyinfo.did}`"
                      type="danger"
                      plain
                      size="small"
                      @click="divorceDating(propertyinfo)"
                    >离婚</el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧面板 -->
      <div class="panel-main-right">
        <div class="relations-box" data-testid="dating-relations-panel">
          <div class="section-header">
            <span class="section-title">我的关系</span>
            <div class="section-controls">
              <el-tag type="success" size="small">{{ activeRelationships.length }}位</el-tag>
            </div>
          </div>
          <div class="relations-content">
            <div v-if="activeRelationships.length === 0" class="section-empty">
              <span class="section-empty-icon">💌</span>
              <span class="section-empty-text">暂无稳定关系</span>
            </div>
            <div v-else class="owned-relations-list">
              <div
                v-for="dating in activeRelationships"
                :key="dating.did"
                class="owned-property-card"
                :class="{ spouse: dating.daffinitylevel === '已婚' }"
              >
                <div class="owned-property-left">
                  <el-image
                    :src="dating.dimage"
                    fit="contain"
                    class="owned-avatar"
                    @click.stop="openDatingPortraitPreview(dating)"
                  >
                    <template #error>{{ dating.dsex ? '🕺' : '💃' }}</template>
                  </el-image>
                  <div class="property-details">
                    <span class="property-name">{{ dating.dname }}</span>
                    <span class="property-meta-info">{{ dating.doccup }} · {{ dating.dage }}岁</span>
                  </div>
                </div>
                <el-tag size="small" :type="getRelationshipTagType(dating.daffinitylevel)">
                  {{ dating.daffinitylevel }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
        <LogPanel title="📋 约会记录" :items="runLogInfo" empty-icon="💌" empty-text="暂无约会记录，快去约会吧！" @clear="clearLogs" />
      </div>
    </div>
  </div>

  <!-- 约会场景选择对话框 -->
  <DialogDatingScene
    v-model="datingDialogVisible"
    :dating="selectedDating"
    @confirm="handleDatingConfirm"
  />
  <DialogDatingPortraitPreview
    v-model="portraitPreviewVisible"
    :dating="previewDating"
  />
  <DialogSpouseMoment
    v-model="spouseMomentVisible"
    :spouse="momentSpouse"
    :kind="momentKind"
    :location="momentLocation"
    :variant="momentVariant"
    :outfit="momentOutfit"
    :outfit-variant="momentOutfitVariant"
    :character-image="momentCharacterImage"
    :interactive="momentInteractive"
    :interaction-pending="momentInteractionPending"
    :interaction-completed="momentInteractionCompleted"
    :player-age="userInfo.uage"
    :gift="momentGift"
    :event="momentEvent"
    @interact="handleMomentInteraction"
  />
  <DialogDatingGift
    v-model="giftDialogVisible"
    :dating="giftDating"
    :options="giftOptions"
    :loading="giftOptionsLoading"
    :pending="giftPending === giftDating?.did"
    @confirm="handleGiftConfirm"
  />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useGameStore } from '@/src/stores/game'
import { useLogList } from '@/src/composables/useLogList'
import { BatheWithSpouse, BreakUpDating, DivorceDating, GetDatingGiftOptions, GetDatingInfo, GiveDatingGift, DoDating, DoDatingInteraction, MarryDating, VisitDatingScene } from "@/wailsjs/go/services/App.js"
import LogPanel from '@/src/components/Common/LogPanel.vue'
import DialogDatingScene from '@/src/components/Dialog/DialogDatingScene.vue'
import DialogDatingPortraitPreview from '@/src/components/Dialog/DialogDatingPortraitPreview.vue'
import DialogSpouseMoment from '@/src/components/Dialog/DialogSpouseMoment.vue'
import DialogDatingGift from '@/src/components/Dialog/DialogDatingGift.vue'
import { datingOutfitImage, findDatingOutfit } from '@/src/utils/datingOutfits'

const gameStore = useGameStore()
const userInfo = computed(() => gameStore.userInfo)
const partnerLabel = computed(() => userInfo.value.usex ? '女友' : '男友')
const partnerPronoun = computed(() => userInfo.value.usex ? '她' : '他')
const personalizePartnerText = text => String(text || '')
  .replaceAll('她', partnerPronoun.value)
  .replaceAll('女性', userInfo.value.usex ? '女性' : '男性')

const { addLog, clearLogs, logs: runLogInfo } = useLogList({ success: '💕', fail: '💔', gift: '🎁' })

// 加载状态
const loading = ref(false)

// 可用的约会对象
const availablePropertys = ref([])

// 约会对话框
const datingDialogVisible = ref(false)
const selectedDating = ref(null)
const meetingScenes = ref([])
const selectedMeetingScene = ref('')
const visitingScene = ref(false)
const spouseMomentVisible = ref(false)
const momentSpouse = ref(null)
const momentKind = ref('date')
const momentLocation = ref('')
const momentVariant = ref('')
const momentOutfit = ref('')
const momentOutfitVariant = ref('career')
const momentBaseImage = ref('')
const momentCharacterImage = ref('')
const momentInteractive = ref(false)
const momentInteractionPending = ref(false)
const momentInteractionCompleted = ref(false)
const spouseInteractionPending = ref(0)
const giftPending = ref(0)
const giftDialogVisible = ref(false)
const giftDating = ref(null)
const giftOptions = ref([])
const giftOptionsLoading = ref(false)
const momentGift = ref('')
const momentEvent = ref('')
const portraitPreviewVisible = ref(false)
const previewDating = ref(null)

const openDatingPortraitPreview = dating => {
  if (!dating?.dimage) return
  previewDating.value = dating
  portraitPreviewVisible.value = true
}

// 处理约会确认
const handleDatingConfirm = async ({ dating, location }) => {
  await doDating(dating.did, location)
}

// 解锁条件文本映射
const conditionTexts = {
  'fame': (val) => `名声达到 ${val}`,
  'cash': (val) => `现金达到 ${val.toLocaleString()}`,
  'bank': (val) => `存款达到 ${val.toLocaleString()}`,
  'house': (val) => `拥有房子 ${val} 套`,
  'car': (val) => `拥有车子 ${val} 辆`,
  'play_game': (val, target) => `玩过游戏${target ? `(${target})` : ''} ${val} 次`,
  'win_game': (val, target) => `游戏获胜${target ? `(${target})` : ''} ${val} 次`,
  'work_count': (val) => `打工次数达到 ${val}`,
  'age': (val) => `年龄达到 ${val} 岁`,
  'random': () => `随机遇见`,
  'date_count': (val) => `累计约会 ${val} 次`,
  'immunity': (val) => `健康值达到 ${val}`,
  'antique_rare': (val) => `拥有真古董 ${val} 个`,
  'lottery_win': () => `彩票中大奖`,
  'company_founder': () => `成为创业者`,
  'item_own': (val, target) => `拥有物品${target ? `(${target})` : ''} ${val} 个`,
  'stock_profit': (val) => `当年股票盈利达到 ${val.toLocaleString()}`
}

const sumItems = (items = {}) => Object.values(items || {}).reduce((sum, item) => sum + Number(item?.uitemnum || 0), 0)

const getMiniGameRecord = (target = '') => {
  const records = userInfo.value.uminigamerecords || {}
  if (target) return records[target] || { playcount: 0, wincount: 0 }
  return Object.values(records).reduce((total, record) => ({
    playcount: total.playcount + Number(record?.playcount || 0),
    wincount: total.wincount + Number(record?.wincount || 0),
  }), { playcount: 0, wincount: 0 })
}

const getOwnedItemCount = (target = '') => {
  if (target) {
    return Number(userInfo.value.uitemins?.[target]?.uitemnum || 0) + Number(userInfo.value.uitemout?.[target]?.uitemnum || 0)
  }
  return sumItems(userInfo.value.uitemins) + sumItems(userInfo.value.uitemout)
}

const getRareAntiqueCount = () => (userInfo.value.uantique || []).filter(antique => antique.aiidisplay === 1 && antique.aiamaterial >= 4).length


// 获取解锁条件
const getUnlockConditions = (date) => {
  if (!date.dmeetconditions || date.dmeetconditions.length === 0) {
    return []
  }

  return date.dmeetconditions.map(condition => {
    const type = condition.ctype
    const value = condition.cvalue
    const target = condition.ctarget

    let met = false
    let current = 0

    // 检查条件是否满足
    switch (type) {
      case 'fame':
        current = userInfo.value.ufame
        met = current >= value
        break
      case 'cash':
        current = userInfo.value.ucash
        met = current >= value
        break
      case 'bank':
        current = userInfo.value.ubank
        met = current >= value
        break
      case 'house':
        current = Object.keys(userInfo.value.uhouse || {}).length
        met = current >= value
        break
      case 'car':
        current = Object.keys(userInfo.value.ucar || {}).length       
        met = current >= value
        break
      case 'play_game':
        current = getMiniGameRecord(target).playcount
        met = current >= value
        break
      case 'win_game':
        current = getMiniGameRecord(target).wincount
        met = current >= value
        break

      case 'work_count':
        let total = 0
        Object.values(userInfo.value.uminigamerecords || {}).forEach((value) => {
          if (value.mgrtype === 'work') {
            total += value.playcount
          }
        })
        current = total
        met = total >= value
        break
      case 'age':
        current = userInfo.value.uage
        met = current >= value
        break
      case 'immunity':
        current = userInfo.value.uimmunity
        met = current >= value
        break
      case 'date_count':
        current = Object.values(userInfo.value.udating || {}).reduce((sum, d) => sum + (d.dcount || 0), 0)
        met = current >= value
        break
      case 'company_founder':
        met = userInfo.value.ucompany && Object.keys(userInfo.value.ucompany).length > 0
        break
      case 'antique_rare':
        current = getRareAntiqueCount()
        met = current >= value
        break
      case 'lottery_win':
        current = getMiniGameRecord('lottery').wincount
        met = current >= value
        break
      
      case 'item_own':
        current = getOwnedItemCount(target)
        met = current >= value
        break
      case 'stock_profit':
        current = userInfo.value.ustockprofit || 0
        met = current >= value
        break
      default:
        met = false
    }

    // 生成条件文本
    const textFunc = conditionTexts[type]
    let text = textFunc ? textFunc(value, target) : `${type}: ${value}`

    // 添加当前进度
    if (type === 'fame' || type === 'cash' || type === 'bank' || type === 'work_count' || type === 'age' || type === 'health' || type === 'date_count' || type === 'car' || type === 'house' || type === 'antique_rare' || type === 'immunity' || type === 'play_game' || type === 'win_game' || type === 'lottery_win' || type === 'item_own' || type === 'stock_profit') {
      text += ` (当前: ${current})`
    }

    return { text, met }
  })
}

// 加载约会信息
const loadDatingInfo = async () => {
  loading.value = true
  try {
    const result = await GetDatingInfo()
    if (result.code === 200) {
      
      availablePropertys.value = result.datinglist
      meetingScenes.value = result.meetingscenes || []
      if (!meetingScenes.value.includes(selectedMeetingScene.value)) {
        selectedMeetingScene.value = meetingScenes.value[0] || ''
      }
      gameStore.applyUserInfo(result.userinfo)
    } else {
      ElMessage.error(result.msg || '加载约会信息失败')
    }
  } catch (err) {
    ElMessage.error('加载约会信息失败: ' + err.message)
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  loadDatingInfo()
})

// 处理卡片点击
const handleCardClick = (date) => {
  if (!date.dunlocked) {
    ElMessage.warning(`该${partnerLabel.value}尚未解锁，继续提升自己吧！`)
    return
  }

  // 打开约会场景选择对话框
  selectedDating.value = date
  datingDialogVisible.value = true
}

const visitMeetingScene = async () => {
  if (!selectedMeetingScene.value) return
  visitingScene.value = true
  try {
    const result = await VisitDatingScene(selectedMeetingScene.value)
    if (result?.code !== 200) {
      ElMessage.error(result?.msg || '前往场景失败')
      return
    }
    if (result.userinfo) gameStore.applyUserInfo(result.userinfo)
    if ((result.met || []).length > 0) ElMessage.success(result.msg)
    else ElMessage.info(result.msg)
    await loadDatingInfo()
  } catch (error) {
    ElMessage.error('前往场景失败: ' + error.message)
  } finally {
    visitingScene.value = false
  }
}

const romanticStatuses = new Set(['暧昧中', '交往中', '恋人', '专属恋人', '爱人'])
const canBreakUp = dating => romanticStatuses.has(dating.daffinitylevel)

const relationshipStatuses = new Set(['朋友', ...romanticStatuses, '已婚', '前任'])
const relationshipRank = {
  '已婚': 0,
  '爱人': 1,
  '专属恋人': 2,
  '恋人': 3,
  '交往中': 4,
  '暧昧中': 5,
  '朋友': 6,
  '前任': 7,
}

const activeRelationships = computed(() => availablePropertys.value
  .filter(dating => dating.dunlocked && relationshipStatuses.has(dating.daffinitylevel))
  .slice()
  .sort((left, right) => (
    (relationshipRank[left.daffinitylevel] ?? 99) - (relationshipRank[right.daffinitylevel] ?? 99)
    || left.did - right.did
  )))

const getRelationshipTagType = status => ({
  '已婚': 'danger',
  '爱人': 'success',
  '专属恋人': 'success',
  '恋人': 'success',
  '交往中': 'warning',
  '暧昧中': 'warning',
  '朋友': 'info',
  '前任': 'info',
}[status] || 'info')

const runRelationshipAction = async (action, dating, confirmText) => {
  try {
    await ElMessageBox.confirm(confirmText, '关系确认', { type: 'warning' })
    const result = await action(dating.did)
    if (result?.code !== 200) {
      ElMessage.error(result?.msg || '操作失败')
      return
    }
    if (result.userinfo) gameStore.applyUserInfo(result.userinfo)
    ElMessage.success(result.msg)
    await loadDatingInfo()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error('操作失败: ' + (error?.message || error))
  }
}

const marryDating = dating => runRelationshipAction(MarryDating, dating, `确定要和${dating.dname}结婚吗？婚姻关系只能同时保留一位配偶。`)
const breakUpDating = dating => runRelationshipAction(BreakUpDating, dating, `确定要和${dating.dname}分手吗？`)
const divorceDating = dating => runRelationshipAction(DivorceDating, dating, `确定要和${dating.dname}离婚吗？`)

const openGiftDialog = async dating => {
  giftDating.value = dating
  giftOptions.value = []
  giftDialogVisible.value = true
  giftOptionsLoading.value = true
  try {
    const result = await GetDatingGiftOptions(dating.did)
    if (result?.code !== 200) {
      giftDialogVisible.value = false
      ElMessage.error(result?.msg || '获取礼物选项失败')
      return
    }
    giftOptions.value = result.options || []
  } catch (error) {
    giftDialogVisible.value = false
    ElMessage.error('获取礼物选项失败: ' + (error?.message || error))
  } finally {
    giftOptionsLoading.value = false
  }
}

const handleGiftConfirm = async option => {
  const dating = giftDating.value
  if (!dating || !option?.name) return
  giftDialogVisible.value = false
  await giveDatingGift(dating, option.name)
}

const giveDatingGift = async (dating, gift) => {
  if (giftPending.value) return
  giftPending.value = dating.did
  try {
    const result = await GiveDatingGift(dating.did, gift)
    if (result?.code !== 200) {
      ElMessage.error(result?.msg || '送礼失败')
      return
    }
    if (result.userinfo) gameStore.applyUserInfo(result.userinfo)
    const eventText = result.event || (result.success
      ? `${dating.dname}收下了${gift}。`
      : `${dating.dname}这次没有收下${gift}。`)
    addLog(`送给${dating.dname}「${gift}」：${eventText} ${result.msg}`, result.success ? 'gift' : 'fail')
    await loadDatingInfo()
    if (result.success) ElMessage.success(result.msg)
    else ElMessage.warning(result.msg)
    openSpouseMoment(
      'gift',
      {
        ...dating,
        daffinitylevel: result.datinginfo?.dstatus || dating.daffinitylevel,
      },
      '',
      result.outcome || (result.preferred ? 'favorite' : 'disliked'),
      false,
      result.gift || gift,
      eventText,
    )
  } catch (error) {
    ElMessage.error('送礼失败: ' + (error?.message || error))
  } finally {
    giftPending.value = 0
  }
}

const imageForMomentVariant = (image, kind, variant) => {
  if (kind === 'bath') return datingOutfitImage(image, 'swimwear')
  if (variant === 'intimacy') return datingOutfitImage(image, 'romantic')
  return image || ''
}

const openSpouseMoment = (kind, spouse, location = '', variant = '', interactive = false, gift = '', event = '') => {
  momentKind.value = kind
  momentSpouse.value = spouse
  momentLocation.value = location
  momentVariant.value = variant
  momentOutfit.value = ''
  momentOutfitVariant.value = kind === 'bath' ? 'swimwear' : variant === 'intimacy' ? 'romantic' : 'career'
  momentBaseImage.value = spouse?.dimage || ''
  momentCharacterImage.value = imageForMomentVariant(momentBaseImage.value, kind, variant)
  momentInteractive.value = interactive
  momentGift.value = gift
  momentEvent.value = event
  momentInteractionPending.value = false
  momentInteractionCompleted.value = false
  spouseMomentVisible.value = true
}

const handleMomentInteraction = async (action, outfitKey = '') => {
  if (!momentSpouse.value || momentInteractionPending.value) return
  // 首次互动已结算后，后续点击只切换短片分支，不再调用后端或修改属性。
  if (momentInteractionCompleted.value) {
    if (action === 'outfit') {
      const outfit = findDatingOutfit(outfitKey)
      if (!outfit) return
      momentVariant.value = action
      momentOutfitVariant.value = outfit.key
      momentOutfit.value = outfit.label
      momentCharacterImage.value = datingOutfitImage(momentBaseImage.value, outfit.key)
    } else {
      momentVariant.value = action
      momentOutfit.value = ''
      momentOutfitVariant.value = action === 'intimacy' ? 'romantic' : 'career'
      momentCharacterImage.value = imageForMomentVariant(momentBaseImage.value, momentKind.value, action)
    }
    return
  }
  momentInteractionPending.value = true
  try {
    const result = await DoDatingInteraction(momentSpouse.value.did, action, outfitKey)
    if (result?.code !== 200) {
      ElMessage.warning(result?.msg || '互动失败')
      return
    }
    if (result.userinfo) gameStore.applyUserInfo(result.userinfo)
    momentVariant.value = result.outcome || action
    momentOutfit.value = result.outfit || ''
    momentOutfitVariant.value = result.outfitvariant || 'career'
    momentCharacterImage.value = result.outfitimage || imageForMomentVariant(
      momentBaseImage.value,
      momentKind.value,
      result.outcome || action,
    )
    momentInteractionCompleted.value = true
    momentSpouse.value = {
      ...momentSpouse.value,
      daffinitylevel: result.datinginfo?.dstatus || momentSpouse.value.daffinitylevel,
    }
    const changeText = `${result.affinitychange > 0 ? '+' : ''}${result.affinitychange || 0}`
    addLog(`和${momentSpouse.value.dname}进行「${result.label || '约会'}」互动：${result.event || result.msg}，好感${changeText}`, result.outcome === 'argument' ? 'fail' : 'success')
    if (result.outcome === 'argument') {
      ElMessage.warning(`${result.msg}，好感${changeText}`)
    } else {
      ElMessage.success(`${result.msg}，好感${changeText}`)
    }
    await loadDatingInfo()
  } catch (error) {
    ElMessage.error('互动失败: ' + (error?.message || error))
  } finally {
    momentInteractionPending.value = false
  }
}

const batheWithSpouse = async dating => {
  try {
    await ElMessageBox.confirm(
      `确定要和${dating.dname}一起洗澡吗？将花费200元并消耗1次年度关系互动机会。`,
      '婚后互动',
      { type: 'info' },
    )
    spouseInteractionPending.value = dating.did
    const result = await BatheWithSpouse(dating.did)
    if (result?.code !== 200) {
      ElMessage.error(result?.msg || '互动失败')
      return
    }
    if (result.userinfo) gameStore.applyUserInfo(result.userinfo)
    addLog(`和${dating.dname}一起泡澡，好感+${result.affinitychange || 0}，健康+${result.healthchange || 0}`, 'success')
    ElMessage.success(result.msg)
    await loadDatingInfo()
    openSpouseMoment('bath', dating)
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') ElMessage.error('互动失败: ' + (error?.message || error))
  } finally {
    spouseInteractionPending.value = 0
  }
}

// 约会
const doDating = async (datingId, location = '') => {
  try {
    const dating = availablePropertys.value.find(item => item.did === datingId)
    // 获取约会结果
    const result = await DoDating(datingId, location)
    if (result.code === 200) {
      // 更新userinfo信息
      gameStore.applyUserInfo(result.userinfo)

      // 生成约会场景描述
      // 根据地点生成约会场景描述
      const locationMessages = {
        '公园': ['散步聊天，享受惬意时光', '坐在长椅上，她靠在你的肩膀', '一起放风筝，她笑得很开心', '湖边漫步，风景如画', '草坪野餐，阳光正好'],
        '电影院': ['看了一场感人的电影', '分享爆米花，她笑得很甜', '电影很精彩，她一直握着你的手', '黑暗中她悄悄靠近你', '电影结束还意犹未尽'],
        '咖啡厅': ['品尝香浓咖啡，聊了很多', '窗边闲聊，时光静好', '她分享了许多有趣的故事', '咖啡香气中，你们相视而笑', '点了她最爱的甜点'],
        '西餐厅': ['烛光晚餐，浪漫温馨', '品尝美食，她对你的品味赞不绝口', '优雅的氛围让感情升温', '她穿着精心挑选的礼服', '举杯共饮，眼神交汇'],
        '艺术馆': ['欣赏艺术作品，她的见解让你刮目相看', '一起讨论画作，发现共同审美', '她在一幅画前驻足良久', '艺术氛围中，你们的心更近了'],
        '音乐厅': ['聆听优美旋律，心灵共鸣', '她闭上眼睛沉浸在音乐中', '音乐会让她感动落泪', '浪漫的旋律让心更近了'],
        '私人会所': ['享受VIP服务，她很开心', '私密空间让你们更加亲密', '她感叹这样的约会太特别了', '奢华的环境让她眼前一亮'],
        '游艇': ['海上风光无限好', '她站在船头，海风吹拂秀发', '一起看海豚跃出水面', '夕阳西下，浪漫至极', '她感叹这是梦一般的约会'],
        '豪华酒店': ['享受顶级服务', '她被奢华的装潢惊艳到', '一起品尝精致下午茶', '窗外的夜景美不胜收'],
        '私人庄园': ['庄园漫步，仿佛置身童话', '她惊叹于这里的美丽', '一起骑马，她笑得很开心', '花园里拍照留念'],
        '海岛': ['沙滩漫步，海浪轻拍脚踝', '一起潜水，探索海底世界', '看日落，她依偎在你怀里', '椰林树影，浪漫无限'],
        '度假村': ['享受悠闲假期', '泳池边晒太阳，惬意无比', 'SPA放松身心', '她说是最难忘的约会'],
        '图书馆': ['一起挑选书籍，发现共同喜好', '安静阅读，偶尔相视一笑', '她推荐了一本好书给你', '书香氛围中，时光静好'],
        '书店': ['逛了一下午，收获满满', '她找到心仪已久的书', '一起在角落看书聊天', '分享彼此喜欢的作家'],
        '茶馆': ['品茶论道，相谈甚欢', '她学会了品茶的技巧', '茶香袅袅，岁月静好', '她喜欢这里的宁静'],
        '电竞馆': ['一起打游戏，配合默契', '她操作很厉害，让你刮目相看', '赢了比赛，她开心得跳起来', '游戏中的浪漫时刻'],
        '游戏展': ['试玩新游戏，乐趣无穷', '她cosplay了喜欢的角色', '一起排队体验热门游戏', '买了限量周边，她很开心'],
        '健身房': ['一起锻炼，互相鼓励', '她教你正确的健身姿势', '运动后的她格外迷人', '一起挑战高难度动作'],
        '海边': ['沙滩漫步，留下两行脚印', '看日落，天边霞光万丈', '她捡起一枚贝壳送给你', '海风吹拂，她笑得很美', '一起堆沙堡，童心未泯'],
        '户外': ['徒步探险，发现美景', '她勇敢地挑战了自己', '山顶远眺，心旷神怡', '野餐休息，简单而幸福'],
        '画廊': ['欣赏画作，交流心得', '她买了一幅喜欢的画', '艺术家现场作画，你们驻足观看', '她给你讲解画作背后的故事'],
        '展览': ['一起探索新奇事物', '她拍了很多照片留念', '发现共同的兴趣爱好', '展览让她大开眼界'],
        '艺术区': ['文艺气息浓厚，她很喜欢', '逛创意小店，买了纪念品', '街头艺人表演，你们驻足欣赏', '一起拍照打卡'],
        '商场': ['逛街购物，她试了很多衣服', '你送了她心仪的礼物', '一起品尝各种美食', '她挽着你的手臂不愿松开'],
        '网红地': ['打卡拍照，她很开心', '一起拍了很多美照', '她发了朋友圈炫耀', '这里确实很出片'],
        '直播间': ['她带你体验直播', '粉丝们都在祝福你们', '她在镜头前介绍你', '一起互动，乐趣无穷'],
        '医院花园': ['花园散步，空气清新', '她分享工作中的趣事', '安静的环境让你们更亲密', '她穿着白大褂也很美'],
        '温泉': ['泡温泉，身心放松', '她脸颊微红，格外迷人', '温泉边的风景很美', '一起享受这份宁静'],
        '餐厅': ['美食当前，胃口大开', '她尝了你推荐的菜', '烛光摇曳，氛围浪漫', '她对你的安排很满意'],
        '律所': ['参观她的工作环境', '她给你讲解有趣的案例', '看她认真工作的样子很迷人', '一起喝咖啡休息'],
        '商务餐厅': ['商务午餐，优雅得体', '她展现职场女性魅力', '讨论事业规划，相谈甚欢', '她对你的见解很欣赏'],
        '会所': ['私密空间，享受二人世界', '她喜欢这里的氛围', '一起品酒聊天', '她说是最放松的约会'],
        '书房': ['一起阅读，偶尔交流', '她朗读了一段喜欢的文字', '安静的氛围很温馨', '她推荐了书架上的好书'],
        '咖啡馆': ['咖啡香气中，聊了很多', '她喜欢这里的装修风格', '一起品尝手冲咖啡', '窗边闲聊，时光飞逝'],
        '沙龙': ['参加文化沙龙，收获满满', '她积极参与讨论', '认识了很多有趣的人', '她对你的见解很欣赏'],
        'KTV': ['她为你唱了一首情歌', '一起合唱，默契十足', '她唱歌很好听', '点了她最爱的歌'],
        '音乐节': ['现场气氛嗨翻天', '她跟着音乐舞动', '一起挥舞荧光棒', '她说是最棒的约会'],
        '录音棚': ['体验录音，很有趣', '她教你唱歌技巧', '一起录制了纪念歌曲', '她展现专业的一面'],
        '剧院': ['看话剧，剧情感人', '她被表演深深打动', '讨论剧情，观点一致', '她喜欢这种艺术形式'],
        '片场': ['参观电影拍摄', '她给你介绍幕后故事', '偶遇明星，她很激动', '体验了一把当演员的感觉'],
        '首映礼': ['走红毯，她光彩照人', '提前看新电影，很兴奋', '她穿着礼服很美', '难忘的夜晚'],
        '高尔夫': ['她教你打高尔夫', '一起挥杆，享受阳光', '她球技不错', '球场风景优美'],
        '酒会': ['她穿着晚礼服很美', '一起品酒社交', '她介绍你认识朋友', '优雅的夜晚'],
        '校园': ['漫步校园，回忆青春', '她分享学生时代的故事', '在操场跑步，她笑得很开心', '食堂吃饭，简单幸福'],
        '游乐园': ['坐过山车，她尖叫着抓紧你', '一起玩旋转木马', '她赢了游戏奖品', '看烟花，浪漫至极', '她像孩子一样开心'],
        '奶茶店': ['她点了最爱的奶茶', '一起尝试新品', '她分享奶茶给你喝', '简单的小确幸'],
        '舞蹈室': ['她教你跳舞', '看她跳舞很迷人', '一起练习，默契十足', '她展现优美的舞姿'],
        '舞台': ['看她排练，很专业', '她邀请你上台互动', '灯光下她格外耀眼', '一起即兴表演'],
        '训练场': ['看她训练，很认真', '她教你一些技巧', '一起挑战体能', '她鼓励你坚持'],
        '烧烤': ['一起烤串，乐趣无穷', '她烤的串很好吃', '烟火气中，感情升温', '她喜欢这种接地气的约会'],
        '拍卖会': ['参加拍卖，很新鲜', '她给你讲解拍品', '一起举牌竞拍', '她展现专业眼光'],
        '珠宝店': ['一起欣赏珠宝', '她试戴了项链', '你送了她小礼物', '她喜欢那款设计'],
        '医院': ['她带你参观', '看她工作很认真', '一起在休息室聊天', '她穿着护士服很可爱'],
        '花园': ['赏花散步，心情愉悦', '她喜欢这里的玫瑰', '一起拍照留念', '花香四溢，浪漫满溢'],
        '音乐厅': ['听音乐会，心灵享受', '她沉浸在旋律中', '一起鼓掌喝彩', '她说是最优雅的约会'],
        '琴房': ['她为你弹了一首曲子', '你尝试学钢琴', '一起四手联弹', '她教你弹简单的曲子'],
        '街拍': ['她当你的模特', '一起发现城市的美', '她很会摆pose', '拍了很多好看的照片'],
        '摄影棚': ['体验专业拍摄', '她教你摄影技巧', '一起拍写真', '她镜头感很好'],
        '大学': ['漫步校园，感受学术氛围', '她分享大学时光', '一起听讲座', '图书馆自习，安静美好'],
        '讲座': ['听知名学者演讲', '她认真做笔记', '讨论讲座内容', '她对你的见解很欣赏'],
        '机场': ['送她出差，依依不舍', '一起看飞机起落', '机场咖啡厅聊天', '她邀请你下次一起旅行'],
        '酒店': ['享受舒适的住宿', '她喜欢房间的布置', '一起在餐厅用餐', '她说是难忘的约会'],
        '景点': ['游览名胜古迹', '她给你当导游', '一起拍照打卡', '她喜欢这里的历史'],
        '家居城': ['一起逛家具', '她畅想未来的家', '挑选喜欢的风格', '她很有品味'],
        '厨房': ['一起做饭，乐趣无穷', '她教你做菜', '她厨艺很好', '一起品尝劳动成果'],
        '超市': ['一起买菜购物', '她挑选食材很认真', '讨论今晚吃什么', '简单的生活很幸福'],
        '瑜伽馆': ['一起练瑜伽', '她教你动作', '身心放松', '她展现柔韧的身姿'],
        '森林': ['森林浴，空气清新', '一起徒步探险', '她喜欢大自然', '发现小动物，她很开心'],
        '静修中心': ['冥想放松', '她教你呼吸法', '远离喧嚣，内心平静', '她说是最特别的约会'],
        '赌场': ['一起玩两把，很刺激', '她教你玩牌', '小赢一把，她很开心', '她牌技不错'],
        'VIP厅': ['享受VIP待遇', '私密空间很舒适', '她喜欢这种氛围', '一起玩牌聊天'],
        '夜场': ['感受夜生活', '她喜欢这里的音乐', '一起跳舞', '她今晚很迷人'],
        'T台': ['看她走秀，很专业', '她穿着设计师作品', '后台探班，她很开心', '她邀请你参加after party'],
        '摄影棚': ['看她拍摄', '她镜头感很好', '一起拍合照', '她教你摆pose'],
        '时尚派对': ['她穿着礼服很美', '认识时尚圈朋友', '一起品酒聊天', '她介绍你给大家']
      }

      // 获取对应地点的消息，如果没有则使用默认消息
      const locationSpecificMessages = locationMessages[location] || [
        '度过了愉快的时光',
        '你们聊得很开心',
        '她笑容满面',
        '她靠在你的肩膀上',
        '美好回忆让她更加喜欢你了'
      ]

      const fallbackMessage = locationSpecificMessages[Math.floor(Math.random() * locationSpecificMessages.length)]
      const randomMessage = personalizePartnerText(result.scene?.event || fallbackMessage)
      const momentMessage = personalizePartnerText(result.scene?.momentevent || '')
      const locationText = `和${result.datinginfo.dname}在${location}`
      const routeText = result.scene?.rewardtier === 'high-risk' ? '冒险路线' : '稳妥路线'

      let detail = `【${result.scene?.label || '约会场景'} · ${routeText}】${locationText}约会${result.success ? '成功' : '失败'}，名声${result.famechange > 0 ? '+' : ''}${result.famechange}，好感${result.affinitychange > 0 ? '+' : ''}${result.affinitychange}`
      if (result.scene?.momentlabel) detail += `，关系事件：${result.scene.momentlabel}`
      // 记录日志
      addLog(detail, result.success ? 'success' : 'fail')

      // 简单提示
      const messages = []
      if (result.famechange !== 0) messages.push(`名声${result.famechange > 0 ? '+' : ''}${result.famechange}`)
      if (result.healthchange !== 0) messages.push(`健康${result.healthchange > 0 ? '+' : ''}${result.healthchange}`)
      if (result.affinitychange !== 0) messages.push(`好感${result.affinitychange > 0 ? '+' : ''}${result.affinitychange}`)

      if (result.success) {
        const successTitle = result.scene?.rewardtier === 'high-risk' ? '冒险约会成功，高回报！' : '稳妥约会成功！'
        ElMessage.success(`${successTitle}${locationText}，${randomMessage}${momentMessage ? `；${momentMessage}` : ''}；${messages.join('，')} 💚`)
      } else {
        ElMessage.warning(`${routeText}约会失败！${randomMessage} ${messages.join('，')}`)
      }

      // 刷新列表
      await loadDatingInfo()
      if (result.success && dating) {
        openSpouseMoment('date', {
          ...dating,
          daffinitylevel: result.datinginfo?.dstatus || dating.daffinitylevel,
        }, location, result.scene?.moment || 'chat', true)
      }
    } else {
      ElMessage.error(result.msg || '约会失败')
    }
  } catch (err) {
    ElMessage.error('约会失败: ' + err.message)
  }
}
</script>

<style scoped>
.meeting-scene-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meeting-scene-label {
  color: var(--font-secondary);
  font-size: 12px;
}

.meeting-scene-select {
  width: 130px;
}

/* ==================== 主内容区域 ==================== */
.panel-main {
  display: flex;
  gap: 5px;
  margin-top: 10px;
  flex: 1;
  min-height: 0;
}

/* ==================== 左侧：约会对象列表 ==================== */
.panel-main-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.section-content {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.section-content::-webkit-scrollbar {
  display: none;
}

.property-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 5px;
}

/* ==================== 约会卡片 ==================== */
.property-card {
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.property-card.locked {
  opacity: 0.68;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 8px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
}

.card-avatar {
  width: 45px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-image {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid var(--border-color);
  background: var(--panel-color);
  overflow: hidden;
}


.avatar-loading {
  /* 关键代码：开启 Flex 布局实现居中 */
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center;     /* 垂直居中 */
  
  /* 确保容器填满图片区域 */
  width: 100%;
  height: 100%;
  
  /* 背景色和字体样式（可选） */
  background-color: var(--panel-color); /* 主题颜色 */
  font-size: 22px;
  color: var(--border-color);
}

.card-heading {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  min-width: 0;
}

.property-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--font-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-header :deep(.el-tag) {
  flex-shrink: 0;
  max-width: 62px;
}

.scene-condition {
  color: var(--warning-color);
}

.card-body {
  padding: 8px 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.card-body-content {
  flex: 1;
}

.card-description {
  margin-bottom: 4px;
  overflow: hidden;
  color: var(--font-secondary);
  font-size: 10px;
  line-height: 13px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.card-meta {
  font-size: 10px;
  color: var(--font-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 解锁条件或者喜好 */
.unlock-or-like {
  padding: 0;
  border-radius: 6px;
}

.unlock-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--font-secondary);
  margin-bottom: 4px;
}

.pref-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--font-secondary);
  margin-bottom: 4px;
}

.ul-content {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.ul-item {
  display: flex;
  align-items: flex-start;
  gap: 4px;
  font-size: 10px;
  background: var(--panel-color);
  border-radius: 4px;
}

.ul-item.met {
  color: var(--success-color);
}

.ul-icon {
  font-size: 11px;
  flex-shrink: 0;
  color: var(--font-secondary);
}

.ul-value {
  color: var(--font-light);
  flex: 1;
  word-break: break-all;
  display: -webkit-box;
  overflow: hidden;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

/* ==================== 卡片底部操作区 ==================== */
.card-footer {
  margin-top: 6px;
  padding-top: 5px;
  border-top: 1px solid var(--border-color);
}

.card-actions {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(0, 1fr);
  gap: 4px;
}

.card-actions :deep(.el-button) {
  width: 100%;
  height: 26px;
  min-width: 0;
  margin-left: 0;
  padding: 4px;
  font-size: 11px;
}

/* ==================== 右侧关系面板 ==================== */
.panel-main-right {
  width: 260px;
  display: flex;
  flex-direction: column;
  gap: 5px;
  flex-shrink: 0;
}

.relations-box {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
}

.relations-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.relations-content::-webkit-scrollbar {
  display: none;
}

.owned-relations-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.owned-property-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 5px;
  padding: 7px 8px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 5px;
}

.owned-property-card.spouse {
  border-color: var(--warning-color);
}

.owned-property-left {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.owned-avatar {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid var(--border-color);
  border-radius: 50%;
}

.property-details {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.property-meta-info {
  overflow: hidden;
  color: var(--font-secondary);
  font-size: 10px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ==================== 日志 ==================== */
.log-panel {
  height: 140px;
  flex: none;
}

</style>
