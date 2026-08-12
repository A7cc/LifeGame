import { expect } from '@playwright/test'

export async function dismissAnnouncement(page) {
  const title = page.getByText('📋 公告', { exact: true })
  try {
    await title.waitFor({ state: 'visible', timeout: 2_000 })
    await page.keyboard.press('Escape')
    await title.waitFor({ state: 'hidden' })
  } catch (error) {
    if (await title.isVisible()) throw error
  }
}

export async function startNewGame(page, {
  name = 'GUI自动测试',
  gender = 'male',
  difficulty = 0,
} = {}) {
  await page.goto('/')
  await expect(page.getByTestId('start-screen')).toBeVisible()
  await page.getByTestId('player-name').fill(name)
  await page.getByTestId(`gender-${gender}`).click()
  await page.getByTestId(`difficulty-${difficulty}`).click()
  await page.getByTestId('start-game').click()
  await expect(page.getByTestId('game-main')).toBeVisible()
  await expect(page.getByTestId('player-summary')).toContainText(name)
  await dismissAnnouncement(page)
}
