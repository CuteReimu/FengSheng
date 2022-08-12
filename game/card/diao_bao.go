package card

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type DiaoBao struct {
	interfaces.BaseCard
}

func (card *DiaoBao) String() string {
	return utils.CardColorToString(card.Color...) + "调包"
}

func (card *DiaoBao) GetType() protos.CardType {
	return protos.CardType_Diao_Bao
}

func (card *DiaoBao) CanUse(g interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *DiaoBao) Execute(g interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *DiaoBao) CanUse2(g interfaces.IGame, user interfaces.IPlayer, args ...interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (card *DiaoBao) Execute2(g interfaces.IGame, user interfaces.IPlayer, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (card *DiaoBao) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
