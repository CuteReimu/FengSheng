package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIFightPhase[protos.CardType_Wu_Dao] = wuDao
}

func wuDao(e *game.FightPhaseIdle, card game.ICard) bool {
	player := e.WhoseFightTurn
	colors := e.MessageCard.GetColors()
	if e.InFrontOfWhom.Location() == player.Location() && (e.IsMessageCardFaceUp || player.Location() == e.WhoseTurn.Location()) && len(colors) == 1 && colors[0] != protos.Color_Black {
		return false
	}
	players := player.GetGame().GetPlayers()
	var target game.IPlayer
	switch utils.Random.Intn(4) {
	case 0:
		for left := e.InFrontOfWhom.Location() - 1; left != e.InFrontOfWhom.Location(); left-- {
			if left < 0 {
				left += len(players)
			}
			if player.GetGame().GetPlayers()[left].IsAlive() {
				target = players[left]
				break
			}
		}
	case 1:
		for right := e.InFrontOfWhom.Location() + 1; right != e.InFrontOfWhom.Location(); right++ {
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
