package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type PingHeng struct {
	game.BaseCard
}

func (card *PingHeng) String() string {
	return utils.CardColorToString(card.Color...) + "平衡"
}

func (card *PingHeng) GetType() protos.CardType {
	return protos.CardType_Ping_Heng
}

func (card *PingHeng) CanUse(g *game.Game, r game.IPlayer, args ...interface{}) bool {
	fsm, ok := g.GetFsm().(*game.MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		logger.Error("平衡的使用时机不对")
		return false
	}
	target := args[0].(game.IPlayer)
	if target.Location() == r.Location() {
		logger.Error("平衡不能对自己使用")
		return false
	}
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	return true
}

func (card *PingHeng) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		for _, p := range g.GetPlayers() {
			if player, ok := p.(*game.HumanPlayer); ok {
				msg := &protos.UsePingHengToc{
					PlayerId:       p.GetAlternativeLocation(r.Location()),
					TargetPlayerId: p.GetAlternativeLocation(target.Location()),
					PingHengCard:   card.ToPbCard(),
				}
				player.Send(msg)
			}
		}
		var discardCards, targetDiscardCards []game.ICard
		for _, c := range r.GetCards() {
			discardCards = append(discardCards, c)
		}
		g.PlayerDiscardCard(r, discardCards...)
		for _, c := range target.GetCards() {
			targetDiscardCards = append(targetDiscardCards, c)
		}
		g.PlayerDiscardCard(target, targetDiscardCards...)
		r.Draw(3)
		target.Draw(3)
		g.GetDeck().Discard(card)
		return &game.MainPhaseIdle{Player: r}, true
	}
	g.Resolve(&game.OnUseCard{
		WhoseTurn:   r,
		Player:      r,
		Card:        card,
		AskWhom:     r,
		ResolveFunc: resolveFunc,
	})
}

func (card *PingHeng) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
