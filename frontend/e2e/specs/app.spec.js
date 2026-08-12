import { test, expect } from '@playwright/test'
import { dismissAnnouncement, startNewGame } from './helpers.js'

test('启动页校验并创建新游戏', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByTestId('start-screen')).toBeVisible()

  await page.getByTestId('start-game').click()
  await expect(page.getByText('请输入你的名字', { exact: true })).toBeVisible()

  await page.getByTestId('player-name').fill('小林')
  await page.getByTestId('gender-female').click()
  await page.getByTestId('difficulty-2').click()
  await expect(page.getByTestId('gender-female')).toHaveClass(/selected/)
  await expect(page.getByTestId('difficulty-2')).toHaveClass(/selected/)

  await page.getByTestId('start-game').click()
  await expect(page.getByTestId('game-main')).toBeVisible()
  await expect(page.getByTestId('player-summary')).toContainText('小林（女）')
  await dismissAnnouncement(page)

  await page.getByTestId('menu-dating').click()
  await expect(page.getByTestId('dating-list-title')).toHaveText('男友列表')
  await expect(page.locator('[data-testid^="dating-card-"]')).toHaveCount(50)
  await expect(page.locator('[data-testid^="dating-card-"][data-dating-sex="female"]')).toHaveCount(0)
  await expectDatingPortrait(page, 'male')
})

test('全部主菜单页面可以加载', async ({ page }) => {
  await startNewGame(page)

  const pages = [
    ['domestic', '/menu-domestic'],
    ['foreign', '/menu-foreign'],
    ['company', '/menu-company'],
    ['stock', '/menu-stock'],
    ['antique', '/menu-antique'],
    ['entertainment', '/menu-entertainment'],
    ['bank', '/menu-bank'],
    ['hospital', '/menu-hospital'],
    ['dating', '/menu-dating'],
    ['carshop', '/menu-carshop'],
    ['realestate', '/menu-realestate'],
    ['profile', '/menu-profile'],
  ]

  for (const [name, path] of pages) {
    await page.getByTestId(`menu-${name}`).click()
    await expect(page).toHaveURL(new RegExp(`${path}$`))
    await expect(page.getByTestId(`page-${name}`)).toBeVisible()
    if (name === 'antique') {
      await page.getByRole('button', { name: '下一场竞拍', exact: true }).click()
      const antiqueImage = page.locator('.antique-img')
      await expect(antiqueImage).toHaveAttribute('src', /\/images\/antiqueinfo\/\d{2}-[a-z0-9-]+\.webp/)
      await expect.poll(() => antiqueImage.evaluate(element => element.naturalWidth)).toBe(1024)
    }
    if (name === 'dating') {
      await expect(page.getByTestId('dating-list-title')).toHaveText('女友列表')
      await expectDatingPortrait(page, 'female')
      await expectDatingLayout(page, 'female')
      const visitScene = page.getByTestId('visit-dating-scene')
      await expect(visitScene).toBeEnabled()
      await visitScene.click()
      await expect(page.locator('.el-message').last()).toContainText(/认识了|这次没有遇到/)
    }
  }
})

test('浅色与暗色主题使用统一且可读的语义色', async ({ page }) => {
  await startNewGame(page, { name: '主题测试' })

  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light')
  let colors = await page.locator('html').evaluate(element => {
    const style = getComputedStyle(element)
    return {
      background: style.getPropertyValue('--background-color').trim(),
      panel: style.getPropertyValue('--panel-color').trim(),
      secondary: style.getPropertyValue('--font-secondary').trim(),
    }
  })
  expect(colors).toEqual({ background: '#f1f5f9', panel: '#f8fafc', secondary: '#596579' })

  await page.getByTestId('theme-toggle').click()
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark')
  colors = await page.locator('html').evaluate(element => {
    const style = getComputedStyle(element)
    return {
      background: style.getPropertyValue('--background-color').trim(),
      panel: style.getPropertyValue('--panel-color').trim(),
      secondary: style.getPropertyValue('--font-secondary').trim(),
    }
  })
  expect(colors).toEqual({ background: '#0b1220', panel: '#152238', secondary: '#b8c4d6' })
})

