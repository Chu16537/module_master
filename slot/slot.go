package slot

type ISlot interface {
	// 取得隨機盤面
	CreateRandReel(reelName string, cheat []int) (*ReelInfo, error)

	// 取得結果資料
	GetResultInfo(betCredit uint64, reelInfo ReelInfo, opt *ResultOpt) *ResultInfo
}

type ResultOpt struct {
	PayLines [][]int
}

// 分成 MainGame 跟 SubGame 可以讓別人知道 原本轉盤是多少結果，
type Result struct {
	BaseGame []ResultInfo // 主要輪盤
	FreeGame []ResultInfo // 免費遊戲
}

type ResultInfo struct {
	ReelInfo          ReelInfo     // 輪帶資料
	PayoutResult      PayoutResult // 派彩結果
	CustomData        interface{}  // 額外資料 (ex.4800功能轉盤)
	TotalPayoutCredit uint64       // 總派獎金額
}

type PayoutResult struct {
	Lines []PayoutResultLine //  line
}
