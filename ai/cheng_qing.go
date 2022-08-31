package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIMainPhase[protos.CardType_Cheng_Qing] = chengQing
}

type playerAndCard struct {
	player game.IPlayer
	card   game.ICard
}

func chengQing(e *game.MainPhaseIdle, card game.ICard) bool {
	player := e.Player
	var playerAndCards []playerAndCard
	for _, c := range player.GetMessageCards() {
		if utils.IsColorIn(protos.Color_Black, c.GetColors()) {
			playerAndCards = append(playerAndCards, playerAndCard{player, c})
		}
	}
	if len(playerAndCards) == 0 {
		return false
	}
	p := playerAndCards[utils.Random.Intn(len(playerAndCards))]
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, p.player, p.card.GetId()) })
	})
	return true
}
