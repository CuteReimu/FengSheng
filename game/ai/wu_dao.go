package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIFightPhase[protos.CardType_Wu_Dao] = wuDao
}

func wuDao(player interfaces.IPlayer, card interfaces.ICard) bool {
	colors := player.GetGame().GetCurrentMessageCard().GetColor()
	if player.GetGame().GetWhoseSendTurn() == player.Location() && (player.GetGame().IsMessageCardFaceUp() || player.Location() == player.GetGame().GetWhoseTurn()) && len(colors) == 1 && colors[0] != protos.Color_Black {
		return false
	}
	players := player.GetGame().GetPlayers()
	var target interfaces.IPlayer
	switch utils.Random.Intn(4) {
	case 0:
		for left := player.GetGame().GetWhoseSendTurn() - 1; left != player.GetGame().GetWhoseSendTurn(); left-- {
			if left < 0 {
				left += len(players)
			}
			if player.GetGame().GetPlayers()[left].IsAlive() {
				target = players[left]
				break
			}
		}
	case 1:
		for right := player.GetGame().GetWhoseSendTurn() + 1; right != player.GetGame().GetWhoseSendTurn(); right++ {
			if right >= len(players) {
				right -= len(players)
			}
			if player.GetGame().GetPlayers()[right].IsAlive() {
				target = players[right]
				break
			}
		}
	default:
		return false
	}
	if target == nil {
		return false
	}
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, target) })
	})
	return true
}
