package role

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
)

func init() {
	game.RoleCache = append(game.RoleCache, &game.RoleSkillsData{
		Name:   "老鳖",
		Role:   protos.Role_lao_bie,
		FaceUp: true,
		Skills: []game.ISkill{&LianLuo{}, &MingEr{}},
	})
}

type LianLuo struct {
}

func (l *LianLuo) Init(*game.Game) {
}

func (l *LianLuo) GetSkillId() game.SkillId {
	return game.SkillIdLianLuo
}

func (l *LianLuo) Execute(*game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
	return
}

type MingEr struct {
}

func (m *MingEr) Init(g *game.Game) {
	g.AddListeningSkill(m)
}

func (m *MingEr) GetSkillId() game.SkillId {
	return game.SkillIdMingEr
}

func (m *MingEr) Execute(g *game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
	fsm, ok := g.GetFsm().(*game.ReceivePhaseSenderSkill)
	if !ok || !fsm.WhoseTurn.HasSkill(m.GetSkillId()) || !fsm.WhoseTurn.IsAlive() {
		return nil, false, false
	}
	if fsm.WhoseTurn.GetSkillUseCount(m.GetSkillId()) > 0 {
		return nil, false, false
	}
	logger.Info("[老鳖]发动了[明饵]")
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.SkillMingErToc{
				PlayerId: player.GetAlternativeLocation(fsm.WhoseTurn.Location()),
			})
		}
	}
	if fsm.WhoseTurn == fsm.InFrontOfWhom {
		fsm.WhoseTurn.Draw(2)
	} else {
		fsm.WhoseTurn.Draw(1)
		fsm.InFrontOfWhom.Draw(1)
	}
	return nil, false, false
}
