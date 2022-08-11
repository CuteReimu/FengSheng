package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
	"time"
)

type HumanPlayer struct {
	interfaces.BasePlayer
	cellnet.Session
	seq   uint32
	timer *time.Timer
}

func (r *HumanPlayer) Init(game interfaces.IGame, location int) {
	r.BasePlayer.Init(game, location)
	msg := &protos.InitToc{
		PlayerCount: uint32(len(r.GetGame().GetPlayers())),
		Identity:    protos.Color_Red,
	}
	r.Send(msg)
	r.seq++
}

func (r *HumanPlayer) NotifyAddHandCard(location int, unknownCount int, cards ...interfaces.ICard) {
	msg := &protos.AddCardToc{
		PlayerId:         r.GetAlternativeLocation(location),
		UnknownCardCount: uint32(unknownCount),
	}
	for _, card := range cards {
		msg.Cards = append(msg.Cards, card.ToPbCard())
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyDrawPhase(location int) {
	playerId := r.GetAlternativeLocation(location)
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Draw_Phase,
		WaitingPlayerId: playerId,
	}
	r.Send(msg)
}

func (r *HumanPlayer) NotifyMainPhase(location int, waitSecond uint32) {
	playerId := r.GetAlternativeLocation(location)
	msg := &protos.NotifyPhaseToc{
		CurrentPlayerId: playerId,
		CurrentPhase:    protos.Phase_Main_Phase,
		WaitingPlayerId: playerId,
		WaitingSecond:   waitSecond,
	}
	r.Send(msg)
	seq := r.seq
	time.AfterFunc(time.Second*time.Duration(waitSecond), func() {
		r.GetGame().Post(func() {
			if seq == r.seq {
				r.GetGame().Post(r.GetGame().SendPhase)
				r.seq++
			}
		})
	})
}
