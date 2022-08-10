package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type LiYou struct {
	interfaces.BaseCard
}

func (card *LiYou) String() string {
	return utils.CardColorToString(card.Color...) + "利诱"
}

func (card *LiYou) GetType() protos.CardType {
	return protos.CardType_Li_You
}

func (card *LiYou) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *LiYou) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *LiYou) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
