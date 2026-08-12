/**
 * ============================================================
 * LifeGame 人生模拟器 - 应用入口文件
 * ============================================================
 *
 * 本文件是 Vue 应用的入口点，负责：
 * 1. 创建 Vue 应用实例
 * 2. 注册全局插件（Element Plus、Pinia、Router）
 * 3. 挂载应用到 DOM
 *
 * ============================================================
 */

// ==================== 核心依赖导入 ====================
import {createApp} from 'vue'           // Vue 3 创建应用函数
import { createPinia } from 'pinia'     // Pinia 状态管理库
import {
  ElButton,
  ElButtonGroup,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElImage,
  ElImageViewer,
  ElInput,
  ElInputNumber,
  ElMenu,
  ElMenuItem,
  ElOption,
  ElRadioButton,
  ElRadioGroup,
  ElSelect,
  ElSlider,
  ElTable,
  ElTableColumn,
  ElTag,
  ElLoading,
} from 'element-plus'

// ==================== 应用组件与配置 ====================
import App from './App.vue'             // 根组件
// 只加载实际使用组件的样式，避免把完整 Element Plus 样式表打入首屏。
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/button-group/style/css'
import 'element-plus/es/components/descriptions/style/css'
import 'element-plus/es/components/dialog/style/css'
import 'element-plus/es/components/image/style/css'
import 'element-plus/es/components/input/style/css'
import 'element-plus/es/components/input-number/style/css'
import 'element-plus/es/components/menu/style/css'
import 'element-plus/es/components/menu-item/style/css'
import 'element-plus/es/components/radio-button/style/css'
import 'element-plus/es/components/radio-group/style/css'
import 'element-plus/es/components/select/style/css'
import 'element-plus/es/components/slider/style/css'
import 'element-plus/es/components/table/style/css'
import 'element-plus/es/components/table-column/style/css'
import 'element-plus/es/components/tag/style/css'
import 'element-plus/es/components/loading/style/css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'
import 'element-plus/theme-chalk/dark/css-vars.css' // Element Plus 暗黑主题样式
import router from './router'           // Vue Router 路由配置
import './style-variables.css'          // 全局样式变量（含暗黑主题）

// ==================== 应用初始化 ====================
// 创建 Vue 应用实例
const app = createApp(App)

// ==================== 插件注册 ====================
const elementComponents = [
  ElButton,
  ElButtonGroup,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElImage,
  ElImageViewer,
  ElInput,
  ElInputNumber,
  ElMenu,
  ElMenuItem,
  ElOption,
  ElRadioButton,
  ElRadioGroup,
  ElSelect,
  ElSlider,
  ElTable,
  ElTableColumn,
  ElTag,
]

elementComponents.forEach(component => app.use(component))
app.use(ElLoading)
app.use(createPinia())  // 注册 Pinia 状态管理
app.use(router)         // 注册 Vue Router 路由

// ==================== 应用挂载 ====================
// 将应用挂载到 index.html 中 id 为 "app" 的 DOM 元素上
app.mount('#app')
