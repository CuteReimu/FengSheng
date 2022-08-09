package main

import (
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/game"
)

func main() {
	totalCount := config.GlobalConfig.GetInt("player.total_count")
	robotCount := config.GlobalConfig.GetInt("player.robot_count")
	g := &game.Game{}
	g.Start(totalCount, robotCount)
}
