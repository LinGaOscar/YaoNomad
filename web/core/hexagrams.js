'use strict'

/**
 * 在卦象列表中以上下卦三位碼查找對應卦象
 * @param {import('./models.js').Hexagram[]} hexagrams
 * @param {number} upper - 上卦三位碼（0–7）
 * @param {number} lower - 下卦三位碼（0–7）
 * @returns {import('./models.js').Hexagram}
 * @throws {Error} 若找不到對應卦象
 */
export function lookup(hexagrams, upper, lower) {
  const found = hexagrams.find(h => h.upper === upper && h.lower === lower)
  if (!found) {
    throw new Error(`找不到卦象: upper=${upper} lower=${lower}`)
  }
  return found
}
