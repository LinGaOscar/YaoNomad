package model

// Line 代表一次三硬幣擲出的爻結果
type Line struct {
	Value    int  // 三枚硬幣加總：6/7/8/9
	Yin      bool // true = 陰爻（偶數）
	Changing bool // true = 老陰(6) 或 老陽(9)，為變爻
}

// Hexagram 代表 64 卦之一
type Hexagram struct {
	Number   int
	Name     string
	Upper    int       // 上卦三位碼 0-7
	Lower    int       // 下卦三位碼 0-7
	Judgment string
	Lines    [6]string // index 0 = 初爻
}

// CastResult 代表一次完整卜卦結果
type CastResult struct {
	Question string
	Lines    [6]Line
	Present  Hexagram
	Future   *Hexagram // 無變爻時為 nil
}
