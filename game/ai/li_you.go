package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIMainPhase[protos.CardType_Li_You] = liYou
}

func liYou(player interfaces.IPlayer, card interfaces.ICard) bool {
	var chooseHuman bool
	n := utils.Random.Intn(3)
	if n == 0 {
		time.AfterFunc(time.Second, func() {
			game.Post(func() { card.Execute(player.GetGame(), player, player) })
		})
		return true
	} else {
		chooseHuman = n == 1
	}
	var players []interfaces.IPlayer
	for _, p := range player.GetGame().GetPlayers() {
		if p.IsAlive() {
			if _, ok := p.(*game.HumanPlayer); ok == chooseHuman {
				players = append(players, p)
			}
		}
	}
	if len(players) == 0 {
		players = player.GetGame().GetPlayers()
	}
	if len(players) == 0 {
		return false
	}
	p := players[utils.Random.Intn(len(players))]
	time.AfterFunc(time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, p) })
	})
	return true
}
