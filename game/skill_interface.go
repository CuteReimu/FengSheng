package game

import (
	"github.com/CuteReimu/FengSheng/protos"
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
}

type SkillId uint32

const (
	SkillIdLianLuo SkillId = iota
	SkillIdMingEr
)
