package card

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/CuteReimu/FengSheng/utils"
)

type WuDao struct {
	game.BaseCard
}

func (card *WuDao) String() string {
	return utils.CardColorToString(card.Color...) + "误导"
}

func (card *WuDao) GetType() protos.CardType {
	return protos.CardType_Wu_Dao
}

func (card *WuDao) CanUse(g *game.Game, _ game.IPlayer, args ...interface{}) bool {
	target := args[0].(game.IPlayer)
	fsm, ok := g.GetFsm().(*game.FightPhaseIdle)
	if !ok {
		logger.Error("误导的使用时机不对")
		return false
	}
	var left, right int
	for left = fsm.InFrontOfWhom.Location() - 1; left != fsm.InFrontOfWhom.Location(); left-- {
		if left < 0 {
			left += len(g.GetPlayers())
		}
		if g.GetPlayers()[left].IsAlive() {
			break
		}
	}
	for right = fsm.InFrontOfWhom.Location() + 1; right != fsm.InFrontOfWhom.Location(); right++ {
		if right >= len(g.GetPlayers()) {
			right -= len(g.GetPlayers())
		}
		if g.GetPlayers()[right].IsAlive() {
			break
		}
	}
	if target.Location() == fsm.InFrontOfWhom.Location() || target.Location() != left && target.Location() != right {
		logger.Error("误导只能选择情报当前人左右两边的人作为目标")
		return false
	}
	return true
}

func (card *WuDao) Execute(g *game.Game, r game.IPlayer, args ...interface{}) {
	target := args[0].(game.IPlayer)
	logger.Info(r, "对", target, "使用了", card)
	fsm := g.GetFsm().(*game.FightPhaseIdle)
	r.DeleteCard(card.GetId())
	resolveFunc := func() (next game.Fsm, continueResolve bool) {
		fsm.InFrontOfWhom = target
		fsm.WhoseFightTurn = fsm.InFrontOfWhom
		g.GetDeck().Discard(card)
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
		return fsm, true
	}
	g.Resolve(&game.OnUseCard{
		WhoseTurn:   fsm.WhoseTurn,
		Player:      r,
		Card:        card,
		AskWhom:     r,
		ResolveFunc: resolveFunc,
	})
}

func (card *WuDao) ToPbCard() *protos.Card {
	pb := card.BaseCard.ToPbCard()
	pb.CardType = card.GetType()
	return pb
}
