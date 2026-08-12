<template>
  <el-dialog
    data-testid="save-dialog"
    v-model="visible"
    title="存档管理"
    width="350px"
    align-center
    class="save-dialog"
    @open="loadSaves"
  >
    <div class="save-container">
      <!-- 左侧：存档列表 -->
      <div class="save-list-section">
        <div class="section-header">
          <span class="section-title">存档列表</span>
          <div class="save-tip">
          <span v-if="saveList.length >= 10" class="warning">存档已达上限（10个），请删除旧存档</span>
          <span v-else>{{ saveList.length }} / 10 个存档</span>
        </div>
        </div>
        <div class="save-list" v-loading="loading">
          <div
            v-for="save in saveList"
            :key="save.id"
            :data-testid="`save-item-${save.id}`"
            class="save-item"
            :class="{ selected: selectedSave?.id === save.id }"
            @click="selectSave(save)"
          >
            <div class="save-info">
              <div class="save-name">{{ save.name }}</div>
              <div class="save-meta">
                <span class="save-age">{{ save.game_year }}岁</span>
                <span class="save-time">{{ save.created_at }}</span>
              </div>
            </div>
            <div class="save-actions">
              <el-button :data-testid="`load-save-${save.id}`" type="primary" size="small" @click.stop="handleLoad(save.id)">加载</el-button>
              <el-button :data-testid="`delete-save-${save.id}`" type="danger" size="small" @click.stop="handleDelete(save.id)">删除</el-button>
            </div>
          </div>
          <div v-if="saveList.length === 0 && !loading" class="no-saves">
            暂无存档，请保存当前进度
          </div>
        </div>
      </div>

      <!-- 右侧：新建存档 -->
      <div>
        <div class="create-form">
          <el-input
            data-testid="save-name"
            v-model="newSaveName"
            placeholder="输入存档名称"
            maxlength="10"
            show-word-limit
          />
          <el-button data-testid="create-save" type="success" :disabled="!newSaveName.trim()" @click="handleSave">
            保存当前进度
          </el-button>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { SaveGame, LoadGame, ListSaves, DeleteSave } from "@/wailsjs/go/services/App.js"
import { useGameStore } from '@/src/stores/game'

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue', 'loaded'])

const gameStore = useGameStore()

// 对话框可见性
const visible = ref(props.modelValue)
watch(() => props.modelValue, val => visible.value = val)
watch(visible, val => emit('update:modelValue', val))

// 存档列表
const saveList = ref([])
const loading = ref(false)
const selectedSave = ref(null)

// 新存档名称
const newSaveName = ref('')

// 加载存档列表
const loadSaves = async () => {
  loading.value = true
  try {
    const result = await ListSaves()
    if (result.code === 200) {
      saveList.value = result.saves || []
    } else {
      ElMessage.error(result.msg || '获取存档列表失败')
    }
  } catch (err) {
    console.error('获取存档列表失败:', err)
    ElMessage.error('获取存档列表失败')
  } finally {
    loading.value = false
  }
}

// 选择存档
const selectSave = (save) => {
  selectedSave.value = save
}

// 保存当前进度
const handleSave = async () => {
  if (!newSaveName.value.trim()) {
    ElMessage.warning('请输入存档名称')
    return
  }
  if (saveList.value.length >= 10) {
    ElMessage.warning('存档已达上限，请删除旧存档')
    return
  }
  try {
    const result = await SaveGame(newSaveName.value.trim())
    if (result.code === 200) {
      ElMessage.success('存档成功')
      newSaveName.value = ''
      loadSaves()
    } else {
      ElMessage.error(result.msg || '存档失败')
    }
  } catch (err) {
    console.error('存档失败:', err)
    ElMessage.error('存档失败')
  }
}

// 加载存档
const handleLoad = async (saveId) => {
  try {
    const result = await LoadGame(saveId)
    if (result.code === 200) {
      gameStore.applyGameData(result)
      ElMessage.success('加载成功')
      emit('loaded', result)
      visible.value = false
    } else {
      ElMessage.error(result.msg || '加载失败')
    }
  } catch (err) {
    console.error('加载存档失败:', err)
    ElMessage.error('加载存档失败')
  }
}

// 删除存档
const handleDelete = async (saveId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个存档吗？', '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const result = await DeleteSave(saveId)
    if (result.code === 200) {
      ElMessage.success('删除成功')
      if (selectedSave.value?.id === saveId) {
        selectedSave.value = null
      }
      loadSaves()
    } else {
      ElMessage.error(result.msg || '删除失败')
    }
  } catch (err) {
    if (err !== 'cancel') {
      console.error('删除存档失败:', err)
      ElMessage.error('删除存档失败')
    }
  }
}
</script>

<style scoped>
.save-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 左侧存档列表 */
.save-list-section {
  flex: 1;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--font-color);
}

.save-list {
  height: 100px;
  overflow-y: auto;
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.save-list::-webkit-scrollbar {
  display: none;
}

.save-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-radius: 6px;
  background: var(--panel-color);
  margin-bottom: 5px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid var(--border-color);
}

.save-item:last-child {
  margin-bottom: 0;
}

.save-item:hover {
  background-color: var(--select-color);         /* 根据主题变化 */
}

.save-info {
  flex: 1;
}

.save-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--font-color);
  margin-bottom: 4px;
}

.save-meta {
  font-size: 12px;
  color: var(--font-secondary);
}

.save-age {
  margin-right: 10px;
}

.save-actions {
  display: flex;
  flex-direction: column;
  gap: 5px;
  align-items: flex-end;
}

.save-actions .el-button {
  width: 60px;
}

.no-saves {
  text-align: center;
  padding: 20px;
  font-size: 13px;
}

/* 右侧新建存档 */
.section-title {
  margin-bottom: 12px;
}

.create-form {
  display: flex;
  flex-direction: row;
  gap: 8px;
}

.save-tip {
  font-size: 12px;
  color: var(--font-secondary);
}

.warning {
  color: var(--warning-color);
}
</style>
