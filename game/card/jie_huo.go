package card

import (
	"github.com/CuteReimu/FengSheng/game"
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

func (card *JieHuo) CanUse(game interfaces.IGame, r interfaces.IPlayer, _ ...interface{}) bool {
	if game.GetCurrentPhase() != protos.Phase_Fight_Phase {
		logger.Error("截获的使用时机不对")
		return false
	}
	if game.GetWhoseSendTurn() == r.Location() {
		logger.Error("情报在自己面前不能使用截获")
		return false
	}
	return true
}

func (card *JieHuo) Execute(g interfaces.IGame, r interfaces.IPlayer, _ ...interface{}) {
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	g.SetWhoseSendTurn(r.Location())
	g.SetWhoseFightTurn(g.GetWhoseSendTurn())
	g.GetDeck().Discard(card)
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.UseJieHuoToc{
				Card:     card.ToPbCard(),
				PlayerId: p.GetAlternativeLocation(r.Location()),
			}
			p.Send(msg)
		}
	}
	for _, p := range g.GetPlayers() {
		p.NotifyFightPhase(20)
	}
}

func (card *JieHuo) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *JieHuo) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *JieHuo) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
