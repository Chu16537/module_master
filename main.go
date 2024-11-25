package main

import (
	"fmt"

	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/proto/db"
)

func main() {
	a := db.GameConfig1{
		EnterBalance:      99,
		MaxPlayerCount:    2,
		Chips:             []uint64{1, 2, 3},
		UpperBetLimitZone: []uint64{4, 5, 6},
		Odds:              []float64{7.7, 8.8, 9.9},
		EffectiveBet:      []float64{100.1, 200.2, 300.3},
		MaxBet:            0,
		Rtp:               98.5,
	}

	b, err := mjson.Marshal(a)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("b", b)
}
