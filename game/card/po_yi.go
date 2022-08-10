package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type PoYi struct {
	interfaces.BaseCard
}

func (card *PoYi) String() string {
	return utils.CardColorToString(card.Color...) + "破译"
}

func (card *PoYi) GetType() protos.CardType {
	return protos.CardType_Po_Yi
}

func (card *PoYi) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *PoYi) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *PoYi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
