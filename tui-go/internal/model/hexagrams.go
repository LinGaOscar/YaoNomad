package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type hexagramJSON struct {
	Number   int       `json:"number"`
	Name     string    `json:"name"`
	Upper    int       `json:"upper"`
	Lower    int       `json:"lower"`
	Judgment string    `json:"judgment"`
	Lines    [6]string `json:"lines"`
}

// LoadHexagrams 從 JSON 檔案載入 64 卦資料
func LoadHexagrams(path string) ([]Hexagram, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀取卦象資料失敗: %w", err)
	}
	var entries []hexagramJSON
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("解析卦象資料失敗: %w", err)
	}
	hexagrams := make([]Hexagram, len(entries))
	for i, e := range entries {
		hexagrams[i] = Hexagram{
			Number:   e.Number,
			Name:     e.Name,
			Upper:    e.Upper,
			Lower:    e.Lower,
			Judgment: e.Judgment,
			Lines:    e.Lines,
		}
	}
	return hexagrams, nil
}

// Lookup 在卦象列表中查找指定上下卦的卦象
func Lookup(hexagrams []Hexagram, upper, lower int) (Hexagram, error) {
	for _, h := range hexagrams {
		if h.Upper == upper && h.Lower == lower {
			return h, nil
		}
	}
	return Hexagram{}, fmt.Errorf("找不到卦象: upper=%d lower=%d", upper, lower)
}
