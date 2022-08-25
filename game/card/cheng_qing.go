package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type ChengQing struct {
	interfaces.BaseCard
}

func (card *ChengQing) String() string {
	return utils.CardColorToString(card.Color...) + "澄清"
}

func (card *ChengQing) GetType() protos.CardType {
	return protos.CardType_Cheng_Qing
}

func (card *ChengQing) CanUse(game interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	target := args[0].(interfaces.IPlayer)
	targetCardId := args[1].(uint32)
	if game.GetDieState() != interfaces.DieStateWaitForChengQing && (game.GetCurrentPhase() != protos.Phase_Main_Phase || game.GetWhoseTurn() != r.Location() || !game.IsIdleTimePoint()) {
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
	if !utils.IsColorIn(protos.Color_Black, targetCard.GetColor()) {
		logger.Error("澄清只能对黑情报使用")
		return false
	}
	return true
}

func (card *ChengQing) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	targetCardId := args[1].(uint32)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
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
	if g.GetDieState() == interfaces.DieStateWaitForChengQing {
		g.AfterChengQing()
	} else {
		game.Post(g.MainPhase)
	}
}

func (card *ChengQing) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *ChengQing) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *ChengQing) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
