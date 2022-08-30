package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

var weiBiTypes = []protos.CardType{protos.CardType_Cheng_Qing, protos.CardType_Jie_Huo, protos.CardType_Diao_Bao, protos.CardType_Wu_Dao}

func init() {
	game.AIMainPhase[protos.CardType_Wei_Bi] = weiBi
}

func weiBi(player interfaces.IPlayer, card interfaces.ICard) bool {
	chooseHuman := utils.Random.Intn(2) == 0
	var players []interfaces.IPlayer
	for _, p := range player.GetGame().GetPlayers() {
		if p.Location() != player.Location() && p.IsAlive() && len(p.GetCards()) > 0 {
			if _, ok := p.(*game.HumanPlayer); ok == chooseHuman {
				players = append(players, p)
			}
		}
	}
	if len(players) == 0 {
		for _, p := range player.GetGame().GetPlayers() {
			if p.Location() != player.Location() && p.IsAlive() && len(p.GetCards()) > 0 {
				players = append(players, p)
			}
		}
	}
	if len(players) == 0 {
		return false
	}
	p := players[utils.Random.Intn(len(players))]
	cardType := weiBiTypes[utils.Random.Intn(len(weiBiTypes))]
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, p, cardType) })
	})
	return true
}
