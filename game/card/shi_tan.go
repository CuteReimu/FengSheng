package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
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
	if len(args) != 1 {
		logger.Error("参数错误")
		return false
	}
	target, ok := args[0].(interfaces.IPlayer)
	if !ok {
		logger.Error("参数错误")
		return false
	}
	if game.GetCurrentPhase() != protos.Phase_Main_Phase || r.GetGame().GetWhoseTurn() != r.Location() {
		logger.Error("试探的使用时机不对")
		return false
	}
	if r.Location() == target.Location() {
		logger.Error("试探不能对自己使用")
		return false
	}
	return true
}

func (card *ShiTan) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.UseShiTanToc{
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
			})
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
				msg.Seq = player.Seq
				fallthrough
			case r.Location():
				msg.Card = card.ToPbCard()
			}
			player.Send(msg)
		}
	}
}

func (card *ShiTan) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	pb.WhoDrawCard = append(pb.WhoDrawCard, card.WhoDrawCard...)
	return pb
}
