import assert from 'node:assert/strict'
import { createHash } from 'node:crypto'
import { readFileSync } from 'node:fs'
import test from 'node:test'

import {
  DATING_CAREER_SCENE_COUNT,
  getDatingCareerScene,
} from '../src/utils/datingCareerScenes.js'

test('all 50 dating career environments are unique optimized wide WebP assets', () => {
  assert.equal(DATING_CAREER_SCENE_COUNT, 50)
  const hashes = new Set()

  for (let index = 1; index <= DATING_CAREER_SCENE_COUNT; index += 1) {
    const id = String(index).padStart(2, '0')
    const image = readFileSync(new URL(`../public/images/datinginfo/dating-careers/${id}.webp`, import.meta.url))

    assert.equal(image.subarray(0, 4).toString(), 'RIFF')
    assert.equal(image.subarray(8, 12).toString(), 'WEBP')
    assert.equal(image.subarray(12, 16).toString(), 'VP8 ')
    const width = image.readUInt16LE(26) & 0x3fff
    const height = image.readUInt16LE(28) & 0x3fff
    assert.ok(width >= 1600, `${id}.webp must remain a wide high-resolution asset`)
    assert.ok(height >= 900, `${id}.webp must remain a wide high-resolution asset`)
    assert.ok(width / height > 1.7 && width / height < 1.85, `${id}.webp must remain close to 16:9`)
    assert.ok(image.length > 80_000, `${id}.webp is unexpectedly small`)

    hashes.add(createHash('sha256').update(image).digest('hex'))
  }

  assert.equal(hashes.size, DATING_CAREER_SCENE_COUNT)
})

test('female and male targets with the same base id share the matching environment', () => {
  for (let index = 1; index <= DATING_CAREER_SCENE_COUNT; index += 1) {
    const expected = `/images/datinginfo/dating-careers/${String(index).padStart(2, '0')}.webp`
    assert.equal(getDatingCareerScene({ did: index, doccup: '测试职业' }).image, expected)
    assert.equal(getDatingCareerScene({ did: index + 50, doccup: '测试职业' }).image, expected)
  }
})
