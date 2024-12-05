package slot_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Chu16537/module_master/myaml"
	"github.com/Chu16537/module_master/slot"
)

type userData struct {
	TotalBetCredit    uint64
	TotalPayoutCredit uint64
	GreenScarab       int
	BlueScarab        int
}

var (
	c         *slot.Config
	symbolMap map[int]slot.Symbol
	times            = 100
	betCredit uint64 = 100
	ud               = &userData{
		GreenScarab: 0,
		BlueScarab:  0,
	}
)

func newConfig() (*slot.Config, error) {
	c = new(slot.Config)
	err := myaml.Read("config.yaml", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestXxx(t *testing.T) {
	st := time.Now()
	fmt.Println("開始時間", st)

	newConfig()
	slot.Init(c)

	symbolMap = c.GameSettings[0].SymbolToMap()

	totalBet := betCredit * uint64(times)
	var totalPayout uint64
	for i := 0; i < times; i++ {
		_, total := game1()
		totalPayout += total
		// r, total := game1()
		// fmt.Println(totalPayout, r.BaseGame[0])
	}

	fmt.Println("RTP:", float64(totalPayout*100)/float64(totalBet), "%", totalPayout, totalBet)

	fmt.Println("總計時間時間", time.Now().Sub(st))
}

func game1() (*slot.Result, uint64) {
	game, err := slot.CreateSlotGame(1)
	if err != nil {
		fmt.Println("CreateSlotGame", err)
		return nil, 0
	}

	// base
	baseReelInfo, err := game.CreateRandReel("base", nil)
	if err != nil {
		fmt.Println("CreateRandReel", err)
		return nil, 0
	}

	// 結果
	gameResultOpt := &slot.ResultOpt{
		PayLines: c.GameSettings[0].PayLine,
	}

	baseResult := game.GetResultInfo(betCredit, *baseReelInfo, gameResultOpt)

	// fmt.Println("game1 baseResult", baseResult)

	// 更換 scarab
	transScarab(baseReelInfo)

	// scarab
	scarabResult := scarabGame(game, ud, betCredit, baseReelInfo, false)

	if scarabResult != nil {
		baseResult.CustomData = scarabResult
	}

	// free
	freeResults := freeGame(game, betCredit, baseReelInfo)

	gameResult := &slot.Result{
		BaseGame: []slot.ResultInfo{*baseResult},
		FreeGame: freeResults,
	}

	var totalPayout uint64 = 0
	totalPayout += baseResult.TotalPayoutCredit
	fmt.Println("totalPayout 1", baseResult, totalPayout)
	if scarabResult != nil {
		totalPayout += scarabResult.TotalPayoutCredit
		fmt.Println("totalPayout 2", scarabResult, totalPayout)
	}
	for _, v := range gameResult.FreeGame {
		totalPayout += v.TotalPayoutCredit
		fmt.Println("totalPayout 3", v, totalPayout)
	}

	return gameResult, totalPayout
}

// 更換scarab
func transScarab(reelInfo *slot.ReelInfo) {
	for row, reel := range reelInfo.Reels {
		for col, symbolID := range reel {
			if symbolID != 100 {
				continue
			}

			randNum := rand.Intn(100)
			if randNum >= 50 {
				reelInfo.UpdatePosSymbolPos(row, col, symbolMap[101].ID)
			}
		}
	}
}

// scarab 遊戲
func scarabGame(game slot.ISlot, userData *userData, betCredit uint64, ri *slot.ReelInfo, isFree bool) *slot.ResultInfo {
	// 更新 userdata
	for _, reel := range ri.Reels {
		for _, symbolID := range reel {
			if symbolID == 100 {
				userData.GreenScarab++
			}
			if symbolID == 101 {
				userData.BlueScarab++
			}

		}
	}

	// 檢查是否可以執行遊戲
	if userData.BlueScarab < 3 || userData.GreenScarab < 3 {
		return nil
	}

	reelInfo, err := game.CreateRandReel("scarab", nil)
	if err != nil {
		return nil
	}

	// 亂數wild
	if userData.GreenScarab >= 3 {
		userData.GreenScarab = 0

		// 取得個數
		randomWildCount := c.CustomDataGame1.RandomWildCount
		// 取得機率
		randomWildRate := c.CustomDataGame1.BaseGameRandomWildRateBase
		if isFree {
			randomWildRate = c.CustomDataGame1.BaseGameRandomWildRateFree
		}

		greenScarab(reelInfo, randomWildCount, randomWildRate)
	}

	// wild x2
	isBlueScarab := false
	if userData.BlueScarab >= 3 {
		userData.BlueScarab = 0
		isBlueScarab = true
		blueScarab(reelInfo)
	}

	// 結果
	gameResultOpt := &slot.ResultOpt{
		PayLines: c.GameSettings[0].PayLine,
	}

	result := game.GetResultInfo(betCredit, *reelInfo, gameResultOpt)

	if isBlueScarab {
		for _, lineResult := range result.PayoutResult.Lines {
			for _, symbolID := range lineResult.SymbolIDs {
				// 代表是x2的
				if symbolID == 102 {
					result.TotalPayoutCredit += lineResult.PayoutCredit
					lineResult.PayoutCredit += lineResult.PayoutCredit
					break
				}
			}
		}
	}

	return result
}

func greenScarab(ri *slot.ReelInfo, randomWildCount []int, randomWildRate []int) {
	// 計算wild 要出現的數量
	randomNum := rand.Intn(randomWildRate[len(randomWildRate)-1])
	idx := 0
	for i, rate := range randomWildRate {
		if randomNum <= rate {
			idx = i
			break
		}
	}

	// 取得數量
	count := randomWildCount[idx]

	// 產生的 wild位置
	var wildPos [][]int
	for row, reel := range ri.Reels {
		// 第1軸不可以有wild
		if row == 0 {
			continue
		}

		for col := range reel {
			wildPos = append(wildPos, []int{row, col})
		}
	}

	// 隨機換位
	for i := range wildPos {
		j := rand.Intn(i + 1)
		wildPos[i], wildPos[j] = wildPos[j], wildPos[i]
	}

	for i := 0; i < count; i++ {
		// 要更換的位置
		ri.UpdatePosSymbolPos(wildPos[i][0], wildPos[i][1], symbolMap[slot.WildID].ID)
	}
}

func blueScarab(ri *slot.ReelInfo) {
	ri.UpdateSymbolSameID(slot.WildID, 102)
}

func freeGame(game slot.ISlot, betCredit uint64, ri *slot.ReelInfo) []slot.ResultInfo {
	var r []slot.ResultInfo

	count := ri.GetScatterCount()

	// 沒有Scatter
	if count <= 0 {
		return r
	}

	// 假如超過上限 就取得上限
	if count >= len(symbolMap[slot.ScatterID].FreeTimes) {
		count = len(symbolMap[slot.ScatterID].FreeTimes) - 1
	}

	// -1 是因為陣列從0開始
	freeTimes := symbolMap[slot.ScatterID].FreeTimes[count-1]

	for i := 0; i < freeTimes; i++ {
		reelInfo, err := game.CreateRandReel("free", nil)
		if err != nil {
			continue
		}

		// 假如有甲蟲 免費遊戲+1
		for _, reel := range reelInfo.Reels {
			for _, symbolID := range reel {
				if symbolID == 100 {
					freeTimes++
				}
			}
		}

		// 結果
		gameResultOpt := &slot.ResultOpt{
			PayLines: c.GameSettings[0].PayLine,
		}

		result := game.GetResultInfo(betCredit, *reelInfo, gameResultOpt)

		r = append(r, *result)
	}

	return r
}
