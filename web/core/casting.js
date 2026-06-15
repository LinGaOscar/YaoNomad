'use strict'

import { lookup } from './hexagrams.js'

/**
 * 執行一次三硬幣卜卦
 * @param {import('./models.js').Hexagram[]} hexagrams - 預先載入的 64 卦資料
 * @param {string} question - 使用者問題
 * @param {() => number} [coinFn] - 回傳 2 或 3 的函式；預設以 Math.random 模擬
 * @returns {import('./models.js').CastResult}
 */
export function castHexagram(hexagrams, question, coinFn = () => (Math.random() < 0.5 ? 2 : 3)) {
  const lines = []
  for (let i = 0; i < 6; i++) {
    const sum = coinFn() + coinFn() + coinFn()
    lines.push({
      value: sum,
      yin: sum % 2 === 0,
      changing: sum === 6 || sum === 9,
    })
  }

  const lower = trigramCode(lines.slice(0, 3))
  const upper = trigramCode(lines.slice(3, 6))
  const present = lookup(hexagrams, upper, lower)

  const future = lines.some(l => l.changing)
    ? computeFuture(hexagrams, lines)
    : null

  return { question, lines, present, future }
}

/**
 * 計算三條爻的 trigram 三位碼（bit0 = 底爻，yang=1，yin=0）
 * @param {import('./models.js').Line[]} threeLines
 * @returns {number}
 */
function trigramCode(threeLines) {
  return threeLines.reduce((code, line, i) => {
    const yang = line.yin ? 0 : 1
    return code | (yang << i)
  }, 0)
}

/**
 * 計算變卦（將變爻的陰陽反轉後查卦）
 * @param {import('./models.js').Hexagram[]} hexagrams
 * @param {import('./models.js').Line[]} lines
 * @returns {import('./models.js').Hexagram}
 */
function computeFuture(hexagrams, lines) {
  const futureLines = lines.map(l => ({
    value: l.value,
    yin: l.changing ? !l.yin : l.yin,
    changing: l.changing,
  }))
  const lower = trigramCode(futureLines.slice(0, 3))
  const upper = trigramCode(futureLines.slice(3, 6))
  return lookup(hexagrams, upper, lower)
}
