package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIMainPhase[protos.CardType_Shi_Tan] = shiTan
}

func shiTan(e *game.MainPhaseIdle, card game.ICard) bool {
	player := e.Player
	var players []game.IPlayer
	for _, p := range player.GetGame().GetPlayers() {
		if p.IsAlive() {
			players = append(players, p)
		}
	}
	if len(players) == 0 {
		return false
	}
	p := players[utils.Random.Intn(len(players))]
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, p) })
	})
	return true
}
