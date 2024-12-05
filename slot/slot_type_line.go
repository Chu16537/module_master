package slot

import (
	"fmt"
)

type slotLine struct {
	config    *GameSetting
	reelMap   map[string][][]int
	symbolMap map[int]Symbol
}

type PayoutResultLine struct {
	ID           int    // 中獎線 ID
	SymbolIDs    []int  // 該中獎線的圖案
	PayoutCredit uint64 // 中獎金額
}

func newGameLine(config *GameSetting) (ISlot, error) {
	s := &slotLine{
		config:    config,
		reelMap:   config.ReelSettingToMap(),
		symbolMap: config.SymbolToMap(),
	}

	for _, reels := range s.reelMap {
		// 檢查 reelsize 是否與 reel 尺寸一致
		if len(reels) != len(config.ReelSize) {
			return nil, fmt.Errorf("newGameLine error: len fail reel:%v reelsize:%v", len(reels), len(config.ReelSize))
		}

		// 檢查每個輪帶是否足夠長
		for i, symbols := range reels {
			if len(symbols) < config.ReelSize[i] {
				return nil, fmt.Errorf("NewGameLine error: symbols%d length:%v is less than reelsize%d:%v", i, len(symbols), i, config.ReelSize[i])
			}
		}
	}

	return s, nil
}

// 取得隨機盤面
func (s *slotLine) CreateRandReel(reelName string, cheat []int) (*ReelInfo, error) {
	r, err := createRandReel(s.reelMap, reelName, s.config.ReelSize, cheat)

	return r, err
}

// 取得盤面結果資料
func (s *slotLine) GetResultInfo(betCredit uint64, reelInfo ReelInfo, opt *ResultOpt) *ResultInfo {
	var (
		ls                []PayoutResultLine
		totalPayoutCredit uint64
	)

	for lineID, linePoss := range opt.PayLines {
		symbols := make([]Symbol, len(reelInfo.Reels))

		// 取得指定位置上的symbol
		for row, col := range linePoss {
			symbols[row] = s.symbolMap[reelInfo.Reels[row][col]]
		}

		// 取得中獎個數
		_, odd, count := s.getPayLineInfo(symbols)

		if count > 0 {
			// 記錄中獎圖標
			symbolIDs := make([]int, count)
			for i, v := range symbols {
				if i >= count {
					break
				}
				symbolIDs[i] = v.ID
			}

			r := PayoutResultLine{
				ID:           lineID,
				SymbolIDs:    symbolIDs,
				PayoutCredit: getPayoutCredit(betCredit, odd),
			}

			ls = append(ls, r)

			totalPayoutCredit += r.PayoutCredit
		}

	}

	r := &ResultInfo{
		ReelInfo: reelInfo,
		PayoutResult: PayoutResult{
			Lines: ls,
		},
		TotalPayoutCredit: totalPayoutCredit,
	}

	return r
}

// 取得中獎個數
func (s *slotLine) getPayLineInfo(symbols []Symbol) (int, float64, int) {
	if len(symbols) == 0 {
		return 0, 0, 0
	}

	var baseSymbol *Symbol
	for _, symbol := range symbols {
		if symbol.IsWild() {
			// 如果是 wild，繼續往後找
			continue
		}
		baseSymbol = &symbol
		break
	}

	// 如果全是 wild
	if baseSymbol == nil {
		return 0, 0, 0
	}

	count := 0
	// 判斷其他符號是否與 baseSymbol 一致（或是 wild）
	for _, symbol := range symbols {
		if symbol.IsWild() || symbol.Equal(baseSymbol) {
			count++
			continue
		}
		break
	}

	// 防呆
	if count >= len(baseSymbol.Odds) {
		count = len(baseSymbol.Odds)
	}

	// 減1 因為 array 從0開始 假如count=2 代表只有中2個 但是用baseSymbol.Odds()[count] 會變成中3個的獎品
	odd := baseSymbol.Odds[count-1]
	// 沒有中獎
	if odd == 0.0 {
		return 0, 0, 0
	}

	return baseSymbol.ID, odd, count
}
