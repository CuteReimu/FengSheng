package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
)

type HumanPlayer struct {
	interfaces.BasePlayer
	cellnet.Session
}

func (r *HumanPlayer) Init(game interfaces.IGame, location int) {
	r.BasePlayer.Init(game, location)
	msg := &protos.InitToc{
		PlayerCount: uint32(len(r.GetGame().GetPlayers())),
		Identity:    protos.Color_Red,
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyAddHandCard(cards ...interfaces.ICard) {
	msg := &protos.AddCardToc{}
	for _, card := range cards {
		msg.Cards = append(msg.Cards, card.ToPbCard())
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyOtherAddHandCard(location int, count int) {
	msg := &protos.AddCardToc{
		PlayerId:         r.getAlternativeLocation(location),
		UnknownCardCount: uint32(count),
	}
	r.Send(msg)
}

func (r *HumanPlayer) getAlternativeLocation(location int) uint32 {
	if location == 99999 {
		return 99999
	}
	location -= r.Location()
	totalPlayerCount := len(r.GetGame().GetPlayers())
	if location < 0 {
		location += totalPlayerCount
	}
	return uint32(location % totalPlayerCount)
}
