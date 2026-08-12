import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const variablesCss = readFileSync(new URL('../src/style-variables.css', import.meta.url), 'utf8')
const appVue = readFileSync(new URL('../src/App.vue', import.meta.url), 'utf8')
const menuVue = readFileSync(new URL('../src/components/GameMenu.vue', import.meta.url), 'utf8')

function ruleVariables(selector) {
  const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const match = variablesCss.match(new RegExp(`${escaped}\\s*\\{([\\s\\S]*?)\\n\\}`))
  assert.ok(match, `missing ${selector} theme rule`)
  return Object.fromEntries(
    [...match[1].matchAll(/(--[\w-]+)\s*:\s*([^;]+);/g)].map(([, key, value]) => [key, value.trim()]),
  )
}

function luminance(hex) {
  const normalized = hex.replace('#', '')
  const value = normalized.length === 3
    ? normalized.split('').map(character => character.repeat(2)).join('')
    : normalized
  const channels = value.match(/.{2}/g).map(channel => Number.parseInt(channel, 16) / 255)
  return channels
    .map(channel => channel <= 0.03928 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4)
    .reduce((sum, channel, index) => sum + channel * [0.2126, 0.7152, 0.0722][index], 0)
}

function contrast(first, second) {
  const values = [luminance(first), luminance(second)].sort((a, b) => b - a)
  return (values[0] + 0.05) / (values[1] + 0.05)
}

test('light and dark themes define complete semantic color tokens', () => {
  const light = ruleVariables(':root')
  const dark = ruleVariables('[data-theme="dark"]')
  const required = [
    '--primary-color', '--accent-strong-color', '--background-color', '--panel-main-color',
    '--panel-color', '--panel-sub-color', '--panel-log-color', '--border-color',
    '--font-color', '--font-secondary', '--font-light', '--success-color-rgb',
    '--error-color-rgb', '--gradient-primary', '--select-color', '--el-color-primary',
  ]

  for (const token of required) {
    assert.ok(light[token], `light theme is missing ${token}`)
    assert.ok(dark[token], `dark theme is missing ${token}`)
  }
})

test('small semantic text meets AA contrast in both themes', () => {
  const light = ruleVariables(':root')
  const dark = ruleVariables('[data-theme="dark"]')

  for (const [name, theme] of [['light', light], ['dark', dark]]) {
    assert.ok(contrast(theme['--font-secondary'], theme['--panel-color']) >= 4.5, `${name} secondary text contrast is too low`)
    assert.ok(contrast(theme['--font-light'], theme['--panel-color']) >= 4.5, `${name} light text contrast is too low`)
    assert.ok(contrast(theme['--primary-color'], theme['--panel-color']) >= 4.5, `${name} primary accent contrast is too low`)
  }

  assert.ok(contrast(light['--el-color-primary'], light['--on-primary-color']) >= 4.5, 'primary button text contrast is too low')
})

test('brand surfaces use the blue and warm-gold palette without purple controls', () => {
  assert.doesNotMatch(appVue, /#845ef7|#5f3dc4|#9775fa/i)
  assert.match(appVue, /#2563eb/)
  assert.match(appVue, /#f59e0b/)
  assert.doesNotMatch(menuVue, /menuBgColor|menuTextColor/)
  assert.match(menuVue, /--el-menu-bg-color:\s*var\(--panel-main-color\)/)
})
