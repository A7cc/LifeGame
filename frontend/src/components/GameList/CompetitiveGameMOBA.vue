<template>
  <el-dialog v-model="visible" title="MOBA对战" width="600px" @close="handleClose">
    <div class="game-dialog">
      <div v-if="!gameStarted" class="game-intro">
        <div class="intro-text">5v5团队竞技，配合是制胜关键</div>
        <div class="game-info">
          <span class="info-item">💰 报名费：{{ props.config?.entryCost || 1000 }}元</span>
        </div>
        <div class="hero-select">
          <div class="select-title">选择你的英雄：</div>
          <div class="hero-list">
            <div v-for="hero in heroes" :key="hero.id" class="hero-item" :class="{ selected: selectedHero === hero.id }" @click="selectedHero = hero.id">
              <span class="hero-icon">{{ hero.icon }}</span>
              <span class="hero-name">{{ hero.name }}</span>
              <span class="hero-role">{{ hero.role }}</span>
              <div class="hero-stats">
                <span class="stat">⚔️ {{ hero.attack }}</span>
                <span class="stat">🛡️ {{ hero.defense }}</span>
              </div>
            </div>
          </div>
        </div>
        <el-button type="primary" @click="handleStartGame" style="width: 100%; margin-top: 15px" :disabled="!selectedHero">开始对战</el-button>
      </div>

      <div v-else-if="!gameEnded" class="game-playing">
        <!-- 顶部信息栏 -->
        <div class="top-bar">
          <div class="team-info team-blue">
            <span class="team-name">🔵 我方</span>
            <span class="team-score">{{ blueScore }}</span>
            <div class="team-gold">💰 {{ blueGold }}</div>
          </div>
          <div class="match-info">
            <div class="match-time">{{ formatTime(matchTime) }}</div>
            <div class="match-phase">{{ phase }}</div>
          </div>
          <div class="team-info team-red">
            <div class="team-gold">💰 {{ redGold }}</div>
            <span class="team-score">{{ redScore }}</span>
            <span class="team-name">敌方 🔴</span>
          </div>
        </div>

        <!-- 游戏地图 -->
        <div class="moba-map">
          <!-- 小地图 -->
          <div class="minimap">
            <div class="map-lanes">
              <div class="lane top-lane"></div>
              <div class="lane mid-lane"></div>
              <div class="lane bot-lane"></div>
              <div class="jungle"></div>
            </div>
            <!-- 英雄位置 -->
            <div class="map-hero player-hero" :style="{ left: playerPos.x + '%', top: playerPos.y + '%' }">
              {{ selectedHeroObj?.icon }}
            </div>
            <div v-for="ally in allies" :key="ally.id" class="map-hero ally-hero" :style="{ left: ally.x + '%', top: ally.y + '%' }">
              {{ ally.icon }}
            </div>
            <div v-for="enemy in enemies" :key="enemy.id" class="map-hero enemy-hero" :style="{ left: enemy.x + '%', top: enemy.y + '%' }">
              {{ enemy.icon }}
            </div>
          </div>

          <!-- 主游戏区域 -->
          <div class="main-area">
            <!-- 玩家信息 -->
            <div class="player-panel">
              <div class="hero-portrait">{{ selectedHeroObj?.icon }}</div>
              <div class="hero-details">
                <div class="hero-level">Lv. {{ level }}</div>
                <div class="exp-bar">
                  <div class="exp-fill" :style="{ width: (exp / expToLevel * 100) + '%' }"></div>
                </div>
              </div>
              <div class="hero-stats">
                <div class="stat-bar">
                  <span class="stat-icon">💚</span>
                  <div class="stat-bar-bg">
                    <div class="stat-bar-fill hp" :style="{ width: (hp / maxHp * 100) + '%' }"></div>
                  </div>
                  <span class="stat-text">{{ Math.round(hp) }}/{{ maxHp }}</span>
                </div>
                <div class="stat-bar">
                  <span class="stat-icon">💙</span>
                  <div class="stat-bar-bg">
                    <div class="stat-bar-fill mp" :style="{ width: (mp / maxMp * 100) + '%' }"></div>
                  </div>
                  <span class="stat-text">{{ Math.round(mp) }}/{{ maxMp }}</span>
                </div>
              </div>
            </div>

            <!-- 战斗区域 -->
            <div class="battle-area">
              <div class="battle-log">
                <div v-for="log in battleLogs" :key="log.id" class="log-entry" :class="log.type">
                  {{ log.text }}
                </div>
              </div>

              <!-- 技能区域 -->
              <div class="skills-panel">
                <div class="skill" v-for="(skill, index) in skills" :key="index"
                     :class="{ ready: !skill.cooldown, on_cooldown: skill.cooldown }"
                     @click="useSkill(index)"
                     :title="skill.name">
                  <span class="skill-icon">{{ skill.icon }}</span>
                  <span class="skill-key">{{ ['Q', 'W', 'E', 'R'][index] }}</span>
                  <span v-if="skill.cooldown" class="skill-cd">{{ skill.cooldown }}</span>
                  <span class="skill-cost">{{ skill.cost }}</span>
                </div>
              </div>
            </div>

            <!-- 团队信息 -->
            <div class="team-panel">
              <div class="team-header">🔵 我方阵容</div>
              <div class="team-members">
                <div v-for="member in blueTeam" :key="member.id" class="member" :class="{ dead: member.dead }">
                  <span class="member-icon">{{ member.icon }}</span>
                  <span class="member-kda">{{ member.kills }}/{{ member.deaths }}/{{ member.assists }}</span>
                  <span v-if="member.dead" class="dead-timer">{{ member.respawn }}s</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <el-button @click="moveLane('top')" :disabled="processing" size="small">上路 📍</el-button>
          <el-button @click="moveLane('mid')" :disabled="processing" size="small">中路 📍</el-button>
          <el-button @click="moveLane('bot')" :disabled="processing" size="small">下路 📍</el-button>
          <el-button @click="moveJungle" :disabled="processing" size="small">打野 🌲</el-button>
          <el-button @click="teamFight" :disabled="processing" type="danger" size="small">团战 ⚔️</el-button>
          <el-button @click="recall" :disabled="processing || recalling" type="warning" size="small">回城 🏠</el-button>
        </div>
      </div>

      <div v-else class="game-result">
        <div class="result-icon">{{ resultIcon }}</div>
        <div class="result-title">{{ resultTitle }}</div>
        <div class="result-stats">
          <div class="stat-item">等级：{{ level }}</div>
          <div class="stat-item">K/D/A：{{ kills }}/{{ deaths }}/{{ assists }}</div>
          <div class="stat-item">击杀数：{{ totalKills }}</div>
          <div class="stat-item">金钱：{{ blueGold }}</div>
        </div>
        <div class="result-reward">{{ resultReward }}</div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useMiniGameBase } from '@/src/composables/useMiniGameBase'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  config: {
    type: Object,
    default: () => ({ id: 'moba', name: 'MOBA对战', entryCost: 300 })
  }
})

