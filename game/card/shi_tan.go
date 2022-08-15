package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
	"time"
)

type ShiTan struct {
	interfaces.BaseCard
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

func (card *ShiTan) CanUse(game interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	target := args[0].(interfaces.IPlayer)
	if game.GetCurrentPhase() != protos.Phase_Main_Phase || game.GetWhoseTurn() != r.Location() || game.IsIdleTimePoint() {
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

func (card *ShiTan) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	r.DeleteCard(card.GetId())
	g.SetCurrentCard(card)
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
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			msg := &protos.ShowShiTanToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
				WaitingSecond:  10,
			}
			switch p.Location() {
			case target.Location():
				player.Timer = time.AfterFunc(time.Second*time.Duration(msg.WaitingSecond), func() {
					g.Post(func() { card.autoSelect(g, player) })
				})
				msg.Seq = player.Seq
				fallthrough
			case r.Location():
				msg.Card = card.ToPbCard()
			}
			player.Send(msg)
		}
	}
	if _, ok := target.(*game.RobotPlayer); ok {
		time.AfterFunc(time.Second, func() {
			g.Post(func() { card.autoSelect(g, target) })
		})
	}
}

func (card *ShiTan) CanUse2(_ interfaces.IGame, r interfaces.IPlayer, args ...interface{}) bool {
	cardIds := args[0].([]uint32)
	if card.checkDrawCard(r) || len(r.GetCards()) == 0 {
		if len(cardIds) != 0 {
			logger.Error(r, "被使用", card, "时不应该弃牌")
			return false
		}
	} else {
		if len(cardIds) != 1 {
			logger.Error(r, "被使用", card, "时应该弃一张牌")
			return false
		}
		if r.FindCard(cardIds[0]) == nil {
			logger.Error("没有这张牌")
			return false
		}
	}
	return true
}

func (card *ShiTan) Execute2(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	cardIds := args[0].([]uint32)
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
	g.SetCurrentCard(nil)
	g.Post(g.MainPhase)
}

func (card *ShiTan) notifyResult(g interfaces.IGame, target interfaces.IPlayer, draw bool) {
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			p.Send(&protos.ExecuteShiTanToc{
				PlayerId:   p.GetAlternativeLocation(target.Location()),
				IsDrawCard: draw,
			})
		}
	}
}

func (card *ShiTan) checkDrawCard(target interfaces.IPlayer) bool {
	identity, _ := target.GetIdentity()
	for _, i := range card.WhoDrawCard {
		if i == identity {
			return true
		}
	}
	return false
}

func (card *ShiTan) autoSelect(g interfaces.IGame, target interfaces.IPlayer) {
	var discardCardIds []uint32
	if !card.checkDrawCard(target) && len(target.GetCards()) > 0 {
		for cardId := range target.GetCards() {
			discardCardIds = append(discardCardIds, cardId)
			break
		}
	}
	card.Execute2(g, target, discardCardIds)
}

func (card *ShiTan) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	pb.WhoDrawCard = append(pb.WhoDrawCard, card.WhoDrawCard...)
	return pb
}