async function expectDatingLayout(page, sex) {
  const card = page.locator(`[data-testid^="dating-card-"][data-dating-sex="${sex}"]`).first()
  const header = card.locator('.card-header')
  const footer = card.locator('.card-footer')
  const relations = page.getByTestId('dating-relations-panel')

  await expect(relations).toBeVisible()
  await expect(footer).toBeVisible()
  await expect(card.locator('.property-price')).toHaveCount(0)

  const [cardBox, headerBox, footerBox, relationsBox] = await Promise.all([
    card.boundingBox(),
    header.boundingBox(),
    footer.boundingBox(),
    relations.boundingBox(),
  ])
  expect(cardBox).not.toBeNull()
  expect(headerBox).not.toBeNull()
  expect(footerBox).not.toBeNull()
  expect(relationsBox).not.toBeNull()
  expect(cardBox.width).toBeGreaterThanOrEqual(200)
  expect(cardBox.width).toBeLessThan(300)
  expect(footerBox.y).toBeGreaterThan(headerBox.y + headerBox.height)
  expect(relationsBox.x).toBeGreaterThan(cardBox.x + cardBox.width)
}

async function expectDatingPortrait(page, sex) {
  const cards = page.locator(`[data-testid^="dating-card-"][data-dating-sex="${sex}"]`)
  await expect(cards).toHaveCount(50)

  const image = cards.first().locator('.avatar-image img')
  await expect(image).toBeVisible()
  await expect.poll(() => image.evaluate(element => element.naturalWidth)).toBeGreaterThan(0)
  const imageData = await image.evaluate(element => {
    const canvas = document.createElement('canvas')
    canvas.width = element.naturalWidth
    canvas.height = element.naturalHeight
    const context = canvas.getContext('2d')
    context.drawImage(element, 0, 0)
    const corners = [[0, 0], [canvas.width - 1, 0], [0, canvas.height - 1], [canvas.width - 1, canvas.height - 1]]
    return {
      width: element.naturalWidth,
      height: element.naturalHeight,
      cornerAlpha: corners.map(([x, y]) => context.getImageData(x, y, 1, 1).data[3]),
    }
  })
  expect(imageData).toEqual({ width: 1024, height: 1024, cornerAlpha: [0, 0, 0, 0] })

  const box = await image.boundingBox()
  expect(box).not.toBeNull()
  expect(Math.abs(box.width - box.height)).toBeLessThan(1)

  await image.click()
  const preview = page.getByTestId('dating-portrait-preview')
  await expect(preview).toBeVisible()
  const backdrop = preview.locator('.portrait-backdrop')
  await expect(backdrop).toHaveAttribute('src', /\/images\/datinginfo\/dating-careers\/01\.webp/)
  await expect.poll(() => backdrop.evaluate(element => element.naturalWidth)).toBeGreaterThan(0)
  const character = preview.getByTestId('dating-portrait-character')
  await expect(character).toHaveAttribute('src', /\/images\/datinginfo\/dating-partner\/(female|male)\/\d{2}\.webp/)
  await expect.poll(() => character.evaluate(element => element.naturalWidth)).toBeGreaterThan(0)
  await preview.getByTestId('close-dating-portrait-preview').click()
  await expect(preview).toBeHidden()
}

