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
	Init()
	GetSkillId() SkillId
	Execute(g *Game, r IPlayer, message proto.Message)
}

type SkillId uint32
