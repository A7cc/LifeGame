import assert from 'node:assert/strict'
import { createHash } from 'node:crypto'
import { existsSync, readFileSync } from 'node:fs'
import test from 'node:test'

const readText = path => readFileSync(new URL(path, import.meta.url), 'utf8')

test('car and house pages share an optimized purchase-delivery moment', () => {
  const carMenu = readText('../src/components/Menu/MenuCarshop.vue')
  const houseMenu = readText('../src/components/Menu/MenuRealEstate.vue')

  for (const menu of [carMenu, houseMenu]) {
    assert.match(menu, /DialogAssetPurchaseMoment/)
    assert.match(menu, /mode="preview"/)
    assert.match(menu, /assetPreviewVisible/)
    assert.match(menu, /openAssetPreview/)
    assert.match(menu, /panel-subtitle/)
    assert.match(menu, /width: 40px/)
    assert.match(menu, /height: 40px/)
    assert.match(menu, /purchaseMomentVisible\.value = true/)
    assert.match(menu, /:owned="Boolean\(userInfo\.(?:ucar|uhouse)\?\.\[previewed(?:Car|House)\?\.(?:ciid|hiid)\]\)"/)
    assert.doesNotMatch(menu, /交付 · 永久资产/)
    assert.doesNotMatch(menu, /preview-src-list/)
  }

  assert.match(carMenu, /kind="car"/)
  assert.match(carMenu, /data-testid="`owned-car-image-\$\{carId\}`"/)
  assert.match(carMenu, /getCarData\(carId, 'img'\)/)
  assert.match(carMenu, /purchasedCar\.value = \{ \.\.\.car \}/)
  assert.match(houseMenu, /kind="house"/)
  assert.match(houseMenu, /data-testid="`owned-house-image-\$\{houseId\}`"/)
  assert.match(houseMenu, /getHouseData\(houseId, 'img'\)/)
  assert.match(houseMenu, /purchasedHouse\.value = \{ \.\.\.house \}/)
})