const emit = defineEmits(['update:modelValue', 'complete'])

const visible = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) resetGameState()
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

const { gameStarted, gameEnded, processing, startGame: startMiniGame, endGame: endMiniGame, reset } = useMiniGameBase(props.config)

const heroes = [
  { id: 1, name: '狂战士', icon: '🪓', role: '战士', attack: 8, defense: 6,
    skills: [
      { name: '旋风斩', icon: '🌀', cost: 40, damage: 30, cooldown: 0, maxCd: 3 },
      { name: '冲撞', icon: '💨', cost: 35, damage: 25, cooldown: 0, maxCd: 4 },
      { name: '狂暴', icon: '😤', cost: 50, damage: 40, cooldown: 0, maxCd: 8 },
      { name: '终结', icon: '💀', cost: 100, damage: 80, cooldown: 0, maxCd: 30 }
    ]
  },
  { id: 2, name: '神射手', icon: '🏹', role: '射手', attack: 9, defense: 3,
    skills: [
      { name: '连射', icon: '🎯', cost: 35, damage: 20, cooldown: 0, maxCd: 2 },
      { name: '陷阱', icon: '🪤', cost: 40, damage: 35, cooldown: 0, maxCd: 6 },
      { name: '闪避', icon: '💨', cost: 30, damage: 15, cooldown: 0, maxCd: 5 },
      { name: '狙击', icon: '🔭', cost: 100, damage: 100, cooldown: 0, maxCd: 25 }
    ]
  },
  { id: 3, name: '大法师', icon: '🧙', role: '法师', attack: 7, defense: 4,
    skills: [
      { name: '火球', icon: '🔥', cost: 45, damage: 35, cooldown: 0, maxCd: 3 },
      { name: '冰冻', icon: '❄️', cost: 40, damage: 25, cooldown: 0, maxCd: 5 },
      { name: '护盾', icon: '🛡️', cost: 50, damage: 10, cooldown: 0, maxCd: 7 },
      { name: '陨石', icon: '☄️', cost: 100, damage: 90, cooldown: 0, maxCd: 35 }
    ]
  },
  { id: 4, name: '圣骑士', icon: '🛡️', role: '辅助', attack: 4, defense: 9,
    skills: [
      { name: '圣光', icon: '✨', cost: 30, damage: 15, cooldown: 0, maxCd: 4 },
      { name: '治疗', icon: '💚', cost: 40, damage: 0, cooldown: 0, maxCd: 6 },
      { name: '祝福', icon: '🙏', cost: 35, damage: 10, cooldown: 0, maxCd: 5 },
      { name: '守护', icon: '⭐', cost: 100, damage: 20, cooldown: 0, maxCd: 40 }
    ]
  },
  { id: 5, name: '暗影者', icon: '🗡️', role: '刺客', attack: 10, defense: 2,
    skills: [
      { name: '背刺', icon: '🔪', cost: 40, damage: 40, cooldown: 0, maxCd: 4 },
      { name: '隐身', icon: '👻', cost: 35, damage: 20, cooldown: 0, maxCd: 8 },
      { name: '暗影', icon: '🌑', cost: 45, damage: 35, cooldown: 0, maxCd: 6 },
      { name: '暗杀', icon: '💀', cost: 100, damage: 120, cooldown: 0, maxCd: 28 }
    ]
  },
]

