package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type DiaoBao struct {
	game.BaseCard
}

func (card *DiaoBao) String() string {
	return utils.CardColorToString(card.Color...) + "调包"
}

func (card *DiaoBao) GetType() protos.CardType {
	return protos.CardType_Diao_Bao
}

func (card *DiaoBao) CanUse(g *game.Game, r game.IPlayer, _ ...interface{}) bool {
	fsm, ok := g.GetFsm().(*game.FightPhaseIdle)
	if !ok || r.Location() != fsm.WhoseFightTurn.Location() {
		logger.Error("调包的使用时机不对")
		return false
	}
	return true
}

func (card *DiaoBao) Execute(g *game.Game, r game.IPlayer, _ ...interface{}) {
	fsm := g.GetFsm().(*game.FightPhaseIdle)
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		oldCard := fsm.MessageCard
		g.GetDeck().Discard(oldCard)
		fsm.MessageCard = card
		fsm.IsMessageCardFaceUp = false
		fsm.WhoseFightTurn = fsm.InFrontOfWhom
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
		return fsm, true
	}
	g.Resolve(&game.OnUseCard{
		WhoseTurn:   fsm.WhoseTurn,
		Player:      r,
		Card:        card,
		AskWhom:     r,
		ResolveFunc: resolveFunc,
	})
}

func (card *DiaoBao) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
