package slot

// 獎金計算
func getPayoutCredit(betCredit uint64, odds float64) uint64 {
	return uint64(float64(betCredit) * odds)
}
