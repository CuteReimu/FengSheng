package main

import (
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/game"
	_ "github.com/CuteReimu/FengSheng/game/ai"
	_ "github.com/CuteReimu/FengSheng/game/card"
)

func main() {
	config.Init()
	totalCount := config.GetTotalPlayerCount()
	robotCount := config.GetRobotPlayerCount()
	g := &game.Game{}
	g.Start(totalCount, robotCount)
}
