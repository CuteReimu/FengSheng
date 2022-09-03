package role

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"google.golang.org/protobuf/proto"
)

func init() {
	game.RoleCache = append(game.RoleCache, &game.RoleSkillsData{
		Name:   "毛不拔",
		Role:   protos.Role_mao_bu_ba,
		FaceUp: true,
		Skills: []game.ISkill{&QiHuoKeJu{}},
	})
}

type QiHuoKeJu struct {
}

func (m *QiHuoKeJu) Init(g *game.Game) {
	g.AddListeningSkill(m)
}

func (m *QiHuoKeJu) GetSkillId() game.SkillId {
	return game.SkillIdQiHuoKeJu
}

func (m *QiHuoKeJu) Execute(g *game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
	fsm, ok := g.GetFsm().(*game.ReceivePhaseReceiverSkill)
	if !ok || fsm.InFrontOfWhom.FindSkill(m.GetSkillId()) == nil {
		return nil, false, false
	}
	if fsm.InFrontOfWhom.GetSkillUseCount(m.GetSkillId()) > 0 {
		return nil, false, false
	}
	if len(fsm.MessageCard.GetColors()) < 2 {
		return nil, false, false
	}
	fsm.InFrontOfWhom.AddSkillUseCount(m.GetSkillId())
	return &executeQiHuoKeJu{fsm: fsm}, true, ok
}

func (m *QiHuoKeJu) ExecuteProtocol(*game.Game, game.IPlayer, proto.Message) {
}

type executeQiHuoKeJu struct {
	fsm *game.ReceivePhaseReceiverSkill
}

func (e *executeQiHuoKeJu) Resolve() (next game.Fsm, continueResolve bool) {
	for _, p := range e.fsm.WhoseTurn.GetGame().GetPlayers() {
		p.NotifyReceivePhaseWithWaiting(e.fsm.WhoseTurn, e.fsm.InFrontOfWhom, e.fsm.MessageCard, e.fsm.InFrontOfWhom, 20)
	}
	return e, false
}

func (e *executeQiHuoKeJu) ResolveProtocol(player game.IPlayer, message proto.Message) (next game.Fsm, continueResolve bool) {
	if player.Location() != e.fsm.InFrontOfWhom.Location() {
		logger.Error("不是你发技能的时机")
		return e, false
	}
	if _, ok := message.(*protos.EndReceivePhaseTos); ok {
		if player.Location() == e.fsm.InFrontOfWhom.Location() {
			player.IncrSeq()
			return e.fsm, true
		} else {
			logger.Error("还没轮到你")
			return e, false
		}
	}
	pb, ok := message.(*protos.SkillQiHuoKeJuTos)
	if !ok {
		logger.Error("错误的协议")
		return e, false
	}
	r := e.fsm.WhoseTurn
	g := r.GetGame()
	if humanPlayer, ok := r.(*game.HumanPlayer); ok && pb.Seq != humanPlayer.Seq {
		logger.Error("操作太晚了, required Seq: ", humanPlayer.Seq, ", actual Seq: ", pb.Seq)
		return e, false
	}
	card := r.FindMessageCard(pb.CardId)
	if card == nil {
		logger.Error("没有这张卡")
		return e, false
	}
	r.IncrSeq()
	logger.Info("[毛不拔]发动了[奇货可居]")
	r.DeleteMessageCard(card.GetId())
	r.AddCards(card)
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.SkillQiHuoKeJuToc{
				PlayerId: player.GetAlternativeLocation(r.Location()),
				CardId:   card.GetId(),
			})
		}
	}
	return e.fsm, true
}
