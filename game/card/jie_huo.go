package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type JieHuo struct {
	interfaces.BaseCard
}

func (card *JieHuo) String() string {
	return utils.CardColorToString(card.Color...) + "截获"
}

func (card *JieHuo) GetType() protos.CardType {
	return protos.CardType_Jie_Huo
}

func (card *JieHuo) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *JieHuo) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *JieHuo) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
