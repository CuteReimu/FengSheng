package game

import (
	"github.com/CuteReimu/FengSheng/game/interfaces"
	"github.com/CuteReimu/FengSheng/protos"
	"github.com/davyxu/cellnet"
	"github.com/sirupsen/logrus"
	"time"
)

type HumanPlayer struct {
	interfaces.BasePlayer
	cellnet.Session
	Seq    uint32
	timer  *time.Timer
	logger logrus.FieldLogger
}

func (r *HumanPlayer) Init(game interfaces.IGame, location int) {
	r.logger = logrus.WithField("human_player", r.Location())
	r.BasePlayer.Init(game, location)
	msg := &protos.InitToc{
		PlayerCount: uint32(len(r.GetGame().GetPlayers())),
		Identity:    protos.Color_Red,
	}
	r.Send(msg)
	r.Seq++
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
	if r.Location() == r.GetGame().GetWhoseTurn() {
		msg.Seq = r.Seq
	}
	r.Send(msg)
	if r.Location() == r.GetGame().GetWhoseTurn() {
		seq := r.Seq
		time.AfterFunc(time.Second*time.Duration(waitSecond), func() {
			r.GetGame().Post(func() {
				if seq == r.Seq {
					r.GetGame().Post(r.GetGame().SendPhase)
					r.Seq++
				}
			})
		})
	}
}

func (r *HumanPlayer) onUseShiTan(pb *protos.UseShiTanTos) {
	if pb.Seq != r.Seq {
		r.logger.Error("操作太晚了, required Seq: ", r.Seq, ", actual Seq: ", pb.Seq)
		return
	}
	card := r.FindCard(pb.CardId)
	if card == nil {
		r.logger.Error("没有这张牌")
		return
	}
	if card.GetType() != protos.CardType_Shi_Tan {
		r.logger.Error("这张牌不是试探，而是", card)
		return
	}
	if pb.PlayerId >= uint32(len(r.GetGame().GetPlayers())) {
		r.logger.Error("目标错误: ", pb.PlayerId)
	}
	target := r.GetGame().GetPlayers()[r.GetAbstractLocation(int(pb.PlayerId))]
	if card.CanUse(r.GetGame(), r, target) {
		r.Seq++
		if r.timer != nil {
			r.timer.Stop()
		}
		card.Execute(r.GetGame(), r, target)
	}
}
