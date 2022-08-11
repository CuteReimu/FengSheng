package main

import (
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/game"
)

func main() {
	config.Init()
	totalCount := config.GetTotalPlayerCount()
	robotCount := config.GetRobotPlayerCount()
	g := &game.Game{}
	g.Start(totalCount, robotCount)
}
