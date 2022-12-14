package role

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"google.golang.org/protobuf/proto"
)

func init() {
	game.RoleCache = append(game.RoleCache, &game.RoleSkillsData{
		Name:   "金生火",
		Role:   protos.Role_jin_sheng_huo,
		FaceUp: true,
		Skills: []game.ISkill{&JinShen{}},
	})
}

type JinShen struct {
}

func (m *JinShen) Init(g *game.Game) {
	g.AddListeningSkill(m)
}

func (m *JinShen) GetSkillId() game.SkillId {
	return game.SkillIdJinShen
}

func (m *JinShen) Execute(g *game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
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
	return &executeJinShen{fsm: fsm}, true, ok
}

func (m *JinShen) ExecuteProtocol(*game.Game, game.IPlayer, proto.Message) {
}

type executeJinShen struct {
	fsm *game.ReceivePhaseReceiverSkill
}

func (e *executeJinShen) Resolve() (next game.Fsm, continueResolve bool) {
	for _, p := range e.fsm.WhoseTurn.GetGame().GetPlayers() {
		p.NotifyReceivePhaseWithWaiting(e.fsm.WhoseTurn, e.fsm.InFrontOfWhom, e.fsm.MessageCard, e.fsm.InFrontOfWhom, 20)
	}
	return e, false
}

func (e *executeJinShen) ResolveProtocol(player game.IPlayer, message proto.Message) (next game.Fsm, continueResolve bool) {
	if player.Location() != e.fsm.InFrontOfWhom.Location() {
		logger.Error("不是你发技能的时机")
		return e, false
	}
	if pb, ok := message.(*protos.EndReceivePhaseTos); ok {
		if r, ok := player.(*game.HumanPlayer); ok && r.Seq != pb.Seq {
			logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
			return e, false
		}
		if player.Location() == e.fsm.InFrontOfWhom.Location() {
			player.IncrSeq()
			return e.fsm, true
		} else {
			logger.Error("还没轮到你")
			return e, false
		}
	}
	pb, ok := message.(*protos.SkillJinShenTos)
	if !ok {
		logger.Error("错误的协议")
		return e, false
	}
	r := e.fsm.InFrontOfWhom
	g := r.GetGame()
	if humanPlayer, ok := r.(*game.HumanPlayer); ok && pb.Seq != humanPlayer.Seq {
		logger.Error("操作太晚了, required Seq: ", humanPlayer.Seq, ", actual Seq: ", pb.Seq)
		return e, false
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		logger.Error("没有这张卡")
		return e, false
	}
	r.IncrSeq()
	logger.Info("[金生火]发动了[谨慎]")
	messageCard := e.fsm.MessageCard
	r.DeleteCard(card.GetId())
	r.DeleteMessageCard(messageCard.GetId())
	r.AddMessageCards(card)
	r.AddCards(messageCard)
	for _, p := range g.GetPlayers() {
		if player, ok := p.(*game.HumanPlayer); ok {
			player.Send(&protos.SkillJinShenToc{
				PlayerId: player.GetAlternativeLocation(r.Location()),
				Card:     card.ToPbCard(),
			})
		}
	}
	return e.fsm, true
}
