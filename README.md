# YaoNomad

以《易經》三枚硬幣卜卦法為核心的多平台系統。

## 平台

| 平台 | 目錄 | 技術 |
|------|------|------|
| Web PWA | `/web` | HTML / CSS / JavaScript（無框架） |
| JS TUI | `/tui-js` | Node.js（無框架） |

## 架構

```
/core-spec/        演算法規格與 64 卦資料正本
/web/              PWA 介面（離線可用）
/tui-js/           終端機介面
/docs/             設計文件
```

核心邏輯（卜卦演算法、資料模型）在兩端各自獨立實作，遵循 `/core-spec` 定義的規格。

## 卜卦流程

問題輸入 → 卜卦動畫（逐爻顯示） → 結果呈現（本卦 + 變卦）

## 快速開始

### JS TUI

```bash
node tui-js/main.js
```

### Web

直接開啟 `web/index.html`，或啟動任一靜態伺服器。

## 開發

詳見 [CLAUDE.md](./CLAUDE.md) 與 [設計規格](./docs/superpowers/specs/2026-06-15-core-design.md)。
