package role

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"google.golang.org/protobuf/proto"
)

func init() {
	game.RoleCache = append(game.RoleCache, &game.RoleSkillsData{
		Name:   "邵秀",
		Role:   protos.Role_shao_xiu,
		FaceUp: true,
		Skills: []game.ISkill{&MianLiCangZhen{}},
	})
}

type MianLiCangZhen struct {
}

func (m *MianLiCangZhen) Init(g *game.Game) {
	g.AddListeningSkill(m)
}

func (m *MianLiCangZhen) GetSkillId() game.SkillId {
	return game.SkillIdMianLiCangZhen
}

func (m *MianLiCangZhen) Execute(g *game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
	fsm, ok := g.GetFsm().(*game.ReceivePhaseSenderSkill)
	if !ok || fsm.WhoseTurn.FindSkill(m.GetSkillId()) == nil {
		return nil, false, false
	}
	if fsm.WhoseTurn.GetSkillUseCount(m.GetSkillId()) > 0 {
		return nil, false, false
	}
	fsm.WhoseTurn.AddSkillUseCount(m.GetSkillId())
	return &executeMianLiCangZhen{fsm: fsm}, true, ok
}

func (m *MianLiCangZhen) ExecuteProtocol(*game.Game, game.IPlayer, proto.Message) {
}

type executeMianLiCangZhen struct {
	fsm *game.ReceivePhaseSenderSkill
}

func (e *executeMianLiCangZhen) Resolve() (next game.Fsm, continueResolve bool) {
	for _, p := range e.fsm.WhoseTurn.GetGame().GetPlayers() {
		p.NotifyReceivePhaseWithWaiting(e.fsm.WhoseTurn, e.fsm.InFrontOfWhom, e.fsm.MessageCard, e.fsm.WhoseTurn, 20)
	}
	return e, false
}

func (e *executeMianLiCangZhen) ResolveProtocol(player game.IPlayer, message proto.Message) (next game.Fsm, continueResolve bool) {
	if player.Location() != e.fsm.WhoseTurn.Location() {
		logger.Error("不是你发技能的时机")
		return e, false
	}
	if _, ok := message.(*protos.EndReceivePhaseTos); ok && player.Location() == e.fsm.WhoseTurn.Location() {
		player.IncrSeq()
		return e.fsm, true
	}
	r := e.fsm.WhoseTurn
	g := r.GetGame()
	pb := message.(*protos.SkillMianLiCangZhenTos)
	if humanPlayer, ok := r.(*game.HumanPlayer); ok && pb.Seq != humanPlayer.Seq {
		logger.Error("操作太晚了, required Seq: ", humanPlayer.Seq, ", actual Seq: ", pb.Seq)
		return e, false
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		logger.Error("没有这张卡")
		return e, false
	}
	if pb.TargetPlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		logger.Error("目标错误")
		return e, false
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.TargetPlayerId))]
	if !target.IsAlive() {
		logger.Error("目标已死亡")
		return e, false
	}
	r.IncrSeq()
	logger.Info("[邵秀]发动了[绵里藏针]")
	r.DeleteCard(card.GetId())
	target.AddMessageCards(card)
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.SkillMianLiCangZhenToc{
				PlayerId:       player.GetAlternativeLocation(r.Location()),
				Card:           card.ToPbCard(),
				TargetPlayerId: player.GetAlternativeLocation(target.Location()),
			})
		}
	}
	r.Draw(1)
	return e.fsm, true
}
