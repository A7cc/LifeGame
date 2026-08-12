<template>
  <div class="home" data-testid="game-main">
    <!-- 顶部信息栏 -->
    <div class="top-bar">
      <GameTopBar @open-announcement="showAnnouncement = true" @exit="emit('exit')" @open-save="showSaveDialog = true" />
    </div>

    <!-- 主体布局 -->
    <div class="main">
      <!-- 左侧菜单栏 -->
      <div class="menu">
        <GameMenu />
      </div>

      <!-- 中间内容及底部状态栏 -->
      <div class="menu-content">
        <!-- 中间内容，路由视图 -->
        <div class="game-content" data-testid="game-content">
          <router-view />
        </div>
        <!-- 底部 -->
        <div class="bottom-bar">
          <GameBottomBar @updateGameData="updateGameData" @exit="emit('exit')"/>
        </div>
      </div>
    </div>
  </div>
  <!-- 公告弹窗 -->
  <DialogAnnouncement v-model="showAnnouncement" :announcements="announcements" />
  <!-- 存档弹窗 -->
  <DialogSave v-model="showSaveDialog" @loaded="handleSaveLoaded" />
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useGameStore } from '@/src/stores/game'
import { useStockTicker } from '@/src/composables/useStockTicker'
import DialogAnnouncement from './Dialog/DialogAnnouncement.vue'
import DialogSave from './Dialog/DialogSave.vue'
import GameTopBar from './GameTopBar.vue'
import GameBottomBar from './GameBottomBar.vue'
import GameMenu from './GameMenu.vue'

// 用pinia存放数据
const gameStore = useGameStore()
useStockTicker()
// 公告信息直接跟随 store，避免复制状态后过期。
const announcements = computed(() => gameStore.announcements)
// 控制公告显示
const showAnnouncement = ref(false)
// 控制存档对话框显示
const showSaveDialog = ref(false)

// 初始化时从 store 获取数据
onMounted(() => {
  // 加载完自动弹出公告
  showAnnouncement.value = true
})

// 更新父组件中的gamedata
const updateGameData = (data) => {
  gameStore.applyGameData(data)
  // 加载完自动弹出公告
  showAnnouncement.value = true
}

// 存档加载完成后的处理
const handleSaveLoaded = (data) => {
  updateGameData(data)
}

// 退出
const emit = defineEmits(['exit'])
</script>

<style scoped>
/* 页面 */
.home {
  height: 100vh;                        /* 高度占据整个可视窗口 */
  display: flex;                        /* 启用 Flex 弹性布局 */
  flex-direction: column;              /* 主轴为垂直方向，从上往下排列子元素 */
  background: var(--background-color);                 /* 使用主题变量 */
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; /* 设置字体族优先级 */
  animation: fadeIn 0.6s ease-in-out;  /* 页面加载时淡入动画，持续 0.6 秒，缓动效果为 ease-in-out */
}
/* top-bar */
.top-bar {
  display: flex;                        /* 使用 Flex 布局 */
  justify-content: space-between;      /* 子元素两端对齐，中间自动分配空隙 */
  align-items: center;                 /* 子元素垂直居中 */
	background: var(--gradient-primary);
	border-bottom: 1px solid rgb(255 255 255 / 18%);
	box-shadow: 0 2px 10px var(--panel-shadow);
  padding: 0 20px;                     /* 上下为 0，左右内边距为 20px */
  height: 60px;                        /* 固定高度 60px */
}

.main {
  display: flex;                      /* Flex 布局 */
  flex: 1;                            /* 占据剩余空间 */
  overflow: hidden;                  /* 隐藏超出容器的内容（防止滚动） */
}

/* 菜单栏 */
.menu {
  width: 180px;                       /* 菜单栏宽度 */
  height: 100%;                       /* 高度占满父容器 */
	background: var(--panel-main-color);
	border-right: 1px solid var(--border-color);
  transition: all 0.3s ease-in-out;  /* 所有属性过渡效果 */
  overflow-y: auto; /* 新增：允许纵向滚动 */
  /* 移除滚动条预留空间 */
  scrollbar-width: none;              /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;           /* IE/Edge 浏览器隐藏滚动条 */
}

.menu::-webkit-scrollbar {
  display: none;          /* Chrome/Safari 隐藏滚动条 */
}

.menu-content {
  display: flex;                      /* Flex 布局 */
  flex-direction: column;            /* 垂直排列子元素 */
  flex: 1;                            /* 占据剩余空间 */
	background: var(--background-color);
  justify-content: space-between;   /* 子元素垂直两端对齐 */
}

.game-content {
  display: flex;                      /* Flex 横向布局 */
  padding: 10px;                      /* 四周内边距 */
  gap: 5px;                          /* 子元素间间距 */
  flex: 1;                            /* 占据剩余空间 */
  overflow-y: auto;                   /* 垂直方向滚动 */
  /* 移除滚动条预留空间 */
  scrollbar-width: none;              /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;           /* IE/Edge 浏览器隐藏滚动条 */
}

.game-content::-webkit-scrollbar {
  display: none;          /* Chrome/Safari 隐藏滚动条 */
}

/* 底部栏美化 */
.bottom-bar {
  display: flex;                              /* Flex 布局 */
  justify-content: space-between;            /* 两端对齐 */
	border-top: 1px solid rgb(255 255 255 / 18%);
  font-size: 15px;                            /* 字体大小 */
  background: var(--gradient-primary); /* 使用主题变量 */
  background-color: var(--primary-color);    /* 备用背景色 */
	color: var(--on-primary-color);
  padding: 12px 24px;                         /* 内边距 */
  font-weight: 600;                           /* 字体加粗 */
	border-top-left-radius: 12px;
	border-top-right-radius: 12px;
	box-shadow: 0 -2px 10px var(--panel-shadow);
}


/* 动画 keyframes */
@keyframes fadeIn {
  0% { opacity: 0; transform: scale(0.96); }  /* 初始：透明+略缩小 */
  100% { opacity: 1; transform: scale(1); }   /* 结束：完全显示 */
}

@keyframes slideIn {
  from { transform: translateX(-10px); opacity: 0; }
  to   { transform: translateX(0); opacity: 1; }
}
</style>
