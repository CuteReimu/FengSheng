package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"google.golang.org/protobuf/proto"
	"time"
)

type ShiTan struct {
	game.BaseCard
	WhoDrawCard []protos.Color
}

func (card *ShiTan) String() string {
	color := utils.CardColorToString(card.Color...)
	if len(card.WhoDrawCard) == 1 {
		return color + utils.IdentityColorToString(card.WhoDrawCard[0]) + "+1试探"
	}
	m := map[protos.Color]bool{protos.Color_Black: true, protos.Color_Red: true, protos.Color_Blue: true}
	delete(m, card.WhoDrawCard[0])
	delete(m, card.WhoDrawCard[1])
	for whoDiscardCard := range m {
		return color + utils.IdentityColorToString(whoDiscardCard) + "-1试探"
	}
	panic("unreachable code")
}

func (card *ShiTan) GetType() protos.CardType {
	return protos.CardType_Shi_Tan
}

func (card *ShiTan) CanUse(g *game.Game, r game.IPlayer, args ...interface{}) bool {
	target := args[0].(game.IPlayer)
	fsm, ok := g.GetFsm().(*game.MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		logger.Error("试探的使用时机不对")
		return false
	}
	if r.Location() == target.Location() {
		logger.Error("试探不能对自己使用")
		return false
	}
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return false
	}
	return true
}

func (card *ShiTan) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	g.Resolve(&executeShiTan{
		player: r,
		target: target,
		card:   card,
	})
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			msg := &protos.UseShiTanToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
			}
			if p.Location() == r.Location() {
				msg.CardId = card.GetId()
			}
			player.Send(msg)
		}
	}
}

type executeShiTan struct {
	player game.IPlayer
	target game.IPlayer
	card   *ShiTan
}

func (e *executeShiTan) Resolve() (next game.Fsm, continueResolve bool) {
	r, target, card := e.player, e.target, e.card
	g := r.GetGame()
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			msg := &protos.ShowShiTanToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
				WaitingSecond:  20,
			}
			switch p.Location() {
			case target.Location():
				seq := player.Seq
				msg.Seq = player.Seq
				player.Timer = time.AfterFunc(time.Second*time.Duration(msg.WaitingSecond+2), func() {
					game.Post(func() {
						if player.Seq == seq {
							player.Seq++
							if player.Timer != nil {
								player.Timer.Stop()
							}
							e.autoSelect()
						}
					})
				})
				fallthrough
			case r.Location():
				msg.Card = card.ToPbCard()
			}
			player.Send(msg)
		}
	}
	if _, ok := target.(*game.RobotPlayer); ok {
		time.AfterFunc(4*time.Second, func() {
			game.Post(func() {
				e.autoSelect()
				g.Resolve(&game.MainPhaseIdle{Player: r})
			})
		})
	}
	return e, false
}

func (e *executeShiTan) ResolveProtocol(player game.IPlayer, pb proto.Message) (next game.Fsm, continueResolve bool) {
	msg, ok := pb.(*protos.ExecuteShiTanTos)
	if !ok {
		logger.Error("现在正在结算试探", e.card)
		return e, false
	}
	if e.target.Location() != player.Location() {
		logger.Error("你不是试探的目标", e.card)
		return e, false
	}
	r, card := e.target, e.card
	cardIds := msg.CardId
	g := r.GetGame()
	if card.checkDrawCard(r) || len(r.GetCards()) == 0 {
		if len(cardIds) != 0 {
			logger.Error(r, "被使用", card, "时不应该弃牌")
			return e, false
		}
	} else {
		if len(cardIds) != 1 {
			logger.Error(r, "被使用", card, "时应该弃一张牌")
			return e, false
		}
		if r.FindCard(cardIds[0]) == nil {
			logger.Error("没有这张牌")
			return e, false
		}
	}
	player.IncrSeq()
	if card.checkDrawCard(r) {
		logger.Info(r, "选择了[摸一张牌]")
		card.notifyResult(g, r, true)
		r.Draw(1)
	} else {
		logger.Info(r, "选择了[弃一张牌]")
		card.notifyResult(g, r, false)
		if len(cardIds) > 0 {
			g.PlayerDiscardCard(r, r.FindCard(cardIds[0]))
		}
	}
	return &game.MainPhaseIdle{Player: e.player}, true
}

func (card *ShiTan) notifyResult(g *game.Game, target game.IPlayer, draw bool) {
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			p.Send(&protos.ExecuteShiTanToc{
				PlayerId:   p.GetAlternativeLocation(target.Location()),
				IsDrawCard: draw,
			})
		}
	}
}

func (card *ShiTan) checkDrawCard(target game.IPlayer) bool {
	identity, _ := target.GetIdentity()
	for _, i := range card.WhoDrawCard {
		if i == identity {
			return true
		}
	}
	return false
}

func (e *executeShiTan) autoSelect() {
	card, target := e.card, e.target
	var discardCardIds []uint32
	if !card.checkDrawCard(target) && len(target.GetCards()) > 0 {
		for cardId := range target.GetCards() {
			discardCardIds = append(discardCardIds, cardId)
			break
		}
	}
	e.ResolveProtocol(e.target, &protos.ExecuteShiTanTos{CardId: discardCardIds})
}

func (card *ShiTan) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	pb.WhoDrawCard = append(pb.WhoDrawCard, card.WhoDrawCard...)
	return pb
}
