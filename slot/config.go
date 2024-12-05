package slot

type Config struct {
	GameSettings    []*GameSetting
	CustomDataGame1 *CustomDataGame1
}

type GameSetting struct {
	GameID       int
	ReelSize     []int
	ReelSettings []ReelSetting
	Symbols      []*Symbol // 圖標資料
	PayLine      [][]int   // 中獎線
}

type ReelSetting struct {
	Name  string
	Reels [][]int
}

func (g *GameSetting) ReelSettingToMap() map[string][][]int {
	m := map[string][][]int{}
	for _, v := range g.ReelSettings {
		m[v.Name] = v.Reels
	}
	return m
}

func (g *GameSetting) SymbolToMap() map[int]Symbol {
	m := map[int]Symbol{}
	for _, v := range g.Symbols {
		m[v.ID] = *v
	}
	return m
}

type CustomDataGame1 struct {
	RandomWildCount            []int
	BaseGameRandomWildRateBase []int
	BaseGameRandomWildRateFree []int
	FreeGameRandomWildRateBase []int
	FreeGameRandomWildRateFree []int
}
