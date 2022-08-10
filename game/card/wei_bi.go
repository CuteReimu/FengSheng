package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type WeiBi struct {
	interfaces.BaseCard
}

func (card *WeiBi) String() string {
	return utils.CardColorToString(card.Color...) + "威逼"
}

func (card *WeiBi) GetType() protos.CardType {
	return protos.CardType_Wei_Bi
}

func (card *WeiBi) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *WeiBi) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *WeiBi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
