package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type ChengQing struct {
	interfaces.BaseCard
}

func (card *ChengQing) String() string {
	return utils.CardColorToString(card.Color...) + "澄清"
}

func (card *ChengQing) GetType() protos.CardType {
	return protos.CardType_Cheng_Qing
}

func (card *ChengQing) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *ChengQing) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *ChengQing) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
