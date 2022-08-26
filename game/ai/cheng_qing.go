package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIMainPhase[protos.CardType_Cheng_Qing] = chengQing
}

type playerAndCard struct {
	player interfaces.IPlayer
	card   interfaces.ICard
}

func chengQing(player interfaces.IPlayer, card interfaces.ICard) bool {
	var playerAndCards []playerAndCard
	for _, c := range player.GetMessageCards() {
		if utils.IsColorIn(protos.Color_Black, c.GetColor()) {
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
