package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"google.golang.org/protobuf/proto"
	"time"
)

type PoYi struct {
	game.BaseCard
}

func (card *PoYi) String() string {
	return utils.CardColorToString(card.Color...) + "破译"
}

func (card *PoYi) GetType() protos.CardType {
	return protos.CardType_Po_Yi
}

func (card *PoYi) CanUse(g *game.Game, r game.IPlayer, _ ...interface{}) bool {
	fsm, ok := g.GetFsm().(*game.SendPhaseIdle)
	if !ok || r.Location() != fsm.InFrontOfWhom.Location() {
		logger.Error("破译的使用时机不对")
		return false
	}
	if fsm.IsMessageCardFaceUp {
		logger.Error("破译不能对已翻开的情报使用")
		return false
	}
	return true
}

func (card *PoYi) Execute(g *game.Game, r game.IPlayer, _ ...interface{}) {
	fsm := g.GetFsm().(*game.SendPhaseIdle)
	logger.Info(r, "使用了", card)
	r.DeleteCard(card.GetId())
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		return &executePoYi{
			card:      card,
			sendPhase: fsm,
		}, true
	}
	g.Resolve(&game.OnUseCard{
		WhoseTurn:   r,
		Player:      r,
		Card:        card,
		AskWhom:     r,
		ResolveFunc: resolveFunc,
	})
}

type executePoYi struct {
	card      *PoYi
	sendPhase *game.SendPhaseIdle
}

func (e *executePoYi) Resolve() (next game.Fsm, continueResolve bool) {
	r, card := e.sendPhase.InFrontOfWhom, e.card
	g := r.GetGame()
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.UsePoYiToc{
				Card:     card.ToPbCard(),
				PlayerId: p.GetAlternativeLocation(r.Location()),
			}
			if p.Location() == r.Location() {
				msg.MessageCard = e.sendPhase.MessageCard.ToPbCard()
				msg.WaitingSecond = 20
				msg.Seq = p.Seq
			}
			p.Send(msg)
		}
	}
	if _, ok := r.(*game.RobotPlayer); ok {
		time.AfterFunc(2*time.Second, func() {
			game.Post(func() {
				e.showAndDrawCard(utils.IsColorIn(protos.Color_Black, e.sendPhase.MessageCard.GetColors()))
				g.Resolve(e.sendPhase)
			})
		})
	}
	return e, false
}

func (e *executePoYi) ResolveProtocol(player game.IPlayer, pb proto.Message) (next game.Fsm, continueResolve bool) {
	msg, ok := pb.(*protos.PoYiShowTos)
	if !ok {
		logger.Error("现在正在结算破译")
		return e, false
	}
	if player.Location() != e.sendPhase.InFrontOfWhom.Location() {
		logger.Error("你不是破译的使用者")
		return e, false
	}
	if msg.Show && !utils.IsColorIn(protos.Color_Black, e.sendPhase.MessageCard.GetColors()) {
		logger.Error("非黑牌不能翻开")
		return e, false
	}
	player.IncrSeq()
	e.showAndDrawCard(msg.Show)
	return e.sendPhase, true
}

func (card *PoYi) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}

func (e *executePoYi) showAndDrawCard(show bool) {
	card, r := e.card, e.sendPhase.InFrontOfWhom
	if show {
		logger.Info(e.sendPhase.MessageCard, "被翻开了")
		e.sendPhase.IsMessageCardFaceUp = true
		r.Draw(1)
	}
	for _, player := range r.GetGame().GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.PoYiShowToc{
				PlayerId: p.GetAlternativeLocation(r.Location()),
				Show:     show,
			}
			if show {
				msg.MessageCard = e.sendPhase.MessageCard.ToPbCard()
			}
			p.Send(msg)
		}
	}
	r.GetGame().GetDeck().Discard(card)
}
