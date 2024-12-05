package slot

import (
	"fmt"
	"math/rand"
)

type ReelInfo struct {
	Reels [][]int
	Idxs  []int
}

// 取得隨機盤面
func createRandReel(reelMap map[string][][]int, reelName string, reelSize []int, cheat []int) (*ReelInfo, error) {
	// 檢查盤面是否存在
	reels, ok := reelMap[reelName]
	if !ok {
		return nil, fmt.Errorf("CreateRandReel error: reelName:%v nil", reelName)
	}

	isCheat := isCheat(reelSize, cheat)

	r := &ReelInfo{
		Reels: make([][]int, len(reelSize)),
		Idxs:  make([]int, len(reelSize)),
	}

	// 遍歷每個轉輪
	for reelIdx, symbols := range reels {
		if isCheat {
			r.Idxs[reelIdx] = cheat[reelIdx]
		} else {
			// 生成隨機起始索引
			startIndex := rand.Intn(len(symbols))
			r.Idxs[reelIdx] = startIndex
		}

		// 提取 reelsize[reelIdx] 個連續值（支援循環）
		r.Reels[reelIdx] = make([]int, reelSize[reelIdx])
		for j := 0; j < reelSize[reelIdx]; j++ {
			r.Reels[reelIdx][j] = symbols[(r.Idxs[reelIdx]+j)%len(symbols)]
		}
	}

	return r, nil
}

// 作弊檢查
func isCheat(reelsize []int, cheat []int) bool {
	// 沒有作弊資料
	if cheat == nil {
		return false
	}

	// 檢查作弊資料
	if len(cheat) != len(reelsize) {
		return false
	}

	return true
}

// 更新 / 轉換 id
func (r *ReelInfo) UpdatePosSymbolPos(row, col int, id int) {
	r.Reels[row][col] = id
}

// 更新 相同id 為新的id
func (r *ReelInfo) UpdateSymbolSameID(tagID int, newID int) {
	for row, reel := range r.Reels {
		for col, sID := range reel {
			if sID == tagID {
				r.Reels[row][col] = newID
			}
		}
	}
}

// 取的Scatter數量
func (r *ReelInfo) GetScatterCount() int {
	count := 0
	for _, reel := range r.Reels {
		for _, sID := range reel {
			if sID == ScatterID {
				count++
			}
		}
	}

	return count
}
