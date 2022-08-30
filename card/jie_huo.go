package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type JieHuo struct {
	game.BaseCard
}

func (card *JieHuo) String() string {
	return utils.CardColorToString(card.Color...) + "截获"
}

func (card *JieHuo) GetType() protos.CardType {
	return protos.CardType_Jie_Huo
}

func (card *JieHuo) CanUse(g *game.Game, r game.IPlayer, _ ...interface{}) bool {
	fsm, ok := g.GetFsm().(*game.FightPhaseIdle)
	if !ok || r.Location() != fsm.WhoseFightTurn.Location() {
		logger.Error("截获的使用时机不对")
		return false
	}
	if r.Location() == fsm.InFrontOfWhom.Location() {
		logger.Error("情报在自己面前不能使用截获")
		return false
	}
	return true
}

func (card *JieHuo) Execute(g *game.Game, r game.IPlayer, _ ...interface{}) {
	fsm := g.GetFsm().(*game.FightPhaseIdle)
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	fsm.InFrontOfWhom = r
	fsm.WhoseFightTurn = fsm.InFrontOfWhom
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
	g.ContinueResolve()
}

func (card *JieHuo) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