test('送礼会触发人物事件并写入约会记录', async ({ page }) => {
  await startNewGame(page, { name: '送礼测试', gender: 'female' })

  await page.getByTestId('menu-dating').click()
  await page.locator('.meeting-scene-select').click()
  await page.getByRole('option', { name: '公园', exact: true }).click()
  await page.getByTestId('visit-dating-scene').click()
  await expect(page.locator('.el-message').last()).toContainText('认识了')

  const datingCard = page.getByTestId('dating-card-51')
  await expect(datingCard).not.toHaveClass(/locked/)
  const datingName = (await datingCard.locator('.property-name').textContent()).trim()
  const giftButton = page.getByTestId('gift-dating-51')
  await expect(giftButton).toBeEnabled()
  await giftButton.click()
  await expect(page.getByTestId('dating-gift-dialog')).toBeVisible()
  const giftOptions = page.getByTestId('dating-gift-options').locator('.gift-item')
  await expect(giftOptions).toHaveCount(3)
  await expect(giftOptions.locator('.gift-price')).toHaveCount(0)
  await expect(giftOptions.locator('.gift-effect')).toHaveCount(0)
  const giftOption = page.getByTestId('gift-option-51-0')
  const giftName = (await giftOption.locator('.gift-name').textContent()).trim()
  await giftOption.click()
  await expect(giftOption).toHaveClass(/selected/)
  await page.getByTestId('confirm-dating-gift').click()

  const giftMoment = page.getByTestId('spouse-moment-dialog')
  await expect(giftMoment).toBeVisible()
  await expect(giftMoment).toContainText(/收到了一份惊喜|大胆的选择带来了意外惊喜|这份心意没有在今天抵达/)
  await expect(giftMoment).toContainText(datingName)
  await expect(giftMoment).toContainText(giftName)
  await expect(giftMoment.locator('.moment-stage')).toHaveClass(/moment-gift-(favorite|risky-success|rejected)/)
  await expect(page.locator('.log-panel')).toContainText(`送给${datingName}「${giftName}」`)
})

test('推进时间、保存、加载和结束游戏', async ({ page }, testInfo) => {
  const saveName = `自动存档-${testInfo.repeatEachIndex}`
  await startNewGame(page, { name: '存档测试' })
  await expect(page.getByTestId('player-summary')).toContainText('18岁')

  await page.getByTestId('open-save-dialog').click()
  await expect(page.getByText('存档管理', { exact: true })).toBeVisible()
  await page.getByPlaceholder('输入存档名称').fill(saveName)
  await page.getByTestId('create-save').click()
  await expect(page.getByText('存档成功', { exact: true })).toBeVisible()
  await expect(page.locator('.save-item').filter({ hasText: saveName })).toBeVisible()
  await page.keyboard.press('Escape')

  await page.getByTestId('next-year').click()
  await expect(page.getByTestId('player-summary')).toContainText('19岁')
  await dismissAnnouncement(page)

  await page.getByTestId('open-save-dialog').click()
  const saved = page.locator('.save-item').filter({ hasText: saveName })
  await saved.getByRole('button', { name: '加载' }).click()
  await expect(page.getByText('加载成功', { exact: true })).toBeVisible()
  await expect(page.getByTestId('player-summary')).toContainText('18岁')
  await dismissAnnouncement(page)

  await page.getByTestId('end-game').click()
  const evaluationDialog = page.getByTestId('evaluation-dialog')
  await expect(evaluationDialog).toBeVisible()
  await expect(evaluationDialog.getByText('年龄评分')).toHaveCount(0)
  await expect(evaluationDialog).toHaveCSS('border-top-width', '2px')
  await expect(evaluationDialog).toHaveCSS('border-top-style', 'solid')
  await page.getByTestId('restart-game').click()
  await expect(page.getByTestId('start-screen')).toBeVisible()
})

test('小游戏关闭后会取消会话并允许重新开始', async ({ page }) => {
  await startNewGame(page, { name: '小游戏测试' })
  await page.getByTestId('menu-entertainment').click()
  await expect(page.getByTestId('page-entertainment')).toBeVisible()

  await page.getByText('休闲游戏', { exact: true }).click()
  await page.getByText('猜数字', { exact: true }).click()

  let dialog = page.getByRole('dialog', { name: '猜数字游戏' })
  await dialog.getByRole('button', { name: '开始游戏' }).click()
  await expect(dialog.getByRole('button', { name: '猜！' })).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dialog).toBeHidden()

  await page.waitForTimeout(300)
  await page.getByText('猜数字', { exact: true }).click()
  dialog = page.getByRole('dialog', { name: '猜数字游戏' })
  await dialog.getByRole('button', { name: '开始游戏' }).click()
  await expect(dialog.getByRole('button', { name: '猜！' })).toBeVisible()
  await dialog.getByRole('button', { name: '猜！' }).click()
  await expect(dialog).toContainText('正确答案')
})