const selectedHero = ref(null)

// 游戏状态
const blueScore = ref(0)
const redScore = ref(0)
const blueGold = ref(500)
const redGold = ref(500)
const matchTime = ref(0)
const phase = ref('对线期')

// 玩家状态
const level = ref(1)
const exp = ref(0)
const expToLevel = ref(100)
const hp = ref(100)
const maxHp = ref(100)
const mp = ref(100)
const maxMp = ref(100)
const kills = ref(0)
const deaths = ref(0)
const assists = ref(0)
const totalKills = ref(0)
const skills = ref([])

// 位置
const playerPos = ref({ x: 50, y: 75 })
const allies = ref([])
const enemies = ref([])

// 队伍
const blueTeam = ref([])
const redTeam = ref([])

// 其他
const battleLogs = ref([])
const recalling = ref(false)
const resultIcon = ref('')
const resultTitle = ref('')
const resultReward = ref('')

let logId = 0
let gameTimer = null
let cooldownTimer = null

const selectedHeroObj = computed(() => heroes.find(h => h.id === selectedHero.value))

const resetGameState = () => {
  reset()
  selectedHero.value = null
  blueScore.value = 0
  redScore.value = 0
  blueGold.value = 500
  redGold.value = 500
  matchTime.value = 0
  phase.value = '对线期'
  level.value = 1
  exp.value = 0
  expToLevel.value = 100
  hp.value = 100
  maxHp.value = 100
  mp.value = 100
  maxMp.value = 100
  kills.value = 0
  deaths.value = 0
  assists.value = 0
  totalKills.value = 0
  playerPos.value = { x: 50, y: 75 }
  allies.value = []
  enemies.value = []
  blueTeam.value = []
  redTeam.value = []
  battleLogs.value = []
  recalling.value = false
  skills.value = []
  resultIcon.value = ''
  resultTitle.value = ''
  resultReward.value = ''
  if (gameTimer) clearInterval(gameTimer)
  if (cooldownTimer) clearInterval(cooldownTimer)
}

