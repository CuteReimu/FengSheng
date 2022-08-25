package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

type PoYi struct {
	interfaces.BaseCard
}

func (card *PoYi) String() string {
	return utils.CardColorToString(card.Color...) + "破译"
}

func (card *PoYi) GetType() protos.CardType {
	return protos.CardType_Po_Yi
}

func (card *PoYi) CanUse(game interfaces.IGame, r interfaces.IPlayer, _ ...interface{}) bool {
	if game.GetCurrentPhase() != protos.Phase_Send_Phase || game.GetWhoseSendTurn() != r.Location() {
		logger.Error("破译的使用时机不对")
		return false
	}
	return true
}

func (card *PoYi) Execute(g interfaces.IGame, r interfaces.IPlayer, _ ...interface{}) {
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	g.SetCurrentCard(&interfaces.CurrentCard{Card: card, Player: r.Location(), TargetPlayer: r.Location()})
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.UsePoYiToc{
				Card:     card.ToPbCard(),
				PlayerId: p.GetAlternativeLocation(r.Location()),
			}
			if p.Location() == r.Location() {
				msg.MessageCard = g.GetCurrentMessageCard().ToPbCard()
				msg.WaitingSecond = 20
				msg.Seq = p.Seq
			}
			p.Send(msg)
		}
	}
	if _, ok := r.(*game.RobotPlayer); ok {
		time.AfterFunc(time.Second, func() {
			game.Post(func() {
				card.showAndDrawCard(g, r, utils.IsColorIn(protos.Color_Black, g.GetCurrentMessageCard().GetColor()))
			})
		})
	}
}

func (card *PoYi) CanUse2(g interfaces.IGame, _ interfaces.IPlayer, args ...interface{}) bool {
	show := args[0].(bool)
	if show && !utils.IsColorIn(protos.Color_Black, g.GetCurrentMessageCard().GetColor()) {
		logger.Error("非黑牌不能翻开：", show)
		return false
	}
	return true
}

func (card *PoYi) Execute2(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	show := args[0].(bool)
	card.showAndDrawCard(g, r, show)
}

func (card *PoYi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}

func (card *PoYi) showAndDrawCard(g interfaces.IGame, r interfaces.IPlayer, show bool) {
	if show {
		logger.Info(g.GetCurrentMessageCard(), "被翻开了")
		g.SetMessageCardFaceUp(true)
		r.Draw(1)
	}
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.PoYiShowToc{
				PlayerId: p.GetAlternativeLocation(r.Location()),
				Show:     show,
			}
			if show {
				msg.MessageCard = g.GetCurrentMessageCard().ToPbCard()
			}
			p.Send(msg)
		}
	}
	g.SetCurrentCard(nil)
	g.GetDeck().Discard(card)
	for _, p := range g.GetPlayers() {
		p.NotifySendPhase(20)
	}
}
