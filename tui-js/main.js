#!/usr/bin/env node
'use strict'

import { readFileSync } from 'node:fs'
import { createInterface } from 'node:readline'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { castHexagram } from '../web/core/casting.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
const hexagrams = JSON.parse(readFileSync(join(__dirname, '../web/data/hexagrams.json'), 'utf8'))

const rl = createInterface({ input: process.stdin, output: process.stdout })

rl.question('請輸入您的問題：', question => {
  rl.close()
  if (!question.trim()) {
    console.error('問題不得為空')
    process.exit(1)
  }
  printResult(castHexagram(hexagrams, question.trim()))
})

function printResult(r) {
  console.log(`\n問題：${r.question}\n`)
  console.log('六爻（初爻在下）：')
  for (let i = 5; i >= 0; i--) {
    console.log(`  第${i + 1}爻  ${lineSymbol(r.lines[i])}`)
  }
  console.log(`\n本卦：第 ${r.present.number} 卦  ${r.present.name}`)
  console.log(`卦辭：${r.present.judgment}`)
  if (r.future) {
    console.log(`\n變卦：第 ${r.future.number} 卦  ${r.future.name}`)
    console.log(`卦辭：${r.future.judgment}`)
  }
}

function lineSymbol(l) {
  switch (l.value) {
    case 9: return '━━━○━━━  老陽（變）'
    case 6: return '━━━×━━━  老陰（變）'
    case 7: return '━━━━━━━  少陽'
    case 8: return '━━━  ━━━  少陰'
  }
}
