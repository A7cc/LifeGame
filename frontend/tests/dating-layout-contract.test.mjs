import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const readText = path => readFileSync(new URL(path, import.meta.url), 'utf8')

test('dating center follows the car and real-estate main-screen layout', () => {
  const dating = readText('../src/components/Menu/MenuDating.vue')
  const carshop = readText('../src/components/Menu/MenuCarshop.vue')
  const realEstate = readText('../src/components/Menu/MenuRealEstate.vue')
  const sharedGrid = /grid-template-columns:\s*repeat\(auto-fill,\s*minmax\(200px,\s*1fr\)\)/

  assert.match(dating, sharedGrid)
  assert.match(carshop, sharedGrid)
  assert.match(realEstate, sharedGrid)
  assert.match(dating, /class="panel-main-left"/)
  assert.match(dating, /class="panel-main-right"/)
  assert.match(dating, /data-testid="dating-relations-panel"/)
  assert.match(dating, /<span class="section-title">我的关系<\/span>/)
})

test('dating cards show each unlocked character preferred locations directly', () => {
  const dating = readText('../src/components/Menu/MenuDating.vue')

  assert.match(dating, /propertyinfo\.dlocations\.join\('、'\)/)
  assert.doesNotMatch(dating, /偏好需要在约会选择中观察/)
})

test('all dating and relationship controls live in the card footer', () => {
  const dating = readText('../src/components/Menu/MenuDating.vue')
  const footerStart = dating.indexOf('<div class="card-footer">')
  const rightPanelStart = dating.indexOf('<!-- 右侧面板 -->')

  assert.ok(footerStart > -1, 'dating card footer was not found')
  assert.ok(rightPanelStart > footerStart, 'dating right panel was not found')
  const cardFooterArea = dating.slice(footerStart, rightPanelStart)
  for (const control of [
    'date-dating-',
    'marry-dating-',
    'breakup-dating-',
    'bath-with-spouse-',
    'divorce-dating-',
  ]) {
    assert.ok(cardFooterArea.includes(control), `${control} must remain in the card footer`)
  }
  assert.match(cardFooterArea, /class="card-actions"/)
  assert.doesNotMatch(cardFooterArea, /property-price/, 'dating cards should not display prices on the main screen')
  assert.match(dating, /grid-auto-columns:\s*minmax\(0,\s*1fr\)/)
  assert.match(dating, /\.card-actions :deep\(\.el-button\)[\s\S]*?height:\s*26px/)
})

test('dating portraits match the compact car-card size and remain previewable', () => {
  const dating = readText('../src/components/Menu/MenuDating.vue')
  const preview = readText('../src/components/Dialog/DialogDatingPortraitPreview.vue')

  assert.match(dating, /\.avatar-image\s*{[\s\S]*?width:\s*40px;[\s\S]*?height:\s*40px;/)
  assert.match(dating, /@click\.stop="openDatingPortraitPreview\(propertyinfo\)"/)
  assert.match(dating, /@click\.stop="openDatingPortraitPreview\(dating\)"/)
  assert.match(dating, /DialogDatingPortraitPreview/)
  assert.doesNotMatch(dating, /preview-src-list|dating-image-preview-open/)
  assert.match(preview, /data-testid="dating-portrait-preview"/)
  assert.match(preview, /:show-close="false"/)
  assert.match(preview, /aria-label="约会对象写真"/)
  assert.match(preview, /:global\(\.dating-portrait-preview-dialog\.el-dialog\)[\s\S]*?background:\s*transparent;[\s\S]*?box-shadow:\s*none;/)
  assert.match(preview, /\.portrait-preview-stage\s*\{[\s\S]*?border:\s*4px solid #d5a43a;[\s\S]*?0 0 0 3px #f1cc72/)
  assert.match(preview, /\.portrait-preview-stage::after\s*\{[\s\S]*?inset:\s*5px;[\s\S]*?border:\s*1px solid rgb\(255 235 169 \/ 88%\);[\s\S]*?inset 0 0 8px/)
  assert.match(preview, /data-testid="dating-portrait-character"/)
  assert.match(preview, /data-testid="dating-portrait-overlay-actions"/)
  assert.match(preview, /data-testid="close-dating-portrait-preview"/)
  assert.doesNotMatch(preview, /class="portrait-preview-actions"/)
  assert.match(preview, /\.portrait-overlay-actions\s*\{[\s\S]*?top:\s*14px;[\s\S]*?right:\s*14px;/)
  assert.match(preview, /data-testid="dating-relationship-status"/)
  assert.match(preview, /class="portrait-relationship-dot"/)
  assert.match(preview, /class="\{ established:/)
  assert.match(preview, /daffinitylevel \|\| '陌生人'/)
  assert.match(preview, /getDatingCareerScene\(props\.dating\)/)
  assert.match(preview, /scene\.occupation/)
  assert.match(preview, /data-testid="dating-scene-hint"/)
  const sceneHintStyle = preview.match(/\.portrait-scene-hint\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(sceneHintStyle, /border|background/)
  assert.match(sceneHintStyle, /portrait-hint-lifecycle/)
  assert.match(preview, /@keyframes portrait-hint-lifecycle/)
  assert.doesNotMatch(preview, /getDatingScene|sceneLocation/)
  assert.match(preview, /class="portrait-backdrop"/)
  assert.match(preview, /class="portrait-character"/)
  const descriptionStyle = preview.match(/\.portrait-description\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(descriptionStyle, /border|background|backdrop-filter/)
})

test('dating videos embed the transparent character directly into the scene', () => {
  const moment = readText('../src/components/Dialog/DialogSpouseMoment.vue')

  assert.match(moment, /data-testid="dating-scene-character"/)
  assert.match(moment, /\.spouse-frame\s*{[\s\S]*?width:\s*54%;/)
  assert.match(moment, /\.spouse-frame img,[\s\S]*?height:\s*100%;[\s\S]*?object-fit:\s*contain;/)
  assert.match(moment, /object-position:\s*center bottom;/)
  assert.doesNotMatch(moment, /height:\s*140px;/)
})

test('all 100 active dating portraits are square WebP assets with transparency', () => {
  for (const gender of ['female', 'male']) {
    for (let index = 1; index <= 50; index += 1) {
      const id = String(index).padStart(2, '0')
      const image = readFileSync(new URL(`../public/images/datinginfo/dating-partner/${gender}/${id}.webp`, import.meta.url))
      assert.equal(image.subarray(0, 4).toString(), 'RIFF')
      assert.equal(image.subarray(8, 12).toString(), 'WEBP')
      assert.equal(image.subarray(12, 16).toString(), 'VP8X')
      assert.ok((image[20] & 0x10) !== 0, `${gender}/${id}.webp must declare alpha`)
      assert.ok(image.includes(Buffer.from('ALPH')), `${gender}/${id}.webp must contain an alpha chunk`)
      const width = image[24] + (image[25] << 8) + (image[26] << 16) + 1
      const height = image[27] + (image[28] << 8) + (image[29] << 16) + 1
      assert.deepEqual([width, height], [1024, 1024])
    }
  }
})
