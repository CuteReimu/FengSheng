package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type ChengQing struct {
	game.BaseCard
}

func (card *ChengQing) String() string {
	return utils.CardColorToString(card.Color...) + "澄清"
}

func (card *ChengQing) GetType() protos.CardType {
	return protos.CardType_Cheng_Qing
}

func (card *ChengQing) CanUse(g *game.Game, r game.IPlayer, args ...interface{}) bool {
	target := args[0].(game.IPlayer)
	targetCardId := args[1].(uint32)
	switch fsm := g.GetFsm().(type) {
	case *game.MainPhaseIdle:
		if r.Location() != fsm.Player.Location() {
			logger.Error("澄清的使用时机不对")
			return false
		}
	case *game.WaitForChengQing:
		if r.Location() != fsm.AskWhom.Location() {
			logger.Error("澄清的使用时机不对")
			return false
		}
	default:
		logger.Error("澄清的使用时机不对")
		return false
	}
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	targetCard := target.FindMessageCard(targetCardId)
	if targetCard == nil {
		logger.Error("没有这张情报")
		return false
	}
	if !utils.IsColorIn(protos.Color_Black, targetCard.GetColors()) {
		logger.Error("澄清只能对黑情报使用")
		return false
	}
	return true
}

func (card *ChengQing) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	targetCardId := args[1].(uint32)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	curFsm := g.GetFsm()
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		targetCard := target.FindMessageCard(targetCardId)
		logger.Info(target, "面前的", targetCard, "被置入弃牌堆")
		target.DeleteMessageCard(targetCardId)
		g.GetDeck().Discard(targetCard)
		for _, p := range g.GetPlayers() {
			if player, ok := p.(*game.HumanPlayer); ok {
				msg := &protos.UseChengQingToc{
					Card:           card.ToPbCard(),
					PlayerId:       p.GetAlternativeLocation(r.Location()),
					TargetPlayerId: p.GetAlternativeLocation(target.Location()),
					TargetCardId:   targetCardId,
				}
				player.Send(msg)
			}
		}
		g.GetDeck().Discard(card)
		return curFsm, true
	}
	switch fsm := curFsm.(type) {
	case *game.MainPhaseIdle:
		g.Resolve(&game.OnUseCard{
			WhoseTurn:   fsm.Player,
			Player:      r,
			Card:        card,
			AskWhom:     r,
			ResolveFunc: resolveFunc,
		})
	case *game.WaitForChengQing:
		g.Resolve(&game.OnUseCard{
			WhoseTurn:   fsm.WhoseTurn,
			Player:      r,
			Card:        card,
			AskWhom:     r,
			ResolveFunc: resolveFunc,
		})
	}
}

func (card *ChengQing) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
