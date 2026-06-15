# YaoNomad 卜卦系統建置計劃

## 一、專案概述

YaoNomad 是一套以《易經》卜卦為核心的多平台系統，在「最小化依賴」的前提下，提供穩定可攜的卜卦體驗。系統分為核心卜卦邏輯與兩種介面實作：Web PWA 與 Go TUI。

---

## 二、設計原則

1. 單一核心邏輯，負責卜卦與卦象計算，與任何 UI 無關。
2. 介面層各自實作動畫與互動，透過核心 API 取得結果。
3. 避免不必要框架與依賴，保持程式碼簡單可讀。
4. 所有功能須在 Windows、Linux、macOS 上可執行。
5. 優先確保離線可用。

---

## 三、系統架構與目錄規劃

```
/core-spec/
  README.md          ← 語言無關的演算法規格（單一真實來源）
  hexagrams.json     ← 64 卦資料正本

/web/
  index.html
  styles.css
  main.js
  manifest.json
  service-worker.js
  core/
    models.js        ← Line、Hexagram、CastResult JSDoc 型別
    casting.js       ← castHexagram(question, randFn?)
    hexagrams.js     ← loadHexagrams() + lookup(lines)
  data/
    hexagrams.json   ← 從 core-spec 複製

/tui-go/
  go.mod             ← module github.com/LinGaOscar/yaonomad
  cmd/yaonomad-tui/
    main.go
  internal/
    model/
      types.go       ← Line、Hexagram、CastResult
      casting.go     ← CastHexagram(question, rand)
      hexagrams.go   ← LoadHexagrams(path) + Lookup(lines)
    tui/
      model.go
      update.go
      view.go
      styles.go
  data/
    hexagrams.json   ← 從 core-spec 複製
```

---

## 四、核心資料模型

### Line（爻）
- `value`：6 / 7 / 8 / 9
- `yin`：偶數為陰（true）
- `changing`：老陰(6) 或 老陽(9) 為變爻

### Hexagram（卦）
- `number`：文王序列 1–64
- `name`：卦名
- `upper` / `lower`：trigram 三位碼（0–7）
- `judgment`：卦辭
- `lines[6]`：爻辭（index 0 = 初爻）

### CastResult（卜卦結果）
- `question`：使用者問題
- `lines[6]`：六爻
- `present`：本卦
- `future`：變卦（無變爻時為 null）

---

## 五、卜卦演算法

三枚硬幣法：正面(字面)=2，背面=3。擲三枚加總得 6/7/8/9，自下而上重複六次。

| 爻值 | 爻型 | 變爻 |
|------|------|------|
| 6 | 老陰，陰爻 | ✓ → 陽 |
| 7 | 少陽，陽爻 | ✗ |
| 8 | 少陰，陰爻 | ✗ |
| 9 | 老陽，陽爻 | ✓ → 陰 |

查卦：線性掃描 64 筆，以 `upper` 與 `lower` trigram 碼比對。

---

## 六、Web PWA 介面（/web）

### 使用流程
1. **問題輸入**：顯示輸入框與「開始卜卦」按鈕。
2. **卜卦動畫**：狀態機控制六次擲卦，逐爻顯示，CSS 動畫營造儀式感。
3. **結果呈現**：顯示本卦、變卦（若有）、卦辭與爻辭。提供重新卜卦操作。

### PWA 設定
- `manifest.json`：應用名稱、圖示、standalone 模式。
- `service-worker.js`：安裝階段快取所有靜態資源，攔截請求優先回傳快取。

---

## 七、Go TUI 介面（/tui-go）

### 狀態設計（tui/model.go）
- `StateInput`：問題輸入
- `StateCasting`：卜卦動畫（逐爻計時更新）
- `StateResult`：結果顯示

### 畫面規劃（tui/view.go）
1. **輸入畫面**：標題列 + 問題輸入框。
2. **動畫畫面**：六爻逐步顯示，未完成爻以佔位符表示。
3. **結果畫面**：本卦 + 變卦（若有）+ 卦辭 + 爻辭，按鍵返回輸入。

---

## 八、開發時程

1. **第 1 週**：撰寫 `/core-spec/README.md`，建立 `hexagrams.json`（AI 產生基礎版本）。
2. **第 2 週**：Go 核心（types.go、casting.go、hexagrams.go）+ 測試。
3. **第 3 週**：Web 核心（models.js、casting.js、hexagrams.js）+ 測試。
4. **第 4 週**：Go TUI 三段畫面與狀態機。
5. **第 5 週**：Web PWA UI + Service Worker。
6. **第 6 週**：整合測試、跨平台驗證、重構。

---

## 九、品質策略

- Go 與 JS 兩端均須對 `castHexagram` 撰寫測試，使用固定亂數種子驗證結果一致。
- 核心與介面修改不得互相影響。
- `/core-spec` 為唯一真實來源，變更須同步兩端副本。
