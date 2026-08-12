<template>
  <el-dialog
    v-model="visible"
    width="550px"
    align-center
    :show-close="false"
    class="modern-dialog"
  >
    <template #header>
      <div class="dialog-header">
        <div class="dialog-title">📋 公告</div>
        <el-select
          v-model="activeType"
          class="type-select"
          style="width: 180px;"
          popper-class="custom-select-dropdown"
        >
          <el-option
            v-for="(label, type) in marketLabels"
            :key="type"
            :label="label"
            :value="type"
          />
        </el-select>
      </div>
    </template>

    <div class="dialog-body">
      <div class="announcement-box">
        <div
          v-for="(item, index) in currentAnnouncements"
          :key="index"
          class="announcement-item"
        >
          {{ item }}
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  announcements: Object,
})

// 映射市场类型为中文标签
const marketLabels = {
  'announceins': '国内市场',
  'announceout': '国外市场',
  'announcecompany': '创业公告',
  'announcegame': '游戏公告',
  'announcehealthy': '身体健康情况',
}

const emit = defineEmits(['update:modelValue'])

const visible = ref(props.modelValue)
const activeType = ref('announceins')  // ✅ 默认值改成内部键名

watch(() => props.modelValue, val => {
  visible.value = val
})

watch(visible, val => {
  emit('update:modelValue', val)
})

const currentAnnouncements = computed(() => {
  return props.announcements[activeType.value] || []
})
</script>

<style scoped>
/* 对话框整体样式 */
.modern-dialog {
  border-radius: 12px;
  background-color: var(--panel-color);
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.12);
}

/* 标题区域 */
.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end; /* 原为 center，现在改为底部对齐 */
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-title {
  font-size: 22px;
  font-weight: bold;
  color: var(--font-color);
  padding-bottom: 2px; /* 可根据字体微调标题基线对齐 */
}

/* 美化下拉框 */
.type-select :deep(.el-input) {
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--panel-color);
  box-shadow: 0 2px 4px rgba(0,0,0,0.06);
  transition: box-shadow 0.3s ease;
}
.type-select :deep(.el-input:hover) {
  box-shadow: 0 3px 6px rgba(0,0,0,0.12);
}


/* 自定义下拉菜单浮层样式 */
::v-deep(.custom-select-dropdown) {
  border-radius: 10px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  font-weight: 500;
  padding: 4px 0;
}

/* 内容区域 */
.dialog-body {
  padding: 20px 24px;
  border-radius: 0 0 12px 12px;
}

/* 公告列表容器 */
.announcement-box {
  padding: 16px;
  border: 1px solid var(--border-color);
  background-color: var(--panel-color);
  border-radius: 10px;
  box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.03);
  max-height: 280px;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.announcement-box::-webkit-scrollbar {
  display: none;
}

/* 公告条目样式 */
.announcement-item {
  padding: 10px 12px;
  margin-bottom: 8px;
  border-left: 4px solid #409EFF;
  background: var(--panel-color);
  border-radius: 6px;
  font-size: 15px;
  line-height: 1.6;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  transition: background 0.3s;
}
.announcement-item:hover {
  background-color: var(--panel-color);
}
.announcement-item:last-child {
  margin-bottom: 0;
}
</style>