const handleClose = () => {
  visible.value = false
  resetGameState()
}

const handleStartGame = async () => {
  const success = await startMiniGame()
  if (!success) return

  // 初始化技能
  skills.value = JSON.parse(JSON.stringify(selectedHeroObj.value.skills))

  // 初始化队伍
  const allyIcons = ['⚔️', '🏹', '🛡️', '🧙']
  const enemyIcons = ['👹', '👺', '💀', '👾', '😈']

  blueTeam.value = [
    { id: 1, icon: selectedHeroObj.value.icon, kills: 0, deaths: 0, assists: 0, dead: false, respawn: 0 },
    ...allyIcons.map((icon, i) => ({ id: i + 2, icon, kills: 0, deaths: 0, assists: 0, dead: false, respawn: 0 }))
  ]

  redTeam.value = enemyIcons.map((icon, i) => ({
    id: i + 1,
    icon,
    kills: Math.floor(Math.random() * 3),
    deaths: Math.floor(Math.random() * 2),
    assists: Math.floor(Math.random() * 4),
    dead: false,
    respawn: 0
  }))

  // 初始化地图上的单位
  allies.value = [
    { id: 1, icon: '⚔️', x: 45, y: 70 },
    { id: 2, icon: '🏹', x: 55, y: 70 },
    { id: 3, icon: '🛡️', x: 50, y: 80 }
  ]

  enemies.value = [
    { id: 1, icon: '👹', x: 30, y: 30 },
    { id: 2, icon: '👺', x: 50, y: 25 },
    { id: 3, icon: '💀', x: 70, y: 30 },
    { id: 4, icon: '👾', x: 40, y: 40 },
    { id: 5, icon: '😈', x: 60, y: 35 }
  ]

  addLog('比赛开始！前往线上发育', 'info')

  // 游戏计时
  gameTimer = setInterval(() => {
    matchTime.value++
    mp.value = Math.min(maxMp.value, mp.value + 2)

    // 游戏阶段
    if (matchTime.value < 60) phase.value = '对线期'
    else if (matchTime.value < 120) phase.value = '中期'
    else phase.value = '后期'

    // AI行动
    aiAction()

    // 检查游戏结束
    if (blueScore.value >= 5 || redScore.value >= 5) {
      clearInterval(gameTimer)
      clearInterval(cooldownTimer)
      handleEndGame()
    }
  }, 1000)

  // 技能冷却计时
  cooldownTimer = setInterval(() => {
    skills.value.forEach(skill => {
      if (skill.cooldown > 0) skill.cooldown--
    })
  }, 1000)
}

const useSkill = (index) => {
  const skill = skills.value[index]
  if (!skill || skill.cooldown > 0) return
  if (mp.value < skill.cost) {
    ElMessage.warning('法力不足！')
    return
  }

  mp.value -= skill.cost
  skill.cooldown = skill.maxCd

  // 随机命中敌人
  if (Math.random() < 0.6) {
    const damage = skill.damage + (level.value * 5)
    const enemyKilled = dealDamageToEnemy(damage)

    if (enemyKilled) {
      kills.value++
      totalKills.value++
      blueGold.value += 300
      addLog(`你使用${skill.name}击杀了敌人！`, 'kill')
    } else {
      addLog(`你使用${skill.name}造成${damage}伤害`, 'player')
    }
  } else {
    addLog(`${skill.name}未命中`, 'miss')
  }
}

const dealDamageToEnemy = (damage) => {
  // 根据伤害值计算击杀概率
  const killChance = Math.min(0.5, damage * 0.01)
  if (Math.random() < killChance) {
    redGold.value -= 50
    return true
  }
  return false
}

