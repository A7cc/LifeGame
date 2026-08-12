import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'
import {
  DATING_LOCATION_SCENES,
  DATING_SCENE_GROUPS,
  createDatingSceneOptions,
  getDatingScene,
} from '../src/utils/datingScenes.js'

const read = path => readFileSync(new URL(path, import.meta.url))
const readText = path => read(path).toString('utf8')

const getSeedLocations = () => {
  const source = readText('../../internal/db/seeds.go')
  const datingSeed = source.slice(
    source.indexOf('func seedDatingInfo'),
    source.indexOf('for _, d := range datings'),
  )
  const locations = new Set()
  const datingRow = /`\[[^`]*\]`\s*,\s*`\[[^`]*\]`\s*,\s*`(\[[^`]*\])`\s*\}/g
  for (const match of datingSeed.matchAll(datingRow)) {
    for (const location of JSON.parse(match[1])) locations.add(location)
  }
  return locations
}

test('all 108 dating locations have an explicit visual scene mapping', () => {
  const locations = getSeedLocations()

  assert.equal(locations.size, 108)
  assert.equal(Object.keys(DATING_LOCATION_SCENES).length, locations.size)
  for (const location of locations) {
    assert.ok(DATING_LOCATION_SCENES[location], `${location} is missing a scene group`)
    const scene = getDatingScene(location)
    assert.equal(scene.location, location)
    assert.ok(scene.image)
    assert.ok(scene.icon)
  }
})

test('all generated dating scene backgrounds are optimized production WebP assets', () => {
  const generatedImages = new Set(Object.values(DATING_SCENE_GROUPS)
    .map(scene => scene.image)
    .filter(image => image.startsWith('/images/datinginfo/dating-scenes/')))

  assert.equal(generatedImages.size, 5)
  for (const image of generatedImages) {
    const bytes = read(`../public${image}`)
    assert.equal(bytes.subarray(0, 4).toString(), 'RIFF', `${image} is not a WebP`)
    assert.equal(bytes.subarray(8, 12).toString(), 'WEBP', `${image} is not a WebP`)
    assert.ok(bytes.length > 100_000, `${image} is unexpectedly small`)
  }
})

test('random date choices contain three distinct locations and exactly one preference', () => {
  const preferences = ['公园', '电影院', '咖啡厅']
  const first = createDatingSceneOptions(preferences, undefined, () => 0)
  const second = createDatingSceneOptions(preferences, undefined, () => 0.99)

  for (const choices of [first, second]) {
    assert.equal(choices.length, 3)
    assert.equal(new Set(choices).size, 3)
    assert.equal(choices.filter(location => preferences.includes(location)).length, 1)
    assert.ok(choices.every(location => DATING_LOCATION_SCENES[location]))
  }
  assert.notDeepEqual(first, second)
})

test('location chooser is text-only and only successful dates open a mapped short', () => {
  const chooser = readText('../src/components/Dialog/DialogDatingScene.vue')
  const moment = readText('../src/components/Dialog/DialogSpouseMoment.vue')
  const menu = readText('../src/components/Menu/MenuDating.vue')

  assert.match(chooser, /createDatingSceneOptions\(props\.dating\?\.dlocations\)/)
  assert.match(chooser, /data-testid="dating-scene-options"/)
  assert.doesNotMatch(chooser, /<img|<el-image/)
  assert.match(chooser, /非偏好场景成功率低，但成功回报更高/)
  assert.match(moment, /getDatingScene\(props\.location\)/)
  assert.match(menu, /DoDating\(datingId, location\)/)
  assert.match(menu, /result\.scene\?\.event/)
  assert.match(menu, /rewardtier === 'high-risk'/)
  assert.match(menu, /冒险约会成功，高回报/)
  assert.match(menu, /稳妥约会成功/)
  assert.match(menu, /if \(result\.success && dating\)/)
})
