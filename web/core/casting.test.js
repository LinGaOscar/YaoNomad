import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { castHexagram } from './casting.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
const hexagrams = JSON.parse(readFileSync(join(__dirname, '../data/hexagrams.json'), 'utf-8'))

test('全背面(3) → 本卦乾(#1)，變卦坤(#2)', () => {
  const result = castHexagram(hexagrams, '測試', () => 3)
  assert.strictEqual(result.present.number, 1)
  assert.strictEqual(result.present.name, '乾')
  assert.ok(result.future !== null, '應有變卦')
  assert.strictEqual(result.future.number, 2)
  assert.strictEqual(result.future.name, '坤')
})

test('全正面(2) → 本卦坤(#2)，變卦乾(#1)', () => {
  const result = castHexagram(hexagrams, '測試', () => 2)
  assert.strictEqual(result.present.number, 2)
  assert.strictEqual(result.present.name, '坤')
  assert.ok(result.future !== null, '應有變卦')
  assert.strictEqual(result.future.number, 1)
  assert.strictEqual(result.future.name, '乾')
})

test('少陽序列(2,2,3=7) → 無變爻 → future 為 null', () => {
  const coins = [2, 2, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++ % 3])
  for (const l of result.lines) {
    assert.strictEqual(l.value, 7)
    assert.strictEqual(l.yin, false)
    assert.strictEqual(l.changing, false)
  }
  assert.strictEqual(result.future, null)
})

test('少陰序列(2,3,3=8) → 無變爻 → future 為 null', () => {
  const coins = [2, 3, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++ % 3])
  for (const l of result.lines) {
    assert.strictEqual(l.value, 8)
    assert.strictEqual(l.yin, true)
    assert.strictEqual(l.changing, false)
  }
  assert.strictEqual(result.future, null)
})

test('question 欄位正確帶入', () => {
  const result = castHexagram(hexagrams, '今日運勢如何？', () => 3)
  assert.strictEqual(result.question, '今日運勢如何？')
})

test('爻值均在合法範圍 6–9', () => {
  const coins = [3, 3, 3, 2, 2, 2, 3, 2, 3, 2, 3, 2, 3, 3, 2, 2, 2, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++])
  for (const l of result.lines) {
    assert.ok([6, 7, 8, 9].includes(l.value), `爻值 ${l.value} 不合法`)
  }
})
