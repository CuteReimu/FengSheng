package game

import (
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
)

type HumanPlayer struct {
	basePlayer
	cellnet.Session
}

func (r *HumanPlayer) Init(game *Game, location int) {
	r.basePlayer.Init(game, location)
	msg := &protos.InitToc{
		PlayerCount: uint32(r.game.TotalPlayerCount),
		Identity:    protos.Color_Red,
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyAddHandCard(cards ...*protos.Card) {
	msg := &protos.AddCardToc{Cards: cards}
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
	if location < 0 {
		location += r.game.TotalPlayerCount
	}
	return uint32(location % r.game.TotalPlayerCount)
}
