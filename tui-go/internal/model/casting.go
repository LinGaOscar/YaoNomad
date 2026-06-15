package model

import "math/rand"

// CastHexagram 執行一次三硬幣卜卦
// r 為隨機源，測試時注入固定種子以確保可重現
func CastHexagram(question string, hexagrams []Hexagram, r *rand.Rand) (CastResult, error) {
	var lines [6]Line
	for i := range lines {
		sum := coin(r) + coin(r) + coin(r)
		lines[i] = Line{
			Value:    sum,
			Yin:      sum%2 == 0,
			Changing: sum == 6 || sum == 9,
		}
	}

	upper, lower := trigramCodes(lines)
	present, err := Lookup(hexagrams, upper, lower)
	if err != nil {
		return CastResult{}, err
	}

	future, err := computeFuture(hexagrams, lines)
	if err != nil {
		return CastResult{}, err
	}

	return CastResult{
		Question: question,
		Lines:    lines,
		Present:  present,
		Future:   future,
	}, nil
}

// coin 模擬一枚硬幣：正面(字面)=2，背面=3
func coin(r *rand.Rand) int {
	if r.Intn(2) == 0 {
		return 2
	}
	return 3
}

// trigramCodes 計算六爻的上下卦三位碼
// bit0 = 底爻，yang=1，yin=0
func trigramCodes(lines [6]Line) (upper, lower int) {
	yang := func(l Line) int {
		if l.Yin {
			return 0
		}
		return 1
	}
	lower = yang(lines[0]) | yang(lines[1])<<1 | yang(lines[2])<<2
	upper = yang(lines[3]) | yang(lines[4])<<1 | yang(lines[5])<<2
	return upper, lower
}

// computeFuture 計算變卦；若無變爻則回傳 nil
func computeFuture(hexagrams []Hexagram, lines [6]Line) (*Hexagram, error) {
	hasChanging := false
	for _, l := range lines {
		if l.Changing {
			hasChanging = true
			break
		}
	}
	if !hasChanging {
		return nil, nil
	}

	var futureLines [6]Line
	for i, l := range lines {
		yin := l.Yin
		if l.Changing {
			yin = !l.Yin
		}
		futureLines[i] = Line{Value: l.Value, Yin: yin, Changing: l.Changing}
	}

	upper, lower := trigramCodes(futureLines)
	h, err := Lookup(hexagrams, upper, lower)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
