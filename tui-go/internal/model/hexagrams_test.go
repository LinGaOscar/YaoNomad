package model

import (
	"testing"
)

func loadTestHexagrams(t *testing.T) []Hexagram {
	t.Helper()
	hexagrams, err := LoadHexagrams("../../data/hexagrams.json")
	if err != nil {
		t.Fatalf("載入卦象失敗: %v", err)
	}
	return hexagrams
}

func TestLoadHexagrams_count(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	if len(hexagrams) != 64 {
		t.Errorf("卦象數量應為 64，實得 %d", len(hexagrams))
	}
}

func TestLoadHexagrams_uniquePairs(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	seen := make(map[[2]int]bool)
	for _, h := range hexagrams {
		key := [2]int{h.Upper, h.Lower}
		if seen[key] {
			t.Errorf("重複的上下卦組合: upper=%d lower=%d", h.Upper, h.Lower)
		}
		seen[key] = true
	}
}

func TestLookup_found(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	h, err := Lookup(hexagrams, 7, 7)
	if err != nil {
		t.Fatalf("查乾卦失敗: %v", err)
	}
	if h.Name != "乾" || h.Number != 1 {
		t.Errorf("乾卦資料錯誤: name=%s number=%d", h.Name, h.Number)
	}
}

func TestLookup_notFound(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	_, err := Lookup(hexagrams, 8, 8)
	if err == nil {
		t.Error("查詢無效碼（8）應回傳錯誤")
	}
}
