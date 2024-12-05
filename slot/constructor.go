package slot

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	c           *Config
	gameSetting map[int]*GameSetting
)

func Init(config *Config) {
	c = config
	gameSetting = map[int]*GameSetting{}

	for _, v := range config.GameSettings {
		gameSetting[v.GameID] = v
	}

	// 設定亂數種子
	rand.Seed(time.Now().UnixNano())

}

func GetGameSetting(gameID int) (*GameSetting, error) {
	if config, ok := gameSetting[gameID]; ok {
		return config, nil
	}
	return nil, fmt.Errorf("gameID:%v not found", gameID)
}

func GetCustomData(gameID int) (interface{}, error) {
	switch gameID {
	case 1:
		return c.CustomDataGame1, nil
	default:
		return nil, fmt.Errorf("GetCustomData gameID:%v is nil", gameID)
	}
}

func CreateSlotGame(gameID int) (ISlot, error) {
	config, ok := gameSetting[gameID]
	if !ok {
		return nil, fmt.Errorf("gameID:%v not found", gameID)
	}

	switch gameID {
	case 1:
		return newGameLine(config)
	default:
		return nil, fmt.Errorf("CreateSlot gameID:%v is nil", gameID)
	}
}
