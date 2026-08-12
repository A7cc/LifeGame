/**
 * ============================================================
 * LifeGame 人生模拟器 - 公共格式化方法
 * ============================================================
 */

/**
 * 格式化价格（根据数值大小自动选择单位）
 * @param {number} price - 价格数值
 * @returns {string} 格式化后的价格字符串
 */
export function formatPrice(price) {
  if (price >= 10000000000000000) {
    return (price / 10000000000000000).toFixed(2).replace(/\.0$/, '') + ' 京'
  } else if (price >= 1000000000000000) {
    return (price / 1000000000000000).toFixed(2).replace(/\.0$/, '') + ' 千万亿'
  } else if (price >= 100000000000000) {
    return (price / 100000000000000).toFixed(2).replace(/\.0$/, '') + ' 百万亿'
  } else if (price >= 10000000000000) {
    return (price / 10000000000000).toFixed(2).replace(/\.0$/, '') + ' 十万亿'
  } else if (price >= 1000000000000) {
    return (price / 1000000000000).toFixed(2).replace(/\.0$/, '') + ' 万亿'
  } else if (price >= 100000000000) {
    return (price / 100000000000).toFixed(2).replace(/\.0$/, '') + ' 千亿'
  } else if (price >= 10000000000) {
    return (price / 10000000000).toFixed(2).replace(/\.0$/, '') + ' 百亿'
  } else if (price >= 1000000000) {
    return (price / 1000000000).toFixed(2).replace(/\.0$/, '') + ' 十亿'
  } else if (price >= 100000000) {
    return (price / 100000000).toFixed(1).replace(/\.0$/, '') + ' 亿'
  } else if (price >= 10000000) {
    return (price / 10000000).toFixed(1).replace(/\.0$/, '') + ' 千万'
  } else if (price >= 1000000) {
    return (price / 1000000).toFixed(1).replace(/\.0$/, '') + ' 百万'
  } else if (price >= 10000) {
    return (price / 10000).toFixed(1).replace(/\.0$/, '') + ' 万'
  } else if (price >= 1000) {
    return (price / 1000).toFixed(1).replace(/\.0$/, '') + ' 千'
  } else {
    return price + ' 元'
  }
}
