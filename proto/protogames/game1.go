package protogames

// 遊戲1 結果
type Game1GameResult struct {
	Cards   []int `json:"cards" bson:"cards"`       // 牌 0龍 1虎
	WinZone []int `json:"win_zone" bson:"win_zone"` // 贏分結果
}
