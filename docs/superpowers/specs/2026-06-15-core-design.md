# YaoNomad 核心設計規格

**日期**：2026-06-15
**範圍**：core-spec 規格 + Go TUI 核心 + Web 核心（純 JS）

---

## 一、架構概覽

```
/core-spec/
  README.md          ← 語言無關的演算法規格（本文件）
  hexagrams.json     ← 64 卦正本資料

/tui-go/
  go.mod             ← module github.com/LinGaOscar/yaonomad
  internal/model/
    types.go         ← Line、Hexagram、CastResult 型別定義
    casting.go       ← CastHexagram(question, rand) CastResult
    hexagrams.go     ← LoadHexagrams(path) + Lookup(lines)
  data/
    hexagrams.json   ← 從 core-spec 複製

/web/
  core/
    models.js        ← Line、Hexagram、CastResult JSDoc 型別
    casting.js       ← castHexagram(question, randFn?)
    hexagrams.js     ← loadHexagrams() + lookup(lines)
  data/
    hexagrams.json   ← 從 core-spec 複製
```

**同步規則**：`/core-spec/hexagrams.json` 為唯一正本。任何資料異動須先修改正本，再同步兩端副本。

---

## 二、資料模型

### 2.1 Line（爻）

| 欄位 | Go 型別 | JS 型別 | 說明 |
|------|---------|---------|------|
| `value` | `int` | `number` | 三枚硬幣加總：6 / 7 / 8 / 9 |
| `yin` | `bool` | `boolean` | `true` = 陰爻（value 為偶數） |
| `changing` | `bool` | `boolean` | `true` = 老陰(6) 或 老陽(9)，為變爻 |

### 2.2 Hexagram（卦）

| 欄位 | Go 型別 | JS 型別 | 說明 |
|------|---------|---------|------|
| `number` | `int` | `number` | 文王序列 1–64 |
| `name` | `string` | `string` | 卦名（如「乾」） |
| `upper` | `int` | `number` | 上卦三位碼（0–7，見§4） |
| `lower` | `int` | `number` | 下卦三位碼（0–7，見§4） |
| `judgment` | `string` | `string` | 卦辭 |
| `lines` | `[6]string` | `string[]` | 爻辭，index 0 = 初爻，index 5 = 上爻 |

### 2.3 CastResult（卜卦結果）

| 欄位 | Go 型別 | JS 型別 | 說明 |
|------|---------|---------|------|
| `question` | `string` | `string` | 使用者問題 |
| `lines` | `[6]Line` | `Line[]` | 六爻，index 0 = 初爻 |
| `present` | `Hexagram` | `Hexagram` | 本卦 |
| `future` | `*Hexagram` | `Hexagram\|null` | 變卦；無變爻時為 null |

---

## 三、hexagrams.json 格式

```json
[
  {
    "number": 1,
    "name": "乾",
    "upper": 7,
    "lower": 7,
    "judgment": "元亨利貞。",
    "lines": [
      "初九：潛龍，勿用。",
      "九二：見龍在田，利見大人。",
      "九三：君子終日乾乾，夕惕若，厲，無咎。",
      "九四：或躍在淵，無咎。",
      "九五：飛龍在天，利見大人。",
      "上九：亢龍有悔。"
    ]
  }
]
```

**八卦三位碼對照表**（bit0 = 底爻，yang=1，yin=0）：

| 卦名 | 符號 | 碼值 |
|------|------|------|
| 乾 | ☰ | 7 (111) |
| 兌 | ☱ | 3 (011) |
| 離 | ☲ | 5 (101) |
| 震 | ☳ | 1 (001) |
| 巽 | ☴ | 6 (110) |
| 坎 | ☵ | 2 (010) |
| 艮 | ☶ | 4 (100) |
| 坤 | ☷ | 0 (000) |

---

## 四、卜卦演算法

### 4.1 擲幣產生爻值

重複 6 次（i = 0..5，自下而上）：

```
coin() → 正面(字面) = 2，背面 = 3
sum = coin() + coin() + coin()   // 結果為 6 / 7 / 8 / 9

lines[i] = {
  value:    sum,
  yin:      sum % 2 == 0,        // 6,8 為陰；7,9 為陽
  changing: sum == 6 || sum == 9
}
```

### 4.2 計算卦象

```
yang(line) = line.yin ? 0 : 1

lowerCode = yang(lines[0]) | yang(lines[1])<<1 | yang(lines[2])<<2
upperCode = yang(lines[3]) | yang(lines[4])<<1 | yang(lines[5])<<2

presentHex = lookup(hexagrams, upperCode, lowerCode)
```

### 4.3 計算變卦

```
若 lines 中存在任一 changing == true：
  futureLines[i] = {
    ...lines[i],
    yin: lines[i].changing ? !lines[i].yin : lines[i].yin
  }
  futureHex = lookup(hexagrams, futureUpperCode, futureLowerCode)
否則：
  futureHex = null
```

### 4.4 查卦（線性掃描）

```
lookup(hexagrams, upper, lower):
  for each h in hexagrams:
    if h.upper == upper && h.lower == lower:
      return h
  error("hexagram not found")
```

---

## 五、隨機源設計

為使測試可重現，**隨機源以參數注入**：

- **Go**：`CastHexagram(question string, r *rand.Rand) CastResult`
  - 正式使用：`rand.New(rand.NewSource(time.Now().UnixNano()))`
  - 測試：`rand.New(rand.NewSource(42))`

- **JS**：`castHexagram(question, coinFn = () => Math.random() < 0.5 ? 2 : 3)`
  - `coinFn` 每次呼叫回傳硬幣值 2 或 3
  - 測試：注入固定序列函式

---

## 六、測試策略

固定種子 → 驗證六爻值 → 驗證本卦編號 → 驗證變卦編號（若有）。

**Go 測試範例**：
```go
r := rand.New(rand.NewSource(42))
result := CastHexagram("測試問題", r)
// assert result.Lines[0].Value in {6,7,8,9}
// assert result.Present.Number in 1..64
```

**JS 測試範例**：
```js
const coins = [2,3,3, 2,2,3, 3,3,3, 2,3,2, 3,2,3, 2,2,2]  // 18 個硬幣值（6爻×3枚）
let i = 0
const result = castHexagram("測試問題", () => coins[i++])
```
