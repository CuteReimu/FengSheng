package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"time"
)

func init() {
	game.AI[protos.CardType_Ping_Heng] = pingHeng
}

func pingHeng(player interfaces.IPlayer, card interfaces.ICard) bool {
	var players []interfaces.IPlayer
	for _, p := range player.GetGame().GetPlayers() {
		if p.Location() != player.Location() {
			if h, ok := p.(*game.HumanPlayer); ok {
				players = append(players, h)
			}
		}
	}
	if len(players) == 0 {
		for _, p := range player.GetGame().GetPlayers() {
			if p.Location() != player.Location() {
				players = append(players, p)
			}
		}
	}
	p := players[player.GetGame().GetRandom().Intn(len(players))]
	if card.CanUse(player.GetGame(), player, p) {
		time.AfterFunc(time.Second, func() {
			player.GetGame().Post(func() { card.Execute(player.GetGame(), player, p) })
		})
		return true
	}
	return false
}
