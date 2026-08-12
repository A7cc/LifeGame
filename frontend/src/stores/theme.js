/**
 * ============================================================
 * LifeGame 人生模拟器 - 主题状态管理 Store
 * ============================================================
 *
 * 本文件负责管理应用的主题状态，支持明暗主题切换
 * - 使用 Pinia 进行状态管理
 * - 使用 localStorage 持久化主题偏好
 * - 使用 data-theme 属性控制主题
 *
 * ============================================================
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  // ==================== 状态定义 ====================
  // 当前主题：'light' 或 'dark'
  const theme = ref('light')
  let mediaQuery = null
  let mediaChangeHandler = null

  // ==================== 主题应用 ====================
  /**
   * 应用主题到 DOM
   * - 设置 html 元素的 data-theme 属性
   * - 为 Element Plus 添加/移除 dark 类
   */
  const applyTheme = (newTheme) => {
    const html = document.documentElement

    if (newTheme === 'dark') {
      html.setAttribute('data-theme', 'dark')
      html.classList.add('dark')
    } else {
      html.setAttribute('data-theme', 'light')
      html.classList.remove('dark')
    }
  }

  // ==================== 主题保存 ====================
  /**
   * 保存主题偏好到 localStorage
   */
  const saveTheme = (newTheme) => {
    localStorage.setItem('lifegame-theme', newTheme)
  }

  // ==================== 主题切换 ====================
  /**
   * 切换主题（亮色 <-> 暗黑）
   */
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    applyTheme(theme.value)
    saveTheme(theme.value)
  }

  // ==================== 主题初始化 ====================
  /**
   * 初始化主题
   * - 优先从 localStorage 读取用户偏好
   * - 其次跟随系统偏好
   * - 默认使用亮色主题
   */
  const initTheme = () => {
    // 尝试从 localStorage 获取保存的主题
    const savedTheme = localStorage.getItem('lifegame-theme')

    if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
      theme.value = savedTheme
    } else {
      // 检测系统主题偏好
      mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      const prefersDark = mediaQuery.matches
      theme.value = prefersDark ? 'dark' : 'light'
    }

    applyTheme(theme.value)

    // 监听系统主题变化
    if (!mediaQuery) {
      mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    }
    if (mediaChangeHandler) {
      mediaQuery.removeEventListener('change', mediaChangeHandler)
    }
    mediaChangeHandler = (e) => {
      // 仅当用户未手动设置主题时，跟随系统
      if (!localStorage.getItem('lifegame-theme')) {
        theme.value = e.matches ? 'dark' : 'light'
        applyTheme(theme.value)
      }
    }
    mediaQuery.addEventListener('change', mediaChangeHandler)
  }

  const disposeThemeListener = () => {
    if (mediaQuery && mediaChangeHandler) {
      mediaQuery.removeEventListener('change', mediaChangeHandler)
    }
    mediaChangeHandler = null
  }

  // ==================== 设置指定主题 ====================
  /**
   * 设置指定主题
   * @param {string} newTheme - 'light' 或 'dark'
   */
  const setTheme = (newTheme) => {
    if (newTheme === 'light' || newTheme === 'dark') {
      theme.value = newTheme
      applyTheme(theme.value)
      saveTheme(theme.value)
    }
  }

  return {
    theme,
    initTheme,
    toggleTheme,
    setTheme,
    applyTheme,
    saveTheme,
    disposeThemeListener,
  }
})
