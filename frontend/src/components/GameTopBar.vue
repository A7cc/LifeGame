<template>
  <div class="top-left" data-testid="player-summary">
    <span class="game-title">LifeGame</span>
    <span class="user-name">{{ gameStore.userInfo?.uname }}（{{ gameStore.userInfo?.usex ? '男' : '女' }}）</span>
    <span class="user-info">{{ gameStore.userInfo?.uage }}岁</span>
  </div>
  <div class="top-right">
    <el-button data-testid="open-announcements" type="info" size="small" @click="openAnnouncement">📊 市场行情</el-button>
    <el-button data-testid="open-save-dialog" type="success" size="small" @click="emit('open-save')">💾 存档</el-button>
    <el-button
      data-testid="background-music-toggle"
      size="small"
      class="music-toggle"
      :aria-pressed="musicEnabled"
      :title="musicEnabled && !musicPlaying ? '播放背景音乐' : (musicEnabled ? '关闭背景音乐' : '开启背景音乐')"
      @click="toggleMusic"
    >
      {{ musicEnabled ? (musicPlaying ? '🎵 音乐' : '▶️ 播放') : '🔇 静音' }}
    </el-button>
    <el-button data-testid="theme-toggle" color="#b45309"
      :icon="themeStore.theme === 'light' ? Moon : Sunny"
      size="small"
      @click="themeStore.toggleTheme()"
      :title="themeStore.theme === 'light' ? '切换到暗黑模式' : '切换到亮色模式'"
    >
      {{ themeStore.theme === 'light' ? '暗黑' : '亮色' }}
    </el-button>
    <el-button data-testid="end-game" type="danger" size="small" @click="$emit('exit')">提前退休</el-button>
  </div>
</template>

<script setup>
import { useGameStore } from '@/src/stores/game'
import { useThemeStore } from '@/src/stores/theme'
import { useBackgroundMusic } from '@/src/composables/useBackgroundMusic'
import { Moon, Sunny } from '@element-plus/icons-vue'

const gameStore = useGameStore()
const themeStore = useThemeStore()
const { enabled: musicEnabled, playing: musicPlaying, toggle: toggleMusic } = useBackgroundMusic()

// 退出和弹窗
const emit = defineEmits(['exit', 'open-announcement', 'open-save'])

const openAnnouncement = () => {
  emit('open-announcement')
}
</script>

<style scoped>
.top-left {
  display: flex;
  align-items: center;
  gap: 5px;
}

.top-right {
  display: flex;
  gap: 5px;
}

.music-toggle {
  min-width: 72px;
}

.game-title {
  font-size: 22px;
  font-weight: bold;
  color: #ffffff;
  margin-right: 10px;
}

.user-name {
  font-weight: bold;
  color: rgba(255, 255, 255, 0.95);;
}

.user-info {
  color: rgba(255, 255, 255, 0.95);;
}
</style>
