package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type LiYou struct {
	interfaces.BaseCard
}

func (card *LiYou) String() string {
	return utils.CardColorToString(card.Color...) + "利诱"
}

func (card *LiYou) GetType() protos.CardType {
	return protos.CardType_Li_You
}

func (card *LiYou) CanUse(game interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	if game.GetCurrentPhase() != protos.Phase_Main_Phase || game.GetWhoseTurn() != r.Location() || game.IsIdleTimePoint() {
		logger.Error("利诱的使用时机不对")
		return false
	}
	target := args[0].(interfaces.IPlayer)
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	return true
}

func (card *LiYou) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	deckCards := g.GetDeck().Draw(1)
	var joinIntoHand bool
	if len(deckCards) > 0 {
		target.AddMessageCards(deckCards...)
		if target.CheckThreeSameMessageCard(deckCards[0].GetColor()...) {
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
	g.Post(g.MainPhase)
}

func (card *LiYou) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *LiYou) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *LiYou) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
