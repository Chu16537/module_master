package slot

type Symbol struct {
	ID   int       // 圖騰編號
	Flag uint32    // 圖騰類型
	Odds []float64 // 倍率表

	// 以下是客製資料
	FreeTimes  []int   // free game 次數
	Multiplier float64 // 倍率 比如x2
}

// 檢查圖騰是否符合指定的 flag
func (s *Symbol) Match(mask uint32) bool {
	return s != nil && s.Flag == (s.Flag&mask)
}

// 檢查兩個圖騰是否相等
func (s *Symbol) Equal(s2 *Symbol) bool {
	return s.ID == s2.ID && s.Flag == s2.Flag
}

// 判斷flag是不是wild
func (s *Symbol) IsWild() bool {
	return s.Flag == WildFlag
}
