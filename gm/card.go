package gm

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/card"
	"github.com/CuteReimu/FengSheng/game/interfaces"
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
	for _, g := range game.Cache {
		game.Post(func() {
			if playerId < len(g.Players) && g.Players[playerId].IsAlive() {
				var cards []interfaces.ICard
				for i := 0; i < count; i++ {
					bc := interfaces.BaseCard{
						Id:        g.GetDeck().GetNextId(),
						Direction: protos.Direction(utils.Random.Intn(len(protos.Direction_name))),
						Lockable:  utils.Random.Intn(2) == 1,
					}
					color := utils.Random.Intn(5) + 1
					bc.Color = append(bc.Color, protos.Color(color%3))
					if color > 3 {
						bc.Color = append(bc.Color, protos.Color(color-3))
					}
					var c interfaces.ICard
					switch protos.CardType(cardType) {
					case protos.CardType_Cheng_Qing:
						c = &card.ChengQing{BaseCard: bc}
					case protos.CardType_Shi_Tan:
						shiTan := &card.ShiTan{BaseCard: bc}
						color = utils.Random.Intn(8)
						if color&0x01 != 0 {
							shiTan.WhoDrawCard = append(shiTan.WhoDrawCard, protos.Color_Red)
						}
						if color&0x02 != 0 {
							shiTan.WhoDrawCard = append(shiTan.WhoDrawCard, protos.Color_Blue)
						}
						if color&0x04 != 0 && len(shiTan.WhoDrawCard) < 2 {
							shiTan.WhoDrawCard = append(shiTan.WhoDrawCard, protos.Color_Black)
						}
						c = shiTan
					case protos.CardType_Wei_Bi:
						c = &card.WeiBi{BaseCard: bc}
					case protos.CardType_Li_You:
						c = &card.LiYou{BaseCard: bc}
					case protos.CardType_Ping_Heng:
						c = &card.PingHeng{BaseCard: bc}
					case protos.CardType_Po_Yi:
						c = &card.PoYi{BaseCard: bc}
					case protos.CardType_Jie_Huo:
						c = &card.JieHuo{BaseCard: bc}
					case protos.CardType_Diao_Bao:
						c = &card.DiaoBao{BaseCard: bc}
					case protos.CardType_Wu_Dao:
						c = &card.WuDao{BaseCard: bc}
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
