package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type ShiTan struct {
	interfaces.BaseCard
	WhoDrawCard []protos.Color
}

func (card *ShiTan) String() string {
	color := utils.CardColorToString(card.Color...)
	if len(card.WhoDrawCard) == 1 {
		return color + utils.IdentityColorToString(card.WhoDrawCard[0]) + "+1试探"
	}
	m := map[protos.Color]bool{protos.Color_Black: true, protos.Color_Red: true, protos.Color_Blue: true}
	delete(m, card.WhoDrawCard[0])
	delete(m, card.WhoDrawCard[1])
	for whoDiscardCard := range m {
		return color + utils.IdentityColorToString(whoDiscardCard) + "-1试探"
	}
	panic("unreachable code")
}

func (card *ShiTan) GetType() protos.CardType {
	return protos.CardType_Shi_Tan
}

func (card *ShiTan) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *ShiTan) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *ShiTan) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	pb.WhoDrawCard = append(pb.WhoDrawCard, card.WhoDrawCard...)
	return pb
}
