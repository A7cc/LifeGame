import assert from 'node:assert/strict'
import test from 'node:test'
import { createPinia, setActivePinia } from 'pinia'
import { useGameStore } from '../src/stores/game.js'
import { formatPrice } from '../src/utils/format.js'

function newStore() {
  setActivePinia(createPinia())
  return useGameStore()
}

test('stock updates reject stale versions and foreign epochs', () => {
  const store = newStore()
  store.applyGameData({
    gameinfo: { ...store.gameInfo, gtime: 100, gstockinfo: [{ siid: 1, siprice: 100 }], gstocknews: ['initial'] },
    stockepoch: 'game-a',
    stockversion: 1,
  })

  assert.equal(store.applyStockUpdate({
    epoch: 'game-a', version: 2, updatedAt: 123,
    stocks: [{ siid: 1, siprice: 120 }], news: ['new'],
  }), true)
  assert.equal(store.gameInfo.gstockinfo[0].siprice, 120)
  assert.equal(store.stockUpdatedAt, 123)

  assert.equal(store.applyStockUpdate({
    epoch: 'game-a', version: 1, stocks: [{ siid: 1, siprice: 90 }],
  }), false)
  assert.equal(store.applyStockUpdate({
    epoch: 'game-b', version: 99, stocks: [{ siid: 1, siprice: 999 }],
  }), false)
  assert.equal(store.gameInfo.gstockinfo[0].siprice, 120)
})

test('a delayed full response preserves newer stock data', () => {
  const store = newStore()
  store.applyGameData({
    gameinfo: { ...store.gameInfo, gstockinfo: [{ siid: 1, siprice: 100 }], gstocknews: [] },
    stockepoch: 'game-a',
    stockversion: 1,
  })
  store.applyStockUpdate({
    epoch: 'game-a', version: 3, stocks: [{ siid: 1, siprice: 130 }], news: ['latest'],
  })

  store.applyGameData({
    gameinfo: { ...store.gameInfo, gtime: 101, gstockinfo: [{ siid: 1, siprice: 90 }], gstocknews: ['stale'] },
    stockepoch: 'game-a',
    stockversion: 2,
  })

  assert.equal(store.gameInfo.gtime, 101)
  assert.equal(store.gameInfo.gstockinfo[0].siprice, 130)
  assert.deepEqual(store.gameInfo.gstocknews, ['latest'])
  assert.equal(store.stockVersion, 3)
})

test('full game replacement, derived inventory and reset stay consistent', () => {
  const store = newStore()
  store.applyGameData({
    gameinfo: {
      ...store.gameInfo,
      giteminsinfo: {
        1: { iiname: '商品A', iiprice: 20, iidisplay: true },
        2: { iiname: '商品B', iiprice: 30, iidisplay: false },
      },
      gitemoutinfo: {},
      gstockinfo: [{ siid: 1, siprice: 200 }],
    },
    userinfo: {
      ...store.userInfo,
      uitemins: {
        1: { uicostprice: 10, uitemnum: 2 },
        2: { uicostprice: 15, uitemnum: 1 },
      },
    },
    stockepoch: 'game-b',
    stockversion: 0,
  })

  assert.equal(store.userItemInsCount, 3)
  assert.deepEqual(store.userItemInsData.map(item => item.iiprice), [20, 15])
  store.resetGameState()
  assert.equal(store.stockEpoch, '')
  assert.equal(store.stockVersion, 0)
  assert.equal(store.userItemInsCount, 0)
})

test('price formatting covers base and compact units', () => {
  assert.equal(formatPrice(999), '999 元')
  assert.equal(formatPrice(1_000), '1 千')
  assert.equal(formatPrice(10_000), '1 万')
  assert.equal(formatPrice(100_000_000), '1 亿')
})
