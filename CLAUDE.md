# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 專案概述

YaoNomad 是以《易經》三枚硬幣卜卦法為核心的多平台系統，包含 Web PWA（`/web`）與 Go TUI（`/tui-go`）兩種介面實作。核心規格與介面實作嚴格分離。

**目前進度**：`core-spec` 與兩端核心卜卦邏輯（model/casting/hexagrams）已完成並通過測試。`tui-go/cmd/yaonomad-tui/` 與 `tui-go/internal/tui/`（Bubble Tea 畫面）尚未建立；`web/` 僅有最小可用版的 `index.html`，`manifest.json`／`service-worker.js` 尚未建立。詳細規劃見 `PLAN.md`。

## 架構原則

1. **`/core-spec` 為單一真實來源**：演算法或資料模型變更，必須先更新 `/core-spec/README.md`，再同步更新兩端實作。
2. **核心與介面分離**：`castHexagram` 邏輯不得依賴任何 UI 框架；兩個平台須各自獨立實作，共用相同的演算法與資料模型。
3. **離線優先**：Web 端透過 Service Worker 快取所有靜態資源，保證離線可卜卦。

## 核心資料模型

三個核心型別（Go 與 JS 端均需對應實作）：

- `Line`：爻值（6/7/8/9）、`yin`（偶數為陰）、`changing`（6 或 9 為變爻）
- `Hexagram`：`number`（1–64）、`name`、`upper`/`lower`（trigram 0–7）、`judgment`、`lines[6]`
- `CastResult`：`question`、`lines[6]`、`present`、`future`（無變爻時為 null）

## 狀態機

所有介面實作均須遵循此三段流程：

**問題輸入** → **卜卦動畫（逐爻顯示）** → **結果呈現**

## 開發指令

### Go（`/tui-go`，獨立 module：`github.com/LinGaOscar/yaonomad`）

`tui-go` 有自己的 `go.mod`，須先 `cd tui-go` 再執行 Go 指令（從 repo 根目錄直接 `go test ./tui-go/...` 會失敗）。

```bash
cd tui-go && go test ./...                                  # 全部測試
cd tui-go && go test ./internal/model/ -run TestName -v     # 執行單一測試
cd tui-go && go build ./...                                 # 編譯檢查（目前無 main package，不會產出執行檔）
```

### Web（`/web`，純 JavaScript，無編譯步驟）

測試使用 Node 內建的 `node:test`（`web/package.json` 設定 `"type": "module"` 以支援 `import`/`export`）。

```bash
node --test web/core/*.test.js                                           # 全部測試
node --test --test-name-pattern="<測試名稱>" web/core/casting.test.js   # 執行單一測試
```

直接開啟 `web/index.html`，或啟動任一靜態伺服器即可在瀏覽器中執行；push 到 `main` 會由 `.github/workflows/deploy-pages.yml` 自動將 `/web` 部署到 GitHub Pages。

## 卜卦演算法（三硬幣法）

| 硬幣面 | 陰陽 | 數值 |
|--------|------|------|
| 正面（字面） | 陰 | 2 |
| 背面 | 陽 | 3 |

三枚硬幣加總得爻值：

| 結果 | 爻值 | 爻型 | 是否變爻 |
|------|------|------|----------|
| 3正（2+2+2） | 6 | 老陰，陰爻 | ✓ → 陽 |
| 2正1背（2+2+3） | 7 | 少陽，陽爻 | ✗ |
| 1正2背（2+3+3） | 8 | 少陰，陰爻 | ✗ |
| 3背（3+3+3） | 9 | 老陽，陽爻 | ✓ → 陰 |

- 重複六次，**自下而上**排列（index 0 = 初爻）
- 查卦：線性掃描 64 筆，以 `upper` 與 `lower` trigram 碼比對
- 有變爻則反轉對應爻的陰陽，查得變卦；無變爻則 `future = null`

## 八卦三位碼（bit0 = 底爻，yang=1）

乾=7，坤=0，震=1，巽=6，坎=2，離=5，艮=4，兌=3

## 測試策略

固定亂數種子 → 驗證六爻值（均在 6–9）→ 驗證本卦編號（1–64）→ 驗證變卦編號（若有）。
Go 與 JS 兩端均須撰寫此項測試以確保跨平台一致性。

詳細規格見：`docs/superpowers/specs/2026-06-15-core-design.md`
