package ai

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

func init() {
	game.AIMainPhase[protos.CardType_Wei_Bi] = weiBi
}

func weiBi(e *game.MainPhaseIdle, card game.ICard) bool {
	player := e.Player
	identity1, _ := player.GetIdentity()
	var players []game.IPlayer
	for _, p := range player.GetGame().GetPlayers() {
		if p.Location() != player.Location() && p.IsAlive() && len(p.GetCards()) > 0 {
			if identity2, _ := p.GetIdentity(); identity1 == protos.Color_Black || identity1 != identity2 {
				if func(p game.IPlayer) bool {
					for _, c := range p.GetCards() {
						switch c.GetType() {
						case protos.CardType_Cheng_Qing, protos.CardType_Jie_Huo, protos.CardType_Diao_Bao, protos.CardType_Wu_Dao:
							return true
						}
					}
					return false
				}(p) {
					players = append(players, p)
				}
			}
		}
	}
	if len(players) == 0 {
		return false
	}
	p := players[utils.Random.Intn(len(players))]
	var cardTypes []protos.CardType
	for _, c := range p.GetCards() {
		switch c.GetType() {
		case protos.CardType_Cheng_Qing, protos.CardType_Jie_Huo, protos.CardType_Diao_Bao, protos.CardType_Wu_Dao:
			cardTypes = append(cardTypes, c.GetType())
		}
	}
	cardType := cardTypes[utils.Random.Intn(len(cardTypes))]
	time.AfterFunc(2*time.Second, func() {
		game.Post(func() { card.Execute(player.GetGame(), player, p, cardType) })
	})
	return true
}
