package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AISendPhase[protos.CardType_Po_Yi] = poYi
}

func poYi(e *game.SendPhaseIdle, card game.ICard) bool {
	player := e.InFrontOfWhom
	if player.Location() == e.WhoseTurn.Location() {
		return false
	}
	if e.IsMessageCardFaceUp {
		return false
	}
	if utils.Random.Intn(2) == 0 {
		return false
	}
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player) })
	})
	return true
}
