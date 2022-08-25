package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AISendPhase[protos.CardType_Po_Yi] = poYi
}

func poYi(player interfaces.IPlayer, card interfaces.ICard) bool {
	if player.Location() == player.GetGame().GetWhoseTurn() {
		return false
	}
	if player.GetGame().IsMessageCardFaceUp() {
		return false
	}
	if utils.Random.Intn(2) == 0 {
		return false
	}
	time.AfterFunc(time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player) })
	})
	return true
}
