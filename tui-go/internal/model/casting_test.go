package model

import (
	"math/rand"
	"testing"
)

// fixedSource 提供固定輸出的隨機源，用於可重現測試
type fixedSource struct {
	seq []int64
	idx int
}

func (fs *fixedSource) Int63() int64 {
	v := fs.seq[fs.idx%len(fs.seq)]
	fs.idx++
	return v
}
func (fs *fixedSource) Seed(int64) {}

// newFixed 建立固定輸出的 *rand.Rand
// val=1 → Int63 返回 MaxInt64 → Intn(2)=1 → coin()=3（背面）
// val=0 → Int63 返回 0 → Intn(2)=0 → coin()=2（正面）
func newFixed(val int64) *rand.Rand {
	var seq []int64
	if val == 0 {
		seq = []int64{0}
	} else {
		seq = []int64{9223372036854775807} // math.MaxInt64
	}
	return rand.New(&fixedSource{seq: seq})
}

func TestCastHexagram_allTails(t *testing.T) {
	// 全背面(3)：sum=9 老陽 → 本卦乾(#1)，變卦坤(#2)
	hexagrams := loadTestHexagrams(t)
	result, err := CastHexagram("測試", hexagrams, newFixed(1))
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		if l.Value != 9 {
			t.Errorf("爻 %d 應為 9，實得 %d", i, l.Value)
		}
		if l.Yin {
			t.Errorf("爻 %d 應為陽爻（Yin=false）", i)
		}
		if !l.Changing {
			t.Errorf("爻 %d 應為變爻（Changing=true）", i)
		}
	}
	if result.Present.Number != 1 {
		t.Errorf("本卦應為乾(#1)，實得 #%d %s", result.Present.Number, result.Present.Name)
	}
	if result.Future == nil {
		t.Fatal("六爻全為老陽，應有變卦")
	}
	if result.Future.Number != 2 {
		t.Errorf("變卦應為坤(#2)，實得 #%d %s", result.Future.Number, result.Future.Name)
	}
}

func TestCastHexagram_allHeads(t *testing.T) {
	// 全正面(2)：sum=6 老陰 → 本卦坤(#2)，變卦乾(#1)
	hexagrams := loadTestHexagrams(t)
	result, err := CastHexagram("測試", hexagrams, newFixed(0))
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	if result.Present.Number != 2 {
		t.Errorf("本卦應為坤(#2)，實得 #%d", result.Present.Number)
	}
	if result.Future == nil || result.Future.Number != 1 {
		t.Error("變卦應為乾(#1)")
	}
}

func TestCastHexagram_noChanging(t *testing.T) {
	// seq [0,0,MaxInt64] → coins 2,2,3 → sum=7 少陽，重複六次
	maxInt64 := int64(9223372036854775807)
	hexagrams := loadTestHexagrams(t)
	r := rand.New(&fixedSource{seq: []int64{0, 0, maxInt64}})
	result, err := CastHexagram("測試", hexagrams, r)
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		if l.Value != 7 {
			t.Errorf("爻 %d 應為 7，實得 %d", i, l.Value)
		}
		if l.Changing {
			t.Errorf("爻 %d 不應為變爻", i)
		}
	}
	if result.Future != nil {
		t.Error("無變爻時 Future 應為 nil")
	}
}

func TestCastHexagram_lineFields(t *testing.T) {
	// 驗證 Yin 與 Changing 欄位與 Value 一致
	hexagrams := loadTestHexagrams(t)
	r := rand.New(rand.NewSource(42))
	result, err := CastHexagram("測試", hexagrams, r)
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		wantYin := l.Value%2 == 0
		if l.Yin != wantYin {
			t.Errorf("爻 %d Yin=%v 應為 %v（value=%d）", i, l.Yin, wantYin, l.Value)
		}
		wantChanging := l.Value == 6 || l.Value == 9
		if l.Changing != wantChanging {
			t.Errorf("爻 %d Changing=%v 應為 %v（value=%d）", i, l.Changing, wantChanging, l.Value)
		}
	}
}
