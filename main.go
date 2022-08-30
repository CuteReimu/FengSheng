package main

import (
	_ "github.com/CuteReimu/FengSheng/ai"
	_ "github.com/CuteReimu/FengSheng/card"
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/gm"
)

func main() {
	config.Init()
	gm.Init()
	totalCount := config.GetTotalPlayerCount()
	game.Start(totalCount)
}