const aiAction = () => {
  // 己方AI行动
  blueTeam.value.forEach((member, index) => {
    if (index === 0 || member.dead) return
    if (Math.random() < 0.15) {
      member.kills++
      totalKills.value++
      blueGold.value += 200
      addLog(`${member.icon}获得击杀！`, 'ally')
    }
  })

  // 敌方AI行动
  redTeam.value.forEach(enemy => {
    if (enemy.dead) return
    if (Math.random() < 0.2) {
      // 对方获得击杀
      const targetIndex = Math.floor(Math.random() * blueTeam.value.length)
      const target = blueTeam.value[targetIndex]

      if (target && !target.dead) {
        if (targetIndex === 0) {
          deaths.value++
          hp.value = 0
          addLog('你被击杀了！', 'death')
        } else {
          target.deaths++
          addLog(`${target.icon}被击杀`, 'enemy')
        }

        enemy.kills++
        redScore.value++
        redGold.value += 300
      }
    }
  })

  // 敌人位置变化
  enemies.value.forEach(enemy => {
    enemy.x = Math.max(20, Math.min(80, enemy.x + (Math.random() - 0.5) * 10))
    enemy.y = Math.max(20, Math.min(80, enemy.y + (Math.random() - 0.5) * 10))
  })

  // 己方位置变化
  allies.value.forEach(ally => {
    ally.x = Math.max(20, Math.min(80, ally.x + (Math.random() - 0.5) * 8))
    ally.y = Math.max(20, Math.min(80, ally.y + (Math.random() - 0.5) * 8))
  })
}

const moveLane = (lane) => {
  if (processing.value) return

  const positions = {
    top: { x: 25, y: 40 },
    mid: { x: 50, y: 50 },
    bot: { x: 75, y: 60 }
  }

  playerPos.value = positions[lane]
  addLog(`前往${lane === 'top' ? '上路' : lane === 'mid' ? '中路' : '下路'}`, 'info')

  // 有几率遭遇敌人
  if (Math.random() < 0.4) {
    encounterEnemy()
  }
}

const moveJungle = () => {
  if (processing.value) return

  playerPos.value = { x: 50, y: 50 }
  addLog('进入野区打野', 'info')

  // 获得金钱和经验
  blueGold.value += 50
  exp.value += 30
  checkLevelUp()

  addLog('击杀野怪，获得50金钱', 'gold')
}

const encounterEnemy = () => {
  addLog('遭遇敌人！', 'warning')

  if (Math.random() < 0.5) {
    const damage = Math.floor(Math.random() * 20) + 10
    hp.value = Math.max(0, hp.value - damage)
    addLog(`受到${damage}点伤害`, 'damage')
  } else {
    const goldGain = Math.floor(Math.random() * 50) + 30
    blueGold.value += goldGain
    addLog(`成功反打，获得${goldGain}金钱`, 'gold')
  }
}

const teamFight = () => {
  if (processing.value) return

  processing.value = true
  addLog('发起团战！', 'warning')

  setTimeout(() => {
    const fightResult = Math.random()

    if (fightResult < 0.45) {
      // 团战胜利
      blueScore.value++
      const killCount = Math.floor(Math.random() * 3) + 1
      totalKills.value += killCount
      blueGold.value += killCount * 300

      addLog(`团战胜利！击杀${killCount}人`, 'kill')
    } else if (fightResult < 0.7) {
      // 团战失败
      redScore.value++
      deaths.value++
      hp.value = 0

      addLog('团战失败！你被击杀', 'death')
    } else {
      // 互换人头
      blueScore.value++
      redScore.value++
      kills.value++
      deaths.value++
      totalKills.value++

      addLog('团战平局，各损失一人', 'info')
    }

    processing.value = false
  }, 1500)
}

const recall = () => {
  if (recalling.value) return

  recalling.value = true
  addLog('正在回城...', 'info')

  setTimeout(() => {
    hp.value = maxHp.value
    mp.value = maxMp.value
    recalling.value = false
    addLog('回城完成，生命值和法力值已恢复', 'heal')
  }, 3000)
}

const checkLevelUp = () => {
  while (exp.value >= expToLevel.value) {
    exp.value -= expToLevel.value
    level.value++
    expToLevel.value = Math.floor(expToLevel.value * 1.5)
    maxHp.value += 20
    maxMp.value += 15
    hp.value = maxHp.value
    mp.value = maxMp.value
    addLog(`升级了！当前等级${level.value}`, 'level')
  }
}

