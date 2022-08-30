package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type LiYou struct {
	game.BaseCard
}

func (card *LiYou) String() string {
	return utils.CardColorToString(card.Color...) + "利诱"
}

func (card *LiYou) GetType() protos.CardType {
	return protos.CardType_Li_You
}

func (card *LiYou) CanUse(g *game.Game, r game.IPlayer, args ...interface{}) bool {
	fsm, ok := g.GetFsm().(*game.MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		logger.Error("利诱的使用时机不对")
		return false
	}
	target := args[0].(game.IPlayer)
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	return true
}

func (card *LiYou) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	deckCards := g.GetDeck().Draw(1)
	var joinIntoHand bool
	if len(deckCards) > 0 {
		target.AddMessageCards(deckCards...)
		if target.CheckThreeSameMessageCard(deckCards[0].GetColors()...) {
			target.DeleteMessageCard(deckCards[0].GetId())
			joinIntoHand = true
			r.AddCards(deckCards...)
			logger.Info(deckCards, "加入了", r, "的手牌")
		} else {
			logger.Info(deckCards, "加入了", target, "的情报区")
		}
	}
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			msg := &protos.UseLiYouToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
				LiYouCard:      card.ToPbCard(),
				JoinIntoHand:   joinIntoHand,
			}
			if len(deckCards) > 0 {
				msg.MessageCard = deckCards[0].ToPbCard()
			}
			player.Send(msg)
		}
	}
	g.GetDeck().Discard(card)
	g.Resolve(&game.MainPhaseIdle{Player: r})
}

func (card *LiYou) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
