package card

import (
	"github.com/CuteReimu/FengSheng/game"
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

func (card *PingHeng) CanUse(game interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	if game.GetCurrentPhase() != protos.Phase_Main_Phase || game.GetWhoseTurn() != r.Location() || !game.IsIdleTimePoint() {
		logger.Error("平衡的使用时机不对")
		return false
	}
	target := args[0].(interfaces.IPlayer)
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

func (card *PingHeng) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
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
	var discardCards, targetDiscardCards []interfaces.ICard
	for _, c := range r.GetCards() {
		discardCards = append(discardCards, c)
	}
	logger.Info(r, "弃掉了", discardCards)
	g.PlayerDiscardCard(r, discardCards...)
	for _, c := range target.GetCards() {
		targetDiscardCards = append(targetDiscardCards, c)
	}
	logger.Info(target, "弃掉了", targetDiscardCards)
	g.PlayerDiscardCard(target, targetDiscardCards...)
	r.Draw(3)
	target.Draw(3)
	g.GetDeck().Discard(card)
	game.Post(g.MainPhase)
}

func (card *PingHeng) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *PingHeng) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *PingHeng) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