const addLog = (text, type) => {
  battleLogs.value.unshift({
    id: logId++,
    text,
    type
  })
  if (battleLogs.value.length > 8) {
    battleLogs.value.pop()
  }
}

const formatTime = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const handleEndGame = async () => {
  const isWin = blueScore.value > redScore.value
  const winCount = isWin ? 1 : 0

  const resultText = isWin
    ? `获胜，KDA ${kills.value}/${deaths.value}/${assists.value}`
    : '失败'

  const detail = {
    level: level.value,
    kills: kills.value,
    deaths: deaths.value,
    assists: assists.value,
    totalKills: totalKills.value,
    blueGold: blueGold.value
  }

  const gameResult = await endMiniGame(winCount, resultText, detail)

  resultIcon.value = isWin ? '🏆' : '💔'
  resultTitle.value = isWin ? '胜利！' : '失败！'
  resultReward.value = gameResult
    ? `获得 ${gameResult.cashChange} 元，名声 +${gameResult.fameChange}`
    : ''

  emit('complete', {
    game: 'MOBA对战',
    result: resultText,
    ...gameResult
  })
}
</script>

<style scoped>
.game-dialog {
  padding: 10px 0;
}

.game-intro {
  text-align: center;
}

.intro-text {
  font-size: 14px;
  color: var(--font-secondary);
  margin-bottom: 20px;
}

.game-info {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.info-item {
  font-size: 13px;
  padding: 6px 12px;
  background: var(--panel-color);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.hero-select {
  margin-bottom: 15px;
}

.select-title {
  font-size: 13px;
  color: var(--font-secondary);
  margin-bottom: 10px;
}

.hero-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.hero-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 12px;
  background: var(--panel-color);
  border: 2px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  width: 80px;
}

.hero-item:hover {
  border-color: var(--el-color-primary);
  background: #fffbeb;
}

.hero-item.selected {
  border-color: var(--el-color-primary);
  background: #fef3c7;
}

.hero-icon {
  font-size: 28px;
  margin-bottom: 4px;
}

.hero-name {
  font-size: 12px;
  color: var(--font-color);
  font-weight: 600;
  margin-bottom: 2px;
}

.hero-role {
  font-size: 10px;
  color: var(--font-light);
  margin-bottom: 4px;
}

.hero-stats {
  display: flex;
  gap: 4px;
}

.stat {
  font-size: 9px;
  color: var(--font-secondary);
}

.game-playing {
  text-align: center;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background: linear-gradient(135deg, #1e3a5f 0%, #2d3748 100%);
  border-radius: 8px;
  margin-bottom: 10px;
}

.team-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #fff;
}

.team-score {
  font-size: 32px;
  font-weight: 600;
}

.team-gold {
  font-size: 12px;
  color: #ffd700;
}

.match-info {
  text-align: center;
  color: #fff;
}

.match-time {
  font-size: 20px;
  font-weight: 600;
  font-family: monospace;
}

.match-phase {
  font-size: 11px;
  color: #a0aec0;
}

.moba-map {
  background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 10px;
}

.minimap {
  position: relative;
  width: 100%;
  height: 120px;
  background: #0f1419;
  border-radius: 8px;
  margin-bottom: 15px;
  overflow: hidden;
}

.map-lanes {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.lane {
  position: absolute;
  background: rgba(255, 255, 255, 0.05);
}

.top-lane {
  width: 80%;
  height: 20%;
  top: 10%;
  left: 10%;
  border-radius: 4px;
}

.mid-lane {
  width: 20%;
  height: 80%;
  top: 10%;
  left: 40%;
  border-radius: 4px;
}

.bot-lane {
  width: 80%;
  height: 20%;
  top: 70%;
  left: 10%;
  border-radius: 4px;
}

.jungle {
  position: absolute;
  width: 100%;
  height: 100%;
}

.map-hero {
  position: absolute;
  transform: translate(-50%, -50%);
  font-size: 16px;
  transition: all 0.3s;
}

.player-hero {
  z-index: 10;
  filter: drop-shadow(0 0 4px var(--success-color));
}

.ally-hero {
  filter: drop-shadow(0 0 2px rgba(var(--success-color-rgb), 0.5));
}

.enemy-hero {
  filter: drop-shadow(0 0 2px rgba(var(--error-color-rgb), 0.5));
}

.main-area {
  display: flex;
  gap: 15px;
}

.player-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  min-width: 100px;
}

.hero-portrait {
  font-size: 40px;
}

.hero-level {
  font-size: 12px;
  color: #fff;
  font-weight: 600;
}

.exp-bar {
  width: 60px;
  height: 4px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 2px;
  overflow: hidden;
}

.exp-fill {
  height: 100%;
  background: #9333ea;
  transition: width 0.3s;
}

.hero-stats {
  width: 100%;
}

.stat-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 4px;
}

