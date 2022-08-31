package gm

import (
	"github.com/CuteReimu/FengSheng/card"
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"net/url"
	"strconv"
)

func init() {
	handlers["addcard"] = addCard
}

func addCard(values url.Values) []byte {
	playerId, err := strconv.Atoi(values.Get("player"))
	if err != nil {
		return []byte(`{"error": "invalid player"}`)
	}
	cardType, err := strconv.ParseInt(values.Get("card"), 10, 32)
	if err != nil {
		return []byte(`{"error": "invalid card"}`)
	}
	if _, ok := protos.CardType_name[int32(cardType)]; !ok {
		return []byte(`{"error": "invalid card"}`)
	}
	count, _ := strconv.Atoi(values.Get("count"))
	if count <= 0 {
		count = 1
	} else if count > 99 {
		count = 99
	}
	var availableCards []game.ICard
	for _, c := range game.DefaultDeck {
		if c.GetType() == protos.CardType(cardType) {
			availableCards = append(availableCards, c)
		}
	}
	for _, g := range game.Cache {
		game.Post(func() {
			if playerId < len(g.Players) && g.Players[playerId].IsAlive() {
				var cards []game.ICard
				for i := 0; i < count; i++ {
					var c game.ICard
					switch protos.CardType(cardType) {
					case protos.CardType_Cheng_Qing:
						c2 := &card.ChengQing{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.ChengQing).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Shi_Tan:
						c1 := availableCards[utils.Random.Intn(len(availableCards))].(*card.ShiTan)
						c2 := &card.ShiTan{BaseCard: c1.BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c2.WhoDrawCard = append(c2.WhoDrawCard, c1.WhoDrawCard...)
						c = c2
					case protos.CardType_Wei_Bi:
						c2 := &card.WeiBi{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.WeiBi).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Li_You:
						c2 := &card.LiYou{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.LiYou).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Ping_Heng:
						c2 := &card.PingHeng{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.PingHeng).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Po_Yi:
						c2 := &card.PoYi{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.PoYi).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Jie_Huo:
						c2 := &card.JieHuo{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.JieHuo).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Diao_Bao:
						c2 := &card.DiaoBao{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.DiaoBao).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					case protos.CardType_Wu_Dao:
						c2 := &card.WuDao{BaseCard: availableCards[utils.Random.Intn(len(availableCards))].(*card.WuDao).BaseCard}
						c2.BaseCard.Id = g.GetDeck().GetNextId()
						c = c2
					}
					cards = append(cards, c)
				}
				g.Players[playerId].AddCards(cards...)
				logger.Info("由于GM命令，", g.Players[playerId], "摸了", cards, "，现在有", len(g.Players[playerId].GetCards()), "张手牌")
				for _, player := range g.GetPlayers() {
					if player.Location() == playerId {
						player.NotifyAddHandCard(playerId, 0, cards...)
					} else {
						player.NotifyAddHandCard(playerId, len(cards))
					}
				}
			}
		})
	}
	return []byte(`{"msg": "success"}`)
}
