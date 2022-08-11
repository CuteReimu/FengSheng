package main

import (
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/game"
)

func main() {
	totalCount := config.GetTotalPlayerCount()
	robotCount := config.GetRobotPlayerCount()
	g := &game.Game{}
	g.Start(totalCount, robotCount)
}
