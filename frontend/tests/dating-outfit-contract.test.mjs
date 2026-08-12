import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'
import { availableDatingOutfits } from '../src/utils/datingOutfits.js'

const read = path => readFileSync(new URL(path, import.meta.url))
const readText = path => read(path).toString('utf8')

const assertTransparentWebp = (path, width, height) => {
  const image = read(path)
  assert.equal(image.subarray(0, 4).toString(), 'RIFF', `${path} is not WebP`)
  assert.equal(image.subarray(8, 12).toString(), 'WEBP', `${path} is not WebP`)
  assert.equal(image.subarray(12, 16).toString(), 'VP8X', `${path} must use extended WebP`)
  assert.ok((image[20] & 0x10) !== 0, `${path} must declare alpha`)
  assert.ok(image.includes(Buffer.from('ALPH')), `${path} must contain alpha data`)
  const actualWidth = image[24] + (image[25] << 8) + (image[26] << 16) + 1
  const actualHeight = image[27] + (image[28] << 8) + (image[29] << 16) + 1
  assert.deepEqual([actualWidth, actualHeight], [width, height], `${path} dimensions`)
}

test('dating outfits are explicitly selected and unlocked by relationship stage', () => {
  const menu = readText('../src/components/Menu/MenuDating.vue')
  const dialog = readText('../src/components/Dialog/DialogSpouseMoment.vue')
  const outfits = readText('../src/utils/datingOutfits.js')
  const backend = readText('../../services/dating_interactions.go')

  assert.match(menu, /:character-image="momentCharacterImage"/)
  for (const outfit of ['career', 'sleepwear', 'romantic', 'swimwear', 'cosplay', 'qipao', 'homewear']) {
    assert.match(outfits, new RegExp(`key: '${outfit}'`))
  }
  assert.match(menu, /kind === 'bath'.*'swimwear'/s)
  assert.match(menu, /variant === 'intimacy'.*'romantic'/s)
  assert.doesNotMatch(menu, /randomDatingOutfit/)
  assert.doesNotMatch(backend, /rand\.Intn\(len\(datingOutfitVariants\)\)/)
  assert.match(menu, /DoDatingInteraction\(momentSpouse\.value\.did, action, outfitKey\)/)
  assert.match(dialog, /value: 'outfit'/)
  assert.match(dialog, /data-testid="dating-outfit-picker"/)
  assert.doesNotMatch(dialog, /指定本次造型/)
  assert.doesNotMatch(dialog, /不会随机换装/)
  assert.match(dialog, /:title="outfitOption\.label"/)
  assert.match(dialog, /availableDatingOutfits/)

  assert.match(outfits, /career'.*minRank: 1/)
  assert.match(outfits, /homewear'.*minRank: 1/)
  assert.match(outfits, /qipao'.*minRank: 2/)
  assert.match(outfits, /swimwear'.*minRank: 3/)
  assert.match(outfits, /sleepwear'.*minRank: 4/)
  assert.match(outfits, /romantic'.*minRank: 5.*adultOnly: true/)

  const keysFor = (status, playerAge = 25, partnerAge = 25) => (
    availableDatingOutfits(status, playerAge, partnerAge).map(outfit => outfit.key)
  )
  assert.deepEqual(keysFor('陌生人'), [])
  assert.deepEqual(keysFor('朋友'), ['career', 'homewear'])
  assert.deepEqual(keysFor('暧昧中'), ['career', 'homewear', 'qipao'])
  assert.deepEqual(keysFor('交往中'), ['career', 'homewear', 'qipao', 'cosplay', 'swimwear'])
  assert.deepEqual(keysFor('恋人'), ['career', 'homewear', 'qipao', 'cosplay', 'swimwear', 'sleepwear'])
  assert.equal(keysFor('专属恋人', 17, 25).includes('romantic'), false)
  assert.equal(keysFor('专属恋人').includes('romantic'), true)
})

test('all 100 characters ship seven transparent full-body outfits', () => {
  const variants = ['sleepwear', 'romantic', 'swimwear', 'cosplay', 'qipao', 'homewear']
  for (const gender of ['female', 'male']) {
    for (let index = 1; index <= 50; index += 1) {
      const id = String(index).padStart(2, '0')
      assertTransparentWebp(`../public/images/datinginfo/dating-partner/${gender}/${id}.webp`, 1024, 1024)
      for (const variant of variants) {
        assertTransparentWebp(`../public/images/datinginfo/dating-partner/${gender}/${variant}/${id}.webp`, 1024, 1024)
      }
    }
  }
})
