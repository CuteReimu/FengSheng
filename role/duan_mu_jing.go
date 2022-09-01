package role

import (
	"github.com/CuteReimu/FengSheng/game"
	"github.com/CuteReimu/FengSheng/protos"
	"google.golang.org/protobuf/proto"
)

func init() {
	game.RoleCache = append(game.RoleCache, &game.RoleSkillsData{
		Name:   "端木静",
		Role:   protos.Role_duan_mu_jing,
		FaceUp: true,
		Skills: []game.ISkill{&XinSiChao{}},
	})
}

type XinSiChao struct {
}

func (x *XinSiChao) Init(*game.Game) {
}

func (x *XinSiChao) GetSkillId() game.SkillId {
	return game.SkillIdXinSiChao
}

func (x *XinSiChao) Execute(*game.Game) (nextFsm game.Fsm, continueResolve bool, ok bool) {
	return
}

func (x *XinSiChao) ExecuteProtocol(g *game.Game, r game.IPlayer, message proto.Message) {
	fsm, ok := g.GetFsm().(*game.MainPhaseIdle)
	if !ok || r.Location() != fsm.Player.Location() {
		logger.Error("现在不是出牌阶段空闲时点")
		return
	}
	if r.GetSkillUseCount(x.GetSkillId()) > 0 {
		logger.Error("[新思潮]一回合只能发动一次")
		return
	}
	pb := message.(*protos.SkillXinSiChaoTos)
	humanPlayer, ok := r.(*game.HumanPlayer)
	if ok && pb.Seq != humanPlayer.Seq {
		logger.Error("操作太晚了, required Seq: ", humanPlayer.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		logger.Error("没有这张卡")
		return
	}
	humanPlayer.Seq++
	if humanPlayer.Timer != nil {
		humanPlayer.Timer.Stop()
	}
	r.AddSkillUseCount(x.GetSkillId())
	logger.Info("[端木静]发动了[新思潮]")
	g.PlayerDiscardCard(r, card)
	r.Draw(2)
	g.ContinueResolve()
}
