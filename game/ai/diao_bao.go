package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIFightPhase[protos.CardType_Diao_Bao] = diaoBao
}

func diaoBao(player interfaces.IPlayer, card interfaces.ICard) bool {
	colors := player.GetGame().GetCurrentMessageCard().GetColor()
	if player.GetGame().GetWhoseSendTurn() == player.Location() && (player.GetGame().IsMessageCardFaceUp() || player.Location() == player.GetGame().GetWhoseTurn()) && len(colors) == 1 && colors[0] != protos.Color_Black {
		return false
	}
	if utils.Random.Intn(4) != 0 {
		return false
	}
	time.AfterFunc(time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player) })
	})
	return true
}
