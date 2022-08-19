package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type WuDao struct {
	interfaces.BaseCard
}

func (card *WuDao) String() string {
	return utils.CardColorToString(card.Color...) + "误导"
}

func (card *WuDao) GetType() protos.CardType {
	return protos.CardType_Wu_Dao
}

func (card *WuDao) CanUse(game interfaces.IGame, _ interfaces.IPlayer, args ...interface{}) bool {
	target := args[0].(interfaces.IPlayer)
	if game.GetCurrentPhase() != protos.Phase_Fight_Phase {
		logger.Error("误导的使用时机不对")
		return false
	}
	if target.Location() != (game.GetWhoseSendTurn()+1)%len(game.GetPlayers()) && game.GetWhoseSendTurn() != (target.Location()+1)%len(game.GetPlayers()) {
		logger.Error("误导只能选择情报当前人左右两边的人作为目标")
		return false
	}
	return true
}

func (card *WuDao) Execute(g interfaces.IGame, r interfaces.IPlayer, args ...interface{}) {
	target := args[0].(interfaces.IPlayer)
	g.SetWhoseSendTurn(target.Location())
	for _, player := range g.GetPlayers() {
		if p, ok := player.(*game.HumanPlayer); ok {
			msg := &protos.UseWuDaoToc{
				Card:           card.ToPbCard(),
				PlayerId:       p.GetAlternativeLocation(r.Location()),
				TargetPlayerId: p.GetAlternativeLocation(target.Location()),
			}
			p.Send(msg)
		}
	}
	for _, p := range g.GetPlayers() {
		p.NotifyFightPhase(20)
	}
}

func (card *WuDao) CanUse2(interfaces.IGame, interfaces.IPlayer, ...interface{}) bool {
	panic("unreachable code")
}

func (card *WuDao) Execute2(interfaces.IGame, interfaces.IPlayer, ...interface{}) {
	panic("unreachable code")
}

func (card *WuDao) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
