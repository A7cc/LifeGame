import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'
import { getDatingGiftScene } from '../src/utils/datingGifts.js'

const readText = path => readFileSync(new URL(path, import.meta.url), 'utf8')

test('gift chooser mirrors date chooser without revealing price or outcome in each choice', () => {
  const dialog = readText('../src/components/Dialog/DialogDatingGift.vue')
  const menu = readText('../src/components/Menu/MenuDating.vue')

  assert.match(dialog, /width="520px"/)
  assert.match(dialog, /data-testid="dating-gift-options"/)
  assert.match(dialog, /v-for="\(gift, index\) in options"/)
  assert.match(dialog, /selectedOption.*giftCost/s)
  assert.doesNotMatch(dialog, /class="gift-price"/)
  assert.doesNotMatch(dialog, /gift\.effecthint/)
  assert.match(dialog, /非偏好礼物很难送出，但成功回报更高/)
  assert.match(dialog, /selectedGift === index/)
  assert.match(dialog, /data-testid="confirm-dating-gift"/)
  assert.match(menu, /GetDatingGiftOptions\(dating\.did\)/)
  assert.match(menu, /GiveDatingGift\(dating\.did, gift\)/)
})

test('gift settlement opens steady, risky-success, or rejected scene shorts', () => {
  const menu = readText('../src/components/Menu/MenuDating.vue')
  const moment = readText('../src/components/Dialog/DialogSpouseMoment.vue')

  assert.match(menu, /openSpouseMoment\(\s*'gift'/)
  assert.match(menu, /result\.outcome/)
  assert.match(menu, /result\.success/)
  assert.match(moment, /getDatingGiftScene\(props\.gift\)/)
  assert.match(moment, /'gift-favorite'/)
  assert.match(moment, /'gift-risky-success'/)
  assert.match(moment, /'gift-rejected'/)
  assert.doesNotMatch(moment, /送礼惊喜短片|冒险送礼成功短片|送礼被婉拒短片/)
  assert.match(moment, /正中喜好/)
  assert.match(moment, /意外打动/)
  assert.match(moment, /这次没有送出/)
  assert.match(moment, /props\.event/)

  for (const gift of ['鲜花', '巧克力', '游戏机', '乐谱', '书籍']) {
    const scene = getDatingGiftScene(gift)
    assert.ok(scene.image.startsWith('/images/'))
    assert.ok(scene.icon)
    assert.ok(scene.label)
  }
})
