<template>
  <div class="product-panel" data-testid="page-company">
    <!-- 顶部标题与控制按钮 -->
    <div class="panel-header">
      <div class="panel-title">创业</div>
      <div class="panel-controls">
        <el-tag size="small">{{ Object.keys(gameStore.userInfo.ucompany || {}).length }}/{{ gameStore.gameInfo.gmaxholdnum.mcholdnum }}</el-tag>
      </div>
    </div>
    
    <!-- 表头单独放置，固定 -->
    <div class="header-row fixed-header">
      <div class="product-title">公司</div>
      <div class="product-title">单价</div>
      <div class="product-title">风险程度</div>
      <div class="product-title">盈利百分比</div>
      <div class="product-title">创业时间</div>
      <div class="product-title">成本</div>
      <div class="product-title">持有量（手）</div>
      <div class="product-action"></div>
    </div>

    <!-- 滚动区域只包含数据行 -->
    <div class="panel-body scrollable-body">
      <template v-for="(item, index) in gameStore.gameInfo.gcompanyinfo" :key="index">
          <div class="product-row">
          <div class="product-name">{{ item.ciname }}</div>
          <div class="product-other">{{ item.ciprice }}</div>
          <div class="product-other">{{ item.cirisk }}</div>
          <div class="product-other">{{ item.ciprofit }}%</div>
          <div class="product-other">{{ item.citime }}</div>
          <div class="product-other">{{ Object.keys(gameStore.userInfo.ucompany).length == 0 ? 0 : gameStore.userInfo.ucompany[index]?.ucompanycostprice || 0 }}</div>
          <div class="product-other">{{ Object.keys(gameStore.userInfo.ucompany).length == 0 ? 0 : gameStore.userInfo.ucompany[index]?.ucompanynum || 0 }}</div>
            <div class="product-action">
              <el-button v-if="!gameStore.userInfo.ucompany?.[index] && item.cistatus" size="small" type="primary" :disabled="Object.keys(gameStore.userInfo.ucompany).length >= 3" style="width: 100%;" @click="DbsUseItemfunc(index, item, '购买', '创业')">创业</el-button>
              <template v-else-if="item.cistatus">
                <el-button size="small" type="success" style="width: 48%;" @click="DbsUseItemfunc(index, item, '购买', '创业')">投</el-button>
                <el-button size="small" type="danger" :disabled="gameStore.userInfo.ucompany[index]?.ucompanyholdtime < item.citime" style="width: 48%;" @click="DbsUseItemfunc(index, item, '出售', '创业')">出</el-button>
              </template>
              <template v-else>
                <el-button size="small" type="info" disabled style="width: 100%;">破产</el-button>
              </template>
            </div>
          </div>
      </template>
    </div>
  </div>
  <!--买卖弹窗 -->
  <DialogBugSell ref="DbsUseItem"/>
</template>

<script setup>
import { ref } from 'vue'
import { useGameStore } from '@/src/stores/game'
import DialogBugSell from '@/src/components/Dialog/DialogBugSell.vue';

const gameStore = useGameStore()

// 调用子组件的useItem方法
const DbsUseItem = ref(null)

// 调用弹窗子组件的useItem方法
const DbsUseItemfunc = (index, item, type, region) => {
  DbsUseItem.value?.useItem(index, item.ciname, type, region)
}

</script>

<style scoped>
.product-row {
  display: flex;                              /* Flex 布局 */
  align-items: center;                        /* 垂直居中 */
  justify-content: space-between;            /* 子元素两端对齐 */
  background-color: var(--panel-color);               /* 半透明白背景 */
  border-radius: 8px;                         /* 圆角 */
  padding: 8px 12px;                          /* 内边距 */
  box-shadow: 0 1px 4px var(--panel-color);      /* 轻微阴影 */
  font-size: 14px;                            /* 字号 */
  border: 1px solid var(--border-color); /* 边框 */
}


.product-name {  
  background-color: var(--product-name-color); /* 使用主题变量 */
  border: 1px solid var(--product-name-border-color); /* 边框 */
}

.product-title {
  font-weight: bold;                         /* 加粗字体 */
  color: var(--font-color);            /* 使用主题变量 */
}

.product-other, .product-title {
  width: 60px;                                /* 固定宽度 */
}

.product-name,
.product-other {
  font-size: 14px;                            /* 字号 */
  color: var(--font-color);             /* 使用主题变量 */
}

.product-name,
.product-other,
.product-title {
  flex: 1;                                    /* 平均分配宽度 */
  text-align: center;                         /* 文本居中 */
  padding: 5px 0px;                          /* 内边距 */
  border-radius: 6px;                         /* 圆角 */
}

.product-action {
  display: flex;                              /* Flex 布局 */
  justify-content: center;                    /* 居中对齐 */
  align-items: center;                        /* 垂直居中 */
  width: 80px;                                 /* 固定宽度 */
}

/* 按钮统一美化 */
.el-button {
  border-radius: 8px;                         /* 圆角 */
  transition: all 0.2s ease-in-out;           /* 动画过渡 */
}

.el-button:hover {
  transform: scale(1.05);                     /* 鼠标悬浮时略微放大 */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);   /* 增加阴影 */
}

.top-right .el-button {
  background-color: transparent;              /* 背景透明 */
  border: none;                               /* 去掉边框 */
  color: white;                               /* 白色字体 */
}

.top-right .el-button:hover {
  background-color: rgba(255, 255, 255, 0.2);  /* 鼠标悬浮时背景微变 */
}

/* 表头行样式，通常用于列表头部固定项 */
.header-row {
  background-color: var(--accent-color); /* 使用主题变量 */
  color: var(--font-secondary);             /* 使用主题变量 */
}

/* 面板主体区域，支持垂直滚动，适用于显示多个条目 */
.panel-body {
  flex: 1;                                       /* 占满父容器剩余高度 */
  overflow-y: auto;                              /* 垂直方向滚动 */
  display: flex;                                 /* 使用 Flex 布局 */
  flex-direction: column;                        /* 子元素垂直排列 */
  gap: 5px;                                      /* 子元素间间距 */
  /* 移除滚动条预留空间 */
  scrollbar-width: none;                         /* Firefox 浏览器隐藏滚动条 */
  -ms-overflow-style: none;                      /* IE/Edge 浏览器隐藏滚动条 */
  }
  .panel-body::-webkit-scrollbar {
  display: none;                                 /* Chrome/Safari 隐藏滚动条 */
}

/* 固定头部区域样式，常用于列表顶部显示关键字段 */
.fixed-header {
  display: flex;                                 /* Flex 布局 */
  justify-content: space-between;               /* 两端对齐 */
  padding: 8px 10px;                             /* 内边距 */
  font-weight: bold;                             /* 加粗 */
  border-bottom: 2px solid var(--border-color);                 /* 底部分隔线 */
  border-radius: 4px;                            /* 轻微圆角 */
  z-index: 1;                                    /* 保证其处于上层，防止被滚动内容遮挡 */
}

/* 可滚动内容区，通常与 fixed-header 配套使用 */
.scrollable-body {
  flex: 1;                                       /* 占据剩余空间 */
  overflow-y: auto;                              /* 支持垂直滚动 */
  margin-top: 8px;                               /* 顶部留白，避免与固定头部重叠 */
  display: flex;                                 /* Flex 垂直布局 */
  flex-direction: column;
  gap: 5px;                                      /* 子项间距 */
}
</style>