test('purchase delivery short contains animated steps and production backdrops', () => {
  const dialog = readText('../src/components/Dialog/DialogAssetPurchaseMoment.vue')

  assert.match(dialog, /data-testid="asset-purchase-moment"/)
  assert.match(dialog, /:show-close="false"/)
  assert.match(dialog, /aria-label="房车场景"/)
  assert.match(dialog, /:global\(\.asset-purchase-moment-dialog\.el-dialog\)[\s\S]*?background:\s*transparent;[\s\S]*?box-shadow:\s*none;/)
  assert.match(dialog, /\.purchase-stage\s*\{[\s\S]*?border:\s*4px solid #d5a43a;[\s\S]*?0 0 0 3px #f1cc72/)
  assert.match(dialog, /\.purchase-stage::after\s*\{[\s\S]*?inset:\s*5px;[\s\S]*?border:\s*1px solid rgb\(255 235 169 \/ 88%\);[\s\S]*?inset 0 0 8px/)
  assert.match(dialog, /data-testid="asset-purchase-stage"/)
  assert.match(dialog, /v-for="\(beat, index\) in beats"/)
  assert.doesNotMatch(dialog, /purchase-progress|showcase-accent/)
  assert.doesNotMatch(dialog, /新车交付短片|新居交付短片/)
  assert.doesNotMatch(dialog, /约 7 秒|仅在购买成功后播放|场景化放大预览|资产与环境一起展示/)
  assert.match(dialog, /data-testid="asset-purchase-item"/)
  assert.match(dialog, /data-testid="asset-overlay-actions"/)
  assert.match(dialog, /data-testid="replay-purchase-moment"/)
  assert.match(dialog, /data-testid="close-purchase-moment"/)
  assert.ok(dialog.indexOf('data-testid="close-purchase-moment"') < dialog.indexOf('data-testid="replay-purchase-moment"'))
  assert.doesNotMatch(dialog, /class="purchase-actions"/)
  assert.match(dialog, /\.purchase-overlay-actions\s*\{[\s\S]*?top:\s*14px;[\s\S]*?right:\s*14px;/)
  assert.match(dialog, /\.purchase-overlay-actions\s*\{[\s\S]*?flex-direction:\s*column;/)
  assert.match(dialog, /data-testid="asset-showcase-specs"/)
  assert.match(dialog, /asset-showcase-stage/)
  assert.match(dialog, /mode: \{/)
  assert.doesNotMatch(dialog, /座驾全景鉴赏|理想住所全景鉴赏/)
  assert.match(dialog, /props\.item\?\.ciimg/)
  assert.match(dialog, /props\.item\?\.hiimg/)
  assert.doesNotMatch(dialog, /delivery-stamp|车辆交付完成|房产交付完成|已加入我的车库|已加入我的房产/)
  assert.match(dialog, /'asset-ownership-status' : 'asset-purchase-status'/)
  assert.match(dialog, /isPreview \? \(owned \? '已购买' : '未购买'\) : '已购买'/)
  assert.doesNotMatch(dialog, /class="purchase-ending"/)
  assert.doesNotMatch(dialog, /座驾展示中|住所展示中|灯光与环境已就位|采光与空间已就位/)
  assert.match(dialog, /data-testid="asset-showcase-bonuses"/)
  assert.match(dialog, /owned \? '已购买' : '未购买'/)
  assert.doesNotMatch(dialog, /class="purchase-result"/)
  assert.match(dialog, /isPreview \? '参考价' : '成交价'/)
  const showcaseBonusStyle = dialog.match(/\.showcase-bonuses\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(showcaseBonusStyle, /border|background/)
  const showcaseHintStyle = dialog.match(/\.showcase-specs\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(showcaseHintStyle, /border|background/)
  assert.match(showcaseHintStyle, /showcase-hint-lifecycle/)
  assert.match(dialog, /@keyframes showcase-hint-lifecycle/)
  const purchaseBeatStyle = dialog.match(/\.purchase-beat\s*\{([\s\S]*?)\}/)?.[1] || ''
  assert.doesNotMatch(purchaseBeatStyle, /border|background/)
  assert.match(purchaseBeatStyle, /purchase-beat-show/)

  for (const image of [
    '../public/images/carinfo/car-moments/car-showroom.webp',
    '../public/images/houseinfo/house-moments/home-handover.webp',
  ]) {
    const path = new URL(image, import.meta.url)
    assert.equal(existsSync(path), true)
    const webp = readFileSync(path)
    assert.equal(webp.subarray(0, 4).toString(), 'RIFF')
    assert.equal(webp.subarray(8, 12).toString(), 'WEBP')
    assert.ok(webp.length > 100_000)
  }
})

test('all 50 cars and 50 houses use independent transparent WebP product assets', () => {
  const seeds = readText('../../internal/db/seeds.go')
  const houseBlock = seeds.slice(seeds.indexOf('func seedHouses'), seeds.indexOf('func seedCars'))
  const carBlock = seeds.slice(seeds.indexOf('func seedCars'), seeds.indexOf('func seedDatingInfo'))
  const houseImages = [...houseBlock.matchAll(/"(\/images\/houseinfo\/houses\/[^"\s]+\.webp)"/g)].map(match => match[1])
  const carImages = [...carBlock.matchAll(/"(\/images\/carinfo\/cars\/[^"\s]+\.webp)"/g)].map(match => match[1])

  assert.equal(houseImages.length, 50)
  assert.equal(carImages.length, 50)
  assert.equal(new Set(houseImages).size, 50)
  assert.equal(new Set(carImages).size, 50)
  assert.deepEqual(
    houseImages.map(image => image.split('/').at(-1).slice(0, 2)),
    Array.from({ length: 50 }, (_, index) => String(index + 1).padStart(2, '0')),
  )
  assert.deepEqual(
    carImages.map(image => image.split('/').at(-1).slice(0, 2)),
    Array.from({ length: 50 }, (_, index) => String(index + 1).padStart(2, '0')),
  )

  const imageHashes = []
  for (const image of new Set([...houseImages, ...carImages])) {
    const path = new URL(`../public${image}`, import.meta.url)
    assert.equal(existsSync(path), true, `${image} should exist`)
    const webp = readFileSync(path)
    assert.equal(webp.subarray(0, 4).toString(), 'RIFF')
    assert.equal(webp.subarray(8, 12).toString(), 'WEBP')
    assert.ok(webp.length > 10_000)
    assert.notEqual(webp.indexOf(Buffer.from('ALPH')), -1, `${image} should contain transparency`)
    imageHashes.push(createHash('sha256').update(webp).digest('hex'))
  }
  assert.equal(new Set(imageHashes).size, 100, 'every car and house should have distinct image content')
})
