package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"google.golang.org/protobuf/proto"
)

type RoleSkillsData struct {
	Name   string
	Role   protos.Role
	FaceUp bool
	Skills []ISkill
}

var RoleCache []*RoleSkillsData

type ISkill interface {
	Init(g *Game)
	GetSkillId() SkillId
	Execute(g *Game) (nextFsm Fsm, continueResolve bool, ok bool)
	ExecuteProtocol(g *Game, r IPlayer, message proto.Message)
}

type SkillId uint32

const (
	SkillIdLianLuo SkillId = iota
	SkillIdMingEr
	SkillIdXinSiChao
	SkillIdMianLiCangZhen
	SkillIdQiHuoKeJu
	SkillIdJinShen
)