.stat-icon {
  font-size: 12px;
}

.stat-bar-bg {
  flex: 1;
  height: 8px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  overflow: hidden;
}

.stat-bar-fill {
  height: 100%;
  transition: width 0.3s;
}

.stat-bar-fill.hp {
  background: linear-gradient(90deg, #ef4444 0%, #dc2626 100%);
}

.stat-bar-fill.mp {
  background: linear-gradient(90deg, #3b82f6 0%, #2563eb 100%);
}

.stat-text {
  font-size: 9px;
  color: #fff;
  min-width: 50px;
  text-align: right;
}

.battle-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.battle-log {
  flex: 1;
  max-height: 100px;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  padding: 8px;
}

.log-entry {
  font-size: 11px;
  padding: 3px 6px;
  margin-bottom: 2px;
  border-radius: 3px;
  text-align: left;
}

.log-entry.info { color: #a0aec0; }
.log-entry.player { color: var(--success-color); }
.log-entry.kill { color: #ffd700; font-weight: 600; }
.log-entry.ally { color: #60a5fa; }
.log-entry.enemy { color: #f87171; }
.log-entry.death { color: var(--error-color); }
.log-entry.damage { color: #fb923c; }
.log-entry.gold { color: #ffd700; }
.log-entry.warning { color: var(--warning-color); }
.log-entry.heal { color: #34d399; }
.log-entry.level { color: #a78bfa; }
.log-entry.miss { color: #9ca3af; }

.skills-panel {
  display: flex;
  justify-content: center;
  gap: 8px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
}

.skill {
  position: relative;
  width: 48px;
  height: 48px;
  background: rgba(0, 0, 0, 0.5);
  border: 2px solid #4a5568;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.skill.ready {
  border-color: var(--success-color);
  cursor: pointer;
}

.skill.ready:hover {
  background: rgba(var(--success-color-rgb), 0.2);
  transform: scale(1.05);
}

.skill.on_cooldown {
  border-color: #4a5568;
  cursor: not-allowed;
  opacity: 0.6;
}

.skill-icon {
  font-size: 22px;
}

.skill-key {
  position: absolute;
  bottom: 1px;
  left: 4px;
  font-size: 9px;
  color: #a0aec0;
  font-weight: 600;
}

.skill-cost {
  position: absolute;
  top: 1px;
  right: 4px;
  font-size: 9px;
  color: #60a5fa;
  font-weight: 600;
}

.skill-cd {
  position: absolute;
  font-size: 14px;
  color: var(--error-color);
  font-weight: 600;
}

.team-panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  min-width: 100px;
}

.team-header {
  font-size: 11px;
  color: #fff;
  font-weight: 600;
}

.team-members {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.member {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 6px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  position: relative;
}

.member.dead {
  opacity: 0.5;
}

.member-icon {
  font-size: 16px;
}

.member-kda {
  font-size: 10px;
  color: #a0aec0;
}

.dead-timer {
  position: absolute;
  right: 4px;
  font-size: 9px;
  color: var(--error-color);
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
}

.game-result {
  text-align: center;
  padding: 20px 0;
}

.result-icon {
  font-size: 64px;
  margin-bottom: 10px;
}

.result-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 15px;
}

.result-stats {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 15px;
  margin-bottom: 15px;
}

.stat-item {
  font-size: 13px;
  color: var(--font-secondary);
}

.result-reward {
  font-size: 14px;
  color: var(--success-color);
}
</style>
