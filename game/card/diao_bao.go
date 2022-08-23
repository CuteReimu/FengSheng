package card

import (
	"github.com/CuteReimu/FengSheng/game"
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

func (card *DiaoBao) CanUse(game interfaces.IGame, _ interfaces.IPlayer, _ ...interface{}) bool {
	if game.GetCurrentPhase() != protos.Phase_Fight_Phase {
		logger.Error("调包的使用时机不对")
		return false
	}
	return true
}

func (card *DiaoBao) Execute(g interfaces.IGame, r interfaces.IPlayer, _ ...interface{}) {
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	oldCard := g.GetCurrentMessageCard()
	g.GetDeck().Discard(oldCard)
	g.SetCurrentMessageCard(card)
	g.SetWhoseFightTurn(g.GetWhoseSendTurn())
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.UseDiaoBaoToc{
				OldMessageCard: oldCard.ToPbCard(),
				PlayerId:       p.GetAlternativeLocation(r.Location()),
			}
			if p.Location() == r.Location() {
				msg.CardId = card.GetId()
			}
			p.Send(msg)
		}
	}
	for _, p := range g.GetPlayers() {
		p.NotifyFightPhase(20)
	}
}

func (card *DiaoBao) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *DiaoBao) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *DiaoBao) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
