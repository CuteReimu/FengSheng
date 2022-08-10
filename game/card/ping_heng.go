package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type PingHeng struct {
	interfaces.BaseCard
}

func (card *PingHeng) String() string {
	return utils.CardColorToString(card.Color...) + "平衡"
}

func (card *PingHeng) GetType() protos.CardType {
	return protos.CardType_Ping_Heng
}

func (card *PingHeng) CanUse(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *PingHeng) Execute(b interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *PingHeng) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
