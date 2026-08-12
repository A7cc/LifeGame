<template>
  <div class="product-panel" :data-testid="testId">
    <div class="panel-header">
      <div class="panel-title">{{ title }}</div>
    </div>

    <div class="header-row fixed-header">
      <div class="product-title">货物名</div>
      <div class="product-title">价格</div>
      <div class="product-action"></div>
    </div>

    <div class="panel-body scrollable-body">
      <template v-if="visibleMarketItems.length > 0">
        <div class="product-row" v-for="item in visibleMarketItems" :key="item.id" :data-testid="`market-buy-row-${region}-${item.id}`">
          <div class="product-name">{{ item.iiname }}</div>
          <div class="product-other">{{ item.iiprice }}</div>
          <div class="product-action">
            <el-button :data-testid="`market-buy-${region}-${item.id}`" size="small" type="primary" @click="openTradeDialog(item.id, item, '购买')">购买</el-button>
          </div>
        </div>
      </template>
      <div v-else class="empty-tip">暂无可交易货物</div>
    </div>
  </div>

  <div class="product-panel" :data-testid="`${testId}-owned`">
    <div class="panel-header">
      <div class="panel-title">拥有货物 ({{ userItemCount }}/{{ userItemLimit }})</div>
    </div>

    <div class="header-row fixed-header">
      <div class="product-title">货物名</div>
      <div class="product-title">当前价格</div>
      <div class="product-title">买入价格</div>
      <div class="product-title">数量</div>
      <div class="product-action"></div>
    </div>

    <div class="panel-body scrollable-body">
      <template v-if="userItems.length > 0">
        <div class="product-row" v-for="item in userItems" :key="item.id" :data-testid="`market-owned-row-${region}-${item.id}`">
          <div
            class="product-name"
            :style="{
              backgroundColor: item.buyprice == item.iiprice ? 'var(--product-name-color)' : item.buyprice < item.iiprice ? 'var(--error-secondary-color)' : 'var(--success-secondary-color)',
              borderColor: item.buyprice == item.iiprice ? 'var(--product-name-border-color)' : item.buyprice < item.iiprice ? '#ffd0d0' : '#d3ffd0'
            }"
          >
            {{ item.iiname }}
          </div>
          <div class="product-other">{{ item.iiprice }}</div>
          <div class="product-other">{{ item.buyprice }}</div>
          <div class="product-other">{{ item.num }}</div>
          <div class="product-action">
            <el-button :data-testid="`market-sell-${region}-${item.id}`" size="small" type="primary" @click="openTradeDialog(item.id, item, '出售')">出售</el-button>
          </div>
        </div>
      </template>
      <div v-else class="empty-tip">暂未持有货物</div>
    </div>
  </div>
  <!-- 买卖弹窗 -->
  <DialogBugSell ref="tradeDialogRef"/>
</template>

<script setup>
import { computed, ref } from 'vue'
import DialogBugSell from '@/src/components/Dialog/DialogBugSell.vue'

const props = defineProps({
  testId: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  region: {
    type: String,
    required: true,
  },
  marketItems: {
    type: [Array, Object],
    default: () => ({}),
  },
  userItems: {
    type: Array,
    default: () => [],
  },
  userItemCount: {
    type: Number,
    default: 0,
  },
  userItemLimit: {
    type: Number,
    default: 0,
  },
})

const tradeDialogRef = ref(null)

const visibleMarketItems = computed(() => Object.entries(props.marketItems || {})
  .map(([id, item]) => ({
    ...item,
    id: Number(id),
  }))
  .filter(item => item.iidisplay))

// 统一由复用面板打开买卖弹窗，两个市场页只负责传入数据。
const openTradeDialog = (index, item, type) => {
  tradeDialogRef.value?.useItem(index, item.iiname, type, props.region)
}
</script>

<style scoped>
/* 产品面板左侧 */
.product-panel:nth-child(1) {
  flex: 1 1 40%;   /* 调整宽度 */
}

/* 产品面板右侧 */
.product-panel:nth-child(2) {
  flex: 1 1 60%;   /* 调整宽度 */
}

.product-row {
  display: flex;                              /* Flex 布局 */
  align-items: center;                        /* 垂直居中 */
  justify-content: space-between;            /* 子元素两端对齐 */
  background-color: var(--panel-color); /* 使用主题变量 */
  border-radius: 8px;                         /* 圆角 */
  padding: 8px 12px;                          /* 内边距 */
  box-shadow: 0 1px 4px var(--panel-color); /* 使用主题变量 */
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

.empty-tip {
  color: var(--font-secondary);
  font-size: 14px;
  padding: 24px 0;
  text-align: center;
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
  border-bottom: 2px solid var(--border-color); /* 使用主题变量 */
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
