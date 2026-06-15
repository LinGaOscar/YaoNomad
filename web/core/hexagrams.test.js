import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { lookup } from './hexagrams.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
const hexagrams = JSON.parse(readFileSync(join(__dirname, '../data/hexagrams.json'), 'utf-8'))

test('hexagrams.json 應有 64 筆', () => {
  assert.strictEqual(hexagrams.length, 64)
})

test('hexagrams.json 的 upper/lower 組合應唯一', () => {
  const seen = new Set()
  for (const h of hexagrams) {
    const key = `${h.upper},${h.lower}`
    assert.ok(!seen.has(key), `重複的組合: ${key}`)
    seen.add(key)
  }
})

test('lookup 找到乾卦 (upper=7, lower=7)', () => {
  const h = lookup(hexagrams, 7, 7)
  assert.strictEqual(h.number, 1)
  assert.strictEqual(h.name, '乾')
})

test('lookup 找到坤卦 (upper=0, lower=0)', () => {
  const h = lookup(hexagrams, 0, 0)
  assert.strictEqual(h.number, 2)
  assert.strictEqual(h.name, '坤')
})

test('lookup 查無效碼應拋出錯誤', () => {
  assert.throws(() => lookup(hexagrams, 8, 8), /找不到卦象/)
})
